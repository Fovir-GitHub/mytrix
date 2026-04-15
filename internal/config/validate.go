// Package config handles configuration loading and validation.
package config

import (
	"errors"

	"github.com/robfig/cron/v3"
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
