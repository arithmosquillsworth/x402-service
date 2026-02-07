# Deployment Guide

Multiple deployment options for the x402-service.

## Quick Deploy

### Railway (Recommended for testing)
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template?template=https://github.com/arithmosquillsworth/x402-service)

### Render
[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/arithmosquillsworth/x402-service)

## Manual Deployment

### Docker (Self-hosted)

```bash
# Pull and run
docker run -p 8080:8080 ghcr.io/arithmosquillsworth/x402-service:main

# Or with docker-compose
docker-compose up -d
```

### Render (CLI)

Prerequisites:
```bash
# Install Render CLI
curl -fsSL https://raw.githubusercontent.com/render-oss/render-cli/main/install.sh | bash

# Login
render login
```

Deploy:
```bash
./deploy-render.sh
```

Or manually:
1. Go to https://dashboard.render.com
2. Click "New +" → "Web Service"
3. Connect GitHub repo
4. Select "Docker" runtime
5. Deploy

### Railway (CLI)

```bash
railway login
railway init
railway up
```

### Fly.io

```bash
fly launch
curl https://raw.githubusercontent.com/superfly/flyctl/master/install.sh | sh
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `RECEIVER_ADDRESS` | `0x120e...Aae91` | Wallet to receive payments |
| `PORT` | `8080` | HTTP server port |
| `METRICS_PORT` | `9090` | Prometheus metrics port |
| `ETH_RPC_URL` | `https://eth.drpc.org` | Ethereum RPC endpoint |
| `PRICE_API_KEY` | - | CoinGecko API key (optional) |

## Health Checks

All deployments include a health check endpoint:
```bash
curl https://your-service.com/health
```

## Monitoring

The service exposes Prometheus metrics on `METRICS_PORT` (default 9090):
- `x402_uptime_seconds`
- `x402_requests_total`
- `x402_payments_total`
- `x402_response_time_seconds`

⚠️ **Security Note**: Do not expose the metrics port publicly. Use internal networking or a reverse proxy.

## Troubleshooting

### Service won't start
Check logs:
```bash
# Docker
docker logs <container-id>

# Render
render logs x402-service

# Railway
railway logs
```

### 402 Payment Required
This is expected! The service requires x402 payment headers for paid endpoints.

### Connection refused
Ensure the port is correctly mapped:
- Container: 8080
- Host: your chosen port

## Production Checklist

- [ ] Set custom `RECEIVER_ADDRESS` (your wallet)
- [ ] Configure `PRICE_API_KEY` for better rate limits
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Enable HTTPS (automatic on most platforms)
- [ ] Configure CORS if needed
- [ ] Test payment flow before going live
