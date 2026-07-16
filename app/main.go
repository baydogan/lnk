package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/baydogan/lnk/auth"
	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/config"
	"github.com/baydogan/lnk/internal/logger"
	mongorepo "github.com/baydogan/lnk/internal/repository/mongo"
	"github.com/baydogan/lnk/internal/rest"
	"github.com/baydogan/lnk/url"
	"github.com/baydogan/lnk/user"
)

func main() {
	logger.Setup("debug", false)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to read server config")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	db, err := mongorepo.Connect(ctx, cfg.MongoURI)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	keyRepo := mongorepo.NewAPIKeyRepository(db)
	userRepo := mongorepo.NewUserRepository(db)
	urlRepo := mongorepo.NewURLRepository(db)

	authSvc := auth.NewService(keyRepo)
	userSvc := user.NewService(userRepo, keyRepo)
	urlHandler := rest.NewHTTPHandler(url.NewService(urlRepo, cfg.BaseURL))
	userHandler := rest.NewUserHandler(userSvc)

	if err := authSvc.EnsureIndexes(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to ensure indexes")
	}

	var (
		pt      string
		created bool
	)
	if cfg.Mode == domain.ModeMulti {
		if err := userSvc.EnsureIndexes(ctx); err != nil {
			logger.Fatal().Err(err).Msg("failed to ensure user indexes")
		}
		admin := cfg.Admin
		if admin == "" {
			admin = "admin"
		}
		pt, created, err = userSvc.EnsureAdmin(ctx, admin)
	} else {
		pt, created, err = authSvc.EnsureAdminKey(ctx)
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ensure admin key")
	}

	if created {
		fmt.Printf("\nAdmin API key generated. Run:\n\n  lnk login --server %s --api-key %s\n\n", cfg.BaseURL, pt)
	}

	router := rest.NewRouter(cfg.Mode, urlHandler, userHandler, authSvc, userSvc)
	srv := &http.Server{Addr: ":" + port, Handler: router}

	logger.Info().Str("port", port).Msg("lnk server starting")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
