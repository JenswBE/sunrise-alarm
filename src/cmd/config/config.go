package config

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Debug     bool
	LogFormat LogFormat
	Physical  PhysicalConfig
}

type LogFormat string

const (
	LogFormatConsole LogFormat = "CONSOLE"
	LogFormatJSON    LogFormat = "JSON"
)

type PhysicalConfig struct {
	Button struct {
		GPIONum int
	}
	Buzzer struct {
		GPIONum int
	}
	Leds struct {
		SunriseDuration time.Duration
	}
	LightSensor struct {
		// Path to the I2C device, e.g. /dev/i2c-1
		I2CDevice string
	}
	Mocked bool
}

func ParseConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("Debug", false)
	viper.SetDefault("LogFormat", LogFormatJSON)
	viper.SetDefault("Physical.Button.GPIONum", 23) // GPIO23 on pin 16
	viper.SetDefault("Physical.Buzzer.GPIONum", 18) // GPIO18 on pin 12
	viper.SetDefault("Physical.Leds.SunriseDuration", 5*time.Minute)
	viper.SetDefault("Physical.LightSensor.I2CDevice", "/dev/i2c-1")
	viper.SetDefault("Physical.Mocked", false)

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
		{"Debug", "DEBUG"},
		{"LogFormat", "LOG_FORMAT"},
		{"Physical.Button.GPIONum", "PHYSICAL_BUTTON_GPIO_PIN"},
		{"Physical.Buzzer.GPIONum", "PHYSICAL_BUZZER_GPIO_PIN"},
		{"Physical.Leds.SunriseDuration", "PHYSICAL_LEDS_SUNRISE_DURATION"},
		{"Physical.LightSensor.I2CDevice", "PHYSICAL_LIGHT_SENSOR_I2C_DEVICE"},
		{"Physical.Mocked", "PHYSICAL_MOCKED"},
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
