package service

import (
	"log/slog"

	"github.com/Fovir-GitHub/mytrix/internal/http"
	"github.com/Fovir-GitHub/mytrix/internal/matrix"
)

type Service struct {
	Gotify  GotifyService
	Message *MessageService
}

func NewService(httpClient *http.Client, matrixClient *client.Client) *Service {
	slog.Debug("create services")
	return &Service{
		Gotify:  newGotifyService(httpClient),
		Message: newMessageService(matrixClient),
	}
}
