package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/repository/mock"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	key_repo "github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	proj_repo "github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/repository/usergrant"
)

type expect func(mockRepository *mock.MockRepository)

func eventstoreExpect(t *testing.T, expects ...expect) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e(m)
	}
	es := eventstore.NewEventstore(m)
	iam_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	proj_repo.RegisterEventMappers(es)
	usergrant.RegisterEventMappers(es)
	key_repo.RegisterEventMappers(es)
	return es
}

func eventPusherToEvents(eventsPushes ...eventstore.EventPusher) []*repository.Event {
	events := make([]*repository.Event, len(eventsPushes))
	for i, event := range eventsPushes {
		data, err := eventstore.EventData(event)
		if err != nil {
			return nil
		}
		events[i] = &repository.Event{
			AggregateID:   event.Aggregate().ID,
			AggregateType: repository.AggregateType(event.Aggregate().Typ),
			ResourceOwner: event.Aggregate().ResourceOwner,
			EditorService: event.EditorService(),
			EditorUser:    event.EditorUser(),
			Type:          repository.EventType(event.Type()),
			Version:       repository.Version(event.Aggregate().Version),
			Data:          data,
		}
	}
	return events
}

type testRepo struct {
	events            []*repository.Event
	uniqueConstraints []*repository.UniqueConstraint
	sequence          uint64
	err               error
	t                 *testing.T
}

func (repo *testRepo) Health(ctx context.Context) error {
	return nil
}

func (repo *testRepo) Push(ctx context.Context, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) error {
	repo.events = append(repo.events, events...)
	repo.uniqueConstraints = append(repo.uniqueConstraints, uniqueConstraints...)
	return nil
}

func (repo *testRepo) Filter(ctx context.Context, searchQuery *repository.SearchQuery) ([]*repository.Event, error) {
	events := make([]*repository.Event, 0, len(repo.events))
	for _, event := range repo.events {
		for _, filter := range searchQuery.Filters {
			if filter.Field == repository.FieldAggregateType {
				if event.AggregateType != filter.Value {
					continue
				}
			}
		}
		events = append(events, event)
	}
	return repo.events, nil
}

func filterAggregateType(aggregateType string) {

}

func (repo *testRepo) LatestSequence(ctx context.Context, queryFactory *repository.SearchQuery) (uint64, error) {
	if repo.err != nil {
		return 0, repo.err
	}
	return repo.sequence, nil
}

func expectPush(events []*repository.Event, assets []*repository.Asset, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPush(events, assets, uniqueConstraints...)
	}
}

func expectPushFailed(err error, events []*repository.Event, assets []*repository.Asset, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPushFailed(err, events, assets, uniqueConstraints...)
	}
}

func expectFilter(events ...*repository.Event) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}

func expectFilterOrgDomainNotFound() expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterNoEventsNoError()
	}
}

func expectFilterOrgMemberNotFound() expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterNoEventsNoError()
	}
}

func eventFromEventPusher(event eventstore.EventPusher) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:               "",
		Sequence:         0,
		PreviousSequence: 0,
		CreationDate:     time.Time{},
		Type:             repository.EventType(event.Type()),
		Data:             data,
		EditorService:    event.EditorService(),
		EditorUser:       event.EditorUser(),
		Version:          repository.Version(event.Aggregate().Version),
		AggregateID:      event.Aggregate().ID,
		AggregateType:    repository.AggregateType(event.Aggregate().Typ),
		ResourceOwner:    event.Aggregate().ResourceOwner,
	}
}

func eventFromEventPusherWithCreationDateNow(event eventstore.EventPusher) *repository.Event {
	e := eventFromEventPusher(event)
	e.CreationDate = time.Now()
	return e
}

func uniqueConstraintsFromEventConstraint(constraint *eventstore.EventUniqueConstraint) *repository.UniqueConstraint {
	return &repository.UniqueConstraint{
		UniqueType:   constraint.UniqueType,
		UniqueField:  constraint.UniqueField,
		ErrorMessage: constraint.ErrorMessage,
		Action:       repository.UniqueConstraintAction(constraint.Action)}
}

func GetMockSecretGenerator(t *testing.T) crypto.Generator {
	ctrl := gomock.NewController(t)
	alg := crypto.CreateMockEncryptionAlg(ctrl)
	generator := crypto.NewMockGenerator(ctrl)
	generator.EXPECT().Length().Return(uint(1)).AnyTimes()
	generator.EXPECT().Runes().Return([]rune("aa")).AnyTimes()
	generator.EXPECT().Alg().Return(alg).AnyTimes()
	generator.EXPECT().Expiry().Return(time.Hour * 1).AnyTimes()

	return generator
}

func GetMockVerifier(t *testing.T, features ...string) *authz.TokenVerifier {
	return authz.Start(&testVerifier{
		features: features,
	})
}

type testVerifier struct {
	features []string
}

func (v *testVerifier) VerifyAccessToken(ctx context.Context, token, clientID string) (string, string, string, string, error) {
	return "userID", "agentID", "de", "orgID", nil
}
func (v *testVerifier) SearchMyMemberships(ctx context.Context) ([]*authz.Membership, error) {
	return nil, nil
}

func (v *testVerifier) ProjectIDAndOriginsByClientID(ctx context.Context, clientID string) (string, []string, error) {
	return "", nil, nil
}

func (v *testVerifier) ExistsOrg(ctx context.Context, orgID string) error {
	return nil
}

func (v *testVerifier) VerifierClientID(ctx context.Context, appName string) (string, error) {
	return "clientID", nil
}

func (v *testVerifier) CheckOrgFeatures(ctx context.Context, orgID string, requiredFeatures ...string) error {
	for _, feature := range requiredFeatures {
		hasFeature := false
		for _, f := range v.features {
			if f == feature {
				hasFeature = true
				break
			}
		}
		if !hasFeature {
			return errors.ThrowPermissionDenied(nil, "id", "missing feature")
		}
	}
	return nil
}
