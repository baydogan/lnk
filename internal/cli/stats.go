package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/baydogan/lnk/internal/client"
	"github.com/baydogan/lnk/internal/models"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats <code>",
	Short: "Show stats for a short link by code or alias",
	Args:  cobra.ExactArgs(1),
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

		alias := "-"
		if u.Alias != nil {
			alias = *u.Alias
		}
		expires := "never"
		if u.ExpiresAt != nil {
			expires = u.ExpiresAt.Format("2006-01-02 15:04")
		}

		rows := [][]string{
			{"Short URL", u.ShortURL},
			{"Code", u.Code},
			{"Alias", alias},
			{"Original", u.OriginalURL},
			{"Clicks", fmt.Sprintf("%d", u.ClickCount)},
			{"Created", u.CreatedAt.Format("2006-01-02 15:04")},
			{"Expires", expires},
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			Rows(rows...).
			StyleFunc(func(row, col int) lipgloss.Style {
				st := lipgloss.NewStyle().Padding(0, 1)
				if col == 0 {
					st = st.Bold(true)
				}
				return st
			})
		fmt.Println(t)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
