# Start from the official Go image
FROM golang:1.17-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o /go/bin/app

# Start a new stage from scratch
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the previous stage
COPY --from=builder /go/bin/app .

# Expose port 8080 to the outside world
EXPOSE 9080

# Command to run the executable
CMD ["./app"]