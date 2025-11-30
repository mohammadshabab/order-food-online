FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-api ./cmd/api

# --- Runtime Image ---
FROM alpine:3.19

RUN apk add --no-cache mysql-client

WORKDIR /app

COPY --from=builder /app/order-api /app/
COPY start.sh /app/start.sh
COPY migrations /migrations

RUN chmod +x /app/start.sh

CMD ["sh", "/app/start.sh"]
