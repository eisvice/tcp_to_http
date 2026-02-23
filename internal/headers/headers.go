package headers

import (
	"bytes"
	"fmt"
	"slices"
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

	if !validTokens([]byte(key)) {
		return 0, false, fmt.Errorf("invalid characters in header: %s", key)
	}

	h.Set(key, value)
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	v, ok := h[key]
	if ok {
		value = strings.Join([]string{v, value}, ", ")
	}
	h[key] = value
}

func (h Headers) Get(key string) (string, bool) {
	v, ok := h[strings.ToLower(key)]
	return v, ok
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

// validTokens checks if the data contains only valid tokens
// or characters that are allowed in a token
func validTokens(data []byte) bool {
	for _, c := range data {
		if !isTokenChar(c) {
			return false
		}
	}
	return true
}

func isTokenChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}