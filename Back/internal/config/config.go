package config

import "os"

type Config struct {
	Port   string
	DBAddr string
}

func NewConfig() *Config {
	return &Config{
		Port:   os.Getenv("PORT"),
		DBAddr: os.Getenv("DB_ADDR"),
	}
}
