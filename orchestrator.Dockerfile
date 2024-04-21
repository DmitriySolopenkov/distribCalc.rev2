FROM golang:1.22-alpine AS builder
WORKDIR /build
RUN apk update && apk add --no-cache git
COPY .env ./
COPY go.* ./
COPY internal/orchestrator ./internal/orchestrator
COPY pkg ./pkg
COPY docs ./docs
RUN go mod download
COPY cmd/orchestrator ./
RUN go build -o ./orchestrator .

FROM scratch
USER 1000

COPY --from=builder --chown=1000 /build/orchestrator /orchestrator
COPY --from=builder --chown=1000 /build/.env /.env

ENTRYPOINT ["/orchestrator"]

EXPOSE 8000
