package render

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"codeberg.org/Fovir/mytrix/internal/model"
)

type wakapiDataView struct {
	Interval string
	Lang     string
	Total    string
}

type wakapiLangView struct {
	Lang    string
	Text    string
	Percent string
}

func WakapiDataMarkdown(wd *model.WakapiData) string {
	var buf bytes.Buffer
	view := wakapiDataToView(wd)
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

func WakapiLanguageMarkdown(wl *model.WakapiLanguage) string {
	var buf bytes.Buffer
	percent := fmt.Sprintf("%.2f%%", wl.Percent)
	if err := wakapiLangTmpl.Execute(&buf, wakapiLangView{
		Lang:    wl.Name,
		Text:    wl.Text,
		Percent: percent,
	}); err != nil {
		slog.Error(
			"parse wakapi message to markdown failed",
			"name", wl.Name,
			"text", wl.Text,
			"percent", percent,
			"err", err,
		)
		return fmt.Sprintf("Lang: %s\tText: %s\tPercent: %s", wl.Name, wl.Text, percent)
	}
	return buf.String()
}

func wakapiDataToView(wd *model.WakapiData) *wakapiDataView {
	lang := generateLangReport(wd.Langs)
	return &wakapiDataView{
		Interval: wd.ReadableInterval,
		Lang:     lang,
		Total:    wd.TotalTime,
	}
}

func generateLangReport(langs []model.WakapiLanguage) string {
	slog.Debug("generate language report begin", "len", len(langs))
	if len(langs) <= 0 {
		slog.Warn("no lanuage found, return empty report", "len", len(langs))
		return ""
	}

	var msg strings.Builder
	msg.WriteString("```text\n")
	for _, lang := range langs {
		msg.WriteString(WakapiLanguageMarkdown(&lang))
		msg.WriteString("\n")
	}
	msg.WriteString("```")
	return msg.String()
}
