package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/NAEOS-foundation/naeos/internal/devexperience"
)

func newDXCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dx",
		Short: "Developer experience tools",
		Long: `Generate VS Code extensions, CLI completions, and code snippets.

Example:
  naeos dx vscode-gen
  naeos dx completion-bash
  naeos dx completion-zsh
  naeos dx completion-powershell
  naeos dx snippet-list
  naeos dx snippet-get --name project`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newDXVSCodeGenCommand())
	cmd.AddCommand(newDXCompletionBashCommand())
	cmd.AddCommand(newDXCompletionZshCommand())
	cmd.AddCommand(newDXCompletionPSCommand())
	cmd.AddCommand(newDXSnippetListCommand())
	cmd.AddCommand(newDXSnippetGetCommand())

	return cmd
}

func newDXVSCodeGenCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "vscode-gen",
		Short: "Generate VS Code extension package",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ext := devexperience.NewVSCodeExtension(
				"naeos", "1.0.0", "NAEOS project support", "NAEOS",
				[]string{"yaml", "json"},
			)

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "package.json:\n")
			fmt.Fprintf(out, "%s\n\n", ext.GeneratePackageJSON())
			fmt.Fprintf(out, "syntax.json:\n")
			fmt.Fprintf(out, "%s\n", ext.GenerateSyntaxJSON())
			return nil
		},
	}
}

func newDXCompletionBashCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion-bash",
		Short: "Generate bash completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := devexperience.NewCompletionEngine()
			fmt.Fprintln(cmd.OutOrStdout(), engine.GenerateBashCompletion())
			return nil
		},
	}
}

func newDXCompletionZshCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion-zsh",
		Short: "Generate zsh completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := devexperience.NewCompletionEngine()
			fmt.Fprintln(cmd.OutOrStdout(), engine.GenerateZshCompletion())
			return nil
		},
	}
}

func newDXCompletionPSCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "completion-powershell",
		Short: "Generate PowerShell completion script",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := devexperience.NewCompletionEngine()
			fmt.Fprintln(cmd.OutOrStdout(), engine.GeneratePowerShellCompletion())
			return nil
		},
	}
}

func newDXSnippetListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "snippet-list",
		Short: "List available code snippets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sm := devexperience.NewSnippetManager()

			names := sm.List()
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "%-20s\n", "SNIPPET")
			fmt.Fprintf(out, "%-20s\n", "--------")
			for _, name := range names {
				fmt.Fprintf(out, "%-20s\n", name)
			}
			return nil
		},
	}
}

func newDXSnippetGetCommand() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "snippet-get",
		Short: "Get a code snippet",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			sm := devexperience.NewSnippetManager()

			snippet, ok := sm.Get(name)
			if !ok {
				return fmt.Errorf("snippet '%s' not found", name)
			}

			fmt.Fprintln(cmd.OutOrStdout(), snippet)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "snippet name (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
