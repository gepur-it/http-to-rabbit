package application

import (
	"net/http"
)

type Request struct {
	*http.Request
	Params []string
}
