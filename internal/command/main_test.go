package command

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/verifier"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	action_repo "github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	iam_repo "github.com/zitadel/zitadel/internal/repository/instance"
	key_repo "github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	proj_repo "github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type expect func(mockRepository *mock.MockRepository)

func eventstoreExpect(t *testing.T, expects ...expect) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e(m)
	}
	es := eventstore.NewEventstore(eventstore.TestConfig(m))
	iam_repo.RegisterEventMappers(es)
	org.RegisterEventMappers(es)
	usr_repo.RegisterEventMappers(es)
	proj_repo.RegisterEventMappers(es)
	usergrant.RegisterEventMappers(es)
	key_repo.RegisterEventMappers(es)
	action_repo.RegisterEventMappers(es)
	session.RegisterEventMappers(es)
	idpintent.RegisterEventMappers(es)
	authrequest.RegisterEventMappers(es)
	oidcsession.RegisterEventMappers(es)
	return es
}

func expectEventstore(expects ...expect) func(*testing.T) *eventstore.Eventstore {
	return func(t *testing.T) *eventstore.Eventstore {
		return eventstoreExpect(t, expects...)
	}
}

func eventPusherToEvents(eventsPushes ...eventstore.Command) []*repository.Event {
	events := make([]*repository.Event, len(eventsPushes))
	for i, event := range eventsPushes {
		data, err := eventstore.EventData(event)
		if err != nil {
			return nil
		}
		events[i] = &repository.Event{
			AggregateID:   event.Aggregate().ID,
			AggregateType: repository.AggregateType(event.Aggregate().Type),
			ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
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
			for _, f := range filter {
				if f.Field == repository.FieldAggregateType {
					if event.AggregateType != f.Value {
						continue
					}
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

func expectPush(events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPush(events, uniqueConstraints...)
	}
}

func expectPushFailed(err error, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPushFailed(err, events, uniqueConstraints...)
	}
}

func expectRandomPush(events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectRandomPush(events, uniqueConstraints...)
	}
}

func expectRandomPushFailed(err error, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) expect {
	return func(m *mock.MockRepository) {
		m.ExpectRandomPushFailed(err, events, uniqueConstraints...)
	}
}

func expectFilter(events ...*repository.Event) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}
func expectFilterError(err error) expect {
	return func(m *mock.MockRepository) {
		m.ExpectFilterEventsError(err)
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

func eventFromEventPusher(event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:                            "",
		Sequence:                      0,
		PreviousAggregateSequence:     0,
		PreviousAggregateTypeSequence: 0,
		CreationDate:                  time.Time{},
		Type:                          repository.EventType(event.Type()),
		Data:                          data,
		EditorService:                 event.EditorService(),
		EditorUser:                    event.EditorUser(),
		Version:                       repository.Version(event.Aggregate().Version),
		AggregateID:                   event.Aggregate().ID,
		AggregateType:                 repository.AggregateType(event.Aggregate().Type),
		ResourceOwner:                 sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
	}
}

func eventFromEventPusherWithInstanceID(instanceID string, event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:                            "",
		Sequence:                      0,
		PreviousAggregateSequence:     0,
		PreviousAggregateTypeSequence: 0,
		CreationDate:                  time.Time{},
		Type:                          repository.EventType(event.Type()),
		Data:                          data,
		EditorService:                 event.EditorService(),
		EditorUser:                    event.EditorUser(),
		Version:                       repository.Version(event.Aggregate().Version),
		AggregateID:                   event.Aggregate().ID,
		AggregateType:                 repository.AggregateType(event.Aggregate().Type),
		ResourceOwner:                 sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		InstanceID:                    instanceID,
	}
}

func eventFromEventPusherWithCreationDateNow(event eventstore.Command) *repository.Event {
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

func uniqueConstraintsFromEventConstraintWithInstanceID(instanceID string, constraint *eventstore.EventUniqueConstraint) *repository.UniqueConstraint {
	return &repository.UniqueConstraint{
		InstanceID:   instanceID,
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

type mockInstance struct{}

func (m *mockInstance) InstanceID() string {
	return "INSTANCE"
}

func (m *mockInstance) ProjectID() string {
	return "projectID"
}

func (m *mockInstance) ConsoleClientID() string {
	return "consoleID"
}

func (m *mockInstance) ConsoleApplicationID() string {
	return "consoleApplicationID"
}

func (m *mockInstance) DefaultLanguage() language.Tag {
	return language.English
}

func (m *mockInstance) DefaultOrganisationID() string {
	return "defaultOrgID"
}

func (m *mockInstance) RequestedDomain() string {
	return "zitadel.cloud"
}

func (m *mockInstance) RequestedHost() string {
	return "zitadel.cloud:443"
}

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func newMockPermissionCheckAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return nil
	}
}

func newMockPermissionCheckNotAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return errors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")
	}
}

type plainHasher struct {
	x string // arbitrary info that triggers update when different from encoding
}

func (h plainHasher) Hash(password string) (string, error) {
	return strings.Join([]string{"", "plain", h.x, password}, "$"), nil
}

func (h plainHasher) Verify(encoded, password string) (verifier.Result, error) {
	nodes := strings.Split(encoded, "$")
	if len(nodes) != 4 || nodes[1] != "plain" {
		return verifier.Skip, nil
	}
	if nodes[3] != password {
		return verifier.Fail, nil
	}
	if nodes[2] != h.x {
		return verifier.NeedUpdate, nil
	}
	return verifier.OK, nil
}

// mockPasswordHasher creates a swapper for plain (cleartext) password used in tests.
// x can be set to arbitrary info which triggers updates when different from the
// setting in the encoded hashes. (normally cost parameters)
//
// With `x` set to "foo", the following encoded string would be produced by Hash:
// $plain$foo$password
func mockPasswordHasher(x string) *crypto.PasswordHasher {
	return &crypto.PasswordHasher{
		Swapper:  passwap.NewSwapper(plainHasher{x: x}),
		Prefixes: []string{"$plain$"},
	}
}
