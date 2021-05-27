package text

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	text_pb "github.com/caos/zitadel/pkg/grpc/text"
)

func ModelCustomMsgTextToPb(msg *model.MailTextView) *text_pb.MessageCustomText {
	return &text_pb.MessageCustomText{
		Title:      msg.Title,
		PreHeader:  msg.PreHeader,
		Subject:    msg.Subject,
		Greeting:   msg.Greeting,
		Text:       msg.Text,
		ButtonText: msg.ButtonText,
		FooterText: msg.FooterText,
		Details: object.ToViewDetailsPb(
			msg.Sequence,
			msg.CreationDate,
			msg.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}
