package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestLabelPolicyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "org.reduceAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b"}`),
				), org.LabelPolicyAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.label_policies (creation_date, change_date, sequence, id, state, is_default, resource_owner, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.LabelPolicyStatePreview,
								false,
								"ro-id",
								"#5282c1",
								"#141735",
								"#ff3b5b",
								"#ffffff",
								"",
								"",
								"",
								"",
								false,
								false,
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyChangedEventType),
					org.AggregateType,
					[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b"}`),
				), org.LabelPolicyChangedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_primary_color, light_background_color, light_warn_color, light_font_color) = ($1, $2, $3, $4, $5, $6) WHERE (id = $7) AND (state = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"#5282c1",
								"#141735",
								"#ff3b5b",
								"#ffffff",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.LabelPolicyRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.label_policies WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceActivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyActivatedEventType),
					org.AggregateType,
					nil,
				), org.LabelPolicyActivatedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceActivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.label_policies (change_date, sequence, state, creation_date, resource_owner, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url) SELECT $1, $2, $3, creation_date, resource_owner, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url FROM zitadel.projections.label_policies AS copy_table WHERE copy_table.id = $4 AND copy_table.state = $5",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.LabelPolicyStateActive,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceLogoAdded light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyLogoAddedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), org.LabelPolicyLogoAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceLogoAdded dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyLogoDarkAddedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), org.LabelPolicyLogoDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIconAdded light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyIconAddedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), org.LabelPolicyIconAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIconAdded dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyIconDarkAddedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), org.LabelPolicyIconDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceLogoRemoved light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyLogoRemovedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), org.LabelPolicyLogoRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceLogoRemoved dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyLogoDarkRemovedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), org.LabelPolicyLogoDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIconRemoved light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyIconRemovedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), org.LabelPolicyIconRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceIconRemoved dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyIconDarkRemovedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), org.LabelPolicyIconDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceFontAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyFontAddedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/font.ttf"}`),
				), org.LabelPolicyFontAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/font.ttf",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceFontRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyFontRemovedEventType),
					org.AggregateType,
					[]byte(`{"storeKey": "/path/to/font.ttf"}`),
				), org.LabelPolicyFontRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "org.reduceAssetsRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LabelPolicyAssetsRemovedEventType),
					org.AggregateType,
					nil,
				), org.LabelPolicyAssetsRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAssetsRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url, light_icon_url, dark_logo_url, dark_icon_url, font_url) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (state = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								nil,
								nil,
								nil,
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b"}`),
				), iam.LabelPolicyAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.label_policies (creation_date, change_date, sequence, id, state, is_default, resource_owner, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.LabelPolicyStatePreview,
								true,
								"ro-id",
								"#5282c1",
								"#141735",
								"#ff3b5b",
								"#ffffff",
								"",
								"",
								"",
								"",
								false,
								false,
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyChangedEventType),
					iam.AggregateType,
					[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b", "primaryColorDark": "#ffffff","backgroundColorDark": "#ffffff", "warnColorDark": "#ffffff", "fontColorDark": "#ffffff", "hideLoginNameSuffix": true, "errorMsgPopup": true, "disableWatermark": true}`),
				), iam.LabelPolicyChangedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) WHERE (id = $14) AND (state = $15)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"#5282c1",
								"#141735",
								"#ff3b5b",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								"#ffffff",
								true,
								true,
								true,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceActivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyActivatedEventType),
					iam.AggregateType,
					nil,
				), iam.LabelPolicyActivatedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceActivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPSERT INTO zitadel.projections.label_policies (change_date, sequence, state, creation_date, resource_owner, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url) SELECT $1, $2, $3, creation_date, resource_owner, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url FROM zitadel.projections.label_policies AS copy_table WHERE copy_table.id = $4 AND copy_table.state = $5",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.LabelPolicyStateActive,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceLogoAdded light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyLogoAddedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), iam.LabelPolicyLogoAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceLogoAdded dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyLogoDarkAddedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), iam.LabelPolicyLogoDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceIconAdded light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyIconAddedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), iam.LabelPolicyIconAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceIconAdded dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyIconDarkAddedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), iam.LabelPolicyIconDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceLogoRemoved light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyLogoRemovedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), iam.LabelPolicyLogoRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceLogoRemoved dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyLogoDarkRemovedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/logo.png"}`),
				), iam.LabelPolicyLogoDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceIconRemoved light",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyIconRemovedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), iam.LabelPolicyIconRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceIconRemoved dark",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyIconDarkRemovedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/icon.png"}`),
				), iam.LabelPolicyIconDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceFontAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyFontAddedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/font.ttf"}`),
				), iam.LabelPolicyFontAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/font.ttf",
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceFontRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyFontRemovedEventType),
					iam.AggregateType,
					[]byte(`{"storeKey": "/path/to/font.ttf"}`),
				), iam.LabelPolicyFontRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceAssetsRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LabelPolicyAssetsRemovedEventType),
					iam.AggregateType,
					nil,
				), iam.LabelPolicyAssetsRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAssetsRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LabelPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.label_policies SET (change_date, sequence, light_logo_url, light_icon_url, dark_logo_url, dark_icon_url, font_url) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (state = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								nil,
								nil,
								nil,
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
