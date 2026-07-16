package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users (admin, multi-user mode)",
}

var userCreateCmd = &cobra.Command{
	Use:     "create <username>",
	Short:   "Create a user and print their API key",
	Args:    cobra.ExactArgs(1),
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Post("/api/v1/users", map[string]string{"username": args[0]})
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusCreated {
			return apiError(body, status)
		}

		var res models.CreateUserResponse
		if err := json.Unmarshal(body, &res); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}

		fmt.Printf("user %q created (role %s). Share:\n\n  lnk login --server %s --api-key %s\n\n",
			res.User.Username, res.User.Role, serverURL, res.APIKey)
		return nil
	},
}

var userListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List users",
	Args:    cobra.NoArgs,
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Get("/api/v1/users")
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusOK {
			return apiError(body, status)
		}

		var users []models.User
		if err := json.Unmarshal(body, &users); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}
		if len(users) == 0 {
			fmt.Println("no users yet")
			return nil
		}

		rows := make([][]string, 0, len(users))
		for _, u := range users {
			rows = append(rows, []string{u.Username, u.Role, u.CreatedAt.Local().Format("2006-01-02")})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			Headers("USERNAME", "ROLE", "CREATED").
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				st := lipgloss.NewStyle().Padding(0, 1)
				if row == table.HeaderRow {
					st = st.Bold(true)
				}
				return st
			})
		fmt.Println(t)
		return nil
	},
}

var userDeleteCmd = &cobra.Command{
	Use:     "delete <username>",
	Short:   "Delete a user and revoke their API key",
	Args:    cobra.ExactArgs(1),
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Delete("/api/v1/users/" + args[0])
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusNoContent {
			return apiError(body, status)
		}

		fmt.Printf("user %q deleted\n", args[0])
		return nil
	},
}

var userWhoamiCmd = &cobra.Command{
	Use:     "whoami",
	Short:   "Show the currently logged-in user",
	Args:    cobra.NoArgs,
	PreRunE: requireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Get("/api/v1/me")
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusOK {
			return apiError(body, status)
		}

		var u models.User
		if err := json.Unmarshal(body, &u); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}

		fmt.Printf("%s (%s)\n", u.Username, u.Role)
		return nil
	},
}

func init() {
	userCmd.AddCommand(userCreateCmd, userListCmd, userDeleteCmd, userWhoamiCmd)
	rootCmd.AddCommand(userCmd)
}
