package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
	if len(os.Args) < 2 {
		fmt.Println("Usage: verify-payment <jwt-token>")
		fmt.Println("\nExample:")
		fmt.Println("  verify-payment eyJhbGciOiJIUzI1NiIs...")
		os.Exit(1)
	}

	tokenString := os.Args[1]

	// Parse without verification (for inspection)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &PaymentToken{})
	if err != nil {
		fmt.Printf("❌ Failed to parse token: %v\n", err)
		os.Exit(1)
	}

	claims, ok := token.Claims.(*PaymentToken)
	if !ok {
		fmt.Println("❌ Invalid token claims")
		os.Exit(1)
	}

	// Display token info
	fmt.Println("🔍 x402 Payment Token Details")
	fmt.Println("==============================")
	fmt.Printf("Amount:   %s %s\n", claims.Payment.Amount, claims.Payment.Asset)
	fmt.Printf("Receiver: %s\n", claims.Payment.Receiver)
	fmt.Printf("Network:  %s\n", claims.Payment.Network)
	
	if claims.Subject != "" {
		fmt.Printf("Subject:  %s\n", claims.Subject)
	}
	if claims.ID != "" {
		fmt.Printf("ID:       %s\n", claims.ID)
	}
	if claims.ExpiresAt != nil {
		fmt.Printf("Expires:  %s\n", claims.ExpiresAt.Time)
	}
	if claims.IssuedAt != nil {
		fmt.Printf("Issued:   %s\n", claims.IssuedAt.Time)
	}

	// Show raw header
	parts := strings.Split(tokenString, ".")
	if len(parts) == 3 {
		fmt.Println("\n📋 Token Structure:")
		fmt.Printf("  Header:    %s...\n", parts[0][:20])
		fmt.Printf("  Payload:   %s...\n", parts[1][:20])
		fmt.Printf("  Signature: %s...\n", parts[2][:20])
	}

	fmt.Println("\n⚠️  Note: This tool only inspects the token.")
	fmt.Println("   To fully verify, the signature must be checked against the network.")
}
