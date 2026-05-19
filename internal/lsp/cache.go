package lsp

import (
	"path/filepath"
	"sync"
)

type ModuleEntry struct {
	Implementations map[string][]Location
	References      map[string][]Location
	CodeLens        map[string][]CodeLens
}

type ModuleCache struct {
	entries map[string]ModuleEntry
	mu      sync.RWMutex
}

func NewModuleCache() *ModuleCache {
	return &ModuleCache{entries: make(map[string]ModuleEntry)}
}

func (c *ModuleCache) Get(path string) (ModuleEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	return e, ok
}

func (c *ModuleCache) Set(path string, e ModuleEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = e
}

func (c *ModuleCache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

func (c *ModuleCache) InvalidatePath(path string) {
	root, err := findModuleRoot(filepath.Dir(path))
	if err != nil {
		return
	}
	c.Invalidate(root)
}

func (srv *Server) getOrSetModuleEntry(fileURI string) (*ModuleEntry, error) {
	filePath := URIToPath(fileURI)
	fileDir := filepath.Dir(filePath)

	// Find the module root so package loading with "./..." covers the entrie module
	moduleRoot, err := findModuleRoot(fileDir)
	if err != nil {
		moduleRoot = fileDir
	}

	if entry, ok := srv.cache.Get(moduleRoot); ok {
		return &entry, nil
	}
	entry, err := LoadModule(moduleRoot)
	if err != nil {
		return nil, err
	}
	srv.cache.Set(moduleRoot, *entry)
	return entry, nil
}
