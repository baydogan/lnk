package setup

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	neon   = lipgloss.AdaptiveColor{Light: "#0CA678", Dark: "#00F0B5"}
	violet = lipgloss.AdaptiveColor{Light: "#7048E8", Dark: "#A98BFF"}
	pink   = lipgloss.AdaptiveColor{Light: "#E64980", Dark: "#FF7AC6"}
	text   = lipgloss.AdaptiveColor{Light: "#1E1E2E", Dark: "#E6E6F0"}
	subtle = lipgloss.AdaptiveColor{Light: "#868E96", Dark: "#7C7F96"}
	faint  = lipgloss.AdaptiveColor{Light: "#ADB5BD", Dark: "#494D64"}
	danger = lipgloss.AdaptiveColor{Light: "#E03131", Dark: "#FF6B6B"}
	ink    = lipgloss.Color("#11111B") // text on top of a neon fill
)

func Theme() *huh.Theme {
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderForeground(neon).
		PaddingLeft(1)
	t.Blurred.Base = t.Focused.Base.BorderForeground(faint)

	t.Focused.Title = t.Focused.Title.Foreground(neon).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(neon).Bold(true)
	t.Focused.Description = t.Focused.Description.Foreground(subtle)
	t.Blurred.Title = t.Blurred.Title.Foreground(subtle)
	t.Blurred.Description = t.Blurred.Description.Foreground(faint)

	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(pink).SetString("❯ ")
	t.Focused.Option = t.Focused.Option.Foreground(text)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(neon).Bold(true)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(neon).SetString("◉ ")
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(faint).SetString("○ ")
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(pink).SetString("❯ ")

	t.Focused.FocusedButton = t.Focused.FocusedButton.
		Foreground(ink).Background(neon).Bold(true).Padding(0, 2)
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(text).Background(faint).Padding(0, 2)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(pink)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(violet)
	t.Focused.TextInput.Text = t.Focused.TextInput.Text.Foreground(text)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(faint)

	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(danger)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(danger)

	t.Help.Ellipsis = t.Help.Ellipsis.Foreground(faint)
	t.Help.ShortKey = t.Help.ShortKey.Foreground(violet)
	t.Help.ShortDesc = t.Help.ShortDesc.Foreground(subtle)
	t.Help.ShortSeparator = t.Help.ShortSeparator.Foreground(faint)
	t.Help.FullKey = t.Help.FullKey.Foreground(violet)
	t.Help.FullDesc = t.Help.FullDesc.Foreground(subtle)
	t.Help.FullSeparator = t.Help.FullSeparator.Foreground(faint)

	return t
}
