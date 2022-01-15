package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Button struct {
		GPIONum int
	}
	MQTT struct {
		BrokerHost string
		BrokerPort int
	}
	Server struct {
		Debug          bool
		Mocked         bool
		Port           int
		TrustedProxies []string
	}
}

func ParseConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("Button.GPIONum", 23) // GPIO23 on pin 16
	viper.SetDefault("MQTT.BrokerHost", "localhost")
	viper.SetDefault("MQTT.BrokerPort", 1883)
	viper.SetDefault("Server.Debug", false)
	viper.SetDefault("Server.Mocked", false)
	viper.SetDefault("Server.Port", 8080)
	viper.SetDefault("Server.TrustedProxies", []string{"172.16.0.0/16"}) // Default Docker IP range

	// Parse config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed reading config file: %w", err)
		}
		log.Warn().Err(err).Msg("No config file found, expecting configuration through ENV variables")
	}

	// Bind ENV variables
	err = bindEnvs([]envBinding{
		{"Button.GPIONum", "BUTTON_GPIO_PIN"},
		{"MQTT.BrokerHost", "MQTT_BROKER_HOST"},
		{"MQTT.BrokerPort", "MQTT_BROKER_PORT"},
		{"Server.Debug", "SRV_PHYSICAL_DEBUG"},
		{"Server.Mocked", "SRV_PHYSICAL_MOCKED"},
		{"Server.Port", "SRV_PHYSICAL_PORT"},
	})
	if err != nil {
		return nil, err
	}

	// Unmarshal to Config struct
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	return &config, nil
}

type envBinding struct {
	configPath string
	envName    string
}

func bindEnvs(bindings []envBinding) error {
	for _, binding := range bindings {
		err := viper.BindEnv(binding.configPath, binding.envName)
		if err != nil {
			return fmt.Errorf("failed to bind env var %s to %s: %w", binding.envName, binding.configPath, err)
		}
	}
	return nil
}
