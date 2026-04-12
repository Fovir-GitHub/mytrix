// Package service contains service implementations for various integrations.
package service

import (
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/matrix"
	"github.com/Fovir-GitHub/mytrix/internal/scheduler"
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
