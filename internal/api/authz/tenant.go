package authz

import "context"

type Tenant struct {
	ID string
}

func WithTenant(parent context.Context, id string) context.Context {
	return context.WithValue(parent, tenantKey, Tenant{ID: id})
}

func GetTenant(ctx context.Context) string {
	tenant, _ := ctx.Value(tenantKey).(Tenant)
	return tenant.ID
}
