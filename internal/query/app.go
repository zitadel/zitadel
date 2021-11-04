package query

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query/projection"
)

type Apps struct {
	SearchResponse
	Apps []*App
}

type App struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.AppState
	Sequence      uint64

	ProjectID string
	Name      string

	OIDCConfig *OIDCApp
	APIConfig  *APIApp
}

type OIDCApp struct {
	RedirectURIs           []string
	ResponseTypes          []domain.OIDCResponseType
	GrantTypes             []domain.OIDCGrantType
	AppType                domain.OIDCApplicationType
	ClientID               string
	ClientSecret           *crypto.CryptoValue
	AuthMethodType         domain.OIDCAuthMethodType
	PostLogoutRedirectURIs []string
	Version                domain.OIDCVersion
	NoneCompliant          bool
	ComplianceProblems     []string
	IsDevMode              bool
	AccessTokenType        domain.OIDCTokenType
	AssertAccessTokenRole  bool
	AssertIDTokenRole      bool
	AssertIDTokenUserinfo  bool
	ClockSkew              time.Duration
	AdditionalOrigins      []string
	AllowedOrigins         []string
}

type APIApp struct {
	ClientID       string
	ClientSecret   *crypto.CryptoValue
	AuthMethodType domain.APIAuthMethodType
}

type AppSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

var (
	appsTable = table{
		name: projection.AppProjectionTable,
	}
	AppColumnID = Column{
		name:  projection.AppColumnID,
		table: appsTable,
	}
	AppColumnName = Column{
		name:  projection.AppColumnName,
		table: appsTable,
	}
	AppColumnProjectID = Column{
		name:  projection.AppColumnProjectID,
		table: appsTable,
	}
	AppColumnCreationDate = Column{
		name:  projection.AppColumnCreationDate,
		table: appsTable,
	}
	AppColumnChangeDate = Column{
		name:  projection.AppColumnChangeDate,
		table: appsTable,
	}
	AppColumnResourceOwner = Column{
		name:  projection.AppColumnResourceOwner,
		table: appsTable,
	}
	AppColumnState = Column{
		name:  projection.AppColumnState,
		table: appsTable,
	}
	AppColumnSequence = Column{
		name:  projection.AppColumnSequence,
		table: appsTable,
	}
)

var (
	appAPIConfigsTable = table{
		name: projection.AppAPITable,
	}
	AppAPIConfigColumnAppID = Column{
		name:  projection.AppAPIConfigColumnAppID,
		table: appAPIConfigsTable,
	}
	AppAPIConfigColumnClientID = Column{
		name:  projection.AppAPIConfigColumnClientID,
		table: appAPIConfigsTable,
	}
	AppAPIConfigColumnClientSecret = Column{
		name:  projection.AppAPIConfigColumnClientSecret,
		table: appAPIConfigsTable,
	}
	AppAPIConfigColumnAuthMethod = Column{
		name:  projection.AppAPIConfigColumnAuthMethod,
		table: appAPIConfigsTable,
	}
)

var (
	appOIDCConfigsTable = table{
		name: projection.AppOIDCTable,
	}
	AppOIDCConfigColumnAppID = Column{
		name:  projection.AppOIDCConfigColumnAppID,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnVersion = Column{
		name:  projection.AppOIDCConfigColumnVersion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnClientID = Column{
		name:  projection.AppOIDCConfigColumnClientID,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnClientSecret = Column{
		name:  projection.AppOIDCConfigColumnClientSecret,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnRedirectUris = Column{
		name:  projection.AppOIDCConfigColumnRedirectUris,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnResponseTypes = Column{
		name:  projection.AppOIDCConfigColumnResponseTypes,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnGrantTypes = Column{
		name:  projection.AppOIDCConfigColumnGrantTypes,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnApplicationType = Column{
		name:  projection.AppOIDCConfigColumnApplicationType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAuthMethodType = Column{
		name:  projection.AppOIDCConfigColumnAuthMethodType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnPostLogoutRedirectUris = Column{
		name:  projection.AppOIDCConfigColumnPostLogoutRedirectUris,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnDevMode = Column{
		name:  projection.AppOIDCConfigColumnDevMode,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAccessTokenType = Column{
		name:  projection.AppOIDCConfigColumnAccessTokenType,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAccessTokenRoleAssertion = Column{
		name:  projection.AppOIDCConfigColumnAccessTokenRoleAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnIDTokenRoleAssertion = Column{
		name:  projection.AppOIDCConfigColumnIDTokenRoleAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnIDTokenUserinfoAssertion = Column{
		name:  projection.AppOIDCConfigColumnIDTokenUserinfoAssertion,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnClockSkew = Column{
		name:  projection.AppOIDCConfigColumnClockSkew,
		table: appOIDCConfigsTable,
	}
	AppOIDCConfigColumnAdditionalOrigins = Column{
		name:  projection.AppOIDCConfigColumnAdditionalOrigins,
		table: appOIDCConfigsTable,
	}
)

func (q *Queries) AppByProjectAndAppID(ctx context.Context, projectID, appID string) (*App, error) {
	return nil, nil
}

func (q *Queries) AppByID(ctx context.Context, appID string) (*App, error) {
	return nil, nil
}

func (q *Queries) ProjectIDFromAppID(ctx context.Context, appID string) (string, error) {
	return "", nil
}

func (q *Queries) ProjectByAppID(ctx context.Context, appID string) (*Project, error) {
	return nil, nil
}

func (q *Queries) AppByOIDCClientID(ctx context.Context, clientID string) (*App, error) {
	return nil, nil
}

func (q *Queries) SearchApps(ctx context.Context, queries *AppSearchQueries) (*Apps, error) {
	return nil, nil
}

func (q *Queries) SearchAppIDs(ctx context.Context, queries *AppSearchQueries) ([]string, error) {
	return nil, nil
}

func NewAppNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(AppColumnName, value, method)
}

func NewAppProjectIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(AppColumnProjectID, id, TextEquals)
}
