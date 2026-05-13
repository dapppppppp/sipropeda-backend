package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Env  string `mapstructure:"SERVER.ENV"`
		Port string `mapstructure:"SERVER.PORT"`
	}
	DB struct {
		PostgreSQL struct {
			Read struct {
				Host     string `mapstructure:"DB_HOST"`
				Port     string `mapstructure:"DB_PORT"`
				Name     string `mapstructure:"DB_NAME"`
				Username string `mapstructure:"DB_USER"`
				Password string `mapstructure:"DB_PASSWORD"`
				Timezone string `mapstructure:"DB_TIMEZONE"`
			}
			Write struct {
				Host     string `mapstructure:"DB_HOST"`
				Port     string `mapstructure:"DB_PORT"`
				Name     string `mapstructure:"DB_NAME"`
				Username string `mapstructure:"DB_USER"`
				Password string `mapstructure:"DB_PASSWORD"`
				Timezone string `mapstructure:"DB_TIMEZONE"`
			}
		}
	}
}

var config *Config

// Get will return the config instance
func Get() *Config {
	if config == nil {
		viper.SetConfigFile(".env")
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("Warning: .env file not found or error reading it.")
		}

		config = &Config{}
		err = viper.Unmarshal(&config.Server)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
		
		err = viper.Unmarshal(&config.DB.PostgreSQL.Read)
		err = viper.Unmarshal(&config.DB.PostgreSQL.Write)
		if err != nil {
			log.Fatalf("unable to decode database config into struct, %v", err)
		}
	}
	return config
}