package filter

import (
	"reflect"
	"testing"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
)

type attributeResolverTestType struct {
	StringValue string
	Nested      struct {
		IntValue int
	}
}

func TestAttributeResolver_resolveAttrPath(t *testing.T) {
	tests := []struct {
		name         string
		schema       schemas.ScimSchemaType
		item         interface{}
		attrPath     *AttrPath
		wantSegments []string
		wantValue    interface{}
		wantErr      bool
	}{
		{
			name:   "simple",
			schema: "fooBar",
			item:   &attributeResolverTestType{StringValue: "value"},
			attrPath: &AttrPath{
				AttrName: "StringValue",
			},
			wantSegments: []string{"stringvalue"},
			wantValue:    "value",
		},
		{
			name:   "simple case insensitive",
			schema: "fooBar",
			item:   &attributeResolverTestType{StringValue: "value"},
			attrPath: &AttrPath{
				AttrName: "stringvalue",
			},
			wantSegments: []string{"stringvalue"},
			wantValue:    "value",
		},
		{
			name:   "invalid attribute",
			schema: "fooBar",
			item:   &attributeResolverTestType{StringValue: "value"},
			attrPath: &AttrPath{
				AttrName: "StringValue2",
			},
			wantErr: true,
		},
		{
			name:   "with SubAttr",
			schema: "fooBar",
			item:   &attributeResolverTestType{Nested: struct{ IntValue int }{IntValue: 42}},
			attrPath: &AttrPath{
				AttrName: "Nested",
				SubAttr:  gu.Ptr("IntValue"),
			},
			wantSegments: []string{"nested", "intvalue"},
			wantValue:    42,
		},
		{
			name:   "with case-insensitive SubAttr",
			schema: "fooBar",
			item:   &attributeResolverTestType{Nested: struct{ IntValue int }{IntValue: 42}},
			attrPath: &AttrPath{
				AttrName: "NESTED",
				SubAttr:  gu.Ptr("intvalue"),
			},
			wantSegments: []string{"nested", "intvalue"},
			wantValue:    42,
		},
		{
			name:   "with invalid SubAttr",
			schema: "fooBar",
			item:   &attributeResolverTestType{Nested: struct{ IntValue int }{IntValue: 42}},
			attrPath: &AttrPath{
				AttrName: "NESTED",
				SubAttr:  gu.Ptr("xxxxx"),
			},
			wantErr: true,
		},
		{
			name:   "with SubAttr and urn",
			schema: "fooBar",
			item:   &attributeResolverTestType{Nested: struct{ IntValue int }{IntValue: 42}},
			attrPath: &AttrPath{
				UrnAttributePrefix: gu.Ptr("fooBar:"),
				AttrName:           "Nested",
				SubAttr:            gu.Ptr("IntValue"),
			},
			wantSegments: []string{"nested", "intvalue"},
			wantValue:    42,
		},
		{
			name:   "with SubAttr and invalid urn",
			schema: "fooBar",
			item:   &attributeResolverTestType{Nested: struct{ IntValue int }{IntValue: 42}},
			attrPath: &AttrPath{
				UrnAttributePrefix: gu.Ptr("fooBaz:"),
				AttrName:           "Nested",
				SubAttr:            gu.Ptr("IntValue"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newAttributeResolver(tt.schema)

			gotSegments, gotValue, err := r.resolveAttrPath(reflect.ValueOf(tt.item), tt.attrPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveAttrPath() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if !reflect.DeepEqual(gotSegments, tt.wantSegments) {
				t.Errorf("resolveAttrPath() gotSegments = %#v, want %#v", gotSegments, tt.wantSegments)
			}

			if !reflect.DeepEqual(gotValue.Interface(), tt.wantValue) {
				t.Errorf("resolveAttrPath() gotValue = %#v, want %#v", gotValue, tt.wantValue)
			}
		})
	}
}
