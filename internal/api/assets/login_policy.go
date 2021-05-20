package assets

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/management/repository"
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

	preview = "/_preview"
)

type labelPolicyLogoUploader struct {
	idGenerator   id.Generator
	darkMode      bool
	defaultPolicy bool
}

func (l *labelPolicyLogoUploader) ObjectName(_ authz.CtxData) (string, error) {
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

func (l *labelPolicyLogoUploader) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
}

func (l *labelPolicyLogoUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
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

type labelPolicyLogoDownloader struct {
	org           repository.OrgRepository
	darkMode      bool
	defaultPolicy bool
	preview       bool
}

func (l *labelPolicyLogoDownloader) ObjectName(ctx context.Context) (string, error) {
	policy, err := getLabelPolicy(ctx, l.defaultPolicy, l.preview, l.org)
	if err != nil {
		return "", nil
	}
	if l.darkMode {
		return policy.LogoDarkURL, nil
	}
	return policy.LogoURL, nil
}

func (l *labelPolicyLogoDownloader) BucketName(ctx context.Context) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return authz.GetCtxData(ctx).OrgID
}

type labelPolicyIconUploader struct {
	idGenerator   id.Generator
	darkMode      bool
	defaultPolicy bool
}

func (l *labelPolicyIconUploader) ObjectName(_ authz.CtxData) (string, error) {
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

func (l *labelPolicyIconUploader) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
}

func (l *labelPolicyIconUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
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

type labelPolicyIconDownloader struct {
	org           repository.OrgRepository
	darkMode      bool
	defaultPolicy bool
	preview       bool
}

func (l *labelPolicyIconDownloader) ObjectName(ctx context.Context) (string, error) {
	policy, err := getLabelPolicy(ctx, l.defaultPolicy, l.preview, l.org)
	if err != nil {
		return "", nil
	}
	if l.darkMode {
		return policy.IconDarkURL, nil
	}
	return policy.IconURL, nil
}

func (l *labelPolicyIconDownloader) BucketName(ctx context.Context) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return authz.GetCtxData(ctx).OrgID
}

type labelPolicyFontUploader struct {
	idGenerator   id.Generator
	defaultPolicy bool
}

func (l *labelPolicyFontUploader) ObjectName(_ authz.CtxData) (string, error) {
	suffixID, err := l.idGenerator.Next()
	if err != nil {
		return "", err
	}
	prefix := domain.LabelPolicyFontPath
	return prefix + "-" + suffixID, nil
}

func (l *labelPolicyFontUploader) BucketName(ctxData authz.CtxData) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return ctxData.OrgID
}

func (l *labelPolicyFontUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	if l.defaultPolicy {
		_, err := commands.AddFontDefaultLabelPolicy(ctx, info.Key)
		return err
	}
	_, err := commands.AddFontLabelPolicy(ctx, orgID, info.Key)
	return err
}

type labelPolicyFontDownloader struct {
	org           repository.OrgRepository
	defaultPolicy bool
	preview       bool
}

func (l *labelPolicyFontDownloader) ObjectName(ctx context.Context) (string, error) {
	policy, err := getLabelPolicy(ctx, l.defaultPolicy, l.preview, l.org)
	if err != nil {
		return "", nil
	}
	return policy.FontURL, nil
}

func (l *labelPolicyFontDownloader) BucketName(ctx context.Context) string {
	if l.defaultPolicy {
		return domain.IAMID
	}
	return authz.GetCtxData(ctx).OrgID
}

func getLabelPolicy(ctx context.Context, defaultPolicy, preview bool, orgRepo repository.OrgRepository) (*model.LabelPolicyView, error) {
	if defaultPolicy {
		if preview {
			return orgRepo.GetPreviewDefaultLabelPolicy(ctx)
		}
		return orgRepo.GetDefaultLabelPolicy(ctx)
	}
	if preview {
		return orgRepo.GetPreviewLabelPolicy(ctx)
	}
	return orgRepo.GetLabelPolicy(ctx)
}
