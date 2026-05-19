package lsp

import (
	"encoding/json"
	"fmt"
)

func (srv *Server) listLenses(fileURI string) ([]CodeLens, error) {
	entry, err := srv.getOrSetModuleEntry(fileURI)
	if err != nil {
		return nil, err
	}
	return entry.CodeLens[URIToPath(fileURI)], nil
}

func (srv *Server) resolveLens(in CodeLens) (*CodeLens, error) {
	if in.Data == nil {
		return nil, fmt.Errorf("missing codelens data")
	}

	var d CodeLensData
	if err := json.Unmarshal(*in.Data, &d); err != nil {
		return nil, err
	}

	entry, err := srv.getOrSetModuleEntry(d.URI)
	if err != nil {
		return nil, err
	}

	var (
		locs []Location
		kind string
	)
	switch d.Kind {
	case "impl":
		locs, kind = entry.Implementations[d.Key], "implementation"
	case "ref":
		locs, kind = entry.References[d.Key], "reference"
	default:
		return nil, fmt.Errorf("unknown codelens kind: %s", d.Key)
	}

	var title string
	if len(locs) == 1 {
		title = fmt.Sprintf("%d %s", len(locs), kind)
	} else {
		title = fmt.Sprintf("%d %ss", len(locs), kind)
	}

	in.Command = &Command{
		Title:   title,
		Command: "editor.action.peekLocations",
		Args:    []any{d.URI, in.Range.Start, locs},
	}
	return &in, nil
}
