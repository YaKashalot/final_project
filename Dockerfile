FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

RUN mkdir -p /app/data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
ENV TODO_PASSWORD="12345"

CMD ["./scheduler"]