package config

import (
	"errors"

	"github.com/robfig/cron/v3"
)

func (mc *MytrixConfig) validate() error {
	var errs []error
	validators := []func() error{
		mc.validateGotify,
		mc.validateWakapi,
		mc.validateUmami,
	}
	for _, validator := range validators {
		if err := validator(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (mc MytrixConfig) validateCrons(crons []string) error {
	var errs []error
	for _, c := range crons {
		if _, err := cron.ParseStandard(c); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
