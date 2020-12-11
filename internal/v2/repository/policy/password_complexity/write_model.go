package password_complexity

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type WriteModel struct {
	eventstore.WriteModel

	MinLength    uint64
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.MinLength = e.MinLength
			wm.HasLowercase = e.HasLowercase
			wm.HasUpperCase = e.HasUpperCase
			wm.HasNumber = e.HasNumber
			wm.HasSymbol = e.HasSymbol
		case *ChangedEvent:
			wm.MinLength = e.MinLength
			wm.HasLowercase = e.HasLowercase
			wm.HasUpperCase = e.HasUpperCase
			wm.HasNumber = e.HasNumber
			wm.HasSymbol = e.HasSymbol
		}
	}
	return wm.WriteModel.Reduce()
}
