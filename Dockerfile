# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o discord-bot cmd/bot/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates ffmpeg python3 py3-pip

# Install yt-dlp
RUN pip3 install yt-dlp

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/discord-bot .

# Copy config files
COPY .env.example .env

# Expose port (if needed for any web interfaces)
EXPOSE 8080

# Run the application
CMD ["./discord-bot"]