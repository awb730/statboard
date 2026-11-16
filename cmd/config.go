package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"statboard/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the list of tracked repositories",
}

var configAddCmd = &cobra.Command{
	Use:   "add <path>",
	Short: "Add a repository to the tracked list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Add(args[0]); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Println("Added:", args[0])
		return nil
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a repository from the tracked list",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if err := cfg.Remove(args[0]); err != nil {
			return err
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Println("Removed:", args[0])
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tracked repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Repos) == 0 {
			fmt.Println("No repositories tracked yet. Add one with: statboard config add <path>")
			return nil
		}
		for _, r := range cfg.Repos {
			fmt.Println(r)
		}
		return nil
	},
}

func init() {
	configCmd.AddCommand(configAddCmd, configRemoveCmd, configListCmd)
	rootCmd.AddCommand(configCmd)
}
