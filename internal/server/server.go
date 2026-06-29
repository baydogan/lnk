package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/baydogan/lnk/internal/config"
	"github.com/baydogan/lnk/internal/container"
	"github.com/baydogan/lnk/internal/database"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/baydogan/lnk/internal/models"
)

func Run() error {
	logger.Setup("debug", false)

	cfg := loadConfig()
	port := getenv("PORT", "8080")

	if err := database.Connect(cfg.MongoURI); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	c := container.New(cfg)

	if err := c.AuthService.EnsureIndexes(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("failed to ensure indexes")
	}

	pt, created, err := c.AuthService.EnsureAdminKey()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ensure admin key")
	}

	if created {
		fmt.Printf("\nAdmin API key generated. Run:\n\n  lnk login --server %s --api-key %s\n\n", cfg.BaseURL, pt)
	}

	router := NewRouter(c.URLHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info().Str("port", port).Msg("lnk server starting")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("server failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	return nil
}

func loadConfig() models.ServerConfig {
	cfg, _, err := config.ReadServerConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read server config")
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
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
