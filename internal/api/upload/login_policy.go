package upload

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
)

const (
	defaultLabelPolicyLogoURL     = "/" + domain.DefaultLabelPolicyLogoPath
	defaultLabelPolicyLogoDarkURL = "/" + domain.DefaultLabelPolicyLogoPath + domain.Dark
	defaultLabelPolicyIconURL     = "/" + domain.DefaultLabelPolicyIconPath
	defaultLabelPolicyIconDarkURL = "/" + domain.DefaultLabelPolicyIconPath + domain.Dark
	defaultLabelPolicyFontURL     = "/" + domain.DefaultLabelPolicyFontPath

	orgLabelPolicyLogoURL     = "/" + domain.OrgLabelPolicyLogoPath
	orgLabelPolicyLogoDarkURL = "/" + domain.OrgLabelPolicyLogoPath + domain.Dark
	orgLabelPolicyIconURL     = "/" + domain.OrgLabelPolicyIconPath
	orgLabelPolicyIconDarkURL = "/" + domain.OrgLabelPolicyIconPath + domain.Dark
	orgLabelPolicyFontURL     = "/" + domain.OrgLabelPolicyFontPath
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
	prefix := domain.OrgLabelPolicyLogoPath
	if l.defaultPolicy {
		prefix = domain.DefaultLabelPolicyLogoPath
	}
	if l.darkMode {
		return prefix + "-" + domain.Dark + "-" + suffixID, nil
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
	prefix := domain.OrgLabelPolicyIconPath
	if l.defaultPolicy {
		prefix = domain.DefaultLabelPolicyIconPath
	}
	if l.darkMode {
		return prefix + "-" + domain.Dark + "-" + suffixID, nil
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
	prefix := domain.OrgLabelPolicyFontPath
	if l.defaultPolicy {
		prefix = domain.DefaultLabelPolicyFontPath
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
