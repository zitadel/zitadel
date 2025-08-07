package test

import (
	"sort"
	"strings"
	"testing"
	"time"
)

func TestCheckAllFieldsIgnores(t *testing.T) {
	type args struct {
		object        interface{}
		ignoredFields []string
	}
	type res struct {
		mock *mappedExpecter
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "simple struct",
			args: args{
				object: &struct {
					Company string
				}{
					Company: "caos AG",
				},
				ignoredFields: []string{"Company"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Company"},
					nil,
				),
			},
		},
		{
			name: "simple struct with private",
			args: args{
				object: &struct {
					Company string
					priv    bool
				}{
					Company: "caos AG",
					priv:    true,
				},
				ignoredFields: []string{"Company"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Company"},
					nil,
				),
			},
		},
		{
			name: "simple struct length",
			args: args{
				object: &struct {
					Company map[string]string
					priv    bool
				}{
					Company: map[string]string{"caos AG": "ZITADEL"},
					priv:    true,
				},
				ignoredFields: []string{"Company"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Company"},
					nil,
				),
			},
		},
		{
			name: "ignore empty nested field",
			args: args{
				object: &struct {
					Company *struct {
						Name    string
						Founded time.Time
					}
					priv bool
				}{
					Company: &struct {
						Name    string
						Founded time.Time
					}{},
					priv: true,
				},
				ignoredFields: []string{"Company"},
			},

			res: res{
				mock: newMappedExpeter(
					nil,
					nil,
				),
			},
		},
		{
			name: "ignore nested field",
			args: args{
				object: &struct {
					Company *struct {
						Name    string
						Founded time.Time
					}
					priv bool
				}{
					Company: &struct {
						Name    string
						Founded time.Time
					}{
						Founded: time.Date(2019, time.April, 1, 1, 1, 1, 1, time.Local),
					},
					priv: true,
				},
				ignoredFields: []string{"Company.Founded"},
			},

			res: res{
				mock: newMappedExpeter(
					[]string{"Company.Founded"},
					[]string{"Company.Name"},
				),
			},
		},
		{
			name: "ignore nested fields",
			args: args{
				object: &struct {
					Company *struct {
						Name    string
						Founded time.Time
					}
					priv bool
				}{
					Company: &struct {
						Name    string
						Founded time.Time
					}{
						Name:    "caos AG",
						Founded: time.Date(2019, time.April, 1, 1, 1, 1, 1, time.Local),
					},
					priv: true,
				},
				ignoredFields: []string{"Company"},
			},

			res: res{
				mock: newMappedExpeter(
					[]string{"Company.Name", "Company.Founded"},
					nil,
				),
			},
		},
		{
			name: "ignore nested struct first field filled",
			args: args{
				object: &struct {
					Nested
					Priv bool
				}{
					Nested: Nested{
						ID: "1",
					},
					Priv: true,
				},
				ignoredFields: []string{"Nested"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Nested.ID"},
					nil,
				),
			},
		},
		{
			name: "ignore nested struct second field filled",
			args: args{
				object: &struct {
					Nested
					Priv bool
				}{
					Nested: Nested{
						ID:  "",
						Seq: 134,
					},
					Priv: true,
				},
				ignoredFields: []string{"Nested"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Nested.Seq"},
					nil,
				),
			},
		},
		{
			name: "ignore nested struct last field filled",
			args: args{
				object: &struct {
					Nested
					Priv bool
				}{
					Nested: Nested{
						Seq2: 134,
					},
					Priv: true,
				},
				ignoredFields: []string{"Nested"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Nested.Seq2"},
					nil,
				),
			},
		},
		{
			name: "ignore not nested field",
			args: args{
				object: &struct {
					Number  int
					Company *struct {
						Name    string
						Founded time.Time
					}
					priv bool
				}{
					Number: 13,
					Company: &struct {
						Name    string
						Founded time.Time
					}{
						Name:    "caos AG",
						Founded: time.Date(2019, time.April, 1, 0, 0, 0, 0, time.Local),
					},
					priv: true,
				},
				ignoredFields: []string{"Number"},
			},
			res: res{
				mock: newMappedExpeter(
					[]string{"Number"},
					nil,
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertFieldsMapped(tt.res.mock, tt.args.object, tt.args.ignoredFields...)
			tt.res.mock.expectationsMet(t)
		})
	}
}

type Nested struct {
	ID   string
	Seq  uint64
	Seq2 uint64
}

func newMappedExpeter(expectedNotEmptyFields, expectedNotFilledFields []string) *mappedExpecter {
	sort.Strings(expectedNotEmptyFields)
	sort.Strings(expectedNotFilledFields)
	return &mappedExpecter{
		expectedNotEmptyFields:  strings.Join(expectedNotEmptyFields, ",\n\t"),
		expectedNotFilledFields: strings.Join(expectedNotFilledFields, ",\n\t"),
	}
}

type mappedExpecter struct {
	expectedNotEmptyFields  string
	expectedNotFilledFields string

	notEmptyFields  string
	notFilledFields string
}

func (e *mappedExpecter) Errorf(format string, args ...interface{}) {
	if format == "expected fields are not empty:\n\t%s\n" {
		e.notEmptyFields = args[0].(string)
	}
	if format == "unexpected empty fields:\n\t%s\n" {
		e.notFilledFields = args[0].(string)
	}
}

func (e *mappedExpecter) Helper() {}

func (e *mappedExpecter) expectationsMet(t *testing.T) {
	if e.notEmptyFields != e.expectedNotEmptyFields {
		t.Errorf("not empty fields not matched: \n expected:\n\t%s\n got:\n\t%s", e.expectedNotEmptyFields, e.notEmptyFields)
	}

	if e.notFilledFields != e.expectedNotFilledFields {
		t.Errorf("not filled fields not matched: \n expected:\n\t%s\n got:\n\t%s", e.expectedNotFilledFields, e.notFilledFields)
	}
}
