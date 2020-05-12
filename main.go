package main

import (
	"fmt"
	"os"
	"path"

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
	var (
		flagVerbose bool
	)

	cmd := &cobra.Command{Use: "goone"}
	cmd.AddCommand(
		newPackCommand(),
	)
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "enable verbose log")
	return cmd
}

func newPackCommand() *cobra.Command {
	var pkgName string
	cmd := &cobra.Command{
		Use:   "pack",
		Short: "Pack target source code file with depended files share the package name into single file",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// use case
			filePath := args[0]
			dir := path.Dir(filePath)
			a, err := getASF(dir)
			if err != nil {
				return err
			}
			if pkgName == "" {
				pkgName = path.Base(dir)
			}
			g, err := a.parseAsDefGraph(pkgName)
			if err != nil {
				return err
			}
			members := reachableFiles(g, filePath)
			return a.WriteToPackedCode(os.Stdout, pkgName, members)
		},
	}
	cmd.PersistentFlags().StringVarP(&pkgName, "package", "p", "", "set target package name (default: directory name)")
	return cmd
}
