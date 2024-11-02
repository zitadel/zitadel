package command

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type InstanceFeatures struct {
	LoginDefaultOrg                 *bool
	TriggerIntrospectionProjections *bool
	LegacyIntrospection             *bool
	UserSchema                      *bool
	TokenExchange                   *bool
	Actions                         *bool
	ImprovedPerformance             []feature.ImprovedPerformanceType
	WebKey                          *bool
	DebugOIDCParentError            *bool
	OIDCSingleV1SessionTermination  *bool
	DisableUserTokenEvent           *bool
	EnableBackChannelLogout         *bool
}

func (m *InstanceFeatures) isEmpty() bool {
	return m.LoginDefaultOrg == nil &&
		m.TriggerIntrospectionProjections == nil &&
		m.LegacyIntrospection == nil &&
		m.UserSchema == nil &&
		m.TokenExchange == nil &&
		m.Actions == nil &&
		// nil check to allow unset improvements
		m.ImprovedPerformance == nil &&
		m.WebKey == nil &&
		m.DebugOIDCParentError == nil &&
		m.OIDCSingleV1SessionTermination == nil &&
		m.DisableUserTokenEvent == nil &&
		m.EnableBackChannelLogout == nil
}

func (c *Commands) SetInstanceFeatures(ctx context.Context, f *InstanceFeatures) (*domain.ObjectDetails, error) {
	if f.isEmpty() {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Vigh1", "Errors.NoChangesFound")
	}
	wm := NewInstanceFeaturesWriteModel(authz.GetInstance(ctx).InstanceID())
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	if err := c.setupWebKeyFeature(ctx, wm, f); err != nil {
		return nil, err
	}
	commands := wm.setCommands(ctx, f)
	if len(commands) == 0 {
		return writeModelToObjectDetails(wm.WriteModel), nil
	}
	events, err := c.eventstore.Push(ctx, commands...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}

func prepareSetFeatures(instanceID string, f *InstanceFeatures) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, _ preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			wm := NewInstanceFeaturesWriteModel(instanceID)
			return wm.setCommands(ctx, f), nil
		}, nil
	}
}

// setupWebKeyFeature generates the initial web keys for the instance,
// if the feature is enabled in the request and the feature wasn't enabled already in the writeModel.
// [Commands.GenerateInitialWebKeys] checks if keys already exist and does nothing if that's the case.
// The default config of a RSA key with 2048 and the SHA256 hasher is assumed.
// Users can customize this after using the webkey/v3 API.
func (c *Commands) setupWebKeyFeature(ctx context.Context, wm *InstanceFeaturesWriteModel, f *InstanceFeatures) error {
	if !gu.Value(f.WebKey) || gu.Value(wm.WebKey) {
		return nil
	}
	return c.GenerateInitialWebKeys(ctx, &crypto.WebKeyRSAConfig{
		Bits:   crypto.RSABits2048,
		Hasher: crypto.RSAHasherSHA256,
	})
}

func (c *Commands) ResetInstanceFeatures(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	wm := NewInstanceFeaturesWriteModel(instanceID)
	if err := c.eventstore.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	if wm.isEmpty() {
		return writeModelToObjectDetails(wm.WriteModel), nil
	}
	aggregate := feature_v2.NewAggregate(instanceID, instanceID)
	events, err := c.eventstore.Push(ctx, feature_v2.NewResetEvent(ctx, aggregate, feature_v2.InstanceResetEventType))
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(events), nil
}
