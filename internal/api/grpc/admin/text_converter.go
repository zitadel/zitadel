package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func textToModel(text *admin.DefaultMailTextUpdate) *iam_model.MailText {
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

func textFromModel(text *iam_model.MailText) *admin.DefaultMailText {
	creationDate, err := ptypes.TimestampProto(text.CreationDate)
	logging.Log("ADMIN-Jlzsj").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(text.ChangeDate)
	logging.Log("ADMIN-mw5b8").OnError(err).Debug("date parse failed")

	return &admin.DefaultMailText{
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

func textsViewFromModel(textsin *iam_model.MailTextsView) *admin.DefaultMailTextsView {
	return &admin.DefaultMailTextsView{
		Texts: textsViewToModel(textsin.Texts),
	}
}

func textsViewToModel(queries []*iam_model.MailTextView) []*admin.DefaultMailTextView {
	modelQueries := make([]*admin.DefaultMailTextView, len(queries))
	for i, query := range queries {
		modelQueries[i] = textViewFromModel(query)
	}

	return modelQueries
}

func textViewFromModel(text *iam_model.MailTextView) *admin.DefaultMailTextView {
	creationDate, err := ptypes.TimestampProto(text.CreationDate)
	logging.Log("ADMIN-7RyJc").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(text.ChangeDate)
	logging.Log("ADMIN-fTFgY").OnError(err).Debug("date parse failed")

	return &admin.DefaultMailTextView{
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
