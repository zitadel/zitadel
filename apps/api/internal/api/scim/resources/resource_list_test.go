package resources

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ListRequest
		want    *ListRequest
		wantErr bool
	}{
		{
			name: "valid",
			req: &ListRequest{
				SortOrder: ListRequestSortOrderAsc,
			},
		},
		{
			name: "invalid sort order",
			req: &ListRequest{
				SortOrder: "fooBar",
			},
			wantErr: true,
		},
		{
			name: "count too big",
			req: &ListRequest{
				Count:     99999999,
				SortOrder: ListRequestSortOrderAsc,
			},
			wantErr: true,
		},
		{
			name: "negative start index",
			req: &ListRequest{
				StartIndex: -1,
				Count:      10,
				SortOrder:  ListRequestSortOrderAsc,
			},
			want: &ListRequest{
				StartIndex: 1,
				Count:      10,
				SortOrder:  ListRequestSortOrderAsc,
			},
		},
		{
			name: "negative count",
			req: &ListRequest{
				StartIndex: 10,
				Count:      -1,
				SortOrder:  ListRequestSortOrderAsc,
			},
			want: &ListRequest{
				StartIndex: 10,
				Count:      0,
				SortOrder:  ListRequestSortOrderAsc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.want != nil && !reflect.DeepEqual(tt.req, tt.want) {
				t.Errorf("got: %#v, want: %#v", tt.req, tt.want)
			}
		})
	}
}
