package filter

import (
	"reflect"
	"testing"

	"github.com/muhlemmer/gu"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Path
		wantErr bool
	}{
		{
			name: "empty",
			path: "",
		},
		{
			name:    "too long",
			path:    longString,
			wantErr: true,
		},
		{
			name:    "invalid syntax",
			path:    "fooBar[['baz']]",
			wantErr: true,
		},

		// test cases from https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2
		{
			name: "simple",
			path: "members",
			want: &Path{
				AttrPath: &AttrPath{
					AttrName: "members",
				},
			},
		},
		{
			name: "nested",
			path: "name.familyName",
			want: &Path{
				AttrPath: &AttrPath{
					AttrName: "name",
					SubAttr:  gu.Ptr("familyName"),
				},
			},
		},
		{
			name: "with filter",
			path: `addresses[type eq "work"]`,
			want: &Path{
				ValuePath: &ValuePathWithSubAttr{
					ValuePath: ValuePath{
						AttrPath: AttrPath{
							AttrName: "addresses",
						},
						ValFilter: OrLogExp{
							Left: AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												AttrName: "type",
											},
											Operator: CompareOp{
												Equal: true,
											},
											Right: CompValue{
												StringValue: gu.Ptr("work"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "with filter and submember",
			path: `members[value pr].displayName`,
			want: &Path{
				ValuePath: &ValuePathWithSubAttr{
					ValuePath: ValuePath{
						AttrPath: AttrPath{
							AttrName: "members",
						},
						ValFilter: OrLogExp{
							Left: AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										UnaryCondition: &UnaryCondition{
											Left: AttrPath{
												AttrName: "value",
											},
											Operator: UnaryConditionOperator{
												Present: true,
											},
										},
									},
								},
							},
						},
					},
					SubAttr: gu.Ptr("displayName"),
				},
			},
		},
		{
			name: "with binary filter",
			path: `entitlements[primary eq true]`,
			want: &Path{
				ValuePath: &ValuePathWithSubAttr{
					ValuePath: ValuePath{
						AttrPath: AttrPath{
							AttrName: "entitlements",
						},
						ValFilter: OrLogExp{
							Left: AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												AttrName: "primary",
											},
											Operator: CompareOp{
												Equal: true,
											},
											Right: CompValue{
												BooleanTrue: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
