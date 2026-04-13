package service

import (
	"codeberg.org/Fovir/mytrix/internal/model"
)

type NoopUmamiService struct {
	err error
}

func (nu *NoopUmamiService) getToken() (string, error) {
	return "", nu.err
}

func (nu *NoopUmamiService) fetchWebsites() ([]*model.UmamiWebsite, error) {
	return nil, nu.err
}

func (nu *NoopUmamiService) fetchWebsiteStat(*model.UmamiWebsite, *model.UmamiInterval) (*model.UmamiWebsiteStat, error) {
	return nil, nu.err
}

func (nu *NoopUmamiService) fetchWebsiteData(*model.UmamiInterval) ([]*model.UmamiWebsite, error) {
	return nil, nu.err
}

func (nu *NoopUmamiService) generateReport([]*model.UmamiWebsite) string {
	return nu.err.Error()
}

func (nu *NoopUmamiService) FetchReport(*model.UmamiInterval) string {
	return nu.err.Error()
}
