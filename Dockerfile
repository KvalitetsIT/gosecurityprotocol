FROM golang:1.12.7 as builder
ENV GO111MODULE=on

# Prepare for custom caddy build
RUN mkdir /securityprotocol
WORKDIR /securityprotocol
RUN go mod init securityprotocol
RUN go get gopkg.in/mgo.v2

# Kitcaddy module source
COPY . /securityprotocol/
RUN go test securityprotocol
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/securityprotocol .
