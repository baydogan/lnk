package cli

import (
	"github.com/spf13/cobra"
)

var serverURL string

var rootCmd = &cobra.Command{
	Use:   "lnk",
	Short: "Terminal-based URL shortener",
	Long:  `lnk is a CLI tool that shortens URLs via a cloud backend. Shorten, list, track, and manage your links from the terminal.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "localhost:8080", "URL to connect to")
}
