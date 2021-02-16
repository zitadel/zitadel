package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) MailTemplateByAggregateID(aggregateID string, mailTemplateTableVar string) (*model.MailTemplateView, error) {
	return view.GetMailTemplateByAggregateID(v.Db, mailTemplateTableVar, aggregateID)
}
