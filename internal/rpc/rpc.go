package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const Version = "2.0"

// http://www.jsonrpc.org/specification#request_object
type Request struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params,omitempty"`
}

// https://www.jsonrpc.org/specification#response_object
type Response struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

func SuccessResponse(id int, v any) Response {
	return Response{
		JSONRPC: Version,
		ID:      id,
		Result:  v,
	}
}

func ErrorResponse(id int, code int, err error) Response {
	return Response{
		JSONRPC: Version,
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: err.Error(),
		},
	}
}

// http://www.jsonrpc.org/specification#error_object
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", strconv.Itoa(e.Code), e.Message)
}

// http://www.jsonrpc.org/specification#error_object
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603
)
