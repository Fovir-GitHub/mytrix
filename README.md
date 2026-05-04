# Mytrix

A Matrix bot written in Go.

## Table of Contents

<!-- toc -->

- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Using Docker Compose](#using-docker-compose)
- [Configuration](#configuration)
  - [General Settings](#general-settings)
  - [Bot Configuration](#bot-configuration)
  - [Message Configuration](#message-configuration)
  - [Gotify Configuration](#gotify-configuration)
  - [Wakapi Configuration](#wakapi-configuration)
  - [Umami Configuration](#umami-configuration)
  - [RSS Configuration](#rss-configuration)
- [Command List](#command-list)
  - [Message](#message)
  - [Umami](#umami)
  - [Wakapi](#wakapi)
  - [RSS](#rss)
- [Development](#development)
- [Acknowledgement](#acknowledgement)

<!-- tocstop -->

## Features

- Matrix bot functionality
- [Gotify](https://gotify.net/) notification forwarding
- [Wakapi](https://github.com/muety/wakapi) report generation.
- [Umami](https://github.com/umami-software/umami) report generation.
- RSS integration.
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
    image: fovir/mytrix:latest
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
```

Then start the service:

```bash
docker-compose up -d
```

## Configuration

All configuration is done via environment variables.

### General Settings

| Variable            | Description                              | Default      |
| ------------------- | ---------------------------------------- | ------------ |
| `MYTRIX_LOG_LEVEL`  | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO`       |
| `MYTRIX_HOMESERVER` | Matrix homeserver URL                    | (required)   |
| `MYTRIX_ROOM_ID`    | Matrix room ID                           | (required)   |
| `MYTRIX_DATA_DIR`   | Data directory for storing sessions      | `data`       |
| `MYTRIX_TIMEOUT`    | HTTP request timeout in seconds          | `10`         |
| `MYTRIX_TZ`         | Timezone                                 | `time.Local` |

### Bot Configuration

| Variable                  | Description                        | Example                 |
| ------------------------- | ---------------------------------- | ----------------------- |
| `MYTRIX_BOT_USERNAME`     | Matrix username of the bot account | `mytrix`                |
| `MYTRIX_BOT_PASSWORD`     | Password for the bot account       | `123456`                |
| `MYTRIX_BOT_RECOVERY_KEY` | Encryption recovery key            | `abcd efgh ijkl`        |
| `MYTRIX_BOT_PICKLE_KEY`   | Key for encrypting crypto storage  | `random-32-byte-string` |

### Message Configuration

| Variable                    | Description                        | Default |
| --------------------------- | ---------------------------------- | ------- |
| `MYTRIX_MSG_ALLOW_MARKDOWN` | Allow Markdown in message contents | `true`  |
| `MYTRIX_MSG_ALLOW_HTML`     | Allow HTML in message contents     | `false` |

### Gotify Configuration

| Variable                | Description                              | Default                                                      |
| ----------------------- | ---------------------------------------- | ------------------------------------------------------------ |
| `MYTRIX_GOTIFY_ENABLED` | Enable Gotify forwarding                 | `false`                                                      |
| `MYTRIX_GOTIFY_SERVER`  | Gotify server URL (no scheme)            | required if `MYTRIX_GOTIFY_ENABLED=true`                     |
| `MYTRIX_GOTIFY_TOKEN`   | Gotify API token                         | required if `MYTRIX_GOTIFY_ENABLED=true`                     |
| `MYTRIX_GOTIFY_FORMAT`  | Gotify message format (support Markdown) | see [internal/config/gotify.go](./internal/config/gotify.go) |

### Wakapi Configuration

| Variable                            | Description                     | Default                                                      |
| ----------------------------------- | ------------------------------- | ------------------------------------------------------------ |
| `MYTRIX_WAKAPI_ENABLED`             | Enable Wakapi integration       | `false`                                                      |
| `MYTRIX_WAKAPI_SERVER`              | Wakapi server (no scheme)       | required if `MYTRIX_WAKAPI_ENABLED=true`                     |
| `MYTRIX_WAKAPI_API_KEY`             | API key to access Wakapi server | required if `MYTRIX_WAKAPI_ENABLED=true`                     |
| `MYTRIX_WAKAPI_USER_ID`             | User ID of Wakapi user          | `current`                                                    |
| `MYTRIX_WAKAPI_DAILY_REPORT_CRON`   | Time to send daily report       | `0 9 * * *`                                                  |
| `MYTRIX_WAKAPI_WEEKLY_REPORT_CRON`  | Time to send weekly report      | `0 9 * * 1`                                                  |
| `MYTRIX_WAKAPI_MONTHLY_REPORT_CRON` | Time to send monthly report     | `0 9 1 * *`                                                  |
| `MYTRIX_WAKAPI_YEARLY_REPORT_CRON`  | Time to send yearly report      | `0 9 1 1 *`                                                  |
| `MYTRIX_WAKAPI_LANG_FORMAT`         | Template of language format     | see [internal/config/wakapi.go](./internal/config/wakapi.go) |
| `MYTRIX_WAKAPI_DATA_FORMAT`         | Template of Wakapi data format  | see [internal/config/wakapi.go](./internal/config/wakapi.go) |

> _Tip:_ To disable reports, the cron can be set to `0 0 31 2 *` so it will not be triggered.

### Umami Configuration

| Variable                           | Description                                                     | Default                                                    |
| ---------------------------------- | --------------------------------------------------------------- | ---------------------------------------------------------- |
| `MYTRIX_UMAMI_ENABLED`             | Enable Umami integration                                        | `false`                                                    |
| `MYTRIX_UMAMI_SERVER`              | Server of Umami                                                 | required if `MYTRIX_UMAMI_ENABLED=true`                    |
| `MYTRIX_UMAMI_USERNAME`            | Umami username                                                  | required if `MYTRIX_UMAMI_ENABLED=true`                    |
| `MYTRIX_UMAMI_PASSWORD`            | Password of Umami user                                          | required if `MYTRIX_UMAMI_ENABLED=true`                    |
| `MYTRIX_UMAMI_DEFAULT_INTERVAL`    | Default query interval when no argument passed to Umami command | `daily`                                                    |
| `MYTRIX_UMAMI_FORMAT`              | Template of report format                                       | see [internal/config/umami.go](./internal/config/umami.go) |
| `MYTRIX_UMAMI_DAILY_REPORT_CRON`   | Time to send daily report                                       | `0 9 * * *`                                                |
| `MYTRIX_UMAMI_WEEKLY_REPORT_CRON`  | Time to send weekly report                                      | `0 9 * * 1`                                                |
| `MYTRIX_UMAMI_MONTHLY_REPORT_CRON` | Time to send monthly report                                     | `0 9 1 * *`                                                |
| `MYTRIX_UMAMI_YEARLY_REPORT_CRON`  | Time to send yearly report                                      | `0 9 1 1 *`                                                |

> _Tip:_ To disable reports, the cron can be set to `0 0 31 2 *` so it will not be triggered.

### RSS Configuration

| Variable                      | Description                              | Default                                                |
| ----------------------------- | ---------------------------------------- | ------------------------------------------------------ |
| `MYTRIX_RSS_ENABLED`          | Enable RSS integration                   | `false`                                                |
| `MYTRIX_RSS_FEED_FORMAT`      | Template of RSS feeds                    | see [internal/config/rss.go](./internal/config/rss.go) |
| `MYTRIX_RSS_ITEM_FORMAT`      | Template of RSS items                    | see [internal/config/rss.go](./internal/config/rss.go) |
| `MYTRIX_RSS_CRON`             | Interval of fetching RSS feeds           | `0 * * * *`                                            |
| `MYTRIX_RSS_UPDATE_AFTER_ADD` | Whether to update feeds after adding one | `true`                                                 |

## Command List

### Message

- Ping-pong

  ```txt
  !ping
  ```

- Current version

  ```txt
  !version
  ```

### Umami

```txt
!umami <interval>
```

Available intervals:

- `daily`
- `weekly`
- `monthly`
- `yearly`

> If the interval is not provided, it will fall back to the default interval defined in environment variables.

### Wakapi

```txt
!wakapi <interval>
```

Available intervals:

- `today`
- `yesterday`
- `weekly`
- `monthly`
- `yearly`
- `7d`
- `30d`
- `6m`
- `12m`
- `1y`
- `all`

> If the interval is not provided, it will fall back to the default interval defined in environment variables.

### RSS

- Add an RSS feed:

  ```txt
  !rss add <url>
  ```

- Update all RSS feeds manually:

  ```txt
  !rss update
  ```

- List all RSS feeds:

  ```txt
  !rss list
  ```

- Delete an RSS feed:

  ```txt
  !rss delete <id>
  ```

- Export RSS feeds:

  ```txt
  !rss export
  ```

- Show help information:

  ```txt
  !rss help
  ```

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
- [`gorilla/websocket`](https://github.com/gorilla/websocket): A fast, well-tested and widely used WebSocket implementation for Go.
- [`robfig/cron`](https://github.com/robfig/cron): A cron library for go
- [`go-gorm/gorm`](https://github.com/go-gorm/gorm): The fantastic ORM library for Golang, aims to be developer friendly
- [`mmcdole/gofeed`](https://github.com/mmcdole/gofeed): Parse RSS, Atom and JSON feeds in Go
