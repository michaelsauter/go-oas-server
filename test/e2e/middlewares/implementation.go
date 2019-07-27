package middlewares

import (
	"net/http"
)

var calledGlobal bool
var calledLocal bool
var calledHandler bool

type server struct {
}

func (s *server) global(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		calledGlobal = true
		next(w, r)
	}
}

func (s *server) local(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		calledLocal = true
		next(w, r)
	}
}

func (s *server) Middlewares() Middlewares {
	return Middlewares{s.global}
}

func (s *server) HandlePetIndex() (Middlewares, APIOperationPetIndex) {
	return Middlewares{s.local}, func(w http.ResponseWriter, r *http.Request) {
		calledHandler = true
	}
}
