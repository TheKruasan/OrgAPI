# build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/app

# runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/app .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./app"]