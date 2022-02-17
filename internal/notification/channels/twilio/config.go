package twilio

type TwilioConfig struct {
	SID        string
	Token      string
	SenderName string
}

func (t *TwilioConfig) IsValid() bool {
	return t.SID != "" && t.Token != "" && t.SenderName != ""
}
