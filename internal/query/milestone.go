package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/repository/milestone"

	"github.com/zitadel/zitadel/internal/api/authz"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Milestones struct {
	SearchResponse
	Milestones []*Milestone
}

type Milestone struct {
	InstanceID    string
	MilestoneType milestone.PushedEventType
	ReachedDate   time.Time
	PushedDate    time.Time
	PrimaryDomain string
}

type MilestonesSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *MilestonesSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

var (
	milestonesTable = table{
		name:          projection.MilestonesProjectionTable,
		instanceIDCol: projection.MilestoneColumnInstanceID,
	}
	MilestoneInstanceIDColID = Column{
		name:  projection.MilestoneColumnInstanceID,
		table: milestonesTable,
	}
	MilestoneTypeColID = Column{
		name:  projection.MilestoneColumnMilestoneType,
		table: milestonesTable,
	}
	MilestonePrimaryDomainColID = Column{
		name:  projection.MilestoneColumnPrimaryDomain,
		table: milestonesTable,
	}
	MilestoneReachedDateColID = Column{
		name:  projection.MilestoneColumnReachedDate,
		table: milestonesTable,
	}
	MilestonePushedDateColID = Column{
		name:  projection.MilestoneColumnPushedDate,
		table: milestonesTable,
	}
)

// SearchMilestones tries to defer the instanceID from the passed context if no instanceIDs are passed
func (q *Queries) SearchMilestones(ctx context.Context, instanceIDs []string, queries *MilestonesSearchQueries) (_ *Milestones, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	query, scan := prepareMilestonesQuery(ctx, q.client)
	if len(instanceIDs) == 0 {
		instanceIDs = []string{authz.GetInstance(ctx).InstanceID()}
	}
	stmt, args, err := queries.toQuery(query).
		Where(sq.Eq{
			MilestoneInstanceIDColID.identifier(): fmt.Sprintf("IN (%s)", strings.Join(instanceIDs, ",")),
		}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-A9i5k", "Errors.Query.SQLStatement")
	}
	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	milestones, err := scan(rows)
	if err != nil {
		return nil, err
	}
	milestones.LatestSequence, err = q.latestSequence(ctx, milestonesTable)
	return milestones, err

}

func prepareMilestonesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*Milestones, error)) {
	return sq.Select(
			MilestonePrimaryDomainColID.identifier(),
			MilestoneReachedDateColID.identifier(),
			MilestonePushedDateColID.identifier(),
			MilestoneTypeColID.identifier(),
			countColumn.identifier(),
		).
			From(notificationPolicyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Milestones, error) {
			milestones := make([]*Milestone, 0)
			var count uint64
			for rows.Next() {
				m := new(Milestone)
				err := rows.Scan(
					&m.PrimaryDomain,
					&m.ReachedDate,
					&m.MilestoneType,
					&count,
				)
				if err != nil {
					return nil, err
				}
				milestones = append(milestones, m)
			}
			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-CK9mI", "Errors.Query.CloseRows")
			}
			return &Milestones{
				Milestones: milestones,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
