FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates docker-cli docker-cli-compose git && \
    git config --global --add safe.directory '*'
WORKDIR /app
COPY --from=builder /server .
RUN mkdir -p /app/data
EXPOSE 3001
ENTRYPOINT ["./server"]
