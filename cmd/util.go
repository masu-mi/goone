package cmd

import (
	"fmt"
	"path"
)

func outputFilePath(dir, prefix, srcFile string) string {
	baseName := path.Clean(path.Base(srcFile))
	return path.Join(dir, fmt.Sprintf("%s%s", prefix, baseName))
}
