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
		newGenerateTemplates(),
	)
	cmd.PersistentFlags().BoolVarP(&flagVerbose, "verbose", "v", false, "enable verbose log")
	return cmd
}

func newGenerateTemplates() *cobra.Command {
	var pkgName string
	var outputDir string
	var prefix string

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate Packed source files",
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
			fs, err := a.packageFiles(pkgName)
			if err != nil {
				return err
			}
			os.MkdirAll(outputDir, os.ModeDir|0755)
			for _, filePath := range fs {
				e := func() error {
					members := reachableFiles(g, filePath)

					w := os.Stdout
					if outputDir != "" {
						var e error
						w, e = os.Create(outputFilePath(outputDir, prefix, filePath))
						if e != nil {
							return e
						}
						defer w.Close()
					}
					if err := a.WriteToPackedCode(w, pkgName, members); err != nil {
						return err
					}
					return nil
				}()
				if e != nil {
					return e
				}
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&outputDir, "out", "o", "", "output directory path")
	cmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "generated file's prefix")
	cmd.PersistentFlags().StringVar(&pkgName, "package", "", "set target package name (default: directory name)")
	return cmd
}

func outputFilePath(dir, prefix, srcFile string) string {
	baseName := path.Clean(srcFile)
	return path.Join(dir, fmt.Sprintf("%s%s", prefix, baseName))
}

func newPackCommand() *cobra.Command {

	var pkgName string
	var outputFile string

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

			w := os.Stdout
			if outputFile != "" {
				var e error
				w, e = os.Create(outputFile)
				if e != nil {
					return e
				}
				defer w.Close()
			}
			return a.WriteToPackedCode(w, pkgName, members)
		},
	}
	cmd.PersistentFlags().StringVarP(&outputFile, "out", "o", "", "output file path (default: stdout)")
	cmd.PersistentFlags().StringVar(&pkgName, "package", "", "set target package name (default: directory name)")
	return cmd
}
