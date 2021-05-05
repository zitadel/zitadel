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
	iamPrefix             = "iam"
	policyPrefix          = "/policy"
	labelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = labelPolicyPrefix + "/logo"
	labelPolicyIconPrefix = labelPolicyPrefix + "/icon"
	labelPolicyFontPrefix = labelPolicyPrefix + "/font"

	defaultLabelPolicyLogoPrefix = iamPrefix + labelPolicyLogoPrefix
	defaultLabelPolicyIconPrefix = iamPrefix + labelPolicyIconPrefix
	defaultLabelPolicyFontPrefix = iamPrefix + labelPolicyFontPrefix

	orgLabelPolicyLogoPrefix = orgPrefix + labelPolicyLogoPrefix
	orgLabelPolicyIconPrefix = orgPrefix + labelPolicyIconPrefix
	orgLabelPolicyFontPrefix = orgPrefix + labelPolicyFontPrefix
)

type labelPolicyLogo struct {
	idGenerator   id.Generator
	darkMode      bool
	defaultPolicy bool
}

func (l *labelPolicyLogo) ObjectName(_ authz.CtxData) (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	prefix := orgLabelPolicyLogoPrefix
	if l.defaultPolicy {
		prefix = defaultLabelPolicyLogoPrefix
	}
	if l.darkMode {
		return prefix + "-" + dark + "-" + suffixID, nil
	}
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyLogo) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.defaultPolicy {
		if l.darkMode {
			_, err := commands.AddLogoDarkDefaultLabelPolicy(ctx, info.Key)
			return err
		}
		_, err := commands.AddLogoDefaultLabelPolicy(ctx, info.Key)
		return err
	}
	if l.darkMode {
		_, err := commands.AddLogoDarkLabelPolicy(ctx, orgID, info.Key)
		return err
	}
	_, err := commands.AddLogoLabelPolicy(ctx, orgID, info.Key)
	return err
}

type labelPolicyIcon struct {
	idGenerator   id.Generator
	darkMode      bool
	defaultPolicy bool
}

func (l *labelPolicyIcon) ObjectName(_ authz.CtxData) (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	prefix := orgLabelPolicyIconPrefix
	if l.defaultPolicy {
		prefix = defaultLabelPolicyIconPrefix
	}
	if l.darkMode {
		return prefix + "-" + dark + "-" + suffixID, nil
	}
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyIcon) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.defaultPolicy {
		if l.darkMode {
			_, err := commands.AddIconDarkDefaultLabelPolicy(ctx, info.Key)
			return err
		}
		_, err := commands.AddIconDefaultLabelPolicy(ctx, info.Key)
		return err
	}

	if l.darkMode {
		_, err := commands.AddIconDarkLabelPolicy(ctx, orgID, info.Key)
		return err
	}
	_, err := commands.AddIconLabelPolicy(ctx, orgID, info.Key)
	return err
}

type labelPolicyFont struct {
	idGenerator   id.Generator
	defaultPolicy bool
}

func (l *labelPolicyFont) ObjectName(_ authz.CtxData) (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	prefix := orgLabelPolicyFontPrefix
	if l.defaultPolicy {
		prefix = defaultLabelPolicyFontPrefix
	}
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyFont) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.defaultPolicy {
		_, err := commands.AddFontDefaultLabelPolicy(ctx, info.Key)
		return err
	}
	_, err := commands.AddFontLabelPolicy(ctx, orgID, info.Key)
	return err
}
