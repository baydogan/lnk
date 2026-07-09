package server

import (
	"path/filepath"
	"testing"

	"github.com/baydogan/lnk/internal/config"
	"github.com/baydogan/lnk/internal/models"
)

func writeConfigFile(t *testing.T, cfg *models.ServerConfig) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "server.yaml")
	t.Setenv("LNK_SERVER_CONFIG", path)
	if _, err := config.WriteServerConfig(cfg); err != nil {
		t.Fatalf("WriteServerConfig: %v", err)
	}
}

func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv("MONGO_URI", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("BASE_URL", "")
	t.Setenv("MODE", "")
}

func TestLoadConfigYamlWins(t *testing.T) {
	writeConfigFile(t, &models.ServerConfig{
		Mode:      "single",
		BaseURL:   "http://yaml.example",
		MongoURI:  "mongodb://yaml",
		RedisAddr: "yaml:6379",
	})
	clearEnv(t)

	cfg := loadConfig()
	if cfg.Mode != "single" || cfg.BaseURL != "http://yaml.example" ||
		cfg.MongoURI != "mongodb://yaml" || cfg.RedisAddr != "yaml:6379" {
		t.Fatalf("loadConfig = %+v, want yaml values", cfg)
	}
}

func TestLoadConfigEnvOverridesYaml(t *testing.T) {
	writeConfigFile(t, &models.ServerConfig{
		Mode:     "single",
		BaseURL:  "http://yaml.example",
		MongoURI: "mongodb://yaml",
	})
	t.Setenv("MONGO_URI", "mongodb://env")
	t.Setenv("BASE_URL", "http://env.example")
	t.Setenv("MODE", "multi")
	t.Setenv("REDIS_ADDR", "env:6379")

	cfg := loadConfig()
	if cfg.Mode != "multi" || cfg.BaseURL != "http://env.example" ||
		cfg.MongoURI != "mongodb://env" || cfg.RedisAddr != "env:6379" {
		t.Fatalf("loadConfig = %+v, want env values", cfg)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("LNK_SERVER_CONFIG", filepath.Join(t.TempDir(), "absent.yaml"))
	clearEnv(t)

	cfg := loadConfig()
	if cfg.MongoURI != "mongodb://lnk:lnk@localhost:27017/lnk?authSource=admin" {
		t.Fatalf("MongoURI default = %q", cfg.MongoURI)
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("BaseURL default = %q", cfg.BaseURL)
	}
	if cfg.Mode != "" {
		t.Fatalf("Mode default = %q, want empty", cfg.Mode)
	}
}
