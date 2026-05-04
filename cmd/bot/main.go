package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"codeberg.org/Fovir/mytrix/internal/bot"
	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/logger"
	"codeberg.org/Fovir/mytrix/internal/model"
	_ "maunium.net/go/mautrix/crypto/goolm"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logger.Init()
	slog.Info("mytrix started", "version", "v0.2.6")
	config.SetTimeZone()
	model.InitTemplates()

	b, err := bot.New()
	if err != nil {
		slog.Error(
			"failed to create bot",
			"err", err,
		)
		os.Exit(1)
	}

	ctx := context.Background()
	if err = b.Start(ctx); err != nil {
		slog.Error(
			"fail to run bot",
			"err", err,
		)
		os.Exit(1)
	}
}
