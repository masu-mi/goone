package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/masu-mi/goone/model"
	"github.com/spf13/cobra"
)

func NewGenerateTemplates() *cobra.Command {
	var pkgName string
	var outputDir string
	var prefix string
	var includeTest bool

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate Packed source files",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			// use case
			filePath := args[0]
			dir := path.Dir(filePath)
			a, err := model.GetASF(dir)
			if err != nil {
				return err
			}
			if pkgName == "" {
				pkgName = path.Base(dir)
			}
			g, err := a.ParseAsDefGraph(pkgName)
			if err != nil {
				return err
			}
			fs, err := a.PackageFiles(pkgName)
			if err != nil {
				return err
			}
			os.MkdirAll(outputDir, os.ModeDir|0755)
			for _, filePath := range fs {
				if !includeTest && strings.HasSuffix(filePath, "_test.go") {
					continue
				}
				e := func() (e error) {
					members := model.ReachableFiles(g, filePath)

					w := os.Stdout
					if outputDir != "" {
						w, e = os.Create(outputFilePath(outputDir, prefix, filePath))
						if e != nil {
							return e
						}
						defer w.Close()
					}
					if e = a.WriteToPackedCode(w, pkgName, members); e != nil {
						return e
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
	cmd.PersistentFlags().BoolVar(&includeTest, "include-test", false, "Include test file as target (default: false)")
	return cmd
}
