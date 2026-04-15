package model

import (
	"log/slog"
	"text/template"

	"codeberg.org/Fovir/mytrix/internal/config"
)

var (
	gotifyTmpl     *template.Template
	wakapiLangTmpl *template.Template
	wakapiDataTmpl *template.Template
	umamiDataTmpl  *template.Template
	rssItemTmpl    *template.Template
)

func InitTemplates() {
	cfg := config.Config
	gotifyTmpl = createTmpl("gotify", cfg.Gotify.Format)
	wakapiLangTmpl = createTmpl("wakapi_lang", cfg.Wakapi.LangFormat)
	wakapiDataTmpl = createTmpl("wakapi_data", cfg.Wakapi.DataFormat)
	umamiDataTmpl = createTmpl("umami_stat", cfg.Umami.Format)
	rssItemTmpl = createTmpl("rss_item", cfg.RSS.Format)

	slog.Info("templates initiailized")
}

func createTmpl(name, format string) *template.Template {
	return template.Must(template.New(name).Parse(format))
}
