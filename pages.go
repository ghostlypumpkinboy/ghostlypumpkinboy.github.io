package main

import (
	"github.com/BurntSushi/toml"
	"time"
)

type SiteHistoryType string

const (
	CreatedSlug SiteHistoryType = "created_slug"
	UpdatedSlug = "updated_slug"
	DestroyedSlug = "destroyed_slug"

	CreatedPage = "created_page"
	UpdatedPage = "updated_page"
	DestroyedPage = "destroyed_page"

	CreatedSite = "created_site"
	UpdatedSite = "updated_site"
	DestroyedSite = "destroyed_site"
)

type SiteHistoryEntry struct {
	Type SiteHistoryType
	Site string
	Page string
	Slug string
	Date time.Time
	Comment string
	Content map[string]any
}

type SiteHistory struct {
	Entries []SiteHistoryEntry
}

