package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Fovir-GitHub/mytrix/internal/bot"
	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/logger"
	"github.com/Fovir-GitHub/mytrix/internal/model"
	_ "maunium.net/go/mautrix/crypto/goolm"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logger.Init()
	slog.Info("main start")
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
