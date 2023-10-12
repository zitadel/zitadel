package smtp

type Config struct {
	ConfigID       string
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
	IsActive     bool
	ProviderType uint32
}

func (smtp *SMTP) HasAuth() bool {
	return smtp.User != "" && smtp.Password != ""
}
