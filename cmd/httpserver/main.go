package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		fmt.Println("your problem")
		return &server.HandlerError{StatusCode: 400, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		fmt.Println("my problem")
		return &server.HandlerError{StatusCode: 500, Message: "Woopsie, my bad\n"}
	}

	fmt.Println("nobody's problem")
	fmt.Fprintf(w, "All good, frfr\n")
	return nil
}