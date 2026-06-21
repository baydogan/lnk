package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/baydogan/lnk/internal/server"
)

func main() {

	router := server.NewRouter()

	srv := &http.Server{
		Addr:    ":" + "8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Print("log")
	}

}
