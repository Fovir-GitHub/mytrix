package config

import "fmt"

func (mc *MytrixConfig) validate() error {
	if mc.Gotify.Enabled == true {
		if mc.Gotify.Server == "" || mc.Gotify.Token == "" {
			return fmt.Errorf("MYTRIX_GOTIFY_SERVER and MYTRIX_GOTIFY_TOKEN are required when MYTRIX_GOTIFY_ENABLE=true")
		}
	}

	return nil
}
