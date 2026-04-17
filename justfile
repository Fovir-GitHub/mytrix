run:
  # Run the bot.
  go run -tags goolm cmd/bot/main.go

lint:
  # Run lint.
  golangci-lint run

db path:
  # Exec into database.
  litecli {{path}}

test:
  # Run testings.
  go test ./... -v

toc:
  # Generate ToC in README.md
  markdown-toc -i README.md
  prettier --write README.md
