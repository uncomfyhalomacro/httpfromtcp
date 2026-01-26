package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		data := make([]byte, 8)
		currentLine := ""
		defer f.Close()
		defer close(ch)
		for {
			n, err := f.Read(data)
			if err == io.EOF {
				break
			}
			start := 0
			for i := 0; i < n; i++ {
				if data[i] == '\n' {
					part := string(data[start:i])
					currentLine += part
					ch <- currentLine
					currentLine = ""
					start = i + 1
				}
			}
			lastPart := string(data[start:n])
			currentLine += lastPart
		}
		if len(currentLine) > 0 {
			ch <- currentLine
		}
	}()
	return ch
}

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
		defer conn.Close()
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Printf("read: %s\n", line)
		}
	}

	// lines := getLinesChannel(f)
	// for line := range lines {
	// 	fmt.Printf("read: %s\n", line)
	// }

}
