package metadata

import (
	"context"
)

type provisioningDomainKeyType struct{}

var provisioningDomainKey provisioningDomainKeyType

type ScimContextData struct {
	ProvisioningDomain          string
	ExternalIDScopedMetadataKey ScopedKey
}

func SetScimContextData(ctx context.Context, data ScimContextData) context.Context {
	return context.WithValue(ctx, provisioningDomainKey, data)
}

func GetScimContextData(ctx context.Context) ScimContextData {
	data, _ := ctx.Value(provisioningDomainKey).(ScimContextData)
	return data
}
