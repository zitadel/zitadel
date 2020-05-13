package twilio

type TwilioMessage struct {
	SenderPhoneNumber    string
	RecipientPhoneNumber string
	Content              string
}

func (msg TwilioMessage) GetContent() string {
	return msg.Content
}
