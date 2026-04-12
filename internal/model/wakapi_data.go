package model

import (
	"bytes"
	"fmt"
	"log/slog"
)

// WakapiData represents Wakapi API response data.
// It contains the total time, readable interval, and language statistics.
type WakapiData struct {
	TotalTime        string           `json:"human_readable_total"`
	ReadableInterval string           `json:"human_readable_range"`
	Langs            []WakapiLanguage `json:"languages"`
}

// dataView is a helper struct for templating Wakapi data.
// It contains the formatted fields needed for the Wakapi template.
type dataView struct {
	Interval string
	Lang     string
	Total    string
}

// toView converts WakapiData to a dataView for templating.
// It generates the language report and returns a dataView with formatted fields.
func (wd WakapiData) toView() *dataView {
	lang := generateLangReport(wd.Langs)
	return &dataView{
		Interval: wd.ReadableInterval,
		Lang:     lang,
		Total:    wd.TotalTime,
	}
}

// ToMarkdown converts WakapiData to a markdown formatted string.
// It executes the Wakapi data template with the formatted view data.
// If template execution fails, it returns a fallback formatted string.
func (wd WakapiData) ToMarkdown() string {
	var buf bytes.Buffer
	view := wd.toView()
	slog.Debug("wakapi data to view", "view", view)
	if err := wakapiDataTmpl.Execute(&buf, view); err != nil {
		slog.Error(
			"wakapi data to markdown failed",
			"interval", view.Interval,
			"lang", view.Lang,
			"total", view.Total,
			"err", err,
		)
		return fmt.Sprintf("Interval: %s\nLang Report: %s\nTotal: %s\n", view.Interval, view.Lang, view.Total)
	}
	return buf.String()
}
