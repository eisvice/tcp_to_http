package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError 

type Server struct {
	inShutdown atomic.Bool
	listener net.Listener
	handler Handler
}

func (he HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, he.StatusCode)
	messageBytes := []byte(he.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(w, headers)
	w.Write(messageBytes)
}


func Serve(port int, handler Handler) (*Server, error) {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: l,
		handler: handler,
	}
	go s.listen()
	return  s, nil
}

func (s *Server) Close() error {
	s.inShutdown.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {

	for {
		// Wait for a connection.
		conn, err := s.listener.Accept()
		if err != nil {
			if s.inShutdown.Load() {
				return
			}
			log.Printf("Error accepting connection: %s", err)
			continue
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go s.handle(conn) 
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	fmt.Printf("target: %v\n", req.RequestLine.RequestTarget)
	if err != nil {
		handlerError := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Message: err.Error(),
		}
		handlerError.Write(conn)
	}
	
	buf := bytes.NewBuffer([]byte{})
	handlerError := s.handler(buf, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}


	if err := response.WriteStatusLine(conn, response.StatusOK); err != nil {
		fmt.Printf("%v", err)
	}
	headers := response.GetDefaultHeaders(buf.Len())
	if err := response.WriteHeaders(conn, headers); err != nil {
		fmt.Printf("error in handle: %v\n", err)
	}
	if err := response.WriteBody(conn, buf.Bytes()); err != nil {
		fmt.Printf("%v", err)
	}
}
