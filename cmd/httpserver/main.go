package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/headers"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/server"
)

const port = 42069

func processVideoExample() server.Handler {
	file, err := os.Open("./assets/vim.mp4")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return func(w *response.Writer, r *request.Request) {
			h := response.GetDefaultHeaders(0)
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(h)
			return
		}
	}
	return func(w *response.Writer, r *request.Request) {
		h := response.GetDefaultHeaders(0)
		if err != nil {
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(h)
			return
		}
		delete(h, "content-length")
		h["transfer-encoding"] = "chunked"
		h["content-type"] = "video/mp4"
		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(h)
		indexOf := 0
		chunk := make([]byte, 1024)
		var bodyBuf []byte
		for {
			n, err := file.Read(chunk)
			indexOf += n
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("%v", err)
				return
			}
			if n > 0 {
				if indexOf > len(chunk) {
					newChunk := make([]byte, len(chunk)*2)
					copy(newChunk, chunk)
					chunk = newChunk
				}
				fmt.Fprintf(w.F, "%x\r\n", n)
				w.WriteChunkedBody(chunk[:n])
				bodyBuf = append(bodyBuf, chunk[:n]...)
				w.F.Write([]byte("\r\n"))
			}
		}
		w.WriteChunkedBodyDone()
		trailers := headers.NewHeaders()
		hasher := sha256.New()
		hasher.Write(bodyBuf)
		hash := fmt.Sprintf("%x", hasher.Sum(nil))
		trailers["X-Content-Sha256"] = hash
		trailers["X-Content-Length"] = fmt.Sprintf("%d", len(bodyBuf))
		w.WriteTrailers(trailers)
		file.Close()
	}

}

func handler(cT string) func(p string) server.Handler {
	return func(p string) server.Handler {
		msg := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`

		statusCode := response.Ok
		if p == "/yourproblem" {
			msg = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
			statusCode = response.BadRequest

		}
		if p == "/myproblem" {
			msg = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
			statusCode = response.InternalServerError
		}

		if strings.HasPrefix(p, "/httpbin") {
			suffix := strings.TrimPrefix(p, "/httpbin")
			return server.Proxy("https://httpbin.org" + suffix)
		}
		msgBytes := []byte(msg)
		h := response.GetDefaultHeaders(len(msgBytes))

		if strings.HasPrefix(p, "/video") {
			return processVideoExample()
		}

		return func(w *response.Writer, r *request.Request) {
			h["content-type"] = cT
			err := w.WriteStatusLine(statusCode)
			if err != nil {
				h = response.GetDefaultHeaders(0)
				w.WriteStatusLine(response.InternalServerError)
				w.WriteHeaders(h)
			}
			err = w.WriteHeaders(h)
			if err != nil {
				h = response.GetDefaultHeaders(0)
				w.WriteStatusLine(response.InternalServerError)
				w.WriteHeaders(h)
			}
			_, err = w.WriteBody(msgBytes)
			if err != nil {
				h = response.GetDefaultHeaders(0)
				w.WriteStatusLine(response.InternalServerError)
				w.WriteHeaders(h)
			}
		}
	}
}

func main() {
	server, err := server.Serve("", port, handler("text/html"))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
