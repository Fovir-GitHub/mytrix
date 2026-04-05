package model

import (
	"text/template"

	"github.com/Fovir-GitHub/mytrix/internal/config"
)

var (
	gotifyTmpl     *template.Template
	wakapiLangTmpl *template.Template
	wakapiDataTmpl *template.Template
)

func InitTemplates() {
	cfg := config.Config

	gotifyTmpl = createTmpl("gotify", cfg.Gotify.Format)
	wakapiLangTmpl = createTmpl("wakapi_lang", cfg.Wakapi.LangFormat)
	wakapiDataTmpl = createTmpl("wakapi_data", cfg.Wakapi.DataFormat)
}

func createTmpl(name, format string) *template.Template {
	return template.Must(template.New(name).Parse(format))
}
