package cli

import (
	"errors"
	"fmt"

	"github.com/baydogan/lnk/internal/config"
	"github.com/baydogan/lnk/internal/errs"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove the saved client config (server URL + API key)",
	Long:  "Deletes ~/.lnk/config.yaml so this machine no longer authenticates automatically.\nThis only clears the local session; it does not revoke the key on the server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := config.RemoveClientConfig()
		if err != nil {
			if errors.Is(err, errs.ErrNotLoggedIn) {
				fmt.Println("not logged in — nothing to remove")
				return nil
			}
			return err
		}
		fmt.Println("logged out — removed", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
