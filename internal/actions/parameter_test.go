package actions

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"
)

func TestSetFields(t *testing.T) {
	testFn := func(a string) { fmt.Println(a) }
	tests := []struct {
		name        string
		setFields   FieldOption
		want        fields
		shouldPanic bool
	}{
		{
			name:      "field is simple value",
			setFields: SetFields("value", 5),
			want: fields{
				"value": 5,
			},
		},
		{
			name:      "field is method",
			setFields: SetFields("value", testFn),
			want: fields{
				"value": testFn,
			},
		},
		{
			name:      "field without value",
			setFields: SetFields("value"),
			want:      fields{},
		},
		{
			name:      "field with empty value",
			setFields: SetFields("value", ""),
			want: fields{
				"value": "",
			},
		},
		{
			name: "nested simple value",
			setFields: SetFields(
				"field",
				SetFields("sub", 5),
			),
			want: fields{
				"field": fields{
					"sub": 5,
				},
			},
		},
		{
			name: "nested multiple fields",
			setFields: SetFields(
				"field",
				SetFields("sub1", 5),
				SetFields("sub2", "asdf"),
				SetFields("sub3", testFn),
			),
			want: fields{
				"field": fields{
					"sub1": 5,
					"sub2": "asdf",
					"sub3": testFn,
				},
			},
		},
		{
			name: "try to overwrite field primitives",
			setFields: SetFields(
				"field",
				SetFields("sub", 5),
				SetFields("sub", testFn),
			),
			shouldPanic: true,
		},
		{
			name: "try to overwrite primitives with fields",
			setFields: SetFields(
				"field",
				SetFields("sub", 5),
				SetFields("sub", SetFields("please", "panic")),
			),
			shouldPanic: true,
		},
		{
			name: "try to overwrite fields with primitive",
			setFields: SetFields(
				"field",
				SetFields("sub", SetFields("please", "panic")),
				SetFields("sub", 5),
			),
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			panicked := false
			if tt.shouldPanic {
				defer func() {
					if panicked != tt.shouldPanic {
						t.Errorf("wanted panic: %v got %v", tt.shouldPanic, panicked)
					}
				}()
				defer func() {
					recover()
					panicked = true
				}()
			}
			config := &FieldConfig{
				Runtime: goja.New(),
				fields:  fields{},
			}
			tt.setFields(config)
			if !tt.shouldPanic && fmt.Sprint(config.fields) != fmt.Sprint(tt.want) {
				t.Errorf("SetFields() = %v, want %v", fmt.Sprint(config.fields), fmt.Sprint(tt.want))
			}
		})
	}
}
