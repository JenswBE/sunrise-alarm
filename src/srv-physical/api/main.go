package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/config"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/handler"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/gpiobutton"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/mockbutton"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/mockleds"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/p9813leds"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/pahomqtt"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/leds"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/mqtt"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/utils/buttonpoller"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	// Setup logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	gin.SetMode(gin.ReleaseMode)

	// Parse config
	apiConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	// Setup Debug logging if enabled
	if apiConfig.Server.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.Debug().Msg("Debug logging enabled")
	}

	// Setup devices
	var devButton repositories.Button
	var devLeds repositories.Leds
	buttonChannel := make(chan buttonpoller.ButtonPress)
	if !apiConfig.Server.Mocked {
		// Init real devices
		if err := rpio.Open(); err != nil {
			log.Fatal().Err(err).Msg("RPIO: Failed to initialize GPIO library")
		}
		defer rpio.Close()
		devButton, err = gpiobutton.NewGPIOButton(apiConfig.Button.GPIONum, true)
		if err != nil {
			log.Fatal().Err(err).Msg("Button: Failed to initialize GPIO button")
		}
		p9813Leds, err := p9813leds.NewP9813Leds()
		if err != nil {
			log.Fatal().Err(err).Msg("LED: Failed to initialize P9813 led driver on SPI0")
		}
		defer p9813Leds.Close()
		devLeds = p9813Leds
	} else {
		// Init mocked devices
		devButton = &mockbutton.MockButton{}
		devLeds = &mockleds.MockLeds{}
	}

	// Setup repositories
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mqttClient, err := pahomqtt.NewPahoMQTT(ctx, apiConfig.MQTT.BrokerHost, apiConfig.MQTT.BrokerPort)
	if err != nil {
		log.Fatal().Err(err).Msg("MQTT: Creating client returned error")
	}

	// Setup services
	ledsService := leds.NewService(devLeds)
	mqttService := mqtt.NewService(mqttClient)

	// Start polling button
	buttonpoller.NewButtonPoller(devButton, buttonChannel)
	go func() {
		for {
			var err error
			switch <-buttonChannel {
			case buttonpoller.ButtonPressShort:
				err = mqttService.PublishButtonPressed(ctx)
			case buttonpoller.ButtonPressLong:
				err = mqttService.PublishButtonLongPressed(ctx)
			}
			if err != nil {
				log.Error().Err(err).Msg("Failed to publish button push to MQTT")
			}
		}
	}()

	// Setup Gin
	router := gin.Default()
	err = router.SetTrustedProxies(apiConfig.Server.TrustedProxies)
	if err != nil {
		log.Fatal().Err(err).Strs("trusted_proxies", apiConfig.Server.TrustedProxies).Msg("Failed to set trusted proxies")
	}
	router.StaticFile("/", "../docs/index.html")
	router.StaticFile("/index.html", "../docs/index.html")
	router.StaticFile("/oauth2-redirect.html", "../docs/oauth2-redirect.html")
	router.StaticFile("/openapi.yml", "../docs/openapi.yml")

	// Setup handlers
	backlightHandler := handler.NewBacklightHandler()
	mockHandler := handler.NewMockHandler(mqttService)
	ledsHandler := handler.NewLedsHandler(ledsService)

	// Public routes
	public := router.Group("/")
	backlightHandler.RegisterRoutes(public)
	ledsHandler.RegisterRoutes(public)
	mockHandler.RegisterRoutes(public)

	// Start Gin
	port := strconv.Itoa(apiConfig.Server.Port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal().Err(err).Int("port", apiConfig.Server.Port).Msg("Failed to start Gin server")
	}
}
