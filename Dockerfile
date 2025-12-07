# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage with Tesseract
FROM alpine:3.19

WORKDIR /app

# Install Tesseract OCR and required languages
RUN apk add --no-cache \
    tesseract-ocr \
    tesseract-ocr-data-eng \
    tesseract-ocr-data-ind \
    ca-certificates \
    tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
