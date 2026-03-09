# ============================================
# Stage 1: Build binary menggunakan golang:alpine
# ============================================
FROM golang:alpine AS builder

WORKDIR /build

# Copy go module files dan download dependencies
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || true

# Copy source code dan build static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o zenvault main.go

# ============================================
# Stage 2: Scratch image (super kecil)
# ============================================
FROM scratch

# Copy binary dari builder stage
COPY --from=builder /build/zenvault /zenvault

# Expose port
EXPOSE 8080

# Run
ENTRYPOINT ["/zenvault"]
