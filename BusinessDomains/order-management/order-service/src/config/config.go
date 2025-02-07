package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	AppPort    string `mapstructure:"APP_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // Permite que Viper lea variables de entorno

	if err = viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
		log.Printf("Unable to decode into struct, %v", err)
		return
	}

	return config, nil
}
