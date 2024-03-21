package dsn

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
)

func GetDSN() (string, error) {
	err := config.GetConfig()
	if err != nil {
		return "", fmt.Errorf("[config.GetConfig]: Can't read config: %w", err)
	}

	url := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", viper.GetString("DATABASE_USERNAME"), viper.GetString("DATABASE_PASSWORD"), viper.GetString("DATABASE_PORT"), viper.GetString("DATABASE_NAME"))
	return url, nil
}
