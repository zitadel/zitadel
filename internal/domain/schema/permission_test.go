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

//go:embed zitadel.schema.json
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
			"invalid permission string, compilation err",
			args{
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"admin": "read"
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
										"admin": true
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
			"invalid permission, validation err",
			args{
				role: roleSelf,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"admin": "rw",
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
			"valid permission, ok",
			args{
				role: roleAdmin,
				schema: `{
							"type": "object",
							"properties": {
								"name": {
									"type": "string",
									"urn:zitadel:schema:permission": {
										"admin": "rw",
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
			"no permission, ok",
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
			err := c.AddResource("urn:zitadel:schema:permission-schema", strings.NewReader(permissionJSON))
			require.NoError(t, err)
			err = c.AddResource("urn:zitadel:schema", strings.NewReader(zitadelJSON))
			require.NoError(t, err)
			c.RegisterExtension("urn:zitadel:schema:permission-schema", permissionMeta, permissionExtension{
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
