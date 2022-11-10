package main

import (
	"os"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/cmd/config"
	"github.com/JenswBE/sunrise-alarm/src/services/alarm"
	"github.com/JenswBE/sunrise-alarm/src/services/audio"
	"github.com/JenswBE/sunrise-alarm/src/services/gui"
	"github.com/JenswBE/sunrise-alarm/src/services/physical"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logging
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.Output(output)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Parse config
	svcConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}

	// Setup Debug logging if enabled
	if svcConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Debug logging enabled")
	}

	// Setup event bus
	eventBus := &pubsub.EventBus{}

	// Setup services
	audioSvc, err := audio.NewAudioService()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init audio service")
	}
	physicalSvc := physical.NewPhysicalService(svcConfig.Physical, eventBus)
	alarmSvc := alarm.NewAlarmService(physicalSvc, audioSvc, eventBus, 10*time.Minute, 5*time.Minute) // TODO: Put in config?
	defer func() {
		if err := alarmSvc.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to cleanly close alarm service")
		}
	}()

	// Start GUI
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.RedirectTrailingSlash = true
	guiHandler := gui.NewHandler(eventBus, alarmSvc, svcConfig.Debug)
	guiHandler.RegisterRoutes(router)
	router.HTMLRender = guiHandler.NewRenderer()
	err = router.Run(":8123")
	if err != nil {
		log.Fatal().Err(err).Int("port", 8123).Msg("Failed to start GUI")
	}
}
