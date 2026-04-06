package config

import "errors"

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
