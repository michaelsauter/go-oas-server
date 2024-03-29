// Code generated by go generate; DO NOT EDIT.
package parameters

import (
	"net/http"
)

// OpenAPIServer wraps the custom server.
type OpenAPIServer struct {
	Server      CustomServer
	catIndex    APIEndpointCatIndex
	petIndex    APIEndpointPetIndex
	petShow     APIEndpointPetShow
	middlewares Middlewares
}

// Middleware.
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// Ordered collection of middlewares.
type Middlewares []Middleware

// CustomServer forces that all API operations and middlewares are implemented on the custom server.
type CustomServer interface {
	HandleCatIndex() (Middlewares, APIOperationCatIndex)
	HandlePetIndex() (Middlewares, APIOperationPetIndex)
	HandlePetShow() (Middlewares, APIOperationPetShow)
	Middlewares() Middlewares
}

// NewOpenAPIServer returns an OpenAPIServer with initialised middlewares.
func NewOpenAPIServer(s CustomServer) *OpenAPIServer {
	return &OpenAPIServer{Server: s, middlewares: s.Middlewares()}
}

// ServeHTTP implements the http.Handler interface.
func (oas *OpenAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routingHandler := func(w http.ResponseWriter, r *http.Request) {
		oas.routeLevel0Root(w, r, r.URL.Path)
	}
	// Global middleware
	for i := range oas.middlewares {
		routingHandler = oas.middlewares[len(oas.middlewares)-1-i](routingHandler)
	}
	routingHandler(w, r)
}

// Boot boots the server.
func (oas *OpenAPIServer) Boot() *OpenAPIServer {
	oas.BootCatIndex()
	oas.BootPetIndex()
	oas.BootPetShow()
	return oas
}

// Serve applies middleware for the endpoint.
func (oas *OpenAPIServer) serve(h http.HandlerFunc, m Middlewares, w http.ResponseWriter, r *http.Request) {
	for i := range m {
		h = m[len(m)-1-i](h)
	}
	h(w, r)
}

// And adds all middlewares to the end of the current middlewares.
func (m Middlewares) And(h ...Middleware) Middlewares {
	m = append(m, h...)
	return m
}
