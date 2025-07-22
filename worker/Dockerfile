# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the worker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker .

# Final stage
FROM alpine:latest

# Install runtime dependencies including FFmpeg
RUN apk --no-cache add ca-certificates ffmpeg

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/worker .

# Run the worker
CMD ["./worker"] 