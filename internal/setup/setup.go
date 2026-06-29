package setup

import (
	"fmt"
	"strings"

	"github.com/baydogan/lnk/internal/models"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultBaseURL   = "http://localhost:8080"
	defaultMongoURI  = "mongodb://localhost:27017/lnk"
	defaultRedisAddr = "localhost:6379"
)

type Provided struct {
	Mode    bool
	BaseURL bool
	Admin   bool
	Mongo   bool
	Redis   bool
}

func Run(cfg *models.ServerConfig, p Provided) (bool, error) {
	// A mode supplied via flag bypasses the select, so validate it here; the
	// interactive picker can only ever yield a valid value.
	if p.Mode && cfg.Mode != "single" && cfg.Mode != "multi" {
		return false, fmt.Errorf("invalid mode %q: must be \"single\" or \"multi\"", cfg.Mode)
	}
	if !p.Mode && cfg.Mode == "" {
		cfg.Mode = "single"
	}
	if !p.BaseURL && cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if !p.Mongo && cfg.MongoURI == "" {
		cfg.MongoURI = defaultMongoURI
	}
	if !p.Redis && cfg.RedisAddr == "" {
		cfg.RedisAddr = defaultRedisAddr
	}

	var configureClient bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Deployment mode").
				Description("How will this server hand out access?").
				Options(
					huh.NewOption("Single-user   ·   one owner, one key", "single"),
					huh.NewOption("Multi-user    ·   admin + per-user keys", "multi"),
				).
				Value(&cfg.Mode),
		).WithHideFunc(func() bool { return p.Mode }),

		huh.NewGroup(
			huh.NewInput().
				Title("Public base URL").
				Description("Origin used in generated short links, e.g. https://shorturl.com").
				Placeholder(defaultBaseURL).
				Value(&cfg.BaseURL).
				Validate(required("public base URL")),
		).WithHideFunc(func() bool { return p.BaseURL }),

		huh.NewGroup(
			huh.NewInput().
				Title("Admin username").
				Description("Owner of the first API key.").
				Placeholder("admin").
				Value(&cfg.Admin).
				Validate(required("admin username")),
		).WithHideFunc(func() bool { return p.Admin || cfg.Mode != "multi" }),

		huh.NewGroup(
			huh.NewInput().
				Title("MongoDB URI").
				Description("Where links, users and keys are stored.").
				Value(&cfg.MongoURI).
				Validate(required("MongoDB URI")),
		).WithHideFunc(func() bool { return p.Mongo }),

		huh.NewGroup(
			huh.NewInput().
				Title("Redis address").
				Description("Backs click counters and rate limiting.").
				Value(&cfg.RedisAddr).
				Validate(required("Redis address")),
		).WithHideFunc(func() bool { return p.Redis }),

		huh.NewGroup(
			huh.NewConfirm().
				Title("Configure local lnk client now?").
				Description("Writes ~/.lnk/config.yaml for this machine so you\ncan use lnk here without copy-pasting the admin key.").
				Affirmative("Yes, set it up").
				Negative("No, remote only").
				Value(&configureClient),
		).WithHideFunc(func() bool { return cfg.Mode != "multi" }),
	).WithTheme(Theme()).WithShowHelp(true)

	// Only prompt when at least one group is visible. When every value arrives
	// via flags, there is nothing to ask, so we skip the form entirely — this
	// lets CI / Docker run `lnkd init` non-interactively (no TTY required).
	hasPrompts := !p.Mode ||
		!p.BaseURL ||
		(!p.Admin && cfg.Mode == "multi") ||
		!p.Mongo ||
		!p.Redis ||
		cfg.Mode == "multi"

	if hasPrompts {
		fmt.Println(banner())

		if err := form.Run(); err != nil {
			return false, err
		}
	}

	// Single-user always configures the local client.
	if cfg.Mode != "multi" {
		configureClient = true
	}

	return configureClient, nil
}

// ConfirmOverwrite asks whether to overwrite an existing config file at path.
// It uses the same theme as the wizard for a consistent look.
func ConfirmOverwrite(path string) (bool, error) {
	var overwrite bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Config already exists").
				Description(fmt.Sprintf("%s already exists.\nOverwrite it with the new configuration?", path)).
				Affirmative("Yes, overwrite").
				Negative("No, keep it").
				Value(&overwrite),
		),
	).WithTheme(Theme()).WithShowHelp(true)

	if err := form.Run(); err != nil {
		return false, err
	}
	return overwrite, nil
}

func Summary(cfg *models.ServerConfig, configureClient bool) {
	label := lipgloss.NewStyle().Foreground(subtle).Width(8)
	val := lipgloss.NewStyle().Foreground(text).Bold(true)

	rows := []string{
		label.Render("mode") + val.Render(cfg.Mode),
		label.Render("url") + val.Render(cfg.BaseURL),
		label.Render("mongo") + val.Render(cfg.MongoURI),
		label.Render("redis") + val.Render(cfg.RedisAddr),
	}
	if cfg.Mode == "multi" {
		rows = append(rows, label.Render("admin")+val.Render(cfg.Admin))
	}
	rows = append(rows, label.Render("client")+val.Render(yesno(configureClient)))

	title := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("✔ configuration ready")
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neon).
		Padding(0, 2).
		Render(title + "\n\n" + strings.Join(rows, "\n"))

	fmt.Println(card)
}

// Saved prints a styled confirmation that a config file was written to path.
func Saved(label, path string) {
	mark := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("✔")
	desc := lipgloss.NewStyle().Foreground(subtle).Render(label)
	loc := lipgloss.NewStyle().Foreground(violet).Bold(true).Render(path)
	fmt.Println(mark + " " + desc + " " + loc)
}

// Kept prints a styled note that an existing config file was left untouched.
func Kept(label, path string) {
	dot := lipgloss.NewStyle().Foreground(faint).Render("•")
	desc := lipgloss.NewStyle().Foreground(subtle).Render(label)
	loc := lipgloss.NewStyle().Foreground(violet).Render(path)
	fmt.Println(dot + " " + desc + " " + loc)
}

func banner() string {
	logo := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("lnk")
	dot := lipgloss.NewStyle().Foreground(faint).Render("·")
	tag := lipgloss.NewStyle().Foreground(violet).Render("server setup")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(violet).
		Padding(0, 3)
	return box.Render(logo + "  " + dot + "  " + tag)
}

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func required(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s is required", label)
		}
		return nil
	}
}
