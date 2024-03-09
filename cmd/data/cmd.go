package data

import (
	"github.com/spf13/cobra"
	"github.com/xuxant/Darwin_API/pkg/config"
	"github.com/xuxant/Darwin_API/pkg/db"
	"github.com/xuxant/Darwin_API/pkg/logs"
	"github.com/xuxant/Darwin_API/pkg/service"
)

func runCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "run data management service",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfig()
			logger := logs.InitializeLogger(*cfg)
			initializer := service.Initialize(*cfg, logger)
			initializer.ListenTillDeath(*cfg)
		},
	}
	return runCmd
}

func migrateCommand() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "database migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.GetConfig()
			dbconn := db.GetConnection(*cfg)
			err := db.MigrateDatabase(dbconn)
			return err
		},
	}
	return migrateCmd
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "data-manager",
		Short: "data manager app for darwin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(migrateCommand())
	return rootCmd
}
