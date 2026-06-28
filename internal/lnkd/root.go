// Package lnkd wires the lnk server daemon's CLI (start + init).
package lnkd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "lnkd",
	Short: "lnk server daemon",
	Long:  "lnkd is the lnk server daemon: it serves redirects and the API, and handles interactive setup.",
}

func Execute() error {
	return rootCmd.Execute()
}
