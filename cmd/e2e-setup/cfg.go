package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/caos/logging"
	internal_authz "github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
)

var (
	e2eSetupPaths = config.NewArrayFlags("authz.yaml", "system-defaults.yaml", "setup.yaml", "e2e.yaml")
)

type setupConfig struct {
	E2E E2EConfig

	Log logging.Config

	Eventstore     types.SQL
	SystemDefaults sd.SystemDefaults
	InternalAuthZ  internal_authz.Config
}

type E2EConfig struct {
	Org                            string
	MachineKeyPath                 string
	ZitadelProjectResourceID       string
	APIURL                         string
	IssuerURL                      string
	OrgOwnerPassword               string
	OrgOwnerViewerPassword         string
	OrgProjectCreatorPassword      string
	PasswordComplexityUserPassword string
	LoginPolicyUserPassword        string
}

func (e E2EConfig) validate() (err error) {
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
	matched, err := regexp.MatchString("number-[0-9]{17}", e.ZitadelProjectResourceID)
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
