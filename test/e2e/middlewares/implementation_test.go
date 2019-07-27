package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerMiddlewares(t *testing.T) {
	s := NewOpenAPIServer(&server{})
	s.BootPetIndex()
	ts := httptest.NewServer(s)
	tests := map[string]struct {
		URL             string
		expectedStatus  int
		expectedLocal   bool
		expectedHandler bool
	}{
		"existing URL": {
			URL:             ts.URL + "/pets",
			expectedStatus:  http.StatusOK,
			expectedLocal:   true,
			expectedHandler: true,
		},
		"wrong URL": {
			URL:             ts.URL + "/petz",
			expectedStatus:  http.StatusNotFound,
			expectedLocal:   false,
			expectedHandler: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			calledGlobal = false
			calledLocal = false
			calledHandler = false
			res, err := http.Get(tc.URL)
			if err != nil {
				t.Fatal(err)
				return
			}
			actualStatus := res.StatusCode
			if tc.expectedStatus != actualStatus {
				t.Fatalf("Got status %d, want %d", actualStatus, tc.expectedStatus)
			}

			if !calledGlobal {
				t.Fatal("Global Middleware was not called")
			}
			if tc.expectedLocal != calledLocal {
				t.Fatalf("Expected calledLocal=%v, got %v", tc.expectedLocal, calledLocal)
			}
			if tc.expectedHandler != calledHandler {
				t.Fatalf("Expected calledHandler=%v, got %v", tc.expectedHandler, calledHandler)
			}
		})
	}
}
