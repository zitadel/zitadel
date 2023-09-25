package mock

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
	quota_repo "github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/repository/session"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
)

type Expecter interface {
	Expect(mockRepository *mock.MockRepository)
}

// ExpectFunc implements the Expecter interface
type ExpectFunc func(mockRepository *mock.MockRepository)

func (e ExpectFunc) Expect(mockRepository *mock.MockRepository) {
	e(mockRepository)
}

func EventstoreExpect(t *testing.T, expects ...Expecter) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e.Expect(m)
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
	quota_repo.RegisterEventMappers(es)
	return es
}

func ExpectEventstore(expects ...Expecter) func(*testing.T) *eventstore.Eventstore {
	return func(t *testing.T) *eventstore.Eventstore {
		return EventstoreExpect(t, expects...)
	}
}

func EventPusherToEvents(eventsPushes ...eventstore.Command) []*repository.Event {
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

type TestRepo struct {
	events            []*repository.Event
	uniqueConstraints []*repository.UniqueConstraint
	sequence          uint64
	err               error
	t                 *testing.T
}

func (repo *TestRepo) Health(ctx context.Context) error {
	return nil
}

func (repo *TestRepo) Push(ctx context.Context, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) error {
	repo.events = append(repo.events, events...)
	repo.uniqueConstraints = append(repo.uniqueConstraints, uniqueConstraints...)
	return nil
}

func (repo *TestRepo) Filter(ctx context.Context, searchQuery *repository.SearchQuery) ([]*repository.Event, error) {
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

func (repo *TestRepo) LatestSequence(ctx context.Context, queryFactory *repository.SearchQuery) (uint64, error) {
	if repo.err != nil {
		return 0, repo.err
	}
	return repo.sequence, nil
}

func ExpectPush(events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectPush(events, uniqueConstraints...)
	})
}

func ExpectPushFailed(err error, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectPushFailed(err, events, uniqueConstraints...)
	})
}

func ExpectRandomPush(events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectRandomPush(events, uniqueConstraints...)
	})
}

func ExpectRandomPushFailed(err error, events []*repository.Event, uniqueConstraints ...*repository.UniqueConstraint) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectRandomPushFailed(err, events, uniqueConstraints...)
	})
}

func ExpectFilter(events ...*repository.Event) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	})
}
func ExpectFilterError(err error) Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectFilterEventsError(err)
	})
}

func ExpectFilterOrgDomainNotFound() Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectFilterNoEventsNoError()
	})
}

func ExpectFilterOrgMemberNotFound() Expecter {
	return ExpectFunc(func(m *mock.MockRepository) {
		m.ExpectFilterNoEventsNoError()
	})
}

func EventFromEventPusher(event eventstore.Command) *repository.Event {
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

func EventFromEventPusherWithInstanceID(instanceID string, event eventstore.Command) *repository.Event {
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

func EventFromEventPusherWithCreationDateNow(event eventstore.Command) *repository.Event {
	e := EventFromEventPusher(event)
	e.CreationDate = time.Now()
	return e
}

func UniqueConstraintsFromEventConstraint(constraint *eventstore.EventUniqueConstraint) *repository.UniqueConstraint {
	return &repository.UniqueConstraint{
		UniqueType:   constraint.UniqueType,
		UniqueField:  constraint.UniqueField,
		ErrorMessage: constraint.ErrorMessage,
		Action:       repository.UniqueConstraintAction(constraint.Action)}
}

func UniqueConstraintsFromEventConstraintWithInstanceID(instanceID string, constraint *eventstore.EventUniqueConstraint) *repository.UniqueConstraint {
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

type MockInstance struct{}

func (m *MockInstance) InstanceID() string {
	return "INSTANCE"
}

func (m *MockInstance) ProjectID() string {
	return "projectID"
}

func (m *MockInstance) ConsoleClientID() string {
	return "consoleID"
}

func (m *MockInstance) ConsoleApplicationID() string {
	return "consoleApplicationID"
}

func (m *MockInstance) DefaultLanguage() language.Tag {
	return language.English
}

func (m *MockInstance) DefaultOrganisationID() string {
	return "defaultOrgID"
}

func (m *MockInstance) RequestedDomain() string {
	return "zitadel.cloud"
}

func (m *MockInstance) RequestedHost() string {
	return "zitadel.cloud:443"
}

func (m *MockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func NewMockPermissionCheckAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return nil
	}
}

func NewMockPermissionCheckNotAllowed() domain.PermissionCheck {
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

// MockPasswordHasher creates a swapper for plain (cleartext) password used in tests.
// x can be set to arbitrary info which triggers updates when different from the
// setting in the encoded hashes. (normally cost parameters)
//
// With `x` set to "foo", the following encoded string would be produced by Hash:
// $plain$foo$password
func MockPasswordHasher(x string) *crypto.PasswordHasher {
	return &crypto.PasswordHasher{
		Swapper:  passwap.NewSwapper(plainHasher{x: x}),
		Prefixes: []string{"$plain$"},
	}
}
