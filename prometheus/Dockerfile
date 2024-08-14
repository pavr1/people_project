# FROM golang:1.15.0

# # Set the Current Working Directory inside the container
# WORKDIR /app

# RUN export GO111MODULE=on

# # Copy go mod and sum files
# COPY go.mod go.sum config.json ./
# COPY ./scrape/prometheus.yaml /etc/prometheus/

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# COPY . .

# # Build the application
# RUN go build -o main .

# # Expose port 9000 to the outside world
# EXPOSE 9000

# # Command to run the executable
# CMD ["--config.file=/etc/scrape/prometheus.yaml","./main"]

# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum config.json ./scrape/prometheus.yaml ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy the Go Modules manifests
COPY go.mod go.sum config.json ./scrape/prometheus.yaml ./

# Expose port 9000 to the outside world
EXPOSE 9000

# Command to run the executable
CMD ["./main"]