package config

import (
	"errors"
	"fmt"
	"slices"
)

type UmamiConfig struct {
	// Enabled determines whether to enable Umami integration.
	Enabled bool `env:"UMAMI_ENABLED" envDefault:"false"`
	// Server sets the server of Umami.
	Server string `env:"UMAMI_SERVER"`
	// Username sets the user to fetch data from.
	Username string `env:"UMAMI_USERNAME"`
	// Password stores the user's password for auth.
	Password string `env:"UMAMI_PASSWORD"`
	// DefaultInterval sets the default query interval used in commands.
	DefaultInterval string `env:"UMAMI_DEFAULT_INTERVAL" envDefault:"daily"`
	// Format sets the format of Umami reports.
	Format string `env:"UMAMI_FORMAT" envDefault:"- {{.Name}} - {{.Domain}}\n\t- Visitors: {{.Visitors}}\n\t- Visits: {{.Visits}}\n\tBounces Rate: {{.BouncesRate}}"`
	// DefaultInterval sets the time to send Umami daily report.
	DailyReportCron string `env:"UMAMI_DAILY_REPORT_CRON" envDefault:"0 9 * * *"`
	// WeeklyReportCron sets the time to send Umami weekly report.
	WeeklyReportCron string `env:"UMAMI_WEEKLY_REPORT_CRON" envDefault:"0 9 * * 1"`
	// MonthlyReportCron sets the time to send Umami monthly report.
	MonthlyReportCron string `env:"UMAMI_MONTHLY_REPORT_CRON" envDefault:"0 9 1 * *"`
	// YearlyReportCron sets the time to send Umami yearly report.
	YearlyReportCron string `env:"UMAMI_YEARLY_REPORT_CRON" envDefault:"0 9 1 1 *"`
}

func (mc *MytrixConfig) validateUmami() error {
	cfg := mc.Umami
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

	if cfg.Server == "" || cfg.Username == "" || cfg.Password == "" {
		errs = append(errs, fmt.Errorf("MYTRIX_UMAMI_SERVER, MYTRIX_UMAMI_USERNAME and MYTRIX_UMAMI_UMAMI_PASSWORD are required when MYTRIX_UMAMI_ENABLED=true"))
	}
	if !validUmamiInterval(cfg.DefaultInterval) {
		errs = append(errs, fmt.Errorf("invalid umami default interval: %s", cfg.DefaultInterval))
	}
	errs = append(errs, mc.validateCrons(crons))
	return errors.Join(errs...)
}

func validUmamiInterval(interval string) bool {
	validIntervals := []string{"daily", "weekly", "monthly", "yearly"}
	return slices.Contains(validIntervals, interval)
}
