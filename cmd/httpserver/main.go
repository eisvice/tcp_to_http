package main

import (
	"crypto/sha256"
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

func handler(w *response.Writer, req *request.Request) {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		handlerHttpbin(w, req)
		return
	}
	handler200(w, req)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
}

func handlerHttpbin(w *response.Writer, req *request.Request) {
	contentSha := strings.ToLower("X-Content-Sha256")
	contentLength := strings.ToLower("X-Content-Length")

	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	fmt.Printf("target: %s\n", target)
	url := fmt.Sprintf("https://httpbin.org/%s", target)
	fmt.Printf("Proxying to %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error in response from the api server: %v", err)
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Override("Transfer-Encoding", "chunked")
	h.Set("Trailer", contentSha)
	h.Set("Trailer", contentLength)
	w.WriteHeaders(h)
	
	fullResponse := []byte{}
	length := 0
	
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		fmt.Printf("read to buffer %d bytes\n", n)

		if n > 0 {
			_, err := w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Printf("error while writing cnunk: %v\n", err)
				break
			}
			fullResponse = append(fullResponse, buf[:n]...)
			length += n
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("error reading chunk: %v\n", err)
			break
		}

	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Printf("error while finishing writing chunk: %v", err)
	}

	trailers := headers.NewHeaders()
	trailers.Set(contentSha, fmt.Sprintf("%x", sha256.Sum256(fullResponse)))
	trailers.Set(contentLength, fmt.Sprintf("%d", length))
	err = w.WriteTrailers(trailers)
	if err != nil {
		fmt.Printf("error while writing trailers: %v\n", err)
	}
}
