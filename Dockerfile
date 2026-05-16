FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8084


CMD ["./server"]
