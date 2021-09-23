package query

import "time"

type PasswordComplexityPolicyProjection struct {
	ID            string
	Sequence      uint64
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string

	MinLength    uint64
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}
