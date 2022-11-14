package build

import "time"

var (
	version  = ""
	commit   = ""
	date     = ""
	dateTime time.Time
)

func Version() string {
	if version != "" {
		return version
	}
	version = Date().Format(time.RFC3339)
	return version
}

func Commit() string {
	return commit
}

func Date() time.Time {
	if !dateTime.IsZero() {
		return dateTime
	}
	dateTime, _ = time.Parse(time.RFC3339, date)
	if dateTime.IsZero() {
		dateTime = time.Now()
	}
	return dateTime
}
