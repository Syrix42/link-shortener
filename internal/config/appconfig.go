package config

// Impelimentaion of App Configuration Resides Here

import (
	"fmt"
	"os"
)

type AppConfig struct {
	Port string
}

func LoadAppConfig() AppConfig {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return AppConfig{Port: port}
}

func (c AppConfig) ListenAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}
