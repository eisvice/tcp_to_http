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

func NewWriter(w io.Writer) *Writer {
	writer := &Writer{
		Writer: w,
		state: writerStateStatusLine,
	}
	return writer
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

func (w Writer) WriteChunkedBody(p []byte) (int, error) {
	n, err := fmt.Fprintf(w.Writer, "%x\r\n%s\r\n", len(p), p)
	if err != nil {
		return n, fmt.Errorf("error while writing chunk: %v", err)
	}

	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := fmt.Fprintf(w.Writer, "0\r\n\r\n")
	if err != nil {
		return n, fmt.Errorf("error while ending writing body: %v", err)
	}
	return n, nil
}