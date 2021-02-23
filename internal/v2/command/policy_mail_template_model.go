package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type MailTemplateWriteModel struct {
	eventstore.WriteModel

	Template []byte

	State domain.PolicyState
}

func (wm *MailTemplateWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.MailTemplateAddedEvent:
			wm.Template = e.Template
			wm.State = domain.PolicyStateActive
		case *policy.MailTemplateChangedEvent:
			if e.Template != nil {
				wm.Template = *e.Template
			}
		case *policy.MailTemplateRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
