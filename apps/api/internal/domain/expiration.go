package domain

import (
	"time"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//most of us won't survive until 12-31-9999 23:59:59, maybe ZITADEL does
	defaultExpDate = time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)
)

type expiration interface {
	GetExpirationDate() time.Time
	SetExpirationDate(time.Time)
}

func EnsureValidExpirationDate(key expiration) error {
	date, err := ValidateExpirationDate(key.GetExpirationDate())
	if err != nil {
		return err
	}
	key.SetExpirationDate(date)
	return nil
}

func ValidateExpirationDate(date time.Time) (time.Time, error) {
	if date.IsZero() {
		return defaultExpDate, nil
	}
	if date.Before(time.Now()) {
		return time.Time{}, zerrors.ThrowInvalidArgument(nil, "DOMAIN-dv3t5", "Errors.AuthNKey.ExpireBeforeNow")
	}
	return date, nil
}
