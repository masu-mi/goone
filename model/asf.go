package model

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
)

type asf struct { // abstruct syntax forest!!
	dir  string
	fst  *token.FileSet
	pkgs map[string]*ast.Package
}

func GetASF(dir string) (a asf, err error) {
	a = asf{
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

func (a asf) ParseAsDefGraph(pkgName string) (*defGraph, error) {
	pkg, ok := a.pkgs[pkgName]
	if !ok {
		return nil, ErrNotExists(pkgName)
	}
	g := NewDefGraph()
	for name, f := range pkg.Files {
		for _, d := range definitions(f) {
			g.addDef(name, d)
		}
		for _, u := range f.Unresolved {
			g.addRef(name, u.Name)
		}
	}
	return g, nil
}

func (a asf) WriteToPackedCode(w io.Writer, pkgName string, members []string) (err error) {
	decls := a.PackedDecls(pkgName, members)
	_, err = fmt.Fprintf(w, "// packed from %v with goone.\n\n", members)
	if err != nil {
		return err
	}
	return format.Node(w, a.fst, decls)
}

func (a asf) PackageFiles(pkgName string) (files []string, err error) {
	pkg, ok := a.pkgs[pkgName]
	if !ok {
		return nil, ErrNotExists(pkgName)
	}
	for name := range pkg.Files {
		files = append(files, name)
	}
	return files, nil
}

func (a asf) PackedDecls(pkgName string, files []string) (output []ast.Decl) {
	for _, m := range files {
		f, ok := a.pkgs[pkgName].Files[m]
		if !ok {
			continue
		}
		appendDecls(&output, f.Decls)
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

func appendDecls(base *[]ast.Decl, items []ast.Decl) {
	b := *base
	for _, d := range items {
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.IMPORT {
				continue
			}
		}
		b = append(b, d)
	}
	*base = b
}
