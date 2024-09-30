package util

import (
	"time"

	"github.com/spf13/viper"
)

// store all configurations
// read by viper
type Config struct {
	DBDriver          string        `mapstructure:"DB_DRIVER"`
	DBSourse          string        `mapstructure:"DB_SOURSE"`
	ServerAddress     string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DURATION"`
}

// load config from file or environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv() //override with environmental var if they exists

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
