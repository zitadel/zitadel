package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type CustomTextWriteModel struct {
	eventstore.WriteModel

	Key      string
	Language language.Tag
	Text     string
	State    domain.CustomTextState
}

func (wm *CustomTextWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.CustomTextSetEvent:
			if wm.Key != e.Key || wm.Language != e.Language {
				continue
			}
			wm.Text = e.Text
			wm.State = domain.CustomTextStateActive
		case *policy.MailTextRemovedEvent:
			wm.State = domain.CustomTextStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
