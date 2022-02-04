package twilio

type TwilioConfig struct {
	SID   string
	Token string
	From  string
}

func (t *TwilioConfig) IsValid() bool {
	return t.SID != "" && t.Token != "" && t.From != ""
}
