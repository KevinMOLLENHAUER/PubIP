package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Struct representing our tests cases
var ipTests = []struct {
	addr     string // input
	expected int    // expected error
}{
	{"127.0.0.1:1234", 200},
	{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", 400},
	{"[3fff::]:1234", 200},
	{"", 400},
	{":123", 400},
	{"string", 400},
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIPHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/ip", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPubIPHandler)
	// Loop over test cases
	for _, ipt := range ipTests {
		// Overwrite remote addr with test cases addr
		req.RemoteAddr = ipt.addr
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != ipt.expected {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, ipt.expected)
		}
	}
}
