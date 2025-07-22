package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":8080"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("New connection from", port)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection:", err)
		}

		getLinesChannel := getLinesChannel(connection)
		for line := range getLinesChannel {
			println(line)
		}
		fmt.Println("Connection to ", connection.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	currentLineChannel := make(chan string)

	go func() {
		defer f.Close()
		defer close(currentLineChannel)
		currentLine := ""
		for {
			data := make([]byte, 8)
			count, err := f.Read(data)
			if err != nil {
				if currentLine != "" {
					currentLineChannel <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
			}
			dataString := string(data[:count])
			parts := strings.Split(dataString, "\n")

			currentLine += parts[0]
			if len(parts) > 1 {
				currentLineChannel <- currentLine
				currentLine = parts[1]
			}
		}

		if currentLine != "" {
			currentLineChannel <- currentLine
		}

	}()
	return currentLineChannel
}
