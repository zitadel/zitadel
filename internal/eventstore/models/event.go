package models

import (
	"time"

	"github.com/caos/zitadel/internal/errors"
)

type Event struct {
	//ID is set by eventstore
	ID               string
	CreationDate     time.Time
	Typ              string
	Sequence         uint64
	PreviousSequence uint64
	Data             []byte
	ModifierService  string
	ModifierTenant   string
	ModifierUser     string
	ResourceOwner    string
	AggregateType    string
	AggregateID      string
	AggregateVersion version
}

func (e *Event) Validate() error {
	if e.Typ == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-R2sB0", "type not defined")
	}
	if e.ModifierService == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-iGnu0", "modifier service not defined")
	}
	if e.ModifierUser == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-uZcBF", "modifier user not defined")
	}
	if e.ResourceOwner == "" {
		return errors.ThrowPreconditionFailed(nil, "MODEL-Bv0we", "resource owner not defined")
	}
	return nil
}
