package model

import (
	"bytes"
	"fmt"
	"log/slog"
)

type UmamiWebsite struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Stat   *UmamiWebsiteStat
}

type UmamiWebsiteStat struct {
	Visitors int `json:"visitors"`
	Visits   int `json:"visits"`
	Bounces  int `json:"bounces"`
}

type umamiWebsiteView struct {
	Name        string
	Domain      string
	Visitors    int
	Visits      int
	BouncesRate string
}

func (u *UmamiWebsite) toView() *umamiWebsiteView {
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

func (u *UmamiWebsite) ToMarkdown() string {
	var buf bytes.Buffer
	view := u.toView()
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
