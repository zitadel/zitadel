package domain

import "time"

type AssetInfo struct {
	Bucket          string
	Key             string
	ETag            string
	Size            int64
	LastModified    time.Time
	Location        string
	VersionID       string
	Expiration      time.Time
	AutheticatedURL string
}
