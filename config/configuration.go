package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Log struct {
	Level                      string   `mapstructure:"LEVEL"`
	Env                        string   `mapstructure:"ENV"`
	GoCommonLoggingEnable      bool     `mapstructure:"GO_COMMON_LOGGING_ENABLE"`
	CtxKeys                    []string `mapstructure:"CTX_KEYS"`
	IsDebugLevelLoggingEnabled bool     `mapstructure:"DEBUG_LOGGING_ENABLED"`
}

type Configurations struct {
	Environment string
	AppName     string `mapstructure:"APP_NAME"`
	LogConfig   Log    `mapstructure:"LOG"`
}

var Configuration Configurations

func InitConfigurations() Configurations {
	fmt.Println("Initializing configurations...")
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %w", err))
	}
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %w", err))
	}

	if err = viper.Unmarshal(&Configuration); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return Configuration
}
