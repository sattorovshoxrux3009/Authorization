package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port  string
	Mysql Mysql
}
type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Load(path string) Config {
	// .env faylni o'qish
	err := godotenv.Load(path + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // Xatolikni qaytarish
	}

	conf := viper.New()
	conf.AutomaticEnv()
	cfg := Config{
		Port: conf.GetString("PORT"),
		Mysql: Mysql{
			Host:     conf.GetString("MYSQL_HOST"),
			Port:     conf.GetString("MYSQL_PORT"),
			User:     conf.GetString("MYSQL_USER"),
			Password: conf.GetString("MYSQL_PASSWORD"),
			Database: conf.GetString("MYSQL_DATABASE"),
		},
	}

	return cfg
}
