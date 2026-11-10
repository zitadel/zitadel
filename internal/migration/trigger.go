package migration

import (
	"context"
	"embed"
	"fmt"
	"reflect"
	"strings"
	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	countTriggerTmpl       = "count_trigger"
	deleteParentCountsTmpl = "delete_parent_counts_trigger"
)

var (
	//go:embed *.sql
	templateFS embed.FS
	templates  = template.Must(template.ParseFS(templateFS, "*.sql"))
)

// CountTrigger registers the existing projections.count_trigger function.
// The trigger than takes care of keeping count of existing
// rows in the source table.
// It also pre-populates the projections.resource_counts table with
// the counts for the given table.
//
// During the population of the resource_counts table,
// the source table is share-locked to prevent concurrent modifications.
// Projection handlers will be halted until the lock is released.
// SELECT statements are not blocked by the lock.
//
// This migration repeats when any of the arguments are changed,
// such as renaming of a projection table.
func CountTrigger(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
) RepeatableMigration {
	return CountTriggerConditional(
		db,
		table,
		parentType,
		instanceIDColumn,
		parentIDColumn,
		resource,
		false,
		nil,
	)
}

// CountTriggerConditional registers the existing projections.count_trigger function
// with conditions for specific column values and will only count rows that meet the conditions.
// The trigger than takes care of keeping count of existing
// rows in the source table.
// Additionally, if trackChange is true, the trigger will also keep track of
// updates to the rows that meet the conditions in case the values of the
// specified columns change.
// It also pre-populates the projections.resource_counts table with
// the counts for the given table.
//
// During the population of the resource_counts table,
// the source table is share-locked to prevent concurrent modifications.
// Projection handlers will be halted until the lock is released.
// SELECT statements are not blocked by the lock.
//
// This migration repeats when any of the arguments are changed,
// such as renaming of a projection table.
func CountTriggerConditional(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
	trackChange bool,
	conditions TriggerConditions,
) RepeatableMigration {
	return &triggerMigration{
		triggerConfig: triggerConfig{
			Table:            table,
			ParentType:       parentType.String(),
			InstanceIDColumn: instanceIDColumn,
			ParentIDColumn:   parentIDColumn,
			Resource:         resource,
			Conditions:       conditions,
			TrackChange:      trackChange,
		},
		db:           db,
		templateName: countTriggerTmpl,
	}
}

type TriggerConditions interface {
	ToSQL(table string, conditionsMet bool) string
}

type TriggerCondition struct {
	Column string `json:"column" mapstructure:"column"`
	Value  any    `json:"value" mapstructure:"value"`
}

// ToSQL implements the [TriggerConditions] interface.
// If conditionsMet is true, the SQL will be built to match when the condition is met.
// If conditionsMet is false, the SQL will be built to match when the condition is not met.
// e.g. col='value' vs col<>'value'
func (t TriggerCondition) ToSQL(table string, conditionsMet bool) string {
	value := fmt.Sprintf("%v", t.Value)
	if reflect.TypeOf(t.Value).Kind() == reflect.String {
		value = fmt.Sprintf("'%s'", t.Value)
	}

	operator := "="
	if !conditionsMet {
		operator = "<>"
	}

	return fmt.Sprintf("%s.%s %s %s", table, t.Column, operator, value)
}

type OrCondition struct {
	Conditions []TriggerCondition `json:"orConditions" mapstructure:"orConditions"`
}

// ToSQL implements the [TriggerConditions] interface.
// If conditionsMet is true, the SQL will be built to match when any of the conditions are met (OR).
// If conditionsMet is false, the SQL will be built to match when none of the conditions are met (AND).
// e.g. col1='value' OR col2='value' vs col1<>'value' AND col2<>'value'
func (t OrCondition) ToSQL(table string, conditionsMet bool) string {
	separator := " OR "
	if !conditionsMet {
		separator = " AND "
	}
	return toSQL(t.Conditions, table, separator, conditionsMet)
}

type AndCondition struct {
	Conditions []TriggerCondition `json:"andConditions" mapstructure:"andConditions"`
}

// ToSQL implements the [TriggerConditions] interface.
// If conditionsMet is true, the SQL will be built to check if all conditions are met (AND).
// If conditionsMet is false, the SQL will be built to check if any condition is not met (OR).
// e.g. col1='value' AND col2='value' vs col1<>'value' OR col2<>'value'
func (t AndCondition) ToSQL(table string, conditionsMet bool) string {
	separator := " AND "
	if !conditionsMet {
		separator = " OR "
	}
	return toSQL(t.Conditions, table, separator, conditionsMet)
}

func toSQL(conditions []TriggerCondition, table, separator string, conditionsMet bool) string {
	if len(conditions) == 0 {
		return ""
	}

	parts := make([]string, len(conditions))
	for i, condition := range conditions {
		parts[i] = condition.ToSQL(table, conditionsMet)
	}

	return "(" + strings.Join(parts, separator) + ")"
}

// DeleteParentCountsTrigger
//
// This migration repeats when any of the arguments are changed,
// such as renaming of a projection table.
func DeleteParentCountsTrigger(
	db *database.DB,
	table string,
	parentType domain.CountParentType,
	instanceIDColumn string,
	parentIDColumn string,
	resource string,
) RepeatableMigration {
	return &triggerMigration{
		triggerConfig: triggerConfig{
			Table:            table,
			ParentType:       parentType.String(),
			InstanceIDColumn: instanceIDColumn,
			ParentIDColumn:   parentIDColumn,
			Resource:         resource,
		},
		db:           db,
		templateName: deleteParentCountsTmpl,
	}
}

type triggerMigration struct {
	triggerConfig
	db           *database.DB
	templateName string
}

// String implements [Migration] and [fmt.Stringer].
func (m *triggerMigration) String() string {
	return fmt.Sprintf("repeatable_%s_%s", m.Resource, m.templateName)
}

// Execute implements [Migration]
func (m *triggerMigration) Execute(ctx context.Context, _ eventstore.Event) error {
	var query strings.Builder
	err := templates.ExecuteTemplate(&query, m.templateName, m.triggerConfig)
	if err != nil {
		return fmt.Errorf("%s: execute trigger template: %w", m, err)
	}
	_, err = m.db.ExecContext(ctx, query.String())
	if err != nil {
		return fmt.Errorf("%s: exec trigger query: %w", m, err)
	}
	return nil
}

type triggerConfig struct {
	Table            string            `json:"table,omitempty" mapstructure:"table"`
	ParentType       string            `json:"parent_type,omitempty" mapstructure:"parent_type"`
	InstanceIDColumn string            `json:"instance_id_column,omitempty" mapstructure:"instance_id_column"`
	ParentIDColumn   string            `json:"parent_id_column,omitempty" mapstructure:"parent_id_column"`
	Resource         string            `json:"resource,omitempty" mapstructure:"resource"`
	Conditions       TriggerConditions `json:"conditions,omitempty" mapstructure:"conditions"`
	TrackChange      bool              `json:"track_change,omitempty" mapstructure:"track_change"`
}

// Check implements [RepeatableMigration].
func (c *triggerConfig) Check(lastRun map[string]any) bool {
	var dst triggerConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       DecodeTriggerConditionsHook(),
		WeaklyTypedInput: true,
		Result:           &dst,
	})
	if err != nil {
		panic(err)
	}
	if err = decoder.Decode(lastRun); err != nil {
		return true
	}
	return !reflect.DeepEqual(dst, *c)
}

// DecodeTriggerConditionsHook returns a mapstructure.DecodeHookFunc that can decode
// a map into the correct concrete type implementing [TriggerConditions].
func DecodeTriggerConditionsHook() mapstructure.DecodeHookFunc {
	return func(
		from reflect.Type,
		to reflect.Type,
		data interface{},
	) (interface{}, error) {
		if to != reflect.TypeOf((*TriggerConditions)(nil)).Elem() {
			return data, nil
		}

		mapData, ok := data.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected a map for TriggerConditions, but got %T", data)
		}

		var result TriggerConditions

		if _, ok := mapData["orConditions"]; ok {
			result = &OrCondition{}
		} else if _, ok := mapData["andConditions"]; ok {
			result = &AndCondition{}
		} else if _, ok := mapData["column"]; ok {
			result = &TriggerCondition{}
		} else {
			return data, nil
		}

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:     result,
			DecodeHook: DecodeTriggerConditionsHook(),
		})
		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(mapData); err != nil {
			return nil, err
		}

		return result, nil
	}
}
