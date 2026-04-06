package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Fovir-GitHub/mytrix/internal/http"
)

func (ru *RealUmamiService) createURL(path string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   ru.server,
		Path:   path,
	}
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
		return "", fmt.Errorf("marshal payload failed: %w", err)
	}

	u := ru.createURL("/api/auth/login")
	req, err := http.NewRequest(
		http.MethodPost,
		u.String(),
		bytes.NewReader(bodyData),
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return "", fmt.Errorf("umami create get token request failed: %w", err)
	}
	if err := ru.c.DoJSON(req, &data); err != nil {
		return "", fmt.Errorf("get json failed: %w", err)
	}
	slog.Debug("get umami token", "token", data.Token)
	return data.Token, nil
}
