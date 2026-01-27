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
		log.Println("About to read from conn")
		req, err := request.RequestFromReader(conn)
		log.Println("Finished read from conn")
		if err != nil {
			log.Fatalf("error parsing request: %s\n", err.Error())
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

	}

	// lines := getLinesChannel(f)
	// for line := range lines {
	// 	fmt.Printf("read: %s\n", line)
	// }

}
