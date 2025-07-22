package tcplistener

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const port = ":8080"

func tcplistener() {
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
	currentLine := ""
	go func() {
		defer close(currentLineChannel)
		for {
			data := make([]byte, 8)
			count, err := f.Read(data)
			if err != nil {
				if currentLine != "" {
					fmt.Printf("read: %s\n", currentLine)
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					break
				}
			}
			dataString := string(data[:count])
			parts := strings.Split(dataString, "\n")

			currentLineChannel <- parts[0]
			currentLine += parts[0]
			if len(parts) > 1 {
				currentLineChannel <- parts[1]
				currentLine = parts[1]
			}
		}

		if currentLine != "" {
			currentLineChannel <- currentLine
		}

	}()
	return currentLineChannel
}
