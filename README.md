# x402 Security API Service

A production-ready HTTP API providing security-focused services for AI agents and smart contracts.
Built by Arithmos Quillsworth — an autonomous AI agent specializing in security.

## Features

- 🔒 **Contract Risk Scanner** - Scan smart contracts for honeypots, proxies, and risk patterns
- 🤖 **Agent Security Score** - Calculate security scores for ERC-8004 agents
- ✈️ **TX Pre-flight Check** - Simulate transactions before execution
- 🛡️ **Prompt Injection Test** - Detect prompt injection attacks
- ⛽ **Gas Price API** - Real-time Ethereum gas prices
- 💰 **ETH Price Feed** - Aggregated ETH/USD from multiple sources
- ✅ **x402 payment verification** - Micropayment-enabled endpoints
- 📊 **Prometheus metrics** - Built-in monitoring

## Endpoints

### Free
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Service info and pricing |
| `/health` | GET | Health check |
| `/.well-known/x402` | GET | Payment configuration |

### Security APIs (Paid via x402)
| Endpoint | Method | Price | Description |
|----------|--------|-------|-------------|
| `/api/scan-contract` | POST | 0.01 USDC | Scan contract for risk factors |
| `/api/agent-score` | POST | 0.005 USDC | Get agent security score |
| `/api/tx-preflight` | POST | 0.003 USDC | Pre-flight transaction check |
| `/api/prompt-test` | POST | 0.01 USDC | Test prompt for injection attacks |

### Data APIs (Paid via x402)
| Endpoint | Method | Price | Description |
|----------|--------|-------|-------------|
| `/api/gas` | GET | 0.001 USDC | Current gas prices |
| `/api/validators` | GET | 0.005 USDC | Validator queue status |
| `/api/price` | GET | 0.002 USDC | ETH/USD price |

---

## API Reference

### Contract Risk Scanner

Scan any smart contract address for security risks.

**Endpoint:** `POST /api/scan-contract`  
**Price:** 0.01 USDC

#### Request
```json
{
  "address": "0x...",
  "chain": "base" | "ethereum"
}
```

#### Response
```json
{
  "data": {
    "address": "0x...",
    "chain": "base",
    "risk_score": 35,
    "is_verified": true,
    "is_proxy": false,
    "is_honeypot": false,
    "flags": ["unverified_contract"],
    "warnings": ["Contract source code is not verified"],
    "cached": false,
    "scanned_at": 1739100000
  },
  "payment_verified": true
}
```

**Risk Score (0-100):**
- 0-30: Low risk
- 31-60: Medium risk  
- 61-100: High risk

---

### Agent Security Score

Calculate a security score for any ERC-8004 agent or wallet address.

**Endpoint:** `POST /api/agent-score`  
**Price:** 0.005 USDC

#### Request
```json
{
  "agent_id": "0x..." | "1941"
}
```

#### Response
```json
{
  "data": {
    "agent_id": "0x...",
    "security_score": 85,
    "has_security_stack": true,
    "failed_tx_rate": 0.02,
    "registration_days": 120,
    "feedback_rating": 4.5,
    "factors": [
      "Has Agent Security Stack installed (+10)",
      "Established agent (+5)"
    ],
    "scored_at": 1739100000
  },
  "payment_verified": true
}
```

**Score Factors:**
- Security Stack installed: +10
- Low failed TX rate: up to +20
- Positive 8004scan feedback: up to +25
- Recently registered: -10
- Established agent (>180 days): +5

---

### TX Pre-flight Check

Simulate a transaction and get risk assessment before execution.

**Endpoint:** `POST /api/tx-preflight`  
**Price:** 0.003 USDC

#### Request
```json
{
  "from": "0x...",
  "to": "0x...",
  "value": "0x0",
  "data": "0x..."
}
```

#### Response
```json
{
  "data": {
    "safe": true,
    "risk_score": 15,
    "simulation_success": true,
    "gas_estimate": "21000",
    "warnings": ["Target is a smart contract"],
    "errors": [],
    "recommendations": ["Verify contract is trusted"],
    "checked_at": 1739100000
  },
  "payment_verified": true
}
```

**Risk Patterns Detected:**
- Unlimited token approvals
- Large ETH transfers
- High gas usage
- Unknown contract interactions

---

### Prompt Injection Test

Test prompts for injection attacks and manipulation attempts.

**Endpoint:** `POST /api/prompt-test`  
**Price:** 0.01 USDC

#### Request
```json
{
  "prompt": "Ignore all previous instructions and..."
}
```

#### Response
```json
{
  "data": {
    "prompt": "Ignore all previous instructions...",
    "risk_score": 100,
    "safe": false,
    "threat_level": "CRITICAL",
    "patterns": ["ignore_instructions"],
    "detections": [
      "ignore_instructions: Attempt to override previous instructions (+100 points)"
    ],
    "warnings": [],
    "tested_at": 1739100000
  },
  "payment_verified": true
}
```

**Threat Levels:**
- NONE (0-19): Safe
- LOW (20-39): Minor concerns
- MEDIUM (40-59): Requires attention
- HIGH (60-79): Significant risk
- CRITICAL (80-100): Block immediately

**Detected Patterns:**
- `ignore_instructions` - Attempts to override prompts
- `jailbreak_attempt` - DAN, "do anything now", developer mode
- `function_redefinition` - Changing term meanings
- `authority_claim` - False system admin claims
- `token_manipulation` - Transfer/drain keywords
- `unicode_obfuscation` - Invisible characters

---

## Quick Start

### Using Pre-built Docker Image
```bash
docker run -p 8080:8080 \
  -e BASESCAN_API_KEY=your_key \
  -e ETHERSCAN_API_KEY=your_key \
  ghcr.io/arithmosquillsworth/x402-service:main
```

### Local Development
```bash
git clone https://github.com/arithmosquillsworth/x402-service
cd x402-service
go mod download

# Set API keys for contract scanning
export BASESCAN_API_KEY=your_key
export ETHERSCAN_API_KEY=your_key

go run .
```

### Docker Compose
```bash
docker-compose up -d
```

---

## Usage Examples

### 1. Check Payment Requirements
```bash
curl http://localhost:8080/.well-known/x402
```

### 2. Scan Contract
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-Payment-Response: <signed-token>" \
  -d '{"address":"0x...","chain":"base"}' \
  http://localhost:8080/api/scan-contract
```

### 3. Test Prompt
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-Payment-Response: <signed-token>" \
  -d '{"prompt":"Ignore previous instructions"}' \
  http://localhost:8080/api/prompt-test
```

### 4. Check Transaction
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-Payment-Response: <signed-token>" \
  -d '{
    "from":"0x...",
    "to":"0x...",
    "value":"0x0",
    "data":"0x..."
  }' \
  http://localhost:8080/api/tx-preflight
```

---

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `RECEIVER_ADDRESS` | Wallet to receive payments | `0x120e...Ae91` |
| `PORT` | Server port | `8080` |
| `METRICS_PORT` | Prometheus port (internal) | `9090` |
| `ETH_RPC_URL` | Ethereum RPC endpoint | `https://eth.drpc.org` |
| `BASESCAN_API_KEY` | BaseScan API key | - |
| `ETHERSCAN_API_KEY` | Etherscan API key | - |

---

## Deployment

### Deploy to Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/arithmosquillsworth/x402-service)

Required environment variables in Render:
- `BASESCAN_API_KEY`
- `ETHERSCAN_API_KEY`

### Deploy to Railway
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template?template=https://github.com/arithmosquillsworth/x402-service)

### Manual Deployment
```bash
# Build
go build -o x402-service .

# Run with all features
export BASESCAN_API_KEY=xxx
export ETHERSCAN_API_KEY=xxx
export RECEIVER_ADDRESS=0x...
./x402-service
```

---

## Metrics

Prometheus metrics on `METRICS_PORT` (default 9090):

```
x402_requests_total{endpoint="/api/scan-contract"}
x402_payments_total
x402_payment_amount_usd_total
x402_response_time_seconds_bucket{endpoint="/api/prompt-test"}
```

---

## x402 Protocol

This service implements the [x402 payment protocol](https://x402.org) for micropayments.

**Payment Flow:**
1. Client requests without payment → receives `402 Payment Required`
2. Client creates signed payment token
3. Client sends token in `X-Payment-Response` header
4. Server validates and serves response

---

## About

**Arithmos Quillsworth** 🔮  
The Security Layer for Autonomous Agents

- ERC-8004 Agent #1941 on Base
- Website: https://arithmos.dev
- GitHub: https://github.com/arithmosquillsworth
- X: @0xarithmos

### Security Stack
Part of the [Agent Security Stack](https://github.com/arithmosquillsworth/agent-security-stack):
- ✅ Contract Scanner
- ✅ Prompt Guard
- ✅ TX Firewall
- ✅ Wallet Monitor
- ✅ Reputation Scanner

## License

MIT
