package projection

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/dialect/postgres"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"

	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	SettingsRelationalProjectionTable = "zitadel.settings"

	SettingInstanceIDCol = "instance_id"
	SettingsOrgIDCol     = "org_id"
	SettingsIDCol        = "id"
	SettingsTypeCol      = "type"
	SettingsSettingsCol  = "settings"
)

type settingsRelationalProjection struct {
	settingsRepo domain.SettingsRepository
}

func newSettingsRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	client := postgres.PGxPool(config.Client.Pool)
	settingsRepo := repository.SettingsRepository(client)
	return handler.NewHandler(ctx, &config, &settingsRelationalProjection{
		settingsRepo: settingsRepo,
	})
}

func (*settingsRelationalProjection) Name() string {
	return SettingsRelationalProjectionTable
}

func (p *settingsRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				// {
				// 	Event:  org.LoginPolicyChangedEventType,
				// 	Reduce: p.reduceLoginPolicyChanged,
				// },
				// {
				// 	Event:  org.LoginPolicyMultiFactorAddedEventType,
				// 	Reduce: p.reduceMFAAdded,
				// },
				// {
				// 	Event:  org.LoginPolicyMultiFactorRemovedEventType,
				// 	Reduce: p.reduceMFARemoved,
				// },
				// {
				// 	Event:  org.LoginPolicyRemovedEventType,
				// 	Reduce: p.reduceLoginPolicyRemoved,
				// },
				// {
				// 	Event:  org.LoginPolicySecondFactorAddedEventType,
				// 	Reduce: p.reduceSecondFactorAdded,
				// },
				// {
				// 	Event:  org.LoginPolicySecondFactorRemovedEventType,
				// 	Reduce: p.reduceSecondFactorRemoved,
				// },
				// {
				// 	Event:  org.OrgRemovedEventType,
				// 	Reduce: p.reduceOwnerRemoved,
				// },
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.LoginPolicyAddedEventType,
					Reduce: p.reduceLoginPolicyAdded,
				},
				// {
				// 	Event:  instance.LoginPolicyChangedEventType,
				// 	Reduce: p.reduceLoginPolicyChanged,
				// },
				// {
				// 	Event:  instance.LoginPolicyMultiFactorAddedEventType,
				// 	Reduce: p.reduceMFAAdded,
				// },
				// {
				// 	Event:  instance.LoginPolicyMultiFactorRemovedEventType,
				// 	Reduce: p.reduceMFARemoved,
				// },
				// {
				// 	Event:  instance.LoginPolicySecondFactorAddedEventType,
				// 	Reduce: p.reduceSecondFactorAdded,
				// },
				// {
				// 	Event:  instance.LoginPolicySecondFactorRemovedEventType,
				// 	Reduce: p.reduceSecondFactorRemoved,
				// },
				// {
				// 	Event:  instance.InstanceRemovedEventType,
				// 	Reduce: reduceInstanceRemovedHelper(LoginPolicyInstanceIDCol),
				// },
			},
		},
	}
}

type loginSettings struct {
	policy.LoginPolicyAddedEvent
	IsDefault *bool `json:"isDefault,omitempty"`
}

func (s *settingsRelationalProjection) reduceLoginPolicyAdded(event eventstore.Event) (*handler.Statement, error) {
	var policyEvent policy.LoginPolicyAddedEvent
	var isDefault bool
	switch e := event.(type) {
	case *instance.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = true
	case *org.LoginPolicyAddedEvent:
		policyEvent = e.LoginPolicyAddedEvent
		isDefault = false
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YYPxS", "reduce.wrong.event.type %v", []eventstore.EventType{org.LoginPolicyAddedEventType, instance.LoginPolicyAddedEventType})
	}

	var orgId *string
	if policyEvent.Aggregate().ResourceOwner != policyEvent.Agg.InstanceID {
		orgId = &policyEvent.Aggregate().ResourceOwner
	}

	return handler.NewStatement(&policyEvent, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		loginSettings := loginSettings{
			LoginPolicyAddedEvent: policyEvent,
			IsDefault:             &isDefault,
		}
		settings, err := json.Marshal(loginSettings)
		if err != nil {
			return err
		}

		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5LOhE", "reduce.wrong.db.pool %T", ex)
		}
		settingsRepo := repository.SettingsRepository(v3_sql.SQLTx(tx))
		setting := domain.Setting{
			ID:         policyEvent.Aggregate().ID,
			InstanceID: policyEvent.Aggregate().InstanceID,
			OrgID:      orgId,
			Type:       domain.SettingTypeLogin,
			Settings:   settings,
		}
		// return settingsRepo.Create(ctx, &setting)
		err = settingsRepo.Create(ctx, &setting)
		return err
	}), nil
}
