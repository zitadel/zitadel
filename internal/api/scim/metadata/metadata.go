package metadata

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/query"
)

type Key string
type ScopedKey string

const (
	externalIdProvisioningDomainPlaceholder = "{provisioningDomain}"

	KeyPrefix                     = "urn:zitadel:scim:"
	KeyProvisioningDomain     Key = KeyPrefix + "provisioningDomain"
	KeyIgnorePasswordOnCreate Key = KeyPrefix + "ignorePasswordOnCreate"

	KeyExternalId               Key = KeyPrefix + "externalId"
	keyScopedExternalIdTemplate     = KeyPrefix + externalIdProvisioningDomainPlaceholder + ":externalId"
	KeyMiddleName               Key = KeyPrefix + "name.middleName"
	KeyHonorificPrefix          Key = KeyPrefix + "name.honorificPrefix"
	KeyHonorificSuffix          Key = KeyPrefix + "name.honorificSuffix"
	KeyProfileUrl               Key = KeyPrefix + "profileUrl"
	KeyTitle                    Key = KeyPrefix + "title"
	KeyLocale                   Key = KeyPrefix + "locale"
	KeyTimezone                 Key = KeyPrefix + "timezone"
	KeyIms                      Key = KeyPrefix + "ims"
	KeyPhotos                   Key = KeyPrefix + "photos"
	KeyAddresses                Key = KeyPrefix + "addresses"
	KeyEntitlements             Key = KeyPrefix + "entitlements"
	KeyRoles                    Key = KeyPrefix + "roles"
	KeyEmails                   Key = KeyPrefix + "emails"
)

var (
	ScimUserRelevantMetadataKeys = []Key{
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
		KeyEmails,
	}

	AttributePathToMetadataKeys = map[string][]Key{
		"externalid":           {KeyExternalId},
		"name":                 {KeyMiddleName, KeyHonorificPrefix, KeyHonorificSuffix},
		"name.middlename":      {KeyMiddleName},
		"name.honorificprefix": {KeyHonorificPrefix},
		"name.honorificsuffix": {KeyHonorificSuffix},
		"profileurl":           {KeyProfileUrl},
		"title":                {KeyTitle},
		"locale":               {KeyLocale},
		"timezone":             {KeyTimezone},
		"ims":                  {KeyIms},
		"photos":               {KeyPhotos},
		"addresses":            {KeyAddresses},
		"entitlements":         {KeyEntitlements},
		"roles":                {KeyRoles},
		"emails":               {KeyEmails},
	}
)

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

func MapToScopedKeyMap(md map[string][]byte) map[ScopedKey][]byte {
	result := make(map[ScopedKey][]byte, len(md))
	for k, v := range md {
		result[ScopedKey(k)] = v
	}

	return result
}

func MapListToScopedKeyMap(metadataList []*query.UserMetadata) map[ScopedKey][]byte {
	metadataMap := make(map[ScopedKey][]byte, len(metadataList))
	for _, entry := range metadataList {
		metadataMap[ScopedKey(entry.Key)] = entry.Value
	}

	return metadataMap
}
