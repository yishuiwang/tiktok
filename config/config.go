package config

import (
	"github.com/spf13/viper"
	"tiktok/logger"
)

func GetConfig(key string) string {
	Viper := viper.New()
	Viper.SetConfigName("app")
	Viper.AddConfigPath("config")
	err := Viper.ReadInConfig()
	if err != nil {
		logger.ZapLogger.Error("read config failed", logger.Error(err))
	}
	return Viper.GetString(key)
}
