package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	DataBase DataBaseConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Port         int
	Env          string
	JWTSecretKey string
}

type DataBaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMODE  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Setup() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.env", "development")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
