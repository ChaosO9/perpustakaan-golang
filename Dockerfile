# Use a base image with Go
FROM golang:1.23.6-alpine3.21 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest

RUN mkdir /app
COPY --from=builder /app/main /app/

RUN apk add --no-cache postgresql-client redis

EXPOSE 9000
WORKDIR /app

CMD ["./main"]
