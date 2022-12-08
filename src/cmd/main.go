package main

import (
	"io"
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
	"github.com/samber/lo"
)

func main() {
	// Setup logging
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Parse config
	svcConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Main: Failed to parse config")
	}

	// Setup log format
	ginLogger := gin.Logger()
	switch svcConfig.LogFormat {
	case config.LogFormatConsole:
		// Gin already defaults to Console
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = log.Output(output)
	case config.LogFormatJSON:
		// Zerolog already defaults to JSON
		gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, handlerCount int) {
			log.Debug().Str("method", httpMethod).Str("path", absolutePath).Str("handler", handlerName).Int("handler_count", handlerCount).Msg("Registered new gin handler")
		}
		ginLogger = gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				log.Debug().
					Time("timestamp", params.TimeStamp).
					Int("status", params.StatusCode).
					Stringer("latency", params.Latency).
					Str("client_ip", params.ClientIP).
					Msgf("%s %s", params.Method, params.Path)
				return "" // Outputs will be discarded anyway
			},
			Output: io.Discard,
		})
	}

	// Setup Trace/Debug logging if enabled
	switch {
	case svcConfig.Trace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Trace().Msg("Main: Trace logging enabled")
	case svcConfig.Debug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("Main: Debug logging enabled")
	}

	// Setup event bus
	eventBus := &pubsub.EventBus{}

	// Setup services
	audioSvc, err := audio.NewAudioService()
	if err != nil {
		log.Fatal().Err(err).Msg("Main: Failed to init audio service")
	}
	physicalSvc := physical.NewPhysicalService(svcConfig.Physical, eventBus)
	defer physicalSvc.Close()
	alarmSvc, err := alarm.NewAlarmService(physicalSvc, audioSvc, eventBus)
	if err != nil {
		log.Fatal().Err(err).Msg("Main: Failed to create the alarm service")
	}
	defer func() {
		if err := alarmSvc.Close(); err != nil {
			log.Fatal().Err(err).Msg("Main: Failed to cleanly close alarm service")
		}
	}()

	// Start GUI
	router := gin.New()
	router.Use(ginLogger, gin.Recovery())
	lo.Must0(router.SetTrustedProxies(nil)) // nil can never return a parsing error
	router.RedirectTrailingSlash = true
	guiHandler := gui.NewHandler(eventBus, alarmSvc, svcConfig.Debug)
	guiHandler.RegisterRoutes(router)
	router.HTMLRender = guiHandler.NewRenderer()
	err = router.Run(":8123")
	if err != nil {
		log.Fatal().Err(err).Int("port", 8123).Msg("Main: Failed to start GUI")
	}
}
