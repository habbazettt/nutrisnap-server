# 🥗 NutriSnap Backend

Backend API untuk **NutriSnap** - platform yang memproses foto nutrition facts atau barcode untuk menghasilkan informasi nutrisi, skor kesehatan, insight, dan perbandingan antar produk.

## Tech Stack

- **Language**: Go 1.24+
- **Framework**: [Fiber](https://gofiber.io/) v2
- **Database**: PostgreSQL + GORM
- **Storage**: MinIO (S3-compatible)
- **OCR**: Tesseract
- **Containerization**: Docker & Docker Compose

## Project Structure

```
nutrisnap-server/
├── cmd/
│   └── api/
│       └── main.go          # Entry point
├── config/                   # Configuration files
├── internal/
│   ├── controllers/          # HTTP handlers
│   ├── middleware/           # Custom middleware
│   ├── models/               # Database models
│   ├── repositories/         # Data access layer
│   ├── routes/               # Route definitions
│   └── services/             # Business logic
├── pkg/
│   └── utils/                # Shared utilities
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose (optional, for full stack)

### Installation

1. Clone the repository

```bash
git clone https://github.com/habbazettt/nutrisnap-server.git
cd nutrisnap-server
```

2. Install dependencies

```bash
go mod tidy
```

3. Run the server

```bash
go run ./cmd/api/main.go
```

Server will start at `http://localhost:3000`

## API Endpoints

### Health Check

- `GET /healthz` - System health check
- `GET /v1/health` - API v1 health check

### Coming Soon

- `POST /v1/scan` - Upload nutrition facts image
- `GET /v1/scan/:id` - Get scan result
- `GET /v1/product/:barcode` - Get product by barcode
- `POST /v1/compare` - Compare two products
- `GET /v1/history` - Get scan history

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3000` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `nutrisnap` |
| `DB_PASSWORD` | PostgreSQL password | `nutrisnap_secret` |
| `DB_NAME` | PostgreSQL database | `nutrisnap_db` |
| `MINIO_ENDPOINT` | MinIO endpoint | `localhost:9000` |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` |
| `MINIO_BUCKET` | MinIO bucket name | `nutrisnap` |

## Docker Compose

### Quick Start with Docker

1. Copy environment file

```bash
cp .env.example .env
```

2. Start all services

```bash
docker-compose up -d
```

3. Check running services

```bash
docker-compose ps
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| `app` | 3000 | NutriSnap API |
| `postgres` | 5432 | PostgreSQL 15 Database |
| `adminer` | 8080 | Database Management UI |
| `minio` | 9000, 9001 | S3-Compatible Object Storage |

### Access URLs

- **API**: <http://localhost:3000>
- **API Health**: <http://localhost:3000/healthz>
- **Adminer**: <http://localhost:8080>
- **MinIO Console**: <http://localhost:9001>

## License

MIT License
