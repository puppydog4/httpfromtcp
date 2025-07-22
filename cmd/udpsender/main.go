package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpEndpoint, err := net.ResolveUDPAddr("udp", "localhost:8080")

	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", nil, udpEndpoint)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		println("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading input:", err.Error())
		}

		conn.Write([]byte(text))
		fmt.Printf("Sent: %s", text)
	}

}
