package query

import (
	"context"
	"slices"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	organizationSettingsTable = table{
		name:          projection.OrganizationSettingsTable,
		instanceIDCol: projection.OrganizationSettingsInstanceIDCol,
	}
	OrganizationSettingsColumnID = Column{
		name:  projection.OrganizationSettingsIDCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnCreationDate = Column{
		name:  projection.OrganizationSettingsCreationDateCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnChangeDate = Column{
		name:  projection.OrganizationSettingsChangeDateCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnResourceOwner = Column{
		name:  projection.OrganizationSettingsResourceOwnerCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnInstanceID = Column{
		name:  projection.OrganizationSettingsInstanceIDCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnSequence = Column{
		name:  projection.OrganizationSettingsSequenceCol,
		table: organizationSettingsTable,
	}
	OrganizationSettingsColumnOrganizationScopedUsernames = Column{
		name:  projection.OrganizationSettingsOrganizationScopedUsernamesCol,
		table: organizationSettingsTable,
	}
)

type OrganizationSettingsList struct {
	SearchResponse
	OrganizationSettingsList []*OrganizationSettings
}

func organizationSettingsListCheckPermission(ctx context.Context, organizationSettingsList *OrganizationSettingsList, permissionCheck domain.PermissionCheck) {
	organizationSettingsList.OrganizationSettingsList = slices.DeleteFunc(organizationSettingsList.OrganizationSettingsList,
		func(organizationSettings *OrganizationSettings) bool {
			return organizationSettingsCheckPermission(ctx, organizationSettings.ResourceOwner, organizationSettings.ID, permissionCheck) != nil
		},
	)
}

func organizationSettingsCheckPermission(ctx context.Context, resourceOwner string, id string, permissionCheck domain.PermissionCheck) error {
	return permissionCheck(ctx, domain.PermissionPolicyRead, resourceOwner, id)
}

type OrganizationSettings struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	Sequence      uint64

	OrganizationScopedUsernames bool
}

type OrganizationSettingsSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *OrganizationSettingsSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func organizationSettingsPermissionCheckV2(ctx context.Context, query sq.SelectBuilder, enabled bool, queries *OrganizationSettingsSearchQueries) sq.SelectBuilder {
	if !enabled {
		return query
	}
	join, args := PermissionClause(
		ctx,
		OrganizationSettingsColumnID,
		domain.PermissionPolicyRead,
		SingleOrgPermissionOption(queries.Queries),
	)
	return query.JoinClause(join, args...)
}

func (q *Queries) SearchOrganizationSettings(ctx context.Context, queries *OrganizationSettingsSearchQueries, permissionCheck domain.PermissionCheck) (*OrganizationSettingsList, error) {
	permissionCheckV2 := PermissionV2(ctx, permissionCheck)
	settings, err := q.searchOrganizationSettings(ctx, queries, permissionCheckV2)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil && !authz.GetFeatures(ctx).PermissionCheckV2 {
		organizationSettingsListCheckPermission(ctx, settings, permissionCheck)
	}
	return settings, nil
}

func (q *Queries) searchOrganizationSettings(ctx context.Context, queries *OrganizationSettingsSearchQueries, permissionCheckV2 bool) (settingsList *OrganizationSettingsList, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareOrganizationSettingsListQuery()
	query = organizationSettingsPermissionCheckV2(ctx, query, permissionCheckV2, queries)
	eq := sq.Eq{OrganizationSettingsColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-qNPeOXlMwj", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows database.Rows) error {
		settingsList, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, err
	}
	return settingsList, nil
}

func NewOrganizationSettingsOrganizationIDSearchQuery(ids []string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(OrganizationSettingsColumnID, list, ListIn)
}

func NewOrganizationSettingsOrganizationScopedUsernamesSearchQuery(organizationScopedUsernames bool) (SearchQuery, error) {
	return NewBoolQuery(OrganizationSettingsColumnOrganizationScopedUsernames, organizationScopedUsernames)
}

func prepareOrganizationSettingsListQuery() (sq.SelectBuilder, func(database.Rows) (*OrganizationSettingsList, error)) {
	return sq.Select(
			OrganizationSettingsColumnID.identifier(),
			OrganizationSettingsColumnCreationDate.identifier(),
			OrganizationSettingsColumnChangeDate.identifier(),
			OrganizationSettingsColumnResourceOwner.identifier(),
			OrganizationSettingsColumnSequence.identifier(),
			OrganizationSettingsColumnOrganizationScopedUsernames.identifier(),
			countColumn.identifier(),
		).From(organizationSettingsTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows database.Rows) (*OrganizationSettingsList, error) {
			settingsList := make([]*OrganizationSettings, 0)
			var (
				count uint64
			)
			for rows.Next() {
				settings := new(OrganizationSettings)
				err := rows.Scan(
					&settings.ID,
					&settings.CreationDate,
					&settings.ChangeDate,
					&settings.ResourceOwner,
					&settings.Sequence,
					&settings.OrganizationScopedUsernames,
					&count,
				)
				if err != nil {
					return nil, err
				}
				settingsList = append(settingsList, settings)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-mmC1K0t5Fq", "Errors.Query.CloseRows")
			}

			return &OrganizationSettingsList{
				OrganizationSettingsList: settingsList,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
