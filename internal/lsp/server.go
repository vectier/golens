package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"github.com/vectier/golens/internal/rpc"
)

// Version is the current version of the language server.
// Set this using ldflags at build time.
var Version = "0.0.0-dev"

type Server struct {
	handlers HandlerRegistry
	cache    *ModuleCache
	nextID   atomic.Int64
}

func NewServer() *Server {
	srv := &Server{
		handlers: make(HandlerRegistry),
		cache:    NewModuleCache(),
	}

	RegisterHandler(srv, "initialize", srv.Initialize)
	RegisterHandler(srv, "initialized", srv.Initialized)
	RegisterHandler(srv, "textDocument/codeLens", srv.TextDocumentCodeLens)
	RegisterHandler(srv, "codeLens/resolve", srv.CodeLensResolve)
	RegisterHandler(srv, "textDocument/didSave", srv.TextDocumentDidSave)
	RegisterHandler(srv, "workspace/didChangeWatchedFiles", srv.WorkspaceDidChangeWatchedFiles)
	RegisterHandler(srv, "shutdown", srv.Shutdown)
	RegisterHandler(srv, "exit", srv.Exit)

	return srv
}

func (srv *Server) Handle(ctx context.Context, r *bufio.Reader, w io.Writer) {
	for {
		s, err := NewSession(srv, r, w)
		if err != nil {
			log.Printf("failed to create a session: %s\n", err)
			continue
		}

		if err := srv.Dispatch(s); err != nil {
			log.Printf("failed to process request: %s\n", err)
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (srv *Server) Dispatch(s *Session) error {
	handler, ok := srv.handlers[s.Method]
	if !ok {
		return Respond(
			s.responseWriter,
			rpc.ErrorResponse(
				s.ID,
				rpc.CodeMethodNotFound,
				fmt.Errorf("unknown method: %s", s.Method),
			),
		)
	}
	return handler(s)
}

type Handler[T any, R any] func(*Session, T) (R, error)
type HandlerRegistry map[string]func(*Session) error

func RegisterHandler[T any, R any](
	s *Server,
	method string,
	fn Handler[T, R],
) {
	s.handlers[method] = func(s *Session) error {
		var param T
		if s.Params != nil {
			if err := json.Unmarshal(*s.Params, &param); err != nil {
				return Respond(s.responseWriter, rpc.ErrorResponse(s.ID, rpc.CodeParseError, err))
			}
		}
		v, err := fn(s, param)
		if err != nil {
			// TODO: unwrap error to get a correct rpc error code
			return Respond(s.responseWriter, rpc.ErrorResponse(s.ID, rpc.CodeInternalError, err))
		}
		return Respond(s.responseWriter, rpc.SuccessResponse(s.ID, v))
	}
}

type Session struct {
	srv            *Server
	responseWriter io.Writer
	rpc.Request
}

func NewSession(srv *Server, r *bufio.Reader, w io.Writer) (*Session, error) {
	body, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("malformed jsonrpc2 request: %w", err)
	}
	var req rpc.Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("failed to unmarshal jsonrpc2 request: %w", err)
	}
	return &Session{srv: srv, responseWriter: w, Request: req}, nil
}

func (s *Session) Callback(method string, params any) error {
	id := int(s.srv.nextID.Add(1))
	req, err := rpc.NewRequest(id, method, params)
	if err != nil {
		return err
	}
	return Respond(s.responseWriter, req)
}
 