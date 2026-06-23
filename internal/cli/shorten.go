package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/models"
	"github.com/spf13/cobra"
)

var (
	shortenAlias   string
	shortenExpires string
)

var shortenCmd = &cobra.Command{
	Use:   "shorten [url]",
	Short: "Shorten a URL",
	Long:  `Create a short URL that redirects to the given long URL.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		req := models.ShortenRequest{
			URL:     args[0],
			Alias:   shortenAlias,
			Expires: shortenExpires,
		}

		body, status, err := client.DefaultClient.Post("/api/v1/shorten", req)

		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}

		if status != http.StatusCreated {
			var errResp map[string]string
			if err := json.Unmarshal(body, &errResp); err != nil {
				return fmt.Errorf("server error (status %d): %w", status, err)
			}
			return fmt.Errorf(errResp["error"])
		}

		var resp models.ShortenResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return fmt.Errorf("server error (status %d): %w", status, err)
		}

		fmt.Printf("✓ %s\n", resp.ShortURL)
		return nil
	},
}
