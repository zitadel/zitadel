//go:build integration

package user_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	object_v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	user_v2beta "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func detailsV2ToV2beta(obj *object.Details) *object_v2beta.Details {
	return &object_v2beta.Details{
		Sequence:      obj.GetSequence(),
		ChangeDate:    obj.GetChangeDate(),
		ResourceOwner: obj.GetResourceOwner(),
	}
}

var (
	permissionCheckV2SetFlagInital bool
	permissionCheckV2SetFlag       bool
)

type permissionCheckV2SettingsStruct struct {
	TestNamePrependString string
	SetFlag               bool
}

var permissionCheckV2Settings []permissionCheckV2SettingsStruct = []permissionCheckV2SettingsStruct{
	{
		SetFlag:               false,
		TestNamePrependString: "permission_check_v2 IS NOT SET" + " ",
	},
	{
		SetFlag:               true,
		TestNamePrependString: "permission_check_v2 IS SET" + " ",
	},
}

func setPermissionCheckV2Flag(t *testing.T, setFlag bool) {
	if permissionCheckV2SetFlagInital && permissionCheckV2SetFlag == setFlag {
		return
	}

	_, err := Instance.Client.FeatureV2.SetInstanceFeatures(IamCTX, &feature.SetInstanceFeaturesRequest{
		PermissionCheckV2: &setFlag,
	})
	require.NoError(t, err)

	var flagSet bool
	for i := 0; !flagSet || i < 6; i++ {
		res, err := Instance.Client.FeatureV2.GetInstanceFeatures(IamCTX, &feature.GetInstanceFeaturesRequest{})
		require.NoError(t, err)
		if res.PermissionCheckV2.Enabled == setFlag {
			flagSet = true
			continue
		}
		time.Sleep(10 * time.Second)
	}

	if !flagSet {
		require.NoError(t, errors.New("unable to set permission_check_v2 flag"))
	}
	permissionCheckV2SetFlagInital = true
	permissionCheckV2SetFlag = setFlag
}

func TestServer_GetUserByID(t *testing.T) {
	orgResp := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	type args struct {
		ctx context.Context
		req *user_v2beta.GetUserByIDRequest
		dep func(ctx context.Context, request *user_v2beta.GetUserByIDRequest) *userAttr
	}
	tests := []struct {
		name    string
		args    args
		want    *user_v2beta.GetUserByIDResponse
		wantErr bool
	}{
		{
			name: "user by ID, no id provided",
			args: args{
				IamCTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: "",
				},
				func(ctx context.Context, request *user_v2beta.GetUserByIDRequest) *userAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "user by ID, not found",
			args: args{
				IamCTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: "unknown",
				},
				func(ctx context.Context, request *user_v2beta.GetUserByIDRequest) *userAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "user by ID, ok",
			args: args{
				IamCTX,
				&user_v2beta.GetUserByIDRequest{},
				func(ctx context.Context, request *user_v2beta.GetUserByIDRequest) *userAttr {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.UserId = info.UserID
					return &info
				},
			},
			want: &user_v2beta.GetUserByIDResponse{
				User: &user_v2beta.User{
					State:              user_v2beta.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user_v2beta.User_Human{
						Human: &user_v2beta.HumanUser{
							Profile: &user_v2beta.HumanProfile{
								GivenName:         "Mickey",
								FamilyName:        "Mouse",
								NickName:          gu.Ptr("Mickey"),
								DisplayName:       gu.Ptr("Mickey Mouse"),
								PreferredLanguage: gu.Ptr("nl"),
								Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								AvatarUrl:         "",
							},
							Email: &user_v2beta.HumanEmail{
								IsVerified: true,
							},
							Phone: &user_v2beta.HumanPhone{
								IsVerified: true,
							},
						},
					},
				},
				Details: &object_v2beta.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: orgResp.OrganizationId,
				},
			},
		},
		{
			name: "user by ID, passwordChangeRequired, ok",
			args: args{
				IamCTX,
				&user_v2beta.GetUserByIDRequest{},
				func(ctx context.Context, request *user_v2beta.GetUserByIDRequest) *userAttr {
					info := createUser(ctx, orgResp.OrganizationId, true)
					request.UserId = info.UserID
					return &info
				},
			},
			want: &user_v2beta.GetUserByIDResponse{
				User: &user_v2beta.User{
					State:              user_v2beta.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user_v2beta.User_Human{
						Human: &user_v2beta.HumanUser{
							Profile: &user_v2beta.HumanProfile{
								GivenName:         "Mickey",
								FamilyName:        "Mouse",
								NickName:          gu.Ptr("Mickey"),
								DisplayName:       gu.Ptr("Mickey Mouse"),
								PreferredLanguage: gu.Ptr("nl"),
								Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								AvatarUrl:         "",
							},
							Email: &user_v2beta.HumanEmail{
								IsVerified: true,
							},
							Phone: &user_v2beta.HumanPhone{
								IsVerified: true,
							},
							PasswordChangeRequired: true,
							PasswordChanged:        timestamppb.Now(),
						},
					},
				},
				Details: &object_v2beta.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: orgResp.OrganizationId,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userAttr := tt.args.dep(IamCTX, tt.args.req)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetUserByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				if !assert.NoError(ttt, err) {
					return
				}

				tt.want.User.Details = detailsV2ToV2beta(userAttr.Details)
				tt.want.User.UserId = userAttr.UserID
				tt.want.User.Username = userAttr.Username
				tt.want.User.PreferredLoginName = userAttr.Username
				tt.want.User.LoginNames = []string{userAttr.Username}
				if human := tt.want.User.GetHuman(); human != nil {
					human.Email.Email = userAttr.Username
					human.Phone.Phone = userAttr.Phone
					if tt.want.User.GetHuman().GetPasswordChanged() != nil {
						human.PasswordChanged = userAttr.Changed
					}
				}
				assert.Equal(ttt, tt.want.User, got.User)
				integration.AssertDetails(ttt, tt.want, got)
			}, retryDuration, tick)
		})
	}
}

func TestServer_GetUserByID_Permission(t *testing.T) {
	timeNow := time.Now().UTC()
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	newUserID := newOrg.CreatedAdmins[0].GetUserId()
	type args struct {
		ctx context.Context
		req *user_v2beta.GetUserByIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user_v2beta.GetUserByIDResponse
		wantErr bool
	}{
		{
			name: "System, ok",
			args: args{
				SystemCTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			want: &user_v2beta.GetUserByIDResponse{
				User: &user_v2beta.User{
					State:              user_v2beta.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user_v2beta.User_Human{
						Human: &user_v2beta.HumanUser{
							Profile: &user_v2beta.HumanProfile{
								GivenName:         "firstname",
								FamilyName:        "lastname",
								NickName:          gu.Ptr(""),
								DisplayName:       gu.Ptr("firstname lastname"),
								PreferredLanguage: gu.Ptr("und"),
								Gender:            user_v2beta.Gender_GENDER_UNSPECIFIED.Enum(),
								AvatarUrl:         "",
							},
							Email: &user_v2beta.HumanEmail{
								Email: newOrgOwnerEmail,
							},
							Phone: &user_v2beta.HumanPhone{},
						},
					},
				},
				Details: &object_v2beta.Details{
					ChangeDate:    timestamppb.New(timeNow),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Instance, ok",
			args: args{
				IamCTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			want: &user_v2beta.GetUserByIDResponse{
				User: &user_v2beta.User{
					State:              user_v2beta.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user_v2beta.User_Human{
						Human: &user_v2beta.HumanUser{
							Profile: &user_v2beta.HumanProfile{
								GivenName:         "firstname",
								FamilyName:        "lastname",
								NickName:          gu.Ptr(""),
								DisplayName:       gu.Ptr("firstname lastname"),
								PreferredLanguage: gu.Ptr("und"),
								Gender:            user_v2beta.Gender_GENDER_UNSPECIFIED.Enum(),
								AvatarUrl:         "",
							},
							Email: &user_v2beta.HumanEmail{
								Email: newOrgOwnerEmail,
							},
							Phone: &user_v2beta.HumanPhone{},
						},
					},
				},
				Details: &object_v2beta.Details{
					ChangeDate:    timestamppb.New(timeNow),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Org, error",
			args: args{
				CTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			wantErr: true,
		},
		{
			name: "User, error",
			args: args{
				UserCTX,
				&user_v2beta.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetUserByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				if !assert.NoError(ttt, err) {
					return
				}

				tt.want.User.UserId = tt.args.req.GetUserId()
				tt.want.User.Username = newOrgOwnerEmail
				tt.want.User.PreferredLoginName = newOrgOwnerEmail
				tt.want.User.LoginNames = []string{newOrgOwnerEmail}
				if human := tt.want.User.GetHuman(); human != nil {
					human.Email.Email = newOrgOwnerEmail
				}
				// details tested in GetUserByID
				tt.want.User.Details = got.User.GetDetails()

				assert.Equal(ttt, tt.want.User, got.User)
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

type userAttrs []userAttr

func (u userAttrs) userIDs() []string {
	ids := make([]string, len(u))
	for i := range u {
		ids[i] = u[i].UserID
	}
	return ids
}

func (u userAttrs) emails() []string {
	emails := make([]string, len(u))
	for i := range u {
		emails[i] = u[i].Username
	}
	return emails
}

type userAttr struct {
	UserID   string
	Username string
	Phone    string
	Changed  *timestamppb.Timestamp
	Details  *object.Details
}

func createUsers(ctx context.Context, orgID string, count int, passwordChangeRequired bool) userAttrs {
	infos := make([]userAttr, count)
	for i := 0; i < count; i++ {
		infos[i] = createUser(ctx, orgID, passwordChangeRequired)
	}
	slices.Reverse(infos)
	return infos
}

func createUser(ctx context.Context, orgID string, passwordChangeRequired bool) userAttr {
	username := integration.Email()
	// used as default country prefix
	phone := integration.Phone()
	resp := Instance.CreateHumanUserVerified(ctx, orgID, username, phone)
	info := userAttr{resp.GetUserId(), username, phone, nil, resp.GetDetails()}
	if passwordChangeRequired {
		details := Instance.SetUserPassword(ctx, resp.GetUserId(), integration.UserPassword, true)
		info.Changed = details.GetChangeDate()
	}
	return info
}

func TestServer_ListUsers(t *testing.T) {
	t.Cleanup(func() {
		_, err := Instance.Client.FeatureV2.ResetInstanceFeatures(IamCTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	})

	orgResp := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	type args struct {
		ctx context.Context
		req *user_v2beta.ListUsersRequest
		dep func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs
	}
	tt := []struct {
		name    string
		args    args
		want    *user_v2beta.ListUsersResponse
		wantErr bool
	}{
		{
			name: "list user by id, no permission machine user",
			args: args{
				UserCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user_v2beta.User{},
			},
		},
		{
			name: "list user by id, no permission human user",
			args: func() args {
				info := createUser(IamCTX, orgResp.OrganizationId, true)
				// create session to get token
				userID := info.UserID
				createResp, err := Instance.Client.SessionV2.CreateSession(IamCTX, &session.CreateSessionRequest{
					Checks: &session.Checks{
						User: &session.CheckUser{
							Search: &session.CheckUser_UserId{UserId: userID},
						},
						Password: &session.CheckPassword{
							Password: integration.UserPassword,
						},
					},
				})
				if err != nil {
					require.NoError(t, err)
				}
				// use token to get ctx
				HumanCTX := integration.WithAuthorizationToken(IamCTX, createResp.GetSessionToken())
				return args{
					HumanCTX,
					&user_v2beta.ListUsersRequest{},
					func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
						return []userAttr{info}
					},
				}
			}(),
			want: &user_v2beta.ListUsersResponse{ // human user should return itself when calling ListUsers() even if it has no permissions
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
								PasswordChangeRequired: true,
								PasswordChanged:        timestamppb.Now(),
							},
						},
					},
				},
			},
		},
		{
			name: "list user by id, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{info}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user by id, passwordChangeRequired, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, true)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{info}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
								PasswordChangeRequired: true,
								PasswordChanged:        timestamppb.Now(),
							},
						},
					},
				},
			},
		},
		{
			name: "list user by id multiple, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery(infos.userIDs()))
					return infos
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user by username, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, UsernameQuery(info.Username))
					return []userAttr{info}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user in emails, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{
					Queries: []*user_v2beta.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, InUserEmailsQuery([]string{info.Username}))
					return []userAttr{info}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user in emails multiple, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))
					return infos
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user in emails no found, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{
					Queries: []*user_v2beta.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
						InUserEmailsQuery([]string{"notfound"}),
					},
				},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					return []userAttr{}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user_v2beta.User{},
			},
		},
		{
			name: "list user phone, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, PhoneQuery(info.Phone))
					return []userAttr{info}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user in emails no found, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{
					Queries: []*user_v2beta.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
						InUserEmailsQuery([]string{"notfound"}),
					},
				},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					return []userAttr{}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user_v2beta.User{},
			},
		},
		{
			name: "list user resourceowner multiple, ok",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())

					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))
					return infos
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user with org query",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					orgRespForOrgTests := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					info := createUser(ctx, orgRespForOrgTests.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgRespForOrgTests.OrganizationId))
					return []userAttr{info, {}}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 2,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					{
						State: user_v2beta.UserState_USER_STATE_ACTIVE,
						Type: &user_v2beta.User_Human{
							Human: &user_v2beta.HumanUser{
								Profile: &user_v2beta.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user_v2beta.Gender_GENDER_MALE.Enum(),
								},
								Email: &user_v2beta.HumanEmail{
									IsVerified: true,
								},
								Phone: &user_v2beta.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
					// this is the admin of the org craated in Instance.CreateOrganization()
					nil,
				},
			},
		},
		{
			name: "list user with wrong org query",
			args: args{
				IamCTX,
				&user_v2beta.ListUsersRequest{},
				func(ctx context.Context, request *user_v2beta.ListUsersRequest) userAttrs {
					orgRespForOrgTests := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					orgRespForOrgTests2 := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					// info := createUser(ctx, orgRespForOrgTests.OrganizationId, false)
					createUser(ctx, orgRespForOrgTests.OrganizationId, false)
					request.Queries = []*user_v2beta.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgRespForOrgTests2.OrganizationId))
					return []userAttr{{}}
				},
			},
			want: &user_v2beta.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user_v2beta.User{
					// this is the admin of the org craated in Instance.CreateOrganization()
					nil,
				},
			},
		},
	}
	for _, f := range permissionCheckV2Settings {
		for _, tc := range tt {
			t.Run(f.TestNamePrependString+tc.name, func(t1 *testing.T) {
				setPermissionCheckV2Flag(t1, f.SetFlag)
				infos := tc.args.dep(IamCTX, tc.args.req)

				// retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, 10*time.Minute)
				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tc.args.ctx, 20*time.Second)
				require.EventuallyWithT(t1, func(ttt *assert.CollectT) {
					got, err := Client.ListUsers(tc.args.ctx, tc.args.req)
					if tc.wantErr {
						require.Error(ttt, err)
						return
					}
					require.NoError(ttt, err)

					// always only give back dependency infos which are required for the response
					require.Len(ttt, tc.want.Result, len(infos))
					// always first check length, otherwise its failed anyway
					if assert.Len(ttt, got.Result, len(tc.want.Result)) {
						// totalResult is unrelated to the tests here so gets carried over, can vary from the count of results due to permissions
						tc.want.Details.TotalResult = got.Details.TotalResult

						// fill in userid and username as it is generated
						for i := range infos {
							if tc.want.Result[i] == nil {
								continue
							}
							tc.want.Result[i].UserId = infos[i].UserID
							tc.want.Result[i].Username = infos[i].Username
							tc.want.Result[i].PreferredLoginName = infos[i].Username
							tc.want.Result[i].LoginNames = []string{infos[i].Username}
							if human := tc.want.Result[i].GetHuman(); human != nil {
								human.Email.Email = infos[i].Username
								human.Phone.Phone = infos[i].Phone
								if tc.want.Result[i].GetHuman().GetPasswordChanged() != nil {
									human.PasswordChanged = infos[i].Changed
								}
							}
							tc.want.Result[i].Details = detailsV2ToV2beta(infos[i].Details)
						}
						for i := range tc.want.Result {
							if tc.want.Result[i] == nil {
								continue
							}
							assert.EqualExportedValues(ttt, got.Result[i], tc.want.Result[i])
						}
					}
					integration.AssertListDetails(ttt, tc.want, got)
				}, retryDuration, tick, "timeout waiting for expected user result")
			})
		}
	}
}

func InUserIDsQuery(ids []string) *user_v2beta.SearchQuery {
	return &user_v2beta.SearchQuery{
		Query: &user_v2beta.SearchQuery_InUserIdsQuery{
			InUserIdsQuery: &user_v2beta.InUserIDQuery{
				UserIds: ids,
			},
		},
	}
}

func InUserEmailsQuery(emails []string) *user_v2beta.SearchQuery {
	return &user_v2beta.SearchQuery{
		Query: &user_v2beta.SearchQuery_InUserEmailsQuery{
			InUserEmailsQuery: &user_v2beta.InUserEmailsQuery{
				UserEmails: emails,
			},
		},
	}
}

func PhoneQuery(number string) *user_v2beta.SearchQuery {
	return &user_v2beta.SearchQuery{
		Query: &user_v2beta.SearchQuery_PhoneQuery{
			PhoneQuery: &user_v2beta.PhoneQuery{
				Number: number,
			},
		},
	}
}

func UsernameQuery(username string) *user_v2beta.SearchQuery {
	return &user_v2beta.SearchQuery{
		Query: &user_v2beta.SearchQuery_UserNameQuery{
			UserNameQuery: &user_v2beta.UserNameQuery{
				UserName: username,
			},
		},
	}
}

func OrganizationIdQuery(resourceowner string) *user_v2beta.SearchQuery {
	return &user_v2beta.SearchQuery{
		Query: &user_v2beta.SearchQuery_OrganizationIdQuery{
			OrganizationIdQuery: &user_v2beta.OrganizationIdQuery{
				OrganizationId: resourceowner,
			},
		},
	}
}
