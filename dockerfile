# Use the official Golang image for building the application
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from SETA to /app
COPY SETA/go.mod SETA/go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the application code from the SETA directory to /app
COPY SETA/ .

# Build the Go application
RUN go build -o main ./main.go

# Use the official PostgreSQL image
FROM postgres:14-alpine AS db

# Copy the schema.sql file into the Docker image
COPY schema/schema.sql /docker-entrypoint-initdb.d/

# Use a minimal Alpine image for the final container
FROM alpine:latest

# Set the working directory for the final image
WORKDIR /root/

# Copy the built application from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 8080

# Set environment variables for the gateways and database DSN
ENV GATEWAY_A_ENDPOINT="https://eac422c8-0b3c-400a-ba18-eab7b2b66050.mock.pstmn.io"
ENV GATEWAY_B_ENDPOINT="https://eac422c8-0b3c-400a-ba18-eab7b2b66050.mock.pstmn.io"
ENV DATABASE_DSN="postgresql://postgres@db:5432/seta?sslmode=disable"

# Run the application
CMD ["./main"]
