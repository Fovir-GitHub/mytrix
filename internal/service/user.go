package service

import (
	"context"
	"fmt"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/db"
	"codeberg.org/Fovir/mytrix/internal/render"
)

type UserService struct {
	q *db.Queries
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{q: q}
}

func (u *UserService) ListUsers(ctx context.Context) (string, error) {
	users, err := u.q.AllUsers(ctx)
	if err != nil {
		return "", fmt.Errorf("query all user failed: %w", err)
	}

	var res strings.Builder
	for _, u := range users {
		res.WriteString(render.UserMarkdown(&u))
		res.WriteString("\n")
	}
	return res.String(), nil
}

func (u *UserService) IsUserAdmin(ctx context.Context, id string) (bool, error) {
	return u.q.IsUserAdmin(ctx, id)
}
