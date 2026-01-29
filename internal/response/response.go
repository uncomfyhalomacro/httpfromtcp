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

type WriterState int

const (
	InitializedWriter WriterState = iota
	StatusLineDone
	HeadersDone
	BodyDone
)

type Writer struct {
	writerState WriterState
	F           io.WriteCloser
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != InitializedWriter {
		return fmt.Errorf("status line should be written first")
	}
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
	_, err := w.F.Write([]byte(statusLine))
	w.writerState = StatusLineDone
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["content-length"] = fmt.Sprintf("%d", contentLen)
	h["connection"] = "close"
	h["content-type"] = "text/plain"
	return h
}

func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.writerState != StatusLineDone {
		return fmt.Errorf("status line should be written first before headers")
	}

	allHeaders := ""
	for key, val := range h {
		headerLine := fmt.Sprintf("%s: %s\r\n", key, val)
		allHeaders += headerLine
	}
	allHeaders += "\r\n"
	_, err := w.F.Write([]byte(allHeaders))
	w.writerState = HeadersDone
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != HeadersDone {
		return 0, fmt.Errorf("headers should be written first before request body")
	}
	n, err := w.F.Write(p)
	w.writerState = BodyDone
	if err != nil {
		return 0, err
	}
	return n, nil
}
