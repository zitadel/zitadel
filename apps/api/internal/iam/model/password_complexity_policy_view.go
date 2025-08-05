package model

import (
	"time"
)

type PasswordComplexityPolicyView struct {
	AggregateID  string
	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
	Default      bool

	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     uint64
}
