package main

import (
	"os"

	"github.com/baydogan/lnk/internal/lnkd"
)

func main() {
	if err := lnkd.Execute(); err != nil {
		os.Exit(1)
	}
}
