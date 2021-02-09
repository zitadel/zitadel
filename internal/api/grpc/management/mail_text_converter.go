package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mailTextRequestToDomain(mailText *management.MailTextUpdate) *domain.MailText {
	return &domain.MailText{
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

func mailTextFromDoamin(mailText *domain.MailText) *management.MailText {
	return &management.MailText{
		MailTextType: mailText.MailTextType,
		Language:     mailText.Language,
		Title:        mailText.Title,
		PreHeader:    mailText.PreHeader,
		Subject:      mailText.Subject,
		Greeting:     mailText.Greeting,
		Text:         mailText.Text,
		ButtonText:   mailText.ButtonText,
		CreationDate: timestamppb.New(mailText.CreationDate),
		ChangeDate:   timestamppb.New(mailText.ChangeDate),
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
