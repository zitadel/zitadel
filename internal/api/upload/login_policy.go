package upload

import (
	"context"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
)

const (
	orgPrefix             = "org"
	policyPrefix          = orgPrefix + "/policy"
	labelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = labelPolicyPrefix + "/logo"
)

type labelPolicyLogo struct {
	idGenerator id.Generator
}

func (l *labelPolicyLogo) ObjectName() (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	return labelPolicyLogoPrefix + "-" + suffixID, nil
}

func (l *labelPolicyLogo) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	_, err := commands.AddLogoLabelPolicy(ctx, orgID, info.Key)
	return err
}
