package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Debug     bool
	LogFormat LogFormat
	Alarms    AlarmsConfig
	Physical  PhysicalConfig
}

type LogFormat string

const (
	LogFormatConsole LogFormat = "CONSOLE"
	LogFormatJSON    LogFormat = "JSON"
)

type AlarmsConfig struct {
	// Number of minutes the light will start before alarm time.
	LightDurationInMinutes uint
	// Number of minutes the sound will start before alarm time.
	SoundDurationInMinutes uint
}

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
	Mocked bool
}

func ParseConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("Debug", false)
	viper.SetDefault("LogFormat", LogFormatJSON)
	viper.SetDefault("Alarms.LightDurationInMinutes", 10) // Light will start 10 minutes before alarm time
	viper.SetDefault("Alarms.SoundDurationInMinutes", 7)  // Sound will start 10 minutes before alarm time
	viper.SetDefault("Physical.Button.GPIONum", 23)       // GPIO23 on pin 16
	viper.SetDefault("Physical.Buzzer.GPIONum", 18)       // GPIO18 on pin 12
	viper.SetDefault("Physical.Leds.SunriseDuration", 5*time.Minute)
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
		{"Alarms.LightDurationInMinutes", "ALARMS_LIGHT_DURATION_IN_MINUTES"},
		{"Alarms.SoundDurationInMinutes", "ALARMS_SOUND_DURATION_IN_MINUTES"},
		{"Physical.Button.GPIONum", "PHYSICAL_BUTTON_GPIO_PIN"},
		{"Physical.Buzzer.GPIONum", "PHYSICAL_BUZZER_GPIO_PIN"},
		{"Physical.Leds.SunriseDuration", "PHYSICAL_LEDS_SUNRISE_DURATION"},
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
