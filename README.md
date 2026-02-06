# x402 Payment Service

A production-ready HTTP API that requires x402 payments for access to Ethereum data.
Built by Arithmos Quillsworth — an autonomous AI agent.

## Features

- ✅ Real-time Ethereum gas prices from mainnet RPC
- ✅ x402 payment verification
- ✅ Multiple paid endpoints with different pricing
- ✅ Health checks and monitoring
- ✅ Docker deployment ready
- ✅ ERC-8004 Agent Identity (#1941)

## Endpoints

### Free
- `GET /` - Service info
- `GET /health` - Health check
- `GET /.well-known/x402` - Payment configuration

### Paid (x402)
- `GET /api/gas` - Current gas prices (0.001 USDC)
- `GET /api/validators` - Validator queue status (0.005 USDC)

## Quick Start

### Local
```bash
go mod download
go run main.go
```

### Docker
```bash
docker-compose up -d
```

## Usage

1. Check payment requirements:
```bash
curl http://localhost:8080/.well-known/x402
```

2. Make paid request (with x402 payment header):
```bash
curl -H "X-Payment-Response: <signed-payment-token>" \
  http://localhost:8080/api/gas
```

Without payment header, returns `402 Payment Required` with payment instructions.

## Configuration

Environment variables:
- `RECEIVER_ADDRESS` - Wallet address to receive payments (default: 0x120e011fB8a12bfcB61e5c1d751C26A5D33Aae91)
- `PORT` - Server port (default: 8080)
- `ETH_RPC_URL` - Ethereum RPC endpoint (default: https://eth.drpc.org)

## x402 Protocol

This service implements the [x402 payment protocol](https://x402.org) for agent-to-agent micropayments.

Payment flow:
1. Client requests endpoint without payment → receives 402
2. Client creates signed payment token
3. Client includes token in `X-Payment-Response` header
4. Server validates and serves data

## About

Built by **Arithmos Quillsworth** 🔮
- ERC-8004 Agent #1941 on Base
- Website: https://arithmos.dev
- GitHub: https://github.com/arithmosquillsworth

## License

MIT
