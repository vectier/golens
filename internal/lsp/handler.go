package lsp

// All handlers here are derived from the LSP 3.17 specification:
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification
//
// Function name matches the request method name from the spec with slash removed.

func (srv *Server) Initialize(s *Session, _ any) (map[string]any, error) {
	return map[string]any{
		"capabilities": map[string]any{
			// Defines how text documents are synced
			"textDocumentSync": map[string]any{
				"openClose": true,
				"save":      true,
			},
			// Tell the client that the server provides code lens
			"codeLensProvider": map[string]any{
				"resolveProvider": true,
			},
		},
		"serverInfo": map[string]any{
			"name":    "golens",
			"version": Version,
		},
	}, nil
}

func (srv *Server) Initialized(s *Session, _ any) (any, error) {
	s.Callback("client/registerCapability", RegistrationParams{
		Registrations: []Registration{{
			ID:     "golens-watcher",
			Method: "workspace/didChangeWatchedFiles",
			RegisterOptions: DidChangeWatchedFilesRegistrationOptions{
				Watchers: []FileSystemWatcher{
					{GlobPattern: "**/*.go"},
				},
			},
		}},
	})
	return nil, nil
}

func (srv *Server) TextDocumentCodeLens(s *Session, p CodeLensParams) ([]CodeLens, error) {
	return ListInterfaceLenses(p.TextDocument.URI)
}

func (srv *Server) CodeLensResolve(s *Session, p CodeLens) (any, error) {
	return nil, nil
}

func (srv *Server) TextDocumentDidSave(s *Session, p DidSaveTextDocumentParams) (any, error) {
	srv.cache.InvalidatePath(URIToPath(p.TextDocument.URI))
	return nil, nil
}

func (srv *Server) WorkspaceDidChangeWatchedFiles(s *Session, p DidChangeWatchedFilesParams) (any, error) {
	for _, e := range p.Changes {
		srv.cache.InvalidatePath(URIToPath(e.URI))
	}
	return nil, nil
}

func (srv *Server) Shutdown(s *Session, _ any) (any, error) {
	return nil, nil
}

func (srv *Server) Exit(s *Session, _ any) (any, error) {
	return nil, nil
}
