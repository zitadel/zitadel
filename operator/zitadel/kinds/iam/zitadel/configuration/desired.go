package configuration

import (
	"github.com/caos/orbos/pkg/secret"
)

type Configuration struct {
	Tracing       *Tracing       `yaml:"tracing,omitempty"`
	Cache         *Cache         `yaml:"cache,omitempty"`
	Secrets       *Secrets       `yaml:"secrets,omitempty"`
	Notifications *Notifications `yaml:"notifications,omitempty"`
	Passwords     *Passwords     `yaml:"passwords,omitempty"`
	DebugMode     bool           `yaml:"debugMode"`
	LogLevel      string         `yaml:"logLevel"`
	Ingress       *Ingress       `yaml:"ingress"`
	ClusterDNS    string         `yaml:"clusterdns"`
}

type Ingress struct {
	Domain              string                 `yaml:"domain"`
	TlsSecret           string                 `yaml:"tlsSecret"`
	Subdomains          *Subdomains            `yaml:"subdomains"`
	Controller          string                 `yaml:"controller"`
	ControllerSpecifics map[string]interface{} `yaml:"controllerSpecifics,omitempty"`
}

type Subdomains struct {
	Accounts string `yaml:"accounts"`
	API      string `yaml:"api"`
	Console  string `yaml:"console"`
	Issuer   string `yaml:"issuer"`
}
type Passwords struct {
	Migration    *secret.Secret `yaml:"migration"`
	Management   *secret.Secret `yaml:"management"`
	Auth         *secret.Secret `yaml:"auth"`
	Authz        *secret.Secret `yaml:"authz"`
	Adminapi     *secret.Secret `yaml:"adminapi"`
	Notification *secret.Secret `yaml:"notification"`
	Eventstore   *secret.Secret `yaml:"eventstore"`
}

type Secrets struct {
	Keys                    *secret.Secret `yaml:"keys,omitempty"`
	UserVerificationID      string         `yaml:"userVerificationID,omitempty"`
	OTPVerificationID       string         `yaml:"otpVerificationID,omitempty"`
	OIDCKeysID              string         `yaml:"oidcKeysID,omitempty"`
	CookieID                string         `yaml:"cookieID,omitempty"`
	CSRFID                  string         `yaml:"csrfID,omitempty"`
	DomainVerificationID    string         `yaml:"domainVerificationID,omitempty"`
	IDPConfigVerificationID string         `yaml:"idpConfigVerificationID,omitempty"`
}

type Notifications struct {
	GoogleChatURL *secret.Secret `yaml:"googleChatURL,omitempty"`
	Email         *Email         `yaml:"email,omitempty"`
	Twilio        *Twilio        `yaml:"twilio,omitempty"`
}

type Tracing struct {
	ServiceAccountJSON *secret.Secret `yaml:"serviceAccountJSON,omitempty"`
	ProjectID          string         `yaml:"projectID,omitempty"`
	Fraction           string         `yaml:"fraction,omitempty"`
	Type               string         `yaml:"type,omitempty"`
}

type Twilio struct {
	SenderName string         `yaml:"senderName,omitempty"`
	AuthToken  *secret.Secret `yaml:"authToken,omitempty"`
	SID        *secret.Secret `yaml:"sid,omitempty"`
}

type Email struct {
	SMTPHost      string         `yaml:"smtpHost,omitempty"`
	SMTPUser      string         `yaml:"smtpUser,omitempty"`
	SenderAddress string         `yaml:"senderAddress,omitempty"`
	SenderName    string         `yaml:"senderName,omitempty"`
	TLS           bool           `yaml:"tls,omitempty"`
	AppKey        *secret.Secret `yaml:"appKey,omitempty"`
}

type Cache struct {
	MaxAge            string `yaml:"maxAge,omitempty"`
	SharedMaxAge      string `yaml:"sharedMaxAge,omitempty"`
	ShortMaxAge       string `yaml:"shortMaxAge,omitempty"`
	ShortSharedMaxAge string `yaml:"shortSharedMaxAge,omitempty"`
}
