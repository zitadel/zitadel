package provider

import "github.com/caos/zitadel/internal/eventstore/v2"

type WriteModel struct {
	eventstore.WriteModel

	IDPConfigID     string
	IDPProviderType Type
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.IDPProviderType = e.IDPProviderType
		}
	}
	return wm.WriteModel.Reduce()
}
