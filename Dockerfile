FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base olm-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o bot ./cmd/bot

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates olm

COPY --from=builder /app/bot /app/bot

ENTRYPOINT ["/app/bot"]
