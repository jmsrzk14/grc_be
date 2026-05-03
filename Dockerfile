FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/grc_be/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Copy the configs directory
COPY --from=builder /app/configs ./configs

# Expose the port (Railway will provide the PORT env var)
EXPOSE 8000

# Run the binary
CMD ["./main", "-conf", "configs"]
