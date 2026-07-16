package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/config"
	"github.com/spf13/cobra"
)

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage your API key",
}

var keyRotateCmd = &cobra.Command{
	Use:     "rotate",
	Short:   "Rotate your API key — the old one stops working immediately",
	Args:    cobra.NoArgs,
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Post("/api/v1/keys/rotate", nil)
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusOK {
			return apiError(body, status)
		}

		var res domain.KeyResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}

		cfg, err := config.ReadClientConfig()
		if err != nil {
			return err
		}
		cfg.APIKey = res.APIKey
		if _, err := config.WriteClientConfig(&cfg); err != nil {
			return err
		}

		fmt.Println("key rotated — new key saved to client config; the old key no longer works")
		return nil
	},
}

func init() {
	keyCmd.AddCommand(keyRotateCmd)
	rootCmd.AddCommand(keyCmd)
}
