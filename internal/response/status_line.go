package response

import "fmt"

type StatusCode int

const (
	StatusOK StatusCode = 200
	StatusBadRequest StatusCode = 400 
	StatusServerError StatusCode = 500
)


func GetStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusServerError:
		reasonPhrase = "Internal Server Error"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}