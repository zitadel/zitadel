package upload

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
)

const (
	dark                  = "dark"
	orgPrefix             = "org"
	policyPrefix          = orgPrefix + "/policy"
	labelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = labelPolicyPrefix + "/logo"
)

type labelPolicyLogo struct {
	idGenerator id.Generator
	darkMode    bool
}

func (l *labelPolicyLogo) ObjectName(_ authz.CtxData) (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	if l.darkMode {
		return labelPolicyLogoPrefix + "-" + dark + "-" + suffixID, nil
	}
	return labelPolicyLogoPrefix + "-" + suffixID, nil
}

func (l *labelPolicyLogo) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.darkMode {
		_, err := commands.AddLogoDarkLabelPolicy(ctx, orgID, info.Key)
		return err
	}
	_, err := commands.AddLogoLabelPolicy(ctx, orgID, info.Key)
	return err
}
