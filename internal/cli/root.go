package cli

import (
	"errors"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/config"
	"github.com/spf13/cobra"
)

var serverURL string

var rootCmd = &cobra.Command{
	Use:   "lnk",
	Short: "Terminal-based URL shortener",
	Long:  `lnk is a CLI tool that shortens URLs via local or cloud backend. Shorten, list, track, and manage your links from the terminal.`,
}

func requireClient(cmd *cobra.Command, args []string) error {
	url := serverURL
	var token string

	if cfg, err := config.ReadClientConfig(); err == nil {
		if !cmd.Flags().Changed("server") && cfg.Server != "" {
			url = cfg.Server
		}
		token = cfg.APIKey
	} else if !errors.Is(err, domain.ErrNotLoggedIn) {
		return err
	}

	client.DefaultClient = client.New(url)
	if token != "" {
		client.DefaultClient.SetToken(token)
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "URL to connect to")

	rootCmd.AddCommand(shortenCmd)
}
