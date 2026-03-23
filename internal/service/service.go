package service

import (
	"github.com/Fovir-GitHub/mytrix/internal/client"
	"github.com/Fovir-GitHub/mytrix/internal/http"
)

type Service struct {
	Gotify  GotifyService
	Message *MessageService
}

func NewService(httpClient *http.Client, matrixClient *client.MatrixClient) *Service {
	return &Service{
		Gotify:  newGotifyService(httpClient),
		Message: newMessageService(matrixClient),
	}
}
