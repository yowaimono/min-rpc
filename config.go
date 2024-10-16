package minrpc

import (
    "flag"
    "fmt"
    "os"

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

    flag.String("server_address", "localhost:8080", "Server address")
    flag.Int("max_concurrent", 100, "Max concurrent requests")
    flag.String("log_level", "info", "Log level")

    flag.Parse()

    viper.BindPFlags(flag.CommandLine)

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