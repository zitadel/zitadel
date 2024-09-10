package readmodel

// import (
// 	"context"

// 	"github.com/zitadel/zitadel/internal/feature"
// 	"github.com/zitadel/zitadel/internal/v2/eventstore"
// 	"github.com/zitadel/zitadel/internal/v2/projection"
// )

// func NewSystemFeatureCache(ctx context.Context, eventStore eventstore.EventStore) {
// 	cache := NewCachedReadModel[*SystemFeatures, *SystemFeatures](ctx, eventStore, NewSystemFeatures().Reduce)
// }

// type SystemFeatures struct {
// 	LoginDefaultOrg                 *projection.Feature[bool]
// 	TriggerIntrospectionProjections *projection.Feature[bool]
// 	LegacyIntrospection             *projection.Feature[bool]
// 	UserSchema                      *projection.Feature[bool]
// 	TokenExchange                   *projection.Feature[bool]
// 	Actions                         *projection.Feature[bool]
// 	ImprovedPerformance             *projection.Feature[[]ImprovedPerformanceType]
// 	OIDCSingleV1SessionTermination  *projection.Feature[bool]
// }

// func NewSystemFeatures() *SystemFeatures {
// 	return &SystemFeatures{
// 		LoginDefaultOrg:                 projection.NewFeature[bool](feature.LevelSystem, feature.KeyLoginDefaultOrg),
// 		TriggerIntrospectionProjections: projection.NewFeature[bool](feature.LevelSystem, feature.KeyTriggerIntrospectionProjections),
// 		LegacyIntrospection:             projection.NewFeature[bool](feature.LevelSystem, feature.KeyLegacyIntrospection),
// 		UserSchema:                      projection.NewFeature[bool](feature.LevelSystem, feature.KeyUserSchema),
// 		TokenExchange:                   projection.NewFeature[bool](feature.LevelSystem, feature.KeyTokenExchange),
// 		Actions:                         projection.NewFeature[bool](feature.LevelSystem, feature.KeyActions),
// 		ImprovedPerformance:             projection.NewFeature[[]ImprovedPerformanceType](feature.LevelSystem, feature.KeyImprovedPerformance),
// 		OIDCSingleV1SessionTermination:  projection.NewFeature[bool](feature.LevelSystem, feature.KeyOIDCSingleV1SessionTermination),
// 	}
// }

// func (f *SystemFeatures) Reduce(events ...*eventstore.StorageEvent) error {
// 	for _, event := range events {
// 		for _, reduce := range []eventstore.Reduce{
// 			f.LoginDefaultOrg.Reduce,
// 			f.TriggerIntrospectionProjections.Reduce,
// 			f.LegacyIntrospection.Reduce,
// 			f.UserSchema.Reduce,
// 			f.TokenExchange.Reduce,
// 			f.Actions.Reduce,
// 			f.ImprovedPerformance.Reduce,
// 			f.OIDCSingleV1SessionTermination.Reduce,
// 		} {
// 			if err := reduce(event); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// type InstanceFeatures struct {
// 	model

// 	parent *SystemFeatures

// 	LoginDefaultOrg                 *projection.Feature[bool]
// 	TriggerIntrospectionProjections *projection.Feature[bool]
// 	LegacyIntrospection             *projection.Feature[bool]
// 	UserSchema                      *projection.Feature[bool]
// 	TokenExchange                   *projection.Feature[bool]
// 	Actions                         *projection.Feature[bool]
// 	ImprovedPerformance             *projection.Feature[[]feature.ImprovedPerformanceType]
// 	WebKey                          *projection.Feature[bool]
// 	DebugOIDCParentError            *projection.Feature[bool]
// 	OIDCSingleV1SessionTermination  *projection.Feature[bool]
// 	InMemoryProjections             *projection.Feature[bool]
// }

// func (f *InstanceFeatures) Reduce(events ...*eventstore.StorageEvent) error {
// 	for _, event := range events {
// 		for _, reduce := range []eventstore.Reduce{
// 			f.LoginDefaultOrg.Reduce,
// 			f.TriggerIntrospectionProjections.Reduce,
// 			f.LegacyIntrospection.Reduce,
// 			f.UserSchema.Reduce,
// 			f.TokenExchange.Reduce,
// 			f.Actions.Reduce,
// 			f.ImprovedPerformance.Reduce,
// 			f.WebKey.Reduce,
// 			f.DebugOIDCParentError.Reduce,
// 			f.OIDCSingleV1SessionTermination.Reduce,
// 			f.InMemoryProjections.Reduce,
// 		} {
// 			if err := reduce(event); err != nil {
// 				return err
// 			}
// 		}
// 		f.model.latestPosition = event.Position
// 	}
// 	return f.parent.Reduce(events...)
// }

// type ImprovedPerformanceType = feature.ImprovedPerformanceType
