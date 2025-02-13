package command

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/verifier"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type expect func(mockRepository *mock.MockRepository)

// Deprecated: use expectEventstore
func eventstoreExpect(t *testing.T, expects ...expect) *eventstore.Eventstore {
	m := mock.NewRepo(t)
	for _, e := range expects {
		e(m)
	}
	es := eventstore.NewEventstore(
		&eventstore.Config{
			Querier: m.MockQuerier,
			Pusher:  m.MockPusher,
		},
	)
	return es
}

// expectEventstore defines expectations for the Eventstore and is initialized within the scope of a (sub) test.
// This allows proper reporting of the test name, instead of reporting on the top-level
// of the Test function being run.
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
			AggregateType: event.Aggregate().Type,
			ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
			EditorUser:    event.Creator(),
			Typ:           event.Type(),
			Version:       event.Aggregate().Version,
			Data:          data,
			Constraints:   event.UniqueConstraints(),
		}
	}
	return events
}

func expectPush(commands ...eventstore.Command) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPush(commands, 0)
	}
}

func expectPushSlow(sleep time.Duration, commands ...eventstore.Command) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPush(commands, sleep)
	}
}

func expectPushFailed(err error, commands ...eventstore.Command) expect {
	return func(m *mock.MockRepository) {
		m.ExpectPushFailed(err, commands)
	}
}

func expectRandomPush(events []eventstore.Command) expect {
	return func(m *mock.MockRepository) {
		m.ExpectRandomPush(events)
	}
}

func expectRandomPushFailed(err error, events []eventstore.Command) expect {
	return func(m *mock.MockRepository) {
		m.ExpectRandomPushFailed(err, events)
	}
}

func expectFilter(events ...eventstore.Event) expect {
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
		InstanceID:    event.Aggregate().InstanceID,
		ID:            "",
		Seq:           0,
		CreationDate:  time.Time{},
		Typ:           event.Type(),
		Data:          data,
		EditorUser:    event.Creator(),
		Version:       event.Aggregate().Version,
		AggregateID:   event.Aggregate().ID,
		AggregateType: event.Aggregate().Type,
		ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		Constraints:   event.UniqueConstraints(),
	}
}

func eventFromEventPusherWithInstanceID(instanceID string, event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		ID:            "",
		Seq:           0,
		CreationDate:  time.Time{},
		Typ:           event.Type(),
		Data:          data,
		EditorUser:    event.Creator(),
		Version:       event.Aggregate().Version,
		AggregateID:   event.Aggregate().ID,
		AggregateType: event.Aggregate().Type,
		ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		InstanceID:    instanceID,
	}
}

func eventFromEventPusherWithCreationDateNow(event eventstore.Command) *repository.Event {
	e := eventFromEventPusher(event)
	e.CreationDate = time.Now()
	return e
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

func (m *mockInstance) Block() *bool {
	panic("shouldn't be called here")
}

func (m *mockInstance) AuditLogRetention() *time.Duration {
	panic("shouldn't be called here")
}

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
	return AllowedLanguage
}

func (m *mockInstance) DefaultOrganisationID() string {
	return "defaultOrgID"
}

func (m *mockInstance) SecurityPolicyAllowedOrigins() []string {
	return nil
}

func (m *mockInstance) EnableImpersonation() bool {
	return false
}

func (m *mockInstance) Features() feature.Features {
	return feature.Features{}
}

func newMockPermissionCheckAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return nil
	}
}

func newMockPermissionCheckNotAllowed() domain.PermissionCheck {
	return func(ctx context.Context, permission, orgID, resourceID string) (err error) {
		return zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied")
	}
}

func newMockProjectPermissionCheckAllowed() domain.ProjectPermissionCheck {
	return func(ctx context.Context, clientID, userID string) (err error) {
		return nil
	}
}

func newMockProjectPermissionCheckOIDCNotAllowed() domain.ProjectPermissionCheck {
	return func(ctx context.Context, clientID, userID string) (err error) {
		return zerrors.ThrowPermissionDenied(nil, "OIDC-foSyH49RvL", "Errors.PermissionDenied")
	}
}

func newMockProjectPermissionCheckSAMLNotAllowed() domain.ProjectPermissionCheck {
	return func(ctx context.Context, clientID, userID string) (err error) {
		return zerrors.ThrowPermissionDenied(nil, "SAML-foSyH49RvL", "Errors.PermissionDenied")
	}
}

func newMockTokenVerifierValid() func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		return nil
	}
}
func newMockTokenVerifierInvalid() func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
	return func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
		return zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
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
func mockPasswordHasher(x string) *crypto.Hasher {
	return &crypto.Hasher{
		Swapper:  passwap.NewSwapper(plainHasher{x: x}),
		Prefixes: []string{"$plain$"},
	}
}
