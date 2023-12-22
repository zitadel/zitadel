package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "org reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyAddedEventType,
						org.AggregateType,
						[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b", "themeMode": 1}`),
					), org.LabelPolicyAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.label_policies3 (creation_date, change_date, sequence, id, state, is_default, resource_owner, instance_id, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled, theme_mode) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.LabelPolicyStatePreview,
								false,
								"ro-id",
								"instance-id",
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
								domain.LabelPolicyThemeLight,
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceChanged",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyChangedEventType,
						org.AggregateType,
						[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b", "themeMode": 1}`),
					), org.LabelPolicyChangedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_primary_color, light_background_color, light_warn_color, light_font_color, theme_mode) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (state = $9) AND (instance_id = $10)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"#5282c1",
								"#141735",
								"#ff3b5b",
								"#ffffff",
								domain.LabelPolicyThemeLight,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.LabelPolicyRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.label_policies3 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(LabelPolicyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.label_policies3 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceActivated",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyActivatedEventType,
						org.AggregateType,
						nil,
					), org.LabelPolicyActivatedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.label_policies3 (change_date, sequence, state, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode) SELECT $1, $2, $3, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode FROM projections.label_policies3 AS copy_table WHERE (copy_table.id = $4) AND (copy_table.state = $5) AND (copy_table.instance_id = $6) ON CONFLICT (instance_id, id, state) DO UPDATE SET (change_date, sequence, state, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode) = ($1, $2, $3, EXCLUDED.creation_date, EXCLUDED.resource_owner, EXCLUDED.instance_id, EXCLUDED.id, EXCLUDED.is_default, EXCLUDED.hide_login_name_suffix, EXCLUDED.font_url, EXCLUDED.watermark_disabled, EXCLUDED.should_error_popup, EXCLUDED.light_primary_color, EXCLUDED.light_warn_color, EXCLUDED.light_background_color, EXCLUDED.light_font_color, EXCLUDED.light_logo_url, EXCLUDED.light_icon_url, EXCLUDED.dark_primary_color, EXCLUDED.dark_warn_color, EXCLUDED.dark_background_color, EXCLUDED.dark_font_color, EXCLUDED.dark_logo_url, EXCLUDED.dark_icon_url, EXCLUDED.theme_mode)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.LabelPolicyStateActive,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLogoAdded light",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyLogoAddedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), org.LabelPolicyLogoAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLogoAdded dark",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyLogoDarkAddedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), org.LabelPolicyLogoDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceIconAdded light",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyIconAddedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), org.LabelPolicyIconAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceIconAdded dark",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyIconDarkAddedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), org.LabelPolicyIconDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLogoRemoved light",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyLogoRemovedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), org.LabelPolicyLogoRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceLogoRemoved dark",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyLogoDarkRemovedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), org.LabelPolicyLogoDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceIconRemoved light",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyIconRemovedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), org.LabelPolicyIconRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceIconRemoved dark",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyIconDarkRemovedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), org.LabelPolicyIconDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceFontAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyFontAddedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/font.ttf"}`),
					), org.LabelPolicyFontAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/font.ttf",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceFontRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyFontRemovedEventType,
						org.AggregateType,
						[]byte(`{"storeKey": "/path/to/font.ttf"}`),
					), org.LabelPolicyFontRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org reduceAssetsRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.LabelPolicyAssetsRemovedEventType,
						org.AggregateType,
						nil,
					), org.LabelPolicyAssetsRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAssetsRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url, light_icon_url, dark_logo_url, dark_icon_url, font_url) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (state = $9) AND (instance_id = $10)",
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
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b", "themeMode": 1}`),
					), instance.LabelPolicyAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.label_policies3 (creation_date, change_date, sequence, id, state, is_default, resource_owner, instance_id, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled, theme_mode) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.LabelPolicyStatePreview,
								true,
								"ro-id",
								"instance-id",
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
								domain.LabelPolicyThemeLight,
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{"backgroundColor": "#141735", "fontColor": "#ffffff", "primaryColor": "#5282c1", "warnColor": "#ff3b5b", "primaryColorDark": "#ffffff","backgroundColorDark": "#ffffff", "warnColorDark": "#ffffff", "fontColorDark": "#ffffff", "hideLoginNameSuffix": true, "errorMsgPopup": true, "disableWatermark": true, "themeMode": 1}`),
					), instance.LabelPolicyChangedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_primary_color, light_background_color, light_warn_color, light_font_color, dark_primary_color, dark_background_color, dark_warn_color, dark_font_color, hide_login_name_suffix, should_error_popup, watermark_disabled, theme_mode) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) WHERE (id = $15) AND (state = $16) AND (instance_id = $17)",
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
								domain.LabelPolicyThemeLight,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceActivated",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyActivatedEventType,
						instance.AggregateType,
						nil,
					), instance.LabelPolicyActivatedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.label_policies3 (change_date, sequence, state, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode) SELECT $1, $2, $3, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode FROM projections.label_policies3 AS copy_table WHERE (copy_table.id = $4) AND (copy_table.state = $5) AND (copy_table.instance_id = $6) ON CONFLICT (instance_id, id, state) DO UPDATE SET (change_date, sequence, state, creation_date, resource_owner, instance_id, id, is_default, hide_login_name_suffix, font_url, watermark_disabled, should_error_popup, light_primary_color, light_warn_color, light_background_color, light_font_color, light_logo_url, light_icon_url, dark_primary_color, dark_warn_color, dark_background_color, dark_font_color, dark_logo_url, dark_icon_url, theme_mode) = ($1, $2, $3, EXCLUDED.creation_date, EXCLUDED.resource_owner, EXCLUDED.instance_id, EXCLUDED.id, EXCLUDED.is_default, EXCLUDED.hide_login_name_suffix, EXCLUDED.font_url, EXCLUDED.watermark_disabled, EXCLUDED.should_error_popup, EXCLUDED.light_primary_color, EXCLUDED.light_warn_color, EXCLUDED.light_background_color, EXCLUDED.light_font_color, EXCLUDED.light_logo_url, EXCLUDED.light_icon_url, EXCLUDED.dark_primary_color, EXCLUDED.dark_warn_color, EXCLUDED.dark_background_color, EXCLUDED.dark_font_color, EXCLUDED.dark_logo_url, EXCLUDED.dark_icon_url, EXCLUDED.theme_mode)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.LabelPolicyStateActive,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLogoAdded light",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyLogoAddedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), instance.LabelPolicyLogoAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLogoAdded dark",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyLogoDarkAddedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), instance.LabelPolicyLogoDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/logo.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceIconAdded light",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyIconAddedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), instance.LabelPolicyIconAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceIconAdded dark",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyIconDarkAddedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), instance.LabelPolicyIconDarkAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/icon.png",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLogoRemoved light",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyLogoRemovedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), instance.LabelPolicyLogoRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLogoRemoved dark",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyLogoDarkRemovedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/logo.png"}`),
					), instance.LabelPolicyLogoDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceLogoRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_logo_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceIconRemoved light",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyIconRemovedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), instance.LabelPolicyIconRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceIconRemoved dark",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyIconDarkRemovedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/icon.png"}`),
					), instance.LabelPolicyIconDarkRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceIconRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, dark_icon_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceFontAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyFontAddedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/font.ttf"}`),
					), instance.LabelPolicyFontAddedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"/path/to/font.ttf",
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceFontRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyFontRemovedEventType,
						instance.AggregateType,
						[]byte(`{"storeKey": "/path/to/font.ttf"}`),
					), instance.LabelPolicyFontRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceFontRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, font_url) = ($1, $2, $3) WHERE (id = $4) AND (state = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								nil,
								"agg-id",
								domain.LabelPolicyStatePreview,
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceAssetsRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.LabelPolicyAssetsRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.LabelPolicyAssetsRemovedEventMapper),
			},
			reduce: (&labelPolicyProjection{}).reduceAssetsRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.label_policies3 SET (change_date, sequence, light_logo_url, light_icon_url, dark_logo_url, dark_icon_url, font_url) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (state = $9) AND (instance_id = $10)",
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
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&labelPolicyProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.label_policies3 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, LabelPolicyTable, tt.want)
		})
	}
}
