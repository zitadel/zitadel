package twilio

type Config struct {
	SID              string
	Token            string
	SenderNumber     string
	VerifyServiceSID string
}

func (t *Config) IsValid() bool {
	return t.SID != "" && t.Token != "" && t.SenderNumber != ""
}
