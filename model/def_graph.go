package model

import "path"

type DefGraph struct {
	defGraph

	Package string
	ASF     *ASF
}

type defGraph struct {
	refs map[string][]string       // file -> ident
	defs map[string]map[int]string // ident -> file
}

func NewDefGraph() *DefGraph {
	return &DefGraph{
		defGraph: defGraph{
			refs: make(map[string][]string),
			defs: make(map[string]map[int]string),
		},
	}
}

func GenDefGraphFromASF(a *ASF, pkgName string) (*DefGraph, error) {
	pkg, ok := a.pkgs[pkgName]
	if !ok {
		return nil, ErrNotExists(pkgName)
	}
	g := NewDefGraph()
	g.ASF = a
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

func (g *DefGraph) addDef(file string, ident definition) {
	m, ok := g.defs[ident.name]
	if !ok {
		m = map[int]string{}
	}
	m[ident.t] = file
	g.defs[ident.name] = m
}

func (g *DefGraph) addRef(file, ident string) {
	g.refs[file] = append(g.refs[file], ident)
}

func (g *DefGraph) dependedFiles(file string) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, id := range g.refs[file] {
			defs, ok := g.defs[id]
			if !ok { // keywords, package name or ...
				continue
			}
			for _, name := range defs {
				ch <- name
			}
		}
	}()
	return ch
}

func ReachableFiles(g *DefGraph, start string) (member []string) {
	start = path.Clean(start)
	visited := newSet()
	visited.add(start)
	_dfs(visited, g, start)
	return visited.members()
}
func _dfs(visited set, g *DefGraph, cur string) {
	for f := range g.dependedFiles(cur) {
		if visited.doesContain(f) {
			continue
		}
		visited.add(f)
		_dfs(visited, g, f)
	}
	return
}

type definition struct {
	t    int
	name string
}

const (
	funcDecl = 0 + iota
	typeDecl
	valueDecl
)

type none struct{}

var mark none

type set map[string]none

func newSet() set {
	return make(map[string]none)
}

func (s set) add(item string) {
	s[item] = mark
}

func (s set) doesContain(item string) bool {
	_, ok := s[item]
	return ok
}

func (s set) size() int {
	return len(s)
}

func (s set) members() (l []string) {
	for k := range s {
		l = append(l, k)
	}
	return l
}
