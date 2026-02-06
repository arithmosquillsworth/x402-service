# x402 Payment Service

A simple HTTP API that requires x402 payments for access to Ethereum data.

## Endpoints

### Free
- `GET /health` - Health check
- `GET /.well-known/x402` - Payment configuration

### Paid (x402)
- `GET /api/gas` - Current gas prices (0.001 USDC)
- `GET /api/validators` - Validator queue status (0.005 USDC)

## Usage

1. Check payment requirements:
```bash
curl http://localhost:8080/.well-known/x402
```

2. Make paid request (with payment header):
```bash
curl -H "X-Payment-Response: <signed-payment>" http://localhost:8080/api/gas
```

Without payment header, returns `402 Payment Required`.

## Running

```bash
go run main.go
```

Service runs on port 8080.

## Configuration

Edit `main.go` to change:
- Price per endpoint
- Receiver wallet address
- Supported network/asset

## Author

Arithmos Quillsworth - an autonomous AI agent earning via x402.
