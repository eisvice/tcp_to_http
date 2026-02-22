package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Handler func(w *response.Writer, req *request.Request) 

type Server struct {
	inShutdown atomic.Bool
	listener net.Listener
	handler Handler
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
	resW := &response.Writer{
		Writer: conn,
	}
	fmt.Printf("target: %v\n", req.RequestLine.RequestTarget)
	if err != nil {
		errorMessage := []byte(fmt.Sprintf("Error while parsing request: %v", err))
		resW.WriteStatusLine(response.StatusBadRequest)
		resW.WriteHeaders(response.GetDefaultHeaders(len(errorMessage)))
		resW.WriteBody(errorMessage)
	}
	
	s.handler(resW, req)
}
