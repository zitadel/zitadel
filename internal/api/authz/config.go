package authz

import "time"

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
	Feature    string
	MaxAge     time.Duration
}

func (a *Config) getPermissionsFromRole(role string) []string {
	for _, roleMap := range a.RolePermissionMappings {
		if roleMap.Role == role {
			return roleMap.Permissions
		}
	}
	return nil
}
