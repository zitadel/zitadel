package command

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddSAMLRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx     context.Context
		request *SAMLRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CurrentSAMLRequest
		wantErr error
	}{
		{
			"already exists error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
				),
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args{
				ctx:     mockCtx,
				request: &SAMLRequest{},
			},
			nil,
			zerrors.ThrowPreconditionFailed(nil, "COMMAND-MO3vmsMLUt", "Errors.SAMLRequest.AlreadyExisting"),
		},
		{
			"added",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"login",
							"application",
							"acs",
							"relaystate",
							"request",
							"binding",
							"issuer",
							"destination",
						),
					),
				),
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "id"),
			},
			args{
				ctx: mockCtx,
				request: &SAMLRequest{
					LoginClient:   "login",
					ApplicationID: "application",
					ACSURL:        "acs",
					RelayState:    "relaystate",
					RequestID:     "request",
					Binding:       "binding",
					Issuer:        "issuer",
					Destination:   "destination",
				},
			},
			&CurrentSAMLRequest{
				SAMLRequest: &SAMLRequest{
					ID:            "V2_id",
					LoginClient:   "login",
					ApplicationID: "application",
					ACSURL:        "acs",
					RelayState:    "relaystate",
					RequestID:     "request",
					Binding:       "binding",
					Issuer:        "issuer",
					Destination:   "destination",
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.AddSAMLRequest(tt.args.ctx, tt.args.request)
			require.ErrorIs(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_LinkSessionToSAMLRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore    func(t *testing.T) *eventstore.Eventstore
		tokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
	}
	type args struct {
		ctx              context.Context
		id               string
		sessionID        string
		sessionToken     string
		checkLoginClient bool
	}
	type res struct {
		details *domain.ObjectDetails
		authReq *CurrentSAMLRequest
		wantErr error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"samlRequest not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowNotFound(nil, "COMMAND-GH3PVLSfXC", "Errors.SAMLRequest.NotExisting"),
			},
		},
		{
			"samlRequest not existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
						eventFromEventPusher(
							samlrequest.NewFailedEvent(mockCtx, &samlrequest.NewAggregate("id", "instanceID").Aggregate,
								domain.SAMLErrorReasonUnspecified,
							),
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-ttPKNdAIFT", "Errors.SAMLRequest.AlreadyHandled"),
			},
		},
		{
			"wrong login client",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:              authz.NewMockContext("instanceID", "orgID", "wrongLoginClient"),
				id:               "id",
				sessionID:        "sessionID",
				sessionToken:     "token",
				checkLoginClient: true,
			},
			res{
				wantErr: zerrors.ThrowPermissionDenied(nil, "COMMAND-KCd48Rxt7x", "Errors.SAMLRequest.WrongLoginClient"),
			},
		},
		{
			"session not existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectFilter(),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:       mockCtx,
				id:        "V2_id",
				sessionID: "sessionID",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
			},
		},
		{
			"session expired",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow.Add(-5*time.Minute), &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow.Add(-5*time.Minute)),
						),
						eventFromEventPusher(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
				),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hkl3d", "Errors.Session.Expired"),
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
					),
				),
				tokenVerifier: newMockTokenVerifierInvalid(),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				wantErr: zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"linked",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPush(
						samlrequest.NewSessionLinkedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"sessionID",
							"userID",
							testNow,
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:          mockCtx,
				id:           "V2_id",
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				authReq: &CurrentSAMLRequest{
					SAMLRequest: &SAMLRequest{
						ID:            "V2_id",
						LoginClient:   "login",
						ApplicationID: "application",
						ACSURL:        "acs",
						RelayState:    "relaystate",
						RequestID:     "request",
						Binding:       "binding",
						Issuer:        "issuer",
						Destination:   "destination",
					},
					SessionID:   "sessionID",
					UserID:      "userID",
					AuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				},
			},
		},
		{
			"linked with login client check",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"loginClient",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(mockCtx,
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewUserCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							session.NewLifetimeSetEvent(mockCtx, &session.NewAggregate("sessionID", "instance1").Aggregate,
								2*time.Minute),
						),
					),
					expectPush(
						samlrequest.NewSessionLinkedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							"sessionID",
							"userID",
							testNow,
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
						),
					),
				),
				tokenVerifier: newMockTokenVerifierValid(),
			},
			args{
				ctx:              authz.NewMockContext("instanceID", "orgID", "loginClient"),
				id:               "V2_id",
				sessionID:        "sessionID",
				sessionToken:     "token",
				checkLoginClient: true,
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				authReq: &CurrentSAMLRequest{
					SAMLRequest: &SAMLRequest{
						ID:            "V2_id",
						LoginClient:   "loginClient",
						ApplicationID: "application",
						ACSURL:        "acs",
						RelayState:    "relaystate",
						RequestID:     "request",
						Binding:       "binding",
						Issuer:        "issuer",
						Destination:   "destination",
					},
					SessionID:   "sessionID",
					UserID:      "userID",
					AuthMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore(t),
				sessionTokenVerifier: tt.fields.tokenVerifier,
			}
			details, got, err := c.LinkSessionToSAMLRequest(tt.args.ctx, tt.args.id, tt.args.sessionID, tt.args.sessionToken, tt.args.checkLoginClient)
			require.ErrorIs(t, err, tt.res.wantErr)
			assertObjectDetails(t, tt.res.details, details)
			if err == nil {
				assert.WithinRange(t, got.AuthTime, testNow, testNow)
				got.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.res.authReq, got)
		})
	}
}

func TestCommands_FailSAMLRequest(t *testing.T) {
	mockCtx := authz.NewMockContext("instanceID", "orgID", "loginClient")
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx         context.Context
		id          string
		reason      domain.SAMLErrorReason
		description string
	}
	type res struct {
		details *domain.ObjectDetails
		samlReq *CurrentSAMLRequest
		wantErr error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"authRequest not existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:    mockCtx,
				id:     "foo",
				reason: domain.SAMLErrorReasonAuthNFailed,
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-32lGj1Fhjt", "Errors.SAMLRequest.AlreadyHandled"),
			},
		}, {
			"already failed",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
						samlrequest.NewFailedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							domain.SAMLErrorReasonAuthNFailed,
						),
					),
				),
			},
			args{
				ctx:         mockCtx,
				id:          "V2_id",
				reason:      domain.SAMLErrorReasonAuthNFailed,
				description: "desc",
			},
			res{
				wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-32lGj1Fhjt", "Errors.SAMLRequest.AlreadyHandled"),
			},
		},
		{
			"failed",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
								"login",
								"application",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
							),
						),
					),
					expectPush(
						samlrequest.NewFailedEvent(mockCtx, &samlrequest.NewAggregate("V2_id", "instanceID").Aggregate,
							domain.SAMLErrorReasonAuthNFailed,
						),
					),
				),
			},
			args{
				ctx:         mockCtx,
				id:          "V2_id",
				reason:      domain.SAMLErrorReasonAuthNFailed,
				description: "desc",
			},
			res{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
				samlReq: &CurrentSAMLRequest{
					SAMLRequest: &SAMLRequest{
						ID:            "V2_id",
						LoginClient:   "login",
						ApplicationID: "application",
						ACSURL:        "acs",
						RelayState:    "relaystate",
						RequestID:     "request",
						Binding:       "binding",
						Issuer:        "issuer",
						Destination:   "destination",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			details, got, err := c.FailSAMLRequest(tt.args.ctx, tt.args.id, tt.args.reason)
			require.ErrorIs(t, err, tt.res.wantErr)
			assertObjectDetails(t, tt.res.details, details)
			assert.Equal(t, tt.res.samlReq, got)
		})
	}
}
