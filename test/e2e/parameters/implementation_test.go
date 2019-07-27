package parameters

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestPathParameters(t *testing.T) {
	s := NewOpenAPIServer(&server{})
	s.BootPetShow()
	ts := httptest.NewServer(s)
	tests := map[string]struct {
		URL                string
		expectedHTTPStatus int
		expectedID         int
	}{
		"with numeric ID": {
			URL:                ts.URL + "/pets/1",
			expectedHTTPStatus: http.StatusOK,
			expectedID:         1,
		},
		"with string ID": {
			URL:                ts.URL + "/pets/abc",
			expectedHTTPStatus: http.StatusBadRequest,
			expectedID:         999,
		},
		"with incorrect URL": {
			URL:                ts.URL + "/pets/1/xyz",
			expectedHTTPStatus: http.StatusNotFound,
			expectedID:         999,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			receivedID = 999
			res, err := http.Get(tc.URL)
			if err != nil {
				t.Fatal(err)
				return
			}
			actualHTTPStatus := res.StatusCode
			if tc.expectedHTTPStatus != actualHTTPStatus {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatalf("Got status %d, want %d. Body was: %s", actualHTTPStatus, tc.expectedHTTPStatus, body)
			}

			if tc.expectedID != receivedID {
				t.Fatalf("Expected ID param value: %d, got: %d", tc.expectedID, receivedID)
			}
		})
	}
}

func TestQueryParameters(t *testing.T) {
	s := NewOpenAPIServer(&server{})
	s.BootPetIndex()
	ts := httptest.NewServer(s)
	tests := map[string]struct {
		URL                 string
		expectedHTTPStatus  int
		expectedParamFoo    string
		expectedParamStatus int
		expectedParamBar    int
	}{
		"foo not given": {
			URL:                 ts.URL + "/pets",
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedParamFoo:    "---",
			expectedParamStatus: 999,
			expectedParamBar:    999,
		},
		"empty foo given": {
			URL:                 ts.URL + "/pets?foo=",
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedParamFoo:    "---",
			expectedParamStatus: 999,
			expectedParamBar:    999,
		},
		"foo given, status given": {
			URL:                 ts.URL + "/pets?foo=abc&status=1",
			expectedHTTPStatus:  http.StatusOK,
			expectedParamFoo:    "abc",
			expectedParamStatus: 1,
			expectedParamBar:    10,
		},
		"foo given, status not given": {
			URL:                 ts.URL + "/pets?foo=abc",
			expectedHTTPStatus:  http.StatusOK,
			expectedParamFoo:    "abc",
			expectedParamStatus: 0,
			expectedParamBar:    10,
		},
		"foo given, status given, bar given": {
			URL:                 ts.URL + "/pets?foo=abc&bar=5",
			expectedHTTPStatus:  http.StatusOK,
			expectedParamFoo:    "abc",
			expectedParamStatus: 0,
			expectedParamBar:    5,
		},
		"with incorrect URL": {
			URL:                 ts.URL + "/petz?foo=abc",
			expectedHTTPStatus:  http.StatusNotFound,
			expectedParamFoo:    "---",
			expectedParamStatus: 999,
			expectedParamBar:    999,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			receivedFoo = "---"
			receivedStatus = 999
			receivedBar = 999
			res, err := http.Get(tc.URL)
			if err != nil {
				t.Fatal(err)
				return
			}
			actualHTTPStatus := res.StatusCode
			if tc.expectedHTTPStatus != actualHTTPStatus {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatalf("Got status %d, want %d. Body was: %s", actualHTTPStatus, tc.expectedHTTPStatus, body)
			}

			if tc.expectedParamStatus != receivedStatus {
				t.Fatalf("Expected status param value: %d, got: %d", tc.expectedParamStatus, receivedStatus)
			}

			if tc.expectedParamBar != receivedBar {
				t.Fatalf("Expected bar param value: %d, got: %d", tc.expectedParamBar, receivedBar)
			}
		})
	}
}

func TestHeaderParameters(t *testing.T) {
	s := NewOpenAPIServer(&server{})
	s.BootCatIndex()
	ts := httptest.NewServer(s)
	fixedClientID := "668c31fb-0cf5-4a72-b8bd-7866bcb7151e"
	emptyClientID := ""
	goodClientID := "21198093-191f-4c27-9874-cd899e94f63c"
	badClientID := "abc"
	fixedClientTime := "0001-01-01T00:00:00+00:00"
	someClientTime := "2012-11-01T22:08:41+00:00"
	tests := map[string]struct {
		URL                     string
		clientID                *string
		clientTime              *string
		expectedHTTPStatus      int
		expectedParamClientID   string
		expectedParamClientTime string
	}{
		"client ID not given": {
			URL:                     ts.URL + "/cats",
			clientID:                nil,
			clientTime:              nil,
			expectedHTTPStatus:      http.StatusBadRequest,
			expectedParamClientID:   fixedClientID,
			expectedParamClientTime: fixedClientTime,
		},
		"empty client ID given": {
			URL:                     ts.URL + "/cats",
			clientID:                &emptyClientID,
			clientTime:              nil,
			expectedHTTPStatus:      http.StatusBadRequest,
			expectedParamClientID:   fixedClientID,
			expectedParamClientTime: fixedClientTime,
		},
		"client ID given": {
			URL:                     ts.URL + "/cats",
			clientID:                &goodClientID,
			clientTime:              nil,
			expectedHTTPStatus:      http.StatusOK,
			expectedParamClientID:   goodClientID,
			expectedParamClientTime: fixedClientTime,
		},
		"malformed client ID given": {
			URL:                     ts.URL + "/cats",
			clientID:                &badClientID,
			clientTime:              nil,
			expectedHTTPStatus:      http.StatusBadRequest,
			expectedParamClientID:   fixedClientID,
			expectedParamClientTime: fixedClientTime,
		},
		"client ID and time given": {
			URL:                     ts.URL + "/cats",
			clientID:                &goodClientID,
			clientTime:              &someClientTime,
			expectedHTTPStatus:      http.StatusOK,
			expectedParamClientID:   goodClientID,
			expectedParamClientTime: someClientTime,
		},
	}

	client := &http.Client{}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			receivedClientID, _ = uuid.Parse(fixedClientID)
			receivedClientTime, _ = time.Parse(time.RFC3339, fixedClientTime)
			req, _ := http.NewRequest("GET", tc.URL, nil)
			if tc.clientID != nil {
				req.Header.Set("X-Client-ID", *tc.clientID)
			}
			if tc.clientTime != nil {
				req.Header.Set("X-Client-Time", *tc.clientTime)
			}
			res, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
				return
			}
			actualHTTPStatus := res.StatusCode
			if tc.expectedHTTPStatus != actualHTTPStatus {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Fatalf("Expected HTTP status %d, got %d. Body was: %s", tc.expectedHTTPStatus, actualHTTPStatus, body)
			}

			if tc.expectedParamClientID != receivedClientID.String() {
				t.Fatalf("Expected X-Client-ID param value: %s, got: %s", tc.expectedParamClientID, receivedClientID)
			}

			if tc.expectedParamClientTime != receivedClientTime.Format("2006-01-02T15:04:05+00:00") {
				t.Fatalf("Expected X-Client-Time param value: %s, got: %s", tc.expectedParamClientTime, receivedClientTime.Format("2006-01-02T15:04:05+00:00"))
			}
		})
	}
}
