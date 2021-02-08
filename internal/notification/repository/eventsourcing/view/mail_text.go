package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) MailTextByIDs(aggregateID string, textType string, language string, mailTextTableVar string) (*model.MailTextView, error) {
	return view.GetMailTextByIDs(v.Db, mailTextTableVar, aggregateID, textType, language)
}
