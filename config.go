package minrpc

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config 是一个配置结构体
type Config struct {
	ServerAddress string `mapstructure:"server_address"`
	MaxConcurrent int    `mapstructure:"max_concurrent"`
	LogLevel      string `mapstructure:"log_level"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/min-rpc/")

	viper.SetEnvPrefix("MIN_RPC")
	viper.AutomaticEnv()

	pflag.String("server_address", "localhost:8080", "Server address")
	pflag.Int("max_concurrent", 100, "Max concurrent requests")
	pflag.String("log_level", "info", "Log level")

	pflag.Parse()

	viper.BindPFlags(pflag.CommandLine)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
