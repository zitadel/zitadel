package schema

import (
	_ "embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed zitadel.schema.v1.json
var zitadelJSON string

func TestPermissionExtension(t *testing.T) {
	type args struct {
		role     role
		schema   string
		instance string
	}
	type want struct {
		compilationErr error
		validationErr  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"invalid permission, compilation err",
			args{
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": "read"
								}
							}
						}`,
			},
			want{
				compilationErr: zerrors.ThrowInvalidArgument(nil, "SCHEMA-WR5gs", "invalid permission"),
			},
		},
		{
			"invalid permission string, compilation err",
			args{
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"self": "read"
									}
								}
							}
						}`,
			},
			want{
				compilationErr: zerrors.ThrowInvalidArgument(nil, "SCHEMA-EZ5zjh", "invalid permission pattern: `e` in (read)"),
			},
		},
		{
			"invalid permission type, compilation err",
			args{
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": true
									}
								}
							}
						}`,
			},
			want{
				compilationErr: zerrors.ThrowInvalidArgument(nil, "SCHEMA-E5h31", "invalid permission type bool (true)"),
			},
		},
		{
			"invalid role, compilation err",
			args{
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"IAM_OWNER": "rw"
									}
								}
							}
						}`,
			},
			want{
				compilationErr: zerrors.ThrowInvalidArgument(nil, "SCHEMA-GFjio", "invalid permission role"),
			},
		},
		{
			"invalid permission self, validation err",
			args{
				role: roleSelf,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": "rw",
										"self": "r"
									}
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: true,
			},
		},
		{
			"invalid permission owner, validation err",
			args{
				role: roleOwner,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": "r",
										"self": "r"
									}
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: true,
			},
		},
		{
			"valid permission self, ok",
			args{
				role: roleSelf,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": "r",
										"self": "rw"
									}
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: false,
			},
		},
		{
			"valid permission owner, ok",
			args{
				role: roleOwner,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": "rw",
										"self": "r"
									}
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: false,
			},
		},
		{
			"no role, validation err",
			args{
				role: roleUnspecified,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"owner": "rw",
										"self": "rw"
									}
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: true,
			},
		},
		{
			"no permission required, ok",
			args{
				role: roleSelf,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string"
								}
							}
						}`,
				instance: `{ "name": "test"}`,
			},
			want{
				validationErr: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := jsonschema.NewCompiler()
			err := c.AddResource("urn:zitadel:schema:permission-schema:v1", strings.NewReader(permissionJSON))
			require.NoError(t, err)
			err = c.AddResource("urn:zitadel:schema:v1", strings.NewReader(zitadelJSON))
			require.NoError(t, err)
			c.RegisterExtension("urn:zitadel:schema:permission-schema:v1", permissionSchema, permissionExtension{
				tt.args.role,
			})
			err = c.AddResource("schema.json", strings.NewReader(tt.args.schema))
			require.NoError(t, err)
			sch, err := c.Compile("schema.json")
			require.ErrorIs(t, err, tt.want.compilationErr)
			if tt.want.compilationErr != nil {
				return
			}

			var v interface{}
			err = json.Unmarshal([]byte(tt.args.instance), &v)
			require.NoError(t, err)

			err = sch.Validate(v)
			if tt.want.validationErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
