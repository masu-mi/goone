package cmd

import (
	"io"
	"os"
	"path"
	"text/template"

	"github.com/masu-mi/goone/model"
	"github.com/spf13/cobra"
)

func NewPackCommand() *cobra.Command {
	var dstPath string
	var templatePath string

	cmd := &cobra.Command{
		Use:   "pack",
		Short: "Pack target source code file with depended files share the package name into single file",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			targetPath, dir, pkgName := targetLocation(c, args)

			cfg := &config{}
			cfg.load(c.Flag("config").Value.String())
			t, err := cfg.NewTemplate(templatePath)
			if err != nil {
				return err
			}

			_, g, err := targetStructure(dir, pkgName)
			if err != nil {
				return err
			}
			return packCode(g, t, dstPath, pkgName, targetPath)
		},
	}
	cmd.PersistentFlags().StringVarP(&dstPath, "out", "o", "", "output file path (default: stdout)")
	cmd.PersistentFlags().String("package", "", "set target package name (default: directory name)")
	cmd.PersistentFlags().StringVarP(&templatePath, "template", "t", "", "generated code's template")
	return cmd
}

func targetLocation(c *cobra.Command, args []string) (filePath, dirName, pkgName string) {
	filePath = args[0]
	dirName = path.Dir(filePath)
	pkgName = c.Flag("package").Value.String()
	if pkgName == "" {
		pkgName = path.Base(dirName)
	}
	return
}

func targetStructure(dir, pkgName string) (a *model.ASF, g *model.DefGraph, err error) {
	a, err = model.GetASF(dir)
	if err != nil {
		return nil, nil, err
	}
	g, err = model.GenDefGraphFromASF(a, pkgName)
	if err != nil {
		return nil, nil, err
	}
	return
}

func packCode(g *model.DefGraph, t *template.Template, dst, pack, src string) (e error) {
	w := os.Stdout
	if dst != "" {
		w, e = os.Create(dst)
		if e != nil {
			return e
		}
	}
	if w != os.Stdout {
		defer w.Close()
	}
	return fprintOutput(w, g, t, pack, src)
}

func fprintOutput(w io.Writer, g *model.DefGraph, t *template.Template, pkgName, src string) error {
	members := model.ReachableFiles(g, src)
	pc, e := g.ASF.PackedCode(pkgName, members)
	if e != nil {
		return e
	}
	if e = ExecutePackedCode(w, t, pc); e != nil {
		return e
	}
	return nil
}
