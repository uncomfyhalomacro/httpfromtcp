package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/server"
)

const port = 42069

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

		return func(w *response.Writer, r *request.Request) {
			msgBytes := []byte(msg)
			h := response.GetDefaultHeaders(len(msgBytes))
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
