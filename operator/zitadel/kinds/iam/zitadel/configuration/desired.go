package configuration

import (
	"errors"
	"fmt"

	"github.com/caos/orbos/pkg/secret"
)

type Configuration struct {
	//Tracing configuration for zitadel
	Tracing *Tracing `yaml:"tracing,omitempty"`
	//Cache configuration for zitadel
	Cache *Cache `yaml:"cache,omitempty"`
	//Secrets used by zitadel
	Secrets *Secrets `yaml:"secrets,omitempty"`
	//Notification configuration for zitadel
	Notifications *Notifications `yaml:"notifications,omitempty"`
	//Passwords used for the maintaining of the users in the database
	Passwords *Passwords `yaml:"passwords,omitempty"`
	//Debug mode for zitadel if notifications should be only sent by chat
	DebugMode bool `yaml:"debugMode"`
	//Log-level for zitadel
	LogLevel            string `yaml:"logLevel"`
	MigrateEventStoreV1 bool   `yaml:"migrateEventstoreV1"`
	//DNS configuration for subdomains
	DNS *DNS `yaml:"dns"`
	//ClusterDNS configuration for db user certificates
	ClusterDNS string `yaml:"clusterdns"`
	//Configuration for asset storage
	AssetStorage *AssetStorage `yaml:"assetStorage,omitempty"`
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
	MultiDelete             bool             `yaml:"multiDelete,omitempty"`
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
	//Password for the User "migration"
	Migration *secret.Secret `yaml:"migration"`
	//Password for the User "management"
	Management *secret.Secret `yaml:"management"`
	//Password for the User "auth"
	Auth *secret.Secret `yaml:"auth"`
	//Password for the User "authz"
	Authz *secret.Secret `yaml:"authz"`
	//Password for the User "adminapi"
	Adminapi *secret.Secret `yaml:"adminapi"`
	//Password for the User "notification"
	Notification *secret.Secret `yaml:"notification"`
	//Password for the User "eventstore"
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
	//Text-file which consists of a list of key/value to provide the keys to encrypt data in zitadel
	Keys *secret.Secret `yaml:"keys,omitempty"`
	//Key used from keys-file for user verification
	UserVerificationID string `yaml:"userVerificationID,omitempty"`
	//Key used from keys-file for OTP verification
	OTPVerificationID string `yaml:"otpVerificationID,omitempty"`
	//Key used from keys-file for OIDC
	OIDCKeysID string `yaml:"oidcKeysID,omitempty"`
	//Key used from keys-file for cookies
	CookieID string `yaml:"cookieID,omitempty"`
	//Key used from keys-file for CSRF
	CSRFID string `yaml:"csrfID,omitempty"`
	//Key used from keys-file for domain verification
	DomainVerificationID string `yaml:"domainVerificationID,omitempty"`
	//Key used from keys-file for IDP configuration verification
	IDPConfigVerificationID string `yaml:"idpConfigVerificationID,omitempty"`
}

type Notifications struct {
	ExistingGoogleChatURL *secret.Existing `yaml:"existingGoogleChatURL,omitempty"`
	//Google chat URL used for notifications
	GoogleChatURL *secret.Secret `yaml:"googleChatURL,omitempty"`
	//Configuration for email notifications
	Email *Email `yaml:"email,omitempty"`
	//Configuration for twilio notifications
	Twilio *Twilio `yaml:"twilio,omitempty"`
}

type Tracing struct {
	//User ServiceAccount to write tracing data to google cloud
	ServiceAccountJSON         *secret.Secret   `yaml:"serviceAccountJSON,omitempty"`
	ExistingServiceAccountJSON *secret.Existing `yaml:"existingServiceAccountJSON,omitempty"`
	//Used projectID to write data to google cloud
	ProjectID string `yaml:"projectID,omitempty"`
	//Fraction for tracing
	Fraction string `yaml:"fraction,omitempty"`
	//Tracing type
	Type string `yaml:"type,omitempty"`
}

type Twilio struct {
	//Sender name for Twilio
	SenderName string `yaml:"senderName,omitempty"`
	//Auth token to connect with Twilio
	AuthToken *secret.Secret `yaml:"authToken,omitempty"`
	//SID to connect with Twilio
	SID               *secret.Secret   `yaml:"sid,omitempty"`
	ExistingAuthToken *secret.Existing `yaml:"existingAuthToken,omitempty"`
	ExistingSID       *secret.Existing `yaml:"ExistingSid,omitempty"`
}

type Email struct {
	//SMTP host used for email notifications
	SMTPHost string `yaml:"smtpHost,omitempty"`
	//SMTP user used for email notifications
	SMTPUser string `yaml:"smtpUser,omitempty"`
	//Sender address from where the emails should get sent
	SenderAddress string `yaml:"senderAddress,omitempty"`
	//Sender name form where the emails should get sent
	SenderName string `yaml:"senderName,omitempty"`
	//Flag if TLS should be used for the communication with the SMTP host
	TLS bool `yaml:"tls,omitempty"`
	//Application-key used for SMTP communication
	AppKey         *secret.Secret   `yaml:"appKey,omitempty"`
	ExistingAppKey *secret.Existing `yaml:"existingAppKey,omitempty"`
}

type Cache struct {
	//Max age for cache records
	MaxAge string `yaml:"maxAge,omitempty"`
	//Max age for the shared cache records
	SharedMaxAge string `yaml:"sharedMaxAge,omitempty"`
	//Max age for the short cache records
	ShortMaxAge string `yaml:"shortMaxAge,omitempty"`
	//Max age for the short shared cache records
	ShortSharedMaxAge string `yaml:"shortSharedMaxAge,omitempty"`
}
