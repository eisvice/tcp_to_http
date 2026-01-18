package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const addr = ":42069"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error while listening on port %s: %s\n", addr, err)
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("error while accepting a connection: %s\n", err)
		}

		fmt.Println("Connection has been accepted!")
		fmt.Println("<=============================>")
		for strReceived := range getLinesChannel(connection) {
			fmt.Println(strReceived)
		}
		fmt.Println("<=============================>")
		fmt.Println("Connection has been closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContent := ""

		for {
			buff := make([]byte, 8) 
			n, err := f.Read(buff)
			if err != nil {
				// first check if the current line contains charachters
				if currentLineContent != "" {
					// and return lines
					lines <- currentLineContent 
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
	
			str := string(buff[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContent, parts[i])
				currentLineContent = ""
			}
			currentLineContent += parts[len(parts)-1]
		}
	}()

	return lines
}