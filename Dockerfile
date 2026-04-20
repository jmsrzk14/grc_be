# Build Stage
FROM golang:1.22-alpine AS builder

# Install build tools jika diperlukan
RUN apk add --no-cache git

WORKDIR /src
COPY . .

# Download dependencies
RUN go mod download

# Build binary dengan flags untuk optimasi ukuran
RUN go build -ldflags="-s -w" -o /app/main ./cmd/grc_be/main.go

# Final Stage (Lightweight image)
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy binary
COPY --from=builder /app/main .
# Copy configs (opsional, karena nanti kita akan gunakan env vars)
COPY --from=builder /src/configs ./configs

# Port standar Railway
ENV PORT=8000
EXPOSE 8000

# Jalankan aplikasi
CMD ["./main", "-conf", "configs"]
