package query

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/errors"
)

var (
	//go:embed testdata/userinfo_not_found.json
	testdataUserInfoNotFound string
	//go:embed testdata/userinfo_human_no_md.json
	testdataUserInfoHumanNoMD string
	//go:embed testdata/userinfo_human.json
	testdataUserInfoHuman string
	//go:embed testdata/userinfo_machine.json
	testdataUserInfoMachine string
)

func TestQueries_GetOIDCUserInfo(t *testing.T) {
	expQuery := regexp.QuoteMeta(oidcUserInfoQuery)
	type args struct {
		userID string
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
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "231965491734773762", "instanceID"),
			wantErr: sql.ErrConnDone,
		},
		{
			name: "unmarshal error",
			args: args{
				userID: "231965491734773762",
			},
			mock:    mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{`~~~`}, "231965491734773762", "instanceID"),
			wantErr: errors.ThrowInternal(nil, "QUERY-Vohs6", "Errors.Internal"),
		},
		{
			name: "user not found",
			args: args{
				userID: "231965491734773762",
			},
			mock:    mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoNotFound}, "231965491734773762", "instanceID"),
			wantErr: errors.ThrowNotFound(nil, "QUERY-ahs4S", "Errors.User.NotFound"),
		},
		{
			name: "human without metadata",
			args: args{
				userID: "231965491734773762",
			},
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoHumanNoMD}, "231965491734773762", "instanceID"),
			want: &OIDCUserInfo{
				User: &User{
					ID:            "231965491734773762",
					CreationDate:  time.Date(2023, time.September, 15, 6, 10, 7, 434142000, time.FixedZone("", 0)),
					ChangeDate:    time.Date(2023, time.November, 14, 13, 27, 2, 72318000, time.FixedZone("", 0)),
					Sequence:      1148,
					State:         1,
					ResourceOwner: "231848297847848962",
					Username:      "tim+tesmail@zitadel.com",
					Human: &Human{
						FirstName:       "Tim",
						LastName:        "Mohlmann",
						NickName:        "muhlemmer",
						DisplayName:     "Tim Mohlmann",
						AvatarKey:       "",
						Email:           "tim+tesmail@zitadel.com",
						IsEmailVerified: true,
						Phone:           "+40123456789",
						IsPhoneVerified: false,
					},
					Machine: nil,
				},
				Org: &userInfoOrg{
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
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoHuman}, "231965491734773762", "instanceID"),
			want: &OIDCUserInfo{
				User: &User{
					ID:            "231965491734773762",
					CreationDate:  time.Date(2023, time.September, 15, 6, 10, 7, 434142000, time.FixedZone("", 0)),
					ChangeDate:    time.Date(2023, time.November, 14, 13, 27, 2, 72318000, time.FixedZone("", 0)),
					Sequence:      1148,
					State:         1,
					ResourceOwner: "231848297847848962",
					Username:      "tim+tesmail@zitadel.com",
					Human: &Human{
						FirstName:       "Tim",
						LastName:        "Mohlmann",
						NickName:        "muhlemmer",
						DisplayName:     "Tim Mohlmann",
						AvatarKey:       "",
						Email:           "tim+tesmail@zitadel.com",
						IsEmailVerified: true,
						Phone:           "+40123456789",
						IsPhoneVerified: false,
					},
					Machine: nil,
				},
				Org: &userInfoOrg{
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: []UserMetadata{
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 26, 3, 553702000, time.FixedZone("", 0)),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 26, 3, 553702000, time.FixedZone("", 0)),
						Sequence:      1147,
						ResourceOwner: "231848297847848962",
						Key:           "bar",
						Value:         []byte("foo"),
					},
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 25, 57, 171368000, time.FixedZone("", 0)),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 25, 57, 171368000, time.FixedZone("", 0)),
						Sequence:      1146,
						ResourceOwner: "231848297847848962",
						Key:           "foo",
						Value:         []byte("bar"),
					},
				},
			},
		},
		{
			name: "machine with metadata",
			args: args{
				userID: "240707570677841922",
			},
			mock: mockQuery(expQuery, []string{"json_build_object"}, []driver.Value{testdataUserInfoMachine}, "240707570677841922", "instanceID"),
			want: &OIDCUserInfo{
				User: &User{
					ID:            "240707570677841922",
					CreationDate:  time.Date(2023, time.November, 14, 13, 34, 52, 473732000, time.FixedZone("", 0)),
					ChangeDate:    time.Date(2023, time.November, 14, 13, 35, 2, 861342000, time.FixedZone("", 0)),
					Sequence:      2,
					State:         1,
					ResourceOwner: "231848297847848962",
					Username:      "tests",
					Human:         nil,
					Machine: &Machine{
						Name:        "tests",
						Description: "My test service user",
					},
				},
				Org: &userInfoOrg{
					Name:          "demo",
					PrimaryDomain: "demo.localhost",
				},
				Metadata: []UserMetadata{
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 35, 30, 126849000, time.FixedZone("", 0)),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 35, 30, 126849000, time.FixedZone("", 0)),
						Sequence:      3,
						ResourceOwner: "231848297847848962",
						Key:           "first",
						Value:         []byte("Hello World!"),
					},
					{
						CreationDate:  time.Date(2023, time.November, 14, 13, 35, 44, 28343000, time.FixedZone("", 0)),
						ChangeDate:    time.Date(2023, time.November, 14, 13, 35, 44, 28343000, time.FixedZone("", 0)),
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

				got, err := q.GetOIDCUserInfo(ctx, tt.args.userID)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
