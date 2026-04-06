package service

import (
	"github.com/Fovir-GitHub/mytrix/internal/model"
)

type NoopUmamiService struct {
	err error
}

func (nu *NoopUmamiService) getToken() (string, error) {
	return "", nu.err
}

func (nu *NoopUmamiService) FetchWebsites() ([]model.UmamiWebsite, error) {
	return nil, nu.err
}
