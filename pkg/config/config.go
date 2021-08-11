package config

import "github.com/spf13/viper"

// InitConfig - load config from config.yml
func InitConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
