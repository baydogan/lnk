package cli

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/baydogan/lnk/internal/client"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var rmYes bool

var rmCmd = &cobra.Command{
	Use:   "rm <code>",
	Short: "Delete a short link by code or alias",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code := args[0]

		if !rmYes {
			var confirmed bool
			if err := huh.NewConfirm().
				Title(fmt.Sprintf("Delete %q?", code)).
				Affirmative("Yes, delete").
				Negative("Cancel").
				Value(&confirmed).
				Run(); err != nil {
				return err
			}
			if !confirmed {
				fmt.Println("cancelled")
				return nil
			}
		}

		body, status, err := client.DefaultClient.Delete("/api/v1/" + url.PathEscape(code))
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusNoContent {
			return apiError(body, status)
		}

		fmt.Printf("✓ deleted %s\n", code)
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&rmYes, "yes", "y", false, "skip the confirmation prompt")
	rootCmd.AddCommand(rmCmd)
}
