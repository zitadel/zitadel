//go:build integration

package user_test

import (
	"context"
	"encoding/base64"
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
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	v2 "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

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
		req *user.GetUserByIDRequest
		dep func(ctx context.Context, request *user.GetUserByIDRequest) *userAttr
	}
	tests := []struct {
		name    string
		args    args
		want    *user.GetUserByIDResponse
		wantErr bool
	}{
		{
			name: "user by ID, no id provided",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{
					UserId: "",
				},
				func(ctx context.Context, request *user.GetUserByIDRequest) *userAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "user by ID, not found",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{
					UserId: "unknown",
				},
				func(ctx context.Context, request *user.GetUserByIDRequest) *userAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "user by ID, ok",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{},
				func(ctx context.Context, request *user.GetUserByIDRequest) *userAttr {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.UserId = info.UserID
					return &info
				},
			},
			want: &user.GetUserByIDResponse{
				User: &user.User{
					State:              user.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user.User_Human{
						Human: &user.HumanUser{
							Profile: &user.HumanProfile{
								GivenName:         "Mickey",
								FamilyName:        "Mouse",
								NickName:          gu.Ptr("Mickey"),
								DisplayName:       gu.Ptr("Mickey Mouse"),
								PreferredLanguage: gu.Ptr("nl"),
								Gender:            user.Gender_GENDER_MALE.Enum(),
								AvatarUrl:         "",
							},
							Email: &user.HumanEmail{
								IsVerified: true,
							},
							Phone: &user.HumanPhone{
								IsVerified: true,
							},
						},
					},
				},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					CreationDate:  timestamppb.Now(),
					ResourceOwner: orgResp.OrganizationId,
				},
			},
		},
		{
			name: "user by ID, passwordChangeRequired, ok",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{},
				func(ctx context.Context, request *user.GetUserByIDRequest) *userAttr {
					info := createUser(ctx, orgResp.OrganizationId, true)
					request.UserId = info.UserID
					return &info
				},
			},
			want: &user.GetUserByIDResponse{
				User: &user.User{
					State:              user.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user.User_Human{
						Human: &user.HumanUser{
							Profile: &user.HumanProfile{
								GivenName:         "Mickey",
								FamilyName:        "Mouse",
								NickName:          gu.Ptr("Mickey"),
								DisplayName:       gu.Ptr("Mickey Mouse"),
								PreferredLanguage: gu.Ptr("nl"),
								Gender:            user.Gender_GENDER_MALE.Enum(),
								AvatarUrl:         "",
							},
							Email: &user.HumanEmail{
								IsVerified: true,
							},
							Phone: &user.HumanPhone{
								IsVerified: true,
							},
							PasswordChangeRequired: true,
							PasswordChanged:        timestamppb.Now(),
						},
					},
				},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					CreationDate:  timestamppb.Now(),
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

				tt.want.User.Details = userAttr.Details
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
				assert.EqualExportedValues(ttt, tt.want.User, got.User)
				integration.AssertDetails(ttt, tt.want, got)
			}, retryDuration, tick)
		})
	}
}

func TestServer_GetUserByID_Permission(t *testing.T) {
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	newUserID := newOrg.CreatedAdmins[0].GetUserId()
	type args struct {
		ctx context.Context
		req *user.GetUserByIDRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.GetUserByIDResponse
		wantErr bool
	}{
		{
			name: "System, ok",
			args: args{
				SystemCTX,
				&user.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			want: &user.GetUserByIDResponse{
				User: &user.User{
					State:              user.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user.User_Human{
						Human: &user.HumanUser{
							Profile: &user.HumanProfile{
								GivenName:         "firstname",
								FamilyName:        "lastname",
								NickName:          gu.Ptr(""),
								DisplayName:       gu.Ptr("firstname lastname"),
								PreferredLanguage: gu.Ptr("und"),
								Gender:            user.Gender_GENDER_UNSPECIFIED.Enum(),
								AvatarUrl:         "",
							},
							Email: &user.HumanEmail{
								Email: newOrgOwnerEmail,
							},
							Phone: &user.HumanPhone{},
						},
					},
				},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					CreationDate:  timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Instance, ok",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			want: &user.GetUserByIDResponse{
				User: &user.User{
					State:              user.UserState_USER_STATE_ACTIVE,
					Username:           "",
					LoginNames:         nil,
					PreferredLoginName: "",
					Type: &user.User_Human{
						Human: &user.HumanUser{
							Profile: &user.HumanProfile{
								GivenName:         "firstname",
								FamilyName:        "lastname",
								NickName:          gu.Ptr(""),
								DisplayName:       gu.Ptr("firstname lastname"),
								PreferredLanguage: gu.Ptr("und"),
								Gender:            user.Gender_GENDER_UNSPECIFIED.Enum(),
								AvatarUrl:         "",
							},
							Email: &user.HumanEmail{
								Email: newOrgOwnerEmail,
							},
							Phone: &user.HumanPhone{},
						},
					},
				},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					CreationDate:  timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Org, error",
			args: args{
				CTX,
				&user.GetUserByIDRequest{
					UserId: newUserID,
				},
			},
			wantErr: true,
		},
		{
			name: "User, error",
			args: args{
				UserCTX,
				&user.GetUserByIDRequest{
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
	return createUserWithUserName(ctx, username, orgID, passwordChangeRequired)
}

func createUserWithUserName(ctx context.Context, username string, orgID string, passwordChangeRequired bool) userAttr {
	// used as default country prefix
	phone := integration.Phone()
	resp := Instance.CreateHumanUserVerified(ctx, orgID, username, phone)
	info := userAttr{resp.GetUserId(), username, phone, nil, resp.GetDetails()}
	// as the change date of the creation is the creation date
	resp.Details.CreationDate = resp.GetDetails().GetChangeDate()
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
		req *user.ListUsersRequest
		dep func(ctx context.Context, request *user.ListUsersRequest) userAttrs
	}
	tt := []struct {
		name    string
		args    args
		want    *user.ListUsersResponse
		wantErr bool
	}{
		{
			name: "list user by id, no permission machine user",
			args: args{
				UserCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user.User{},
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
					&user.ListUsersRequest{},
					func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
						return []userAttr{info}
					},
				}
			}(),
			want: &user.ListUsersResponse{ // human user should return itself when calling ListUsers() even if it has no permissions
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{info}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, true)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery([]string{info.UserID}))
					return []userAttr{info}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
			name: "list user by id and meta key multiple, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserIDsQuery(infos.userIDs()))

					Instance.SetUserMetadata(ctx, infos[0].UserID, "my meta", "my value 1")
					Instance.SetUserMetadata(ctx, infos[1].UserID, "my meta 2", "my value 3")
					Instance.SetUserMetadata(ctx, infos[2].UserID, "my meta", "my value 2")

					request.Queries = append(request.Queries, MetadataKeyContainsQuery("my meta"))
					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, UsernameQuery(info.Username))

					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return []userAttr{info}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery([]string{info.Username}))
					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return []userAttr{info}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "list user by emails and meta value multiple, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))

					Instance.SetUserMetadata(ctx, infos[0].UserID, "my meta 1", "my value")
					Instance.SetUserMetadata(ctx, infos[0].UserID, "my meta 2", "my value")
					Instance.SetUserMetadata(ctx, infos[1].UserID, "my meta 2", "my value")
					Instance.SetUserMetadata(ctx, infos[2].UserID, "my meta", "my value")

					request.Queries = append(request.Queries, MetadataValueQuery("my value"))
					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE

					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
						InUserEmailsQuery([]string{"notfound"}),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					return []userAttr{}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user.User{},
			},
		},
		{
			name: "list user phone, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, PhoneQuery(info.Phone))
					return []userAttr{info}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
						InUserEmailsQuery([]string{"notfound"}),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					return []userAttr{}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result:        []*user.User{},
			},
		},
		{
			name: "list user resourceowner multiple, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())

					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))
					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
									IsVerified: true,
								},
							},
						},
					}, {
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					orgRespForOrgTests := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					info := createUser(ctx, orgRespForOrgTests.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgRespForOrgTests.OrganizationId))
					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return []userAttr{info, {}}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 2,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result: []*user.User{
					{
						State: user.UserState_USER_STATE_ACTIVE,
						Type: &user.User_Human{
							Human: &user.HumanUser{
								Profile: &user.HumanProfile{
									GivenName:         "Mickey",
									FamilyName:        "Mouse",
									NickName:          gu.Ptr("Mickey"),
									DisplayName:       gu.Ptr("Mickey Mouse"),
									PreferredLanguage: gu.Ptr("nl"),
									Gender:            user.Gender_GENDER_MALE.Enum(),
								},
								Email: &user.HumanEmail{
									IsVerified: true,
								},
								Phone: &user.HumanPhone{
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
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					orgRespForOrgTests := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					orgRespForOrgTests2 := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
					createUser(ctx, orgRespForOrgTests.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgRespForOrgTests2.OrganizationId))
					return []userAttr{{}}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: 0,
				Result: []*user.User{
					// this is the admin of the org craated in Instance.CreateOrganization()
					nil,
				},
			},
		},
		{
			name: "when no users matching meta key should return empty list",
			args: args{
				IamCTX,
				&user.ListUsersRequest{},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = []*user.SearchQuery{}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, MetadataKeyContainsQuery("some non-existent meta"))

					request.SortingColumn = user.UserFieldName_USER_FIELD_NAME_CREATION_DATE
					return []userAttr{}
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
					Timestamp:   timestamppb.Now(),
				},
				SortingColumn: user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
				Result:        []*user.User{},
			},
		},
	}
	for _, f := range permissionCheckV2Settings {
		for _, tc := range tt {
			t.Run(f.TestNamePrependString+tc.name, func(t1 *testing.T) {
				setPermissionCheckV2Flag(t, f.SetFlag)
				infos := tc.args.dep(IamCTX, tc.args.req)

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
					if assert.Len(ttt, got.Result, len(tc.want.Result)) {
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
							tc.want.Result[i].Details = infos[i].Details
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

func TestServer_SystemUsers_ListUsers(t *testing.T) {
	defer func() {
		_, err := Instance.Client.FeatureV2.ResetInstanceFeatures(IamCTX, &feature.ResetInstanceFeaturesRequest{})
		require.NoError(t, err)
	}()

	org1 := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	org2 := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), "org2@zitadel.com")
	org3 := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	_ = createUserWithUserName(IamCTX, "Test_SystemUsers_ListUser1@zitadel.com", org1.OrganizationId, false)
	_ = createUserWithUserName(IamCTX, "Test_SystemUsers_ListUser2@zitadel.com", org2.OrganizationId, false)
	_ = createUserWithUserName(IamCTX, "Test_SystemUsers_ListUser3@zitadel.com", org3.OrganizationId, false)

	tests := []struct {
		name                       string
		ctx                        context.Context
		req                        *user.ListUsersRequest
		expectedFoundUsernames     []string
		checkNumberOfUsersReturned bool
	}{
		{
			name: "list users with neccessary permissions",
			ctx:  SystemCTX,
			req:  &user.ListUsersRequest{},
			// the number of users returned will vary from test run to test run,
			// so just check the system user gets back users from different orgs whcih it is not a memeber of
			checkNumberOfUsersReturned: false,
			expectedFoundUsernames:     []string{"Test_SystemUsers_ListUser1@zitadel.com", "Test_SystemUsers_ListUser2@zitadel.com", "Test_SystemUsers_ListUser3@zitadel.com"},
		},
		{
			name: "list users without neccessary permissions",
			ctx:  SystemUserWithNoPermissionsCTX,
			req:  &user.ListUsersRequest{},
			// check no users returned
			checkNumberOfUsersReturned: true,
		},
		{
			name: "list users with neccessary permissions specifying org",
			req: &user.ListUsersRequest{
				Queries: []*user.SearchQuery{OrganizationIdQuery(org2.OrganizationId)},
			},
			ctx:                        SystemCTX,
			expectedFoundUsernames:     []string{"Test_SystemUsers_ListUser2@zitadel.com", "org2@zitadel.com"},
			checkNumberOfUsersReturned: true,
		},
		{
			name: "list users without neccessary permissions specifying org",
			req: &user.ListUsersRequest{
				Queries: []*user.SearchQuery{OrganizationIdQuery(org2.OrganizationId)},
			},
			ctx: SystemUserWithNoPermissionsCTX,
			// check no users returned
			checkNumberOfUsersReturned: true,
		},
	}

	for _, f := range permissionCheckV2Settings {
		f := f
		for _, tt := range tests {
			t.Run(f.TestNamePrependString+tt.name, func(t *testing.T) {
				setPermissionCheckV2Flag(t, f.SetFlag)

				retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.ctx, 1*time.Minute)
				require.EventuallyWithT(t, func(ttt *assert.CollectT) {
					got, err := Client.ListUsers(tt.ctx, tt.req)
					require.NoError(ttt, err)

					if tt.checkNumberOfUsersReturned {
						require.Equal(t, len(tt.expectedFoundUsernames), len(got.Result))
					}

					if tt.expectedFoundUsernames != nil {
						for _, user := range got.Result {
							for i, username := range tt.expectedFoundUsernames {
								if username == user.Username {
									tt.expectedFoundUsernames = tt.expectedFoundUsernames[i+1:]
									break
								}
							}
							if len(tt.expectedFoundUsernames) == 0 {
								return
							}
						}
						require.FailNow(t, "unable to find all users with specified usernames")
					}
				}, retryDuration, tick, "timeout waiting for expected user result")
			})
		}
	}
}

func InUserIDsQuery(ids []string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_InUserIdsQuery{
			InUserIdsQuery: &user.InUserIDQuery{
				UserIds: ids,
			},
		},
	}
}

func InUserEmailsQuery(emails []string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_InUserEmailsQuery{
			InUserEmailsQuery: &user.InUserEmailsQuery{
				UserEmails: emails,
			},
		},
	}
}

func PhoneQuery(number string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_PhoneQuery{
			PhoneQuery: &user.PhoneQuery{
				Number: number,
			},
		},
	}
}

func UsernameQuery(username string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_UserNameQuery{
			UserNameQuery: &user.UserNameQuery{
				UserName: username,
			},
		},
	}
}

func OrganizationIdQuery(resourceowner string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_OrganizationIdQuery{
			OrganizationIdQuery: &user.OrganizationIdQuery{
				OrganizationId: resourceowner,
			},
		},
	}
}

func OrQuery(queries []*user.SearchQuery) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_OrQuery{
			OrQuery: &user.OrQuery{
				Queries: queries,
			},
		},
	}
}

func MetadataKeyContainsQuery(metadataKey string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_MetadataKeyFilter{
			MetadataKeyFilter: &v2.MetadataKeyFilter{
				Key:    metadataKey,
				Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH},
		},
	}
}

func MetakeyEqualsQuery(metaKey string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_MetadataKeyFilter{
			MetadataKeyFilter: &v2.MetadataKeyFilter{
				Key:    metaKey,
				Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS},
		},
	}
}

func MetadataValueQuery(metaValue string) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_MetadataValueFilter{
			MetadataValueFilter: &v2.MetadataValueFilter{
				Value:  []byte(base64.StdEncoding.EncodeToString([]byte(metaValue))),
				Method: filter.ByteFilterMethod_BYTE_FILTER_METHOD_EQUALS,
			},
		},
	}
}
