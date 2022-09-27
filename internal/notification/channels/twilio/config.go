package twilio

type TwilioConfig struct {
	SID          string
	Token        string
	SenderNumber string
}

func (t *TwilioConfig) IsValid() bool {
	return t.SID != "" && t.Token != "" && t.SenderNumber != ""
}
