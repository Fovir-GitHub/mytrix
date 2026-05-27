FROM --platform=linux/amd64 golang:1.26-alpine AS builder

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk add --no-cache git build-base olm-dev musl-dev gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

RUN CGO_ENABLED=1  \
    go build -ldflags "-X codeberg.org/Fovir/mytrix/internal/version.Version=${VERSION}"  \
    -o bot ./cmd/bot

FROM alpine:3.21

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk add --no-cache ca-certificates olm

COPY --from=builder /app/bot /app/bot

ENTRYPOINT ["/app/bot"]
