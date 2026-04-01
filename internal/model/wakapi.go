package model

import (
	"bytes"
	"fmt"
	"log/slog"
	"text/template"
)

type WakapiLanguage struct {
	Name    string  `json:"name"`
	Text    string  `json:"text"`
	Percent float32 `json:"percent"`
}

func (w WakapiLanguage) ToMarkdown() string {
	mdTmpl := "{{.Lang}} {{.Text}} {{.Percent}}"
	tmpl := template.Must(template.New(SourceWakapi).Parse(mdTmpl))
	var buf bytes.Buffer
	percent := fmt.Sprintf("%.2f%%", w.Percent)
	err := tmpl.Execute(&buf, map[string]any{
		"Lang":    w.Name,
		"Text":    w.Text,
		"Percent": percent,
	})
	if err != nil {
		slog.Error(
			"parse wakapi message to markdown failed",
			"name", w.Name,
			"text", w.Text,
			"percent", percent,
		)
		return fmt.Sprintf("Lang: %s\tText: %s\tPercent: %s", w.Name, w.Text, percent)
	}
	return buf.String()
}

type WakapiInterval string

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
