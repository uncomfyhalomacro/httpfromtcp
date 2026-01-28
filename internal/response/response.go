package response

import (
	"fmt"
	"io"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Ok                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := ""
	switch statusCode {
	case Ok:
		statusLine = "HTTP/1.1 200 OK"
	case BadRequest:
		statusLine = "HTTP/1.1 400 Bad Request"
	case InternalServerError:
		statusLine = "HTTP/1.1 500 Internal Server Error"
	default:
		statusLine = ""
	}
	statusLine += "\r\n"
	_, err := w.Write([]byte(statusLine))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = fmt.Sprintf("%d", contentLen) 
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func WriteHeaders(w io.Writer, h headers.Headers) error {
	allHeaders := ""
	for key, val := range h {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, val)
		allHeaders += headerLine
	}
	allHeaders += "\r\n"
	_, err := w.Write([]byte(allHeaders))
	return err
}
