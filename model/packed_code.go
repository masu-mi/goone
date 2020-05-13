package model

// PackedCode is passed value to template
type PackedCode struct {
	// Package is package name (e.g. `main`)
	Package string
	// SrcFiles is used source file name list (e.g. `[src_a.go src_b.go]`)
	SrcFiles string
	// Imports is import declare code of go (e.g. `import "strings"`)
	Imports string
	// Decls is packed code body
	Decls string
}
