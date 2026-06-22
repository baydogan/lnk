package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/baydogan/lnk/internal/database"
	"github.com/baydogan/lnk/internal/logger"
	"github.com/baydogan/lnk/internal/server"
)

func main() {
	logger.Setup("debug", false)

	if err := database.Connect("mongodb://lnk:lnk@localhost:27017/lnk?authSource=admin"); err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	router := server.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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

}
