package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/vectier/golens/internal/rpc"
)

// Version is the current version of the language server.
// Set this using ldflags at build time.
var Version = "0.0.0-dev"

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Handle(ctx context.Context, r *bufio.Reader, w io.Writer) {
	for {
		body, err := Parse(r)
		if err != nil {
			log.Printf("malformed jsonrpc2 request: %s\n", err)
			continue
		}
		var req rpc.Request
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf("failed to unmarshal jsonrpc2 request: %s\n", err)
			continue
		}

		switch req.Method {
		case "initialize":
			err = s.Intialize(w, req.ID)
		case "textDocument/codeLens":
			err = s.handleCodeLens(w, &req)
		case "shutdown":
			err = Respond(w, rpc.SuccessResponse(req.ID, nil))
		case "exit":
			return
		}
		if err != nil {
			log.Printf("failed to process request: %s\n", err)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *Server) handleCodeLens(w io.Writer, req *rpc.Request) error {
	var params CodeLensParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return Respond(w, rpc.ErrorResponse(req.ID, rpc.CodeParseError, err))
	}

	lenses, err := ListInterfaceLenses(params.TextDocument.URI)
	if err != nil {
		return Respond(w, rpc.ErrorResponse(req.ID, rpc.CodeInternalError, err))
	}

	return Respond(w, rpc.SuccessResponse(req.ID, lenses))
}

func (s *Server) Intialize(w io.Writer, id int) error {
	return Respond(w, rpc.Response{
		JSONRPC: "2.0",
		ID:      id,
		// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initializeResult
		Result: map[string]any{
			"capabilities": map[string]any{
				// Tell the client that the server provides code lens
				"codeLensProvider": map[string]any{
					"resolveProvider": false,
				},
				// Tell the client to notify the server when any .go file changes.
				// This is to invalidate the CodeLens cache immediately.
				"workspace": map[string]any{
					"fileOperations": map[string]any{},
				},
			},
			"serverInfo": map[string]any{
				"name":    "golens",
				"version": Version,
			},
		},
	})
}
