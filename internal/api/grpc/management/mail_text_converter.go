package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func mailTextRequestToModel(mailText *management.MailTextUpdate) *iam_model.MailText {
	return &iam_model.MailText{
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
	}
}

func mailTextFromModel(mailText *iam_model.MailText) *management.MailText {
	creationDate, err := ptypes.TimestampProto(mailText.CreationDate)
	logging.Log("MANAG-ULKZ6").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(mailText.ChangeDate)
	logging.Log("MANAG-451rI").OnError(err).Debug("date parse failed")

	return &management.MailText{
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}

func mailTextsViewFromModel(queries []*iam_model.MailTextView) *management.MailTextsView {
	modelQueries := make([]*management.MailTextView, len(queries))
	for i, query := range queries {
		modelQueries[i] = mailTextViewFromModel(query)
	}

	return &management.MailTextsView{
		Texts: modelQueries,
	}
}

func mailTextViewFromModel(mailText *iam_model.MailTextView) *management.MailTextView {
	creationDate, err := ptypes.TimestampProto(mailText.CreationDate)
	logging.Log("MANAG-koQnB").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(mailText.ChangeDate)
	logging.Log("MANAG-ToDhD").OnError(err).Debug("date parse failed")

	return &management.MailTextView{
		Default:      mailText.Default,
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
	}
}
