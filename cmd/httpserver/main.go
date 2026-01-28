package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/server"
)

const port = 42069

func CheckPath(p string) server.Handler {
	msg := "All good, frfr\n"

	statusCode := response.Ok
	if p == "/yourproblem" {
		msg = "Your problem is not my problem\n"
		statusCode = response.BadRequest

	}
	if p == "/myproblem" {
		msg = "Woopsie, my bad\n"
		statusCode = response.InternalServerError
	}

	return func(w io.Writer, r *request.Request) *server.HandlerError {
		msgBytes := []byte(msg)
		h := response.GetDefaultHeaders(len(msgBytes))
		err := response.WriteStatusLine(w, statusCode)
		if err != nil {
			return &server.HandlerError{
				StatusCode: 500,
				Message:    fmt.Sprintf("%v", err),
			}
		}
		err = response.WriteHeaders(w, h)
		if err != nil {
			return &server.HandlerError{
				StatusCode: 500,
				Message:    fmt.Sprintf("%v", err),
			}
		}
		_, err = w.Write(msgBytes)
		if err != nil {
			return &server.HandlerError{
				StatusCode: 500,
				Message:    fmt.Sprintf("%v", err),
			}
		}
		return nil
	}
}

func main() {
	server, err := server.Serve("", port, CheckPath)
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
