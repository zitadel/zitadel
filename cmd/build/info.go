package build

import (
	"log/slog"
	"time"
)

// These variables are set via ldflags in the Makefile
var (
	version = ""
	commit  = ""
	date    = ""
)

// dateTime is the parsed version of [date]
var dateTime time.Time

// init prevents race conditions when accessing dateTime and version.
func init() {
	var err error
	dateTime, err = time.Parse(time.RFC3339, date)
	if err != nil {
		slog.Warn("could not parse build date, using current time instead", "err", err)
		dateTime = time.Now()
		date = dateTime.Format(time.RFC3339)
	}
	if version == "" {
		slog.Warn("no build version set, using timestamp as version")
		version = date
	}
}

// Version returns the current build version of Zitadel
func Version() string {
	return version
}

// Commit returns the git commit hash of the current build of Zitadel
func Commit() string {
	return commit
}

// Date returns the build date of the current build of Zitadel
func Date() time.Time {
	return dateTime
}
