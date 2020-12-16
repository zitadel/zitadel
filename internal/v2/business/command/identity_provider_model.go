package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IdentityProviderWriteModel struct {
	eventstore.WriteModel

	IDPConfigID     string
	IDPProviderType domain.IdentityProviderType
	IsActive        bool
}

func (wm *IdentityProviderWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.IdentityProviderAddedEvent:
			wm.IDPConfigID = e.IDPConfigID
			wm.IDPProviderType = e.IDPProviderType
			wm.IsActive = true
		case *policy.IdentityProviderRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
