package authz

import "context"

type Region struct {
	Region string
}

func WithRegion(parent context.Context, region string) context.Context {
	return context.WithValue(parent, regionKey, Region{Region: region})
}

func GetRegion(ctx context.Context) string {
	region, _ := ctx.Value(tenantKey).(Region)
	if region.Region == "" {
		return "global"
	}
	return region.Region
}
