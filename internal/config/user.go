package config

import (
	"fmt"
	"log/slog"

	"maunium.net/go/mautrix/id"
)

type UserConfig struct {
	// AdminID determines who can invite the bot.
	AdminID string `env:"USER_ADMIN_ID,required"`
	// Format is the message style when listing users.
	Format string `env:"USER_FORMAT" envDefault:"{{.ID}} - {{.Role}}"`
}

func (mc *MytrixConfig) validateAdminID() error {
	idStr := mc.User.AdminID
	if len(idStr) <= 0 {
		slog.Warn("MYTRIX_USER_ADMIN_ID is empty, all invitation will be ignored")
		return nil
	}

	userID := id.UserID(idStr)
	if _, _, err := userID.ParseAndValidateRelaxed(); err != nil {
		return fmt.Errorf("invalid admin ID (admin_id=%v): %w", idStr, err)
	}

	return nil
}
