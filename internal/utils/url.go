package utils

import "net/url"

func ParseURLHost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	return u.Host
}
