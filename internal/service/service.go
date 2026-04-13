// Package service contains service implementations for various integrations.
package service

import (
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/http"
	"codeberg.org/Fovir/mytrix/internal/matrix"
	"codeberg.org/Fovir/mytrix/internal/scheduler"
)

type Service struct {
	Gotify  *GotifyService
	Message *MessageService
	Umami   UmamiService
	Wakapi  WakapiService
}

func NewService(httpClient *http.Client, matrixClient *matrix.Client, schedulerClient *scheduler.Scheduler) *Service {
	slog.Debug("create services")
	return &Service{
		Gotify:  newGotifyService(),
		Message: newMessageService(matrixClient),
		Umami:   newUmamiService(httpClient),
		Wakapi:  newWakapiService(httpClient, schedulerClient),
	}
}
