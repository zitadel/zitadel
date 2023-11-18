package smtp

type Config struct {
	SMTP           SMTP
	Tls            bool
	From           string
	FromName       string
	ReplyToAddress string
}

type SMTP struct {
	Host         string
	User         string
	Password     string
	ProviderType uint32
}

func (smtp *SMTP) HasAuth() bool {
	return smtp.User != "" && smtp.Password != ""
}
