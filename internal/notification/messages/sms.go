package messages

import "github.com/zitadel/zitadel/internal/notification/channels"

var _ channels.Message = (*SMS)(nil)

type SMS struct {
	SenderPhoneNumber    string
	RecipientPhoneNumber string
	Content              string
}

func (msg *SMS) GetContent() string {
	return msg.Content
}
