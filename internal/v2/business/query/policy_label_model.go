package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type LabelPolicyReadModel struct {
	eventstore.ReadModel

	PrimaryColor   string
	SecondaryColor string
}

func (rm *LabelPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.LabelPolicyAddedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		case *policy.LabelPolicyChangedEvent:
			rm.PrimaryColor = e.PrimaryColor
			rm.SecondaryColor = e.SecondaryColor
		}
	}
	return rm.ReadModel.Reduce()
}
