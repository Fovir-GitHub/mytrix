package render

import (
	"bytes"
	"fmt"
	"log/slog"

	"codeberg.org/Fovir/mytrix/internal/model"
)

type umamiWebsiteView struct {
	Name        string
	Domain      string
	Visitors    int
	Visits      int
	BouncesRate string
}

func umamiWebsiteToView(u *model.UmamiWebsite) *umamiWebsiteView {
	var bouncesRate string
	if u.Stat.Visits == 0 {
		bouncesRate = "0%"
	} else {
		bouncesRate = fmt.Sprintf("%.2f%%", float64(u.Stat.Bounces)/float64(u.Stat.Visits)*100)
	}

	return &umamiWebsiteView{
		Name:        u.Name,
		Domain:      u.Domain,
		Visitors:    u.Stat.Visitors,
		Visits:      u.Stat.Visits,
		BouncesRate: bouncesRate,
	}
}

func UmamiWebsiteMarkdown(u *model.UmamiWebsite) string {
	var buf bytes.Buffer
	view := umamiWebsiteToView(u)
	if err := umamiDataTmpl.Execute(&buf, view); err != nil {
		slog.Error(
			"parse umami message to markdown failed",
			"name", view.Name,
			"domain", view.Domain,
			"visitors", view.Visitors,
			"visits", view.Visits,
			"bounces_rate", view.BouncesRate,
			"err", err,
		)
		return fmt.Sprintf("Name: %s\nDomain: %s\nVisitors: %d\nVisits: %d\nBounces Rate: %s", view.Name, view.Domain, view.Visitors, view.Visits, view.BouncesRate)
	}
	return buf.String()
}
