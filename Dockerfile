FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o mcphub ./cmd/mcphub
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o mcphub-mcp ./mcp

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/mcphub /usr/local/bin/mcphub
COPY --from=builder /app/mcphub-mcp /usr/local/bin/mcphub-mcp
ENTRYPOINT ["mcphub"]
