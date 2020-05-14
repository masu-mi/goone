package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func NewGenerateTemplates() *cobra.Command {
	var outputDir string
	var prefix string
	var templatePath string
	var includeTest bool

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate Packed source files",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			_, dir, pkgName := targetLocation(c, args)

			cfg := &config{}
			cfg.load(c.Flag("config").Value.String())
			t, err := cfg.NewTemplate(templatePath)
			if err != nil {
				return err
			}

			a, g, err := targetStructure(dir, pkgName)
			if err != nil {
				return err
			}

			fs, err := a.PackageFiles(pkgName)
			if err != nil {
				return err
			}
			if outputDir != "" {
				os.MkdirAll(outputDir, os.ModeDir|0755)
			}
			for src := range targetFiles(fs, includeTest) {
				dstFile := outputFilePath(outputDir, prefix, src)
				e := packCode(g, t, dstFile, pkgName, src)
				if e != nil {
					return e
				}
			}
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&outputDir, "out", "o", "", "output directory path")
	cmd.PersistentFlags().String("package", "", "set target package name (default: directory name)")
	cmd.PersistentFlags().StringVarP(&prefix, "prefix", "p", "", "generated file's prefix")
	cmd.PersistentFlags().BoolVar(&includeTest, "include-test", false, "Include test file as target (default: false)")
	cmd.PersistentFlags().StringVarP(&templatePath, "template", "t", "", "generated code's template")
	return cmd
}

func targetFiles(srcFiles []string, includeTestFile bool) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, p := range srcFiles {
			if !includeTestFile && strings.HasSuffix(p, "_test.go") {
				continue
			}
			ch <- p
		}
	}()
	return ch
}
