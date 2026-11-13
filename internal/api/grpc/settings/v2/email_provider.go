package settings

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

// emailProviderRestrictions returns the restrictions the operator rules define for the active email provider.
// It returns nil if no provider is active or the active provider is not restricted.
func (s *Server) emailProviderRestrictions(ctx context.Context) (*settings.EmailProviderRestrictions, error) {
	if len(s.smtpRules) == 0 {
		return nil, nil
	}
	config, err := s.query.SMTPConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		if zerrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return emailProviderRestrictionsToPb(config, s.smtpRules), nil
}

func emailProviderRestrictionsToPb(config *query.SMTPConfig, rules smtp.Rules) *settings.EmailProviderRestrictions {
	if config.SMTPConfig == nil {
		return nil
	}
	options := rules.Options(config.SMTPConfig.Host, config.SMTPConfig.User, config.SMTPConfig.SenderAddress)
	if options == (smtp.RuleOptions{}) {
		return nil
	}
	restrictions := &settings.EmailProviderRestrictions{
		CustomHtmlRestricted:               options.RestrictCustomHTML,
		ReservedRecipientDomainsSuppressed: options.SuppressReservedRecipientDomains,
	}
	if options.Limit != nil {
		restrictions.SendingLimit = &settings.EmailProviderSendingLimit{
			// the count is validated to fit when the rules are compiled
			Count:  uint32(options.Limit.Count), //nolint:gosec
			Window: durationpb.New(time.Duration(options.Limit.Window)),
		}
	}
	return restrictions
}
