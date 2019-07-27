// Code generated by go generate; DO NOT EDIT.
package middlewares

import (
	"net/http"
)

// APIEndpointPetIndex describes the petIndex endpoint.
type APIEndpointPetIndex struct {
	handler     APIOperationPetIndex
	middlewares Middlewares
}

// BootPetIndex boots the petIndex endpoint.
func (oas *OpenAPIServer) BootPetIndex() {
	m, h := oas.Server.HandlePetIndex()
	oas.petIndex = APIEndpointPetIndex{
		handler:     h,
		middlewares: m,
	}
}

// APIOperationPetIndex is an alias for the func signature.
type APIOperationPetIndex func(w http.ResponseWriter, r *http.Request)

func (e APIEndpointPetIndex) execute(w http.ResponseWriter, r *http.Request) {

	// Call handler
	e.handler(w, r)
}
