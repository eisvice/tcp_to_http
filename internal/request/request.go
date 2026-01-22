package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)


type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	requestMessage, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error while reading request: %s", err)
	}

	requestLine, err := parseRequestLine(requestMessage)
	if err != nil {
		return nil, fmt.Errorf("error while parsing Request-line: %s", err)
	}

	request := Request{RequestLine: *requestLine}

	return &request, nil
}

func parseRequestLine(message []byte) (*RequestLine, error) {
	messageLines := strings.Split(string (message), "\r\n")
	if len(messageLines) < 1 {
		return nil, errors.New("Empty or Invalid Request")
	}

	requestLineParts := strings.Split(messageLines[0], " ")
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