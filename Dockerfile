FROM golang:1.14.1 as builder
ENV GO111MODULE=on

# Prepare for custom caddy build
RUN mkdir /securityprotocol
WORKDIR /securityprotocol

COPY go.mod go.mod
COPY go.sum go.sum
# Download dependencies
RUN go mod download

COPY . /securityprotocol/
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/securityprotocol .
