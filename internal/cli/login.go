package cli

import (
	"errors"
	"fmt"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/config"
	"github.com/spf13/cobra"
)

var loginAPIKey string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save the server URL and API key to the client config",
	Long:  "Writes ~/.lnk/config.yaml so subsequent commands authenticate automatically.\nRun the command printed by lnkd on first start.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginAPIKey == "" {
			return errors.New("--api-key is required")
		}

		cfg := &domain.ClientConfig{Server: serverURL, APIKey: loginAPIKey}
		path, err := config.WriteClientConfig(cfg)
		if err != nil {
			return err
		}

		fmt.Println("logged in — client config written to", path)
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginAPIKey, "api-key", "", "API key printed by lnkd")
	rootCmd.AddCommand(loginCmd)
}
