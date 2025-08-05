package filter

import (
	"reflect"
	"strings"
	"testing"

	"github.com/muhlemmer/gu"
)

var longString = ""

func init() {
	var sb strings.Builder
	for i := 0; i < maxInputLength+1; i++ {
		sb.WriteRune('x')
	}

	longString = sb.String()
}

func TestParseFilter(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		want    *Filter
		wantErr bool
	}{
		{
			name: "empty",
		},
		{
			name:    "too long",
			filter:  longString,
			wantErr: true,
		},
		{
			name:    "invalid syntax",
			filter:  "fooBar[['baz']]",
			wantErr: true,
		},
		{
			name:    "unknown binary operator",
			filter:  `userName fu "bjensen"`,
			wantErr: true,
		},
		{
			name:    "unknown unary operator",
			filter:  `userName ok`,
			wantErr: true,
		},

		// test cases from https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.2.2
		{
			name:   "negation",
			filter: `not(username pr)`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								Negation: &Segment{
									OrExp: OrLogExp{
										Left: AndLogExp{
											Left: ValueAtom{
												AttrExp: &AttrExp{
													UnaryCondition: &UnaryCondition{
														Left: AttrPath{
															AttrName: "username",
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
							},
						},
					},
				},
			},
		},
		{
			name:   "number",
			filter: `age gt 10`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "age",
										},
										Operator: CompareOp{
											GreaterThan: true,
										},
										Right: CompValue{
											Int: gu.Ptr(10),
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
			name:   "float",
			filter: `age gt 10.5`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "age",
										},
										Operator: CompareOp{
											GreaterThan: true,
										},
										Right: CompValue{
											Float: gu.Ptr(10.5),
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
			name:   "null",
			filter: `age eq null`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "age",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											Null: true,
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
			name:   "simple binary operator",
			filter: `userName eq "bjensen"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userName",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("bjensen"),
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
			name:   "uppercase binary operator",
			filter: `userName EQ "bjensen"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userName",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("bjensen"),
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
			name:   "case-insensitive literals and operators",
			filter: `active Eq TRue`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "active",
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
		{
			name:   "extra whitespace",
			filter: `userName      eq      "bjensen"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userName",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("bjensen"),
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
			name:   "nested attribute binary operator",
			filter: `name.familyName co "O'Malley"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "name",
											SubAttr:  gu.Ptr("familyName"),
										},
										Operator: CompareOp{
											Contains: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("O'Malley"),
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
			name:   "urn prefixed",
			filter: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "J"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											UrnAttributePrefix: gu.Ptr("urn:ietf:params:scim:schemas:core:2.0:User:"),
											AttrName:           "userName",
										},
										Operator: CompareOp{
											StartsWith: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("J"),
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
			name:   "urn prefixed nested",
			filter: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "J" and urn:ietf:params:scim:schemas:core:2.0:User:emails.value ew "@example.com"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											UrnAttributePrefix: gu.Ptr("urn:ietf:params:scim:schemas:core:2.0:User:"),
											AttrName:           "userName",
										},
										Operator: CompareOp{
											StartsWith: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("J"),
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												UrnAttributePrefix: gu.Ptr("urn:ietf:params:scim:schemas:core:2.0:User:"),
												AttrName:           "emails",
												SubAttr:            gu.Ptr("value"),
											},
											Operator: CompareOp{
												EndsWith: true,
											},
											Right: CompValue{
												StringValue: gu.Ptr("@example.com"),
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
			name:   "unary operator",
			filter: `title pr`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									UnaryCondition: &UnaryCondition{
										Left: AttrPath{
											AttrName: "title",
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
			},
		},
		{
			name:   "binary nested date operator",
			filter: `meta.lastModified gt "2011-05-13T04:42:34Z"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "meta",
											SubAttr:  gu.Ptr("lastModified"),
										},
										Operator: CompareOp{
											GreaterThan: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("2011-05-13T04:42:34Z"),
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
			name:   "and logical expression",
			filter: `title pr and userType eq "Employee"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									UnaryCondition: &UnaryCondition{
										Left: AttrPath{
											AttrName: "title",
										},
										Operator: UnaryConditionOperator{
											Present: true,
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												AttrName: "userType",
											},
											Operator: CompareOp{
												Equal: true,
											},
											Right: CompValue{
												StringValue: gu.Ptr("Employee"),
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
			name:   "nested and / or with grouping",
			filter: `userType eq "Employee" and (emails co "example.com" or emails.value co "example.org")`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userType",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("Employee"),
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									SubFilter: &Segment{
										OrExp: OrLogExp{
											Left: AndLogExp{
												Left: ValueAtom{
													AttrExp: &AttrExp{
														BinaryCondition: &BinaryCondition{
															Left: AttrPath{
																AttrName: "emails",
															},
															Operator: CompareOp{
																Contains: true,
															},
															Right: CompValue{
																StringValue: gu.Ptr("example.com"),
															},
														},
													},
												},
											},
											Right: &OrLogExp{
												Left: AndLogExp{
													Left: ValueAtom{
														AttrExp: &AttrExp{
															BinaryCondition: &BinaryCondition{
																Left: AttrPath{
																	AttrName: "emails",
																	SubAttr:  gu.Ptr("value"),
																},
																Operator: CompareOp{
																	Contains: true,
																},
																Right: CompValue{
																	StringValue: gu.Ptr("example.org"),
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
						},
					},
				},
			},
		},
		{
			name:   "nested and / or without grouping",
			filter: `userType eq "Employee" and emails co "example.com" or emails.value co "example2.org"`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userType",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("Employee"),
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												AttrName: "emails",
											},
											Operator: CompareOp{
												Contains: true,
											},
											Right: CompValue{
												StringValue: gu.Ptr("example.com"),
											},
										},
									},
								},
							},
						},
						Right: &OrLogExp{
							Left: AndLogExp{
								Left: ValueAtom{
									AttrExp: &AttrExp{
										BinaryCondition: &BinaryCondition{
											Left: AttrPath{
												AttrName: "emails",
												SubAttr:  gu.Ptr("value"),
											},
											Operator: CompareOp{
												Contains: true,
											},
											Right: CompValue{
												StringValue: gu.Ptr("example2.org"),
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
			name:   "nested and / or with negated grouping",
			filter: `userType ne "Employee" and not (emails co "example.com" or emails.value co "example.org")`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userType",
										},
										Operator: CompareOp{
											NotEqual: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("Employee"),
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									Negation: &Segment{
										OrExp: OrLogExp{
											Left: AndLogExp{
												Left: ValueAtom{
													AttrExp: &AttrExp{
														BinaryCondition: &BinaryCondition{
															Left: AttrPath{
																AttrName: "emails",
															},
															Operator: CompareOp{
																Contains: true,
															},
															Right: CompValue{
																StringValue: gu.Ptr("example.com"),
															},
														},
													},
												},
											},
											Right: &OrLogExp{
												Left: AndLogExp{
													Left: ValueAtom{
														AttrExp: &AttrExp{
															BinaryCondition: &BinaryCondition{
																Left: AttrPath{
																	AttrName: "emails",
																	SubAttr:  gu.Ptr("value"),
																},
																Operator: CompareOp{
																	Contains: true,
																},
																Right: CompValue{
																	StringValue: gu.Ptr("example.org"),
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
						},
					},
				},
			},
		},
		{
			name:   "nested value path path",
			filter: `userType eq "Employee" and emails[type eq "work" and value co "@example.com"]`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								AttrExp: &AttrExp{
									BinaryCondition: &BinaryCondition{
										Left: AttrPath{
											AttrName: "userType",
										},
										Operator: CompareOp{
											Equal: true,
										},
										Right: CompValue{
											StringValue: gu.Ptr("Employee"),
										},
									},
								},
							},
							Right: &AndLogExp{
								Left: ValueAtom{
									ValuePath: &ValuePath{
										AttrPath: AttrPath{
											AttrName: "emails",
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
												Right: &AndLogExp{
													Left: ValueAtom{
														AttrExp: &AttrExp{
															BinaryCondition: &BinaryCondition{
																Left: AttrPath{
																	AttrName: "value",
																},
																Operator: CompareOp{
																	Contains: true,
																},
																Right: CompValue{
																	StringValue: gu.Ptr("@example.com"),
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
						},
					},
				},
			},
		},
		{
			name:   "complex value path filter",
			filter: `emails[type eq "work" and value co "@example.com"] or ims[type eq "xmpp" and value co "@foo.com"]`,
			want: &Filter{
				Root: Segment{
					OrExp: OrLogExp{
						Left: AndLogExp{
							Left: ValueAtom{
								ValuePath: &ValuePath{
									AttrPath: AttrPath{
										AttrName: "emails",
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
											Right: &AndLogExp{
												Left: ValueAtom{
													AttrExp: &AttrExp{
														BinaryCondition: &BinaryCondition{
															Left: AttrPath{
																AttrName: "value",
															},
															Operator: CompareOp{
																Contains: true,
															},
															Right: CompValue{
																StringValue: gu.Ptr("@example.com"),
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
						Right: &OrLogExp{
							Left: AndLogExp{
								Left: ValueAtom{
									ValuePath: &ValuePath{
										AttrPath: AttrPath{
											AttrName: "ims",
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
																StringValue: gu.Ptr("xmpp"),
															},
														},
													},
												},
												Right: &AndLogExp{
													Left: ValueAtom{
														AttrExp: &AttrExp{
															BinaryCondition: &BinaryCondition{
																Left: AttrPath{
																	AttrName: "value",
																},
																Operator: CompareOp{
																	Contains: true,
																},
																Right: CompValue{
																	StringValue: gu.Ptr("@foo.com"),
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
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFilter(tt.filter)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFilter() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFilter() got = %s, want %s", got, tt.want)
			}
		})
	}
}
