# Docker Setup & Usage Guide

This guide explains how to run the Order Food Online API with Docker and Docker Compose.

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
Open a new PowerShell terminal and test:

```powershell
# List products
curl -H "api_key: apitest" http://localhost:8080/product

# Health check
curl http://localhost:8080/health
```

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
  -e DB_PORT=5432 \
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

## Environment Variables
All environment variables can be set in `docker-compose.yml` or passed via `-e` flag:

| Variable | Default | Purpose |
|----------|---------|---------|
| `DB_HOST` | mariadb | Database hostname (use `mariadb` for Docker Compose) |
| `DB_PORT` | 3306 | Database port |
| `DB_USER` | mariadb | Database user |
| `DB_PASSWORD` | mariadb | Database password |
| `DB_NAME` | food_order | Database name |
| `ENV` | dev | Environment (dev/prod) |
| `SERVICE` | order-food-online | Service name |
| `LOG_LEVEL` | info | Log level (debug/info/warn/error) |

---

## Accessing mariadbQL from Host
Once the container is running:

```powershell
# Connect with mysql client
mysql -h localhost -P 3306 -u food_order -p mariadbpassword

# Or use Docker exec
docker exec -it order-food-mariadb mysql -u food_order -p mariadbpassword
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

### Execute command in container
```powershell
docker exec -it order-food-api ./order-api
docker exec -it order-food-mariadb mysql -u food_order -p
```

---

## Troubleshooting

### "Port 8080 already in use"
Change the port in `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Host:Container
```

### "Connection refused" on API startup
- Ensure mariadb is healthy before API starts
- Check `docker-compose logs mariadb`
- Verify environment variables match mariadb config

### Database migrations fail
- Ensure `migrations/` folder exists with `.sql` files
- Check file permissions: `chmod 644 migrations/*.sql`
- Verify migration SQL syntax

### Reset database (delete all data)
```powershell
docker-compose down -v
docker-compose up  # Fresh database
```

