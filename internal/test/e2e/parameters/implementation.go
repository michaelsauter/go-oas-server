package parameters

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var receivedFoo string
var receivedStatus int
var receivedBar int
var receivedID int
var receivedClientID uuid.UUID
var receivedClientTime time.Time

type server struct {
}

func (s *server) Middlewares() Middlewares {
	return Middlewares{}
}

func (s *server) HandlePetIndex() (Middlewares, APIOperationPetIndex) {
	return nil, func(w http.ResponseWriter, r *http.Request, p ParametersPetIndex) {
		receivedStatus = p.Status
		receivedFoo = p.Foo
		receivedBar = p.Bar
	}
}

func (s *server) HandleCatIndex() (Middlewares, APIOperationCatIndex) {
	return nil, func(w http.ResponseWriter, r *http.Request, p ParametersCatIndex) {
		receivedClientID = p.X_Client_ID
		receivedClientTime = p.X_Client_Time
	}
}

func (s *server) HandlePetShow() (Middlewares, APIOperationPetShow) {
	return nil, func(w http.ResponseWriter, r *http.Request, p ParametersPetShow) {
		receivedID = p.Id
	}
}
