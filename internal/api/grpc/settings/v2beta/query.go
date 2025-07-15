package settings

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	filter_pb "github.com/zitadel/zitadel/pkg/grpc/filter/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

func (s *Server) ListOrganizationSettings(ctx context.Context, req *connect.Request[settings.ListOrganizationSettingsRequest]) (*connect.Response[settings.ListOrganizationSettingsResponse], error) {
	queries, err := s.listOrganizationSettingsRequestToModel(req.Msg)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.SearchOrganizationSettings(ctx, queries, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.ListOrganizationSettingsResponse{
		OrganizationSettings: organizationSettingsListToPb(resp.OrganizationSettingsList),
		Pagination:           filter.QueryToPaginationPb(queries.SearchRequest, resp.SearchResponse),
	}), nil
}

func (s *Server) listOrganizationSettingsRequestToModel(req *settings.ListOrganizationSettingsRequest) (*query.OrganizationSettingsSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(s.systemDefaults, req.Pagination)
	if err != nil {
		return nil, err
	}
	queries, err := organizationSettingsFiltersToQuery(req.Filters)
	if err != nil {
		return nil, err
	}
	return &query.OrganizationSettingsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: organizationSettingsFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func organizationSettingsFieldNameToSortingColumn(field *settings.OrganizationSettingsFieldName) query.Column {
	if field == nil {
		return query.OrganizationSettingsColumnCreationDate
	}
	switch *field {
	case settings.OrganizationSettingsFieldName_ORGANIZATION_SETTINGS_FIELD_NAME_CREATION_DATE:
		return query.OrganizationSettingsColumnCreationDate
	case settings.OrganizationSettingsFieldName_ORGANIZATION_SETTINGS_FIELD_NAME_ORGANIZATION_ID:
		return query.OrganizationSettingsColumnID
	case settings.OrganizationSettingsFieldName_ORGANIZATION_SETTINGS_FIELD_NAME_CHANGE_DATE:
		return query.OrganizationSettingsColumnChangeDate
	case settings.OrganizationSettingsFieldName_ORGANIZATION_SETTINGS_FIELD_NAME_UNSPECIFIED:
		return query.OrganizationSettingsColumnCreationDate
	default:
		return query.OrganizationSettingsColumnCreationDate
	}
}

func organizationSettingsFiltersToQuery(queries []*settings.OrganizationSettingsSearchFilter) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, qry := range queries {
		q[i], err = organizationSettingsToModel(qry)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func organizationSettingsToModel(filter *settings.OrganizationSettingsSearchFilter) (query.SearchQuery, error) {
	switch q := filter.Filter.(type) {
	case *settings.OrganizationSettingsSearchFilter_InOrganizationIdsFilter:
		return organizationInIDsFilterToQuery(q.InOrganizationIdsFilter)
	case *settings.OrganizationSettingsSearchFilter_OrganizationScopedUsernamesFilter:
		return organizationScopedUsernamesFilterToQuery(q.OrganizationScopedUsernamesFilter)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "SETTINGS-TODO", "List.Query.Invalid")
	}
}

func organizationInIDsFilterToQuery(q *filter_pb.InIDsFilter) (query.SearchQuery, error) {
	return query.NewOrganizationSettingsOrganizationIDSearchQuery(q.Ids)
}

func organizationScopedUsernamesFilterToQuery(q *settings.OrganizationScopedUsernamesFilter) (query.SearchQuery, error) {
	return query.NewOrganizationSettingsOrganizationScopedUsernamesSearchQuery(q.OrganizationScopedUsernames)
}

func organizationSettingsListToPb(settingsList []*query.OrganizationSettings) []*settings.OrganizationSettings {
	o := make([]*settings.OrganizationSettings, len(settingsList))
	for i, organizationSettings := range settingsList {
		o[i] = organizationSettingsToPb(organizationSettings)
	}
	return o
}

func organizationSettingsToPb(organizationSettings *query.OrganizationSettings) *settings.OrganizationSettings {
	return &settings.OrganizationSettings{
		OrganizationId:              organizationSettings.ID,
		CreationDate:                timestamppb.New(organizationSettings.CreationDate),
		ChangeDate:                  timestamppb.New(organizationSettings.ChangeDate),
		OrganizationScopedUsernames: organizationSettings.OrganizationScopedUsernames,
	}
}
