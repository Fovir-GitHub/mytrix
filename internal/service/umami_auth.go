package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	myhttp "codeberg.org/Fovir/mytrix/internal/http"
)

func (ru *RealUmamiService) createURL(path string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   ru.server,
		Path:   path,
	}
}

func (ru *RealUmamiService) setAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ru.token))
}

func (ru *RealUmamiService) getToken() (string, error) {
	var data struct {
		Token string `json:"token"`
	}
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: ru.username,
		Password: ru.password,
	}
	bodyData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("get umami token failed: %w", err)
	}

	u := ru.createURL("/api/auth/login")
	req, err := myhttp.NewRequest(
		myhttp.MethodPost,
		u.String(),
		bytes.NewReader(bodyData),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return "", fmt.Errorf("fetch umami token failed: %w", err)
	}
	if err := ru.c.DoJSON(req, &data); err != nil {
		return "", fmt.Errorf("fetch umami token failed: %w", err)
	}
	slog.Debug("umami token fetched")
	return data.Token, nil
}

func (ru *RealUmamiService) IsTokenValid() bool {
	u := ru.createURL("/api/auth/verify")
	req, err := myhttp.NewRequest(myhttp.MethodPost,
		u.String(),
		nil,
		map[string]string{"Accept": "application/json"},
	)
	if err != nil {
		return false
	}
	ru.setAuthHeader(req)
	resp, err := ru.c.Do(req)
	if err != nil {
		return false
	}
	valid := resp.StatusCode == http.StatusOK
	return valid
}

func (ru *RealUmamiService) updateToken() {
	t, err := ru.getToken()
	if err != nil {
		return
	}
	ru.token = t
	slog.Info("umami token updated")
}
