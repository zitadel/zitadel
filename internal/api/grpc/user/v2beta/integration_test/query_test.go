//go:build integration

package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	object_v2beta "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func detailsV2ToV2beta(obj *object.Details) *object_v2beta.Details {
	return &object_v2beta.Details{
		Sequence:      obj.GetSequence(),
		ChangeDate:    obj.GetChangeDate(),
		ResourceOwner: obj.GetResourceOwner(),
	}
}

func TestServer_GetUserByID(t *testing.T) {
	orgResp := Instance.CreateOrganization(IamCTX, fmt.Sprintf("GetUserByIDOrg-%s", gofakeit.AppName()), gofakeit.Email())
	type args struct {
		ctx context.Context
		req *user.GetUserByIDRequest
		dep func(ctx context.Context, username string, request *user.GetUserByIDRequest) (*userAttr, error)
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
				func(ctx context.Context, username string, request *user.GetUserByIDRequest) (*userAttr, error) {
					return nil, nil
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
				func(ctx context.Context, username string, request *user.GetUserByIDRequest) (*userAttr, error) {
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "user by ID, ok",
			args: args{
				IamCTX,
				&user.GetUserByIDRequest{},
				func(ctx context.Context, username string, request *user.GetUserByIDRequest) (*userAttr, error) {
					resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
					request.UserId = resp.GetUserId()
					return &userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}, nil
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
								Phone:      "+41791234567",
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
				&user.GetUserByIDRequest{},
				func(ctx context.Context, username string, request *user.GetUserByIDRequest) (*userAttr, error) {
					resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
					request.UserId = resp.GetUserId()
					details := Instance.SetUserPassword(ctx, resp.GetUserId(), integration.UserPassword, true)
					return &userAttr{resp.GetUserId(), username, details.GetChangeDate(), resp.GetDetails()}, nil
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
								Phone:      "+41791234567",
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
			username := gofakeit.Email()
			userAttr, err := tt.args.dep(tt.args.ctx, username, tt.args.req)
			require.NoError(t, err)

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

type userAttr struct {
	UserID   string
	Username string
	Changed  *timestamppb.Timestamp
	Details  *object.Details
}

func TestServer_ListUsers(t *testing.T) {
	orgResp := Instance.CreateOrganization(IamCTX, fmt.Sprintf("ListUsersOrg-%s", gofakeit.AppName()), gofakeit.Email())
	userResp := Instance.CreateHumanUserVerified(IamCTX, orgResp.OrganizationId, gofakeit.Email())
	type args struct {
		ctx   context.Context
		count int
		req   *user.ListUsersRequest
		dep   func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error)
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
				0,
				&user.ListUsersRequest{},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					request.Queries = append(request.Queries, InUserIDsQuery([]string{userResp.UserId}))
					return []userAttr{}, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
				1,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					userIDs := make([]string, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						userIDs[i] = resp.GetUserId()
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
					}
					request.Queries = append(request.Queries, InUserIDsQuery(userIDs))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
				1,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					userIDs := make([]string, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						userIDs[i] = resp.GetUserId()
						details := Instance.SetUserPassword(ctx, resp.GetUserId(), integration.UserPassword, true)
						infos[i] = userAttr{resp.GetUserId(), username, details.GetChangeDate(), resp.GetDetails()}
					}
					request.Queries = append(request.Queries, InUserIDsQuery(userIDs))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
				3,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					userIDs := make([]string, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						userIDs[i] = resp.GetUserId()
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
					}
					request.Queries = append(request.Queries, InUserIDsQuery(userIDs))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
				1,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					userIDs := make([]string, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						userIDs[i] = resp.GetUserId()
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
						request.Queries = append(request.Queries, UsernameQuery(username))
					}
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
				1,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
					}
					request.Queries = append(request.Queries, InUserEmailsQuery(usernames))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
				3,
				&user.ListUsersRequest{
					Queries: []*user.SearchQuery{
						OrganizationIdQuery(orgResp.OrganizationId),
					},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					infos := make([]userAttr, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
					}
					request.Queries = append(request.Queries, InUserEmailsQuery(usernames))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
				3,
				&user.ListUsersRequest{Queries: []*user.SearchQuery{
					OrganizationIdQuery(orgResp.OrganizationId),
					InUserEmailsQuery([]string{"notfound"}),
				},
				},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					return []userAttr{}, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
				3,
				&user.ListUsersRequest{},
				func(ctx context.Context, usernames []string, request *user.ListUsersRequest) ([]userAttr, error) {
					orgResp := Instance.CreateOrganization(ctx, fmt.Sprintf("ListUsersResourceowner-%s", gofakeit.AppName()), gofakeit.Email())

					infos := make([]userAttr, len(usernames))
					for i, username := range usernames {
						resp := Instance.CreateHumanUserVerified(ctx, orgResp.OrganizationId, username)
						infos[i] = userAttr{resp.GetUserId(), username, nil, resp.GetDetails()}
					}
					request.Queries = append(request.Queries, OrganizationIdQuery(orgResp.OrganizationId))
					request.Queries = append(request.Queries, InUserEmailsQuery(usernames))
					return infos, nil
				},
			},
			want: &user.ListUsersResponse{
				Details: &object_v2beta.ListDetails{
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
									Phone:      "+41791234567",
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
			usernames := make([]string, tt.args.count)
			for i := 0; i < tt.args.count; i++ {
				usernames[i] = gofakeit.Email()
			}
			infos, err := tt.args.dep(tt.args.ctx, usernames, tt.args.req)
			require.NoError(t, err)

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
							if tt.want.Result[i].GetHuman().GetPasswordChanged() != nil {
								human.PasswordChanged = infos[i].Changed
							}
						}
						tt.want.Result[i].Details = detailsV2ToV2beta(infos[i].Details)
					}
					for i := range tt.want.Result {
						assert.Contains(ttt, got.Result, tt.want.Result[i])
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
