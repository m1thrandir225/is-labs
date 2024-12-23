package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	JWTKey               string        `mapstructure:"JWT_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	SMTPHost             string        `mapstructure:"SMTP_HOST"`
	SMTPPort             int           `mapstructure:"SMTP_PORT"`
	SMTPUsername         string        `mapstructure:"SMTP_USERNAME"`
	SMTPPassword         string        `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
