package config

import (
	"errors"
	"fmt"
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
	// WeeklyReportCron sets the time to send weekly report.
	WeeklyReportCron string `env:"WAKAPI_WEEKLY_REPORT_CRON" envDefault:"0 9 * * 1"`
	// MonthlyReportCron sets the time to send monthly report.
	MonthlyReportCron string `env:"WAKAPI_MONTHLY_REPORT_CRON" envDefault:"0 9 1 * *"`
	// YearlyReportCron sets the time to send yearly report.
	YearlyReportCron string `env:"WAKAPI_YEARLY_REPORT_CRON" envDefault:"0 9 1 1 *"`
	// LangFormat sets the template of language report.
	LangFormat string `env:"WAKAPI_LANG_FORMAT" envDefault:"{{.Lang}} {{.Text}} {{.Percent}}"`
	// DataFormat sets the template of Wakapi data report.
	DataFormat string `env:"WAKAPI_DATA_FORMAT" envDefault:"{{.Interval}} Report\n\n{{.Lang}}\n\nTotal: {{.Total}}"`
}

func (mc *MytrixConfig) validateWakapi() error {
	cfg := mc.Wakapi
	if !cfg.Enabled {
		return nil
	}

	var errs []error
	crons := []string{
		cfg.DailyReportCron,
		cfg.WeeklyReportCron,
		cfg.MonthlyReportCron,
		cfg.YearlyReportCron,
	}

	if cfg.Server == "" || cfg.APIKey == "" {
		err := fmt.Errorf("MYTRIX_WAKAPI_SERVER and MYTRIX_WAKAPI_API_KEY are required when MYTRIX_WAKAPI_ENABLED=true")
		errs = append(errs, err)
	}
	errs = append(errs, mc.validateCrons(crons))
	return errors.Join(errs...)
}
