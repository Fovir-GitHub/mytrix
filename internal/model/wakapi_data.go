package model

import (
	"bytes"
	"fmt"
	"log/slog"
)

type WakapiData struct {
	TotalTime        string           `json:"human_readable_total"`
	ReadableInterval string           `json:"human_readable_range"`
	Langs            []WakapiLanguage `json:"languages"`
}

type dataView struct {
	Interval string
	Lang     string
	Total    string
}

func (wd WakapiData) toView() *dataView {
	lang := generateLangReport(wd.Langs)
	return &dataView{
		Interval: wd.ReadableInterval,
		Lang:     lang,
		Total:    wd.TotalTime,
	}
}

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
