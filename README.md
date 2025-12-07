# 🥗 NutriSnap Backend

Backend API untuk **NutriSnap** - platform yang memproses foto nutrition facts atau barcode untuk menghasilkan informasi nutrisi, skor kesehatan, insight, dan perbandingan antar produk.

## Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.24+ |
| **Framework** | Fiber v2 |
| **Database** | PostgreSQL 15 + GORM |
| **Storage** | MinIO (S3-compatible) |
| **OCR** | Tesseract |
| **Docs** | Swagger/OpenAPI |
| **Monitoring** | Prometheus + Grafana |
| **Container** | Docker & Docker Compose |

## Project Structure

```
nutrisnap-server/
├── cmd/api/                  # Entry point
├── config/                   # Configuration loader
├── docs/                     # Swagger documentation
├── internal/
│   ├── bootstrap/            # App initialization
│   ├── controllers/          # HTTP handlers
│   ├── dto/                  # Data transfer objects
│   ├── middleware/           # Custom middleware
│   ├── models/               # Database models
│   ├── repositories/         # Data access layer
│   ├── routes/               # Route definitions
│   └── services/             # Business logic
├── monitoring/
│   ├── prometheus/           # Prometheus config
│   └── grafana/              # Grafana provisioning
├── pkg/
│   ├── database/             # Database connection
│   ├── logger/               # Structured logging
│   └── response/             # API response helpers
└── docker-compose.yml
```

## Quick Start

### With Docker (Recommended)

```bash
# Clone repository
git clone https://github.com/habbazettt/nutrisnap-server.git
cd nutrisnap-server

# Copy environment file
cp .env.example .env

# Start all services
docker-compose up -d

# Check services
docker-compose ps
```

### Without Docker

```bash
# Install dependencies
go mod tidy

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs

# Run server
go run ./cmd/api/main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Health check |
| GET | `/api/v1/health` | API v1 health check |
| GET | `/docs/*` | Swagger UI |
| GET | `/metrics` | Prometheus metrics |

### Coming Soon (EPIC 2-9)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | User registration |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/scan` | Upload nutrition image |
| GET | `/api/v1/scan/:id` | Get scan result |
| GET | `/api/v1/product/:barcode` | Get product by barcode |
| POST | `/api/v1/compare` | Compare products |
| GET | `/api/v1/history` | Scan history |

## Services

| Service | Port | Description |
|---------|------|-------------|
| **app** | 3000 | NutriSnap API |
| **postgres** | 5432 | PostgreSQL 15 |
| **adminer** | 8080 | Database UI |
| **minio** | 9010, 9011 | Object Storage |
| **prometheus** | 9090 | Metrics Collection |
| **grafana** | 3001 | Metrics Dashboard |

## Access URLs

| Service | URL |
|---------|-----|
| API | <http://localhost:3000> |
| Swagger Docs | <http://localhost:3000/docs> |
| Metrics | <http://localhost:3000/metrics> |
| Adminer | <http://localhost:8080> |
| MinIO Console | <http://localhost:9011> |
| Prometheus | <http://localhost:9090> |
| Grafana | <http://localhost:3001> |

### Grafana Login

- **Username**: admin
- **Password**: admin

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | Server port |
| `ENV` | Environment (development/production) |
| `LOG_LEVEL` | Log level (debug/info/warn/error) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_USER` | PostgreSQL user |
| `DB_PASSWORD` | PostgreSQL password |
| `DB_NAME` | PostgreSQL database |
| `MINIO_ENDPOINT` | MinIO endpoint |
| `MINIO_ACCESS_KEY` | MinIO access key |
| `MINIO_SECRET_KEY` | MinIO secret key |
| `GRAFANA_USER` | Grafana admin user |
| `GRAFANA_PASSWORD` | Grafana admin password |

## Features

- ✅ Structured logging with slog
- ✅ Rate limiting (100 req/min default)
- ✅ API response envelope
- ✅ Swagger/OpenAPI documentation
- ✅ Prometheus metrics
- ✅ Grafana dashboard (auto-provisioned)
- ✅ GORM with auto-migration
- ✅ Graceful shutdown

## License

MIT License - see [LICENSE](LICENSE)
