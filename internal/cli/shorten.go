package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baydogan/lnk/domain"
	"github.com/baydogan/lnk/internal/client"
	"github.com/spf13/cobra"
)

var (
	shortenAlias   string
	shortenExpires string
	shortenQR      bool
)

var shortenCmd = &cobra.Command{
	Use:     "shorten [url]",
	Short:   "Shorten a URL",
	Long:    `Create a short URL that redirects to the given long URL.`,
	Args:    cobra.ExactArgs(1),
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		req := domain.ShortenRequest{
			URL:     args[0],
			Alias:   shortenAlias,
			Expires: shortenExpires,
		}

		body, status, err := client.DefaultClient.Post("/api/v1/shorten", req)

		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}

		if status != http.StatusCreated {
			return apiError(body, status)
		}

		var resp domain.ShortenResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("server error (status %d): %w", status, err)
		}

		fmt.Printf("✓ %s\n", resp.ShortURL)
		if shortenQR {
			printTerminalQR(resp.ShortURL)
		}
		return nil
	},
}

func init() {
	shortenCmd.Flags().StringVar(&shortenAlias, "alias", "", "custom alias (path) for the short link")
	shortenCmd.Flags().StringVar(&shortenExpires, "expires", "", "expiry TTL, e.g. 30m, 1h, 7d, 2w")
	shortenCmd.Flags().BoolVar(&shortenQR, "qr", false, "also print a QR code for the short link")
}
