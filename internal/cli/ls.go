package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all short links",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		body, status, err := client.DefaultClient.Get("/api/v1/urls")
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}
		if status != http.StatusOK {
			return apiError(body, status)
		}

		var urls []models.URLResponse
		if err := json.Unmarshal(body, &urls); err != nil {
			return fmt.Errorf("bad response (status %d): %w", status, err)
		}

		if len(urls) == 0 {
			fmt.Println("no links yet")
			return nil
		}

		rows := make([][]string, 0, len(urls))
		for _, u := range urls {
			rows = append(rows, []string{
				u.ShortURL,
				fmt.Sprintf("%d", u.ClickCount),
				u.CreatedAt.Local().Format("2006-01-02"),
				u.OriginalURL,
			})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			Headers("SHORT URL", "CLICKS", "CREATED", "ORIGINAL").
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

func apiError(body []byte, status int) error {
	var e map[string]string
	if json.Unmarshal(body, &e) == nil && e["error"] != "" {
		return errors.New(e["error"])
	}
	return fmt.Errorf("request failed with status %d", status)
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
