package cmd

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	var confPath string

	cmd := &cobra.Command{Use: "goone"}
	cmd.AddCommand(
		NewPackCommand(),
		NewGenerateTemplates(),
	)
	cmd.PersistentFlags().StringVarP(&confPath, "config", "c", "", "config file path (default: ~/.config/goone/config.toml)")
	return cmd
}
