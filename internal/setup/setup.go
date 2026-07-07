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
	defaultMongoURI  = "mongodb://lnk:lnk@mongodb:27017/lnk?authSource=admin"
	defaultRedisAddr = "redis:6379"
)

type Provided struct {
	Mode    bool
	BaseURL bool
	Admin   bool
	Mongo   bool
	Redis   bool
}

func Run(cfg *models.ServerConfig, p Provided) error {
	if p.Mode && cfg.Mode != "single" && cfg.Mode != "multi" {
		return fmt.Errorf("invalid mode %q: must be \"single\" or \"multi\"", cfg.Mode)
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
	).WithTheme(Theme()).WithShowHelp(true)

	hasPrompts := !p.Mode ||
		!p.BaseURL ||
		(!p.Admin && cfg.Mode == "multi") ||
		!p.Mongo ||
		!p.Redis

	if hasPrompts {
		fmt.Println(banner())

		if err := form.Run(); err != nil {
			return err
		}
	}

	return nil
}

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

func Summary(cfg *models.ServerConfig) {
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

	title := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("✔ configuration ready")
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neon).
		Padding(0, 2).
		Render(title + "\n\n" + strings.Join(rows, "\n"))

	fmt.Println(card)
}

func Saved(label, path string) {
	mark := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("✔")
	desc := lipgloss.NewStyle().Foreground(subtle).Render(label)
	loc := lipgloss.NewStyle().Foreground(violet).Bold(true).Render(path)
	fmt.Println(mark + " " + desc + " " + loc)
}

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

func required(label string) func(string) error {
	return func(s string) error {
		if strings.TrimSpace(s) == "" {
			return fmt.Errorf("%s is required", label)
		}
		return nil
	}
}
