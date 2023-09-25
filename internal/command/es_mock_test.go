package command

import (
	"github.com/zitadel/zitadel/internal/eventstore/mock"
)

// TODO: Delete this file and reference the mock package directly in the tests

var (
	eventstoreExpect                                   = mock.EventstoreExpect
	expectEventstore                                   = mock.ExpectEventstore
	expectPush                                         = mock.ExpectPush
	expectPushFailed                                   = mock.ExpectPushFailed
	expectRandomPush                                   = mock.ExpectRandomPush
	eventFromEventPusherWithInstanceID                 = mock.EventFromEventPusherWithInstanceID
	eventFromEventPusherWithCreationDateNow            = mock.EventFromEventPusherWithCreationDateNow
	eventPusherToEvents                                = mock.EventPusherToEvents
	expectFilter                                       = mock.ExpectFilter
	expectFilterError                                  = mock.ExpectFilterError
	expectFilterOrgDomainNotFound                      = mock.ExpectFilterOrgDomainNotFound
	expectFilterOrgMemberNotFound                      = mock.ExpectFilterOrgMemberNotFound
	expectRandomPushFailed                             = mock.ExpectRandomPushFailed
	newMockPermissionCheckAllowed                      = mock.NewMockPermissionCheckAllowed
	newMockPermissionCheckNotAllowed                   = mock.NewMockPermissionCheckNotAllowed
	uniqueConstraintsFromEventConstraint               = mock.UniqueConstraintsFromEventConstraint
	uniqueConstraintsFromEventConstraintWithInstanceID = mock.UniqueConstraintsFromEventConstraintWithInstanceID
	eventFromEventPusher                               = mock.EventFromEventPusher
	GetMockSecretGenerator                             = mock.GetMockSecretGenerator
	mockPasswordHasher                                 = mock.MockPasswordHasher
)

type mockInstance struct {
	mock.MockInstance
}

type expect mock.Expecter

func toExpecters(expects ...expect) []mock.Expecter {
	expecters := make([]mock.Expecter, len(expects))
	for i := range expects {
		expecters[i] = mock.Expecter(expects[i])
	}
	return expecters
}
