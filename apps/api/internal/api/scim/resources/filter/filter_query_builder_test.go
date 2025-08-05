package filter

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/test"
)

var fieldPathColumnMapping = FieldPathMapping{
	// a timestamp field
	"meta.lastmodified": {
		Column:    query.UserChangeDateCol,
		FieldType: FieldTypeTimestamp,
	},
	// a case-insensitive string field
	"username": {
		Column:          query.UserUsernameCol,
		FieldType:       FieldTypeString,
		CaseInsensitive: true,
	},
	// a nested string field
	"name.familyname": {
		Column:    query.HumanLastNameCol,
		FieldType: FieldTypeString,
	},
	// a field which is a list in scim
	"emails": {
		Column:    query.HumanEmailCol,
		FieldType: FieldTypeString,
	},
	// the default value field
	"emails.value": {
		Column:    query.HumanEmailCol,
		FieldType: FieldTypeString,
	},
	// pseudo field to test number queries
	"age": {
		Column:    query.HumanGenderCol,
		FieldType: FieldTypeNumber,
	},
	// pseudo field to test boolean queries
	"locked": {
		Column:    query.HumanPasswordChangeRequiredCol,
		FieldType: FieldTypeBoolean,
	},
	// mapped field
	"active": {
		Column:    query.UserStateCol,
		FieldType: FieldTypeCustom,
		BuildMappedQuery: func(ctx context.Context, compareValue *CompValue, op *CompareOp) (query.SearchQuery, error) {
			// very simple mock implementation
			return query.NewTextQuery(query.UserUsernameCol, "fooBar", query.TextContains)
		},
	},
}

func TestFilter_BuildQuery(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		want    query.SearchQuery
		wantErr bool
	}{
		{
			name:    "unknown attribute",
			filter:  `foobar eq "bjensen"`,
			wantErr: true,
		},
		{
			name:   "simple binary operator",
			filter: `userName eq "bjensen"`,
			want:   test.Must(query.NewTextQuery(query.UserUsernameCol, "bjensen", query.TextEqualsIgnoreCase)),
		},
		{
			name:   "binary operator equals null",
			filter: `userName eq null`,
			want:   test.Must(query.NewIsNullQuery(query.UserUsernameCol)),
		},
		{
			name:   "binary operator not equals null",
			filter: `userName ne null`,
			want:   test.Must(query.NewNotNullQuery(query.UserUsernameCol)),
		},
		{
			name:    "binary number operator on string field",
			filter:  `userName gt 10`,
			wantErr: true,
		},
		{
			name:   "binary number operator greater",
			filter: `age gt 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberGreater)),
		},
		{
			name:   "binary number operator greater equal",
			filter: `age ge 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberGreaterOrEqual)),
		},
		{
			name:   "binary number operator less",
			filter: `age lt 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberLess)),
		},
		{
			name:   "binary number operator less float",
			filter: `age lt 10.5`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10.5, query.NumberLess)),
		},
		{
			name:    "binary number unsupported operator",
			filter:  `age co 10.5`,
			wantErr: true,
		},
		{
			name:    "binary number unsupported comparison value",
			filter:  `age gt "foo"`,
			wantErr: true,
		},
		{
			name:   "binary number operator less equal",
			filter: `age le 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberLessOrEqual)),
		},
		{
			name:   "binary number operator equals",
			filter: `age eq 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberEquals)),
		},
		{
			name:   "binary number operator not equals",
			filter: `age ne 10`,
			want:   test.Must(query.NewNumberQuery(query.HumanGenderCol, 10, query.NumberNotEquals)),
		},
		{
			name:    "binary bool operator equals string",
			filter:  `locked eq "foo"`,
			wantErr: true,
		},
		{
			name:    "binary bool operator startswith bool",
			filter:  `locked sw true`,
			wantErr: true,
		},
		{
			name:   "binary bool operator equals",
			filter: `locked eq true`,
			want:   test.Must(query.NewBoolQuery(query.HumanPasswordChangeRequiredCol, true)),
		},
		{
			name:   "binary bool operator not equals",
			filter: `locked ne true`,
			want:   test.Must(query.NewBoolQuery(query.HumanPasswordChangeRequiredCol, false)),
		},
		{
			name:   "binary bool operator not equals false",
			filter: `locked ne false`,
			want:   test.Must(query.NewBoolQuery(query.HumanPasswordChangeRequiredCol, true)),
		},
		{
			name:    "binary string invalid operator",
			filter:  `username gt "test"`,
			wantErr: true,
		},
		{
			name:   "nested attribute binary operator",
			filter: `name.familyName co "O'Malley"`,
			want:   test.Must(query.NewTextQuery(query.HumanLastNameCol, "O'Malley", query.TextContains)),
		},
		{
			name:   "urn prefixed binary operator",
			filter: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "J"`,
			want:   test.Must(query.NewTextQuery(query.UserUsernameCol, "J", query.TextStartsWithIgnoreCase)),
		},
		{
			name:   "urn prefixed nested binary operator",
			filter: `urn:ietf:params:scim:schemas:core:2.0:User:emails[value sw "hans.peter@"]`,
			want:   test.Must(query.NewTextQuery(query.HumanEmailCol, "hans.peter@", query.TextStartsWith)),
		},
		{
			name:    "invalid urn prefixed nested binary operator",
			filter:  `urn:ietf:params:scim:schemas:core:2.0:UserFoo:emails[value sw "hans.peter@"]`,
			wantErr: true,
		},
		{
			name:   "unary operator",
			filter: `name.familyName pr`,
			want:   test.Must(query.NewNotNullQuery(query.HumanLastNameCol)),
		},
		{
			name:   "and logical expression",
			filter: `name.familyName pr and userName eq "bjensen"`,
			want:   test.Must(query.NewAndQuery(test.Must(query.NewNotNullQuery(query.HumanLastNameCol)), test.Must(query.NewTextQuery(query.UserUsernameCol, "bjensen", query.TextEqualsIgnoreCase)))),
		},
		{
			name:   "timestamp condition equal",
			filter: `meta.lastModified eq "2011-05-13T04:42:34Z"`,
			want:   test.Must(query.NewTimestampQuery(query.UserChangeDateCol, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC), query.TimestampEquals)),
		},
		{
			name:   "timestamp condition greater equals",
			filter: `meta.lastModified ge "2011-05-13T04:42:34Z"`,
			want:   test.Must(query.NewTimestampQuery(query.UserChangeDateCol, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC), query.TimestampGreaterOrEquals)),
		},
		{
			name:   "timestamp condition greater",
			filter: `meta.lastModified gt "2011-05-13T04:42:34Z"`,
			want:   test.Must(query.NewTimestampQuery(query.UserChangeDateCol, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC), query.TimestampGreater)),
		},
		{
			name:   "timestamp condition less equals",
			filter: `meta.lastModified le "2011-05-13T04:42:34Z"`,
			want:   test.Must(query.NewTimestampQuery(query.UserChangeDateCol, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC), query.TimestampLessOrEquals)),
		},
		{
			name:   "timestamp condition less",
			filter: `meta.lastModified lt "2011-05-13T04:42:34Z"`,
			want:   test.Must(query.NewTimestampQuery(query.UserChangeDateCol, time.Date(2011, time.May, 13, 4, 42, 34, 0, time.UTC), query.TimestampLess)),
		},
		{
			name:    "timestamp condition invalid operator",
			filter:  `meta.lastModified ew "2011-05-13T04:42:34Z"`,
			wantErr: true,
		},
		{
			name:    "timestamp condition invalid format",
			filter:  `meta.lastModified ge "2011-05-13T0:34Z"`,
			wantErr: true,
		},
		{
			name:    "timestamp condition invalid comparison value",
			filter:  `meta.lastModified ge 15`,
			wantErr: true,
		},
		{
			name:   "nested and / or without grouping",
			filter: `userName eq "rudolpho" and emails co "example.com" or emails.value co "example2.org"`,
			want: test.Must(query.NewOrQuery(
				test.Must(query.NewAndQuery(
					test.Must(query.NewTextQuery(query.UserUsernameCol, "rudolpho", query.TextEqualsIgnoreCase)),
					test.Must(query.NewTextQuery(query.HumanEmailCol, "example.com", query.TextContains))),
				),
				test.Must(query.NewTextQuery(query.HumanEmailCol, "example2.org", query.TextContains)))),
		},
		{
			name:   "nested and / or with grouping",
			filter: `userName ne "rudolpho" and (emails co "example.com" or emails.value co "example.org")`,
			want: test.Must(query.NewAndQuery(
				test.Must(query.NewTextQuery(query.UserUsernameCol, "rudolpho", query.TextNotEqualsIgnoreCase)),
				test.Must(query.NewOrQuery(
					test.Must(query.NewTextQuery(query.HumanEmailCol, "example.com", query.TextContains)),
					test.Must(query.NewTextQuery(query.HumanEmailCol, "example.org", query.TextContains)),
				)),
			)),
		},
		{
			name:   "nested value path path",
			filter: `userName eq "Hans" and emails[value ew "@example.org" or value ew "@example.com"]`,
			want: test.Must(query.NewAndQuery(
				test.Must(query.NewTextQuery(query.UserUsernameCol, "Hans", query.TextEqualsIgnoreCase)),
				test.Must(query.NewOrQuery(
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.org", query.TextEndsWith)),
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.com", query.TextEndsWith)),
				)),
			)),
		},
		{
			name:   "or value path filter",
			filter: `emails[value ew "@example.org" and value co "@example.com"] or emails[value sw "hans" or value sw "peter"]`,
			want: test.Must(query.NewOrQuery(
				test.Must(query.NewAndQuery(
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.org", query.TextEndsWith)),
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.com", query.TextContains)),
				)),
				test.Must(query.NewTextQuery(query.HumanEmailCol, "hans", query.TextStartsWith)),
				test.Must(query.NewTextQuery(query.HumanEmailCol, "peter", query.TextStartsWith)),
			)),
		},
		{
			name:   "and value path filter",
			filter: `emails[value ew "@example.com"] and name.familyname co "hans" and username co "peter"`,
			want: test.Must(query.NewAndQuery(
				test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.com", query.TextEndsWith)),
				test.Must(query.NewTextQuery(query.HumanLastNameCol, "hans", query.TextContains)),
				test.Must(query.NewTextQuery(query.UserUsernameCol, "peter", query.TextContainsIgnoreCase)),
			)),
		},
		{
			name:   "negation",
			filter: `not(username eq "foo")`,
			want:   test.Must(query.NewNotQuery(test.Must(query.NewTextQuery(query.UserUsernameCol, "foo", query.TextEqualsIgnoreCase)))),
		},
		{
			name:   "negation with complex filter",
			filter: `not(emails[value ew "@example.com"])`,
			want:   test.Must(query.NewNotQuery(test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.com", query.TextEndsWith)))),
		},
		{
			name:   "nested negation",
			filter: `emails[not(value ew "@example.org" or value ew "@example.com")]`,
			want: test.Must(query.NewNotQuery(
				test.Must(query.NewOrQuery(
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.org", query.TextEndsWith)),
					test.Must(query.NewTextQuery(query.HumanEmailCol, "@example.com", query.TextEndsWith)),
				)),
			)),
		},
		{
			name:   "mapped field",
			filter: `active eq true`,
			want:   test.Must(query.NewTextQuery(query.UserUsernameCol, "fooBar", query.TextContains)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := ParseFilter(tt.filter)
			require.NoError(t, err)

			got, err := f.BuildQuery(context.Background(), schemas.IdUser, fieldPathColumnMapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildQuery() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_queryBuilder_reduceAttrPaths(t *testing.T) {
	tests := []struct {
		name          string
		schema        string
		attrPaths     []*AttrPath
		wantFieldPath string
		wantErr       bool
	}{
		{
			name:      "empty",
			attrPaths: []*AttrPath{},
			wantErr:   true,
		},
		{
			name: "simple",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
				},
			},
			wantFieldPath: "foo",
		},
		{
			name: "multiple simple",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
				},
				{
					AttrName: "bar",
				},
			},
			wantFieldPath: "foo.bar",
		},
		{
			name: "with sub attr",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
					SubAttr:  gu.Ptr("bar"),
				},
			},
			wantFieldPath: "foo.bar",
		},
		{
			name: "multiple with sub attr",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
					SubAttr:  gu.Ptr("bar"),
				},
				{
					AttrName: "baz",
					SubAttr:  gu.Ptr("woo"),
				},
			},
			wantFieldPath: "foo.bar.baz.woo",
		},
		{
			name:   "with urn and sub attr",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:bar:"),
					AttrName:           "foo",
					SubAttr:            gu.Ptr("bar"),
				},
			},
			wantFieldPath: "foo.bar",
		},
		{
			name:   "multiple with urn and sub attr",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:bar:"),
					AttrName:           "foo",
					SubAttr:            gu.Ptr("bar"),
				},
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:bar:"),
					AttrName:           "foo2",
					SubAttr:            gu.Ptr("bar2"),
				},
			},
			wantFieldPath: "foo.bar.foo2.bar2",
		},
		{
			name:   "secondary with urn and sub attr",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
					SubAttr:  gu.Ptr("bar"),
				},
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:bar:"),
					AttrName:           "foo2",
					SubAttr:            gu.Ptr("bar2"),
				},
			},
			wantFieldPath: "foo.bar.foo2.bar2",
		},
		{
			name:   "urn mismatch",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:baz"),
					AttrName:           "foo",
				},
			},
			wantErr: true,
		},
		{
			name:   "nested urn mismatch",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:bar:"),
					AttrName:           "foo",
				},
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:baz"),
					AttrName:           "foo2",
				},
			},
			wantErr: true,
		},
		{
			name:   "secondary urn mismatch",
			schema: "urn:foo:bar",
			attrPaths: []*AttrPath{
				{
					AttrName: "foo",
				},
				{
					UrnAttributePrefix: gu.Ptr("urn:foo:baz"),
					AttrName:           "foo2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &queryBuilder{
				schema: schemas.ScimSchemaType(tt.schema),
			}
			gotFieldPath, err := b.reduceAttrPaths(tt.attrPaths)
			if (err != nil) != tt.wantErr {
				t.Errorf("reduceAttrPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFieldPath != tt.wantFieldPath {
				t.Errorf("reduceAttrPaths() gotFieldPath = %v, want %v", gotFieldPath, tt.wantFieldPath)
			}
		})
	}
}
