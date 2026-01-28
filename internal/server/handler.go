package server

import (
	"io"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
)

type HandlerError struct {
    StatusCode int
    Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
