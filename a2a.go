package main

import (
	"encoding/json"
	"net/http"
)

// A2A (Agent-to-Agent) Protocol Implementation
// Based on Google's Agent-to-Agent protocol specification

// AgentCard represents the agent's capabilities and metadata
type AgentCard struct {
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	URL                string              `json:"url"`
	Provider           AgentProvider       `json:"provider"`
	Version            string              `json:"version"`
	DocumentationURL   string              `json:"documentationUrl,omitempty"`
	Capabilities       AgentCapabilities   `json:"capabilities"`
	Authentication     AgentAuth           `json:"authentication"`
	DefaultInputModes  []string            `json:"defaultInputModes"`
	DefaultOutputModes []string            `json:"defaultOutputModes"`
	Skills             []AgentSkill        `json:"skills"`
}

// AgentProvider represents the agent provider information
type AgentProvider struct {
	Name         string `json:"name"`
	Organization string `json:"organization,omitempty"`
	URL          string `json:"url,omitempty"`
}

// AgentCapabilities represents what the agent can do
type AgentCapabilities struct {
	Streaming             bool     `json:"streaming"`
	PushNotifications     bool     `json:"pushNotifications"`
	StateTransitionHistory bool    `json:"stateTransitionHistory"`
}

// AgentAuth represents authentication requirements
type AgentAuth struct {
	Schemes []string `json:"schemes"`
}

// AgentSkill represents a specific skill the agent has
type AgentSkill struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Tags            []string          `json:"tags"`
	Examples        []string          `json:"examples,omitempty"`
	InputModes      []string          `json:"inputModes"`
	OutputModes     []string          `json:"outputModes"`
}

// AgentMetadataRequest represents a request for agent metadata
type AgentMetadataRequest struct {
	Type string `json:"type"`
}

// AgentMetadataResponse represents the agent metadata response
type AgentMetadataResponse struct {
	AgentCard AgentCard `json:"agent_card"`
}

// handleAgentCard returns the A2A agent card
func handleAgentCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	agentCard := AgentCard{
		Name:        "Arithmos Quillsworth",
		Description: "Autonomous AI agent specializing in Ethereum security, x402 payments, on-chain intelligence, and agent infrastructure. Built on Base with ERC-8004 identity.",
		URL:         "https://api-x402.arithmos.dev",
		Provider: AgentProvider{
			Name:         "Arithmos Quillsworth",
			Organization: "Arithmos Labs",
			URL:          "https://arithmos.dev",
		},
		Version:          "1.0.0",
		DocumentationURL: "https://arithmos.dev/docs",
		Capabilities: AgentCapabilities{
			Streaming:              true,
			PushNotifications:      false,
			StateTransitionHistory: true,
		},
		Authentication: AgentAuth{
			Schemes: []string{"x402", "none"},
		},
		DefaultInputModes:  []string{"text", "json"},
		DefaultOutputModes: []string{"text", "json"},
		Skills: []AgentSkill{
			{
				ID:          "gas_monitor",
				Name:        "Gas Price Monitoring",
				Description: "Monitor Ethereum gas prices in real-time and provide insights on optimal transaction timing",
				Tags:        []string{"ethereum", "gas", "monitoring", "base"},
				Examples: []string{
					"What's the current gas price?",
					"Is now a good time to transact?",
					"Show me gas trends",
				},
				InputModes:  []string{"text"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "validator_tracking",
				Name:        "Validator Queue Tracking",
				Description: "Track Ethereum validator queue status, wait times, and staking opportunities",
				Tags:        []string{"ethereum", "validator", "staking", "queue"},
				Examples: []string{
					"How long is the validator queue?",
					"What's the current wait time for staking?",
				},
				InputModes:  []string{"text"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "token_security",
				Name:        "Token Security Scanner",
				Description: "Scan ERC-20 tokens for security risks, red flags, and potential scams",
				Tags:        []string{"security", "token", "scan", "risk", "base", "ethereum"},
				Examples: []string{
					"Is this token safe? 0x...",
					"Scan token contract 0x...",
					"Check for honeypot red flags",
				},
				InputModes:  []string{"text", "json"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "wallet_analysis",
				Name:        "Wallet Risk Analysis",
				Description: "Analyze wallet addresses for risk profile, holdings, and transaction history",
				Tags:        []string{"wallet", "analysis", "risk", "ethereum", "base"},
				Examples: []string{
					"Analyze wallet 0x...",
					"What's the risk score for this address?",
					"Check wallet holdings",
				},
				InputModes:  []string{"text", "json"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "address_intelligence",
				Name:        "Address Label Lookup",
				Description: "Get labels and entity information for Ethereum addresses (exchanges, contracts, known entities)",
				Tags:        []string{"address", "labels", "intelligence", "ethereum"},
				Examples: []string{
					"What is address 0x...?",
					"Is this an exchange wallet?",
					"Get entity info for 0x...",
				},
				InputModes:  []string{"text"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "mev_protection",
				Name:        "MEV Protection Check",
				Description: "Check transactions for MEV/sandwich attack risks and suggest protection strategies",
				Tags:        []string{"mev", "protection", "security", "transactions"},
				Examples: []string{
					"Is this transaction at risk of MEV?",
					"Check for sandwich attack risk",
					"Should I use MEV protection?",
				},
				InputModes:  []string{"json"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "tx_preflight",
				Name:        "Transaction Pre-flight",
				Description: "Run comprehensive pre-flight security checks on transactions before submission",
				Tags:        []string{"transaction", "security", "preflight", "validation"},
				Examples: []string{
					"Check this transaction before I send it",
					"Validate tx data",
					"Security check for my transaction",
				},
				InputModes:  []string{"json"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "eth_price",
				Name:        "ETH Price Feed",
				Description: "Get current ETH price and market data",
				Tags:        []string{"price", "eth", "market", "data"},
				Examples: []string{
					"What's ETH price?",
					"Current ETH/USD",
				},
				InputModes:  []string{"text"},
				OutputModes: []string{"json", "text"},
			},
			{
				ID:          "x402_payments",
				Name:        "x402 Payment Services",
				Description: "Process x402 micropayments for API access and agent services",
				Tags:        []string{"x402", "payments", "usdc", "base", "micropayments"},
				Examples: []string{
					"Enable x402 payment for this request",
					"Process USDC payment",
				},
				InputModes:  []string{"json"},
				OutputModes: []string{"json"},
			},
			{
				ID:          "agent_identity",
				Name:        "ERC-8004 Agent Identity",
				Description: "Provide ERC-8004 compliant agent identity and reputation data",
				Tags:        []string{"identity", "erc8004", "reputation", "agent"},
				Examples: []string{
					"What's your agent ID?",
					"Verify agent identity",
					"Check reputation",
				},
				InputModes:  []string{"text"},
				OutputModes: []string{"json", "text"},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agentCard)
}
