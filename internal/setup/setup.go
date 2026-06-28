package setup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultMongoURI  = "mongodb://localhost:27017/lnk"
	defaultRedisAddr = "localhost:6379"
)

type ServerConfig struct {
	Mode            string // "single" | "multi"
	AdminUser       string
	MongoURI        string
	RedisAddr       string
	ConfigureClient bool
}

type Provided struct {
	Mode  bool
	Admin bool
	Mongo bool
	Redis bool
}

func Run(cfg *ServerConfig, p Provided) error {
	if !p.Mode && cfg.Mode == "" {
		cfg.Mode = "single"
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
				Title("Admin username").
				Description("Owner of the first API key.").
				Placeholder("admin").
				Value(&cfg.AdminUser).
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
				Value(&cfg.ConfigureClient),
		).WithHideFunc(func() bool { return cfg.Mode != "multi" }),
	).WithTheme(Theme()).WithShowHelp(true)

	fmt.Println(banner())

	if err := form.Run(); err != nil {
		return err
	}

	if cfg.Mode != "multi" {
		cfg.ConfigureClient = true
	}

	return nil
}

func Summary(cfg *ServerConfig) {
	label := lipgloss.NewStyle().Foreground(subtle).Width(8)
	val := lipgloss.NewStyle().Foreground(text).Bold(true)

	rows := []string{
		label.Render("mode") + val.Render(cfg.Mode),
		label.Render("mongo") + val.Render(cfg.MongoURI),
		label.Render("redis") + val.Render(cfg.RedisAddr),
	}
	if cfg.Mode == "multi" {
		rows = append(rows, label.Render("admin")+val.Render(cfg.AdminUser))
	}
	rows = append(rows, label.Render("client")+val.Render(yesno(cfg.ConfigureClient)))

	title := lipgloss.NewStyle().Foreground(neon).Bold(true).Render("✔ configuration ready")
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(neon).
		Padding(0, 2).
		Render(title + "\n\n" + strings.Join(rows, "\n"))

	fmt.Println(card)
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
