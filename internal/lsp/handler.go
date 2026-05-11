package lsp

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initialize
func (srv *Server) Initialize(s *Session, _ *any) (map[string]any, error) {
	result := map[string]any{
		"capabilities": map[string]any{
			// Tell the client that the server provides code lens
			"codeLensProvider": map[string]any{
				"resolveProvider": false,
			},
		},
		"serverInfo": map[string]any{
			"name":    "golens",
			"version": Version,
		},
	}
	return result, nil
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initialized
func (srv *Server) Initialized(s *Session, _ *any) (any, error) {
	return nil, nil
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_codeLens
func (srv *Server) TextDocumentCodeLens(s *Session, p *CodeLensParams) ([]CodeLens, error) {
	return ListInterfaceLenses(p.TextDocument.URI)
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#shutdown
func (srv *Server) Shutdown(s *Session, _ *any) (any, error) {
	return nil, nil
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#exit
func (srv *Server) Exit(s *Session, _ *any) (any, error) {
	return nil, nil
}
