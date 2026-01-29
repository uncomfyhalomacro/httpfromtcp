package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/uncomfyhalomacro/httpfromtcp/internal/headers"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/response"
)

func Proxy(url string) Handler {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%v", err)
	}
	chunk := make([]byte, 1024)
	return func(w *response.Writer, r *request.Request) {
		h := response.GetDefaultHeaders(0)
		if err != nil {
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(h)
			return
		}
		delete(h, "content-length")
		h["transfer-encoding"] = "chunked"
		w.WriteStatusLine(response.Ok)
		w.WriteHeaders(h)
		indexOf := 0
		var bodyBuf []byte
		for {
			n, err := resp.Body.Read(chunk)
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
		resp.Body.Close()
	}
}
