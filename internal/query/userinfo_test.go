package query

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed testdata/userinfo_not_found.json
	testdataUserInfoNotFound string
	//go:embed testdata/userinfo_human_no_md.json
	testdataUserInfoHumanNoMD string
	//go:embed testdata/userinfo_human.json
	testdataUserInfoHuman string
	//go:embed testdata/userinfo_human_grants.json
	testdataUserInfoHumanGrants string
	//go:embed testdata/userinfo_machine.json
	testdataUserInfoMachine string

	// timeLocation does a single parse of the testdata and extracts a time.Location,
	// so that it may be used during test assertion.
	// Because depending on the environment json.Unmarshal parses
	// the 00:00 timezones differently.
	// On my local system is parses to an empty timezone with 0 offset,
	// but in github action is pares into UTC.
	timeLocation = func() *time.Location {
		referenceInfo := new(OIDCUserInfo)
		err := json.Unmarshal([]byte(testdataUserInfoHumanNoMD), referenceInfo)
		if err != nil {
			panic(err)
		}
		return referenceInfo.User.CreationDate.Location()
	}()
)

func TestQueries_GetOIDCUserInfo(t *testing.T) {
	expQuery := regexp.QuoteMeta(oidcUserInfoQuery)
	type args struct {
		userID       string
		roleAudience []string
	}
	tests := []struct {
		name    string
		args    args
		mock    sqlExpectation
		want    *OIDCUserInfo
		wantErr error
	}{
		{
			name: "query error",
			args: args{
				userID: "231965491734773762",
			},
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "231965491734773762", "instanceID", nil),
			wantErr: sql.ErrConnDone,
		},
		{
			name: "user not found",
			args: args{
				userID: "231965491734773762",
			},
			mock:    mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoNotFound}, "231965491734773762", "instanceID", nil),
			wantErr: zerrors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.User.NotFound"),
		},
		{
			name: "human without metadata",
			args: args{
				userID: "231965491734773762",
			},
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoHumanNoMD}, "231965491734773762", "instanceID", nil),
			want: &OIDCUserInfo{
				User: &User{
					ID:                 "231965491734773762",
					CreationDate:       time.Date(2023, time.September, 15, 6, 10, 7, 434142000, timeLocation),
					ChangeDate:         time.Date(2023, time.November, 14, 13, 27, 2, 72318000, timeLocation),
					Sequence:           1148,
					State:              1,
					ResourceOwner:      "231848297847848962",
					Username:           "tim+tesmail@zitadel.com",
					PreferredLoginName: "tim+tesmail@zitadel.com@demo.localhost",
					Human: &Human{
						FirstName:         "Tim",
						LastName:          "Mohlmann",
						NickName:          "muhlemmer",
						DisplayName:       "Tim Mohlmann",
						AvatarKey:         "",
						PreferredLanguage: language.English,
						Gender:            domain.GenderMale,
						Email:             "tim+tesmail@zitadel.com",
						IsEmailVerified:   true,
						Phone:             "+40123456789",
						IsPhoneVerified:   false,
					},
					Machine: nil,
				},
				Org: &UserInfoOrg{
					ID:            "231848297847848962",
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: nil,
			},
		},
		{
			name: "human with metadata",
			args: args{
				userID: "231965491734773762",
			},
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoHuman}, "231965491734773762", "instanceID", nil),
			want: &OIDCUserInfo{
				User: &User{
					ID:                 "231965491734773762",
					CreationDate:       time.Date(2023, time.September, 15, 6, 10, 7, 434142000, timeLocation),
					ChangeDate:         time.Date(2023, time.November, 14, 13, 27, 2, 72318000, timeLocation),
					Sequence:           1148,
					State:              1,
					ResourceOwner:      "231848297847848962",
					Username:           "tim+tesmail@zitadel.com",
					PreferredLoginName: "tim+tesmail@zitadel.com@demo.localhost",
					Human: &Human{
						FirstName:         "Tim",
						LastName:          "Mohlmann",
						NickName:          "muhlemmer",
						DisplayName:       "Tim Mohlmann",
						AvatarKey:         "",
						PreferredLanguage: language.English,
						Gender:            domain.GenderMale,
						Email:             "tim+tesmail@zitadel.com",
						IsEmailVerified:   true,
						Phone:             "+40123456789",
						IsPhoneVerified:   false,
					},
					Machine: nil,
				},
				Org: &UserInfoOrg{
					ID:            "231848297847848962",
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: []UserMetadata{
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 26, 3, 553702000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 26, 3, 553702000, timeLocation),
						Sequence:      1147,
						ResourceOwner: "231848297847848962",
						Key:           "bar",
						Value:         []byte("foo"),
					},
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 25, 57, 171368000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 25, 57, 171368000, timeLocation),
						Sequence:      1146,
						ResourceOwner: "231848297847848962",
						Key:           "foo",
						Value:         []byte("bar"),
					},
				},
			},
		},
		{
			name: "human with metadata and grants",
			args: args{
				userID:       "231965491734773762",
				roleAudience: []string{"236645808328409090", "240762134579904514"},
			},
			mock: mockQuery(expQuery,
				[]string{"json_build_object"},
				[]driver.Value{testdataUserInfoHumanGrants},
				"231965491734773762", "instanceID", database.TextArray[string]{"236645808328409090", "240762134579904514"},
			),
			want: &OIDCUserInfo{
				User: &User{
					ID:                 "231965491734773762",
					CreationDate:       time.Date(2023, time.September, 15, 6, 10, 7, 434142000, timeLocation),
					ChangeDate:         time.Date(2023, time.November, 14, 13, 27, 2, 72318000, timeLocation),
					Sequence:           1148,
					State:              1,
					ResourceOwner:      "231848297847848962",
					Username:           "tim+tesmail@zitadel.com",
					PreferredLoginName: "tim+tesmail@zitadel.com@demo.localhost",
					Human: &Human{
						FirstName:         "Tim",
						LastName:          "Mohlmann",
						NickName:          "muhlemmer",
						DisplayName:       "Tim Mohlmann",
						AvatarKey:         "",
						PreferredLanguage: language.English,
						Gender:            domain.GenderMale,
						Email:             "tim+tesmail@zitadel.com",
						IsEmailVerified:   true,
						Phone:             "+40123456789",
						IsPhoneVerified:   false,
					},
					Machine: nil,
				},
				Org: &UserInfoOrg{
					ID:            "231848297847848962",
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: []UserMetadata{
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 26, 3, 553702000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 26, 3, 553702000, timeLocation),
						Sequence:      1147,
						ResourceOwner: "231848297847848962",
						Key:           "bar",
						Value:         []byte("foo"),
					},
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 25, 57, 171368000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 25, 57, 171368000, timeLocation),
						Sequence:      1146,
						ResourceOwner: "231848297847848962",
						Key:           "foo",
						Value:         []byte("bar"),
					},
				},
				UserGrants: []UserGrant{
					{
						ID:           "240749256523120642",
						GrantID:      "",
						State:        1,
						CreationDate: time.Date(2023, time.November, 14, 20, 28, 59, 168208000, timeLocation),
						ChangeDate:   time.Date(2023, time.November, 14, 20, 50, 58, 822391000, timeLocation),
						Sequence:     2,
						UserID:       "231965491734773762",
						Roles: []string{
							"role1",
							"role2",
						},
						ResourceOwner:     "231848297847848962",
						ProjectID:         "236645808328409090",
						OrgName:           "demo",
						OrgPrimaryDomain:  "demo.localhost",
						ProjectName:       "tests",
						UserResourceOwner: "231848297847848962",
					},
					{
						ID:           "240762315572510722",
						GrantID:      "",
						State:        1,
						CreationDate: time.Date(2023, time.November, 14, 22, 38, 42, 967317000, timeLocation),
						ChangeDate:   time.Date(2023, time.November, 14, 22, 38, 42, 967317000, timeLocation),
						Sequence:     1,
						UserID:       "231965491734773762",
						Roles: []string{
							"role3",
							"role4",
						},
						ResourceOwner:     "231848297847848962",
						ProjectID:         "240762134579904514",
						OrgName:           "demo",
						OrgPrimaryDomain:  "demo.localhost",
						ProjectName:       "tests2",
						UserResourceOwner: "231848297847848962",
					},
				},
			},
		},
		{
			name: "machine with metadata",
			args: args{
				userID: "240707570677841922",
			},
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoMachine}, "240707570677841922", "instanceID", nil),
			want: &OIDCUserInfo{
				User: &User{
					ID:                 "240707570677841922",
					CreationDate:       time.Date(2023, time.November, 14, 13, 34, 52, 473732000, timeLocation),
					ChangeDate:         time.Date(2023, time.November, 14, 13, 35, 2, 861342000, timeLocation),
					Sequence:           2,
					State:              1,
					ResourceOwner:      "231848297847848962",
					Username:           "tests",
					PreferredLoginName: "tests@demo.localhost",
					Human:              nil,
					Machine: &Machine{
						Name:        "tests",
						Description: "My test service user",
					},
				},
				Org: &UserInfoOrg{
					ID:            "231848297847848962",
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: []UserMetadata{
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 35, 30, 126849000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 35, 30, 126849000, timeLocation),
						Sequence:      3,
						ResourceOwner: "231848297847848962",
						Key:           "first",
						Value:         []byte("Hello World!"),
					},
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 35, 44, 28343000, timeLocation),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 35, 44, 28343000, timeLocation),
						Sequence:      4,
						ResourceOwner: "231848297847848962",
						Key:           "second",
						Value:         []byte("Bye World!"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB:       db,
						Database: &prepareDB{},
					},
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")

				got, err := q.GetOIDCUserInfo(ctx, tt.args.userID, tt.args.roleAudience)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
