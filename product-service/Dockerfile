FROM golang:1.21.4 AS builder


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o product_service ./cmd

EXPOSE 8001

CMD ["./product_service"]
