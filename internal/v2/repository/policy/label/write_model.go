package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type WriteModel struct {
	eventstore.WriteModel

	PrimaryColor   string
	SecondaryColor string
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
		case *ChangedEvent:
			wm.PrimaryColor = e.PrimaryColor
			wm.SecondaryColor = e.SecondaryColor
		}
	}
	return wm.WriteModel.Reduce()
}
