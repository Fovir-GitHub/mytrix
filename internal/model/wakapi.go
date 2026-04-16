package model

import (
	"fmt"
	"log/slog"
	"strings"
)

// WakapiInterval represents a time interval for Wakapi data retrieval.
type WakapiInterval string

// Wakapi interval constants representing different time periods for data retrieval.
// These constants map to the API endpoint paths used by Wakapi.
const (
	WakapiIntervalToday        WakapiInterval = "today"
	WakapiIntervalYesterday    WakapiInterval = "yesterday"
	WakapiIntervalWeek         WakapiInterval = "week"
	WakapiIntervalMonth        WakapiInterval = "month"
	WakapiIntervalYear         WakapiInterval = "year"
	WakapiInterval7Days        WakapiInterval = "7_days"
	WakapiIntervalLast7Days    WakapiInterval = "last_7_days"
	WakapiInterval30Days       WakapiInterval = "30_days"
	WakapiIntervalLast30Days   WakapiInterval = "last_30_days"
	WakapiInterval6Months      WakapiInterval = "6_months"
	WakapiIntervalLast6Months  WakapiInterval = "last_6_months"
	WakapiInterval12Months     WakapiInterval = "12_months"
	WakapiIntervalLast12Months WakapiInterval = "last_12_months"
	WakapiIntervalLastYear     WakapiInterval = "last_year"
	WakapiIntervalAny          WakapiInterval = "any"
	WakapiIntervalAllTime      WakapiInterval = "all_time"
)

// wakapiIntervalMap maps string representations to WakapiInterval constants.
// It is used to parse user-provided interval strings into the correct enum values.
// Note: "monthly" maps to 30 days, and various shorthand notations map to their full representations.
var wakapiIntervalMap = map[string]WakapiInterval{
	"today":     WakapiIntervalToday,
	"yesterday": WakapiIntervalYesterday,
	"weekly":    WakapiIntervalWeek,
	"monthly":   WakapiInterval30Days,
	"yearly":    WakapiIntervalYear,
	"7d":        WakapiIntervalLast7Days,
	"30d":       WakapiInterval30Days,
	"6m":        WakapiIntervalLast6Months,
	"12m":       WakapiIntervalLast12Months,
	"1y":        WakapiIntervalLastYear,
	"all":       WakapiIntervalAllTime,
}

// ParseWakapiInterval parses the given string into a WakapiInterval constant.
// It converts the input to lowercase and trims spaces before looking it up in the interval map.
// If the string matches a known interval, it returns the corresponding WakapiInterval value.
// Otherwise, it returns an error indicating an invalid interval.
func ParseWakapiInterval(s string) (WakapiInterval, error) {
	k := strings.ToLower(strings.TrimSpace(s))
	slog.Debug("parse wakapi interval", "original", s, "key", k)
	if v, ok := wakapiIntervalMap[k]; ok {
		return v, nil
	}
	return "", fmt.Errorf("invalid interval: %s", s)
}
