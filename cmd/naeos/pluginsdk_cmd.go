package main

import (
	"fmt"
	"github.com/NAEOS-foundation/naeos/internal/pluginsdk"
	"github.com/spf13/cobra"
)

func newPluginSDKCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pluginsdk",
		Short: "Plugin SDK management commands",
		Long:  `Manage and develop custom plugins for NAEOS.`,
	}

	cmd.AddCommand(newPluginSDKListCommand())
	cmd.AddCommand(newPluginSDKInfoCommand())

	return cmd
}

func newPluginSDKListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available plugins",
		Run: func(cmd *cobra.Command, args []string) {
			manager := pluginsdk.NewManager()
			plugins := manager.List()
			if len(plugins) == 0 {
				fmt.Println("No plugins registered")
				return
			}
			for _, name := range plugins {
				fmt.Printf("  - %s\n", name)
			}
		},
	}
}

func newPluginSDKInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info [name]",
		Short: "Show plugin information",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			fmt.Printf("Plugin: %s\n", name)
			fmt.Println("Use 'naeos plugin install <name>' to install a plugin")
		},
	}
}
