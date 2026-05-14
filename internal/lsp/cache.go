package lsp

import (
	"path/filepath"
)

type ModuleEntry struct {
	Implementations map[string][]Location
	References      map[string][]Location
	CodeLens        map[string][]CodeLens
}

type ModuleCache map[string]ModuleEntry

func (c ModuleCache) Get(path string) (ModuleEntry, bool) {
	e, ok := c[path]
	return e, ok
}

func (c ModuleCache) Set(path string, e ModuleEntry) {
	c[path] = e
}

func (c ModuleCache) Invalidate(path string) {
	delete(c, path)
}

func (c ModuleCache) InvalidatePath(path string) {
	root, err := findModuleRoot(filepath.Dir(path))
	if err != nil {
		return
	}
	c.Invalidate(root)
}
