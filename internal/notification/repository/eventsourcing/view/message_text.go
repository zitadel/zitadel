package view

import (
	"github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/iam/repository/view/model"
)

func (v *View) MessageTextByIDs(aggregateID, textType, lang, messageTextTableVar string) (*model.MessageTextView, error) {
	return view.GetMessageTextByIDs(v.Db, messageTextTableVar, aggregateID, textType, lang)
}
