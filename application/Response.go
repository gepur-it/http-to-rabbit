package application

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	http.ResponseWriter
}

func (r *Response) Text(code int, body string) {
	r.Header().Set("Content-Type", "text/plain")
	r.WriteHeader(code)

	io.WriteString(r, fmt.Sprintf("\"%s\"", body))
}

func (r *Response) Success() {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusOK)

	io.WriteString(r, "{\"success\":true}")
}

func (r *Response) Error(error string) {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusInternalServerError)

	errorStr := strings.Replace(error, `"`, `\"`, -1)
	io.WriteString(r, fmt.Sprintf("{\"success\":false,\"error\":\"%s\"}", errorStr))
}
