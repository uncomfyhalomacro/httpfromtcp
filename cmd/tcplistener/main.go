package main

import (
	"fmt"
	"github.com/uncomfyhalomacro/httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	// f, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		log.Println("Connection established!")
		r, _ := request.RequestFromReader(conn)
		fmt.Printf(`Request line:
- Method: %s
- Target: %s
- Version: %s
`, r.RequestLine.Method, r.RequestLine.RequestTarget, r.RequestLine.HttpVersion)

	}

	// lines := getLinesChannel(f)
	// for line := range lines {
	// 	fmt.Printf("read: %s\n", line)
	// }

}
