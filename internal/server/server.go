package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/baydogan/lnk/internal/container"
	"github.com/baydogan/lnk/internal/database"
	"github.com/baydogan/lnk/internal/logger"
)

// Run boots the lnk server: connects the database, bootstraps the admin key on
// first run, and serves until shutdown. It owns the full server lifecycle so
// the CLI layer only has to invoke it.
func Run() error {
	logger.Setup("debug", false)

	mongoURI := getenv("MONGO_URI", "mongodb://lnk:lnk@localhost:27017/lnk?authSource=admin")
	port := getenv("PORT", "8080")

	if err := database.Connect(mongoURI); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	c := container.New()

	pt, created, err := c.AuthService.EnsureAdminKey()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ensure admin key")
	}

	if created {
		fmt.Printf("\nAdmin API key generated. Run:\n\n  lnk login --server http://localhost:%s --api-key %s\n\n", port, pt)
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

// getenv returns the env var for key, or def when it is unset/empty.
func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
