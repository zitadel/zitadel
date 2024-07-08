package command

import (
	"context"
	"strings"

	"github.com/zitadel/saml/pkg/provider/xml"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddOrgGenericOAuthProvider(ctx context.Context, resourceOwner string, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGenericOAuthProvider(ctx context.Context, resourceOwner, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOAuthOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgOAuthProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGenericOIDCProvider(ctx context.Context, resourceOwner string, provider GenericOIDCProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOIDCOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgOIDCProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGenericOIDCProvider(ctx context.Context, resourceOwner, id string, provider GenericOIDCProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOIDCOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgOIDCProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) MigrateOrgGenericOIDCToAzureADProvider(ctx context.Context, resourceOwner, id string, provider AzureADProvider) (*domain.ObjectDetails, error) {
	return c.migrateOrgGenericOIDC(ctx, resourceOwner, id, provider)
}

func (c *Commands) MigrateOrgGenericOIDCToGoogleProvider(ctx context.Context, resourceOwner, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	return c.migrateOrgGenericOIDC(ctx, resourceOwner, id, provider)
}

func (c *Commands) migrateOrgGenericOIDC(ctx context.Context, resourceOwner, id string, provider interface{}) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewOIDCOrgIDPWriteModel(resourceOwner, id)

	var validation preparation.Validation
	switch p := provider.(type) {
	case AzureADProvider:
		validation = c.prepareMigrateOrgOIDCToAzureADProvider(orgAgg, writeModel, p)
	case GoogleProvider:
		validation = c.prepareMigrateOrgOIDCToGoogleProvider(orgAgg, writeModel, p)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-s9s2919", "Errors.IDPConfig.NotExisting")
	}

	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgJWTProvider(ctx context.Context, resourceOwner string, provider JWTProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewJWTOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgJWTProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgJWTProvider(ctx context.Context, resourceOwner, id string, provider JWTProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewJWTOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgJWTProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}
func (c *Commands) AddOrgAzureADProvider(ctx context.Context, resourceOwner string, provider AzureADProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewAzureADOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgAzureADProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgAzureADProvider(ctx context.Context, resourceOwner, id string, provider AzureADProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewAzureADOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgAzureADProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitHubProvider(ctx context.Context, resourceOwner string, provider GitHubProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitHubOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitHubProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitHubProvider(ctx context.Context, resourceOwner, id string, provider GitHubProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGitHubOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitHubProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitHubEnterpriseProvider(ctx context.Context, resourceOwner string, provider GitHubEnterpriseProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitHubEnterpriseOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitHubEnterpriseProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitHubEnterpriseProvider(ctx context.Context, resourceOwner, id string, provider GitHubEnterpriseProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGitHubEnterpriseOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitHubEnterpriseProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitLabProvider(ctx context.Context, resourceOwner string, provider GitLabProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitLabOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitLabProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitLabProvider(ctx context.Context, resourceOwner, id string, provider GitLabProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGitLabOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitLabProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGitLabSelfHostedProvider(ctx context.Context, resourceOwner string, provider GitLabSelfHostedProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitLabSelfHostedOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGitLabSelfHostedProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGitLabSelfHostedProvider(ctx context.Context, resourceOwner, id string, provider GitLabSelfHostedProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGitLabSelfHostedOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGitLabSelfHostedProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgGoogleProvider(ctx context.Context, resourceOwner string, provider GoogleProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgGoogleProvider(ctx context.Context, resourceOwner, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewGoogleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgGoogleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgLDAPProvider(ctx context.Context, resourceOwner string, provider LDAPProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgLDAPProvider(ctx context.Context, resourceOwner, id string, provider LDAPProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewLDAPOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgLDAPProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgSAMLProvider(ctx context.Context, resourceOwner string, provider SAMLProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewSAMLOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgSAMLProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgSAMLProvider(ctx context.Context, resourceOwner, id string, provider SAMLProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewSAMLOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgSAMLProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) RegenerateOrgSAMLProviderCertificate(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewSAMLOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareRegenerateOrgSAMLProviderCertificate(orgAgg, writeModel))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddOrgAppleProvider(ctx context.Context, resourceOwner string, provider AppleProvider) (string, *domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewAppleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddOrgAppleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateOrgAppleProvider(ctx context.Context, resourceOwner, id string, provider AppleProvider) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	writeModel := NewAppleOrgIDPWriteModel(resourceOwner, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateOrgAppleProvider(orgAgg, writeModel, provider))
	if err != nil {
		return nil, err
	}
	if len(cmds) == 0 {
		// no change, so return directly
		return &domain.ObjectDetails{
			Sequence:      writeModel.ProcessedSequence,
			EventDate:     writeModel.ChangeDate,
			ResourceOwner: writeModel.ResourceOwner,
		}, nil
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) DeleteOrgProvider(ctx context.Context, resourceOwner, id string) (*domain.ObjectDetails, error) {
	orgAgg := org.NewAggregate(resourceOwner)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareDeleteOrgProvider(orgAgg, resourceOwner, id))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func ExistsOrgIDP(ctx context.Context, filter preparation.FilterToQueryReducer, id, orgID string) (exists bool, err error) {
	writeModel := NewOrgIDPRemoveWriteModel(orgID, id)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return false, nil
	}
	writeModel.AppendEvents(events...)
	if err := writeModel.Reduce(); err != nil {
		return false, err
	}
	return writeModel.State.Exists(), nil
}

func (c *Commands) prepareAddOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-DF4ga", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		if provider.IDAttribute = strings.TrimSpace(provider.IDAttribute); provider.IDAttribute == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sadf3d", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOAuthIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.IDAttribute,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOAuthProvider(a *org.Aggregate, writeModel *OrgOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-asfsa", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Fb8jk", "Errors.Invalid.Argument")
		}
		if provider.IDAttribute = strings.TrimSpace(provider.IDAttribute); provider.IDAttribute == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAe4gh", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-JNsd3", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.IDAttribute,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgOIDCProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Sgtj5", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Hz6zj", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-fb5jm", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Sfdf4", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOIDCIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Issuer,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IsIDTokenMapping,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgOIDCProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAfd3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dvf4f", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-BDfr3", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Db3bs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Dg331", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Issuer,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IsIDTokenMapping,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareMigrateOrgOIDCToAzureADProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdf3g", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Fhbr2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dzh3g", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INST-Dg239201", "Errors.Instance.IDPConfig.NotExisting")
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOIDCIDPMigratedAzureADEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.Tenant,
					provider.EmailVerified,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareMigrateOrgOIDCToGoogleProvider(a *org.Aggregate, writeModel *OrgOIDCIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-W2vqs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INST-x09981", "Errors.Instance.IDPConfig.NotExisting")
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewOIDCIDPMigratedGoogleEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgJWTProvider(a *org.Aggregate, writeModel *OrgJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-JLKef", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-WNJK3", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-NJKSD", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-NJKE3", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-2rlks", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewJWTIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Issuer,
					provider.JWTEndpoint,
					provider.KeyEndpoint,
					provider.HeaderName,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgJWTProvider(a *org.Aggregate, writeModel *OrgJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-HUe3q", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-JKLS2", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-JKs3f", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-NJKS2", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SJk2d", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SJK2f", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Bhju5", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Issuer,
				provider.JWTEndpoint,
				provider.KeyEndpoint,
				provider.HeaderName,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgAzureADProvider(a *org.Aggregate, writeModel *OrgAzureADIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdf3g", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Fhbr2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dzh3g", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewAzureADIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.Tenant,
					provider.EmailVerified,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgAzureADProvider(a *org.Aggregate, writeModel *OrgAzureADIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAgh2", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-fh3h1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-dmitg", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-BHz3q", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.Tenant,
				provider.EmailVerified,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitHubProvider(a *org.Aggregate, writeModel *OrgGitHubIDPWriteModel, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Jdsgf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-dsgz3", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitHubIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitHubProvider(a *org.Aggregate, writeModel *OrgGitHubIDPWriteModel, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdf4h", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-fdh5z", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Dr1gs", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitHubEnterpriseProvider(a *org.Aggregate, writeModel *OrgGitHubEnterpriseIDPWriteModel, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dg4td", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-dgj53", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Ghjjs", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sani2", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-agj42", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sd5hn", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitHubEnterpriseIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.AuthorizationEndpoint,
					provider.TokenEndpoint,
					provider.UserEndpoint,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitHubEnterpriseProvider(a *org.Aggregate, writeModel *OrgGitHubEnterpriseIDPWriteModel, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdfh3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-shj42", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdh73", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-acx2w", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-dgj6q", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-ybj62", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-GBr42", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.AuthorizationEndpoint,
				provider.TokenEndpoint,
				provider.UserEndpoint,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitLabProvider(a *org.Aggregate, writeModel *OrgGitLabIDPWriteModel, provider GitLabProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-adsg2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-GD1j2", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitLabIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitLabProvider(a *org.Aggregate, writeModel *OrgGitLabIDPWriteModel, provider GitLabProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-HJK91", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D12t6", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-HBReq", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGitLabSelfHostedProvider(a *org.Aggregate, writeModel *OrgGitLabSelfHostedIDPWriteModel, provider GitLabSelfHostedProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-jw4ZT", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-AST4S", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-DBZHJ", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SDGJ4", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGitLabSelfHostedIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Issuer,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGitLabSelfHostedProvider(a *org.Aggregate, writeModel *OrgGitLabSelfHostedIDPWriteModel, provider GitLabSelfHostedProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAFG4", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-DG4H", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SD4eb", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-GHWE3", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-D2tg1", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Issuer,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil {
				return nil, err
			}
			if event == nil {
				return nil, nil
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-W2vqs", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewGoogleIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					secret,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgGoogleProvider(a *org.Aggregate, writeModel *OrgGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-S32t1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-ds432", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Dqrg1", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.ClientSecret,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAfdd", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sv31s", "Errors.Invalid.Argument")
		}
		if provider.BindDN = strings.TrimSpace(provider.BindDN); provider.BindDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-sdgf4", "Errors.Invalid.Argument")
		}
		if provider.BindPassword = strings.TrimSpace(provider.BindPassword); provider.BindPassword == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-AEG2w", "Errors.Invalid.Argument")
		}
		if provider.UserBase = strings.TrimSpace(provider.UserBase); provider.UserBase == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAD5n", "Errors.Invalid.Argument")
		}
		if len(provider.Servers) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SAy945n", "Errors.Invalid.Argument")
		}
		if len(provider.UserObjectClasses) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-S1x705n", "Errors.Invalid.Argument")
		}
		if len(provider.UserFilters) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-aAx9x1n", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			secret, err := crypto.Encrypt([]byte(provider.BindPassword), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewLDAPIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Servers,
					provider.StartTLS,
					provider.BaseDN,
					provider.BindDN,
					secret,
					provider.UserBase,
					provider.UserObjectClasses,
					provider.UserFilters,
					provider.Timeout,
					provider.LDAPAttributes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgLDAPProvider(a *org.Aggregate, writeModel *OrgLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Dgdbs", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Sffgd", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-vb3ss", "Errors.Invalid.Argument")
		}
		if provider.BindDN = strings.TrimSpace(provider.BindDN); provider.BindDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-hbere", "Errors.Invalid.Argument")
		}
		if provider.UserBase = strings.TrimSpace(provider.UserBase); provider.UserBase == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-DG45z", "Errors.Invalid.Argument")
		}
		if len(provider.Servers) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Sxx945n", "Errors.Invalid.Argument")
		}
		if len(provider.UserObjectClasses) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-S1p605n", "Errors.Invalid.Argument")
		}
		if len(provider.UserFilters) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-aBx901n", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-ASF3F", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Servers,
				provider.StartTLS,
				provider.BaseDN,
				provider.BindDN,
				provider.BindPassword,
				provider.UserBase,
				provider.UserObjectClasses,
				provider.UserFilters,
				provider.Timeout,
				c.idpConfigEncryption,
				provider.LDAPAttributes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgAppleProvider(a *org.Aggregate, writeModel *OrgAppleIDPWriteModel, provider AppleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-jkn3w", "Errors.IDP.ClientIDMissing")
		}
		if provider.TeamID = strings.TrimSpace(provider.TeamID); provider.TeamID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Ffg32", "Errors.IDP.TeamIDMissing")
		}
		if provider.KeyID = strings.TrimSpace(provider.KeyID); provider.KeyID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-GDjm5", "Errors.IDP.KeyIDMissing")
		}
		if len(provider.PrivateKey) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-GVD4n", "Errors.IDP.PrivateKeyMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			privateKey, err := crypto.Encrypt(provider.PrivateKey, c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewAppleIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.ClientID,
					provider.TeamID,
					provider.KeyID,
					privateKey,
					provider.Scopes,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgAppleProvider(a *org.Aggregate, writeModel *OrgAppleIDPWriteModel, provider AppleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-FRHBH", "Errors.IDMissing")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SFm4l", "Errors.IDP.ClientIDMissing")
		}
		if provider.TeamID = strings.TrimSpace(provider.TeamID); provider.TeamID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-SG34t", "Errors.IDP.TeamIDMissing")
		}
		if provider.KeyID = strings.TrimSpace(provider.KeyID); provider.KeyID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-Gh4z2", "Errors.IDP.KeyIDMissing")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-SG3bh", "Errors.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.ClientID,
				provider.TeamID,
				provider.KeyID,
				provider.PrivateKey,
				c.idpConfigEncryption,
				provider.Scopes,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddOrgSAMLProvider(a *org.Aggregate, writeModel *OrgSAMLIDPWriteModel, provider SAMLProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-957lr0f8u3", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-78isv6m53a", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL != "" {
			data, err := xml.ReadMetadataFromURL(c.httpClient, provider.MetadataURL)
			if err != nil {
				return nil, zerrors.ThrowInvalidArgument(err, "ORG-ipzxvf3cv2", "Errors.Project.App.SAMLMetadataMissing")
			}
			provider.Metadata = data
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			key, cert, err := c.samlCertificateAndKeyGenerator(writeModel.ID)
			if err != nil {
				return nil, err
			}
			keyEnc, err := crypto.Encrypt(key, c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				org.NewSAMLIDPAddedEvent(
					ctx,
					&a.Aggregate,
					writeModel.ID,
					provider.Name,
					provider.Metadata,
					keyEnc,
					cert,
					provider.Binding,
					provider.WithSignedRequest,
					provider.NameIDFormat,
					provider.TransientMappingAttributeName,
					provider.IDPOptions,
				),
			}, nil
		}, nil
	}
}

func (c *Commands) prepareUpdateOrgSAMLProvider(a *org.Aggregate, writeModel *OrgSAMLIDPWriteModel, provider SAMLProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-wwdwdlaya0", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-egixaofgyl", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL != "" {
			data, err := xml.ReadMetadataFromURL(c.httpClient, provider.MetadataURL)
			if err != nil {
				return nil, zerrors.ThrowInvalidArgument(err, "ORG-bkaiyd3rfo", "Errors.Project.App.SAMLMetadataMissing")
			}
			provider.Metadata = data
		}
		if provider.Metadata == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-j6spncd74m", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-z82dddndql", "Errors.Org.IDPConfig.NotExisting")
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				provider.Name,
				provider.Metadata,
				nil,
				nil,
				c.idpConfigEncryption,
				provider.Binding,
				provider.WithSignedRequest,
				provider.NameIDFormat,
				provider.TransientMappingAttributeName,
				provider.IDPOptions,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareRegenerateOrgSAMLProviderCertificate(a *org.Aggregate, writeModel *OrgSAMLIDPWriteModel) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "ORG-arv4vdrb6c", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-4dw21ch9o9", "Errors.Org.IDPConfig.NotExisting")
			}

			key, cert, err := c.samlCertificateAndKeyGenerator(writeModel.ID)
			if err != nil {
				return nil, err
			}
			event, err := writeModel.NewChangedEvent(
				ctx,
				&a.Aggregate,
				writeModel.ID,
				writeModel.Name,
				writeModel.Metadata,
				key,
				cert,
				c.idpConfigEncryption,
				writeModel.Binding,
				writeModel.WithSignedRequest,
				writeModel.NameIDFormat,
				writeModel.TransientMappingAttributeName,
				writeModel.Options,
			)
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareDeleteOrgProvider(a *org.Aggregate, resourceOwner, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewOrgIDPRemoveWriteModel(resourceOwner, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "ORG-Se3tg", "Errors.Org.IDPConfig.NotExisting")
			}
			return []eventstore.Command{org.NewIDPRemovedEvent(ctx, &a.Aggregate, id)}, nil
		}, nil
	}
}
