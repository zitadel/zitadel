package org_iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
)

type ReadModel struct{ org_iam.ReadModel }

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *org_iam.AddedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
