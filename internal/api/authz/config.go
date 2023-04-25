package authz

type Config struct {
	RolePermissionMappings []RoleMapping
}

type RoleMapping struct {
	Role        string
	Permissions []string
}

type MethodMapping map[string]Option

type Option struct {
	Permission string
	CheckParam string
	AllowSelf  bool
}

func getPermissionsFromRole(rolePermissionMappings []RoleMapping, role string) []string {
	for _, roleMap := range rolePermissionMappings {
		if roleMap.Role == role {
			return roleMap.Permissions
		}
	}
	return nil
}
