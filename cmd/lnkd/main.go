package main

import (
	"os"

	"github.com/baydogan/lnk/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		os.Exit(1)
	}
}
