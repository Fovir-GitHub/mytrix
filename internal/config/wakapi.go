package config

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
)

type WakapiConfig struct {
	// Enabled determines whether to enable Wakapi feature.
	Enabled bool `env:"WAKAPI_ENABLED" envDefault:"false"`
	// Server sets the Wakapi server.
	Server string `env:"WAKAPI_SERVER"`
	// APIKey sets the api key to access Wakapi API.
	APIKey string `env:"WAKAPI_API_KEY"`
	// UserID sets the user id.
	UserID string `env:"WAKAPI_USER_ID" envDefault:"current"`
	// DefaultInterval sets the default interval to fetch wakapi analysis.
	DefaultInterval string `env:"WAKAPI_DEFAULT_INTERVAL" envDefault:"today"`
	// DailyReportCron sets the time to send daily report.
	DailyReportCron string `env:"WAKAPI_DAILY_REPORT_CRON" envDefault:"0 9 * * *"`
	// MonthlyReportCron sets the time to send monthly report.
	MonthlyReportCron string `env:"WAKAPI_MONTHLY_REPORT_CRON" envDefault:"0 9 1 * *"`
	// YearlyReportCron sets the time to send yearly report.
	YearlyReportCron string `env:"WAKAPI_YEARLY_REPORT_CRON" envDefault:"0 9 1 1 *"`
}

func (mc *MytrixConfig) validateWakapi() error {
	cfg := mc.Wakapi
	if !cfg.Enabled {
		return nil
	}

	var errs []error
	crons := []string{
		cfg.DailyReportCron,
		cfg.MonthlyReportCron,
		cfg.YearlyReportCron,
	}

	if cfg.Server == "" || cfg.APIKey == "" {
		err := fmt.Errorf("MYTRIX_WAKAPI_SERVER and MYTRIX_WAKAPI_API_KEY are required when MYTRIX_WAKAPI_ENABLED=true")
		errs = append(errs, err)
	}

	for _, c := range crons {
		if _, err := cron.ParseStandard(c); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
