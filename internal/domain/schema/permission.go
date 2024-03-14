package schema

import (
	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed permission.schema.v1.json
	permissionJSON string

	permissionSchema = jsonschema.MustCompileString(PermissionSchemaID, permissionJSON)
)

const (
	PermissionSchemaID = "urn:zitadel:schema:permission-schema:v1"
	PermissionProperty = "urn:zitadel:schema:permission"
)

type role int32

const (
	roleUnspecified role = iota
	roleSelf
	roleOwner
)

type permissionExtension struct {
	role role
}

// Compile implements the [jsonschema.ExtCompiler] interface.
// It parses the permission schema extension / annotation of the passed field.
func (c permissionExtension) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (_ jsonschema.ExtSchema, err error) {
	perm, ok := m[PermissionProperty]
	if !ok {
		return nil, nil
	}
	p, ok := perm.(map[string]interface{})
	if !ok {
		return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-WR5gs", "invalid permission")
	}
	perms := new(permissions)
	for key, value := range p {
		switch key {
		case "self":
			perms.self, err = mapPermission(value)
			if err != nil {
				return
			}
		case "owner":
			perms.owner, err = mapPermission(value)
			if err != nil {
				return
			}
		default:
			return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-GFjio", "invalid permission role")
		}
	}
	return permissionExtensionConfig{c.role, perms}, nil
}

type permissionExtensionConfig struct {
	role        role
	permissions *permissions
}

// Validate implements the [jsonschema.ExtSchema] interface.
// It validates the fields of the json instance according to the permission schema.
func (s permissionExtensionConfig) Validate(ctx jsonschema.ValidationContext, v interface{}) error {
	switch s.role {
	case roleSelf:
		if s.permissions.self == nil || !s.permissions.self.write {
			return ctx.Error("permission", "missing required permission")
		}
		return nil
	case roleOwner:
		if s.permissions.owner == nil || !s.permissions.owner.write {
			return ctx.Error("permission", "missing required permission")
		}
		return nil
	case roleUnspecified:
		fallthrough
	default:
		return ctx.Error("permission", "missing required permission")
	}
}

func mapPermission(value any) (*permission, error) {
	p := new(permission)
	switch v := value.(type) {
	case string:
		for _, s := range v {
			switch s {
			case 'r':
				p.read = true
			case 'w':
				p.write = true
			default:
				return nil, zerrors.ThrowInvalidArgumentf(nil, "SCHEMA-EZ5zjh", "invalid permission pattern: `%s` in (%s)", string(s), v)
			}
		}
		return p, nil
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "SCHEMA-E5h31", "invalid permission type %T (%v)", v, v)
	}
}

type permissions struct {
	self  *permission
	owner *permission
}

type permission struct {
	read  bool
	write bool
}
