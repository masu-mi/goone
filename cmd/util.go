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

const defautTemplate = `package {{ .Package }}
// packed from {{ .SrcFiles }} with goone.

// {{"{{_cursor_}}"}}

{{ .Imports }}

{{ .Decls }}
`

var dt, _ = template.New("PackedCode").Parse(defautTemplate)

func ExecutePackedCode(w io.Writer, t *template.Template, code model.PackedCode) error {
	return t.Execute(w, code)
}
