# Order Food Online — API

This repository provides a small e‑commerce service (product listing + placing orders).
This README documents how to run the service, environment variables, the API endpoints, example requests & responses, and developer notes.

**Project layout (key files):**
- `cmd/api/main.go` — application entry point
- `internal/product` — product models, service, handler, repository
- `internal/order` — order models, service, handler, repository
- `internal/health` — health check endpoints
- `internal/middleware/apikey.go` — API key middleware
- `internal/db/db.go` — Mariadb connection helper
- `migrations/` — SQL migrations

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

Example (PowerShell):

```powershell
$env:DB_HOST = "localhost"
$env:DB_PORT = "3306"
$env:DB_USER = "food_user"
$env:DB_PASSWORD = "mariadbpassword"
$env:DB_NAME = "order_food"
$env:LOG_LEVEL = "debug"
./bin/order-api
```


**API Key**
- Header name: `api_key`
- Expected key (per OpenAPI doc and current default middleware): `apitest`
- Middleware is applied globally. Requests without the correct header receive `401 Unauthorized`.


**Endpoints**
Base URL: `http://localhost:8080`

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
  - Headers: `api_key: apitest`
  - Responses:
    - `200` — product found
      ```json
      { "id": "10", "name": "Chicken Waffle", "price": 9.99, "category": "Waffle" }
      ```
    - `400` — invalid id supplied
      ```json
      { "code": 400, "message": "invalid ID supplied" }
      ```
    - `404` — not found
      ```json
      { "code": 404, "message": "product not found" }
      ```

- **POST /order**
  - Description: place a new order
  - Headers: `api_key: apitest`, `Content-Type: application/json`
  - Request body: `OrderReq` object
    ```json
    {
      "couponCode": "SUMMER10",   // optional
      "items": [
        { "productId": "10", "quantity": 2 }
      ]
    }
    ```
  - Responses:
    - `200` — order created
      ```json
      {
        "id": "0000-0000-0000-0000",
        "items": [ { "productId": "10", "quantity": 2 } ],
        "products": [ /*product refs if available */ ],
        "couponCode": "SUMMER10"
      }
      ```
    - `400` — invalid request body (binding error)
      ```json
      { "code": 400, "message": "invalid request body" }
      ```
    - `422` — validation failure (domain validation)
      ```json
      { "code": 422, "message": "order must have at least one item" }
      ```
    - `401` — missing/invalid API key
      ```json
      { "message": "unauthorized" }
      ```


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
