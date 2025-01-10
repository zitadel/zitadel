//go:build integration

package user_test

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_GetUserByID(t *testing.T) {
	orgResp := Instance.CreateOrganization(IamCTX, fmt.Sprintf("GetUserByIDOrg-%s", gofakeit.AppName()), gofakeit.Email())
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
	newOrgOwnerEmail := gofakeit.Email()
	newOrg := Instance.CreateOrganization(IamCTX, fmt.Sprintf("GetHuman-%s", gofakeit.AppName()), newOrgOwnerEmail)
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
	username := gofakeit.Email()
	// used as default country prefix
	phone := "+41" + gofakeit.Phone()
	resp := Instance.CreateHumanUserVerified(ctx, orgID, username, phone)
	info := userAttr{resp.GetUserId(), username, phone, nil, resp.GetDetails()}
	if passwordChangeRequired {
		details := Instance.SetUserPassword(ctx, resp.GetUserId(), integration.UserPassword, true)
		info.Changed = details.GetChangeDate()
	}
	return info
}

func TestServer_ListUsers(t *testing.T) {
	orgResp := Instance.CreateOrganization(IamCTX, fmt.Sprintf("ListUsersOrg-%s", gofakeit.AppName()), gofakeit.Email())
	type args struct {
		ctx context.Context
		req *user.ListUsersRequest
		dep func(ctx context.Context, request *user.ListUsersRequest) userAttrs
	}
	tests := []struct {
		name    string
		args    args
		want    *user.ListUsersResponse
		wantErr bool
	}{
		{
			name: "list user by id, no permission",
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
			name: "list user by id, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
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
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, true)
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
			name: "list user by id multiple, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = append(request.Queries, InUserIDsQuery(infos.userIDs()))
					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
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
					},
				},
			},
		},
		{
			name: "list user by username, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = append(request.Queries, UsernameQuery(info.Username))
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
			name: "list user in emails, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
					request.Queries = append(request.Queries, InUserEmailsQuery([]string{info.Username}))
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
			name: "list user in emails multiple, ok",
			args: args{
				IamCTX,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))
					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
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
				&user.ListUsersRequest{Queries: []*user.SearchQuery{
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
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, request *user.ListUsersRequest) userAttrs {
					info := createUser(ctx, orgResp.OrganizationId, false)
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
				&user.ListUsersRequest{Queries: []*user.SearchQuery{
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
					orgResp := Instance.CreateOrganization(ctx, fmt.Sprintf("ListUsersResourceowner-%s", gofakeit.AppName()), gofakeit.Email())

					infos := createUsers(ctx, orgResp.OrganizationId, 3, false)
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(infos.emails()))
					return infos
				},
			},
			want: &user.ListUsersResponse{
				Details: &object.ListDetails{
					TotalResult: 3,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos := tt.args.dep(IamCTX, tt.args.req)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListUsers(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				// always only give back dependency infos which are required for the response
				require.Len(ttt, tt.want.Result, len(infos))
				// always first check length, otherwise its failed anyway
				if assert.Len(ttt, got.Result, len(tt.want.Result)) {
					// totalResult is unrelated to the tests here so gets carried over, can vary from the count of results due to permissions
					tt.want.Details.TotalResult = got.Details.TotalResult

					// fill in userid and username as it is generated
					for i := range infos {
						tt.want.Result[i].UserId = infos[i].UserID
						tt.want.Result[i].Username = infos[i].Username
						tt.want.Result[i].PreferredLoginName = infos[i].Username
						tt.want.Result[i].LoginNames = []string{infos[i].Username}
						if human := tt.want.Result[i].GetHuman(); human != nil {
							human.Email.Email = infos[i].Username
							human.Phone.Phone = infos[i].Phone
							if tt.want.Result[i].GetHuman().GetPasswordChanged() != nil {
								human.PasswordChanged = infos[i].Changed
							}
						}
						tt.want.Result[i].Details = infos[i].Details
					}
					for i := range tt.want.Result {
						assert.EqualExportedValues(ttt, got.Result[i], tt.want.Result[i])
					}
				}
				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func InUserIDsQuery(ids []string) *user.SearchQuery {
	return &user.SearchQuery{Query: &user.SearchQuery_InUserIdsQuery{
		InUserIdsQuery: &user.InUserIDQuery{
			UserIds: ids,
		},
	},
	}
}

func InUserEmailsQuery(emails []string) *user.SearchQuery {
	return &user.SearchQuery{Query: &user.SearchQuery_InUserEmailsQuery{
		InUserEmailsQuery: &user.InUserEmailsQuery{
			UserEmails: emails,
		},
	},
	}
}

func PhoneQuery(number string) *user.SearchQuery {
	return &user.SearchQuery{Query: &user.SearchQuery_PhoneQuery{
		PhoneQuery: &user.PhoneQuery{
			Number: number,
		},
	},
	}
}

func UsernameQuery(username string) *user.SearchQuery {
	return &user.SearchQuery{Query: &user.SearchQuery_UserNameQuery{
		UserNameQuery: &user.UserNameQuery{
			UserName: username,
		},
	},
	}
}

func OrganizationIdQuery(resourceowner string) *user.SearchQuery {
	return &user.SearchQuery{Query: &user.SearchQuery_OrganizationIdQuery{
		OrganizationIdQuery: &user.OrganizationIdQuery{
			OrganizationId: resourceowner,
		},
	},
	}
}
