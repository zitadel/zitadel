package models

import (
	"regexp"

	"github.com/caos/zitadel/internal/errors"
)

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

type Version string

func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}

func (v Version) String() string {
	return string(v)
}
