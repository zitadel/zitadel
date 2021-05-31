package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) MessageTextByIDs(aggregateID, textType, lang, mailTextTableVar string) (*model.MessageTextView, error) {
	return view.GetMessageTextByIDs(v.Db, mailTextTableVar, aggregateID, textType, lang)
}
