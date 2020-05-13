package email

type EmailConfig struct {
	SMTP     SMTP
	Tls      bool
	From     string
	FromName string
}

type SMTP struct {
	Host     string
	User     string
	Password string
}
