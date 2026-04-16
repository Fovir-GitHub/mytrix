package service

import (
	"codeberg.org/Fovir/mytrix/internal/model"
)

type NoopUmamiService struct {
	err error
}

func (nu *NoopUmamiService) FetchReport(*model.UmamiInterval) string {
	return nu.err.Error()
}
