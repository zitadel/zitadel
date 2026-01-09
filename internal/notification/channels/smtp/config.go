package smtp

type Config struct {
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

// HasAuth reports whether SMTP authentication should be attempted.
//
// It returns true when a User is configured, regardless of whether Password
// is set. This behavior is intentional to support SMTP servers that allow
// authentication without a password. For servers that require a password,
// this means authentication may still be attempted with an empty password.
func (smtp *SMTP) HasAuth() bool {
	return smtp.User != ""
}

type ConfigHTTP struct {
	Endpoint string
}
