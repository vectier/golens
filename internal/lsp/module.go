package lsp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

func LoadModule(root string) (*ModuleEntry, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		Dir:  root,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}

	types := ListTypes(pkgs)
	implementations := ListImplementations(types)
	references := ListReferences(pkgs)

	lenses := make(map[string][]CodeLens)
	for _, pkg := range pkgs {
		for path, ls := range ListCodeLenses(pkg) {
			lenses[path] = append(lenses[path], ls...)
		}
	}

	return &ModuleEntry{
		Implementations: implementations,
		References:      references,
		CodeLens:        lenses,
	}, nil
}

type NamedWithLocation struct {
	Named    *types.Named
	Location Location
}

func ListTypes(pkgs []*packages.Package) map[string]NamedWithLocation {
	allTypes := make(map[string]NamedWithLocation)
	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}

		scope := pkg.Types.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)

			tn, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}
			named, ok := tn.Type().(*types.Named)
			if !ok {
				continue
			}

			allTypes[typeObjectKey(obj)] = NamedWithLocation{
				Named: named,
				Location: Location{
					URI:   PathToURI(pkg.Fset.Position(obj.Pos()).Filename),
					Range: posRange(pkg.Fset, obj.Pos(), obj.Pos()+token.Pos(len(named.Obj().Name()))),
				},
			}
		}
	}
	return allTypes
}

func ListImplementations(allTypes map[string]NamedWithLocation) map[string][]Location {
	implementations := make(map[string][]Location)
	for key, t := range allTypes {
		iface, ok := t.Named.Underlying().(*types.Interface)
		if !ok || iface.NumMethods() == 0 {
			continue
		}

		for k, v := range allTypes {
			// Skip the interface itself
			if k == key {
				continue
			}
			if types.Implements(v.Named, iface) || types.Implements(types.NewPointer(t.Named), iface) {
				implementations[key] = append(implementations[key], v.Location)
			}
		}
	}
	return implementations
}

func ListReferences(pkgs []*packages.Package) map[string][]Location {
	references := make(map[string][]Location)
	for _, pkg := range pkgs {
		if pkg.Types == nil {
			continue
		}

		for ident, obj := range pkg.TypesInfo.Uses {
			tn, ok := obj.(*types.TypeName)
			if !ok {
				continue
			}
			named, ok := tn.Type().(*types.Named)
			if !ok {
				continue
			}
			if _, ok := named.Underlying().(*types.Interface); !ok {
				continue
			}

			key := typeObjectKey(obj)
			references[key] = append(references[key], Location{
				URI:   PathToURI(pkg.Fset.Position(ident.Pos()).Filename),
				Range: posRange(pkg.Fset, ident.Pos(), ident.End()),
			})
		}
	}
	return references
}

func ListCodeLenses(pkg *packages.Package) map[string][]CodeLens {
	lenses := make(map[string][]CodeLens)
	if pkg.TypesInfo == nil {
		return lenses
	}

	for _, file := range pkg.Syntax {
		tf := pkg.Fset.File(file.Pos())
		if tf == nil {
			continue
		}

		path := tf.Name()
		uri := PathToURI(path)
		ast.Inspect(file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			// Only process non-empty interfaces
			ifaceType, ok := ts.Type.(*ast.InterfaceType)
			if !ok || ifaceType.Methods == nil || len(ifaceType.Methods.List) == 0 {
				return true
			}

			obj := pkg.TypesInfo.Defs[ts.Name]
			if obj == nil {
				return true
			}

			r := posRange(pkg.Fset, ts.Name.Pos(), ts.Name.End())
			key := typeObjectKey(obj)
			lenses[path] = append(lenses[path],
				CodeLensWithData(r, uri, key, "impl"),
				CodeLensWithData(r, uri, key, "ref"),
			)
			return true
		})
	}
	return lenses
}

func findModuleRoot(dir string) (string, error) {
	current := dir
	for {
		if _, err := os.Stat(filepath.Join(current, "go.mod")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return dir, fmt.Errorf("go.mod not found, falling back to dir")
		}
		current = parent
	}
}

func typeObjectKey(obj types.Object) string {
	if obj.Pkg() == nil {
		return obj.Name()
	}
	return obj.Pkg().Path() + "." + obj.Name()
}

func posRange(fset *token.FileSet, start token.Pos, end token.Pos) Range {
	s := fset.Position(start)
	e := fset.Position(end)
	return Range{
		Start: Position{
			Line:      uint(s.Line - 1),
			Character: uint(s.Column - 1),
		},
		End: Position{
			Line:      uint(e.Line - 1),
			Character: uint(e.Column - 1),
		},
	}
}
