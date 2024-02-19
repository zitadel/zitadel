package zerrors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplatableError_GetVars(t *testing.T) {
	type fields struct {
		withVariables func(error) error
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{{
		name: "variables are available",
		fields: fields{
			withVariables: func(err error) error {
				return WithTemplateVariables(err, map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				})
			},
		},
		want: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
	}, {
		name: "variables in the chain are merged",
		fields: fields{
			withVariables: func(err error) error {
				return WithTemplateVariables(WithTemplateVariables(err, map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				}), map[string]interface{}{
					"key3": "value3",
				})
			},
		},
		want: map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}, {
		name: "variables are overwritten by children",
		fields: fields{
			withVariables: func(err error) error {
				return WithTemplateVariables(WithTemplateVariables(err, map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				}), map[string]interface{}{
					"key2": "value3",
				})
			},
		},
		want: map[string]interface{}{
			"key1": "value1",
			"key2": "value3",
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templErr := new(TemplatableError)
			require.True(t, errors.As(tt.fields.withVariables(errors.New("test")), &templErr))
			require.Equal(t, "test", templErr.Error(), "templErr error should not print anything")
			require.Equal(t, templErr.GetVars(), tt.want, "templErr.GetVars()")
		})
	}
}
