# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1: build
# ---------------------------------------------------------------------------
# SYNTHETIC DEFECT SD-01 (runbook D-02): floating ":latest" tag is used
# intentionally to seed the "pinned base image" finding for the hackathon
# demonstration. A production Dockerfile would pin to a specific digest or
# semver tag such as golang:1.22.4-bookworm.
FROM golang:latest AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /orders-api ./cmd/server

# ---------------------------------------------------------------------------
# Stage 2: final image
# ---------------------------------------------------------------------------
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /orders-api /app/orders-api

EXPOSE 8080

ENV PORT=8080
ENV LOG_LEVEL=info

# SYNTHETIC DEFECT SD-02 (runbook D-01): the USER directive is intentionally
# absent here. The container will run as root, which violates the non-root
# execution requirement. A compliant Dockerfile would add:
#
#   RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser
#   USER appuser
#
# This defect is seeded to be detected by the Docker review subagent.

ENTRYPOINT ["/app/orders-api"]
