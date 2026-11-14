package resources

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadListRequest_excludedAttributes(t *testing.T) {
	const search = `{"schemas":["urn:ietf:params:scim:api:messages:2.0:SearchRequest"],"excludedAttributes":`
	tests := []struct {
		name    string
		request *http.Request
	}{
		{"query", httptest.NewRequest(http.MethodGet, "/Groups?excludedAttributes=members,%20displayName", nil)},
		{"json string", httptest.NewRequest(http.MethodPost, "/Groups/.search", strings.NewReader(search+`"members, displayName"}`))},
		{"json array", httptest.NewRequest(http.MethodPost, "/Groups/.search", strings.NewReader(search+`["members","displayName"]}`))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := new(ResourceHandlerAdapter[*ScimGroup]).readListRequest(tt.request)
			require.NoError(t, err)
			assert.Equal(t, Attributes{"members", "displayName"}, got.ExcludedAttributes)
			assert.True(t, excludeGroupMembers(got.ExcludedAttributes))
		})
	}
}

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
