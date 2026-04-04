package handler

import (
	"log/slog"
	"strings"

	"github.com/Fovir-GitHub/mytrix/internal/config"
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

func (h *Handler) fetchWakapiReport(interval model.WakapiInterval) (string, error) {
	langs, err := h.service.Wakapi.FetchLanguages(interval)
	if err != nil {
		slog.Error("fetch wakapi data failed", "err", err)
		return "", err
	}

	slog.Debug("received wakapi languages", "len", len(langs))
	if len(langs) <= 0 {
		return "no data", nil
	}

	var msg string
	for _, lang := range langs {
		msg += (lang.ToMarkdown() + "\n")
	}
	msg = strings.TrimSpace(msg)
	msg = "```text\n" + msg + "\n```"
	return msg, nil
}

func getWakapiInterval(msg string) (model.WakapiInterval, error) {
	defaultInterval := config.Config.Wakapi.DefaultInterval
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		slog.Warn("no interval provided, fallback to default interval", "defaultInterval", defaultInterval)
		return model.WakapiInterval(defaultInterval), nil
	}
	intervalStr := parts[1]
	interval, err := model.ParseWakapiInterval(intervalStr)
	if err != nil {
		slog.Warn(
			"parse wakapi interval failed, use default interval",
			"msg", msg,
			"intervalStr", intervalStr,
			"defaultInterval", defaultInterval,
			"err", err,
		)
		return model.WakapiInterval(defaultInterval), err
	}
	slog.Info("got wakapi interval", "interval", interval)
	return interval, nil
}
