package model

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type UmamiInterval struct {
	Start time.Time
	End   time.Time
}

func getStartDate(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func getEndDate(t *time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

func UmamiIntervalYesterday() *UmamiInterval {
	yesterday := time.Now().AddDate(0, 0, -1)
	return &UmamiInterval{
		Start: getStartDate(&yesterday),
		End:   getEndDate(&yesterday),
	}
}

func UmamiIntervalLastWeek() *UmamiInterval {
	now := time.Now()
	weekday := int(now.Weekday())
	lastMonday := now.AddDate(0, 0, -(weekday-1)-7)
	lastSunday := lastMonday.AddDate(0, 0, 6)
	return &UmamiInterval{
		Start: getStartDate(&lastMonday),
		End:   getEndDate(&lastSunday),
	}
}

func UmamiIntervalLastMonth() *UmamiInterval {
	now := time.Now()
	lastMonthStart := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	lastMonthEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return &UmamiInterval{
		Start: getStartDate(&lastMonthStart),
		End:   getEndDate(&lastMonthEnd),
	}
}

func UmamiIntervalLastYear() *UmamiInterval {
	now := time.Now()
	lastYearStart := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	lastYearEnd := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return &UmamiInterval{
		Start: getStartDate(&lastYearStart),
		End:   getEndDate(&lastYearEnd),
	}
}

var UmamiIntervalMap = map[string]func() *UmamiInterval{
	"daily":   UmamiIntervalYesterday,
	"weekly":  UmamiIntervalLastWeek,
	"monthly": UmamiIntervalLastMonth,
	"yearly":  UmamiIntervalLastYear,
}

func ParseUmamiInterval(s string) (*UmamiInterval, error) {
	k := strings.ToLower(strings.TrimSpace(s))
	slog.Debug("parse umami interval", "original", s, "key", k)
	if v, ok := UmamiIntervalMap[k]; ok {
		return v(), nil
	}
	return nil, fmt.Errorf("invalid interval: %s", s)
}
