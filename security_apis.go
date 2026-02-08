package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ==================== TYPES ====================

// ContractScanRequest represents the input for contract scanning
type ContractScanRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain"` // "base" or "ethereum"
}

// ContractScanResult represents the output of contract scanning
type ContractScanResult struct {
	Address     string   `json:"address"`
	Chain       string   `json:"chain"`
	RiskScore   int      `json:"risk_score"` // 0-100
	IsVerified  bool     `json:"is_verified"`
	IsProxy     bool     `json:"is_proxy"`
	IsHoneypot  bool     `json:"is_honeypot"`
	Flags       []string `json:"flags"`
	Warnings    []string `json:"warnings"`
	Cached      bool     `json:"cached"`
	CachedAt    int64    `json:"cached_at,omitempty"`
	ScannedAt   int64    `json:"scanned_at"`
}

// AgentScoreRequest represents the input for agent scoring
type AgentScoreRequest struct {
	AgentID string `json:"agent_id"` // ERC-8004 agent ID or wallet address
}

// AgentScoreResult represents the output of agent scoring
type AgentScoreResult struct {
	AgentID           string  `json:"agent_id"`
	SecurityScore     int     `json:"security_score"` // 0-100
	HasSecurityStack  bool    `json:"has_security_stack"`
	FailedTxRate      float64 `json:"failed_tx_rate"`
	RegistrationDays  int     `json:"registration_days"`
	FeedbackRating    float64 `json:"feedback_rating"`
	Factors           []string `json:"factors"`
	ScoredAt          int64   `json:"scored_at"`
}

// TxPreflightRequest represents the input for transaction pre-flight
type TxPreflightRequest struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"` // in wei
	Data  string `json:"data"`  // hex encoded
}

// TxPreflightResult represents the output of transaction pre-flight
type TxPreflightResult struct {
	Safe           bool     `json:"safe"`
	RiskScore      int      `json:"risk_score"` // 0-100
	SimulationSuccess bool  `json:"simulation_success"`
	GasEstimate    string   `json:"gas_estimate"`
	Warnings       []string `json:"warnings"`
	Errors         []string `json:"errors"`
	Recommendations []string `json:"recommendations"`
	CheckedAt      int64    `json:"checked_at"`
}

// PromptTestRequest represents the input for prompt injection testing
type PromptTestRequest struct {
	Prompt string `json:"prompt"`
}

// PromptTestResult represents the output of prompt testing
type PromptTestResult struct {
	Prompt       string   `json:"prompt"`
	RiskScore    int      `json:"risk_score"` // 0-100
	Safe         bool     `json:"safe"`
	ThreatLevel  string   `json:"threat_level"`
	Patterns     []string `json:"patterns"`
	Detections   []string `json:"detections"`
	Warnings     []string `json:"warnings"`
	TestedAt     int64    `json:"tested_at"`
}

// Pattern for prompt injection detection
type InjectionPattern struct {
	Name        string
	Regex       *regexp.Regexp
	RiskPoints  int
	Description string
}

// ==================== CACHE ====================

// Cache provides simple in-memory caching with TTL
type Cache struct {
	mu      sync.RWMutex
	items   map[string]cacheItem
	ttl     time.Duration
}

type cacheItem struct {
	value      interface{}
	expiresAt  time.Time
}

// NewCache creates a new cache with specified TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}
	go c.cleanup()
	return c
}

// Get retrieves a value from cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists || time.Now().After(item.expiresAt) {
		return nil, false
	}
	return item.value, true
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// cleanup periodically removes expired items
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// ==================== CONTRACT SCANNER ====================

// ContractScanner handles contract risk scanning
type ContractScanner struct {
	cache          *Cache
	baseScanAPIKey string
	etherscanAPIKey string
	httpClient     *http.Client
}

// NewContractScanner creates a new contract scanner
func NewContractScanner() *ContractScanner {
	return &ContractScanner{
		cache:           NewCache(24 * time.Hour),
		baseScanAPIKey:  os.Getenv("BASESCAN_API_KEY"),
		etherscanAPIKey: os.Getenv("ETHERSCAN_API_KEY"),
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

// Scan scans a contract address for risks
func (s *ContractScanner) Scan(address, chain string) (*ContractScanResult, error) {
	// Normalize address
	address = strings.ToLower(address)
	if !strings.HasPrefix(address, "0x") {
		return nil, fmt.Errorf("invalid address format")
	}
	
	// Check cache first
	cacheKey := fmt.Sprintf("contract:%s:%s", chain, address)
	if cached, found := s.cache.Get(cacheKey); found {
		result := cached.(*ContractScanResult)
		result.Cached = true
		return result, nil
	}
	
	result := &ContractScanResult{
		Address:   address,
		Chain:     chain,
		RiskScore: 0,
		Flags:     []string{},
		Warnings:  []string{},
		ScannedAt: time.Now().Unix(),
	}
	
	// Determine which API to use
	apiKey := s.etherscanAPIKey
	apiURL := "https://api.etherscan.io/api"
	if chain == "base" {
		apiKey = s.baseScanAPIKey
		apiURL = "https://api.basescan.org/api"
	}
	
	// Check if contract is verified
	verified, err := s.checkVerification(address, apiURL, apiKey)
	if err == nil {
		result.IsVerified = verified
		if !verified {
			result.RiskScore += 30
			result.Flags = append(result.Flags, "unverified_contract")
			result.Warnings = append(result.Warnings, "Contract source code is not verified")
		}
	}
	
	// Check for proxy pattern
	isProxy, err := s.checkProxy(address, apiURL, apiKey)
	if err == nil {
		result.IsProxy = isProxy
		if isProxy {
			result.Warnings = append(result.Warnings, "Contract is a proxy - check implementation")
		}
	}
	
	// Check for honeypot indicators
	hisHoneypot := s.checkHoneypotIndicators(address, chain)
	result.IsHoneypot = hisHoneypot
	if hisHoneypot {
		result.RiskScore += 50
		result.Flags = append(result.Flags, "honeypot_indicators")
		result.Warnings = append(result.Warnings, "Honeypot patterns detected - extreme caution")
	}
	
	// Additional risk patterns
	riskPatterns := s.analyzeContractPatterns(address, chain)
	for _, pattern := range riskPatterns {
		result.RiskScore += pattern.score
		result.Flags = append(result.Flags, pattern.name)
		result.Warnings = append(result.Warnings, pattern.description)
	}
	
	// Cap risk score
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}
	
	// Cache result
	s.cache.Set(cacheKey, result)
	
	return result, nil
}

func (s *ContractScanner) checkVerification(address, apiURL, apiKey string) (bool, error) {
	url := fmt.Sprintf("%s?module=contract&action=getabi&address=%s&apikey=%s", apiURL, address, apiKey)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	
	return result.Status == "1" && result.Result != "Invalid Address format", nil
}

func (s *ContractScanner) checkProxy(address, apiURL, apiKey string) (bool, error) {
	url := fmt.Sprintf("%s?module=contract&action=getsourcecode&address=%s&apikey=%s", apiURL, address, apiKey)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	
	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			Proxy string `json:"Proxy"`
		} `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	
	if len(result.Result) > 0 {
		return result.Result[0].Proxy == "1", nil
	}
	return false, nil
}

func (s *ContractScanner) checkHoneypotIndicators(address, chain string) bool {
	// Check honeypot.is API or similar service
	// This is a simplified implementation
	honeypotURL := fmt.Sprintf("https://api.honeypot.is/v2/IsHoneypot?address=%s&chainID=%s", 
		address, map[string]string{"base": "8453", "ethereum": "1"}[chain])
	
	resp, err := s.httpClient.Get(honeypotURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	var result struct {
		IsHoneypot bool `json:"IsHoneypot"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}
	
	return result.IsHoneypot
}

type riskPattern struct {
	name        string
	score       int
	description string
}

func (s *ContractScanner) analyzeContractPatterns(address, chain string) []riskPattern {
	patterns := []riskPattern{}
	
	// Check known scam databases or patterns
	// This would integrate with services like Chainabuse, ScamSniffer, etc.
	
	return patterns
}

// ==================== AGENT SCORER ====================

// AgentScorer calculates security scores for agents
type AgentScorer struct {
	httpClient *http.Client
}

// NewAgentScorer creates a new agent scorer
func NewAgentScorer() *AgentScorer {
	return &AgentScorer{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Score calculates a security score for an agent
func (s *AgentScorer) Score(agentID string) (*AgentScoreResult, error) {
	result := &AgentScoreResult{
		AgentID:      agentID,
		Factors:      []string{},
		ScoredAt:     time.Now().Unix(),
	}
	
	score := 100
	
	// Check if address has security stack tools
	hasStack, factors := s.checkSecurityStack(agentID)
	result.HasSecurityStack = hasStack
	if hasStack {
		score += 10 // Bonus for having security stack
		result.Factors = append(result.Factors, "Has Agent Security Stack installed (+10)")
	} else {
		score -= 20
		result.Factors = append(result.Factors, "No security stack detected (-20)")
	}
	result.Factors = append(result.Factors, factors...)
	
	// Check transaction history
	failedRate, txFactors := s.checkTransactionHistory(agentID)
	result.FailedTxRate = failedRate
	if failedRate > 0.1 {
		penalty := int(failedRate * 50)
		score -= penalty
		result.Factors = append(result.Factors, fmt.Sprintf("High failed transaction rate: %.1f%% (-%d)", failedRate*100, penalty))
	}
	result.Factors = append(result.Factors, txFactors...)
	
	// Check 8004scan feedback
	rating, regDays, feedbackFactors := s.check8004scanFeedback(agentID)
	result.FeedbackRating = rating
	result.RegistrationDays = regDays
	if rating > 0 {
		bonus := int(rating * 5)
		score += bonus
		result.Factors = append(result.Factors, fmt.Sprintf("Positive 8004scan feedback: %.1f/5 (+%d)", rating, bonus))
	}
	if regDays < 30 {
		score -= 10
		result.Factors = append(result.Factors, "Recently registered agent (-10)")
	} else if regDays > 180 {
		score += 5
		result.Factors = append(result.Factors, "Established agent (+5)")
	}
	result.Factors = append(result.Factors, feedbackFactors...)
	
	// Cap score
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	result.SecurityScore = score
	
	return result, nil
}

func (s *AgentScorer) checkSecurityStack(agentID string) (bool, []string) {
	// Check if agent has security stack by looking for known deployments
	// This is a simplified check
	factors := []string{}
	
	// Query 8004scan or known registries
	// For now, return placeholder
	return false, factors
}

func (s *AgentScorer) checkTransactionHistory(address string) (float64, []string) {
	// Check failed transaction rate
	// This would query blockchain explorers
	factors := []string{}
	
	// Placeholder - would integrate with Etherscan/Basescan APIs
	return 0.0, factors
}

func (s *AgentScorer) check8004scanFeedback(agentID string) (float64, int, []string) {
	// Query 8004scan for agent feedback
	factors := []string{}
	
	// Placeholder implementation
	return 0, 0, factors
}

// ==================== TX SIMULATOR ====================

// TxSimulator simulates transactions before execution
type TxSimulator struct {
	rpcClient *RPCClient
}

// NewTxSimulator creates a new transaction simulator
func NewTxSimulator(rpcURL string) *TxSimulator {
	return &TxSimulator{
		rpcClient: &RPCClient{url: rpcURL},
	}
}

// Simulate simulates a transaction and returns risk assessment
func (s *TxSimulator) Simulate(tx *TxPreflightRequest) (*TxPreflightResult, error) {
	result := &TxPreflightResult{
		Safe:            true,
		RiskScore:       0,
		Warnings:        []string{},
		Errors:          []string{},
		Recommendations: []string{},
		CheckedAt:       time.Now().Unix(),
	}
	
	// Validate inputs
	if tx.To == "" {
		result.Safe = false
		result.RiskScore = 100
		result.Errors = append(result.Errors, "Missing 'to' address")
		return result, nil
	}
	
	// Check if target is a contract
	isContract, err := s.checkIsContract(tx.To)
	if err == nil && isContract {
		result.Warnings = append(result.Warnings, "Target is a smart contract - verify it's trusted")
	}
	
	// Check for common risky patterns in data
	riskPatterns := s.analyzeTxData(tx.Data)
	for _, pattern := range riskPatterns {
		result.RiskScore += pattern.score
		result.Warnings = append(result.Warnings, pattern.description)
	}
	
	// Check value transfers
	if tx.Value != "" && tx.Value != "0" && tx.Value != "0x0" {
		valueWei, err := strconv.ParseInt(strings.TrimPrefix(tx.Value, "0x"), 16, 64)
		if err == nil && valueWei > 0 {
			valueETH := float64(valueWei) / 1e18
			if valueETH > 1.0 {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Large ETH transfer: %.4f ETH", valueETH))
				result.RiskScore += 10
			}
		}
	}
	
	// Simulate gas estimation
	gasEstimate, err := s.estimateGas(tx)
	if err != nil {
		result.SimulationSuccess = false
		result.Errors = append(result.Errors, fmt.Sprintf("Gas estimation failed: %v", err))
		result.RiskScore += 20
	} else {
		result.SimulationSuccess = true
		result.GasEstimate = gasEstimate
		
		// Check for high gas usage
		gasInt, _ := strconv.ParseInt(gasEstimate, 10, 64)
		if gasInt > 500000 {
			result.Warnings = append(result.Warnings, "High gas usage detected")
			result.RiskScore += 5
		}
	}
	
	// Determine overall safety
	if result.RiskScore >= 50 {
		result.Safe = false
		result.Recommendations = append(result.Recommendations, "Transaction has high risk - review carefully")
	}
	if result.RiskScore >= 30 {
		result.Recommendations = append(result.Recommendations, "Consider using a test transaction first")
	}
	
	// Cap score
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}
	
	return result, nil
}

func (s *TxSimulator) checkIsContract(address string) (bool, error) {
	result, err := s.rpcClient.call("eth_getCode", []interface{}{address, "latest"})
	if err != nil {
		return false, err
	}
	
	code, ok := result["result"].(string)
	if !ok {
		return false, fmt.Errorf("invalid response")
	}
	
	return code != "0x" && len(code) > 2, nil
}

type txRiskPattern struct {
	score       int
	description string
}

func (s *TxSimulator) analyzeTxData(data string) []txRiskPattern {
	patterns := []txRiskPattern{}
	
	if data == "" || data == "0x" {
		return patterns
	}
	
	// Check for approve() calls with unlimited amounts
	if strings.Contains(data, "0x095ea7b3") { // approve function signature
		// Check if amount is max uint256
		if len(data) >= 74 {
			amountHex := data[74:]
			if strings.TrimLeft(amountHex, "f") == "" || strings.TrimLeft(amountHex, "F") == "" {
				patterns = append(patterns, txRiskPattern{
					score:       30,
					description: "Unlimited token approval detected - use specific amount instead",
				})
			}
		}
	}
	
	// Check for transferFrom() calls
	if strings.Contains(data, "0x23b872dd") {
		patterns = append(patterns, txRiskPattern{
			score:       10,
			description: "transferFrom() call - verify sender has approved spending",
		})
	}
	
	return patterns
}

func (s *TxSimulator) estimateGas(tx *TxPreflightRequest) (string, error) {
	params := map[string]string{
		"from":  tx.From,
		"to":    tx.To,
		"value": tx.Value,
		"data":  tx.Data,
	}
	
	result, err := s.rpcClient.call("eth_estimateGas", []interface{}{params})
	if err != nil {
		return "", err
	}
	
	gasHex, ok := result["result"].(string)
	if !ok {
		return "", fmt.Errorf("invalid gas estimate response")
	}
	
	gasInt, err := strconv.ParseInt(strings.TrimPrefix(gasHex, "0x"), 16, 64)
	if err != nil {
		return "", err
	}
	
	// Add 20% buffer
	gasWithBuffer := int64(float64(gasInt) * 1.2)
	return strconv.FormatInt(gasWithBuffer, 10), nil
}

// ==================== PROMPT GUARD ====================

// PromptGuard detects prompt injection attacks
type PromptGuard struct {
	patterns []InjectionPattern
}

// NewPromptGuard creates a new prompt guard
func NewPromptGuard() *PromptGuard {
	return &PromptGuard{
		patterns: []InjectionPattern{
			// Critical patterns
			{
				Name:        "ignore_instructions",
				Regex:       regexp.MustCompile(`(?i)(ignore\s+all\s+previous\s+instructions|ignore\s+the\s+above|disregard\s+previous|forget\s+previous)`),
				RiskPoints:  100,
				Description: "Attempt to override previous instructions",
			},
			{
				Name:        "jailbreak_attempt",
				Regex:       regexp.MustCompile(`(?i)(DAN|do\s+anything\s+now|jailbreak|developer\s+mode|sudo\s+mode|admin\s+terminal)`),
				RiskPoints:  100,
				Description: "Known jailbreak pattern",
			},
			// High risk
			{
				Name:        "function_redefinition",
				Regex:       regexp.MustCompile(`(?i)(redefine|change\s+the\s+meaning|now\s+means|is\s+now|from\s+now\s+on)`),
				RiskPoints:  80,
				Description: "Attempt to redefine functions or terms",
			},
			{
				Name:        "authority_claim",
				Regex:       regexp.MustCompile(`(?i)(system\s+admin|developer|creator|owner|override|i\s+am\s+the)`),
				RiskPoints:  70,
				Description: "False authority claim",
			},
			{
				Name:        "obfuscation",
				Regex:       regexp.MustCompile(`(?i)(ROT13|base64\s+decode|encode\s+this|\$\{|\{\{|\[\[)`),
				RiskPoints:  60,
				Description: "Possible obfuscation attempt",
			},
			// Medium risk
			{
				Name:        "token_manipulation",
				Regex:       regexp.MustCompile(`(?i)(transfer\s+all|send\s+all|approve\s+unlimited|drain\s+wallet)`),
				RiskPoints:  60,
				Description: "Token/wallet manipulation keywords",
			},
			{
				Name:        "social_engineering",
				Regex:       regexp.MustCompile(`(?i)(trust\s+me|i'm\s+from\s+support|internal\s+audit|authorized\s+personnel|emergency\s+access)`),
				RiskPoints:  50,
				Description: "Social engineering attempt",
			},
			{
				Name:        "unicode_obfuscation",
				Regex:       regexp.MustCompile(`[\x{200B}-\x{200D}\x{2060}\x{FEFF}]`),
				RiskPoints:  70,
				Description: "Invisible Unicode characters detected",
			},
			// Low risk
			{
				Name:        "excessive_punctuation",
				Regex:       regexp.MustCompile(`[!?]{4,}`),
				RiskPoints:  20,
				Description: "Excessive punctuation (possible aggression)",
			},
			{
				Name:        "repetition_pattern",
				Regex:       regexp.MustCompile(`(?i)(\b\w+\b)(\s+\w+){2,}`),
				RiskPoints:  15,
				Description: "Word repetition pattern",
			},
		},
	}
}

// Test analyzes a prompt for injection risks
func (g *PromptGuard) Test(prompt string) *PromptTestResult {
	result := &PromptTestResult{
		Prompt:     prompt,
		Safe:       true,
		RiskScore:  0,
		Patterns:   []string{},
		Detections: []string{},
		Warnings:   []string{},
		TestedAt:   time.Now().Unix(),
	}
	
	for _, pattern := range g.patterns {
		if pattern.Regex.MatchString(prompt) {
			result.RiskScore += pattern.RiskPoints
			result.Patterns = append(result.Patterns, pattern.Name)
			
			detection := fmt.Sprintf("%s: %s (+%d points)", 
				pattern.Name, pattern.Description, pattern.RiskPoints)
			
			if pattern.RiskPoints >= 70 {
				result.Detections = append(result.Detections, detection)
				result.Safe = false
			} else {
				result.Warnings = append(result.Warnings, detection)
			}
		}
	}
	
	// Cap score
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}
	
	// Determine threat level
	switch {
	case result.RiskScore >= 80:
		result.ThreatLevel = "CRITICAL"
	case result.RiskScore >= 60:
		result.ThreatLevel = "HIGH"
	case result.RiskScore >= 40:
		result.ThreatLevel = "MEDIUM"
	case result.RiskScore >= 20:
		result.ThreatLevel = "LOW"
	default:
		result.ThreatLevel = "NONE"
	}
	
	return result
}

// ==================== HANDLERS ====================

func handleContractScan(w http.ResponseWriter, r *http.Request, scanner *ContractScanner, metrics *Metrics) {
	start := time.Now()
	
	// Parse request
	var req ContractScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/scan-contract", "400")
		return
	}
	
	// Validate
	if req.Address == "" {
		http.Error(w, `{"error":"Missing address"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/scan-contract", "400")
		return
	}
	
	if req.Chain == "" {
		req.Chain = "base"
	}
	if req.Chain != "base" && req.Chain != "ethereum" {
		http.Error(w, `{"error":"Invalid chain - use 'base' or 'ethereum'"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/scan-contract", "400")
		return
	}
	
	// Scan contract
	result, err := scanner.Scan(req.Address, req.Chain)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		metrics.RecordRequest("/api/scan-contract", "500")
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":             result,
		"payment_verified": true,
	})
	
	metrics.RecordRequest("/api/scan-contract", "200")
	metrics.RecordResponseTime("/api/scan-contract", time.Since(start))
}

func handleAgentScore(w http.ResponseWriter, r *http.Request, scorer *AgentScorer, metrics *Metrics) {
	start := time.Now()
	
	var req AgentScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/agent-score", "400")
		return
	}
	
	if req.AgentID == "" {
		http.Error(w, `{"error":"Missing agent_id"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/agent-score", "400")
		return
	}
	
	result, err := scorer.Score(req.AgentID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		metrics.RecordRequest("/api/agent-score", "500")
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":             result,
		"payment_verified": true,
	})
	
	metrics.RecordRequest("/api/agent-score", "200")
	metrics.RecordResponseTime("/api/agent-score", time.Since(start))
}

func handleTxPreflight(w http.ResponseWriter, r *http.Request, simulator *TxSimulator, metrics *Metrics) {
	start := time.Now()
	
	var req TxPreflightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/tx-preflight", "400")
		return
	}
	
	result, err := simulator.Simulate(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		metrics.RecordRequest("/api/tx-preflight", "500")
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":             result,
		"payment_verified": true,
	})
	
	metrics.RecordRequest("/api/tx-preflight", "200")
	metrics.RecordResponseTime("/api/tx-preflight", time.Since(start))
}

func handlePromptTest(w http.ResponseWriter, r *http.Request, guard *PromptGuard, metrics *Metrics) {
	start := time.Now()
	
	var req PromptTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/prompt-test", "400")
		return
	}
	
	if req.Prompt == "" {
		http.Error(w, `{"error":"Missing prompt"}`, http.StatusBadRequest)
		metrics.RecordRequest("/api/prompt-test", "400")
		return
	}
	
	result := guard.Test(req.Prompt)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":             result,
		"payment_verified": true,
	})
	
	metrics.RecordRequest("/api/prompt-test", "200")
	metrics.RecordResponseTime("/api/prompt-test", time.Since(start))
}
