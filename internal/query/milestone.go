package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Milestones struct {
	SearchResponse
	Milestones []*Milestone
}

type Milestone struct {
	InstanceID    string
	Type          milestone.Type
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
		name:  projection.MilestoneColumnType,
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
	instanceIDParams := make([]string, len(instanceIDs))
	instanceIDArgs := make([]interface{}, len(instanceIDs))
	for idx := range instanceIDs {
		instanceIDParams[idx] = fmt.Sprintf("$%d", idx+1)
		instanceIDArgs[idx] = instanceIDs[idx]
	}
	expr := fmt.Sprintf("%s IN (%s)", MilestoneInstanceIDColID.name, strings.Join(instanceIDParams, ","))
	stmt, args, err := queries.toQuery(query).Where(sq.Expr(expr, instanceIDArgs...)).ToSql()
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
			MilestoneInstanceIDColID.identifier(),
			MilestonePrimaryDomainColID.identifier(),
			MilestoneReachedDateColID.identifier(),
			MilestonePushedDateColID.identifier(),
			MilestoneTypeColID.identifier(),
			countColumn.identifier(),
		).
			From(milestonesTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Milestones, error) {
			milestones := make([]*Milestone, 0)
			var count uint64
			for rows.Next() {
				m := new(Milestone)
				reachedDate := sql.NullTime{}
				pushedDate := sql.NullTime{}
				primaryDomain := sql.NullString{}
				err := rows.Scan(
					&m.InstanceID,
					&primaryDomain,
					&reachedDate,
					&pushedDate,
					&m.Type,
					&count,
				)
				if err != nil {
					return nil, err
				}
				m.PrimaryDomain = primaryDomain.String
				m.ReachedDate = reachedDate.Time
				m.PushedDate = pushedDate.Time
				milestones = append(milestones, m)
			}
			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-CK9mI", "Errors.Query.CloseRows")
			}
			if err := rows.Err(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-asLsI", "Errors.Internal")
			}
			return &Milestones{
				Milestones: milestones,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
