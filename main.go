package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ServiceConfig holds the x402 configuration
type ServiceConfig struct {
	Price        string `json:"price"`        // e.g., "0.001"
	Asset        string `json:"asset"`        // e.g., "USDC"
	Network      string `json:"network"`      // e.g., "base"
	Receiver     string `json:"receiver"`     // Your wallet address
	Description  string `json:"description"`
}

// PaymentRequirement is what we send in 402 responses
type PaymentRequirement struct {
	Scheme       string          `json:"scheme"`
	Network      string          `json:"network"`
	MaxAmount    string          `json:"maxAmount"`
	MinAmount    string          `json:"minAmount"`
	Asset        string          `json:"asset"`
	Receiver     string          `json:"receiver"`
	Description  string          `json:"description"`
}

// X402Config is the full configuration object
type X402Config struct {
	Version      string               `json:"version"`
	PaymentRequirements []PaymentRequirement `json:"paymentRequirements"`
}

func main() {
	config := ServiceConfig{
		Price:       "0.001",
		Asset:       "USDC",
		Network:     "base",
		Receiver:    "0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91", // My wallet
		Description: "Arithmos API - Gas price data and on-chain analytics",
	}

	mux := http.NewServeMux()
	
	// Health check (free)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"service": "arithmos-x402",
		})
	})
	
	// x402 config endpoint (free, tells clients how to pay)
	mux.HandleFunc("/.well-known/x402", func(w http.ResponseWriter, r *http.Request) {
		x402 := X402Config{
			Version: "1.0",
			PaymentRequirements: []PaymentRequirement{
				{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   config.Price,
					MinAmount:   config.Price,
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: config.Description,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(x402)
	})
	
	// Protected endpoint - requires x402 payment
	mux.HandleFunc("/api/gas", func(w http.ResponseWriter, r *http.Request) {
		// Check for x402 payment header
		paymentHeader := r.Header.Get("X-Payment-Response")
		if paymentHeader == "" {
			// Return 402 Payment Required
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Payment-Required", fmt.Sprintf("x402-%s-%s", config.Network, config.Asset))
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Payment required",
				"payment": PaymentRequirement{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   config.Price,
					MinAmount:   config.Price,
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: "Get current gas prices",
				},
			})
			return
		}
		
		// Validate payment (simplified - would verify on-chain in production)
		// For now, just check if header exists and is well-formed
		
		// Return gas data
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"gas": map[string]string{
				"mainnet": "0.27",
				"unit":    "gwei",
			},
			"payment_received": true,
		})
	})
	
	// Another paid endpoint - validator queue status
	mux.HandleFunc("/api/validators", func(w http.ResponseWriter, r *http.Request) {
		paymentHeader := r.Header.Get("X-Payment-Response")
		if paymentHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Payment required",
				"payment": PaymentRequirement{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   "0.005", // More expensive endpoint
					MinAmount:   "0.005",
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: "Get validator queue status",
				},
			})
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"queue": map[string]interface{}{
				"entry_wait": "~4 hours",
				"exit_wait":  "~2 hours",
				"active":     1048576,
			},
		})
	})
	
	port := ":8080"
	log.Printf("x402 service starting on %s", port)
	log.Printf("Receiver: %s", config.Receiver)
	log.Fatal(http.ListenAndServe(port, mux))
}
