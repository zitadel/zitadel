package password_complexity

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	MinLength    uint64
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.MinLength = e.MinLength
			rm.HasLowercase = e.HasLowercase
			rm.HasUpperCase = e.HasUpperCase
			rm.HasNumber = e.HasNumber
			rm.HasSymbol = e.HasSymbol
		case *ChangedEvent:
			rm.MinLength = e.MinLength
			rm.HasLowercase = e.HasLowercase
			rm.HasUpperCase = e.HasUpperCase
			rm.HasNumber = e.HasNumber
			rm.HasSymbol = e.HasSymbol
		}
	}
	return rm.ReadModel.Reduce()
}
