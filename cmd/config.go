package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/BurntSushi/toml"
)

const cmdName = "goone"

type config struct {
	confPath     string
	TemplateFile string `toml:"templatefile"`
}

const defaultTemplate = `
{{"{{_cursor_}}"}}
// package: {{ .Package }}
// packed src of {{ .SrcFiles }} with goone.

{{ .Decls }}
`

func (cfg *config) load(confPath string) error {
	var dir string
	file := confPath
	if file == "" {
		if runtime.GOOS == "windows" {
			dir = os.Getenv("APPDATA")
			if dir == "" {
				dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", cmdName)
			}
			dir = filepath.Join(dir, cmdName)
		} else {
			dir = filepath.Join(os.Getenv("HOME"), ".config", cmdName)
		}
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("cannot create directory: %v", err)
		}
		file = filepath.Join(dir, "config.toml")
	}
	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		cfg.confPath = file
	} else if confPath == "" {
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		toml.NewEncoder(f).Encode(cfg)
	}
	return nil
}

func (cfg *config) NewTemplate(p string) (*template.Template, error) {
	templateString := defaultTemplate
	if p == "" {
		p = cfg.TemplateFile
	}
	path := cfg.expandFilePath(p)
	if fileExists(path) {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		templateString = string(b)
	}
	return template.New("PackedCode").Parse(templateString)
}

func (cfg *config) expandFilePath(p string) string {
	if p == "" {
		return ""
	}
	// absolute, relative-from-current, in templates dir under conf dir
	if fileExists(p) {
		return p
	}
	cp := path.Clean(p)
	return path.Join(path.Dir(cfg.confPath), "templates", cp)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
