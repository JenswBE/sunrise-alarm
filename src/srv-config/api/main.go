package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/config"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/handler"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/repositories/pahomqtt"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/usecases/mqtt"
	"github.com/dgraph-io/badger"
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

	// Setup DB
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open BadgerDB")
	}
	defer db.Close()

	// Setup repositories
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mqttClient, err := pahomqtt.NewPahoMQTT(ctx, apiConfig.MQTT.BrokerHost, apiConfig.MQTT.BrokerPort)
	if err != nil {
		log.Fatal().Err(err).Msg("MQTT: Creating client returned error")
	}

	// Setup services
	mqttService := mqtt.NewService(mqttClient)

	// Setup Gin
	router := gin.Default()
	err = router.SetTrustedProxies(apiConfig.Server.TrustedProxies)
	if err != nil {
		log.Fatal().Err(err).Strs("trusted_proxies", apiConfig.Server.TrustedProxies).Msg("Failed to set trusted proxies")
	}
	router.StaticFile("/", "../docs/index.html")
	router.StaticFile("/index.html", "../docs/index.html")
	router.StaticFile("/openapi.yml", "../docs/openapi.yml")

	// Setup handlers
	alarmsHandler := handler.NewAlarmsHandler(mqttService)

	// Register routes
	root := router.Group("/")
	alarmsHandler.RegisterRoutes(root)

	// Start Gin
	port := strconv.Itoa(apiConfig.Server.Port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatal().Err(err).Int("port", apiConfig.Server.Port).Msg("Failed to start Gin server")
	}
}
