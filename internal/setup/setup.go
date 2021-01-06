package setup

import (
	"context"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/command"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	OrgOwnerRole                   = "ORG_OWNER"
	SetupUser                      = "SETUP"
	OIDCResponseTypeCode           = "CODE"
	OIDCResponseTypeIDToken        = "ID_TOKEN"
	OIDCResponseTypeToken          = "ID_TOKEN TOKEN"
	OIDCGrantTypeAuthorizationCode = "AUTHORIZATION_CODE"
	OIDCGrantTypeImplicit          = "IMPLICIT"
	OIDCGrantTypeRefreshToken      = "REFRESH_TOKEN"
	OIDCApplicationTypeNative      = "NATIVE"
	OIDCApplicationTypeUserAgent   = "USER_AGENT"
	OIDCApplicationTypeWeb         = "WEB"
	OIDCAuthMethodTypeNone         = "NONE"
	OIDCAuthMethodTypeBasic        = "BASIC"
	OIDCAuthMethodTypePost         = "POST"
)

func Execute(ctx context.Context, setUpConfig IAMSetUp, iamID string, commands *command.CommandSide) error {
	logging.Log("SETUP-JAK2q").Info("starting setup")

	iam, err := commands.GetIAM(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == domain.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("SETUP-VA2k1").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &iam_model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: iamID}}
	}

	steps, err := setUpConfig.Steps(iam.SetUpDone)
	if err != nil || len(steps) == 0 {
		return err
	}

	err = commands.ExecuteSetupSteps(ctx, steps)
	if err != nil {
		return err
	}

	logging.Log("SETUP-ds31h").Info("setup done")
	return nil
}
