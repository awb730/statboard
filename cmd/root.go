package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"statboard/internal/config"
	"statboard/internal/filter"
	"statboard/internal/git"
	"statboard/internal/render"
)

var (
	dirtyOnly  bool
	behindOnly bool
)

func init() {
	rootCmd.Flags().BoolVar(&dirtyOnly, "dirty", false, "only show repos with uncommitted changes")
	rootCmd.Flags().BoolVar(&behindOnly, "behind", false, "only show repos behind their remote-tracking branch")
}

var rootCmd = &cobra.Command{
	Use:          "statboard",
	Short:        "A read-only dashboard for the status of your Git repositories",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Repos) == 0 {
			fmt.Println("No repositories tracked yet. Add one with: statboard config add <path>")
			return nil
		}
		statuses := git.CheckAll(cfg.Repos)
		shown := filter.Apply(statuses, filter.Options{Dirty: dirtyOnly, Behind: behindOnly})
		if len(shown) == 0 {
			fmt.Printf("No repositories match the filters (checked %d).\n", len(statuses))
			return nil
		}
		render.Table(os.Stdout, shown)
		render.Summary(os.Stdout, shown, len(statuses))
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
