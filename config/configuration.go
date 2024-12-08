package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HttpServer struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c *HttpServer) String() string {
	return fmt.Sprintf("Port - %d, readTimeout: %v, WriteTimeout: %v, Idletimeout: %v",
		c.Port, c.ReadTimeout, c.WriteTimeout, c.IdleTimeout)
}

type BigCacheConfig struct {
	TTL              int  `mapstructure:"TTL_SECS"`
	Flag             bool `mapstructure:"FLAG"`
	Shards           int  `mapstructure:"SHARDS"`
	MaxEntrySize     int  `mapstructure:"MAX_ENTRY_SIZE"`
	StatsEnabled     bool `mapstructure:"STATS_ENABLED"`
	HardMaxCacheSize int  `mapstructure:"HARD_MAX_CACHE_SIZE"`
}

type Cors struct {
	Origins []string `mapstructure:"ORIGINS"`
	Methods []string `mapstructure:"METHODS"`
	Headers []string `mapstructure:"HEADERS"`
}

type Log struct {
	Level                      string   `mapstructure:"LEVEL"`
	Env                        string   `mapstructure:"ENV"`
	GoCommonLoggingEnable      bool     `mapstructure:"GO_COMMON_LOGGING_ENABLE"`
	CtxKeys                    []string `mapstructure:"CTX_KEYS"`
	IsDebugLevelLoggingEnabled bool     `mapstructure:"DEBUG_LOGGING_ENABLED"`
}

type Configurations struct {
	Environment string
	AppName     string         `mapstructure:"APP_NAME"`
	LogConfig   Log            `mapstructure:"LOG"`
	BigCache    BigCacheConfig `mapstructure:"BIG_CACHE"`
	HttpServer  *HttpServer    `json:"httpServer"`
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
