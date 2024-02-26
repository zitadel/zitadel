package schema

import (
	_ "embed"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/zitadel/zitadel/internal/zerrors"
)

//go:embed permission.schema.json
var permissionJSON string

var permissionMeta = jsonschema.MustCompileString("urn:zitadel:schema:permission-schema", permissionJSON)

type role int32

const (
	roleUnspecified role = iota
	roleSelf
	roleAdmin
)

type permissionExtension struct {
	role role
}

func (c permissionExtension) Compile(ctx jsonschema.CompilerContext, m map[string]interface{}) (_ jsonschema.ExtSchema, err error) {
	perm, ok := m["urn:zitadel:schema:permission"]
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
		case "admin":
			perms.admin, err = mapPermission(value)
			if err != nil {
				return
			}
		default:
			return nil, zerrors.ThrowInvalidArgument(nil, "SCHEMA-GFjio", "invalid permission")
		}
	}
	return permissionExtensionConfig{c.role, perms}, nil
}

type permissionExtensionConfig struct {
	role        role
	permissions *permissions
}

func (s permissionExtensionConfig) Validate(ctx jsonschema.ValidationContext, v interface{}) error {
	switch s.role {
	case roleSelf:
		if s.permissions.self == nil || !s.permissions.self.write {
			return ctx.Error("permission", "missing required permission")
		}
		return nil
	case roleAdmin:
		if s.permissions.admin == nil || !s.permissions.admin.write {
			return ctx.Error("permission", "missing required permission")
		}
		return nil
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
	admin *permission
}

type permission struct {
	read  bool
	write bool
}
