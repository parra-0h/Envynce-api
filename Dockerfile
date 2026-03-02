# --- Stage 1: Build ---
FROM golang:1.21-alpine AS builder

# Install gcc for cgo (needed by gorm/postgres)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o envynce-api ./cmd/api

# --- Stage 2: Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/envynce-api .

EXPOSE 8080

ENTRYPOINT ["./envynce-api"]
