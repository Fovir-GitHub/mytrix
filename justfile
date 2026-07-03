run:
  # Run the bot.
  mkdir -p $MYTRIX_DATA_DIR
  CGO_ENABLED=0 go run -tags goolm cmd/bot/main.go

lint:
  # Run lint.
  golangci-lint run

db path:
  # Exec into database.
  litecli {{path}}

v-test:
  # Run testings.
  CGO_ENABLED=0 go test -tags goolm ./... -v

test:
  # Run testings.
  CGO_ENABLED=0 go test -tags goolm ./...

toc:
  # Generate ToC in README.md
  markdown-toc -i README.md
  prettier --write README.md

gen:
  # Generate database operation using `sqlc`
  sqlc generate
