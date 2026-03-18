FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

# (Опционально) копируем миграции или статику
# COPY --from=builder /app/migrations ./migrations

EXPOSE ${GO_PORT}

CMD ["./main"]