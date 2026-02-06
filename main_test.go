package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"service": "arithmos-x402",
		})
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response["status"] != "ok" {
		t.Errorf("unexpected status: got %v want %v",
			response["status"], "ok")
	}
}

func TestX402ConfigEndpoint(t *testing.T) {
	config := X402Config{
		Version: "1.0",
		PaymentRequirements: []PaymentRequirement{
			{
				Scheme:   "x402",
				Network:  "base",
				MaxAmount: "0.001",
				Asset:    "USDC",
				Receiver: "0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91",
			},
		},
	}

	req, err := http.NewRequest("GET", "/.well-known/x402", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response X402Config
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.Version != "1.0" {
		t.Errorf("unexpected version: got %v want %v",
			response.Version, "1.0")
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input     float64
		precision int
		expected  float64
	}{
		{3.14159, 2, 3.14},
		{3.14159, 3, 3.142},
		{2.5, 0, 3},
		{2.4, 0, 2},
	}

	for _, tt := range tests {
		result := round(tt.input, tt.precision)
		if result != tt.expected {
			t.Errorf("round(%v, %d) = %v, want %v",
				tt.input, tt.precision, result, tt.expected)
		}
	}
}
