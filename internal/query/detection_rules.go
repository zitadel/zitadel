package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/detection"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	detectionRulesTable = table{
		name:          projection.DetectionRulesProjectionTable,
		instanceIDCol: projection.DetectionRulesColumnInstanceID,
	}
	DetectionRulesColumnID              = Column{name: projection.DetectionRulesColumnID, table: detectionRulesTable}
	DetectionRulesColumnCreationDate    = Column{name: projection.DetectionRulesColumnCreationDate, table: detectionRulesTable}
	DetectionRulesColumnChangeDate      = Column{name: projection.DetectionRulesColumnChangeDate, table: detectionRulesTable}
	DetectionRulesColumnInstanceID      = Column{name: projection.DetectionRulesColumnInstanceID, table: detectionRulesTable}
	DetectionRulesColumnSequence        = Column{name: projection.DetectionRulesColumnSequence, table: detectionRulesTable}
	DetectionRulesColumnDescription     = Column{name: projection.DetectionRulesColumnDescription, table: detectionRulesTable}
	DetectionRulesColumnExpr            = Column{name: projection.DetectionRulesColumnExpr, table: detectionRulesTable}
	DetectionRulesColumnEngine          = Column{name: projection.DetectionRulesColumnEngine, table: detectionRulesTable}
	DetectionRulesColumnFindingName     = Column{name: projection.DetectionRulesColumnFindingName, table: detectionRulesTable}
	DetectionRulesColumnFindingMessage  = Column{name: projection.DetectionRulesColumnFindingMessage, table: detectionRulesTable}
	DetectionRulesColumnFindingBlock    = Column{name: projection.DetectionRulesColumnFindingBlock, table: detectionRulesTable}
	DetectionRulesColumnContextTemplate = Column{name: projection.DetectionRulesColumnContextTemplate, table: detectionRulesTable}
	DetectionRulesColumnRateLimitKey    = Column{name: projection.DetectionRulesColumnRateLimitKey, table: detectionRulesTable}
	DetectionRulesColumnRateLimitWindow = Column{name: projection.DetectionRulesColumnRateLimitWindow, table: detectionRulesTable}
	DetectionRulesColumnRateLimitMax    = Column{name: projection.DetectionRulesColumnRateLimitMax, table: detectionRulesTable}
	DetectionRulesColumnPriority        = Column{name: projection.DetectionRulesColumnPriority, table: detectionRulesTable}
	DetectionRulesColumnStopOnMatch     = Column{name: projection.DetectionRulesColumnStopOnMatch, table: detectionRulesTable}
)

type DetectionRule struct {
	ID              string
	CreationDate    time.Time
	ChangeDate      time.Time
	Sequence        uint64
	Description     string
	Expr            string
	Action          detection.ActionType
	Priority        int64
	StopOnMatch     bool
	FindingName     string
	FindingMessage  string
	FindingBlock    bool
	ContextTemplate string
	RateLimitKey    string
	RateLimitWindow database.Duration
	RateLimitMax    int64
}

func (q *Queries) SearchDetectionRules(ctx context.Context) (_ []*DetectionRule, err error) {
	stmt, scan := prepareDetectionRulesQuery()
	query, args, err := stmt.Where(sq.Eq{DetectionRulesColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-ef9p3", "Errors.Query.InvalidRequest")
	}
	var rules []*DetectionRule
	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		rules, err = scan(rows)
		return err
	}, query, args...)
	return rules, err
}

func (q *Queries) DetectionRule(ctx context.Context, ruleID string) (_ *DetectionRule, err error) {
	stmt, scan := prepareDetectionRuleQuery()
	query, args, err := stmt.Where(sq.Eq{
		DetectionRulesColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
		DetectionRulesColumnID.identifier():         ruleID,
	}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-3v6G7", "Errors.Query.InvalidRequest")
	}
	var rule *DetectionRule
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		rule, err = scan(row)
		return err
	}, query, args...)
	return rule, err
}

func prepareDetectionRulesQuery() (sq.SelectBuilder, func(*sql.Rows) ([]*DetectionRule, error)) {
	return sq.Select(
			DetectionRulesColumnID.identifier(),
			DetectionRulesColumnCreationDate.identifier(),
			DetectionRulesColumnChangeDate.identifier(),
			DetectionRulesColumnSequence.identifier(),
			DetectionRulesColumnDescription.identifier(),
			DetectionRulesColumnExpr.identifier(),
			DetectionRulesColumnEngine.identifier(),
			DetectionRulesColumnPriority.identifier(),
			DetectionRulesColumnStopOnMatch.identifier(),
			DetectionRulesColumnFindingName.identifier(),
			DetectionRulesColumnFindingMessage.identifier(),
			DetectionRulesColumnFindingBlock.identifier(),
			DetectionRulesColumnContextTemplate.identifier(),
			DetectionRulesColumnRateLimitKey.identifier(),
			DetectionRulesColumnRateLimitWindow.identifier(),
			DetectionRulesColumnRateLimitMax.identifier(),
		).From(detectionRulesTable.identifier()).OrderBy(DetectionRulesColumnID.identifier()).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) ([]*DetectionRule, error) {
			rules := make([]*DetectionRule, 0)
			for rows.Next() {
				rule := new(DetectionRule)
				if err := rows.Scan(
					&rule.ID,
					&rule.CreationDate,
					&rule.ChangeDate,
					&rule.Sequence,
					&rule.Description,
					&rule.Expr,
					&rule.Action,
					&rule.Priority,
					&rule.StopOnMatch,
					&rule.FindingName,
					&rule.FindingMessage,
					&rule.FindingBlock,
					&rule.ContextTemplate,
					&rule.RateLimitKey,
					&rule.RateLimitWindow,
					&rule.RateLimitMax,
				); err != nil {
					return nil, zerrors.ThrowInternal(err, "QUERY-3Y0H4", "Errors.Internal")
				}
				rules = append(rules, rule)
			}
			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-rD4gS", "Errors.Query.CloseRows")
			}
			return rules, nil
		}
}

func prepareDetectionRuleQuery() (sq.SelectBuilder, func(*sql.Row) (*DetectionRule, error)) {
	return sq.Select(
			DetectionRulesColumnID.identifier(),
			DetectionRulesColumnCreationDate.identifier(),
			DetectionRulesColumnChangeDate.identifier(),
			DetectionRulesColumnSequence.identifier(),
			DetectionRulesColumnDescription.identifier(),
			DetectionRulesColumnExpr.identifier(),
			DetectionRulesColumnEngine.identifier(),
			DetectionRulesColumnPriority.identifier(),
			DetectionRulesColumnStopOnMatch.identifier(),
			DetectionRulesColumnFindingName.identifier(),
			DetectionRulesColumnFindingMessage.identifier(),
			DetectionRulesColumnFindingBlock.identifier(),
			DetectionRulesColumnContextTemplate.identifier(),
			DetectionRulesColumnRateLimitKey.identifier(),
			DetectionRulesColumnRateLimitWindow.identifier(),
			DetectionRulesColumnRateLimitMax.identifier(),
		).From(detectionRulesTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*DetectionRule, error) {
			rule := new(DetectionRule)
			err := row.Scan(
				&rule.ID,
				&rule.CreationDate,
				&rule.ChangeDate,
				&rule.Sequence,
				&rule.Description,
				&rule.Expr,
				&rule.Action,
				&rule.Priority,
				&rule.StopOnMatch,
				&rule.FindingName,
				&rule.FindingMessage,
				&rule.FindingBlock,
				&rule.ContextTemplate,
				&rule.RateLimitKey,
				&rule.RateLimitWindow,
				&rule.RateLimitMax,
			)
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			if err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-p7dSh", "Errors.Internal")
			}
			return rule, nil
		}
}
