package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type PrivacyPolicyWriteModel struct {
	eventstore.WriteModel

	TOSLink        string
	PrivacyLink    string
	HelpLink       string
	SupportEmail   domain.EmailAddress
	State          domain.PolicyState
	DocsLink       string
	CustomLink     string
	CustomLinkText string
}

func (wm *PrivacyPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PrivacyPolicyAddedEvent:
			wm.TOSLink = e.TOSLink
			wm.PrivacyLink = e.PrivacyLink
			wm.HelpLink = e.HelpLink
			wm.SupportEmail = e.SupportEmail
			wm.State = domain.PolicyStateActive
			wm.DocsLink = e.DocsLink
			wm.CustomLink = e.CustomLink
			wm.CustomLinkText = e.CustomLinkText
		case *policy.PrivacyPolicyChangedEvent:
			if e.PrivacyLink != nil {
				wm.PrivacyLink = *e.PrivacyLink
			}
			if e.TOSLink != nil {
				wm.TOSLink = *e.TOSLink
			}
			if e.HelpLink != nil {
				wm.HelpLink = *e.HelpLink
			}
			if e.SupportEmail != nil {
				wm.SupportEmail = *e.SupportEmail
			}
			if e.DocsLink != nil {
				wm.DocsLink = *e.DocsLink
			}
			if e.CustomLink != nil {
				wm.CustomLink = *e.CustomLink
			}
			if e.CustomLinkText != nil {
				wm.CustomLinkText = *e.CustomLinkText
			}
		case *policy.PrivacyPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
