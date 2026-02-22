package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	Writer io.Writer
	state writerState
}

func (w Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err :=w.Writer.Write(GetStatusLine(statusCode))
	return err
}

func (w Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := fmt.Fprintf(w.Writer, "%s: %s\r\n", k, v)
		if err != nil {
			return err
		}
	}
	_, err := w.Writer.Write([]byte("\r\n"))
	return err
}

func (w Writer) WriteBody(b []byte) error {
	_, err := fmt.Fprintf(w.Writer, "%s", b)
	return err
}