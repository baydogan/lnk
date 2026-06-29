// Command lnkd is the lnk server daemon. It runs the HTTP server and nothing
// else — all client/setup commands live in the separate `lnk` binary.
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
