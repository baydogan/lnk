package config

import (
	"os"

	"github.com/baydogan/lnk/domain"
)

func Load() (domain.ServerConfig, error) {
	cfg, _, err := ReadServerConfig()
	if err != nil {
		return cfg, err
	}

	if v := os.Getenv("MONGO_URI"); v != "" {
		cfg.MongoURI = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		cfg.RedisAddr = v
	}
	if v := os.Getenv("BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("MODE"); v != "" {
		cfg.Mode = v
	}

	if cfg.MongoURI == "" {
		cfg.MongoURI = "mongodb://lnk:lnk@localhost:27017/lnk?authSource=admin"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	return cfg, nil
}
