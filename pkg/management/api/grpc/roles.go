package grpc

import "strings"

func (s *Server) getOrgMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range s.authZ.RolePermissionMappings {
		if strings.HasPrefix(roleMap.Role, "ORG") {
			roles = append(roles, roleMap.Role)
		}
	}
	return roles
}

func (s *Server) getProjectMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range s.authZ.RolePermissionMappings {
		if strings.HasPrefix(roleMap.Role, "PROJECT") && !strings.HasPrefix(roleMap.Role, "PROJECT_GRANT") {
			roles = append(roles, roleMap.Role)
		}
	}
	return roles
}

func (s *Server) getProjectGrantMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range s.authZ.RolePermissionMappings {
		if strings.HasPrefix(roleMap.Role, "PROJECT_GRANT") {
			roles = append(roles, roleMap.Role)
		}
	}
	return roles
}
