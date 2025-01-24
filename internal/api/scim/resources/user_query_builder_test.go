package resources

import (
	"context"
	"reflect"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/test"
)

func Test_buildMetadataQuery(t *testing.T) {
	tests := []struct {
		name    string
		key     metadata.Key
		value   *filter.CompValue
		op      *filter.CompareOp
		want    query.SearchQuery
		wantErr bool
	}{
		{
			name:    "equals",
			key:     "foo",
			value:   &filter.CompValue{StringValue: gu.Ptr("bar")},
			op:      &filter.CompareOp{Equal: true},
			want:    test.Must(query.NewUserMetadataExistsQuery("foo", []byte("bar"), query.TextEquals, query.BytesEquals)),
			wantErr: false,
		},
		{
			name:    "not equals",
			key:     "foo",
			value:   &filter.CompValue{StringValue: gu.Ptr("bar")},
			op:      &filter.CompareOp{NotEqual: true},
			want:    test.Must(query.NewUserMetadataExistsQuery("foo", []byte("bar"), query.TextEquals, query.BytesNotEquals)),
			wantErr: false,
		},
		{
			name:    "unsupported operator",
			key:     "foo",
			value:   &filter.CompValue{StringValue: gu.Ptr("bar")},
			op:      &filter.CompareOp{StartsWith: true},
			wantErr: true,
		},
		{
			name:    "unsupported comparison value",
			key:     "foo",
			value:   &filter.CompValue{Int: gu.Ptr(10)},
			op:      &filter.CompareOp{Equal: true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildMetadataQuery(context.Background(), tt.key, tt.value, tt.op)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildMetadataQuery() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_buildActiveUserStateQuery(t *testing.T) {
	tests := []struct {
		name         string
		compareValue *filter.CompValue
		compOp       *filter.CompareOp
		want         query.SearchQuery
		wantErr      bool
	}{
		{
			name:         "eq true",
			compareValue: &filter.CompValue{BooleanTrue: true},
			compOp:       &filter.CompareOp{Equal: true},
			want: test.Must(query.NewOrQuery(
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberEquals)),
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberEquals)),
			)),
		},
		{
			name:         "eq false",
			compareValue: &filter.CompValue{BooleanFalse: true},
			compOp:       &filter.CompareOp{Equal: true},
			want: test.Must(query.NewAndQuery(
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberNotEquals)),
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberNotEquals)),
			)),
		},
		{
			name:         "ne true",
			compareValue: &filter.CompValue{BooleanTrue: true},
			compOp:       &filter.CompareOp{NotEqual: true},
			want: test.Must(query.NewAndQuery(
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberNotEquals)),
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberNotEquals)),
			)),
		},
		{
			name:         "ne false",
			compareValue: &filter.CompValue{BooleanTrue: true},
			compOp:       &filter.CompareOp{Equal: true},
			want: test.Must(query.NewOrQuery(
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateInitial), query.NumberEquals)),
				test.Must(query.NewNumberQuery(query.UserStateCol, int32(domain.UserStateActive), query.NumberEquals)),
			)),
		},
		{
			name:         "invalid operator",
			compareValue: &filter.CompValue{BooleanTrue: true},
			compOp:       &filter.CompareOp{StartsWith: true},
			wantErr:      true,
		},
		{
			name:         "invalid comp value",
			compareValue: &filter.CompValue{StringValue: gu.Ptr("foo")},
			compOp:       &filter.CompareOp{Equal: true},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildActiveUserStateQuery(context.Background(), tt.compareValue, tt.compOp)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equalf(t, tt.want, got, "buildActiveUserStateQuery(%#v, %#v)", tt.compareValue, tt.compOp)
		})
	}
}
