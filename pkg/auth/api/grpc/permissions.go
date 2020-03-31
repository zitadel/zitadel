package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMyZitadelPermissions(ctx context.Context, _ *empty.Empty) (*MyPermissions, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-or67G", "Not implemented")
	//ctxData := auth.GetCtxData(ctx)
	//
	//grants, err := s.processor.ResolveGrants(ctx, ctxData.UserID, ctxData.OrgID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//permissions := &MyPermissions{Permissions: []string{}}
	//
	//for _, grant := range grants {
	//	for _, role := range grant.Roles {
	//		roleName, ctxID := auth.SplitPermission(role)
	//		for _, mapping := range s.authConf.RolePermissionMappings {
	//			if mapping.Role == roleName {
	//				permissions.appendPermissions(ctxID, mapping.Permissions...)
	//			}
	//		}
	//	}
	//}
	//
	//return permissions, nil
}

func (p *MyPermissions) appendPermissions(ctxID string, permissions ...string) {
	for _, permission := range permissions {
		p.appendPermission(ctxID, permission)
	}
}

func (p *MyPermissions) appendPermission(ctxID, permission string) {
	if ctxID != "" {
		permission = permission + ":" + ctxID
	}
	for _, existingPermission := range p.Permissions {
		if existingPermission == permission {
			return
		}
	}
	p.Permissions = append(p.Permissions, permission)
}
