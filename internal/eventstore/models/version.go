package models

import (
	"fmt"
	"regexp"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

var versionRegexp = regexp.MustCompile(`^v[0-9]+(\.[0-9]+){0,2}$`)

type version string

func Version(major, minor, patch int) (version, error) {
	if major < 0 || minor < 0 || patch < 0 {
		return "", errors.ThrowInvalidArgument(nil, "MODEL-ToKK8", "versions must be >= 0")
	}
	return version(fmt.Sprintf("v%d.%d.%d", major, minor, patch)), nil
}

func MustVersion(major, minor, patch int) version {
	version, err := Version(major, minor, patch)
	logging.Log("MODEL-4hlgF").OnError(err).Fatal("invalid version number")

	return version
}

func (v version) Validate() error {
	if !versionRegexp.MatchString(string(v)) {
		return errors.ThrowPreconditionFailed(nil, "MODEL-luDuS", "version is not semver")
	}
	return nil
}

func (v version) String() string {
	return string(v)
}
