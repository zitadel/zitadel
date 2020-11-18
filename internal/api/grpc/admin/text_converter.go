package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func textToModel(text *admin.DefaultTextUpdate) *iam_model.MailText {
	return &iam_model.MailText{
		MailTextType: text.MailTextType,
		Language:     text.Language,
		Title:        text.Title,
		PreHeader:    text.PreHeader,
		Subject:      text.Subject,
		Greeting:     text.Greeting,
		Text:         text.Text,
		ButtonText:   text.ButtonText,
	}
}

func textFromModel(text *iam_model.MailText) *admin.DefaultText {
	creationDate, err := ptypes.TimestampProto(text.CreationDate)
	logging.Log("ADMIN-Jlzsj").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(text.ChangeDate)
	logging.Log("ADMIN-mw5b8").OnError(err).Debug("date parse failed")

	return &admin.DefaultText{
		MailTextType: text.MailTextType,
		Language:     text.Language,
		Title:        text.Title,
		PreHeader:    text.PreHeader,
		Subject:      text.Subject,
		Greeting:     text.Greeting,
		Text:         text.Text,
		ButtonText:   text.ButtonText,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

// func textsViewFromModel(text *iam_model.MailTextsView) *admin.DefaultTextsView {

// }

func textsViewFromModel(textsin *iam_model.MailTextsView) *admin.DefaultTextsView {
	return &admin.DefaultTextsView{
		Texts: textsViewToModel(textsin.Texts),
	}
}

func textsViewToModel(queries []*iam_model.MailTextView) []*admin.DefaultTextView {
	modelQueries := make([]*admin.DefaultTextView, len(queries))
	for i, query := range queries {
		modelQueries[i] = textViewFromModel(query)
	}

	return modelQueries
}

func textViewFromModel(text *iam_model.MailTextView) *admin.DefaultTextView {
	creationDate, err := ptypes.TimestampProto(text.CreationDate)
	logging.Log("ADMIN-7RyJc").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(text.ChangeDate)
	logging.Log("ADMIN-fTFgY").OnError(err).Debug("date parse failed")

	return &admin.DefaultTextView{
		MailTextType: text.MailTextType,
		Language:     text.Language,
		Title:        text.Title,
		PreHeader:    text.PreHeader,
		Subject:      text.Subject,
		Greeting:     text.Greeting,
		Text:         text.Text,
		ButtonText:   text.ButtonText,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}
