package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleHealth(t *testing.T) {
	// seeding random source
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// metrics map
	metrics = make(map[string]*Metric)
	metrics["health"] = &Metric{0, 0.0, 0.0}
	metrics["ready"] = &Metric{0, 0.0, 0.0}
	metrics["payload"] = &Metric{0, 0.0, 0.0}
	metrics["metrics"] = &Metric{0, 0.0, 0.0}
	// Create a new test server with the desired handler
	ts := httptest.NewServer(http.HandlerFunc(handleHealth))
	defer ts.Close()

	// Send a GET request to the test server
	res, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("error sending request: %v", err)
	}
	defer res.Body.Close()

	// Check the status code is what we expect.
	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", res.StatusCode, http.StatusOK)
	}

	// Decode the response body
	var responseBody string
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	// Check the response body is what we expect.
	expected := "Service is healthy"
	if responseBody != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", responseBody, expected)
	}
}
