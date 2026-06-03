package model

// UmamiWebsite represents an analytics website tracked by Umami.
// It contains the website ID, name, domain, and associated statistics.
type UmamiWebsite struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Stat   *UmamiWebsiteStat
}

// UmamiWebsiteStat represents visitor statistics for an Umami website.
// It contains the number of visitors, total visits, and bounce count.
type UmamiWebsiteStat struct {
	Visitors int `json:"visitors"`
	Visits   int `json:"visits"`
	Bounces  int `json:"bounces"`
}
