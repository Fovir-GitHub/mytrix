// Package config handles configuration loading and validation.
package config

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/robfig/cron/v3"
	"maunium.net/go/mautrix/id"
)

// validate validates the configuration by running all validators.
// It returns an error if any validation fails.
func (mc *MytrixConfig) validate() error {
	var errs []error
	validators := []func() error{
		mc.validateGotify,
		mc.validateWakapi,
		mc.validateUmami,
		mc.validateRSS,
		mc.validateAdminID,
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// validateCrons validates the cron expressions defined in environment variables.
// It parses each cron expression and returns an error if any are invalid.
func (mc MytrixConfig) validateCrons(crons []string) error {
	var errs []error
	for _, c := range crons {
		if _, err := cron.ParseStandard(c); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (mc *MytrixConfig) validateAdminID() error {
	idStr := mc.AdminID
	if len(idStr) <= 0 {
		slog.Warn("MYTRIX_ADMIN_ID is empty, all invitation will be ignored")
		return nil
	}

	userID := id.UserID(idStr)
	if _, _, err := userID.ParseAndValidateRelaxed(); err != nil {
		return fmt.Errorf("invalid admin ID (admin_id=%v): %w", idStr, err)
	}

	return nil
}
