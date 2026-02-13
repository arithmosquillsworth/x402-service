package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MCP (Model Context Protocol) Server Implementation
// Based on the MCP specification for AI agent tool discovery

// MCPTool represents a tool that can be called by AI agents
type MCPTool struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	InputSchema MCPInputSchema  `json:"inputSchema"`
}

// MCPInputSchema defines the expected input for a tool
type MCPInputSchema struct {
	Type       string                     `json:"type"`
	Properties map[string]MCPProperty   `json:"properties"`
	Required   []string                   `json:"required"`
}

// MCPProperty defines a property in the input schema
type MCPProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// MCPServerInfo represents the MCP server metadata
type MCPServerInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Tools       []MCPTool `json:"tools"`
}

// MCPRequest represents a tool call request
type MCPRequest struct {
	Tool      string                 `json:"tool"`
	Arguments map[string]interface{} `json:"arguments"`
}

// MCPResponse represents a tool call response
type MCPResponse struct {
	Content []MCPContent `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

// MCPContent represents the content of a response
type MCPContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// handleMCPInfo returns the MCP server information and available tools
func handleMCPInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	serverInfo := MCPServerInfo{
		Name:        "Arithmos MCP Server",
		Version:     "1.0.0",
		Description: "AI agent services for Ethereum security, x402 payments, and on-chain intelligence",
		Tools: []MCPTool{
			{
				Name:        "get_gas_prices",
				Description: "Get current Ethereum gas prices in gwei",
				InputSchema: MCPInputSchema{
					Type:       "object",
					Properties: map[string]MCPProperty{},
					Required:   []string{},
				},
			},
			{
				Name:        "get_validator_queue",
				Description: "Get Ethereum validator queue status and wait times",
				InputSchema: MCPInputSchema{
					Type:       "object",
					Properties: map[string]MCPProperty{},
					Required:   []string{},
				},
			},
			{
				Name:        "scan_token",
				Description: "Scan an ERC-20 token contract for security risks and red flags",
				InputSchema: MCPInputSchema{
					Type: "object",
					Properties: map[string]MCPProperty{
						"tokenAddress": {
							Type:        "string",
							Description: "The token contract address to scan (0x...)",
						},
						"chain": {
							Type:        "string",
							Description: "Chain ID (1 for Ethereum, 8453 for Base)",
						},
					},
					Required: []string{"tokenAddress"},
				},
			},
			{
				Name:        "scan_wallet",
				Description: "Analyze wallet address for risk profile and holdings",
				InputSchema: MCPInputSchema{
					Type: "object",
					Properties: map[string]MCPProperty{
						"walletAddress": {
							Type:        "string",
							Description: "The wallet address to analyze (0x...)",
						},
						"chain": {
							Type:        "string",
							Description: "Chain ID (1 for Ethereum, 8453 for Base)",
						},
					},
					Required: []string{"walletAddress"},
				},
			},
			{
				Name:        "get_address_labels",
				Description: "Get labels and entity information for an address",
				InputSchema: MCPInputSchema{
					Type: "object",
					Properties: map[string]MCPProperty{
						"address": {
							Type:        "string",
							Description: "The address to lookup (0x...)",
						},
					},
					Required: []string{"address"},
				},
			},
			{
				Name:        "check_mev_risk",
				Description: "Check if a transaction is at risk of MEV/sandwich attacks",
				InputSchema: MCPInputSchema{
					Type: "object",
					Properties: map[string]MCPProperty{
						"txData": {
							Type:        "string",
							Description: "The transaction data/calldata",
						},
						"to": {
							Type:        "string",
							Description: "Target contract address",
						},
						"value": {
							Type:        "string",
							Description: "Transaction value in wei",
						},
					},
					Required: []string{"txData", "to"},
				},
			},
			{
				Name:        "get_eth_price",
				Description: "Get current ETH price in USD",
				InputSchema: MCPInputSchema{
					Type:       "object",
					Properties: map[string]MCPProperty{},
					Required:   []string{},
				},
			},
			{
				Name:        "check_tx_preflight",
				Description: "Run pre-flight security check on a transaction",
				InputSchema: MCPInputSchema{
					Type: "object",
					Properties: map[string]MCPProperty{
						"txData": {
							Type:        "string",
							Description: "The transaction data/calldata",
						},
						"to": {
							Type:        "string",
							Description: "Target contract address",
						},
						"value": {
							Type:        "string",
							Description: "Transaction value in wei",
						},
						"from": {
							Type:        "string",
							Description: "Sender address",
						},
					},
					Required: []string{"txData", "to", "from"},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverInfo)
}

// handleMCPCall handles tool execution requests
func handleMCPCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Invalid request: " + err.Error()}},
			IsError: true,
		})
		return
	}

	// Route to appropriate handler based on tool name
	switch req.Tool {
	case "get_gas_prices":
		handleMCPGasPrices(w, r, req.Arguments)
	case "get_validator_queue":
		handleMCPValidatorQueue(w, r, req.Arguments)
	case "scan_token":
		handleMCPScanToken(w, r, req.Arguments)
	case "scan_wallet":
		handleMCPScanWallet(w, r, req.Arguments)
	case "get_address_labels":
		handleMCPAddressLabels(w, r, req.Arguments)
	case "check_mev_risk":
		handleMCPMEVCheck(w, r, req.Arguments)
	case "get_eth_price":
		handleMCPEthPrice(w, r, req.Arguments)
	case "check_tx_preflight":
		handleMCPPreflight(w, r, req.Arguments)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Unknown tool: " + req.Tool}},
			IsError: true,
		})
	}
}

// Individual tool handlers
func handleMCPGasPrices(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	rpcClient := &RPCClient{url: getEnv("ETH_RPC_URL", "https://ethereum-rpc.publicnode.com")}
	gasData, err := rpcClient.fetchGasPrices()
	
	if err != nil {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Error fetching gas prices: " + err.Error()}},
			IsError: true,
		})
		return
	}

	result, _ := json.MarshalIndent(gasData, "", "  ")
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(result)}},
	})
}

func handleMCPValidatorQueue(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	beaconClient := NewBeaconClient()
	validatorData, err := beaconClient.fetchValidatorData()
	
	if err != nil {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Error fetching validator queue: " + err.Error()}},
			IsError: true,
		})
		return
	}

	result, _ := json.MarshalIndent(validatorData, "", "  ")
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(result)}},
	})
}

func handleMCPScanToken(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	tokenAddress, ok := args["tokenAddress"].(string)
	if !ok || tokenAddress == "" {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Missing required parameter: tokenAddress"}},
			IsError: true,
		})
		return
	}

	chain, _ := args["chain"].(string)
	if chain == "" {
		chain = "base" // Default to Base
	}

	// Use internal scan function
	result := scanToken(tokenAddress, chain)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(resultJSON)}},
	})
}

func handleMCPScanWallet(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	walletAddress, ok := args["walletAddress"].(string)
	if !ok || walletAddress == "" {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Missing required parameter: walletAddress"}},
			IsError: true,
		})
		return
	}

	chain, _ := args["chain"].(string)
	if chain == "" {
		chain = "base"
	}

	result := scanWallet(walletAddress, chain)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(resultJSON)}},
	})
}

func handleMCPAddressLabels(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	address, ok := args["address"].(string)
	if !ok || address == "" {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Missing required parameter: address"}},
			IsError: true,
		})
		return
	}

	result := lookupAddressLabel(address)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(resultJSON)}},
	})
}

func handleMCPMEVCheck(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	txData, _ := args["txData"].(string)
	to, _ := args["to"].(string)
	value, _ := args["value"].(string)
	from, _ := args["from"].(string)

	if txData == "" || to == "" {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Missing required parameters: txData and to"}},
			IsError: true,
		})
		return
	}

	req := MEVCheckRequest{
		From:  from,
		To:    to,
		Value: value,
		Data:  txData,
	}
	result := checkMEVRisk(req)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(resultJSON)}},
	})
}

func handleMCPEthPrice(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	price, err := fetchETHPrice()
	if err != nil {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Error fetching ETH price: " + err.Error()}},
			IsError: true,
		})
		return
	}

	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: fmt.Sprintf(`{"price_usd": %.2f}`, price)}},
	})
}

func handleMCPPreflight(w http.ResponseWriter, r *http.Request, args map[string]interface{}) {
	txData, _ := args["txData"].(string)
	to, _ := args["to"].(string)
	value, _ := args["value"].(string)
	from, _ := args["from"].(string)

	if txData == "" || to == "" || from == "" {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Missing required parameters: txData, to, from"}},
			IsError: true,
		})
		return
	}

	// Create TxSimulator and simulate
	rpcURL := getEnv("ETH_RPC_URL", "https://ethereum-rpc.publicnode.com")
	simulator := NewTxSimulator(rpcURL)
	req := TxPreflightRequest{
		From:  from,
		To:    to,
		Value: value,
		Data:  txData,
	}
	
	result, err := simulator.Simulate(&req)
	if err != nil {
		json.NewEncoder(w).Encode(MCPResponse{
			Content: []MCPContent{{Type: "text", Text: "Error simulating transaction: " + err.Error()}},
			IsError: true,
		})
		return
	}
	
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	
	json.NewEncoder(w).Encode(MCPResponse{
		Content: []MCPContent{{Type: "text", Text: string(resultJSON)}},
	})
}
