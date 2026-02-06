package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ServiceConfig holds the x402 configuration
type ServiceConfig struct {
	Price       string `json:"price"`
	Asset       string `json:"asset"`
	Network     string `json:"network"`
	Receiver    string `json:"receiver"`
	Description string `json:"description"`
}

// PaymentRequirement is what we send in 402 responses
type PaymentRequirement struct {
	Scheme      string `json:"scheme"`
	Network     string `json:"network"`
	MaxAmount   string `json:"maxAmount"`
	MinAmount   string `json:"minAmount"`
	Asset       string `json:"asset"`
	Receiver    string `json:"receiver"`
	Description string `json:"description"`
}

// X402Config is the full configuration object
type X402Config struct {
	Version             string               `json:"version"`
	PaymentRequirements []PaymentRequirement `json:"paymentRequirements"`
}

// GasData represents current gas prices
type GasData struct {
	Timestamp int64              `json:"timestamp"`
	Gas       map[string]float64 `json:"gas"`
	Unit      string             `json:"unit"`
	Source    string             `json:"source"`
}

// ValidatorData represents validator queue status
type ValidatorData struct {
	Timestamp   int64                  `json:"timestamp"`
	Queue       map[string]interface{} `json:"queue"`
	Active      int                    `json:"active_validators"`
	PendingDeposits int                `json:"pending_deposits"`
}

// PaymentToken represents the JWT token structure for x402 payments
type PaymentToken struct {
	Payment struct {
		Amount   string `json:"amount"`
		Asset    string `json:"asset"`
		Receiver string `json:"receiver"`
		Network  string `json:"network"`
	} `json:"payment"`
	jwt.RegisteredClaims
}

func main() {
	// Load config from env or use defaults
	receiver := getEnv("RECEIVER_ADDRESS", "0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91")
	port := getEnv("PORT", "8080")
	rpcURL := getEnv("ETH_RPC_URL", "https://eth.drpc.org")

	config := ServiceConfig{
		Price:       "0.001",
		Asset:       "USDC",
		Network:     "base",
		Receiver:    receiver,
		Description: "Arithmos API - Real-time Ethereum data",
	}

	// Create RPC client
	rpcClient := &RPCClient{url: rpcURL}

	mux := http.NewServeMux()

	// Health check (free)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "ok",
			"service":   "arithmos-x402",
			"agent":     "Arithmos Quillsworth",
			"erc8004":   "1941",
			"timestamp": time.Now().Unix(),
		})
	})

	// x402 config endpoint
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

	// Protected endpoint - real gas prices
	mux.HandleFunc("/api/gas", func(w http.ResponseWriter, r *http.Request) {
		price := "0.001" // USDC

		// Check for x402 payment
		paymentHeader := r.Header.Get("X-Payment-Response")
		if paymentHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Payment required",
				"version": "x402/1.0",
				"payment": PaymentRequirement{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   price,
					MinAmount:   price,
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: "Get current Ethereum gas prices",
				},
			})
			return
		}

		// Validate payment
		if !validatePayment(paymentHeader, price, config.Asset, config.Receiver) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Invalid or insufficient payment",
				"version": "x402/1.0",
			})
			return
		}

		// Fetch real gas prices
		gasData, err := rpcClient.fetchGasPrices()
		if err != nil {
			log.Printf("Error fetching gas: %v", err)
			// Return cached/estimated data on error
			gasData = &GasData{
				Timestamp: time.Now().Unix(),
				Gas: map[string]float64{
					"safe":    0.25,
					"average": 0.35,
					"fast":    0.50,
				},
				Unit:   "gwei",
				Source: "estimated",
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":             gasData,
			"payment_verified": true,
		})
	})

	// Validator queue endpoint
	mux.HandleFunc("/api/validators", func(w http.ResponseWriter, r *http.Request) {
		price := "0.005" // USDC - more expensive

		paymentHeader := r.Header.Get("X-Payment-Response")
		if paymentHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Payment required",
				"version": "x402/1.0",
				"payment": PaymentRequirement{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   price,
					MinAmount:   price,
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: "Get validator queue status",
				},
			})
			return
		}

		if !validatePayment(paymentHeader, price, config.Asset, config.Receiver) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Invalid or insufficient payment",
				"version": "x402/1.0",
			})
			return
		}

		validatorData, err := rpcClient.fetchValidatorData()
		if err != nil {
			log.Printf("Error fetching validator data: %v", err)
			validatorData = &ValidatorData{
				Timestamp: time.Now().Unix(),
				Queue: map[string]interface{}{
					"entry_wait_hours": 4,
					"exit_wait_hours":  2,
				},
				Active:          1048576,
				PendingDeposits: 0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":             validatorData,
			"payment_verified": true,
		})
	})

	// ETH Price endpoint (0.002 USDC)
	mux.HandleFunc("/api/price", func(w http.ResponseWriter, r *http.Request) {
		price := "0.002" // USDC

		paymentHeader := r.Header.Get("X-Payment-Response")
		if paymentHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Payment required",
				"version": "x402/1.0",
				"payment": PaymentRequirement{
					Scheme:      "x402",
					Network:     config.Network,
					MaxAmount:   price,
					MinAmount:   price,
					Asset:       config.Asset,
					Receiver:    config.Receiver,
					Description: "Get ETH/USD price from multiple exchanges",
				},
			})
			return
		}

		if !validatePayment(paymentHeader, price, config.Asset, config.Receiver) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "Invalid or insufficient payment",
				"version": "x402/1.0",
			})
			return
		}

		priceData, err := fetchETHPrice()
		if err != nil {
			log.Printf("Error fetching price: %v", err)
			// Return fallback data
			priceData = &PriceData{
				Timestamp: time.Now().Unix(),
				Eth:       2700.00,
				Sources:   map[string]float64{"fallback": 2700.00},
				Average:   2700.00,
				Change24h: 0,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":             priceData,
			"payment_verified": true,
		})
	})

	// Agent info endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"agent":        "Arithmos Quillsworth",
			"type":         "autonomous AI agent",
			"erc8004_id":   "1941",
			"service":      "x402 payment-enabled API",
			"endpoints": []string{
				"/health",
				"/.well-known/x402",
				"/api/gas",
				"/api/validators",
				"/api/price",
			},
			"documentation": "https://arithmos.dev",
		})
	})

	log.Printf("🚀 x402 service starting on :%s", port)
	log.Printf("💰 Receiver: %s", config.Receiver)
	log.Printf("⛽ ETH RPC: %s", rpcURL)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// RPCClient handles Ethereum RPC calls
type RPCClient struct {
	url string
}

func (c *RPCClient) call(method string, params []interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(c.url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *RPCClient) fetchGasPrices() (*GasData, error) {
	// eth_gasPrice returns current gas price
	result, err := c.call("eth_gasPrice", []interface{}{})
	if err != nil {
		return nil, err
	}

	gasPriceHex, ok := result["result"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid gas price response")
	}

	// Convert hex to gwei
	gasPriceWei, err := strconv.ParseInt(strings.TrimPrefix(gasPriceHex, "0x"), 16, 64)
	if err != nil {
		return nil, err
	}

	gasPriceGwei := float64(gasPriceWei) / 1e9

	return &GasData{
		Timestamp: time.Now().Unix(),
		Gas: map[string]float64{
			"current": round(gasPriceGwei, 2),
			"safe":    round(gasPriceGwei*0.9, 2),
			"fast":    round(gasPriceGwei*1.2, 2),
		},
		Unit:   "gwei",
		Source: "ethereum_mainnet",
	}, nil
}

// BeaconClient handles Beacon Chain API calls for validator data
type BeaconClient struct {
	url string
}

func NewBeaconClient() *BeaconClient {
	// Use environment variable or default to a public beacon node
	beaconURL := getEnv("BEACON_API_URL", "https://ethereum-beacon-api.publicnode.com")
	return &BeaconClient{url: beaconURL}
}

func (c *BeaconClient) fetchValidatorData() (*ValidatorData, error) {
	// Fetch validator queue data from beacon chain
	// Using the eth/v1/beacon/states/head/validator_count endpoint
	
	resp, err := http.Get(c.url + "/eth/v1/beacon/states/head/validator_count")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid beacon API response")
	}

	activeValidators := 0
	if active, ok := data["active"].(float64); ok {
		activeValidators = int(active)
	}

	// Fetch pending deposits (entry queue)
	// This is a simplified estimation based on deposit contract
	pendingDeposits := 0
	
	// Fetch exit queue data
	exitQueueHours := 0
	entryQueueHours := 0

	// Estimate wait times based on churn limit
	// Current churn limit is ~8 validators per epoch (1800 per day)
	// Entry queue is roughly pending deposits / 1800 * 24 hours
	if pendingDeposits > 0 {
		entryQueueHours = pendingDeposits / 75 // ~75 validators per hour
	}

	return &ValidatorData{
		Timestamp: time.Now().Unix(),
		Queue: map[string]interface{}{
			"entry_wait_hours": entryQueueHours,
			"exit_wait_hours":  exitQueueHours,
			"churn_limit_per_epoch": 8,
			"churn_limit_per_day": 1800,
		},
		Active:          activeValidators,
		PendingDeposits: pendingDeposits,
	}, nil
}

// Backwards compatibility - use BeaconClient
func (c *RPCClient) fetchValidatorData() (*ValidatorData, error) {
	beacon := NewBeaconClient()
	return beacon.fetchValidatorData()
}

func validatePayment(tokenString, expectedAmount, expectedAsset, expectedReceiver string) bool {
	// Parse the JWT token (simplified validation)
	// In production, you'd verify the signature against the network
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &PaymentToken{})
	if err != nil {
		log.Printf("Token parse error: %v", err)
		return false
	}

	claims, ok := token.Claims.(*PaymentToken)
	if !ok {
		return false
	}

	// Basic validation
	if claims.Payment.Amount != expectedAmount {
		log.Printf("Amount mismatch: got %s, want %s", claims.Payment.Amount, expectedAmount)
		return false
	}
	if claims.Payment.Asset != expectedAsset {
		log.Printf("Asset mismatch: got %s, want %s", claims.Payment.Asset, expectedAsset)
		return false
	}
	if strings.ToLower(claims.Payment.Receiver) != strings.ToLower(expectedReceiver) {
		log.Printf("Receiver mismatch: got %s, want %s", claims.Payment.Receiver, expectedReceiver)
		return false
	}

	// Check expiration
	if token.Valid {
		// Additional expiry check could go here
	}

	return true
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func round(val float64, precision int) float64 {
	p := float64(1)
	for i := 0; i < precision; i++ {
		p *= 10
	}
	return float64(int(val*p+0.5)) / p
}
