# 1. Tahap Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod dan sum untuk dependency caching
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

# Copy seluruh source code backend
COPY backend/ .

# Build binary mandiri (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

# 2. Tahap Runtime (Image ringan dan aman)
FROM alpine:3.20

WORKDIR /app

# Install sertifikat CA untuk koneksi database SSL/TLS
RUN apk --no-cache add ca-certificates tzdata

# Salin biner server dari tahap builder
COPY --from=builder /app/server /app/server

# Salin seluruh aset antarmuka frontend
COPY frontend/ /app/frontend/

# Hostless akan menentukan variabel PORT secara dinamis
EXPOSE 8080

# Eksekusi aplikasi
CMD ["/app/server"]

