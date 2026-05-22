FROM golang:1.25.7-alpine

WORKDIR /app

RUN apk add --no-cache git bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o app ./cmd/app

EXPOSE 8080

CMD ["/app/app"]