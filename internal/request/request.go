package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)


type Request struct {
	RequestLine RequestLine
	Headers headers.Headers
	state requestState
}


type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialied requestState = iota
	requestStateParsingHeaders 
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize)
	readToIndex := 0
	request := Request{
		state: requestStateInitialied,
	}

	for request.state != requestStateDone {
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff) * 2)
			copy(newBuff, buff)
			buff = newBuff
		}

		numBytesRead, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", request.state, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := request.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[numBytesParsed:])
		readToIndex -= numBytesParsed
		fmt.Println(request)
	}

	return &request, nil
}

func parseRequestLine(message []byte) (*RequestLine, int, error) {
	idx := bytes.Index(message, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineText := string(message[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	requestLineParts := strings.Split(str, " ")
	if len(requestLineParts) != 3 {
		return nil, errors.New("Invalid Request-Line")
	}

	method := requestLineParts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("Invalid method: %s", method)
		}
	}

	versionParts := strings.Split(requestLineParts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("Invalid HTTP part of the Request-Line: %s", requestLineParts[2])
	}

	httpPart := versionParts[0]
	version := versionParts[1]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("Unrecognized HTTP-version: %s", httpPart)
	}
	if version != "1.1" {
		return nil, fmt.Errorf("Unrecognized HTTP-version: %s", httpPart)
	}

	requestLine := RequestLine{
		HttpVersion: version,
		Method: method,
		RequestTarget: requestLineParts[1],
	}
	return &requestLine, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialied:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.Headers = headers.NewHeaders()
		r.state = requestStateParsingHeaders
		return n, nil

	case requestStateParsingHeaders:
		done := false
		totalBytesParsed := 0
		
		for !done {
			n, done, err := r.Headers.Parse(data[totalBytesParsed:])
			totalBytesParsed += n
			if err != nil {
				return 0, err
			}
			if done {
				r.state = requestStateDone
			} 
			if n == 0 {
				break
			}
		}

		return totalBytesParsed, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}	
}
