package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Request struct {
	*http.Request
	Params []string
}

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

	io.WriteString(r, "{\"status\":\"OK\"}")
}

type Handler func(Response, Request)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
}

type App struct {
	Routes       []Route
	DefaultRoute Handler
}

func (a *App) Handle(pattern string, handler Handler) {
	re := regexp.MustCompile(pattern)
	route := Route{Pattern: re, Handler: handler}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{Request: r}
	resp := Response{w}

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
			if len(matches) > 1 {
				req.Params = matches[1:]
			}

			rt.Handler(resp, req)
			return
		}
	}

	a.DefaultRoute(resp, req)
}
