package repository

import (
	"regexp"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore/v3"
)

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

// Version represents the semver of an aggregate
type Version eventstore.Version

// Validate checks if the v is semver
func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}
