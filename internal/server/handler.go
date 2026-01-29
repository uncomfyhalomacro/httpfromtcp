package server

import (
	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)
