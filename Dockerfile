FROM golang:1.12.7 as builder
ENV GO111MODULE=on

# Prepare for custom caddy build
RUN mkdir /securityprotocol
WORKDIR /securityprotocol

COPY . /securityprotocol/
RUN go mod download
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/securityprotocol .
