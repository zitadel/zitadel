package command

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/repository/samlsession"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func mockSAMLRequestComplianceChecker(returnErr error) SAMLRequestComplianceChecker {
	return func(context.Context, *SAMLRequestWriteModel) error {
		return returnErr
	}
}

func TestCommands_CreateSAMLSessionFromSAMLRequest(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id.Generator
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx                  context.Context
		samlRequestID        string
		samlResponseID       string
		complianceCheck      SAMLRequestComplianceChecker
		samlResponseLifetime time.Duration
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing code",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:   "",
				complianceCheck: mockSAMLRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-0LxK6O31wH", "Errors.SAMLRequest.InvalidCode"),
			},
		},
		{
			"filter error",
			fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:   "V2_samlRequestID",
				complianceCheck: mockSAMLRequestComplianceChecker(nil),
			},
			res{
				err: io.ErrClosedPipe,
			},
		},
		{
			"session filter error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"loginClient",
								"applicationId",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
								"responseissuer",
							),
						),
					),
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:   "V2_samlRequestID",
				complianceCheck: mockSAMLRequestComplianceChecker(nil),
			},
			res{
				err: io.ErrClosedPipe,
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"loginClient",
								"applicationId",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
								"responseissuer",
							),
						),
						eventFromEventPusher(
							samlrequest.NewSessionLinkedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectFilter(), // inactive session
				),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:   "V2_samlRequestID",
				complianceCheck: mockSAMLRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
			},
		},
		{
			"user not active",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"loginClient",
								"applicationId",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
								"responseissuer",
							),
						),
						eventFromEventPusher(
							samlrequest.NewSessionLinkedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								testNow),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.Afrikaans,
							domain.GenderUnspecified,
							"email",
							false,
						),
						user.NewUserDeactivatedEvent(
							context.Background(),
							&user.NewAggregate("userID", "org1").Aggregate,
						),
					),
				),
				idGenerator:  mock.NewIDGeneratorExpectIDs(t),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:                  authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:        "V2_samlRequestID",
				samlResponseID:       "samlResponseID",
				samlResponseLifetime: time.Minute * 5,
				complianceCheck:      mockSAMLRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "SAML-1768ZQpmcP", "Errors.User.NotActive"),
			},
		},
		{
			"add successful",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							samlrequest.NewAddedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"loginClient",
								"applicationId",
								"acs",
								"relaystate",
								"request",
								"binding",
								"issuer",
								"destination",
								"responseissuer",
							),
						),
						eventFromEventPusher(
							samlrequest.NewSessionLinkedEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								testNow),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("userID", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.Afrikaans,
							domain.GenderUnspecified,
							"email",
							false,
						),
					),
					expectPush(
						samlsession.NewAddedEvent(context.Background(), &samlsession.NewAggregate("V2_samlSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "issuer", []string{"issuer"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						samlsession.NewSAMLResponseAddedEvent(context.Background(), &samlsession.NewAggregate("V2_samlSessionID", "org1").Aggregate, "samlResponseID", time.Minute*5),
						samlrequest.NewSucceededEvent(context.Background(), &samlrequest.NewAggregate("V2_samlRequestID", "instanceID").Aggregate),
					),
				),
				idGenerator:  mock.NewIDGeneratorExpectIDs(t, "samlSessionID"),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:                  authz.WithInstanceID(context.Background(), "instanceID"),
				samlRequestID:        "V2_samlRequestID",
				samlResponseID:       "samlResponseID",
				samlResponseLifetime: time.Minute * 5,
				complianceCheck:      mockSAMLRequestComplianceChecker(nil),
			},
			res{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore(t),
				idGenerator:  tt.fields.idGenerator,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			c.setMilestonesCompletedForTest("instanceID")
			err := c.CreateSAMLSessionFromSAMLRequest(tt.args.ctx, tt.args.samlRequestID, tt.args.complianceCheck, tt.args.samlResponseID, tt.args.samlResponseLifetime)
			require.ErrorIs(t, err, tt.res.err)
		})
	}
}
