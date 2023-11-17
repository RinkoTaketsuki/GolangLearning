package main

import (
	"net/http"
)

type ServeMux struct {
	*http.ServeMux
	middlewares []Middleware
}

func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: http.DefaultServeMux}
}

func (mux *ServeMux) ApplyMiddlewares(m ...Middleware) {
	mux.middlewares = append(mux.middlewares, m...)
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler) {
	for _, m := range mux.middlewares {
		handler = m(handler)
	}
	mux.ServeMux.Handle(pattern, handler)
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("nil handler")
	}
	mux.Handle(pattern, http.HandlerFunc(handler))
}
