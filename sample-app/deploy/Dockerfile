# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1: build
# ---------------------------------------------------------------------------
FROM golang:1.22.4-bookworm AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /orders-api ./cmd/server

# ---------------------------------------------------------------------------
# Stage 2: final image
# ---------------------------------------------------------------------------
FROM debian:bookworm-slim

RUN addgroup --system appgroup \
    && adduser --system --ingroup appgroup --no-create-home appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /orders-api /app/orders-api

EXPOSE 8080

ENV PORT=8080
ENV LOG_LEVEL=info

USER appuser

ENTRYPOINT ["/app/orders-api"]
