package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// SimpleClient is a basic client for the x402 service
type SimpleClient struct {
	BaseURL string
}

// NewClient creates a new x402 client
func NewClient(baseURL string) *SimpleClient {
	return &SimpleClient{BaseURL: baseURL}
}

// GetHealth checks service health
func (c *SimpleClient) GetHealth() (map[string]interface{}, error) {
	resp, err := http.Get(c.BaseURL + "/health")
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

// GetConfig gets the x402 payment configuration
func (c *SimpleClient) GetConfig() (map[string]interface{}, error) {
	resp, err := http.Get(c.BaseURL + "/.well-known/x402")
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

// GetGasPrice attempts to get gas price (will fail with 402 without payment)
func (c *SimpleClient) GetGasPrice() (*GasData, error) {
	resp, err := http.Get(c.BaseURL + "/api/gas")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusPaymentRequired {
		var paymentReq map[string]interface{}
		if err := json.Unmarshal(body, &paymentReq); err != nil {
			return nil, fmt.Errorf("payment required but couldn't parse response: %v", err)
		}
		return nil, fmt.Errorf("payment required: %v", paymentReq["payment"])
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data           *GasData `json:"data"`
		PaymentVerified bool     `json:"payment_verified"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

func main() {
	client := NewClient(getEnv("X402_API_URL", "http://localhost:8080"))

	// Check health
	fmt.Println("🔍 Checking service health...")
	health, err := client.GetHealth()
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Service healthy: %+v\n\n", health)

	// Get payment config
	fmt.Println("💰 Getting payment configuration...")
	config, err := client.GetConfig()
	if err != nil {
		fmt.Printf("❌ Failed to get config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Payment config: %+v\n\n", config)

	// Try to get gas price (will show payment requirement)
	fmt.Println("⛽ Attempting to get gas price...")
	gas, err := client.GetGasPrice()
	if err != nil {
		fmt.Printf("⚠️  %v\n", err)
		fmt.Println("\nTo access paid endpoints, you need to:")
		fmt.Println("1. Create an x402 payment token")
		fmt.Println("2. Include it in the X-Payment-Response header")
		fmt.Println("\nExample:")
		fmt.Println(`curl -H "X-Payment-Response: <token>" ` + client.BaseURL + `/api/gas`)
	} else {
		fmt.Printf("✅ Gas data: %+v\n", gas)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
