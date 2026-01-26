package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Println("Connection[UDP] established!")
	defer conn.Close()
	userInput := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := userInput.ReadString('\n')
		if err != nil {
			log.Printf("%v\n", err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			log.Printf("%v\n", err)
		}
	}
}
