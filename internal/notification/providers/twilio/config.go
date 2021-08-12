package twilio

type TwilioConfig struct {
	SID   string
	Token string
	From  string
	Proxy *Proxy
}

type Proxy struct {
	HTTP     string
	HTTPS    string
	CertPath string
}
