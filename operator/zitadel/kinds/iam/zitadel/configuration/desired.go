package configuration

import (
	"errors"
	"fmt"

	"github.com/caos/orbos/pkg/secret"
)

type Configuration struct {
	Tracing             *Tracing       `yaml:"tracing,omitempty"`
	Cache               *Cache         `yaml:"cache,omitempty"`
	Secrets             *Secrets       `yaml:"secrets,omitempty"`
	Notifications       *Notifications `yaml:"notifications,omitempty"`
	Passwords           *Passwords     `yaml:"passwords,omitempty"`
	DebugMode           bool           `yaml:"debugMode"`
	LogLevel            string         `yaml:"logLevel"`
	MigrateEventStoreV1 bool           `yaml:"migrateEventstoreV1"`
	DNS                 *DNS           `yaml:"dns"`
	ClusterDNS          string         `yaml:"clusterdns"`
	AssetStorage        *AssetStorage  `yaml:"assetStorage,omitempty"`
}

func (c *Configuration) Validate() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("validating configuration failed: %w", err)
		}
	}()

	return c.DNS.validate()
}

type AssetStorage struct {
	Type                    string           `yaml:"type,omitempty"`
	Endpoint                string           `yaml:"endpoint,omitempty"`
	AccessKeyID             *secret.Secret   `yaml:"accessKeyID,omitempty"`
	ExistingAccessKeyID     *secret.Existing `yaml:"existingAccessKeyID,omitempty"`
	SecretAccessKey         *secret.Secret   `yaml:"secretAccessKey,omitempty"`
	ExistingSecretAccessKey *secret.Existing `yaml:"ExistingSecretAccessKey,omitempty"`
	SSL                     bool             `yaml:"ssl,omitempty"`
	Location                string           `yaml:"location,omitempty"`
	BucketPrefix            string           `yaml:"bucketPrefix,omitempty"`
}

type DNS struct {
	Domain        string      `yaml:"domain"`
	TlsSecret     string      `yaml:"tlsSecret"`
	ACMEAuthority string      `yaml:"acmeAuthority"`
	Subdomains    *Subdomains `yaml:"subdomains"`
}

func (d *DNS) validate() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("validating dns failed: %w", err)
		}
	}()

	if d.TlsSecret != "" && d.ACMEAuthority != "none" && d.ACMEAuthority != "" {
		return errors.New("if tls secret is provided, acme authority must be 'none'")
	}
	return nil
}

type Subdomains struct {
	Accounts     string `yaml:"accounts"`
	API          string `yaml:"api"`
	Console      string `yaml:"console"`
	Issuer       string `yaml:"issuer"`
	Subscription string `yaml:"subscription"`
}
type Passwords struct {
	Migration            *secret.Secret   `yaml:"migration"`
	Management           *secret.Secret   `yaml:"management"`
	Auth                 *secret.Secret   `yaml:"auth"`
	Authz                *secret.Secret   `yaml:"authz"`
	Adminapi             *secret.Secret   `yaml:"adminapi"`
	Notification         *secret.Secret   `yaml:"notification"`
	Eventstore           *secret.Secret   `yaml:"eventstore"`
	Queries              *secret.Secret   `yaml:"queries"`
	ExistingMigration    *secret.Existing `yaml:"existingMigration"`
	ExistingManagement   *secret.Existing `yaml:"existingManagement"`
	ExistingAuth         *secret.Existing `yaml:"existingAuth"`
	ExistingAuthz        *secret.Existing `yaml:"existingAuthz"`
	ExistingAdminapi     *secret.Existing `yaml:"existingAdminapi"`
	ExistingNotification *secret.Existing `yaml:"existingNotification"`
	ExistingEventstore   *secret.Existing `yaml:"existingEventstore"`
	ExistingQueries      *secret.Existing `yaml:"existingQueries"`
}

type Secrets struct {
	Keys                    *secret.Secret   `yaml:"keys,omitempty"`
	ExistingKeys            *secret.Existing `yaml:"existingKeys,omitempty"`
	UserVerificationID      string           `yaml:"userVerificationID,omitempty"`
	OTPVerificationID       string           `yaml:"otpVerificationID,omitempty"`
	OIDCKeysID              string           `yaml:"oidcKeysID,omitempty"`
	CookieID                string           `yaml:"cookieID,omitempty"`
	CSRFID                  string           `yaml:"csrfID,omitempty"`
	DomainVerificationID    string           `yaml:"domainVerificationID,omitempty"`
	IDPConfigVerificationID string           `yaml:"idpConfigVerificationID,omitempty"`
}

type Notifications struct {
	GoogleChatURL         *secret.Secret   `yaml:"googleChatURL,omitempty"`
	ExistingGoogleChatURL *secret.Existing `yaml:"existingGoogleChatURL,omitempty"`
	Email                 *Email           `yaml:"email,omitempty"`
	Twilio                *Twilio          `yaml:"twilio,omitempty"`
}

type Tracing struct {
	ServiceAccountJSON         *secret.Secret   `yaml:"serviceAccountJSON,omitempty"`
	ExistingServiceAccountJSON *secret.Existing `yaml:"existingServiceAccountJSON,omitempty"`
	ProjectID                  string           `yaml:"projectID,omitempty"`
	Fraction                   string           `yaml:"fraction,omitempty"`
	Type                       string           `yaml:"type,omitempty"`
}

type Twilio struct {
	SenderName        string           `yaml:"senderName,omitempty"`
	AuthToken         *secret.Secret   `yaml:"authToken,omitempty"`
	SID               *secret.Secret   `yaml:"sid,omitempty"`
	ExistingAuthToken *secret.Existing `yaml:"existingAuthToken,omitempty"`
	ExistingSID       *secret.Existing `yaml:"ExistingSid,omitempty"`
}

type Email struct {
	SMTPHost       string           `yaml:"smtpHost,omitempty"`
	SMTPUser       string           `yaml:"smtpUser,omitempty"`
	SenderAddress  string           `yaml:"senderAddress,omitempty"`
	SenderName     string           `yaml:"senderName,omitempty"`
	TLS            bool             `yaml:"tls,omitempty"`
	AppKey         *secret.Secret   `yaml:"appKey,omitempty"`
	ExistingAppKey *secret.Existing `yaml:"existingAppKey,omitempty"`
}

type Cache struct {
	MaxAge            string `yaml:"maxAge,omitempty"`
	SharedMaxAge      string `yaml:"sharedMaxAge,omitempty"`
	ShortMaxAge       string `yaml:"shortMaxAge,omitempty"`
	ShortSharedMaxAge string `yaml:"shortSharedMaxAge,omitempty"`
}
