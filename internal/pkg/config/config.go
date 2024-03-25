package config

import "github.com/spf13/viper"

// GetConfig
// Получает значения пременных, занесенных в конфиг
func GetConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
