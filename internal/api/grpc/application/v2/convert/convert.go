package convert

import (
	"net/url"
	"strings"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/filter/v2"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func AppToPb(query_app *query.App) *application.Application {
	if query_app == nil {
		return &application.Application{}
	}

	return &application.Application{
		Id:           query_app.ID,
		CreationDate: timestamppb.New(query_app.CreationDate),
		ChangeDate:   timestamppb.New(query_app.ChangeDate),
		State:        appStateToPb(query_app.State),
		Name:         query_app.Name,
		Config:       appConfigToPb(query_app),
	}
}

func AppsToPb(queryApps []*query.App) []*application.Application {
	pbApps := make([]*application.Application, len(queryApps))

	for i, queryApp := range queryApps {
		pbApps[i] = AppToPb(queryApp)
	}

	return pbApps
}

func ListApplicationsRequestToModel(sysDefaults systemdefaults.SystemDefaults, req *application.ListApplicationsRequest) (*query.AppSearchQueries, error) {
	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	queries, err := appQueriesToModel(req.GetQueries())
	if err != nil {
		return nil, err
	}
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

func appSortingToColumn(sortingCriteria application.ApplicationSorting) query.Column {
	switch sortingCriteria {
	case application.ApplicationSorting_APPLICATION_SORT_BY_CHANGE_DATE:
		return query.AppColumnChangeDate
	case application.ApplicationSorting_APPLICATION_SORT_BY_CREATION_DATE:
		return query.AppColumnCreationDate
	case application.ApplicationSorting_APPLICATION_SORT_BY_NAME:
		return query.AppColumnName
	case application.ApplicationSorting_APPLICATION_SORT_BY_STATE:
		return query.AppColumnState
	case application.ApplicationSorting_APPLICATION_SORT_BY_ID:
		fallthrough
	default:
		return query.AppColumnID
	}
}

func appStateToPb(state domain.AppState) application.ApplicationState {
	switch state {
	case domain.AppStateActive:
		return application.ApplicationState_APPLICATION_STATE_ACTIVE
	case domain.AppStateInactive:
		return application.ApplicationState_APPLICATION_STATE_INACTIVE
	case domain.AppStateRemoved:
		return application.ApplicationState_APPLICATION_STATE_REMOVED
	case domain.AppStateUnspecified:
		fallthrough
	default:
		return application.ApplicationState_APPLICATION_STATE_UNSPECIFIED
	}
}

func appConfigToPb(app *query.App) application.ApplicationConfig {
	if app.OIDCConfig != nil {
		return appOIDCConfigToPb(app.OIDCConfig)
	}
	if app.SAMLConfig != nil {
		return appSAMLConfigToPb(app.SAMLConfig)
	}
	return appAPIConfigToPb(app.APIConfig)
}

func loginVersionToDomain(version *application.LoginVersion) (*domain.LoginVersion, *string, error) {
	switch v := version.GetVersion().(type) {
	case nil:
		return gu.Ptr(domain.LoginVersionUnspecified), gu.Ptr(""), nil
	case *application.LoginVersion_LoginV1:
		return gu.Ptr(domain.LoginVersion1), gu.Ptr(""), nil
	case *application.LoginVersion_LoginV2:
		_, err := url.Parse(v.LoginV2.GetBaseUri())
		return gu.Ptr(domain.LoginVersion2), gu.Ptr(v.LoginV2.GetBaseUri()), err
	default:
		return gu.Ptr(domain.LoginVersionUnspecified), gu.Ptr(""), nil
	}
}

func loginVersionToPb(version domain.LoginVersion, baseURI *string) *application.LoginVersion {
	switch version {
	case domain.LoginVersionUnspecified:
		return nil
	case domain.LoginVersion1:
		return &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}}
	case domain.LoginVersion2:
		return &application.LoginVersion{Version: &application.LoginVersion_LoginV2{LoginV2: &application.LoginV2{BaseUri: baseURI}}}
	default:
		return nil
	}
}

func appQueriesToModel(queries []*application.ApplicationSearchQuery) (toReturn []query.SearchQuery, err error) {
	toReturn = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		toReturn[i], err = appQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return toReturn, nil
}

func appQueryToModel(appQuery *application.ApplicationSearchQuery) (query.SearchQuery, error) {
	switch q := appQuery.GetQuery().(type) {
	case *application.ApplicationSearchQuery_ProjectIdQuery:
		return query.NewAppProjectIDSearchQuery(q.ProjectIdQuery.GetProjectId())
	case *application.ApplicationSearchQuery_NameQuery:
		return query.NewAppNameSearchQuery(filter.TextMethodPbToQuery(q.NameQuery.GetMethod()), q.NameQuery.Name)
	case *application.ApplicationSearchQuery_StateQuery:
		return query.NewAppStateSearchQuery(domain.AppState(q.StateQuery))
	case *application.ApplicationSearchQuery_ApiAppOnly:
		return query.NewNotNullQuery(query.AppAPIConfigColumnAppID)
	case *application.ApplicationSearchQuery_OidcAppOnly:
		return query.NewNotNullQuery(query.AppOIDCConfigColumnAppID)
	case *application.ApplicationSearchQuery_SamlAppOnly:
		return query.NewNotNullQuery(query.AppSAMLConfigColumnAppID)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "CONV-z2mAGy", "List.Query.Invalid")
	}
}

func CreateAPIClientKeyRequestToDomain(key *application.CreateApplicationKeyRequest) *domain.ApplicationKey {
	return &domain.ApplicationKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: strings.TrimSpace(key.GetProjectId()),
		},
		ExpirationDate: key.GetExpirationDate().AsTime(),
		Type:           domain.AuthNKeyTypeJSON,
		ApplicationID:  strings.TrimSpace(key.GetApplicationId()),
	}
}

func ListApplicationKeysRequestToDomain(sysDefaults systemdefaults.SystemDefaults, req *application.ListApplicationKeysRequest) (*query.AuthNKeySearchQueries, error) {
	var queries []query.SearchQuery

	queries, err := ApplicationKeySearchQueriesToQuery(req.GetQueries())
	if err != nil {
		return nil, err
	}

	offset, limit, asc, err := filter.PaginationPbToQuery(sysDefaults, req.GetPagination())
	if err != nil {
		return nil, err
	}

	return &query.AuthNKeySearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: appKeysSortingToColumn(req.GetSortingColumn()),
		},

		Queries: queries,
	}, nil
}

func ApplicationKeySearchQueriesToQuery(queries []*application.ApplicationKeySearchQuery) (_ []query.SearchQuery, err error) {
	searchQueries := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		searchQueries[i], err = ApplicationKeySearchQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return searchQueries, nil
}

func ApplicationKeySearchQueryToQuery(searchQuery *application.ApplicationKeySearchQuery) (query.SearchQuery, error) {
	switch f := searchQuery.GetQuery().(type) {
	case *application.ApplicationKeySearchQuery_ApplicationIdQuery:
		return query.NewAuthNKeyObjectIDQuery(strings.TrimSpace(f.ApplicationIdQuery.GetApplicationId()))
	case *application.ApplicationKeySearchQuery_OrganizationIdQuery:
		return query.NewAuthNKeyResourceOwnerQuery(strings.TrimSpace(f.OrganizationIdQuery.GetOrganizationId()))
	case *application.ApplicationKeySearchQuery_ProjectIdQuery:
		return query.NewAuthNKeyAggregateIDQuery(strings.TrimSpace(f.ProjectIdQuery.GetProjectId()))
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "CONV-t3ENme", "List.Query.Invalid")
	}
}

func appKeysSortingToColumn(sortingCriteria application.ApplicationKeysSorting) query.Column {
	switch sortingCriteria {
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_PROJECT_ID:
		return query.AuthNKeyColumnAggregateID
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_CREATION_DATE:
		return query.AuthNKeyColumnCreationDate
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_EXPIRATION:
		return query.AuthNKeyColumnExpiration
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_ORGANIZATION_ID:
		return query.AuthNKeyColumnResourceOwner
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_TYPE:
		return query.AuthNKeyColumnType
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_APPLICATION_ID:
		return query.AuthNKeyColumnObjectID
	case application.ApplicationKeysSorting_APPLICATION_KEYS_SORT_BY_ID:
		fallthrough
	default:
		return query.AuthNKeyColumnID
	}
}

func ApplicationKeysToPb(keys []*query.AuthNKey) []*application.ApplicationKey {
	pbAppKeys := make([]*application.ApplicationKey, len(keys))

	for i, k := range keys {
		pbKey := &application.ApplicationKey{
			KeyId:          k.ID,
			ApplicationId:  k.ApplicationID,
			ProjectId:      k.AggregateID,
			CreationDate:   timestamppb.New(k.CreationDate),
			OrganizationId: k.ResourceOwner,
			ExpirationDate: timestamppb.New(k.Expiration),
		}
		pbAppKeys[i] = pbKey
	}

	return pbAppKeys
}
