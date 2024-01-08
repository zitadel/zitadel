package eventstore

import (
	"regexp"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Version string

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

func (v Version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return zerrors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}
