package main

import (
	"errors"
	"net/http"

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
}
