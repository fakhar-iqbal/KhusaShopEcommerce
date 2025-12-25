# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Adjust the path to your main file if it's different
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates needed for making HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy environment file if needed (usually env vars are set in Railway dashboard, but keeping .env example is fine)
# COPY .env . 

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
