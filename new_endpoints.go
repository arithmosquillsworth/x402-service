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
)

// ==================== NEW ENDPOINT TYPES ====================

// TokenScanRequest represents the input for token risk scanning
type TokenScanRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain"` // "base" or "ethereum"
}

// TokenScanResult represents the output of token scanning
type TokenScanResult struct {
	Address          string   `json:"address"`
	Chain            string   `json:"chain"`
	RiskScore        int      `json:"risk_score"` // 0-100
	Name             string   `json:"name,omitempty"`
	Symbol           string   `json:"symbol,omitempty"`
	IsHoneypot       bool     `json:"is_honeypot"`
	HasMintFunction  bool     `json:"has_mint_function"`
	HasBlacklist     bool     `json:"has_blacklist"`
	IsProxy          bool     `json:"is_proxy"`
	IsVerified       bool     `json:"is_verified"`
	TotalSupply      string   `json:"total_supply,omitempty"`
	HolderCount      int      `json:"holder_count,omitempty"`
	Flags            []string `json:"flags"`
	Warnings         []string `json:"warnings"`
	ScannedAt        int64    `json:"scanned_at"`
}

// WalletScanRequest represents the input for wallet portfolio scanning
type WalletScanRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain"` // "base" or "ethereum"
}

// TokenHolding represents a single token in a wallet
type TokenHolding struct {
	Address     string  `json:"address"`
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	Balance     string  `json:"balance"`
	USDValue    float64 `json:"usd_value,omitempty"`
	RiskScore   int     `json:"risk_score"`
	IsSuspicious bool   `json:"is_suspicious"`
}

// WalletScanResult represents the output of wallet scanning
type WalletScanResult struct {
	Address         string         `json:"address"`
	Chain           string         `json:"chain"`
	ETHBalance      string         `json:"eth_balance"`
	TotalUSDValue   float64        `json:"total_usd_value"`
	TokenCount      int            `json:"token_count"`
	Holdings        []TokenHolding `json:"holdings"`
	SuspiciousTokens int           `json:"suspicious_tokens"`
	RiskScore       int            `json:"risk_score"` // Aggregate risk
	ScannedAt       int64          `json:"scanned_at"`
}

// AddressLabelRequest represents the input for address label lookup
type AddressLabelRequest struct {
	Address string `json:"address"`
}

// AddressLabelResult represents the output of label lookup
type AddressLabelResult struct {
	Address     string   `json:"address"`
	Labels      []string `json:"labels"`
	Entity      string   `json:"entity,omitempty"`
	Category    string   `json:"category,omitempty"` // "exchange", "contract", "wallet", "agent"
	RiskLevel   string   `json:"risk_level,omitempty"` // "low", "medium", "high"
	Confidence  float64  `json:"confidence"`
	Sources     []string `json:"sources"`
	CheckedAt   int64    `json:"checked_at"`
}

// MEVCheckRequest represents the input for MEV risk checking
type MEVCheckRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

// MEVCheckResult represents the output of MEV checking
type MEVCheckResult struct {
	Safe              bool     `json:"safe"`
	MEVRiskScore      int      `json:"mev_risk_score"` // 0-100
	RiskFactors       []string `json:"risk_factors"`
	SandwichRisk      string   `json:"sandwich_risk"`      // "low", "medium", "high"
	FrontrunRisk      string   `json:"frontrun_risk"`      // "low", "medium", "high"
	GasPriceRisk      string   `json:"gas_price_risk"`     // "low", "medium", "high"
	RecommendedSlippage string `json:"recommended_slippage"`
	ProtectedRPCs     []string `json:"protected_rpcs,omitempty"`
	CheckedAt         int64    `json:"checked_at"`
}

// ==================== TOKEN SCANNER ====================

// handleTokenScan scans a token contract for risks
func handleTokenScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TokenScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Token scan decode error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate address
	if !isValidAddress(req.Address) {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	// Default to base if not specified
	if req.Chain == "" {
		req.Chain = "base"
	}
	if req.Chain != "base" && req.Chain != "ethereum" {
		http.Error(w, "Chain must be 'base' or 'ethereum'", http.StatusBadRequest)
		return
	}

	// Perform scan (mock for now, would integrate with API)
	result := scanToken(req.Address, req.Chain)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":            result,
		"payment_verified": true,
	})
}

// scanToken performs token risk analysis
func scanToken(address, chain string) TokenScanResult {
	result := TokenScanResult{
		Address:   address,
		Chain:     chain,
		ScannedAt: time.Now().Unix(),
		Flags:     []string{},
		Warnings:  []string{},
	}

	// Try to fetch contract info from explorer
	apiKey := getAPIKeyForChain(chain)
	if apiKey != "" {
		// Fetch contract ABI to check for risky functions
		if abi, err := fetchContractABI(address, chain, apiKey); err == nil {
			// Check for mint function
			if strings.Contains(abi, "mint") || strings.Contains(abi, "_mint") {
				result.HasMintFunction = true
				result.Warnings = append(result.Warnings, "Contract has mint function - supply can be inflated")
				result.RiskScore += 20
			}

			// Check for blacklist
			if strings.Contains(abi, "blacklist") || strings.Contains(abi, "blocked") {
				result.HasBlacklist = true
				result.Warnings = append(result.Warnings, "Contract can blacklist addresses")
				result.RiskScore += 15
			}

			// Check for proxy
			if strings.Contains(abi, "delegatecall") || strings.Contains(abi, "implementation") {
				result.IsProxy = true
				result.Flags = append(result.Flags, "proxy_contract")
			}

			result.IsVerified = true
		} else {
			result.IsVerified = false
			result.Flags = append(result.Flags, "unverified_contract")
			result.Warnings = append(result.Warnings, "Contract source code is not verified")
			result.RiskScore += 30
		}
	}

	// Additional heuristics would go here:
	// - Check honeypot.is or similar service
	// - Analyze holder distribution
	// - Check liquidity locked
	// - Verify ownership renounced

	// Cap risk score
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}

	return result
}

// ==================== WALLET SCANNER ====================

// handleWalletScan scans a wallet for portfolio risks
func handleWalletScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WalletScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Wallet scan decode error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isValidAddress(req.Address) {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	if req.Chain == "" {
		req.Chain = "base"
	}

	result := scanWallet(req.Address, req.Chain)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":            result,
		"payment_verified": true,
	})
}

// scanWallet performs wallet portfolio analysis
func scanWallet(address, chain string) WalletScanResult {
	result := WalletScanResult{
		Address:   address,
		Chain:     chain,
		ScannedAt: time.Now().Unix(),
		Holdings:  []TokenHolding{},
	}

	// This would integrate with:
	// - Covalent/Alchemy to get token balances
	// - DexScreener for pricing
	// - Internal risk database for token scoring

	// Mock data structure for now
	result.ETHBalance = "0.5"
	result.TotalUSDValue = 1250.50
	result.TokenCount = 5

	// Example holding
	result.Holdings = append(result.Holdings, TokenHolding{
		Address:      "0x...",
		Symbol:       "EXAMPLE",
		Name:         "Example Token",
		Balance:      "1000",
		USDValue:     100.0,
		RiskScore:    25,
		IsSuspicious: false,
	})

	// Calculate aggregate risk
	for _, h := range result.Holdings {
		if h.IsSuspicious {
			result.SuspiciousTokens++
		}
		result.RiskScore += h.RiskScore
	}

	if result.TokenCount > 0 {
		result.RiskScore = result.RiskScore / result.TokenCount
	}

	if result.SuspiciousTokens > 0 {
		result.RiskScore += result.SuspiciousTokens * 10
	}

	if result.RiskScore > 100 {
		result.RiskScore = 100
	}

	return result
}

// ==================== ADDRESS LABEL LOOKUP ====================

// handleAddressLabel looks up labels for an address
func handleAddressLabel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddressLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Address label decode error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if !isValidAddress(req.Address) {
		http.Error(w, "Invalid address format", http.StatusBadRequest)
		return
	}

	result := lookupAddressLabel(req.Address)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":            result,
		"payment_verified": true,
	})
}

// Known address labels database (expandable)
var knownLabels = map[string]AddressLabelResult{
	"0x4200000000000000000000000000000000000006": { // Base WETH
		Labels:     []string{"weth", "wrapped-ether"},
		Entity:     "Wrapped Ether",
		Category:   "contract",
		RiskLevel:  "low",
		Confidence: 1.0,
		Sources:    []string{"official"},
	},
	"0x833589fcd6edb6e08f4c7c32d4f71b54bda02913": { // Base USDC
		Labels:     []string{"usdc", "stablecoin"},
		Entity:     "USD Coin",
		Category:   "contract",
		RiskLevel:  "low",
		Confidence: 1.0,
		Sources:    []string{"official"},
	},
}

// lookupAddressLabel looks up labels for an address
func lookupAddressLabel(address string) AddressLabelResult {
	address = strings.ToLower(address)

	// Check known labels
	if label, ok := knownLabels[address]; ok {
		label.Address = address
		label.CheckedAt = time.Now().Unix()
		return label
	}

	// Default: unknown address
	return AddressLabelResult{
		Address:    address,
		Labels:     []string{},
		Category:   "unknown",
		RiskLevel:  "medium",
		Confidence: 0.0,
		Sources:    []string{},
		CheckedAt:  time.Now().Unix(),
	}
}

// ==================== MEV PROTECTION CHECK ====================

// handleMEVCheck checks transaction for MEV risks
func handleMEVCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MEVCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("MEV check decode error: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result := checkMEVRisk(req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":            result,
		"payment_verified": true,
	})
}

// checkMEVRisk analyzes transaction for MEV exposure
func checkMEVRisk(req MEVCheckRequest) MEVCheckResult {
	result := MEVCheckResult{
		Safe:              true,
		MEVRiskScore:      0,
		RiskFactors:       []string{},
		SandwichRisk:      "low",
		FrontrunRisk:      "low",
		GasPriceRisk:      "low",
		RecommendedSlippage: "0.5%",
		ProtectedRPCs:     []string{
			"https://rpc.flashbots.net",
			"https://mevblocker.io",
		},
		CheckedAt:         time.Now().Unix(),
	}

	// Parse value
	valueWei := req.Value
	if valueWei == "" {
		valueWei = "0x0"
	}

	// Check for high value transfers
	if valueWei != "0x0" && valueWei != "0" {
		value, _ := strconv.ParseInt(strings.TrimPrefix(valueWei, "0x"), 16, 64)
		if value > 1000000000000000000 { // > 1 ETH
			result.RiskFactors = append(result.RiskFactors, "high_value_transfer")
			result.MEVRiskScore += 20
		}
	}

	// Check transaction data
	if len(req.Data) > 10 {
		// DEX swap signatures
		swapSignatures := []string{
			"0x38ed1739", // swapExactTokensForTokens
			"0x8803dbee", // swapTokensForExactTokens
			"0x7ff36ab5", // swapExactETHForTokens
			"0x18cbafe5", // swapExactTokensForETH
		}

		for _, sig := range swapSignatures {
			if strings.HasPrefix(req.Data, sig) {
				result.RiskFactors = append(result.RiskFactors, "dex_swap_detected")
				result.SandwichRisk = "medium"
				result.MEVRiskScore += 30
				result.RecommendedSlippage = "0.1%"
				break
			}
		}
	}

	// Check gas price risk
	gasPrice, _ := getCurrentGasPrice()
	if gasPrice > 50 {
		result.GasPriceRisk = "high"
		result.MEVRiskScore += 10
	} else if gasPrice > 20 {
		result.GasPriceRisk = "medium"
		result.MEVRiskScore += 5
	}

	// Determine overall safety
	if result.MEVRiskScore >= 50 {
		result.Safe = false
		result.FrontrunRisk = "high"
	} else if result.MEVRiskScore >= 30 {
		result.Safe = true
		result.FrontrunRisk = "medium"
	}

	if result.MEVRiskScore > 100 {
		result.MEVRiskScore = 100
	}

	return result
}

// ==================== HELPERS ====================

// getAPIKeyForChain returns the appropriate API key
func getAPIKeyForChain(chain string) string {
	if chain == "base" {
		return os.Getenv("BASESCAN_API_KEY")
	}
	return os.Getenv("ETHERSCAN_API_KEY")
}

// fetchContractABI fetches contract ABI from explorer
func fetchContractABI(address, chain, apiKey string) (string, error) {
	baseURL := "https://api.basescan.org/api"
	if chain == "ethereum" {
		baseURL = "https://api.etherscan.io/api"
	}

	url := fmt.Sprintf("%s?module=contract&action=getabi&address=%s&apikey=%s",
		baseURL, address, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "1" {
		return "", fmt.Errorf("API error: %s", result.Message)
	}

	return result.Result, nil
}

// getCurrentGasPrice gets current gas price (simplified)
func getCurrentGasPrice() (float64, error) {
	// Would integrate with actual RPC
	return 20.0, nil // Placeholder
}

// isValidAddress validates Ethereum address format
func isValidAddress(addr string) bool {
	return len(addr) == 42 && strings.HasPrefix(addr, "0x")
}
