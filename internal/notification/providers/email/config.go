package email

type EmailConfig struct {
	SMTP SMTP
	Tls  bool
}

type SMTP struct {
	Host     string
	User     string
	Password string
}
