# Multi-stage build untuk menghasilkan image yang lebih kecil dan secure
# Stage 1: Build stage - digunakan untuk compile aplikasi Go
FROM golang:1.25-alpine AS builder

# Install dependencies yang diperlukan untuk build
# git untuk go mod download jika ada private repos
# ca-certificates untuk HTTPS connections
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory di dalam container
WORKDIR /app

# Copy go mod files terlebih dahulu untuk leverage Docker layer caching
# Jika go.mod dan go.sum tidak berubah, layer ini akan di-cache
# Ini mempercepat rebuild ketika hanya source code yang berubah
COPY go.mod go.sum ./

# Download dependencies
# Menggunakan go mod download lebih baik daripada go build langsung
# karena dependencies di-cache sebagai layer terpisah
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi
# CGO_ENABLED=0 untuk static binary yang tidak depend pada libc
# GOOS=linux untuk memastikan binary untuk Linux
# -ldflags="-w -s" untuk strip debug info dan reduce binary size
# -a untuk rebuild semua packages
# -installsuffix cgo untuk membedakan build dengan CGO_ENABLED
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# Stage 2: Runtime stage - image minimal untuk production
# Menggunakan scratch atau alpine minimal untuk keamanan dan ukuran
FROM alpine:latest

# Install ca-certificates untuk HTTPS calls
# Tanpa ini, aplikasi tidak bisa melakukan HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari builder stage
# Hanya copy binary yang sudah di-compile, tidak termasuk source code
COPY --from=builder /app/main .

# Copy timezone data jika aplikasi Anda butuh timezone
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Expose port yang digunakan aplikasi
# Ini hanya dokumentasi, tidak actually open port
EXPOSE 8080

# Health check untuk monitoring
# Docker akan periodically menjalankan command ini untuk cek status container
# Jika gagal berturut-turut, container dianggap unhealthy
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run aplikasi
# Menggunakan array format lebih baik daripada shell form
# karena signal handling lebih baik (penting untuk graceful shutdown)
CMD ["./main"]
