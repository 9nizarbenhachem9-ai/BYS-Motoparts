# 0. Tahap Build Frontend (npm + Vite)
FROM node:22-alpine AS frontend-builder

WORKDIR /frontend

# Copy manifest dependensi dulu agar cache npm tidak terbuang saat source berubah
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci

# Copy source frontend lalu build ke /frontend/dist
COPY frontend/ ./
RUN npm run build

# 1. Tahap Build Backend
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod dan sum untuk dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary mandiri (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .

# 2. Tahap Runtime
FROM alpine:3.20

WORKDIR /app

# Install sertifikat CA untuk koneksi database SSL/TLS
RUN apk --no-cache add ca-certificates tzdata

# Salin biner server dari tahap builder
COPY --from=builder /app/server /app/server

# Salin hasil build frontend (Vite) — disajikan oleh biner Go saat runtime
COPY --from=frontend-builder /frontend/dist /app/frontend/

# Port default
EXPOSE 8080

# Eksekusi aplikasi
CMD ["/app/server"]
