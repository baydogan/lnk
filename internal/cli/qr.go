package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/models"
	"github.com/mdp/qrterminal/v3"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var qrOut string

var qrCmd = &cobra.Command{
	Use:     "qr <code>",
	Short:   "Show a QR code for a short link (terminal, or PNG with -o)",
	Long:    "Generates a QR code that encodes the short URL, so scanning it goes through the shortener (and counts the click). Nothing is stored — the QR is derived on demand.",
	Args:    cobra.ExactArgs(1),
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Get("/api/v1/urls/" + url.PathEscape(args[0]))
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusOK {
			return apiError(body, status)
		}

		var u models.URLResponse
		if err := json.Unmarshal(body, &u); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}

		if qrOut != "" {
			if err := qrcode.WriteFile(u.ShortURL, qrcode.Medium, 256, qrOut); err != nil {
				return fmt.Errorf("could not write PNG: %w", err)
			}
			fmt.Printf("✓ QR for %s written to %s\n", u.ShortURL, qrOut)
			return nil
		}

		fmt.Println(u.ShortURL)
		printTerminalQR(u.ShortURL)
		return nil
	},
}

func printTerminalQR(text string) {
	qrterminal.GenerateHalfBlock(text, qrterminal.M, os.Stdout)
}

func init() {
	qrCmd.Flags().StringVarP(&qrOut, "out", "o", "", "write a PNG file instead of terminal output")
	rootCmd.AddCommand(qrCmd)
}
