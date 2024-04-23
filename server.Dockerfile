FROM golang:1.22 AS orchestrator
WORKDIR /app
COPY .env ./
COPY go.* ./
COPY internal/orchestrator ./internal/orchestrator
COPY pkg ./pkg
COPY docs ./docs
RUN go mod download
COPY cmd/orchestrator ./
RUN go build -o orchestrator .
EXPOSE 8000
CMD ["./orchestrator"]