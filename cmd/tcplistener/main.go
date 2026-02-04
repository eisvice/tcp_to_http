package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

const addr = ":42069"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("error while listening on port %s: %s\n", addr, err)
	}
	defer listener.Close()

	fmt.Printf("listen on port %s\n", addr)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalf("error while accepting a connection: %s\n", err)
		}

		fmt.Println("Connection has been accepted!")
		fmt.Println("<=============================>")
		request, err := request.RequestFromReader(connection)
		if err != nil {
			log.Fatalf("error while reading request: %s", err)
		}
		fmt.Println("Request line:")
		fmt.Println("- Method:", request.RequestLine.Method)
		fmt.Println("- Target:", request.RequestLine.RequestTarget)
		fmt.Println("- Version:", request.RequestLine.HttpVersion)
		fmt.Println("Headers:")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("<=============================>")
		fmt.Println("Connection has been closed!")
	}
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)

// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
// 		currentLineContent := ""

// 		for {
// 			buff := make([]byte, 8) 
// 			n, err := f.Read(buff)
// 			if err != nil {
// 				// first check if the current line contains charachters
// 				if currentLineContent != "" {
// 					// and return lines
// 					lines <- currentLineContent 
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
	
// 			str := string(buff[:n])
// 			parts := strings.Split(str, "\n")
// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- fmt.Sprintf("%s%s", currentLineContent, parts[i])
// 				currentLineContent = ""
// 			}
// 			currentLineContent += parts[len(parts)-1]
// 		}
// 	}()

// 	return lines
// }