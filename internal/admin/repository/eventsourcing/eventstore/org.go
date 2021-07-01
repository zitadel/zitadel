package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
)

type OrgRepo struct {
	Eventstore v1.Eventstore

	View *admin_view.View

	SearchLimit    uint64
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.OrgView, error) {
	org, viewErr := repo.View.OrgByID(id)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		org = new(model.OrgView)
	}

	sequence := org.Sequence
	currentSequence, err := repo.View.GetLatestOrgSequence()
	if err == nil {
		sequence = currentSequence.CurrentSequence
	}

	events, esErr := repo.getOrgEvents(ctx, id, sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-Lsoj7", "Errors.Org.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return model.OrgToModel(org), nil
	}
	orgCopy := *org
	for _, event := range events {
		if err := orgCopy.AppendEvent(event); err != nil {
			return model.OrgToModel(&orgCopy), nil
		}
	}
	return model.OrgToModel(&orgCopy), nil
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	err := query.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-LXo9w").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest iam sequence")
	orgs, count, err := repo.View.SearchOrgs(query)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      query.Offset,
		Limit:       query.Limit,
		TotalResult: count,
		Result:      model.OrgsToModel(orgs),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	var found bool
	err = es_sdk.Filter(ctx, repo.Eventstore.FilterEvents, isUniqueValidation(&found), view.OrgNameUniqueQuery(name))
	if (err != nil && !errors.IsNotFound(err)) || found {
		return false, err
	}

	err = es_sdk.Filter(ctx, repo.Eventstore.FilterEvents, isUniqueValidation(&found), view.OrgDomainUniqueQuery(domain))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return !found, nil
}

func (repo *OrgRepo) GetOrgIAMPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(id)
	if errors.IsNotFound(err) {
		return repo.GetDefaultOrgIAMPolicy(ctx)
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepo) GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy.Default = true
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepo) getOrgEvents(ctx context.Context, orgID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.OrgByIDQuery(orgID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}

func isUniqueValidation(unique *bool) func(events ...*models.Event) error {
	return func(events ...*models.Event) error {
		if len(events) == 0 {
			return nil
		}
		*unique = *unique || events[0].Type == org_es_model.OrgDomainReserved || events[0].Type == org_es_model.OrgNameReserved

		return nil
	}
}
