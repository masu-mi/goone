package main

import (
	"fmt"
	"os"

	"github.com/masu-mi/goone/cmd"
	"github.com/spf13/cobra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cmd := newCommand()
	return cmd.Execute()
}

func newCommand() *cobra.Command {
	c := &cobra.Command{Use: "goone"}
	c.AddCommand(
		cmd.NewPackCommand(),
		cmd.NewGenerateTemplates(),
	)
	return c
}
