FROM golang:1.26-alpine AS builder

RUN apk update && apk upgrade --no-cache

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

RUN CGO_ENABLED=0  \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build -ldflags "-X codeberg.org/Fovir/mytrix/internal/version.Version=${VERSION}"  \
    -tags goolm \
    -o bot ./cmd/bot

FROM alpine:3.21

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bot /app/bot

ENTRYPOINT ["/app/bot"]
