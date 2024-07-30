package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateOTPSMSChallengeReturnCode(t *testing.T) {
	type fields struct {
		userID     string
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode encryptedCodeWithDefaultFunc
	}
	type res struct {
		err        error
		returnCode string
		commands   []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "userID missing, precondition error",
			fields: fields{
				userID:     "",
				eventstore: expectEventstore(),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKL3g", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp not ready, precondition error",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockEncryptedCodeWithDefault("1234567", 5*time.Minute),
			},
			res: res{
				returnCode: "1234567",
				commands: []eventstore.Command{
					session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						true,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				// config will not be actively used for the test (is only for default),
				// but not providing it would result in a nil pointer
				defaultSecretGenerators: &SecretGenerators{
					OTPSMS: emptyConfig,
				},
			}
			var dst string
			cmd := c.CreateOTPSMSChallengeReturnCode(&dst)

			sessionModel := &SessionWriteModel{
				UserID:        tt.fields.userID,
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Empty(t, gotCmds)
			assert.Equal(t, tt.res.returnCode, dst)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPSMSChallenge(t *testing.T) {
	type fields struct {
		userID     string
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode encryptedCodeWithDefaultFunc
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "userID missing, precondition error",
			fields: fields{
				userID:     "",
				eventstore: expectEventstore(),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKL3g", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp not ready, precondition error",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-BJ2g3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockEncryptedCodeWithDefault("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				// config will not be actively used for the test (is only for default),
				// but not providing it would result in a nil pointer
				defaultSecretGenerators: &SecretGenerators{
					OTPSMS: emptyConfig,
				},
			}

			cmd := c.CreateOTPSMSChallenge()

			sessionModel := &SessionWriteModel{
				UserID:        tt.fields.userID,
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Empty(t, gotCmds)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_OTPSMSSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		sessionID     string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "not challenged, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-G3t31", "Errors.User.Code.NotFound"),
		},
		{
			name: "challenged and sent",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewOTPSMSChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("1234567"),
								},
								5*time.Minute,
								false,
							),
						),
					),
					expectPush(
						session.NewOTPSMSSentEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.OTPSMSSent(tt.args.ctx, tt.args.sessionID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_CreateOTPEmailChallengeURLTemplate(t *testing.T) {
	type fields struct {
		userID     string
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode encryptedCodeWithDefaultFunc
	}
	type args struct {
		urlTmpl string
	}
	type res struct {
		templateError error
		err           error
		commands      []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid template, precondition error",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.InvalidField}}",
			},
			fields: fields{
				eventstore: expectEventstore(),
			},
			res: res{
				templateError: zerrors.ThrowInvalidArgument(nil, "DOMAIN-ieYa7", "Errors.User.InvalidURLTemplate"),
			},
		},
		{
			name: "userID missing, precondition error",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}&lang={{.PreferredLanguage}}",
			},
			fields: fields{
				eventstore: expectEventstore(),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JK3gp", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp not ready, precondition error",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}&lang={{.PreferredLanguage}}",
			},
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			args: args{
				urlTmpl: "https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}&lang={{.PreferredLanguage}}",
			},
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockEncryptedCodeWithDefault("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
						"https://example.com/mfa/email?userID={{.UserID}}&code={{.Code}}&lang={{.PreferredLanguage}}",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				// config will not be actively used for the test (is only for default),
				// but not providing it would result in a nil pointer
				defaultSecretGenerators: &SecretGenerators{
					OTPEmail: emptyConfig,
				},
			}

			cmd, err := c.CreateOTPEmailChallengeURLTemplate(tt.args.urlTmpl)
			assert.ErrorIs(t, err, tt.res.templateError)
			if tt.res.templateError != nil {
				return
			}

			sessionModel := &SessionWriteModel{
				UserID:        tt.fields.userID,
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Empty(t, gotCmds)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPEmailChallengeReturnCode(t *testing.T) {
	type fields struct {
		userID     string
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode encryptedCodeWithDefaultFunc
	}
	type res struct {
		err        error
		returnCode string
		commands   []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "userID missing, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JK3gp", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp not ready, precondition error",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockEncryptedCodeWithDefault("1234567", 5*time.Minute),
			},
			res: res{
				returnCode: "1234567",
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						true,
						"",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				// config will not be actively used for the test (is only for default),
				// but not providing it would result in a nil pointer
				defaultSecretGenerators: &SecretGenerators{
					OTPEmail: emptyConfig,
				},
			}
			var dst string
			cmd := c.CreateOTPEmailChallengeReturnCode(&dst)

			sessionModel := &SessionWriteModel{
				UserID:        tt.fields.userID,
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Empty(t, gotCmds)
			assert.Equal(t, tt.res.returnCode, dst)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_CreateOTPEmailChallenge(t *testing.T) {
	type fields struct {
		userID     string
		eventstore func(*testing.T) *eventstore.Eventstore
		createCode encryptedCodeWithDefaultFunc
	}
	type res struct {
		err      error
		commands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "userID missing, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JK3gp", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "otp not ready, precondition error",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-JKLJ3", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "generate code",
			fields: fields{
				userID: "userID",
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org").Aggregate),
						),
					),
				),
				createCode: mockEncryptedCodeWithDefault("1234567", 5*time.Minute),
			},
			res: res{
				commands: []eventstore.Command{
					session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("1234567"),
						},
						5*time.Minute,
						false,
						"",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				// config will not be actively used for the test (is only for default),
				// but not providing it would result in a nil pointer
				defaultSecretGenerators: &SecretGenerators{
					OTPEmail: emptyConfig,
				},
			}

			cmd := c.CreateOTPEmailChallenge()

			sessionModel := &SessionWriteModel{
				UserID:        tt.fields.userID,
				UserCheckedAt: testNow,
				State:         domain.SessionStateActive,
				aggregate:     &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				createCode:        tt.fields.createCode,
				now:               time.Now,
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Empty(t, gotCmds)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCommands_OTPEmailSent(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		sessionID     string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "not challenged, precondition error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SLr02", "Errors.User.Code.NotFound"),
		},
		{
			name: "challenged and sent",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewOTPEmailChallengedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("1234567"),
								},
								5*time.Minute,
								false,
								"",
							),
						),
					),
					expectPush(
						session.NewOTPEmailSentEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				sessionID:     "sessionID",
				resourceOwner: "instanceID",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			err := c.OTPEmailSent(tt.args.ctx, tt.args.sessionID, tt.args.resourceOwner)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCheckOTPSMS(t *testing.T) {
	type fields struct {
		eventstore       func(*testing.T) *eventstore.Eventstore
		userID           string
		otpCodeChallenge *OTPCode
		otpAlg           crypto.EncryptionAlgorithm
	}
	type args struct {
		code string
	}
	type res struct {
		err           error
		commands      []eventstore.Command
		errorCommands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "missing code",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
			},
			args: args{},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty"),
			},
		},
		{
			name: "not set up",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				userID: "userID",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "missing challenge",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
				),
				userID:           "userID",
				otpCodeChallenge: nil,
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								0, 0, false,
							),
						),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
				errorCommands: []eventstore.Command{
					user.NewHumanOTPSMSCheckFailedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
				},
			},
		},
		{
			name: "invalid code, locked",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								0, 1, false,
							),
						),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
				errorCommands: []eventstore.Command{
					user.NewHumanOTPSMSCheckFailedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
					user.NewUserLockedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate),
				},
			},
		},
		{
			name: "check ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				commands: []eventstore.Command{
					user.NewHumanOTPSMSCheckSucceededEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
					session.NewOTPSMSCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						testNow,
					),
				},
			},
		},
		{
			name: "check ok, locked in the meantime",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPSMSAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(
						user.NewUserLockedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S6h4R", "Errors.User.Locked"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CheckOTPSMS(tt.args.code)

			sessionModel := &SessionWriteModel{
				UserID:              tt.fields.userID,
				UserCheckedAt:       testNow,
				State:               domain.SessionStateActive,
				OTPSMSCodeChallenge: tt.fields.otpCodeChallenge,
				aggregate:           &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				otpAlg:            tt.fields.otpAlg,
				now: func() time.Time {
					return testNow
				},
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.errorCommands, gotCmds)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}

func TestCheckOTPEmail(t *testing.T) {
	type fields struct {
		eventstore       func(*testing.T) *eventstore.Eventstore
		userID           string
		otpCodeChallenge *OTPCode
		otpAlg           crypto.EncryptionAlgorithm
	}
	type args struct {
		code string
	}
	type res struct {
		err           error
		commands      []eventstore.Command
		errorCommands []eventstore.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "missing userID",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "missing code",
			fields: fields{
				eventstore: expectEventstore(),
				userID:     "userID",
			},
			args: args{},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty"),
			},
		},
		{
			name: "not set up",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				userID: "userID",
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-d2r52", "Errors.User.MFA.OTP.NotReady"),
			},
		},
		{
			name: "missing challenge",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
				),
				userID:           "userID",
				otpCodeChallenge: nil,
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S34gh", "Errors.User.Code.NotFound"),
			},
		},
		{
			name: "invalid code",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								0, 0, false,
							),
						),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
				errorCommands: []eventstore.Command{
					user.NewHumanOTPEmailCheckFailedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
				},
			},
		},
		{
			name: "invalid code, locked",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(
							org.NewLockoutPolicyAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								0, 1, false,
							),
						),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow.Add(-10 * time.Minute),
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired"),
				errorCommands: []eventstore.Command{
					user.NewHumanOTPEmailCheckFailedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
					user.NewUserLockedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate),
				},
			},
		},
		{
			name: "check ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(), // recheck
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				commands: []eventstore.Command{
					user.NewHumanOTPEmailCheckSucceededEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
					session.NewOTPEmailCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
						testNow,
					),
				},
			},
		},
		{
			name: "check ok, locked in the meantime",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(user.NewHumanOTPEmailAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate)),
					),
					expectFilter(
						user.NewUserLockedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate),
					),
				),
				userID: "userID",
				otpCodeChallenge: &OTPCode{
					Code: &crypto.CryptoValue{
						CryptoType: crypto.TypeEncryption,
						Algorithm:  "enc",
						KeyID:      "id",
						Crypted:    []byte("code"),
					},
					Expiry:       5 * time.Minute,
					CreationDate: testNow,
				},
				otpAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				code: "code",
			},
			res: res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-S6h4R", "Errors.User.Locked"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CheckOTPEmail(tt.args.code)

			sessionModel := &SessionWriteModel{
				UserID:                tt.fields.userID,
				UserCheckedAt:         testNow,
				State:                 domain.SessionStateActive,
				OTPEmailCodeChallenge: tt.fields.otpCodeChallenge,
				aggregate:             &session.NewAggregate("sessionID", "instanceID").Aggregate,
			}
			cmds := &SessionCommands{
				sessionCommands:   []SessionCommand{cmd},
				sessionWriteModel: sessionModel,
				eventstore:        tt.fields.eventstore(t),
				otpAlg:            tt.fields.otpAlg,
				now: func() time.Time {
					return testNow
				},
			}

			gotCmds, err := cmd(context.Background(), cmds)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.errorCommands, gotCmds)
			assert.Equal(t, tt.res.commands, cmds.eventCommands)
		})
	}
}
