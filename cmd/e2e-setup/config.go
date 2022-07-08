package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/zitadel/zitadel/internal/crypto"

	"github.com/zitadel/zitadel/internal/id"

	"github.com/zitadel/zitadel/internal/database"

	internal_authz "github.com/zitadel/zitadel/internal/api/authz"

	static_config "github.com/zitadel/zitadel/internal/static/config"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/config/hook"
)

type Config struct {
	E2E            *E2EConfig
	Log            *logging.Config
	ExternalPort   uint16
	ExternalDomain string
	ExternalSecure bool
	Database       database.Config
	AssetStorage   static_config.AssetStorageConfig
	WebAuthNName   string
	InternalAuthZ  internal_authz.Config
	Machine        *id.Config
	SystemDefaults systemdefaults.SystemDefaults
	EncryptionKeys *encryptionKeyConfig
}

func (c Config) Validate() error {
	if c.E2E == nil {
		return errors.New("no e2e config found")
	}
	return c.E2E.Validate()
}

type E2EConfig struct {
	Org                            string
	MachineKeyPath                 string
	InstanceID                     string
	ZitadelProjectResourceID       string
	APIURL                         string
	IssuerURL                      string
	Audience                       string
	OrgOwnerPassword               string
	OrgOwnerViewerPassword         string
	OrgProjectCreatorPassword      string
	PasswordComplexityUserPassword string
	LoginPolicyUserPassword        string
}

func (e E2EConfig) Validate() (err error) {
	if e.Org == "" {
		return errors.New("field Org is empty")
	}
	if e.MachineKeyPath == "" {
		return errors.New("field MachineKeyPath is empty")
	}
	if e.ZitadelProjectResourceID == "" {
		return errors.New("field ZitadelProjectResourceID is empty")
	}

	audPattern := "number-[0-9]{17}"
	matched, err := regexp.MatchString("bignumber-[0-9]{17}", e.ZitadelProjectResourceID)
	if err != nil {
		return fmt.Errorf("validating ZitadelProjectResourceID failed: %w", err)
	}
	if !matched {
		return fmt.Errorf("ZitadelProjectResourceID doesn't match regular expression %s", audPattern)
	}

	if e.APIURL == "" {
		return errors.New("field APIURL is empty")
	}
	if e.IssuerURL == "" {
		return errors.New("field IssuerURL is empty")
	}
	if e.OrgOwnerPassword == "" {
		return errors.New("field OrgOwnerPassword is empty")
	}
	if e.OrgOwnerViewerPassword == "" {
		return errors.New("field OrgOwnerViewerPassword is empty")
	}
	if e.OrgProjectCreatorPassword == "" {
		return errors.New("field OrgProjectCreatorPassword is empty")
	}
	if e.PasswordComplexityUserPassword == "" {
		return errors.New("field PasswordComplexityUserPassword is empty")
	}
	if e.LoginPolicyUserPassword == "" {
		return errors.New("field LoginPolicyUserPassword is empty")
	}
	return nil
}

type encryptionKeyConfig struct {
	DomainVerification   *crypto.KeyConfig
	IDPConfig            *crypto.KeyConfig
	OIDC                 *crypto.KeyConfig
	OTP                  *crypto.KeyConfig
	SMS                  *crypto.KeyConfig
	SMTP                 *crypto.KeyConfig
	User                 *crypto.KeyConfig
	CSRFCookieKeyID      string
	UserAgentCookieKeyID string
}

func MustNewConfig(v *viper.Viper) *Config {
	config := new(Config)

	err := v.Unmarshal(config,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
			hook.Base64ToBytesHookFunc(),
			hook.TagToLanguageHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)),
	)
	err = config.Log.SetLogger()
	logging.OnError(err).Fatal("unable to set logger")

	return config
}
