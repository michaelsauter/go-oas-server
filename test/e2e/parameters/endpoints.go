// Code generated by go generate; DO NOT EDIT.
package parameters

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// APIEndpointCatIndex describes the catIndex endpoint.
type APIEndpointCatIndex struct {
	handler     APIOperationCatIndex
	middlewares Middlewares
}

// ParametersCatIndex describes the parameters for the catIndex endpoint.
type ParametersCatIndex struct {
	X_Client_ID   uuid.UUID
	X_Client_Time time.Time
}

// BootCatIndex boots the catIndex endpoint.
func (oas *OpenAPIServer) BootCatIndex() {
	m, h := oas.Server.HandleCatIndex()
	oas.catIndex = APIEndpointCatIndex{
		handler:     h,
		middlewares: m,
	}
}

// APIOperationCatIndex is an alias for the func signature.
type APIOperationCatIndex func(w http.ResponseWriter, r *http.Request, p ParametersCatIndex)

func (e APIEndpointCatIndex) execute(w http.ResponseWriter, r *http.Request) {
	p := ParametersCatIndex{}
	// Header Parameters
	var rawHeaderValue string
	rawHeaderValue = r.Header.Get("X-Client-ID")
	if len(rawHeaderValue) == 0 {
		http.Error(w, "Header 'X-Client-ID' required", http.StatusBadRequest)
		return
	}
	if len(rawHeaderValue) > 0 {
		u, err := uuid.Parse(rawHeaderValue)
		if err != nil {
			http.Error(w, "Header 'X-Client-ID' is not a UUID", http.StatusBadRequest)
			return
		}
		p.X_Client_ID = u
	}

	rawHeaderValue = r.Header.Get("X-Client-Time")
	if len(rawHeaderValue) > 0 {
		t, err := time.Parse(time.RFC3339, rawHeaderValue)
		if err != nil {
			http.Error(w, "Header 'X-Client-Time' is not a date-time", http.StatusBadRequest)
			return
		}
		p.X_Client_Time = t
	}

	// Call handler
	e.handler(w, r, p)
}

// APIEndpointPetIndex describes the petIndex endpoint.
type APIEndpointPetIndex struct {
	handler     APIOperationPetIndex
	middlewares Middlewares
}

// ParametersPetIndex describes the parameters for the petIndex endpoint.
type ParametersPetIndex struct {
	Status int
	Foo    string
	Bar    int
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
type APIOperationPetIndex func(w http.ResponseWriter, r *http.Request, p ParametersPetIndex)

func (e APIEndpointPetIndex) execute(w http.ResponseWriter, r *http.Request) {
	p := ParametersPetIndex{}
	// Query parameters
	var rawQueryValue string
	// Handle status param
	rawQueryValue = r.URL.Query().Get("status")
	if len(rawQueryValue) > 0 {
		intValue, err := strconv.Atoi(rawQueryValue)
		if err != nil {
			http.Error(w, "Query 'status' is not an integer", http.StatusBadRequest)
			return
		}
		p.Status = intValue
	}

	// Handle foo param
	rawQueryValue = r.URL.Query().Get("foo")
	if len(rawQueryValue) == 0 {
		http.Error(w, "Query 'foo' required", http.StatusBadRequest)
		return
	}
	if len(rawQueryValue) > 0 {
		p.Foo = rawQueryValue
	}

	// Handle bar param
	rawQueryValue = r.URL.Query().Get("bar")
	p.Bar = 10
	if len(rawQueryValue) > 0 {
		intValue, err := strconv.Atoi(rawQueryValue)
		if err != nil {
			http.Error(w, "Query 'bar' is not an integer", http.StatusBadRequest)
			return
		}
		p.Bar = intValue
	}

	// Call handler
	e.handler(w, r, p)
}

// APIEndpointPetShow describes the petShow endpoint.
type APIEndpointPetShow struct {
	handler     APIOperationPetShow
	middlewares Middlewares
}

// ParametersPetShow describes the parameters for the petShow endpoint.
type ParametersPetShow struct {
	Id int
}

// BootPetShow boots the petShow endpoint.
func (oas *OpenAPIServer) BootPetShow() {
	m, h := oas.Server.HandlePetShow()
	oas.petShow = APIEndpointPetShow{
		handler:     h,
		middlewares: m,
	}
}

// APIOperationPetShow is an alias for the func signature.
type APIOperationPetShow func(w http.ResponseWriter, r *http.Request, p ParametersPetShow)

func (e APIEndpointPetShow) execute(w http.ResponseWriter, r *http.Request) {
	p := ParametersPetShow{}
	// Path parameters
	oasPathParts := strings.Split("/pets/{id}", "/")
	pathParts := strings.Split(r.URL.Path, "/")
	for k, v := range pathParts {
		if oasPathParts[k] == "{id}" {
			rawValue, err := strconv.Atoi(v)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid user id %q", v), http.StatusBadRequest)
				return
			}
			p.Id = rawValue
		}
	}

	// Call handler
	e.handler(w, r, p)
}
