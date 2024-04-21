FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk update && apk add --no-cache git
COPY .env ./
COPY go.* ./
COPY pkg ./pkg
COPY internal/agent ./internal/agent
RUN go mod download
COPY cmd/agent ./cmd/agent
RUN go build -o /build/agent cmd/agent/main.go

FROM scratch
USER 1000

COPY --from=builder --chown=1000 /build/agent /agent
COPY --from=builder --chown=1000 /build/.env /.env

ENTRYPOINT ["/agent"]

EXPOSE 8000
