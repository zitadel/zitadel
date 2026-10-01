package command

import (
	"context"
	"slices"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SecurityPolicy struct {
	EnableIframeEmbedding bool
	AllowedOrigins        []string
	EnableImpersonation   bool

	EnableDynamicClientRegistration               bool
	AllowUnauthenticatedDynamicClientRegistration bool

	EnableClientIDMetadataDocument bool
	// ClientIDMetadataDocumentAllowedURLs are the client_id URLs the instance resolves as Client
	// ID Metadata Documents, within the URLs the system allows.
	ClientIDMetadataDocumentAllowedURLs []string
	// ClientIDMetadataDocumentAllowAnyURL lets the instance resolve every client_id URL the system
	// allows, regardless of ClientIDMetadataDocumentAllowedURLs.
	ClientIDMetadataDocumentAllowAnyURL bool
}

func (c *Commands) SetSecurityPolicy(ctx context.Context, policy *SecurityPolicy) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	validation := c.prepareSetSecurityPolicy(instanceAgg, policy)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, validation)
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreatedAt(),
		ResourceOwner: events[len(events)-1].Aggregate().InstanceID,
	}, nil
}

func (c *Commands) prepareSetSecurityPolicy(a *instance.Aggregate, policy *SecurityPolicy) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		for _, allowedURL := range policy.ClientIDMetadataDocumentAllowedURLs {
			if err := validateClientIDMetadataDocumentAllowedURL(allowedURL); err != nil {
				return nil, err
			}
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := c.getSecurityPolicyWriteModel(ctx, filter)
			if err != nil {
				return nil, err
			}
			cmd, err := writeModel.NewSetEvent(ctx, &a.Aggregate, policy)
			if err != nil {
				return nil, err
			}
			return []eventstore.Command{cmd}, nil
		}, nil
	}
}

func (c *Commands) getSecurityPolicyWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer) (_ *InstanceSecurityPolicyWriteModel, err error) {
	writeModel := NewInstanceSecurityPolicyWriteModel(ctx)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return writeModel, nil
	}
	writeModel.AppendEvents(events...)
	err = writeModel.Reduce()
	return writeModel, err
}

// AddClientIDMetadataDocumentAllowedURL adds a client_id URL the instance resolves as a Client ID
// Metadata Document. An entry ending with a slash allows every client_id it is a prefix of. The
// URL must not be allowed already.
func (c *Commands) AddClientIDMetadataDocumentAllowedURL(ctx context.Context, allowedURL string) (*domain.ObjectDetails, error) {
	if err := validateClientIDMetadataDocumentAllowedURL(allowedURL); err != nil {
		return nil, err
	}
	writeModel := NewInstanceSecurityPolicyWriteModel(ctx)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if slices.Contains(writeModel.ClientIDMetadataDocumentAllowedURLs, allowedURL) {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMA-Vq7tY", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.AlreadyExists")
	}
	err := c.pushAppendAndReduce(ctx, writeModel, instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(ctx, InstanceAggregateFromWriteModel(&writeModel.WriteModel), allowedURL))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

// RemoveClientIDMetadataDocumentAllowedURL removes a client_id URL the instance resolves as a
// Client ID Metadata Document. Removing a URL that is not allowed is not an error and returns no
// details.
func (c *Commands) RemoveClientIDMetadataDocumentAllowedURL(ctx context.Context, allowedURL string) (*domain.ObjectDetails, error) {
	writeModel := NewInstanceSecurityPolicyWriteModel(ctx)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}
	if !slices.Contains(writeModel.ClientIDMetadataDocumentAllowedURLs, allowedURL) {
		return nil, nil
	}
	err := c.pushAppendAndReduce(ctx, writeModel, instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(ctx, InstanceAggregateFromWriteModel(&writeModel.WriteModel), allowedURL))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&writeModel.WriteModel), nil
}

// validateClientIDMetadataDocumentAllowedURL requires allowedURL to be a valid client_id URL, so
// that an entry ending with a slash cannot be widened by a client_id that is not one.
func validateClientIDMetadataDocumentAllowedURL(allowedURL string) error {
	if !domain.IsClientIDMetadataDocumentURL(allowedURL) {
		return zerrors.ThrowInvalidArgument(nil, "COMMA-p4Lz8", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.Invalid")
	}
	return nil
}
