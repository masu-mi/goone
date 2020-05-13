package cmd

import (
	"os"
	"path"

	"github.com/masu-mi/goone/model"
	"github.com/spf13/cobra"
)

func NewPackCommand() *cobra.Command {

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
			members := model.ReachableFiles(g, filePath)

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
