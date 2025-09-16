package command

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_ImportHumanRecoveryCodes(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		codes         []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "user not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // user not found
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				codes:         []string{"code1", "code2"},
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-uXHNj", "Errors.User.NotFound"),
		},
		{
			name: "recovery codes already exist, add to existing",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"hashedcode1", "hashedcode2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodesAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							[]string{"$plain$$code1", "$plain$$code2"},
							nil,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				codes:         []string{"code1", "code2"},
			},
		},
		{
			name: "empty codes, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(), // no existing recovery codes
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				codes:         []string{},
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "DOMAIN-vee93", "Errors.User.MFA.RecoveryCodes.InvalidCount"),
		},
		{
			name: "max count exceeded with existing codes, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"code1", "code2", "code3", "code4", "code5", "code6", "code7", "code8"},
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				codes:         []string{"code9", "code10", "code11"}, // 8 existing + 3 new = 11 > max 10
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-53cjw", "Errors.User.MFA.RecoveryCodes.MaxCountExceeded"),
		},
		{
			name: "successful import",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(), // no existing recovery codes
					expectPush(
						user.NewHumanRecoveryCodesAddedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							[]string{"$plain$$code1", "$plain$$code2"},
							nil,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				codes:         []string{"code1", "code2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore(t),
				secretHasher: mockPasswordHasher(""),
				multifactors: domain.MultifactorConfigs{
					RecoveryCodes: domain.RecoveryCodesConfig{
						MaxCount:   10,
						Format:     domain.RecoveryCodeFormatAlphanumeric,
						Length:     8,
						WithHyphen: false,
					},
				},
			}
			err := c.ImportHumanRecoveryCodes(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.codes)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_GenerateRecoveryCodes(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		userID        string
		count         int
		resourceOwner string
		authRequest   *domain.AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RecoveryCodesDetails
		wantErr error
	}{
		{
			name: "missing userID, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "",
				count:         2,
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-4kje7", "Errors.User.UserIDMissing"),
		},
		{
			name: "user not found, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // user not found
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				count:         2,
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-uXHNj", "Errors.User.NotFound"),
		},
		{
			name: "permission denied, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				count:         2,
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "invalid count (zero), error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				count:         0,
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-7c0nx", "Errors.User.RecoveryCodes.CountInvalid"),
		},
		{
			name: "invalid count (too high), error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				count:         15, // exceeds max count of 10
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-7c0nx", "Errors.User.RecoveryCodes.CountInvalid"),
		},
		{
			name: "max count exceeded with existing codes, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email@test.ch",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"code1", "code2", "code3", "code4", "code5", "code6", "code7", "code8"},
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				count:         5, // 8 existing + 5 new = 13 > max 10
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowAlreadyExists(nil, "COMMAND-8f2k9", "Errors.User.MFA.RecoveryCodes.MaxCountExceeded"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
				secretHasher:    mockPasswordHasher(""),
				multifactors: domain.MultifactorConfigs{
					RecoveryCodes: domain.RecoveryCodesConfig{
						MaxCount:   10,
						Format:     domain.RecoveryCodeFormatAlphanumeric,
						Length:     8,
						WithHyphen: false,
					},
				},
			}
			got, err := c.GenerateRecoveryCodes(tt.args.ctx, tt.args.userID, tt.args.count, tt.args.resourceOwner, tt.args.authRequest)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want != nil {
				require.NotNil(t, got)
				assert.Equal(t, tt.want.ObjectDetails.ResourceOwner, got.ObjectDetails.ResourceOwner)
			}
		})
	}
}

func TestCommands_RemoveRecoveryCodes(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore      func(*testing.T) *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx           context.Context
		userID        string
		resourceOwner string
		authRequest   *domain.AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "missing userID, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-l2n9r", "Errors.User.UserIDMissing"),
		},
		{
			name: "permission denied, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
								[]string{"code1", "code2"},
								nil,
							),
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
		},
		{
			name: "user locked, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"code1", "code2"},
								nil,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-d9u8q", "Errors.User.RecoveryCodes.Locked"),
		},
		{
			name: "recovery codes not added, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // no recovery codes
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotAdded"),
		},
		{
			name: "successful removal",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"code1", "code2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodeRemovedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							nil,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
		{
			name: "successful removal with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"code1", "code2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodeRemovedEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
		{
			name: "successful removal, other user",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user2", "org1").Aggregate,
								[]string{"code1", "code2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodeRemovedEvent(ctx,
							&user.NewAggregate("user2", "org1").Aggregate,
							nil,
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user2",
				resourceOwner: "org1",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore(t),
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.RemoveRecoveryCodes(tt.args.ctx, tt.args.userID, tt.args.resourceOwner, tt.args.authRequest)
			require.ErrorIs(t, err, tt.wantErr)
			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCommands_HumanCheckRecoveryCode(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		userID        string
		code          string
		resourceOwner string
		authRequest   *domain.AuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "missing code, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-u0b6c", "Errors.User.UserIDMissing"),
		},
		{
			name: "missing userID, error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           ctx,
				userID:        "",
				code:          "validcode",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9s2", "Errors.User.UserIDMissing"),
		},
		{
			name: "user locked, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"$plain$$validcode", "$plain$$validcode2"},
								nil,
							),
						),
						eventFromEventPusher(
							user.NewUserLockedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
							),
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowNotFound(nil, "COMMAND-2w6oa", "Errors.User.MFA.RecoveryCodes.Locked"),
		},
		{
			name: "recovery codes not ready, error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // no recovery codes
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotReady"),
		},
		{
			name: "valid code, successful check",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"$plain$$validcode", "$plain$$validcode2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodeCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$$validcode",
							nil,
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
			},
		},
		{
			name: "valid code, successful check with auth request",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanRecoveryCodesAddedEvent(ctx,
								&user.NewAggregate("user1", "org1").Aggregate,
								[]string{"$plain$$validcode", "$plain$$validcode2"},
								nil,
							),
						),
					),
					expectPush(
						user.NewHumanRecoveryCodeCheckSucceededEvent(ctx,
							&user.NewAggregate("user1", "org1").Aggregate,
							"$plain$$validcode",
							&user.AuthRequestInfo{
								ID:          "authRequestID",
								UserAgentID: "userAgentID",
								BrowserInfo: &user.BrowserInfo{
									UserAgent:      "user-agent",
									AcceptLanguage: "en",
									RemoteIP:       net.IP{192, 0, 2, 1},
								},
							},
						),
					),
				),
			},
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
				authRequest: &domain.AuthRequest{
					ID:      "authRequestID",
					AgentID: "userAgentID",
					BrowserInfo: &domain.BrowserInfo{
						UserAgent:      "user-agent",
						AcceptLanguage: "en",
						RemoteIP:       net.IP{192, 0, 2, 1},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore(t),
				secretHasher: mockPasswordHasher(""),
			}
			err := c.HumanCheckRecoveryCode(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.authRequest)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_checkRecoveryCode(t *testing.T) {
	ctx := authz.NewMockContext("inst1", "org1", "user1")
	hasher := mockPasswordHasher("")

	queryReducer := func(ctx context.Context, r eventstore.QueryReducer) error {
		// Basic query reducer that doesn't set any specific state
		return nil
	}
	queryReducerError := func(ctx context.Context, r eventstore.QueryReducer) error {
		return io.ErrClosedPipe
	}

	// Query reducer that sets up a write model with valid recovery codes
	queryReducerWithCodes := func(ctx context.Context, r eventstore.QueryReducer) error {
		switch wm := r.(type) {
		case *HumanRecoveryCodeWriteModel:
			wm.State = domain.MFAStateReady
			wm.codes = []string{"$plain$$validcode", "$plain$$validcode2"}
		case *OrgLockoutPolicyWriteModel:
			// Set up lockout policy - no lockout for this test
			wm.MaxOTPAttempts = 0
			wm.State = domain.PolicyStateActive
		}
		return nil
	}

	// Query reducer that sets up a locked user
	queryReducerUserLocked := func(ctx context.Context, r eventstore.QueryReducer) error {
		switch wm := r.(type) {
		case *HumanRecoveryCodeWriteModel:
			wm.State = domain.MFAStateReady
			wm.codes = []string{"$plain$$validcode", "$plain$$validcode2"}
			wm.userLocked = true
		case *OrgLockoutPolicyWriteModel:
			// Default lockout policy
			wm.State = domain.PolicyStateActive
		}
		return nil
	}

	// Query reducer that sets up recovery codes not ready
	queryReducerNotReady := func(ctx context.Context, r eventstore.QueryReducer) error {
		switch wm := r.(type) {
		case *HumanRecoveryCodeWriteModel:
			wm.State = domain.MFAStateNotReady
		case *OrgLockoutPolicyWriteModel:
			// Default lockout policy
			wm.State = domain.PolicyStateActive
		}
		return nil
	}

	type args struct {
		ctx           context.Context
		userID        string
		code          string
		resourceOwner string
		authRequest   *domain.AuthRequest
		queryReducer  func(ctx context.Context, r eventstore.QueryReducer) error
		secretHasher  *crypto.Hasher
	}
	tests := []struct {
		name         string
		args         args
		wantCommands int
		wantErr      error
	}{
		{
			name: "missing code, error",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "",
				resourceOwner: "org1",
				queryReducer:  queryReducer,
				secretHasher:  hasher,
			},
			wantCommands: 0,
			wantErr:      zerrors.ThrowInvalidArgument(nil, "COMMAND-u0b6c", "Errors.User.UserIDMissing"),
		},
		{
			name: "missing userID, error",
			args: args{
				ctx:           ctx,
				userID:        "",
				code:          "validcode",
				resourceOwner: "org1",
				queryReducer:  queryReducer,
				secretHasher:  hasher,
			},
			wantCommands: 0,
			wantErr:      zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9s2", "Errors.User.UserIDMissing"),
		},
		{
			name: "query reducer error",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
				queryReducer:  queryReducerError,
				secretHasher:  hasher,
			},
			wantCommands: 0,
			wantErr:      io.ErrClosedPipe,
		},
		{
			name: "user locked, error",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
				queryReducer:  queryReducerUserLocked,
				secretHasher:  hasher,
			},
			wantCommands: 0,
			wantErr:      zerrors.ThrowNotFound(nil, "COMMAND-2w6oa", "Errors.User.MFA.RecoveryCodes.Locked"),
		},
		{
			name: "recovery codes not ready, error",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
				queryReducer:  queryReducerNotReady,
				secretHasher:  hasher,
			},
			wantCommands: 0,
			wantErr:      zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotReady"),
		},
		{
			name: "invalid code, returns failed event and error",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "invalidcode",
				resourceOwner: "org1",
				queryReducer:  queryReducerWithCodes,
				secretHasher:  hasher,
			},
			wantCommands: 1, // should return failed event command
			wantErr:      zerrors.ThrowInvalidArgument(nil, "DOMAIN-6uvh0", "Errors.User.MFA.RecoveryCodes.InvalidCode"),
		},
		{
			name: "valid code, returns success event",
			args: args{
				ctx:           ctx,
				userID:        "user1",
				code:          "validcode",
				resourceOwner: "org1",
				queryReducer:  queryReducerWithCodes,
				secretHasher:  hasher,
			},
			wantCommands: 1, // should return success event command
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands, err := checkRecoveryCode(tt.args.ctx, tt.args.userID, tt.args.code, tt.args.resourceOwner, tt.args.authRequest, tt.args.queryReducer, tt.args.secretHasher)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Len(t, commands, tt.wantCommands)
		})
	}
}
