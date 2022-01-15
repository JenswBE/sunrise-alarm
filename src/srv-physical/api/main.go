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
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/pahomqtt"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/mqtt"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/utils/buttonpoller"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// Services
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mqttClient, err := pahomqtt.NewPahoMQTT(ctx, apiConfig.MQTT.BrokerHost, apiConfig.MQTT.BrokerPort)
	if err != nil {
		log.Fatal().Err(err).Msg("MQTT: Creating client returned error")
	}
	mqttService := mqtt.NewService(mqttClient)

	// Setup devices
	var button repositories.Button
	buttonChannel := make(chan buttonpoller.ButtonPress)
	if !apiConfig.Server.Mocked {
		// Init real devices
		button, err = gpiobutton.NewGPIOButton(apiConfig.Button.GPIONum, true)
		if err != nil {
			log.Fatal().Err(err).Msg("Button: Failed to initialize GPIO button")
		}
	} else {
		// Init mocked devices
		button = &mockbutton.MockButton{}
	}
	buttonpoller.NewButtonPoller(button, buttonChannel)
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
	mockHandler := handler.NewMockHandler(mqttService)

	// Public routes
	public := router.Group("/")
	mockHandler.RegisterRoutes(public)

	// Start Gin
	port := strconv.Itoa(apiConfig.Server.Port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal().Err(err).Int("port", apiConfig.Server.Port).Msg("Failed to start Gin server")
	}
}
