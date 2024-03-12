FROM golang:1.21-alpine AS builder
WORKDIR /usr/local/src
COPY ./ ./

RUN go mod download
RUN go build -o main ./cmd/main/main.go

EXPOSE 8081

CMD ["/main"]