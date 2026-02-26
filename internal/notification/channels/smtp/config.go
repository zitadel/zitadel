package smtp

type Config struct {
	SMTP           SMTP
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

type SMTP struct {
	Host        string
	PlainAuth   *PlainAuthConfig
	XOAuth2Auth *XOAuth2AuthConfig
}

type ConfigHTTP struct {
	Endpoint string
}
