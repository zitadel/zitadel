package internal_permission

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	"github.com/zitadel/zitadel/pkg/grpc/internal_permission/v2"
)

func Test_administratorFilterToQuery(t *testing.T) {
	now := time.Now().UTC()
	type args struct {
		filter *internal_permission.AdministratorSearchFilter
		level  uint8
	}
	tests := []struct {
		name    string
		args    args
		want    query.SearchQuery
		wantErr error
	}{
		{
			name: "max nested queries",
			args: args{
				filter: &internal_permission.AdministratorSearchFilter{
					Filter: &internal_permission.AdministratorSearchFilter_And{
						And: &internal_permission.AndFilter{
							Queries: []*internal_permission.AdministratorSearchFilter{
								{
									Filter: &internal_permission.AdministratorSearchFilter_CreationDate{
										CreationDate: &filter.TimestampFilter{
											Timestamp: timestamppb.New(now),
										},
									},
								},
							},
						},
					},
				},
				level: 19,
			},
			want: func() query.SearchQuery {
				dateQuery, _ := query.NewMembershipCreationDateQuery(now, query.TimestampEquals)
				q, _ := query.NewAndQuery(dateQuery)
				return q
			}(),
		},
		{
			name: "too many nested queries",
			args: args{
				filter: &internal_permission.AdministratorSearchFilter{
					Filter: &internal_permission.AdministratorSearchFilter_And{
						And: &internal_permission.AndFilter{
							Queries: []*internal_permission.AdministratorSearchFilter{
								{
									Filter: &internal_permission.AdministratorSearchFilter_And{
										And: &internal_permission.AndFilter{
											Queries: []*internal_permission.AdministratorSearchFilter{
												{
													Filter: &internal_permission.AdministratorSearchFilter_CreationDate{
														CreationDate: &filter.TimestampFilter{
															Timestamp: timestamppb.New(time.Now().UTC()),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				level: 19,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PERM-zsQ97", "Errors.Query.TooManyNestingLevels"),
		},
		{
			name: "invalid filter",
			args: args{
				filter: &internal_permission.AdministratorSearchFilter{
					Filter: nil,
				},
				level: 0,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "PERM-vR9nC", "List.Query.Invalid"),
		},
		// rest is tested in integration tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := administratorFilterToQuery(tt.args.filter, tt.args.level)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
