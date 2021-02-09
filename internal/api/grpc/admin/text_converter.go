package admin

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func textToDomain(text *admin.DefaultMailTextUpdate) *domain.MailText {
	return &domain.MailText{
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

func textFromDomain(text *domain.MailText) *admin.DefaultMailText {
	return &admin.DefaultMailText{
		MailTextType: text.MailTextType,
		Language:     text.Language,
		Title:        text.Title,
		PreHeader:    text.PreHeader,
		Subject:      text.Subject,
		Greeting:     text.Greeting,
		Text:         text.Text,
		ButtonText:   text.ButtonText,
		CreationDate: timestamppb.New(text.CreationDate),
		ChangeDate:   timestamppb.New(text.ChangeDate),
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
	return &admin.DefaultMailTextView{
		MailTextType: text.MailTextType,
		Language:     text.Language,
		Title:        text.Title,
		PreHeader:    text.PreHeader,
		Subject:      text.Subject,
		Greeting:     text.Greeting,
		Text:         text.Text,
		ButtonText:   text.ButtonText,
		CreationDate: timestamppb.New(text.CreationDate),
		ChangeDate:   timestamppb.New(text.ChangeDate),
	}
}
