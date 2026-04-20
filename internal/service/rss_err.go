package service

import "errors"

var (
	ErrRSSFetchFeeds    = errors.New("fetch feeds failed")
	ErrRSSPartialUpdate = errors.New("partial update failed")
	ErrRSSNoUpdate      = errors.New("no feed updated")
)
