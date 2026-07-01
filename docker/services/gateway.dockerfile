FROM ecommerce/ecommerce-base:latest AS build
COPY account account
COPY product product
COPY order order
COPY payment payment
COPY gateway gateway
COPY pkg pkg
RUN GO111MODULE=on go build -mod mod -o /go/bin/app ./gateway/cmd/gateway

FROM alpine:3.20
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]
