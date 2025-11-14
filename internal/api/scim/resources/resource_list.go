package resources

import (
	"net/http"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ListRequest struct {
	Schemas []schemas.ScimSchemaType `json:"schemas"`

	// Count An integer indicating the desired maximum number of query results per page.
	Count int64 `json:"count" schema:"count"`

	// StartIndex An integer indicating the 1-based index of the first query result.
	StartIndex int64 `json:"startIndex" schema:"startIndex"`

	// Filter a scim filter expression to filter the query result.
	Filter *filter.Filter `json:"filter,omitempty" schema:"filter"`

	// SortBy attribute path to the sort attribute
	SortBy    string               `json:"sortBy" schema:"sortBy"`
	SortOrder ListRequestSortOrder `json:"sortOrder" schema:"sortOrder"`
}

type ListResponse[T any] struct {
	Schemas      []schemas.ScimSchemaType `json:"schemas"`
	ItemsPerPage uint64                   `json:"itemsPerPage"`
	TotalResults uint64                   `json:"totalResults"`
	StartIndex   uint64                   `json:"startIndex"`
	Resources    []T                      `json:"Resources"` // according to the rfc this is the only field in PascalCase...
}

type ListRequestSortOrder string

const (
	ListRequestSortOrderAsc ListRequestSortOrder = "ascending"
	ListRequestSortOrderDsc ListRequestSortOrder = "descending"

	defaultListCount = 100
	MaxListCount     = 100
)

var parser = zhttp.NewParser()

func (r *ListRequest) GetSchemas() []schemas.ScimSchemaType {
	return r.Schemas
}

func (o ListRequestSortOrder) isDefined() bool {
	switch o {
	case ListRequestSortOrderAsc, ListRequestSortOrderDsc:
		return true
	default:
		return false
	}
}

func (o ListRequestSortOrder) IsAscending() bool {
	return o == ListRequestSortOrderAsc
}

func NewListResponse[T any](totalResultCount uint64, q query.SearchRequest, resources []T) *ListResponse[T] {
	return &ListResponse[T]{
		Schemas:      []schemas.ScimSchemaType{schemas.IdListResponse},
		ItemsPerPage: q.Limit,
		TotalResults: totalResultCount,
		StartIndex:   q.Offset + 1, // start index is 1 based
		Resources:    resources,
	}
}

func (adapter *ResourceHandlerAdapter[T]) readListRequest(r *http.Request) (*ListRequest, error) {
	request := &ListRequest{
		Count:      defaultListCount,
		StartIndex: 1,
		SortOrder:  ListRequestSortOrderAsc,
	}

	switch r.Method {
	case http.MethodGet:
		if err := parser.Parse(r, request); err != nil {
			err = parser.UnwrapParserError(err)

			if serrors.IsScimOrZitadelError(err) {
				return nil, err
			}

			return nil, zerrors.ThrowInvalidArgument(nil, "SCIM-ullform", "Could not decode form: "+err.Error())
		}
	case http.MethodPost:
		if err := readSchema(r.Body, request, schemas.IdSearchRequest); err != nil {
			return nil, err
		}

		// json deserialization initializes this field if an empty string is provided
		// to not special case this in the resource implementation,
		// set it to nil here.
		if request.Filter.IsZero() {
			request.Filter = nil
		}
	}

	return request, request.validate()
}

func (r *ListRequest) toSearchRequest(defaultSortCol query.Column, fieldPathColumnMapping filter.FieldPathMapping) (query.SearchRequest, error) {
	sr := query.SearchRequest{
		Offset: uint64(r.StartIndex - 1), // start index is 1 based
		Limit:  uint64(r.Count),
		Asc:    r.SortOrder.IsAscending(),
	}

	if r.SortBy == "" {
		// set a default sort to ensure consistent results
		sr.SortingColumn = defaultSortCol
	} else if sortCol, err := fieldPathColumnMapping.Resolve(r.SortBy); err != nil || sortCol.FieldType == filter.FieldTypeCustom {
		return sr, serrors.ThrowInvalidValue(zerrors.ThrowInvalidArgument(err, "SCIM-SRT1", "SortBy field is unknown or not supported"))
	} else {
		sr.SortingColumn = sortCol.Column
	}

	return sr, nil
}

func (r *ListRequest) validate() error {
	// according to the spec values < 1 are treated as 1
	if r.StartIndex < 1 {
		r.StartIndex = 1
	}

	// according to the spec values < 0 are treated as 0
	if r.Count < 0 {
		r.Count = 0
	} else if r.Count > MaxListCount {
		return zerrors.ThrowInvalidArgumentf(nil, "SCIM-ucr", "Limit count exceeded, set a count <= %v", MaxListCount)
	}

	if !r.SortOrder.isDefined() {
		return zerrors.ThrowInvalidArgument(nil, "SCIM-ucx", "Invalid sort order")
	}

	return nil
}
