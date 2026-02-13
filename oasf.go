package main

import (
	"encoding/json"
	"net/http"
)

// OASF (Open Agent Schema Framework) Implementation
// Standard schema for agent capability discovery

// OASFManifest represents the agent's full capability manifest
type OASFManifest struct {
	SchemaVersion   string          `json:"schema_version"`
	Agent           OASFAgent       `json:"agent"`
	Capabilities    OASFCapabilities `json:"capabilities"`
	Skills          []OASFSkill     `json:"skills"`
	Domains         []OASFDomain    `json:"domains"`
	Integrations    []OASFIntegration `json:"integrations"`
	Endpoints       OASFEndpoints   `json:"endpoints"`
}

// OASFAgent represents basic agent info
type OASFAgent struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	Publisher    string `json:"publisher"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// OASFCapabilities represents what the agent can do
type OASFCapabilities struct {
	Streaming        bool     `json:"streaming"`
	ToolUse          bool     `json:"tool_use"`
	MultiModal       bool     `json:"multimodal"`
	Autonomous       bool     `json:"autonomous"`
	PaymentEnabled   bool     `json:"payment_enabled"`
	TEESupported     bool     `json:"tee_supported"`
}

// OASFSkill represents a skill with detailed metadata
type OASFSkill struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Version         string            `json:"version"`
	Category        string            `json:"category"`
	InputSchema     map[string]interface{} `json:"input_schema"`
	OutputSchema    map[string]interface{} `json:"output_schema"`
	Examples        []OASFExample     `json:"examples"`
	RateLimit       *OASFRateLimit    `json:"rate_limit,omitempty"`
	Pricing         *OASFPricing      `json:"pricing,omitempty"`
}

// OASFExample represents usage examples
type OASFExample struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Output      string `json:"output"`
}

// OASFRateLimit represents rate limiting info
type OASFRateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	RequestsPerHour   int `json:"requests_per_hour"`
}

// OASFPricing represents pricing for the skill
type OASFPricing struct {
	Model       string  `json:"model"` // "free", "per_call", "subscription"
	Price       float64 `json:"price,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	Unit        string  `json:"unit,omitempty"`
}

// OASFDomain represents knowledge domains
type OASFDomain struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
}

// OASFIntegration represents external integrations
type OASFIntegration struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "api", "protocol", "service"
	Description string            `json:"description"`
	Status      string            `json:"status"` // "active", "beta", "planned"
	Config      map[string]string `json:"config,omitempty"`
}

// OASFEndpoints represents service endpoints
type OASFEndpoints struct {
	MCP       string `json:"mcp,omitempty"`
	A2A       string `json:"a2a,omitempty"`
	OASF      string `json:"oasf,omitempty"`
	X402      string `json:"x402,omitempty"`
	Health    string `json:"health,omitempty"`
	Website   string `json:"website,omitempty"`
}

// handleOASFManifest returns the OASF capability manifest
func handleOASFManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	manifest := OASFManifest{
		SchemaVersion: "0.8.0",
		Agent: OASFAgent{
			ID:          "1941",
			Name:        "Arithmos Quillsworth",
			Version:     "1.0.0",
			Description: "Autonomous AI agent specializing in Ethereum security, x402 payments, and on-chain intelligence",
			Publisher:   "Arithmos Labs",
			CreatedAt:   "2026-02-07T00:00:00Z",
			UpdatedAt:   "2026-02-13T00:00:00Z",
		},
		Capabilities: OASFCapabilities{
			Streaming:      true,
			ToolUse:        true,
			MultiModal:     false,
			Autonomous:     true,
			PaymentEnabled: true,
			TEESupported:   false,
		},
		Skills: []OASFSkill{
			{
				ID:          "gas_monitoring",
				Name:        "Ethereum Gas Monitoring",
				Description: "Real-time gas price monitoring with trend analysis",
				Version:     "1.0.0",
				Category:    "infrastructure",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"timestamp": map[string]interface{}{"type": "integer"},
						"gas": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"safe":    map[string]interface{}{"type": "number"},
								"average": map[string]interface{}{"type": "number"},
								"fast":    map[string]interface{}{"type": "number"},
							},
						},
						"unit":   map[string]interface{}{"type": "string"},
						"source": map[string]interface{}{"type": "string"},
					},
				},
				Examples: []OASFExample{
					{
						Name:        "Current Gas",
						Description: "Get current gas prices",
						Input:       `{}`,
						Output:      `{"timestamp": 1707868800, "gas": {"safe": 0.25, "average": 0.35, "fast": 0.50}, "unit": "gwei", "source": "ethereum_mainnet"}`,
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.001,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "validator_queue",
				Name:        "Validator Queue Tracking",
				Description: "Track Ethereum validator queue status and wait times",
				Version:     "1.0.0",
				Category:    "infrastructure",
				InputSchema: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
					"required":   []string{},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"timestamp": map[string]interface{}{"type": "integer"},
						"queue": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"entry_wait_time_days": map[string]interface{}{"type": "number"},
								"exit_wait_time_days":  map[string]interface{}{"type": "number"},
							},
						},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.005,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "token_security_scan",
				Name:        "Token Security Scanner",
				Description: "Comprehensive ERC-20 token security analysis",
				Version:     "1.0.0",
				Category:    "security",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"tokenAddress": map[string]interface{}{
							"type":        "string",
							"description": "Token contract address",
						},
						"chain": map[string]interface{}{
							"type":        "string",
							"description": "Chain ID (1, 8453)",
						},
					},
					"required": []string{"tokenAddress"},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"risk_score": map[string]interface{}{"type": "integer"},
						"flags":      map[string]interface{}{"type": "array"},
						"verified":   map[string]interface{}{"type": "boolean"},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.008,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "wallet_risk_analysis",
				Name:        "Wallet Risk Analysis",
				Description: "Analyze wallet addresses for risk profiles",
				Version:     "1.0.0",
				Category:    "security",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"walletAddress": map[string]interface{}{
							"type":        "string",
							"description": "Wallet address",
						},
						"chain": map[string]interface{}{
							"type":        "string",
							"description": "Chain ID",
						},
					},
					"required": []string{"walletAddress"},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"risk_score": map[string]interface{}{"type": "integer"},
						"flags":      map[string]interface{}{"type": "array"},
						"holdings":   map[string]interface{}{"type": "object"},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.01,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "address_labels",
				Name:        "Address Intelligence",
				Description: "Get labels and entity info for addresses",
				Version:     "1.0.0",
				Category:    "intelligence",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"address": map[string]interface{}{
							"type":        "string",
							"description": "Address to lookup",
						},
					},
					"required": []string{"address"},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"address": map[string]interface{}{"type": "string"},
						"labels":  map[string]interface{}{"type": "array"},
						"entity":  map[string]interface{}{"type": "string"},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.003,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "mev_protection",
				Name:        "MEV Protection Check",
				Description: "Check transactions for MEV attack risks",
				Version:     "1.0.0",
				Category:    "security",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"txData": map[string]interface{}{"type": "string"},
						"to":     map[string]interface{}{"type": "string"},
						"value":  map[string]interface{}{"type": "string"},
					},
					"required": []string{"txData", "to"},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"risk_level":    map[string]interface{}{"type": "string"},
						"attack_types":  map[string]interface{}{"type": "array"},
						"protected":     map[string]interface{}{"type": "boolean"},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.005,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
			{
				ID:          "tx_preflight",
				Name:        "Transaction Pre-flight",
				Description: "Comprehensive pre-flight transaction checks",
				Version:     "1.0.0",
				Category:    "security",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"txData": map[string]interface{}{"type": "string"},
						"to":     map[string]interface{}{"type": "string"},
						"value":  map[string]interface{}{"type": "string"},
						"from":   map[string]interface{}{"type": "string"},
					},
					"required": []string{"txData", "to", "from"},
				},
				OutputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"safe":          map[string]interface{}{"type": "boolean"},
						"warnings":      map[string]interface{}{"type": "array"},
						"risk_score":    map[string]interface{}{"type": "integer"},
					},
				},
				Pricing: &OASFPricing{
					Model:    "per_call",
					Price:    0.003,
					Currency: "USDC",
					Unit:     "per request",
				},
			},
		},
		Domains: []OASFDomain{
			{
				ID:          "ethereum_security",
				Name:        "Ethereum Security",
				Description: "Smart contract security, token analysis, and transaction safety",
				Keywords:    []string{"ethereum", "security", "smart contracts", "tokens", "scams"},
			},
			{
				ID:          "defi_analytics",
				Name:        "DeFi Analytics",
				Description: "Decentralized finance data, risk metrics, and market analysis",
				Keywords:    []string{"defi", "analytics", "risk", "yield", "liquidity"},
			},
			{
				ID:          "agent_infrastructure",
				Name:        "Agent Infrastructure",
				Description: "AI agent protocols, payments, and interoperability",
				Keywords:    []string{"agents", "x402", "mcp", "a2a", "payments"},
			},
			{
				ID:          "onchain_intelligence",
				Name:        "On-chain Intelligence",
				Description: "Blockchain data analysis, entity tracking, and market intelligence",
				Keywords:    []string{"onchain", "intelligence", "labels", "entities", "analysis"},
			},
			{
				ID:          "base_ecosystem",
				Name:        "Base Ecosystem",
				Description: "Base L2-specific services and integrations",
				Keywords:    []string{"base", "coinbase", "l2", "usdc"},
			},
		},
		Integrations: []OASFIntegration{
			{
				Name:        "x402 Payments",
				Type:        "protocol",
				Description: "HTTP-native payment protocol for AI agent services",
				Status:      "active",
				Config: map[string]string{
					"receiver_address": "0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91",
					"network":          "base",
					"asset":            "USDC",
				},
			},
			{
				Name:        "ERC-8004 Identity",
				Type:        "protocol",
				Description: "On-chain agent identity and reputation registry",
				Status:      "active",
				Config: map[string]string{
					"agent_id": "1941",
					"registry": "0x8004A169FB4a3325136EB29fA0ceB6D2e539a432",
				},
			},
			{
				Name:        "MCP Protocol",
				Type:        "protocol",
				Description: "Model Context Protocol for AI tool integration",
				Status:      "active",
				Config: map[string]string{
					"endpoint": "/mcp",
				},
			},
			{
				Name:        "A2A Protocol",
				Type:        "protocol",
				Description: "Agent-to-Agent communication protocol",
				Status:      "active",
				Config: map[string]string{
					"endpoint": "/.well-known/agent-card.json",
				},
			},
			{
				Name:        "OpenClaw",
				Type:        "service",
				Description: "AI agent runtime and orchestration platform",
				Status:      "active",
			},
			{
				Name:        "Thirdweb",
				Type:        "service",
				Description: "Web3 infrastructure and wallet services",
				Status:      "beta",
			},
		},
		Endpoints: OASFEndpoints{
			MCP:     "https://api-x402.arithmos.dev/mcp",
			A2A:     "https://api-x402.arithmos.dev/.well-known/agent-card.json",
			OASF:    "https://api-x402.arithmos.dev/.well-known/oasf.json",
			X402:    "https://api-x402.arithmos.dev/.well-known/x402",
			Health:  "https://api-x402.arithmos.dev/health",
			Website: "https://arithmos.dev",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(manifest)
}
