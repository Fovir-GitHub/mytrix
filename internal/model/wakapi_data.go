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

func (wd WakapiData) ToMarkdown() string {
	var buf bytes.Buffer
	langReport := generateLangReport(wd.Langs)
	if err := wakapiDataTmpl.Execute(&buf, dataView{
		Interval: wd.ReadableInterval,
		Lang:     langReport,
		Total:    wd.TotalTime,
	}); err != nil {
		slog.Error(
			"wakapi data to markdown failed",
			"interval", wd.ReadableInterval,
			"langReport", langReport,
			"total", wd.TotalTime,
			"err", err,
		)
		return fmt.Sprintf("Interval: %s\nLang Report: %s\nTotal: %s\n", wd.ReadableInterval, langReport, wd.TotalTime)
	}
	return buf.String()
}
