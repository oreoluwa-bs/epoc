package web

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	PORT        uint16
}

func LoadConfig() Config {
	cfg := Config{
		DatabaseURL: "data/data.db",
		PORT:        8000,
	}

	if dbUrl, exists := os.LookupEnv("DATABASE_URL"); exists {
		cfg.DatabaseURL = dbUrl
	}

	if port, exists := os.LookupEnv("PORT"); exists {

		if p, err := strconv.Atoi(port); err == nil {
			cfg.PORT = uint16(p)
		}
	}

	return cfg
}
