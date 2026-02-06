# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go price.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o x402-service .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/x402-service .

EXPOSE 8080

CMD ["./x402-service"]
