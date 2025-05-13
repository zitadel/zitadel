package info

import (
	"time"

	"github.com/Masterminds/semver/v3"
)

var (
	version string
	commit  string
	date    string
)

func Version() *semver.Version {
	v, _ := semver.NewVersion(version)
	if v != nil {
		return v
	}
	return semver.New(uint64(Date().Year()), uint64(Date().Month()), uint64(Date().Day()), "", "")
}

func Commit() string {
	return commit
}

func Date() time.Time {
	if date == "" {
		return time.Now()
	}
	d, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Now()
	}
	return d
}
