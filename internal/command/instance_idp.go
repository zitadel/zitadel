package command

import (
	"context"
	"strings"

	"github.com/zitadel/saml/pkg/provider/xml"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddInstanceGenericOAuthProvider(ctx context.Context, provider GenericOAuthProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOAuthInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOAuthProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGenericOAuthProvider(ctx context.Context, id string, provider GenericOAuthProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewOAuthInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceOAuthProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGenericOIDCProvider(ctx context.Context, provider GenericOIDCProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewOIDCInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceOIDCProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGenericOIDCProvider(ctx context.Context, id string, provider GenericOIDCProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewOIDCInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceOIDCProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) MigrateInstanceGenericOIDCToAzureADProvider(ctx context.Context, id string, provider AzureADProvider) (*domain.ObjectDetails, error) {
	return c.migrateInstanceGenericOIDC(ctx, id, provider)
}

func (c *Commands) MigrateInstanceGenericOIDCToGoogleProvider(ctx context.Context, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	return c.migrateInstanceGenericOIDC(ctx, id, provider)
}

func (c *Commands) migrateInstanceGenericOIDC(ctx context.Context, id string, provider interface{}) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewOIDCInstanceIDPWriteModel(instanceID, id)

	var validation preparation.Validation
	switch p := provider.(type) {
	case AzureADProvider:
		validation = c.prepareMigrateInstanceOIDCToAzureADProvider(instanceAgg, writeModel, p)
	case GoogleProvider:
		validation = c.prepareMigrateInstanceOIDCToGoogleProvider(instanceAgg, writeModel, p)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-s9219", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) AddInstanceJWTProvider(ctx context.Context, provider JWTProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewJWTInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceJWTProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceJWTProvider(ctx context.Context, id string, provider JWTProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewJWTInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceJWTProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceAzureADProvider(ctx context.Context, provider AzureADProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewAzureADInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceAzureADProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceAzureADProvider(ctx context.Context, id string, provider AzureADProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewAzureADInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceAzureADProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGitHubProvider(ctx context.Context, provider GitHubProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitHubInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGitHubProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGitHubProvider(ctx context.Context, id string, provider GitHubProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewGitHubInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGitHubProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGitHubEnterpriseProvider(ctx context.Context, provider GitHubEnterpriseProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitHubEnterpriseInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGitHubEnterpriseProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGitHubEnterpriseProvider(ctx context.Context, id string, provider GitHubEnterpriseProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewGitHubEnterpriseInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGitHubEnterpriseProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGitLabProvider(ctx context.Context, provider GitLabProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitLabInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGitLabProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGitLabProvider(ctx context.Context, id string, provider GitLabProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewGitLabInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGitLabProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGitLabSelfHostedProvider(ctx context.Context, provider GitLabSelfHostedProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGitLabSelfHostedInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGitLabSelfHostedProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGitLabSelfHostedProvider(ctx context.Context, id string, provider GitLabSelfHostedProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewGitLabSelfHostedInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGitLabSelfHostedProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceGoogleProvider(ctx context.Context, provider GoogleProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewGoogleInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceGoogleProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceGoogleProvider(ctx context.Context, id string, provider GoogleProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewGoogleInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceGoogleProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceLDAPProvider(ctx context.Context, provider LDAPProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewLDAPInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceLDAPProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceLDAPProvider(ctx context.Context, id string, provider LDAPProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewLDAPInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceLDAPProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceAppleProvider(ctx context.Context, provider AppleProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewAppleInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceAppleProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceAppleProvider(ctx context.Context, id string, provider AppleProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewAppleInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceAppleProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) AddInstanceSAMLProvider(ctx context.Context, provider SAMLProvider) (string, *domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	id, err := id_generator.Next()
	if err != nil {
		return "", nil, err
	}
	writeModel := NewSAMLInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareAddInstanceSAMLProvider(instanceAgg, writeModel, provider))
	if err != nil {
		return "", nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return "", nil, err
	}
	return id, pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) UpdateInstanceSAMLProvider(ctx context.Context, id string, provider SAMLProvider) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewSAMLInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareUpdateInstanceSAMLProvider(instanceAgg, writeModel, provider))
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

func (c *Commands) RegenerateInstanceSAMLProviderCertificate(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	instanceAgg := instance.NewAggregate(instanceID)
	writeModel := NewSAMLInstanceIDPWriteModel(instanceID, id)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareRegenerateInstanceSAMLProviderCertificate(instanceAgg, writeModel))
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

func (c *Commands) DeleteInstanceProvider(ctx context.Context, id string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.prepareDeleteInstanceProvider(instanceAgg, id))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func ExistsInstanceIDP(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (exists bool, err error) {
	instanceWriteModel := NewInstanceIDPRemoveWriteModel(authz.GetInstance(ctx).InstanceID(), id)
	events, err := filter(ctx, instanceWriteModel.Query())
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return false, nil
	}
	instanceWriteModel.AppendEvents(events...)
	if err := instanceWriteModel.Reduce(); err != nil {
		return false, err
	}
	return instanceWriteModel.State.Exists(), nil
}

func (c *Commands) prepareAddInstanceOAuthProvider(a *instance.Aggregate, writeModel *InstanceOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-D32ef", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dbgzf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-DF4ga", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-B23bs", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-D2gj8", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Fb8jk", "Errors.Invalid.Argument")
		}
		if provider.IDAttribute = strings.TrimSpace(provider.IDAttribute); provider.IDAttribute == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdf3f", "Errors.Invalid.Argument")
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
				instance.NewOAuthIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceOAuthProvider(a *instance.Aggregate, writeModel *InstanceOAuthIDPWriteModel, provider GenericOAuthProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAffg", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Sf3gh", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SHJ3ui", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SVrgh", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-DJKeio", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-ILSJi", "Errors.Invalid.Argument")
		}
		if provider.IDAttribute = strings.TrimSpace(provider.IDAttribute); provider.IDAttribute == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-JKD3h", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-D3r1s", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceOIDCProvider(a *instance.Aggregate, writeModel *InstanceOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Sgtj5", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Hz6zj", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-fb5jm", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Sfdf4", "Errors.Invalid.Argument")
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
				instance.NewOIDCIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceOIDCProvider(a *instance.Aggregate, writeModel *InstanceOIDCIDPWriteModel, provider GenericOIDCProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAfd3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dvf4f", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-BDfr3", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Db3bs", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-Dg331", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareMigrateInstanceOIDCToAzureADProvider(a *instance.Aggregate, writeModel *InstanceOIDCIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdf3g", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Fhbr2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dzh3g", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-Dg29201", "Errors.IDPConfig.NotExisting")
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				instance.NewOIDCIDPMigratedAzureADEvent(
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

func (c *Commands) prepareMigrateInstanceOIDCToGoogleProvider(a *instance.Aggregate, writeModel *InstanceOIDCIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-W2vqs", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-Dg29202", "Errors.IDPConfig.NotExisting")
			}
			secret, err := crypto.Encrypt([]byte(provider.ClientSecret), c.idpConfigEncryption)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{
				instance.NewOIDCIDPMigratedGoogleEvent(
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

func (c *Commands) prepareAddInstanceJWTProvider(a *instance.Aggregate, writeModel *InstanceJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-JLKef", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-WNJK3", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-NJKSD", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-NJKE3", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-2rlks", "Errors.Invalid.Argument")
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
				instance.NewJWTIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceJWTProvider(a *instance.Aggregate, writeModel *InstanceJWTIDPWriteModel, provider JWTProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-HUe3q", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-JKLS2", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-JKs3f", "Errors.Invalid.Argument")
		}
		if provider.JWTEndpoint = strings.TrimSpace(provider.JWTEndpoint); provider.JWTEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-NJKS2", "Errors.Invalid.Argument")
		}
		if provider.KeyEndpoint = strings.TrimSpace(provider.KeyEndpoint); provider.KeyEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SJk2d", "Errors.Invalid.Argument")
		}
		if provider.HeaderName = strings.TrimSpace(provider.HeaderName); provider.HeaderName == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SJK2f", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-Bhju5", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceAzureADProvider(a *instance.Aggregate, writeModel *InstanceAzureADIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdf3g", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Fhbr2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dzh3g", "Errors.Invalid.Argument")
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
				instance.NewAzureADIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceAzureADProvider(a *instance.Aggregate, writeModel *InstanceAzureADIDPWriteModel, provider AzureADProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAgh2", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-fh3h1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-dmitg", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-BHz3q", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceGitHubProvider(a *instance.Aggregate, writeModel *InstanceGitHubIDPWriteModel, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Jdsgf", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-dsgz3", "Errors.Invalid.Argument")
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
				instance.NewGitHubIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceGitHubProvider(a *instance.Aggregate, writeModel *InstanceGitHubIDPWriteModel, provider GitHubProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdf4h", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-fdh5z", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-Dr1gs", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceGitHubEnterpriseProvider(a *instance.Aggregate, writeModel *InstanceGitHubEnterpriseIDPWriteModel, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dg4td", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-dgj53", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ghjjs", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sani2", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-agj42", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sd5hn", "Errors.Invalid.Argument")
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
				instance.NewGitHubEnterpriseIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceGitHubEnterpriseProvider(a *instance.Aggregate, writeModel *InstanceGitHubEnterpriseIDPWriteModel, provider GitHubEnterpriseProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdfh3", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-shj42", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdh73", "Errors.Invalid.Argument")
		}
		if provider.AuthorizationEndpoint = strings.TrimSpace(provider.AuthorizationEndpoint); provider.AuthorizationEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-acx2w", "Errors.Invalid.Argument")
		}
		if provider.TokenEndpoint = strings.TrimSpace(provider.TokenEndpoint); provider.TokenEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-dgj6q", "Errors.Invalid.Argument")
		}
		if provider.UserEndpoint = strings.TrimSpace(provider.UserEndpoint); provider.UserEndpoint == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-ybj62", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-GBr42", "Errors.IDPConfig.NotExisting")
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
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddInstanceGitLabProvider(a *instance.Aggregate, writeModel *InstanceGitLabIDPWriteModel, provider GitLabProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-adsg2", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-GD1j2", "Errors.Invalid.Argument")
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
				instance.NewGitLabIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceGitLabProvider(a *instance.Aggregate, writeModel *InstanceGitLabIDPWriteModel, provider GitLabProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-HJK91", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-D12t6", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-HBReq", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceGitLabSelfHostedProvider(a *instance.Aggregate, writeModel *InstanceGitLabSelfHostedIDPWriteModel, provider GitLabSelfHostedProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-jw4ZT", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-AST4S", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-DBZHJ", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SDGJ4", "Errors.Invalid.Argument")
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
				instance.NewGitLabSelfHostedIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceGitLabSelfHostedProvider(a *instance.Aggregate, writeModel *InstanceGitLabSelfHostedIDPWriteModel, provider GitLabSelfHostedProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAFG4", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-DG4H", "Errors.Invalid.Argument")
		}
		if provider.Issuer = strings.TrimSpace(provider.Issuer); provider.Issuer == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SD4eb", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-GHWE3", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-D2tg1", "Errors.IDPConfig.NotExisting")
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
			if err != nil || event == nil {
				return nil, err
			}
			return []eventstore.Command{event}, nil
		}, nil
	}
}

func (c *Commands) prepareAddInstanceGoogleProvider(a *instance.Aggregate, writeModel *InstanceGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-D3fvs", "Errors.Invalid.Argument")
		}
		if provider.ClientSecret = strings.TrimSpace(provider.ClientSecret); provider.ClientSecret == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-W2vqs", "Errors.Invalid.Argument")
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
				instance.NewGoogleIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceGoogleProvider(a *instance.Aggregate, writeModel *InstanceGoogleIDPWriteModel, provider GoogleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-S32t1", "Errors.Invalid.Argument")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-ds432", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-D3r1s", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceLDAPProvider(a *instance.Aggregate, writeModel *InstanceLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAfdd", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sv31s", "Errors.Invalid.Argument")
		}
		if provider.BindDN = strings.TrimSpace(provider.BindDN); provider.BindDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-sdgf4", "Errors.Invalid.Argument")
		}
		if provider.BindPassword = strings.TrimSpace(provider.BindPassword); provider.BindPassword == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-AEG2w", "Errors.Invalid.Argument")
		}
		if provider.UserBase = strings.TrimSpace(provider.UserBase); provider.UserBase == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAD5n", "Errors.Invalid.Argument")
		}
		if len(provider.Servers) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAx905n", "Errors.Invalid.Argument")
		}
		if len(provider.UserObjectClasses) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-S1x905n", "Errors.Invalid.Argument")
		}
		if len(provider.UserFilters) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-aAx905n", "Errors.Invalid.Argument")
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
				instance.NewLDAPIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceLDAPProvider(a *instance.Aggregate, writeModel *InstanceLDAPIDPWriteModel, provider LDAPProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Dgdbs", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Sffgd", "Errors.Invalid.Argument")
		}
		if provider.BaseDN = strings.TrimSpace(provider.BaseDN); provider.BaseDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-vb3ss", "Errors.Invalid.Argument")
		}
		if provider.BindDN = strings.TrimSpace(provider.BindDN); provider.BindDN == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-hbere", "Errors.Invalid.Argument")
		}
		if provider.UserBase = strings.TrimSpace(provider.UserBase); provider.UserBase == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-DG45z", "Errors.Invalid.Argument")
		}
		if len(provider.Servers) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SAx945n", "Errors.Invalid.Argument")
		}
		if len(provider.UserObjectClasses) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-S1x605n", "Errors.Invalid.Argument")
		}
		if len(provider.UserFilters) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-aAx901n", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-ASF3F", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceAppleProvider(a *instance.Aggregate, writeModel *InstanceAppleIDPWriteModel, provider AppleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-jkn3w", "Errors.IDP.ClientIDMissing")
		}
		if provider.TeamID = strings.TrimSpace(provider.TeamID); provider.TeamID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Ffg32", "Errors.IDP.TeamIDMissing")
		}
		if provider.KeyID = strings.TrimSpace(provider.KeyID); provider.KeyID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-GDjm5", "Errors.IDP.KeyIDMissing")
		}
		if len(provider.PrivateKey) == 0 {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-GVD4n", "Errors.IDP.PrivateKeyMissing")
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
				instance.NewAppleIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceAppleProvider(a *instance.Aggregate, writeModel *InstanceAppleIDPWriteModel, provider AppleProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-FRHBH", "Errors.IDMissing")
		}
		if provider.ClientID = strings.TrimSpace(provider.ClientID); provider.ClientID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SFm4l", "Errors.IDP.ClientIDMissing")
		}
		if provider.TeamID = strings.TrimSpace(provider.TeamID); provider.TeamID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-SG34t", "Errors.IDP.TeamIDMissing")
		}
		if provider.KeyID = strings.TrimSpace(provider.KeyID); provider.KeyID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-Gh4z2", "Errors.IDP.KeyIDMissing")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-SG3bh", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareAddInstanceSAMLProvider(a *instance.Aggregate, writeModel *InstanceSAMLIDPWriteModel, provider SAMLProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-o07zjotgnd", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL != "" {
			data, err := xml.ReadMetadataFromURL(c.httpClient, provider.MetadataURL)
			if err != nil {
				return nil, zerrors.ThrowInvalidArgument(err, "INST-8vam1khq22", "Errors.Project.App.SAMLMetadataMissing")
			}
			provider.Metadata = data
		}
		if provider.Metadata == nil {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-3bi3esi16t", "Errors.Invalid.Argument")
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
				instance.NewSAMLIDPAddedEvent(
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

func (c *Commands) prepareUpdateInstanceSAMLProvider(a *instance.Aggregate, writeModel *InstanceSAMLIDPWriteModel, provider SAMLProvider) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-7o3rq1owpm", "Errors.Invalid.Argument")
		}
		if provider.Name = strings.TrimSpace(provider.Name); provider.Name == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-q2s9rak7o9", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-iw1rxnf4sf", "Errors.Invalid.Argument")
		}
		if provider.Metadata == nil && provider.MetadataURL != "" {
			data, err := xml.ReadMetadataFromURL(c.httpClient, provider.MetadataURL)
			if err != nil {
				return nil, zerrors.ThrowInvalidArgument(err, "INST-iijz4h01if", "Errors.Project.App.SAMLMetadataMissing")
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
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INST-D3r1s", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareRegenerateInstanceSAMLProviderCertificate(a *instance.Aggregate, writeModel *InstanceSAMLIDPWriteModel) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if writeModel.ID = strings.TrimSpace(writeModel.ID); writeModel.ID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "INST-7de108gqya", "Errors.Invalid.Argument")
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
				return nil, zerrors.ThrowNotFound(nil, "INST-76dbwsv9vm", "Errors.IDPConfig.NotExisting")
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

func (c *Commands) prepareDeleteInstanceProvider(a *instance.Aggregate, id string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceIDPRemoveWriteModel(a.InstanceID, id)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INST-Se3tg", "Errors.IDPConfig.NotExisting")
			}
			return []eventstore.Command{instance.NewIDPRemovedEvent(ctx, &a.Aggregate, id)}, nil
		}, nil
	}
}
