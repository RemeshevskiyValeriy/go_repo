FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd

RUN go build -o tasks ./cmd/tasks

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/tasks /app/tasks

EXPOSE 8082

CMD ["/app/tasks"]