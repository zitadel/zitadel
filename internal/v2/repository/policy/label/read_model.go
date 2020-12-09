package label

import "github.com/caos/zitadel/internal/eventstore/v2"

type LabelPolicyReadModel struct {
	eventstore.ReadModel

	PrimaryColor   string
	SecondaryColor string
}

func (rm *LabelPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		case *LabelPolicyChangedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		}
	}
	return rm.ReadModel.Reduce()
}
