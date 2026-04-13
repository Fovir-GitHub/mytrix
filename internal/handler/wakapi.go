// Package handler handles incoming events and commands.
package handler

import (
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/config"
	"codeberg.org/Fovir/mytrix/internal/model"
)

// getWakapiInterval extracts the Wakapi interval from the given message string.
// It splits the message into fields and checks the second field for a valid interval.
// If no interval is provided or the provided interval is invalid, it falls back to
// the configured default interval. It returns the interval and an error if parsing fails,
// otherwise it returns the interval and nil.
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
	slog.Debug("got wakapi interval", "interval", interval)
	return interval, nil
}
