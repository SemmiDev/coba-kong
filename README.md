# Kong API Gateway with Go Microservice

A comprehensive example project demonstrating how to set up **Kong API Gateway** with a **Go microservice** using Docker Compose. This project includes service registration, rate limiting, CORS, request transformation, and a web-based management UI (Konga).

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Client Request                             │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Kong API Gateway (:8000)                         │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐     │
│  │Rate Limiting│ │    CORS     │ │ Request     │ │   Logging   │     │
│  │   Plugin    │ │   Plugin    │ │ Transformer │ │   Plugin    │     │
│  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘     │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      User Service (:8080)                            │
│                         (Go + Gin)                                   │
│                    REST API: /api/v1/users                           │
└─────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     PostgreSQL Database                              │
└─────────────────────────────────────────────────────────────────────┘
```

## 📦 Services

| Service              | Port | Description                       |
| -------------------- | ---- | --------------------------------- |
| **Kong Proxy**       | 8000 | API Gateway proxy endpoint        |
| **Kong Admin API**   | 8001 | Kong management REST API          |
| **Kong HTTPS Proxy** | 8443 | HTTPS proxy endpoint              |
| **Konga Dashboard**  | 1337 | Web UI for Kong management        |
| **User Service**     | 8080 | Go microservice (direct access)   |
| **PostgreSQL (App)** | 5432 | Database for user service         |
| **Kong Database**    | -    | PostgreSQL for Kong configuration |
| **Konga Database**   | -    | PostgreSQL 9.6 for Konga          |

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- curl or httpie (for API testing)

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd belajar-kong

# Copy environment file
cp .env.example .env
```

### 2. Start All Services

```bash
# Start all containers
docker compose up -d

# Wait for all services to be healthy (about 30-60 seconds)
docker compose ps
```

### 3. Initialize Konga Database

Konga requires database migration before first use:

```bash
docker run --rm --network kong-network \
  pantsel/konga:latest \
  -c prepare \
  -a postgres \
  -u postgresql://konga:kongapass@konga-database:5432/konga
```

### 4. Apply Kong Configuration

Use **decK** to sync your declarative configuration:

```bash
docker run --rm --network kong-network \
  -v $(pwd)/kong:/kong \
  kong/deck:latest gateway sync /kong/kong.yml \
  --kong-addr http://kong-gateway:8001
```

### 5. Verify Setup

```bash
# Check all containers are running
docker compose ps

# Test Kong proxy
curl http://localhost:8000/health

# Test user service through Kong
curl http://localhost:8000/api/v1/users
```

## 📁 Project Structure

```
belajar-kong/
├── compose.yaml          # Docker Compose configuration
├── docker/
│   └── Dockerfile        # Go service Dockerfile
├── kong/
│   └── kong.yml          # Kong declarative configuration
├── main.go               # Application entry point
├── config.go             # Configuration management
├── user.go               # User model
├── user_handler.go       # HTTP handlers
├── logger.go             # Logging middleware
├── response.go           # Response helpers
├── play.http             # HTTP client test file
├── Makefile              # Development commands
├── .env.example          # Environment template
└── README.md             # This file
```

## 🔧 Configuration

### Kong Declarative Configuration (`kong/kong.yml`)

The Kong configuration is managed declaratively via `kong.yml`:

```yaml
_format_version: '3.0'

services:
    - name: user-service
      url: http://user-service:8080
      routes:
          - name: user-api-route
            paths:
                - /api/v1/users
          - name: user-health-route
            paths:
                - /health

plugins:
    - name: rate-limiting
      service: user-service
      config:
          minute: 100
          policy: local
          limit_by: ip

    - name: cors
      service: user-service
      config:
          origins: ['*']
          methods: [GET, POST, PUT, DELETE, PATCH, OPTIONS]
          credentials: true

    - name: request-transformer
      service: user-service
      config:
          add:
              headers:
                  - X-Gateway:Kong
```

### Applying Configuration Changes

After modifying `kong/kong.yml`, apply changes with:

```bash
# Using Docker (recommended)
docker run --rm --network kong-network \
  -v $(pwd)/kong:/kong \
  kong/deck:latest gateway sync /kong/kong.yml \
  --kong-addr http://kong-gateway:8001

# Or if decK is installed locally
deck sync --kong-addr http://localhost:8001 --state kong/kong.yml
```

## 🧪 API Testing

### Using curl

```bash
# Health check
curl http://localhost:8000/health

# Get all users
curl http://localhost:8000/api/v1/users

# Create a user
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get user by ID
curl http://localhost:8000/api/v1/users/{id}

# Update user
curl -X PUT http://localhost:8000/api/v1/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "John Updated", "email": "john.updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8000/api/v1/users/{id}
```

### Using play.http

Open `play.http` in VS Code with [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension or JetBrains IDE to run HTTP requests directly.

### Using Makefile

```bash
make test-api      # Run API tests
make kong-health   # Check Kong health
make kong-services # List registered services
make kong-routes   # List registered routes
```

## 🛠️ Development

### Running Locally (without Docker)

```bash
# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable
export PORT=8080

# Run the application
go run .

# Or use air for hot reload
air
```

### Make Commands

```bash
make help           # Show all available commands
make build          # Build the Go application
make run            # Run locally
make test           # Run tests
make docker-up      # Start all services
make docker-down    # Stop all services
make docker-logs    # View logs
make kong-config    # Apply Kong configuration
make fmt            # Format code
make lint           # Run linter
```

## 🔌 Kong Plugins

### Enabled Plugins

| Plugin                  | Description                     | Configuration                |
| ----------------------- | ------------------------------- | ---------------------------- |
| **rate-limiting**       | Limits requests per time window | 100 req/min per IP           |
| **cors**                | Cross-Origin Resource Sharing   | Allow all origins            |
| **request-transformer** | Modify request headers          | Add `X-Gateway: Kong`        |
| **file-log**            | Log requests to file            | `/tmp/kong-user-service.log` |

### Adding More Plugins

Edit `kong/kong.yml` to add plugins:

```yaml
plugins:
    # Key Authentication
    - name: key-auth
      service: user-service
      config:
          key_names: [apikey]
          hide_credentials: true

    # Basic Authentication
    - name: basic-auth
      service: user-service
      config:
          hide_credentials: true

    # IP Restriction
    - name: ip-restriction
      service: user-service
      config:
          allow:
              - 10.0.0.0/8
              - 172.16.0.0/12
```

## 🖥️ Konga Dashboard

Access the Konga web UI at **http://localhost:1337**

### First-Time Setup

1. Open http://localhost:1337
2. Create an admin account
3. Add a Kong connection:
    - **Name**: `kong-local`
    - **Kong Admin URL**: `http://kong-gateway:8001`

### Features

- Visual service and route management
- Plugin configuration UI
- Consumer management
- Real-time monitoring

## 🔍 Monitoring & Debugging

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f kong
docker compose logs -f user-service
docker compose logs -f konga
```

### Check Service Status

```bash
# Kong status
curl http://localhost:8001/status

# Kong services
curl http://localhost:8001/services | jq

# Kong routes
curl http://localhost:8001/routes | jq

# Kong plugins
curl http://localhost:8001/plugins | jq
```

### Rate Limit Headers

Check rate limiting in response headers:

```bash
curl -i http://localhost:8000/health
# Look for: X-RateLimit-Limit-Minute, X-RateLimit-Remaining-Minute
```

## 🧹 Cleanup

```bash
# Stop all services
docker compose down

# Stop and remove volumes (deletes all data)
docker compose down -v

# Remove unused Docker resources
docker system prune -f
```

## ⚠️ Troubleshooting

### Konga keeps restarting

Konga's pg driver doesn't support PostgreSQL 15's SCRAM-SHA-256 authentication. The solution is to use a separate PostgreSQL 9.6 instance for Konga (already configured in `compose.yaml`).

```bash
# Check Konga logs
docker logs konga-dashboard

# Re-run database preparation
docker run --rm --network kong-network \
  pantsel/konga:latest \
  -c prepare \
  -a postgres \
  -u postgresql://konga:kongapass@konga-database:5432/konga
```

### Kong routes not working

1. Verify routes are registered:

    ```bash
    curl http://localhost:8001/routes | jq
    ```

2. Re-sync configuration:
    ```bash
    docker run --rm --network kong-network \
      -v $(pwd)/kong:/kong \
      kong/deck:latest gateway sync /kong/kong.yml \
      --kong-addr http://kong-gateway:8001
    ```

### Container can't connect to another container

Ensure all containers are on the same network:

```bash
docker network inspect kong-network
```

## 📚 Resources

- [Kong Documentation](https://docs.konghq.com/)
- [decK Documentation](https://docs.konghq.com/deck/)
- [Konga GitHub](https://github.com/pantsel/konga)
- [Gin Web Framework](https://gin-gonic.com/)

## 📝 License

MIT License - feel free to use this project for learning and production.
