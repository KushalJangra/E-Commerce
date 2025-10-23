FROM golang:1.21.4 AS builder

WORKDIR /app

COPY ./pkg ./pkg

COPY ./user-service ./user-service

WORKDIR /app/user-service

RUN go mod download

RUN go build -o user-service ./cmd/main.go

RUN ls -l /app/user-service

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install libc6-compat to run Go binaries
RUN apk add --no-cache libc6-compat

# Copy the binary from the builder stage
COPY --from=builder /app/user-service/user-service .

COPY ./user-service/.env .

# Verify the binary is copied
RUN ls -l /app

# Expose the port
EXPOSE 8000

# Run the application
CMD ["./user-service"]

