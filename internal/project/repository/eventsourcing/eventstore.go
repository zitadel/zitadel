package eventsourcing

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	key_model "github.com/caos/zitadel/internal/key/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

const (
	projectOwnerRole       = "PROJECT_OWNER"
	projectOwnerGlobalRole = "PROJECT_OWNER_GLOBAL"
)

type ProjectEventstore struct {
	es_int.Eventstore
	projectCache  *ProjectCache
	passwordAlg   crypto.HashAlgorithm
	pwGenerator   crypto.Generator
	idGenerator   id.Generator
	ClientKeySize int
}

type ProjectConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartProject(conf ProjectConfig, systemDefaults sd.SystemDefaults) (*ProjectEventstore, error) {
	projectCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	passwordAlg := crypto.NewBCrypt(systemDefaults.SecretGenerators.PasswordSaltCost)
	pwGenerator := crypto.NewHashGenerator(systemDefaults.SecretGenerators.ClientSecretGenerator, passwordAlg)
	return &ProjectEventstore{
		Eventstore:    conf.Eventstore,
		projectCache:  projectCache,
		passwordAlg:   passwordAlg,
		pwGenerator:   pwGenerator,
		idGenerator:   id.SonyFlakeGenerator,
		ClientKeySize: int(systemDefaults.SecretGenerators.ClientKeySize),
	}, nil
}

func (es *ProjectEventstore) ProjectByID(ctx context.Context, id string) (*proj_model.Project, error) {
	project := es.projectCache.getProject(id)

	query, err := ProjectByIDQuery(project.AggregateID, project.Sequence)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Filter(ctx, es.FilterEvents, project.AppendEvents, query)
	if err != nil && !(caos_errs.IsNotFound(err) && project.Sequence != 0) {
		return nil, err
	}
	if project.State == int32(proj_model.ProjectStateRemoved) {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-dG8ie", "Errors.Project.NotFound")
	}
	es.projectCache.cacheProject(project)
	return model.ProjectToModel(project), nil
}

func (es *ProjectEventstore) ProjectEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := ProjectByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
}

func (es *ProjectEventstore) ProjectChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ProjectChanges, error) {
	query := ChangesQuery(id, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*proj_model.ProjectChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &proj_model.ProjectChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			var data interface{}
			if strings.Contains(change.EventType, "application") {
				data = new(model.Application)
			} else {
				data = new(model.Project)
			}
			err = json.Unmarshal(event.Data, data)
			logging.Log("EVENT-NCkpN").OnError(err).Debug("unable to unmarshal data")
			change.Data = data
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ProjectChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func ChangesQuery(projectID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.ProjectAggregate)
	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(projectID).
		SetLimit(limit)
	return query
}

func (es *ProjectEventstore) AddApplication(ctx context.Context, app *proj_model.Application) (*proj_model.Application, error) {
	if app == nil || !app.IsValid(true) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9eidw", "Errors.Project.App.Invalid")
	}
	existingProject, err := es.ProjectByID(ctx, app.AggregateID)
	if err != nil {
		return nil, err
	}
	app.AppID, err = es.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	var stringPw string
	if app.OIDCConfig != nil {
		app.OIDCConfig.AppID = app.AppID
		err := app.OIDCConfig.GenerateNewClientID(es.idGenerator, existingProject)
		if err != nil {
			return nil, err
		}
		stringPw, err = app.OIDCConfig.GenerateClientSecretIfNeeded(es.pwGenerator)
		if err != nil {
			return nil, err
		}
	}
	if app.APIConfig != nil {
		app.APIConfig.AppID = app.AppID
		err := app.APIConfig.GenerateNewClientID(es.idGenerator, existingProject)
		if err != nil {
			return nil, err
		}
		stringPw, err = app.APIConfig.GenerateClientSecretIfNeeded(es.pwGenerator)
		if err != nil {
			return nil, err
		}
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoApp := model.AppFromModel(app)

	addAggregate := ApplicationAddedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoApp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		converted := model.AppToModel(a)
		if converted.OIDCConfig != nil {
			converted.OIDCConfig.ClientSecretString = stringPw
			converted.OIDCConfig.FillCompliance()
		}
		if converted.APIConfig != nil {
			converted.APIConfig.ClientSecretString = stringPw
		}
		return converted, nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-GvPct", "Errors.Internal")
}

func (es *ProjectEventstore) ApplicationChanges(ctx context.Context, projectID string, appID string, lastSequence uint64, limit uint64, sortAscending bool) (*proj_model.ApplicationChanges, error) {
	query := ChangesQuery(projectID, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(ctx, query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-sw6Ku", "Errors.Internal")
	}
	if len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-9IHLP", "Errors.Changes.NotFound")
	}

	result := make([]*proj_model.ApplicationChange, 0)
	for _, event := range events {
		if !strings.Contains(event.Type.String(), "application") || event.Data == nil {
			continue
		}

		app := new(model.Application)
		err := json.Unmarshal(event.Data, app)
		logging.Log("EVENT-GIiKD").OnError(err).Debug("unable to unmarshal data")
		if app.AppID != appID {
			continue
		}

		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-MJzeN").OnError(err).Debug("unable to parse timestamp")

		result = append(result, &proj_model.ApplicationChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
			Data:       app,
		})
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &proj_model.ApplicationChanges{
		Changes:      result,
		LastSequence: lastSequence,
	}, nil
}

func (es *ProjectEventstore) ChangeAPIConfig(ctx context.Context, config *proj_model.APIConfig) (*proj_model.APIConfig, error) {
	if config == nil || !config.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-SDg54", "Errors.Project.APIConfigInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(config.AppID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Rgu63", "Errors.Project.AppNotExisting")
	}
	if app.Type != proj_model.AppTypeAPI {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-RHj63", "Errors.Project.AppIsNotAPI")
	}
	repoProject := model.ProjectFromModel(existingProject)
	repoConfig := model.APIConfigFromModel(config)

	projectAggregate := APIConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)
	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		return model.APIConfigToModel(a.APIConfig), nil
	}
	return nil, caos_errs.ThrowInternal(nil, "EVENT-aebn5", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeOIDCConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.OIDCConfig, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-7ue34", "Errors.Project.App.OIDCConfigInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(appID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-9odi4", "Errors.Project.App.NotExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-dile4", "Errors.Project.App.IsNotOIDC")
	}
	if app.OIDCConfig.AuthMethodType == proj_model.OIDCAuthMethodTypeNone || app.OIDCConfig.AuthMethodType == proj_model.OIDCAuthMethodTypePrivateKeyJWT {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-GDrg2", "Errors.Project.OIDCAuthMethodNoSecret")
	}
	repoProject := model.ProjectFromModel(existingProject)

	stringPw, err := app.OIDCConfig.GenerateNewClientSecret(es.pwGenerator)
	if err != nil {
		return nil, err
	}

	projectAggregate := OIDCConfigSecretChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, appID, app.OIDCConfig.ClientSecret)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		config := model.OIDCConfigToModel(a.OIDCConfig)
		config.ClientSecretString = stringPw
		return config, nil
	}

	return nil, caos_errs.ThrowInternal(nil, "EVENT-dk87s", "Errors.Internal")
}

func (es *ProjectEventstore) ChangeAPIConfigSecret(ctx context.Context, projectID, appID string) (*proj_model.APIConfig, error) {
	if appID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-sdfb3", "Errors.Project.APIConfigInvalid")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(appID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-ADbg3", "Errors.Project.AppNotExisting")
	}
	if app.Type != proj_model.AppTypeAPI {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Ntwqw", "Errors.Project.AppIsNotAPI")
	}
	if app.APIConfig.AuthMethodType != proj_model.APIAuthMethodTypeBasic {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-HW4tw", "Errors.Project.APIAuthMethodNoSecret")
	}
	repoProject := model.ProjectFromModel(existingProject)

	stringPw, err := app.APIConfig.GenerateNewClientSecret(es.pwGenerator)
	if err != nil {
		return nil, err
	}

	projectAggregate := APIConfigSecretChangedAggregate(es.Eventstore.AggregateCreator(), repoProject, appID, app.APIConfig.ClientSecret)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, projectAggregate)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	if _, a := model.GetApplication(repoProject.Applications, app.AppID); a != nil {
		config := model.APIConfigToModel(a.APIConfig)
		config.ClientSecretString = stringPw
		return config, nil
	}

	return nil, caos_errs.ThrowInternal(nil, "EVENT-HBfju", "Errors.Internal")
}

func (es *ProjectEventstore) VerifyOIDCClientSecret(ctx context.Context, projectID, appID string, secret string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if appID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-H3RT2", "Errors.Project.RequiredFieldsMissing")
	}
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(appID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D6hba", "Errors.Project.AppNoExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-huywq", "Errors.Project.App.IsNotOIDC")
	}

	ctx, spanHash := tracing.NewSpan(ctx)
	err = crypto.CompareHash(app.OIDCConfig.ClientSecret, []byte(secret), es.passwordAlg)
	spanHash.EndWithError(err)
	if err == nil {
		err = es.setOIDCClientSecretCheckResult(ctx, existingProject, app.AppID, OIDCClientSecretCheckSucceededAggregate)
		logging.Log("EVENT-AE1vf").OnError(err).Warn("could not push event OIDCClientSecretCheckSucceeded")
		return nil
	}
	err = es.setOIDCClientSecretCheckResult(ctx, existingProject, app.AppID, OIDCClientSecretCheckFailedAggregate)
	logging.Log("EVENT-GD1gh").OnError(err).Warn("could not push event OIDCClientSecretCheckFailed")
	return caos_errs.ThrowInvalidArgument(nil, "EVENT-wg24q", "Errors.Project.OIDCSecretInvalid")
}

func (es *ProjectEventstore) setOIDCClientSecretCheckResult(ctx context.Context, project *proj_model.Project, appID string, check func(*es_models.AggregateCreator, *model.Project, string) es_sdk.AggregateFunc) error {
	repoProject := model.ProjectFromModel(project)
	agg := check(es.AggregateCreator(), repoProject, appID)
	err := es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}

func (es *ProjectEventstore) AddClientKey(ctx context.Context, key *proj_model.ClientKey) (*proj_model.ClientKey, error) {
	existingProject, err := es.ProjectByID(ctx, key.AggregateID)
	if err != nil {
		return nil, err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(key.ApplicationID); app == nil {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dbf32", "Errors.Project.AppNoExisting")
	}
	if (app.OIDCConfig == nil || app.OIDCConfig != nil && app.OIDCConfig.AuthMethodType != proj_model.OIDCAuthMethodTypePrivateKeyJWT) &&
		(app.APIConfig == nil || app.APIConfig != nil && app.APIConfig.AuthMethodType != proj_model.APIAuthMethodTypePrivateKeyJWT) {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "EVENT-Dff54", "Errors.Project.AuthMethodNoPrivateKeyJWT")
	}
	if app.OIDCConfig != nil {
		key.ClientID = app.OIDCConfig.ClientID
	}
	if app.APIConfig != nil {
		key.ClientID = app.APIConfig.ClientID
	}
	key.KeyID, err = es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	if key.ExpirationDate.IsZero() {
		key.ExpirationDate, err = key_model.DefaultExpiration()
		if err != nil {
			logging.Log("EVENT-Adgf2").WithError(err).Warn("unable to set default date")
			return nil, errors.ThrowInternal(err, "EVENT-j68fg", "Errors.Internal")
		}
	}
	if key.ExpirationDate.Before(time.Now()) {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-C6YV5", "Errors.MachineKey.ExpireBeforeNow")
	}

	repoProject := model.ProjectFromModel(existingProject)
	repoKey := model.ClientKeyFromModel(key)
	err = repoKey.GenerateClientKeyPair(es.ClientKeySize)
	if err != nil {
		return nil, err
	}
	agg := OIDCApplicationKeyAddedAggregate(es.AggregateCreator(), repoProject, repoKey)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return nil, err
	}
	es.projectCache.cacheProject(repoProject)

	return model.ClientKeyToModel(repoKey), nil
}

func (es *ProjectEventstore) RemoveApplicationKey(ctx context.Context, projectID, applicationID, keyID string) error {
	existingProject, err := es.ProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	var app *proj_model.Application
	if _, app = existingProject.GetApp(applicationID); app == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-ADfzz", "Errors.Project.AppNotExisting")
	}
	if app.Type != proj_model.AppTypeOIDC {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-ADffh", "Errors.Project.AppIsNotOIDC")
	}
	if _, key := app.GetKey(keyID); key == nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-D2Sff", "Errors.Project.AppKeyNotExisting")
	}
	repoProject := model.ProjectFromModel(existingProject)
	agg := OIDCApplicationKeyRemovedAggregate(es.AggregateCreator(), repoProject, keyID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoProject.AppendEvents, agg)
	if err != nil {
		return err
	}
	es.projectCache.cacheProject(repoProject)
	return nil
}
