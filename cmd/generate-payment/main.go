package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PaymentToken represents the x402 payment token
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
	if len(os.Args) < 4 {
		fmt.Println("Usage: generate-payment <receiver-address> <amount> <network>")
		fmt.Println("\nExample:")
		fmt.Println("  generate-payment 0x123... 0.001 base")
		fmt.Println("\nEnvironment variables:")
		fmt.Println("  X402_SIGNING_KEY - Secret key for signing (required)")
		fmt.Println("  X402_ASSET       - Asset to use (default: USDC)")
		fmt.Println("  X402_EXPIRY_MIN  - Expiry in minutes (default: 5)")
		os.Exit(1)
	}

	receiver := os.Args[1]
	amount := os.Args[2]
	network := os.Args[3]

	// Get config from env
	asset := getEnv("X402_ASSET", "USDC")
	expiryMin := 5
	if e := os.Getenv("X402_EXPIRY_MIN"); e != "" {
		if m, err := time.ParseDuration(e + "m"); err == nil {
			expiryMin = int(m.Minutes())
		}
	}

	// Get signing key
	signingKey := os.Getenv("X402_SIGNING_KEY")
	if signingKey == "" {
		fmt.Println("❌ X402_SIGNING_KEY environment variable required")
		os.Exit(1)
	}

	// Create claims
	claims := PaymentToken{
		Payment: struct {
			Amount   string `json:"amount"`
			Asset    string `json:"asset"`
			Receiver string `json:"receiver"`
			Network  string `json:"network"`
		}{
			Amount:   amount,
			Asset:    asset,
			Receiver: receiver,
			Network:  network,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   receiver,
			ID:        fmt.Sprintf("%d", time.Now().Unix()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMin) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		fmt.Printf("❌ Failed to sign token: %v\n", err)
		os.Exit(1)
	}

	// Display
	fmt.Println("✅ x402 Payment Token Generated")
	fmt.Println("================================")
	fmt.Printf("Amount:   %s %s\n", amount, asset)
	fmt.Printf("Receiver: %s\n", receiver)
	fmt.Printf("Network:  %s\n", network)
	fmt.Printf("Expires:  %d minutes\n", expiryMin)
	fmt.Println("\nToken:")
	fmt.Println(tokenString)
	fmt.Println("\nUsage:")
	fmt.Printf("  curl -H \"X-Payment-Response: %s\" http://localhost:8080/api/gas\n", tokenString[:50]+"...")
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
