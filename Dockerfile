FROM golang:1.26-alpine AS builder

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk add --no-cache git build-base olm-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-X codeberg.org/Fovir/mytrix/internal/version.Version=${VERSION}"  \
    -o bot ./cmd/bot

FROM alpine:3.21

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk add --no-cache ca-certificates olm

COPY --from=builder /app/bot /app/bot

ENTRYPOINT ["/app/bot"]
