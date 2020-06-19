package model

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

type ASF struct { // abstruct syntax forest!!
	dir  string
	fst  *token.FileSet
	pkgs map[string]*ast.Package
}

func GetASF(dir string) (a *ASF, err error) {
	a = &ASF{
		dir: dir,
		fst: token.NewFileSet(),
	}
	a.pkgs, err = parser.ParseDir(a.fst, dir, nil, 0)
	if err != nil {
		return a, err
	}
	return a, err
}

func ErrNotExists(pkgName string) error {
	return fmt.Errorf("package %s don't exist", pkgName)
}

func (a *ASF) PackedCode(pkgName string, members []string) (pc PackedCode, err error) {
	decls, imports := a.PackedDecls(pkgName, members)
	var bd, bi strings.Builder
	err = format.Node(&bd, a.fst, decls)
	if err != nil {
		return pc, err
	}
	err = format.Node(&bi, a.fst, imports)
	if err != nil {
		return pc, err
	}
	return PackedCode{
		Package:  pkgName,
		SrcFiles: fmt.Sprintf("%v", members),
		Imports:  bi.String(),
		Decls:    bd.String(),
	}, nil
}

func (a *ASF) PackageFiles(pkgName string) (files []string, err error) {
	pkg, ok := a.pkgs[pkgName]
	if !ok {
		return nil, ErrNotExists(pkgName)
	}
	for name := range pkg.Files {
		files = append(files, name)
	}
	return files, nil
}

func (a *ASF) PackedDecls(pkgName string, files []string) (defs []ast.Decl, imports *ast.GenDecl) {
	imports = &ast.GenDecl{
		Tok: token.IMPORT,
	}
	for _, m := range files {
		f, ok := a.pkgs[pkgName].Files[m]
		if !ok {
			continue
		}
		appendDecls(&defs, imports, f.Decls)
	}
	return
}

func definitions(f *ast.File) (defs []definition) {
	for _, d := range f.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			if v.Recv != nil {
				continue
			}
			defs = append(defs, definition{funcDecl, v.Name.Name})
		case *ast.GenDecl:
			for _, s := range v.Specs {
				switch sv := s.(type) {
				case *ast.TypeSpec:
					defs = append(defs, definition{typeDecl, sv.Name.Name})
				case *ast.ValueSpec:
					for _, i := range sv.Names {
						defs = append(defs, definition{valueDecl, i.Name})
					}
				default:
				}
			}
		default:
		}
	}
	return defs
}

func appendDecls(base *[]ast.Decl, imports *ast.GenDecl, items []ast.Decl) {
	b := *base
	defer func() {
		*base = b
	}()
	appended := map[string]struct{}{}
	for _, d := range items {
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.IMPORT {
				for _, i := range gd.Specs {
					if is, ok := i.(*ast.ImportSpec); ok {
						importPath := is.Path.Value
						if _, ok := appended[importPath]; !ok {
							imports.Specs = append(imports.Specs, i)
						}
						appended[importPath] = struct{}{}
					}
				}
				continue
			}
		}
		b = append(b, d)
	}
}
