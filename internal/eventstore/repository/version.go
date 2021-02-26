package repository

import (
	"regexp"

	"github.com/caos/zitadel/internal/errors"
)

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

//Version represents the semver of an aggregate
type Version string

//Validate checks if the v is semver
func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}
