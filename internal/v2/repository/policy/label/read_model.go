package label

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	PrimaryColor   string
	SecondaryColor string
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		case *ChangedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		}
	}
	return rm.ReadModel.Reduce()
}
