package convert

import (
	"net/url"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func AppToPb(query_app *query.App) *app.Application {
	if query_app == nil {
		return &app.Application{}
	}

	return &app.Application{
		Id:           query_app.ID,
		CreationDate: timestamppb.New(query_app.CreationDate),
		ChangeDate:   timestamppb.New(query_app.ChangeDate),
		State:        appStateToPb(query_app.State),
		Name:         query_app.Name,
		Config:       appConfigToPb(query_app),
	}
}

func AppsToPb(queryApps []*query.App) []*app.Application {
	pbApps := make([]*app.Application, len(queryApps))

	for i, queryApp := range queryApps {
		pbApps[i] = AppToPb(queryApp)
	}

	return pbApps
}

func ListApplicationsRequestToModel(sysDefaults systemdefaults.SystemDefaults, req *app.ListApplicationsRequest) (*query.AppSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := appQueriesToModel(req.GetFilters())
	if err != nil {
		return nil, err
	}
	projectQuery, err := query.NewAppProjectIDSearchQuery(req.GetProjectId())
	if err != nil {
		return nil, err
	}

	queries = append(queries, projectQuery)
	return &query.AppSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: appSortingToColumn(req.GetSortingColumn()),
		},

		Queries: queries,
	}, nil
}

func appSortingToColumn(sortingCriteria app.AppSorting) query.Column {
	switch sortingCriteria {
	case app.AppSorting_APP_SORT_BY_CHANGE_DATE:
		return query.AppColumnChangeDate
	case app.AppSorting_APP_SORT_BY_CREATION_DATE:
		return query.AppColumnCreationDate
	case app.AppSorting_APP_SORT_BY_NAME:
		return query.AppColumnName
	case app.AppSorting_APP_SORT_BY_STATE:
		return query.AppColumnState
	case app.AppSorting_APP_SORT_BY_ID:
		fallthrough
	default:
		return query.AppColumnID
	}
}

func appStateToPb(state domain.AppState) app.AppState {
	switch state {
	case domain.AppStateActive:
		return app.AppState_APP_STATE_ACTIVE
	case domain.AppStateInactive:
		return app.AppState_APP_STATE_INACTIVE
	case domain.AppStateRemoved:
		return app.AppState_APP_STATE_REMOVED
	case domain.AppStateUnspecified:
		fallthrough
	default:
		return app.AppState_APP_STATE_UNSPECIFIED
	}
}

func appConfigToPb(app *query.App) app.ApplicationConfig {
	if app.OIDCConfig != nil {
		return appOIDCConfigToPb(app.OIDCConfig)
	}
	if app.SAMLConfig != nil {
		return appSAMLConfigToPb(app.SAMLConfig)
	}
	return appAPIConfigToPb(app.APIConfig)
}

func loginVersionToDomain(version *app.LoginVersion) (*domain.LoginVersion, *string, error) {
	switch v := version.GetVersion().(type) {
	case nil:
		return gu.Ptr(domain.LoginVersionUnspecified), gu.Ptr(""), nil
	case *app.LoginVersion_LoginV1:
		return gu.Ptr(domain.LoginVersion1), gu.Ptr(""), nil
	case *app.LoginVersion_LoginV2:
		_, err := url.Parse(v.LoginV2.GetBaseUri())
		return gu.Ptr(domain.LoginVersion2), gu.Ptr(v.LoginV2.GetBaseUri()), err
	default:
		return gu.Ptr(domain.LoginVersionUnspecified), gu.Ptr(""), nil
	}
}

func loginVersionToPb(version domain.LoginVersion, baseURI *string) *app.LoginVersion {
	switch version {
	case domain.LoginVersionUnspecified:
		return nil
	case domain.LoginVersion1:
		return &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}}
	case domain.LoginVersion2:
		return &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: baseURI}}}
	default:
		return nil
	}
}

func appQueriesToModel(queries []*app.ApplicationSearchFilter) (toReturn []query.SearchQuery, err error) {
	toReturn = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		toReturn[i], err = appQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return toReturn, nil
}

func appQueryToModel(appQuery *app.ApplicationSearchFilter) (query.SearchQuery, error) {
	switch q := appQuery.GetFilter().(type) {
	case *app.ApplicationSearchFilter_NameFilter:
		return query.NewAppNameSearchQuery(filter.TextMethodPbToQuery(q.NameFilter.GetMethod()), q.NameFilter.Name)
	case *app.ApplicationSearchFilter_StateFilter:
		return query.NewAppStateSearchQuery(domain.AppState(q.StateFilter))
	case *app.ApplicationSearchFilter_ApiAppOnly:
		return query.NewNotNullQuery(query.AppAPIConfigColumnAppID)
	case *app.ApplicationSearchFilter_OidcAppOnly:
		return query.NewNotNullQuery(query.AppOIDCConfigColumnAppID)
	case *app.ApplicationSearchFilter_SamlAppOnly:
		return query.NewNotNullQuery(query.AppSAMLConfigColumnAppID)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "CONV-z2mAGy", "List.Query.Invalid")
	}
}
