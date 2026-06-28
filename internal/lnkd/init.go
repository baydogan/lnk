package lnkd

import (
	"github.com/baydogan/lnk/internal/setup"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Interactively configure the lnk server",
	Long:  "Launches a TUI wizard that picks the deployment mode and writes server config.\nAny value passed as a flag is used directly and not prompted.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &setup.ServerConfig{}
		cfg.Mode, _ = cmd.Flags().GetString("mode")
		cfg.AdminUser, _ = cmd.Flags().GetString("admin")
		cfg.MongoURI, _ = cmd.Flags().GetString("mongo-uri")
		cfg.RedisAddr, _ = cmd.Flags().GetString("redis-addr")

		p := setup.Provided{
			Mode:  cmd.Flags().Changed("mode"),
			Admin: cmd.Flags().Changed("admin"),
			Mongo: cmd.Flags().Changed("mongo-uri"),
			Redis: cmd.Flags().Changed("redis-addr"),
		}

		if err := setup.Run(cfg, p); err != nil {
			return err
		}

		setup.Summary(cfg)
		// TODO: persist cfg to ~/.lnk/server.yaml and, when ConfigureClient,
		// write ~/.lnk/config.yaml with the generated admin key.
		return nil
	},
}

func init() {
	initCmd.Flags().String("mode", "", "deployment mode: single | multi")
	initCmd.Flags().String("admin", "", "admin username (multi-user)")
	initCmd.Flags().String("mongo-uri", "", "MongoDB connection URI")
	initCmd.Flags().String("redis-addr", "", "Redis address")
	rootCmd.AddCommand(initCmd)
}
