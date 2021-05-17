package upload

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
)

const (
	defaultLabelPolicyLogoURL     = "/iam/" + domain.LabelPolicyLogoPath
	defaultLabelPolicyLogoDarkURL = "/iam/" + domain.LabelPolicyLogoPath + "/" + domain.Dark
	defaultLabelPolicyIconURL     = "/iam/" + domain.LabelPolicyIconPath
	defaultLabelPolicyIconDarkURL = "/iam/" + domain.LabelPolicyIconPath + "/" + domain.Dark
	defaultLabelPolicyFontURL     = "/iam/" + domain.LabelPolicyFontPath

	orgLabelPolicyLogoURL     = "/org/" + domain.LabelPolicyLogoPath
	orgLabelPolicyLogoDarkURL = "/org/" + domain.LabelPolicyLogoPath + "/" + domain.Dark
	orgLabelPolicyIconURL     = "/org/" + domain.LabelPolicyIconPath
	orgLabelPolicyIconDarkURL = "/org/" + domain.LabelPolicyIconPath + "/" + domain.Dark
	orgLabelPolicyFontURL     = "/org/" + domain.LabelPolicyFontPath
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
	prefix := domain.LabelPolicyLogoPath
	if l.darkMode {
		return prefix + "-" + domain.Dark + "-" + suffixID, nil
	}
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyLogo) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
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
	prefix := domain.LabelPolicyIconPath
	if l.darkMode {
		return prefix + "-" + domain.Dark + "-" + suffixID, nil
	}
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyIcon) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
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
	prefix := domain.LabelPolicyFontPath
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyFont) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
}

func (l *labelPolicyFont) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.defaultPolicy {
		_, err := commands.AddFontDefaultLabelPolicy(ctx, info.Key)
		return err
	}
	_, err := commands.AddFontLabelPolicy(ctx, orgID, info.Key)
	return err
}
