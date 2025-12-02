# Order Food Online — API

This repository provides a small e‑commerce service (product listing + placing orders).
This README documents shows how to run the service, environment variables, the API endpoints, example requests & responses, and developer notes.

**Project layout (key files):**
```
order-food-online/
├─ .env                     Environment variables
├─ Dockerfile               Docker image build definition
├─ docker-compose.yml       Docker Compose setup
├─ start.sh                 Waits for MariaDB, then starts API
├─ README.md                Project documentation
├─ cmd/
│ └─ api/
│ └─ main.go                Application entry point
├─ config/
│ └─ env.go # Default environment/config values
├─ coupons/ # .gz coupon files
├─ internal/
│ ├─ apperrors/
│ │ └─ apperrors.go # Custom error handling
│ ├─ db/
│ │ └─ db.go # MariaDB connection helper
│ ├─ event/
│ │ ├─ noop.go # No-op event handler
│ │ ├─ publisher.go # Event publisher
│ │ └─ types.go # Event type definitions
│ ├─ health/
│ │ ├─ handler.go # Health check HTTP handlers
│ │ └─ setup.go # Health check setup/init
│ ├─ logger/
│ │ ├─ context.go # Logger context management
│ │ ├─ logger.go # Core logging implementation
│ │ └─ sensitive.go # Sensitive data handling
│ ├─ middleware/
│ │ └─ apikey.go # API key middleware
│ ├─ product/
│ │ ├─ model.go
│ │ ├─ service.go
│ │ ├─ handler.go
│ │ ├─ repository.go # Generic repository interface
│ │ ├─ mariadb_repository.go # MariaDB implementation
│ │ └─ error.go
│ ├─ order/
│ │ ├─ model.go
│ │ ├─ service.go
│ │ ├─ handler.go
│ │ ├─ repository.go # Generic repository interface
│ │ ├─ mariadb_repository.go # MariaDB implementation
│ │ └─ error.go
│ └─ promo/
│ ├─ cache.go # Coupon caching
│ ├─ loader.go # Loading coupon data
│ ├─ model.go # Coupon models
│ └─ validator.go # Coupon validation logic
└─ migrations/
├─ 0001_create_products.up.sql
├─ 0002_create_orders.up.sql
└─ 0003_seed_products.up.sql

```


**Quick start (local)**
1. Create a Mariadb database and run migrations (see migrations/*.sql).
2. Set environment variables (see "Configuration").
3. From repo root:


```powershell
# ensure dependencies
go mod tidy

# build
go build -o bin/order-api ./cmd/api

# run
./bin/order-api
```

**Configuration (env variables)**
The service uses `internal/config` (envconfig). The following environment variables are supported (defaults shown):

- `ENV` — environment name (`dev`)
- `SERVICE` — service name (`food-order-online`)
- `DB_HOST` — Mariadb host (default `localhost`)
- `DB_PORT` — Mariadb port (default `3306`)
- `DB_USER` — Mariadb user (default `Mariadb`)
- `DB_PASSWORD` — Mariadb password (default `Mariadb`)
- `DB_NAME` — Mariadb database name (default `food_order`)
- `LOG_LEVEL` — log level (default `info`)
- `COUPON_DIR` — coupon directory

Example (PowerShell):

```powershell
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_USER = "food_user"
$env:DB_PASSWORD = "mariadbpassword"
$env:DB_NAME = "order_food"
$env:LOG_LEVEL = "debug"
./app/order-api
```


**API Key**
- Header name: `api_key`
- Expected key (per OpenAPI doc and current default middleware): `apitest`
- Middleware is applied globally. Requests without the correct header receive `401 Unauthorized`.


**Endpoints**

***Base URL:*** `http://localhost:8080`

- **GET /health** (comprehensive)
  - Description: comprehensive health check — checks DB connectivity.
  - Response (healthy): `200` JSON
    ```json
    { "status": "ok", "service": "order-food-online" }
    ```
  - Response (unhealthy): `503`
    ```json
    { "status": "unhealthy", "error": "database connection failed" }
    ```

- **GET /health/ping** (lightweight)
  - Description: simple heartbeat for load balancers.
  - Response: `200`
    ```json
    { "status": "ok" }
    ```

- **GET /product**
  - Description: list all products

  #### Scenarios

| Scenario    | Status | Notes                   |
|------------|--------|------------------------|
| All OK      | 200    | Returns list of products |
| No products | 200    | Returns empty array      |

  - Headers: `api_key: apitest`
  - Response: `200` JSON array of `Product` objects
    ```json
    [
      { "id": "10", "name": "Chicken Waffle", "price": 100.99, "category": "Waffle" }
    ]
    ```

- **GET /product/{productId}**
  - Description: get a single product by id
  - Path parameter: `productId` (string in current implementation; OpenAPI suggests integer format in spec — handler returns `400` if empty)

  #### Scenarios

| Scenario   | Status | Notes              |
|-----------|--------|------------------|
| Valid ID  | 200    | Returns product    |
| Invalid ID| 400    | Non-integer ID     |
| Not found | 404    | Product not found  |

#### Example Requests

Valid ID:
```bash
curl -X GET http://localhost:8080/product/13f6b5b2a-7f66-4b3f-9a1b-000000000000 \
--header 'api_key: mytest'
```

Invalid ID:
```bash
curl -X GET http://localhost:8080/product/3f6b5b2a-7f66-4b3f-9a1b-000000000020 \
--header 'api_key: mytest'
```

Not found:
```bash
curl -X GET http://localhost:8080/product/999 \
--header 'api_key: mytest'
```

#### Responses

Valid ID:
```json
{"id": "1", "name": "Pizza", "price": 299}
```

Invalid ID:
```json
{"error": "invalid product ID"}
```

Not found:
```json
{"error": "product not found"}
```

---

- **POST /order**
  - Description: place a new order

#### Example Requests
  #### Scenarios

| Scenario          | Status | Notes              |
|------------------|--------|------------------|
| Valid order       | 200    | Order created      |
| Missing items     | 400    | Request invalid    |
| Empty items       | 422    | Request invalid    |
| Missing productId | 400    | Request invalid    |
| Missing quantity  | 400    | Request invalid    |
| Invalid productId | 404    | Product not found  |
| Missing API Key   | 401    | Unauthorized       |
| Invalid API Key   | 403    | Forbidden          |

 Valid order:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{
  "items": [
    {"productId": "3f6b5b2a-7f66-4b3f-9a1b-000000000000", "quantity": 2},
    {"productId": "7e2d5a1b-1c3f-4d8b-9f2a-111111111111", "quantity": 1}
  ]
}'
```

Missing items:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{}'
```

Empty items:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{"items": []}'
```

Missing productId:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{"items": [{"quantity": 2}]}'
```

Missing quantity:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{"items": [{"productId": "3f6b5b2a-7f66-4b3f-9a1b-000000000000"}]}'
```

Invalid productId:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: mytest' \
--data '{"items": [{"productId": "999", "quantity": 1}]}'
```

Missing API Key:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--data '{"items": [{"productId": "3f6b5b2a-7f66-4b3f-9a1b-000000000000", "quantity": 2}]}'
```

Invalid API Key:
```bash
curl --location 'http://localhost:8080/order' \
--header 'Content-Type: application/json' \
--header 'api_key: wrongkey' \
--data '{"items": [{"productId": "3f6b5b2a-7f66-4b3f-9a1b-000000000000", "quantity": 2}]}'
```

#### Responses

Valid order:
```json
{"orderId": "12345-uuid-demo", "status": "created", "totalPrice": 747}
```

Missing items:
```json
{"error": "items field is required"}
```

Empty items:
```json
{"error": "items cannot be empty"}
```

Missing productId:
```json
{"error": "productId is required for each item"}
```

Missing quantity:
```json
{"error": "quantity is required for each item"}
```

Invalid productId:
```json
{"error": "product not found"}
```

Missing API Key:
```json
{"error": "API key missing"}
```

Invalid API Key:
```json
{"error": "invalid API key"}
```
---

**Error format**
Application errors use the `apperrors.AppError` marshaller which returns JSON in this shape:
```json
{ "code": <http-status>, "message": "text" }
```
Middleware returns the simple map `{"message":"unauthorized"}` for `401` responses.


**Database & Migrations**
- Migrations are under `migrations/`. Example files:
  - `0001_create_products.up.sql`
  - `0002_create_orders.up.sql`
  - `0003_seed_products.up.sql`

Example (using `psql`):
```powershell
mysql -h localhost -U food_user -d food_order -f migrations/0001_create_products.up.sql
mysql -h localhost -U food_user -d food_order -f migrations/0002_create_orders.up.sql
mysql -h localhost -U food_user -d food_order -f migrations/0003_seed_products.up.sql
```

**Graceful shutdown**
- The server listens for `SIGINT`/`SIGTERM` and calls `echo.Shutdown(ctx)` with a 10 second timeout so ongoing requests can finish. The DB pool is closed via `db.Close()` on exit.


**Development notes**
- To run tests (if added): `go test ./...`
- Keep `go.mod` tidy: `go mod tidy`
- Build: `go build -o bin/order-api ./cmd/api`
- Run locally: `./bin/order-api`


**Examples (curl)**
1. List products:
```bash
curl -s -H "api_key: apitest" http://localhost:8080/product | jq
```

2. Get product:
```bash
curl -s -H "api_key: apitest" http://localhost:8080/product/10 | jq
```

3. Place order:
```bash
curl -s -X POST http://localhost:8080/order \
  -H "api_key: apitest" \
  -H "Content-Type: application/json" \
  -d '{"items":[{"productId":"10","quantity":2}]}' | jq
```

---

# Docker Setup & Usage Guide
## Prerequisites
- Docker Desktop installed ([download here](https://www.docker.com/products/docker-desktop))
- Docker Compose (included with Docker Desktop)
- `curl` or Postman to test API endpoints

---
## Quick Start (Recommended)

### 1. Start Services with Docker Compose
From the repository root, run:

```bash
docker-compose up

```

This will:
- ✅ Build the Go application image
- ✅ Start mariadbQL container
- ✅ Apply database migrations automatically
- ✅ Start the API server on port 8080

You should see:
```
mariadb is running... ✓
Migrations complete. Starting API...
starting server on :8080
```

### 2. Test the API
Test the api using curl provided above

### 3. Stop Services
```powershell
docker-compose down
```

To also remove the database volume (deletes all data):
```powershell
docker-compose down -v
```

---

## Manual Docker Commands

### Build Image Manually
```powershell
docker build -t order-food-api:latest .
```

### Run API Container (with existing mariadb)
```powershell
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=3306 \
  -e DB_USER=mariadb \
  -e DB_PASSWORD=mariadb \
  -e DB_NAME=food_order \
  order-food-api:latest
```

### Start Only MariaDB Container
```powershell
docker run -d `
  --name order-food-mariadb `
  -e MARIADB_ROOT_PASSWORD=rootpassword `
  -e MARIADB_DATABASE=order_food `
  -e MARIADB_USER=food_user `
  -e MARIADB_PASSWORD=mariadbpassword `
  -p 3306:3306 `
  -p 3306:3306
```
---

## Viewing Logs
```powershell
# API logs
docker-compose logs -f api

# mariadb logs
docker-compose logs -f mariadb

# All logs
docker-compose logs -f
```

---

## Rebuilding After Code Changes
```powershell
# Rebuild the image and restart
docker-compose up --build

# Or just the API service
docker-compose up --build api
```

---

## Useful Commands

### List running containers
```powershell
docker ps
```

### Stop specific service
```powershell
docker-compose stop api
docker-compose stop mariadb
```

### Remove unused images/volumes
```powershell
docker system prune
docker volume prune
```