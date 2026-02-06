# x402 Service Client

Client examples for interacting with the x402 payment-enabled API.

## Quick Start

### Using the test script

```bash
# Test local instance
./test.sh

# Test remote instance
X402_API_URL=https://your-deployed-url ./test.sh
```

### Using the Go client

```bash
# Run with local instance
go run example.go

# Run with remote instance
X402_API_URL=https://your-deployed-url go run example.go
```

## Endpoints

| Endpoint | Price | Description |
|----------|-------|-------------|
| `/health` | Free | Health check |
| `/.well-known/x402` | Free | Payment configuration |
| `/api/gas` | 0.001 USDC | Current gas prices |
| `/api/validators` | 0.005 USDC | Validator queue status |
| `/api/price` | 0.002 USDC | ETH/USD price |

## Making Paid Requests

To access paid endpoints, include an x402 payment token:

```bash
curl -H "X-Payment-Response: <your-payment-token>" \
  http://localhost:8080/api/gas
```

See the [main README](../README.md) for payment token format.

## Docker Image

```bash
docker run -p 8080:8080 ghcr.io/arithmosquillsworth/x402-service:main
```
