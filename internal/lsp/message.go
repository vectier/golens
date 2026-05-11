package lsp

import (
	"net/url"
	"strings"
)

// All structs here are derived from the LSP 3.17 specification:
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification

type CodeLensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type CodeLens struct {
	Range   Range    `json:"range"`
	Command *Command `json:"command"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      uint `json:"line"`
	Character uint `json:"character"`
}

type Command struct {
	Title   string `json:"title"`
	Command string `json:"command"`
	Args    []any  `json:"arguments,omitempty"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

func URIToPath(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return strings.TrimPrefix(uri, "file://")
	}
	return u.Path
}

func PathToURI(path string) string {
	return "file://" + path
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

type FileEvent struct {
	URI  string         `json:"uri"`
	Type FileChangeType `json:"type"`
}

type FileChangeType uint

const (
	FileChangeTypeCreated FileChangeType = iota + 1
	FileChangedTypeChanged
	FileChangedTypeDeleted
)

type FileSystemWatcher struct {
	GlobPattern string    `json:"globPattern"`
	Kind        WatchKind `json:"kind,omitempty"`
}

type WatchKind uint

const (
	WatchKindCreate WatchKind = 1
	WatchKindChange WatchKind = 2
	WatchKindDelete WatchKind = 4
)
