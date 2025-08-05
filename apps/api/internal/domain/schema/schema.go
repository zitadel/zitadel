package schema

import (
	_ "embed"
	"io"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed zitadel.schema.v1.json
	zitadelJSON string
)

const (
	MetaSchemaID = "urn:zitadel:schema:v1"
)

func NewSchema(role Role, r io.Reader) (*jsonschema.Schema, error) {
	c := jsonschema.NewCompiler()
	if err := c.AddResource(PermissionSchemaID, strings.NewReader(permissionJSON)); err != nil {
		return nil, err
	}
	if err := c.AddResource(MetaSchemaID, strings.NewReader(zitadelJSON)); err != nil {
		return nil, err
	}
	c.RegisterExtension(PermissionSchemaID, permissionSchema, permissionExtension{
		role,
	})
	if err := c.AddResource("schema.json", r); err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "COMMA-Frh42", "Errors.UserSchema.Invalid")
	}
	schema, err := c.Compile("schema.json")
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "COMMA-W21tg", "Errors.UserSchema.Invalid")
	}
	return schema, nil
}
