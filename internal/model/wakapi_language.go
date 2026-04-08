package model

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
)

type WakapiLanguage struct {
	Name    string  `json:"name"`
	Text    string  `json:"text"`
	Percent float32 `json:"percent"`
}

type langView struct {
	Lang    string
	Text    string
	Percent string
}

func (wl WakapiLanguage) ToMarkdown() string {
	var buf bytes.Buffer
	percent := fmt.Sprintf("%.2f%%", wl.Percent)
	if err := wakapiLangTmpl.Execute(&buf, langView{
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

func generateLangReport(langs []WakapiLanguage) string {
	if len(langs) <= 0 {
		slog.Warn("no lanuage found, return empty report", "len", len(langs))
		return ""
	}

	var msg strings.Builder
	msg.WriteString("```text\n")
	for _, lang := range langs {
		msg.WriteString(lang.ToMarkdown())
		msg.WriteString("\n")
	}
	msg.WriteString("```")
	return msg.String()
}
