package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"
const colon = ":"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	headerName, headerValue, found := bytes.Cut(data[:idx], []byte(":"))
	if !found {
		return idx+2, false, fmt.Errorf("invalid header: %s", data[:idx])
	}

	key := string(headerName)
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	key = strings.TrimSpace(key)
	value := strings.TrimSpace(string(headerValue))

	h[key] = value
	return idx + 2, false, nil
}

