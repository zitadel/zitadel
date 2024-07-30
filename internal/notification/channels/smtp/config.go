package smtp

type Config struct {
	Description    string
	SMTP           SMTP
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

type SMTP struct {
	Host     string
	User     string
	Password string
}

func (smtp *SMTP) HasAuth() bool {
	return smtp.User != "" && smtp.Password != ""
}
