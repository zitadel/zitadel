package query

import (
"context"
"database/sql"
"errors"
"time"

sq "github.com/Masterminds/squirrel"

"github.com/zitadel/zitadel/internal/api/authz"
"github.com/zitadel/zitadel/internal/database"
"github.com/zitadel/zitadel/internal/query/projection"
"github.com/zitadel/zitadel/internal/zerrors"
)

var (
detectionSettingsTable = table{
name:          projection.DetectionSettingsProjectionTable,
instanceIDCol: projection.DetectionSettingsColumnInstanceID,
}
DetectionSettingsColumnCreationDate = Column{name: projection.DetectionSettingsColumnCreationDate, table: detectionSettingsTable}
DetectionSettingsColumnChangeDate = Column{name: projection.DetectionSettingsColumnChangeDate, table: detectionSettingsTable}
DetectionSettingsColumnInstanceID = Column{name: projection.DetectionSettingsColumnInstanceID, table: detectionSettingsTable}
DetectionSettingsColumnSequence = Column{name: projection.DetectionSettingsColumnSequence, table: detectionSettingsTable}
DetectionSettingsColumnEnabled = Column{name: projection.DetectionSettingsColumnEnabled, table: detectionSettingsTable}
DetectionSettingsColumnFailOpen = Column{name: projection.DetectionSettingsColumnFailOpen, table: detectionSettingsTable}
DetectionSettingsColumnFailureBurstThreshold = Column{name: projection.DetectionSettingsColumnFailureBurstThreshold, table: detectionSettingsTable}
DetectionSettingsColumnHistoryWindow = Column{name: projection.DetectionSettingsColumnHistoryWindow, table: detectionSettingsTable}
DetectionSettingsColumnContextChangeWindow = Column{name: projection.DetectionSettingsColumnContextChangeWindow, table: detectionSettingsTable}
DetectionSettingsColumnMaxSignalsPerUser = Column{name: projection.DetectionSettingsColumnMaxSignalsPerUser, table: detectionSettingsTable}
DetectionSettingsColumnMaxSignalsPerSession = Column{name: projection.DetectionSettingsColumnMaxSignalsPerSession, table: detectionSettingsTable}
DetectionSettingsColumnRulesManaged = Column{name: projection.DetectionSettingsColumnRulesManaged, table: detectionSettingsTable}
)

type DetectionSettings struct {
CreationDate          time.Time
ChangeDate            time.Time
InstanceID            string
Sequence              uint64
Enabled               bool
FailOpen              bool
FailureBurstThreshold int64
HistoryWindow         database.Duration
ContextChangeWindow   database.Duration
MaxSignalsPerUser     int64
MaxSignalsPerSession  int64
RulesManaged          bool
}

func (q *Queries) DetectionSettings(ctx context.Context) (_ *DetectionSettings, err error) {
stmt, scan := prepareDetectionSettingsQuery()
query, args, err := stmt.Where(sq.Eq{DetectionSettingsColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}).ToSql()
if err != nil {
return nil, zerrors.ThrowInternal(err, "QUERY-R4kmJ", "Errors.Query.SQLStatement")
}
var settings *DetectionSettings
err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
settings, err = scan(row)
return err
}, query, args...)
return settings, err
}

func prepareDetectionSettingsQuery() (sq.SelectBuilder, func(*sql.Row) (*DetectionSettings, error)) {
return sq.Select(
DetectionSettingsColumnCreationDate.identifier(),
DetectionSettingsColumnChangeDate.identifier(),
DetectionSettingsColumnInstanceID.identifier(),
DetectionSettingsColumnSequence.identifier(),
DetectionSettingsColumnEnabled.identifier(),
DetectionSettingsColumnFailOpen.identifier(),
DetectionSettingsColumnFailureBurstThreshold.identifier(),
DetectionSettingsColumnHistoryWindow.identifier(),
DetectionSettingsColumnContextChangeWindow.identifier(),
DetectionSettingsColumnMaxSignalsPerUser.identifier(),
DetectionSettingsColumnMaxSignalsPerSession.identifier(),
DetectionSettingsColumnRulesManaged.identifier(),
).From(detectionSettingsTable.identifier()).PlaceholderFormat(sq.Dollar),
func(row *sql.Row) (*DetectionSettings, error) {
settings := new(DetectionSettings)
err := row.Scan(
&settings.CreationDate,
&settings.ChangeDate,
&settings.InstanceID,
&settings.Sequence,
&settings.Enabled,
&settings.FailOpen,
&settings.FailureBurstThreshold,
&settings.HistoryWindow,
&settings.ContextChangeWindow,
&settings.MaxSignalsPerUser,
&settings.MaxSignalsPerSession,
&settings.RulesManaged,
)
if errors.Is(err, sql.ErrNoRows) {
return nil, nil
}
if err != nil {
return nil, zerrors.ThrowInternal(err, "QUERY-V0ah5", "Errors.Internal")
}
return settings, nil
}
}
