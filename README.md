# Mytrix

A Matrix bot written in Go.

## Features

- Matrix bot functionality
- [Gotify](https://gotify.net/) notification forwarding
- End-to-end encryption support
- Session persistence

## Installation

### Prerequisites

- Go 1.26+ (for building from source)
- Docker (optional, for containerized deployment)
- Nix (optional, for reproducible builds)

### Using Docker Compose

Create a `compose.yml` file:

```yaml
services:
  mytrix:
    image: ghcr.io/fovir-github/mytrix:main
    container_name: mytrix
    volumes:
      - ./mytrix:/data
    environment:
      - MYTRIX_LOG_LEVEL=INFO
      - MYTRIX_HOMESERVER=https://matrix.example.com
      - MYTRIX_ROOM_ID=!roomid:matrix.example.com
      - MYTRIX_BOT_USERNAME=bot@example.com
      - MYTRIX_BOT_PASSWORD=your-bot-password
      - MYTRIX_BOT_RECOVERY_KEY=your-recovery-key
      - MYTRIX_BOT_PICKLE_KEY=random-32-byte-string
      - MYTRIX_GOTIFY_ENABLED=true
      - MYTRIX_GOTIFY_SERVER=gotify.example.com
      - MYTRIX_GOTIFY_TOKEN=your-gotify-token
```

Then start the service:

```bash
docker-compose up -d
```

## Configuration

All configuration is done via environment variables.

### General Settings

| Variable            | Description                              | Default    |
| ------------------- | ---------------------------------------- | ---------- |
| `MYTRIX_LOG_LEVEL`  | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO`     |
| `MYTRIX_HOMESERVER` | Matrix homeserver URL                    | (required) |
| `MYTRIX_ROOM_ID`    | Matrix room ID                           | (required) |
| `MYTRIX_DATA_DIR`   | Data directory for storing sessions      | `data`     |
| `MYTRIX_TIMEOUT`    | HTTP request timeout in seconds          | `10`       |

### Bot Configuration

| Variable                  | Description                        | Example                 |
| ------------------------- | ---------------------------------- | ----------------------- |
| `MYTRIX_BOT_USERNAME`     | Matrix username of the bot account | `mytrix`                |
| `MYTRIX_BOT_PASSWORD`     | Password for the bot account       | `123456`                |
| `MYTRIX_BOT_RECOVERY_KEY` | Encryption recovery key            | `abcd efgh ijkl`        |
| `MYTRIX_BOT_PICKLE_KEY`   | Key for encrypting crypto storage  | `random-32-byte-string` |

### Message Configuration

| Variable                    | Description                        | Example                     |
| --------------------------- | ---------------------------------- | --------------------------- |
| `MYTRIX_MSG_ALLOW_MARKDOWN` | Allow Markdown in message contents | `true` (default) or `false` |
| `MYTRIX_MSG_ALLOW_HTML`     | Allow HTML in message contents     | `false` (default) or `true` |

### Gotify Configuration

| Variable                | Description                              | Default                                                      |
| ----------------------- | ---------------------------------------- | ------------------------------------------------------------ |
| `MYTRIX_GOTIFY_ENABLED` | Enable Gotify forwarding                 | false                                                        |
| `MYTRIX_GOTIFY_SERVER`  | Gotify server URL (for WebSocket)        | required if `MYTRIX_GOTIFY_ENABLED=true`                     |
| `MYTRIX_GOTIFY_TOKEN`   | Gotify API token                         | required if `MYTRIX_GOTIFY_ENABLED=true`                     |
| `MYTRIX_GOTIFY_FORMAT`  | Gotify message format (support Markdown) | see [internal/config/gotify.go](./internal/config/gotify.go) |

## Development

0. Clone and enter the repository:

```bash
git clone https://github.com/Fovir-GitHub/mytrix.git
cd mytrix
```

1. Allow [`direnv`](https://github.com/direnv/direnv):

```bash
direnv allow
```

2. Start the server using [`just`](https://github.com/casey/just):

```bash
just run
```

## Acknowledgement

- [`mautrix/go`](https://github.com/mautrix/go): A Golang Matrix framework.
- [`caarlos0/env`](https://github.com/caarlos0/env): A simple, zero-dependencies library to parse environment variables into structs
- [Blog by Dominik Chrástecký](https://chrastecky.dev/programming/creating-a-simple-encrypted-matrix-bot-in-go): Tutorial for this project.
- [gorilla/websocket](https://github.com/gorilla/websocket): A fast, well-tested and widely used WebSocket implementation for Go.
