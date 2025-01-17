package metadata

import (
	"context"
	"strings"
)

type Key string
type ScopedKey string

const (
	externalIdProvisioningDomainPlaceholder = "{provisioningDomain}"

	KeyPrefix                 = "urn:zitadel:scim:"
	KeyProvisioningDomain Key = KeyPrefix + "provisioning_domain"

	KeyExternalId               Key = KeyPrefix + "externalId"
	keyScopedExternalIdTemplate     = KeyPrefix + externalIdProvisioningDomainPlaceholder + ":externalId"
	KeyMiddleName               Key = KeyPrefix + "name.middleName"
	KeyHonorificPrefix          Key = KeyPrefix + "name.honorificPrefix"
	KeyHonorificSuffix          Key = KeyPrefix + "name.honorificSuffix"
	KeyProfileUrl               Key = KeyPrefix + "profileURL"
	KeyTitle                    Key = KeyPrefix + "title"
	KeyLocale                   Key = KeyPrefix + "locale"
	KeyTimezone                 Key = KeyPrefix + "timezone"
	KeyIms                      Key = KeyPrefix + "ims"
	KeyPhotos                   Key = KeyPrefix + "photos"
	KeyAddresses                Key = KeyPrefix + "addresses"
	KeyEntitlements             Key = KeyPrefix + "entitlements"
	KeyRoles                    Key = KeyPrefix + "roles"
)

var ScimUserRelevantMetadataKeys = []Key{
	KeyExternalId,
	KeyMiddleName,
	KeyHonorificPrefix,
	KeyHonorificSuffix,
	KeyProfileUrl,
	KeyTitle,
	KeyLocale,
	KeyTimezone,
	KeyIms,
	KeyPhotos,
	KeyAddresses,
	KeyEntitlements,
	KeyRoles,
}

func ScopeExternalIdKey(provisioningDomain string) ScopedKey {
	return ScopedKey(strings.Replace(keyScopedExternalIdTemplate, externalIdProvisioningDomainPlaceholder, provisioningDomain, 1))
}

func ScopeKey(ctx context.Context, key Key) ScopedKey {
	// only the externalID is scoped
	if key == KeyExternalId {
		return GetScimContextData(ctx).ExternalIDScopedMetadataKey
	}

	return ScopedKey(key)
}
