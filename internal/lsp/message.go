package lsp

import (
	"encoding/json"
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
	Range   Range            `json:"range"`
	Command *Command         `json:"command"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

// A data entry field that is preserved on a code lens item
// between a code lens and a code lens resolve request.
//
// This is specific for golens, when the server receives `codeLens/resolve`,
// it reads this field to look up in the cache to fills in `command.title` and `command.arguments`.
type CodeLensData struct {
	URI  string `json:"uri"`
	Key  string `json:"key"`
	Kind string `json:"kind"`
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

type Registration struct {
	ID              string `json:"id"`
	Method          string `json:"method"`
	RegisterOptions any    `json:"registerOptions,omitempty"`
}

type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

type DidChangeWatchedFilesRegistrationOptions struct {
	Watchers []FileSystemWatcher `json:"watchers"`
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
