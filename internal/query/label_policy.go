package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type LabelPolicy struct {
	ID                  string
	CreationDate        time.Time
	ChangeDate          time.Time
	Sequence            uint64
	State               domain.LabelPolicyState
	IsDefault           bool
	ResourceOwner       string
	HideLoginNameSuffix bool
	FontURL             string
	WatermarkDisabled   bool
	ShouldErrorPopup    bool

	Dark  Theme
	Light Theme
}

type Theme struct {
	PrimaryColor    string
	WarnColor       string
	BackgroundColor string
	FontColor       string
	LogoURL         string
	IconURL         string
}

func (q *Queries) ActiveLabelPolicyByOrg(ctx context.Context, orgID string, withOwnerRemoved bool) (_ *LabelPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLabelPolicyQuery(ctx, q.client)
	eq := sq.Eq{
		LabelPolicyColState.identifier():      domain.LabelPolicyStateActive,
		LabelPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[LabelPolicyOwnerRemoved.identifier()] = false
	}
	query, args, err := stmt.Where(
		sq.And{
			sq.Or{
				sq.Eq{LabelPolicyColID.identifier(): orgID},
				sq.Eq{LabelPolicyColID.identifier(): authz.GetInstance(ctx).InstanceID()},
			},
			eq,
		}).
		OrderBy(LabelPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-V22un", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) PreviewLabelPolicyByOrg(ctx context.Context, orgID string) (_ *LabelPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLabelPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(
		sq.And{
			sq.Or{
				sq.Eq{
					LabelPolicyColID.identifier(): orgID,
				},
				sq.Eq{
					LabelPolicyColID.identifier(): authz.GetInstance(ctx).InstanceID(),
				},
			},
			sq.Eq{
				LabelPolicyColState.identifier():      domain.LabelPolicyStatePreview,
				LabelPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
		}).
		OrderBy(LabelPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AG5eq", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultActiveLabelPolicy(ctx context.Context) (_ *LabelPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLabelPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		LabelPolicyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		LabelPolicyColState.identifier():      domain.LabelPolicyStateActive,
		LabelPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(LabelPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-mN0Ci", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) DefaultPreviewLabelPolicy(ctx context.Context) (_ *LabelPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareLabelPolicyQuery(ctx, q.client)
	query, args, err := stmt.Where(sq.Eq{
		LabelPolicyColID.identifier():         authz.GetInstance(ctx).InstanceID(),
		LabelPolicyColState.identifier():      domain.LabelPolicyStatePreview,
		LabelPolicyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}).
		OrderBy(LabelPolicyColIsDefault.identifier()).
		Limit(1).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-B3JQR", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

var (
	labelPolicyTable = table{
		name:          projection.LabelPolicyTable,
		instanceIDCol: projection.LabelPolicyInstanceIDCol,
	}
	LabelPolicyColCreationDate = Column{
		name: projection.LabelPolicyCreationDateCol,
	}
	LabelPolicyColChangeDate = Column{
		name: projection.LabelPolicyChangeDateCol,
	}
	LabelPolicyColSequence = Column{
		name: projection.LabelPolicySequenceCol,
	}
	LabelPolicyColID = Column{
		name: projection.LabelPolicyIDCol,
	}
	LabelPolicyColState = Column{
		name: projection.LabelPolicyStateCol,
	}
	LabelPolicyColIsDefault = Column{
		name: projection.LabelPolicyIsDefaultCol,
	}
	LabelPolicyColResourceOwner = Column{
		name: projection.LabelPolicyResourceOwnerCol,
	}
	LabelPolicyColInstanceID = Column{
		name: projection.LabelPolicyInstanceIDCol,
	}
	LabelPolicyColHideLoginNameSuffix = Column{
		name: projection.LabelPolicyHideLoginNameSuffixCol,
	}
	LabelPolicyColFontURL = Column{
		name: projection.LabelPolicyFontURLCol,
	}
	LabelPolicyColWatermarkDisabled = Column{
		name: projection.LabelPolicyWatermarkDisabledCol,
	}
	LabelPolicyColShouldErrorPopup = Column{
		name: projection.LabelPolicyShouldErrorPopupCol,
	}
	LabelPolicyColLightPrimaryColor = Column{
		name: projection.LabelPolicyLightPrimaryColorCol,
	}
	LabelPolicyColLightWarnColor = Column{
		name: projection.LabelPolicyLightWarnColorCol,
	}
	LabelPolicyColLightBackgroundColor = Column{
		name: projection.LabelPolicyLightBackgroundColorCol,
	}
	LabelPolicyColLightFontColor = Column{
		name: projection.LabelPolicyLightFontColorCol,
	}
	LabelPolicyColLightLogoURL = Column{
		name: projection.LabelPolicyLightLogoURLCol,
	}
	LabelPolicyColLightIconURL = Column{
		name: projection.LabelPolicyLightIconURLCol,
	}
	LabelPolicyColDarkPrimaryColor = Column{
		name: projection.LabelPolicyDarkPrimaryColorCol,
	}
	LabelPolicyColDarkWarnColor = Column{
		name: projection.LabelPolicyDarkWarnColorCol,
	}
	LabelPolicyColDarkBackgroundColor = Column{
		name: projection.LabelPolicyDarkBackgroundColorCol,
	}
	LabelPolicyColDarkFontColor = Column{
		name: projection.LabelPolicyDarkFontColorCol,
	}
	LabelPolicyColDarkLogoURL = Column{
		name: projection.LabelPolicyDarkLogoURLCol,
	}
	LabelPolicyColDarkIconURL = Column{
		name: projection.LabelPolicyDarkIconURLCol,
	}
	LabelPolicyOwnerRemoved = Column{
		name: projection.LabelPolicyOwnerRemovedCol,
	}
)

func prepareLabelPolicyQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*LabelPolicy, error)) {
	return sq.Select(
			LabelPolicyColCreationDate.identifier(),
			LabelPolicyColChangeDate.identifier(),
			LabelPolicyColSequence.identifier(),
			LabelPolicyColID.identifier(),
			LabelPolicyColState.identifier(),
			LabelPolicyColIsDefault.identifier(),
			LabelPolicyColResourceOwner.identifier(),

			LabelPolicyColHideLoginNameSuffix.identifier(),
			LabelPolicyColFontURL.identifier(),
			LabelPolicyColWatermarkDisabled.identifier(),
			LabelPolicyColShouldErrorPopup.identifier(),

			LabelPolicyColLightPrimaryColor.identifier(),
			LabelPolicyColLightWarnColor.identifier(),
			LabelPolicyColLightBackgroundColor.identifier(),
			LabelPolicyColLightFontColor.identifier(),
			LabelPolicyColLightLogoURL.identifier(),
			LabelPolicyColLightIconURL.identifier(),

			LabelPolicyColDarkPrimaryColor.identifier(),
			LabelPolicyColDarkWarnColor.identifier(),
			LabelPolicyColDarkBackgroundColor.identifier(),
			LabelPolicyColDarkFontColor.identifier(),
			LabelPolicyColDarkLogoURL.identifier(),
			LabelPolicyColDarkIconURL.identifier(),
		).
			From(labelPolicyTable.identifier() + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*LabelPolicy, error) {
			policy := new(LabelPolicy)

			var (
				fontURL              = sql.NullString{}
				lightPrimaryColor    = sql.NullString{}
				lightWarnColor       = sql.NullString{}
				lightBackgroundColor = sql.NullString{}
				lightFontColor       = sql.NullString{}
				lightLogoURL         = sql.NullString{}
				lightIconURL         = sql.NullString{}
				darkPrimaryColor     = sql.NullString{}
				darkWarnColor        = sql.NullString{}
				darkBackgroundColor  = sql.NullString{}
				darkFontColor        = sql.NullString{}
				darkLogoURL          = sql.NullString{}
				darkIconURL          = sql.NullString{}
			)

			err := row.Scan(
				&policy.CreationDate,
				&policy.ChangeDate,
				&policy.Sequence,
				&policy.ID,
				&policy.State,
				&policy.IsDefault,
				&policy.ResourceOwner,

				&policy.HideLoginNameSuffix,
				&fontURL,
				&policy.WatermarkDisabled,
				&policy.ShouldErrorPopup,

				&lightPrimaryColor,
				&lightWarnColor,
				&lightBackgroundColor,
				&lightFontColor,
				&lightLogoURL,
				&lightIconURL,

				&darkPrimaryColor,
				&darkWarnColor,
				&darkBackgroundColor,
				&darkFontColor,
				&darkLogoURL,
				&darkIconURL,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-bJEsm", "Errors.Org.PolicyNotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-awLM6", "Errors.Internal")
			}

			policy.FontURL = fontURL.String
			policy.Light.PrimaryColor = lightPrimaryColor.String
			policy.Light.WarnColor = lightWarnColor.String
			policy.Light.BackgroundColor = lightBackgroundColor.String
			policy.Light.FontColor = lightFontColor.String
			policy.Light.LogoURL = lightLogoURL.String
			policy.Light.IconURL = lightIconURL.String
			policy.Dark.PrimaryColor = darkPrimaryColor.String
			policy.Dark.WarnColor = darkWarnColor.String
			policy.Dark.BackgroundColor = darkBackgroundColor.String
			policy.Dark.FontColor = darkFontColor.String
			policy.Dark.LogoURL = darkLogoURL.String
			policy.Dark.IconURL = darkIconURL.String

			return policy, nil
		}
}

func (p *LabelPolicy) ToDomain() *domain.LabelPolicy {
	return &domain.LabelPolicy{
		Default:             p.IsDefault,
		PrimaryColor:        p.Light.PrimaryColor,
		BackgroundColor:     p.Light.BackgroundColor,
		WarnColor:           p.Light.WarnColor,
		FontColor:           p.Light.FontColor,
		LogoURL:             p.Light.LogoURL,
		IconURL:             p.Light.IconURL,
		PrimaryColorDark:    p.Dark.PrimaryColor,
		BackgroundColorDark: p.Dark.BackgroundColor,
		WarnColorDark:       p.Dark.WarnColor,
		FontColorDark:       p.Dark.FontColor,
		LogoDarkURL:         p.Dark.LogoURL,
		IconDarkURL:         p.Dark.IconURL,
		Font:                p.FontURL,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		ErrorMsgPopup:       p.ShouldErrorPopup,
		DisableWatermark:    p.WatermarkDisabled,
	}
}
