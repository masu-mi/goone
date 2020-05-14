package cmd

import (
	"fmt"
	"io"
	"path"
	"text/template"

	"github.com/masu-mi/goone/model"
)

func outputFilePath(dir, prefix, srcFile string) string {
	baseName := path.Clean(path.Base(srcFile))
	return path.Join(dir, fmt.Sprintf("%s%s", prefix, baseName))
}

func ExecutePackedCode(w io.Writer, t *template.Template, code model.PackedCode) error {
	return t.Execute(w, code)
}
