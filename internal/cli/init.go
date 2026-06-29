package cli

import (
	"github.com/baydogan/lnk/internal/config"
	"github.com/baydogan/lnk/internal/models"
	"github.com/baydogan/lnk/internal/setup"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively configure the lnk server",
	Long:  "Launches a TUI wizard that picks the deployment mode and writes server config.\nAny value passed as a flag is used directly and not prompted.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &models.ServerConfig{}
		cfg.Mode, _ = cmd.Flags().GetString("mode")
		cfg.BaseURL, _ = cmd.Flags().GetString("base-url")
		cfg.Admin, _ = cmd.Flags().GetString("admin")
		cfg.MongoURI, _ = cmd.Flags().GetString("mongo-uri")
		cfg.RedisAddr, _ = cmd.Flags().GetString("redis-addr")

		p := setup.Provided{
			Mode:    cmd.Flags().Changed("mode"),
			BaseURL: cmd.Flags().Changed("base-url"),
			Admin:   cmd.Flags().Changed("admin"),
			Mongo:   cmd.Flags().Changed("mongo-uri"),
			Redis:   cmd.Flags().Changed("redis-addr"),
		}

		configureClient, err := setup.Run(cfg, p)
		if err != nil {
			return err
		}

		path, exists, err := config.ServerConfigExists()
		if err != nil {
			return err
		}
		if exists {
			overwrite, err := setup.ConfirmOverwrite(path)
			if err != nil {
				return err
			}
			if !overwrite {
				setup.Kept("kept existing server config at", path)
				return nil
			}
		}

		written, err := config.WriteServerConfig(cfg)
		if err != nil {
			return err
		}

		setup.Summary(cfg, configureClient)
		setup.Saved("server config written to", written)

		// TODO: when configureClient, write ~/.lnk/config.yaml with the
		// generated admin key.
		return nil
	},
}

func init() {
	initCmd.Flags().String("mode", "", "deployment mode: single | multi")
	initCmd.Flags().String("base-url", "", "public base URL for short links, e.g. https://shorturl.com")
	initCmd.Flags().String("admin", "", "admin username (multi-user)")
	initCmd.Flags().String("mongo-uri", "", "MongoDB connection URI")
	initCmd.Flags().String("redis-addr", "", "Redis address")
	rootCmd.AddCommand(initCmd)
}
