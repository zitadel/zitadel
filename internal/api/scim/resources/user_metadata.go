package resources

import (
	"context"
	"encoding/json"
	"time"
	// import timezone database to ensure it is available at runtime
	// data is required to validate time zones.
	_ "time/tzdata"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (h *UsersHandler) queryMetadataForUsers(ctx context.Context, userIds []string) (map[string]map[metadata.ScopedKey][]byte, error) {
	queries := BuildMetadataQueries(ctx, metadata.ScimUserRelevantMetadataKeys)

	md, err := h.query.SearchUserMetadataForUsers(ctx, false, userIds, queries)
	if err != nil {
		return nil, err
	}

	metadataMap := make(map[string]map[metadata.ScopedKey][]byte, len(md.Metadata))
	for _, entry := range md.Metadata {
		userMetadata, ok := metadataMap[entry.UserID]
		if !ok {
			userMetadata = make(map[metadata.ScopedKey][]byte)
			metadataMap[entry.UserID] = userMetadata
		}

		userMetadata[metadata.ScopedKey(entry.Key)] = entry.Value
	}

	return metadataMap, nil
}

func (h *UsersHandler) queryMetadataForUser(ctx context.Context, id string) (map[metadata.ScopedKey][]byte, error) {
	queries := BuildMetadataQueries(ctx, metadata.ScimUserRelevantMetadataKeys)

	md, err := h.query.SearchUserMetadata(ctx, false, id, queries, nil)
	if err != nil {
		return nil, err
	}

	return metadata.MapListToScopedKeyMap(md.Metadata), nil
}

func BuildMetadataQueries(ctx context.Context, metadataKeys []metadata.Key) *query.UserMetadataSearchQueries {
	keyQueries := make([]query.SearchQuery, len(metadataKeys))
	for i, key := range metadataKeys {
		keyQueries[i] = buildMetadataKeyQuery(ctx, key)
	}

	queries := &query.UserMetadataSearchQueries{
		SearchRequest: query.SearchRequest{},
		Queries:       []query.SearchQuery{query.Or(keyQueries...)},
	}
	return queries
}

func buildMetadataKeyQuery(ctx context.Context, key metadata.Key) query.SearchQuery {
	scopedKey := metadata.ScopeKey(ctx, key)
	q, err := query.NewUserMetadataKeySearchQuery(string(scopedKey), query.TextEquals)
	if err != nil {
		logging.Panic("Error build user metadata query for key " + key)
	}

	return q
}

func (h *UsersHandler) mapMetadataToDomain(ctx context.Context, user *ScimUser) (md []*domain.Metadata, skippedMetadata []string, err error) {
	md = make([]*domain.Metadata, 0, len(metadata.ScimUserRelevantMetadataKeys))
	for _, key := range metadata.ScimUserRelevantMetadataKeys {
		var value []byte
		value, err = getValueForMetadataKey(user, key)
		if err != nil {
			return
		}

		if len(value) > 0 {
			md = append(md, &domain.Metadata{
				Key:   string(metadata.ScopeKey(ctx, key)),
				Value: value,
			})
		} else {
			skippedMetadata = append(skippedMetadata, string(metadata.ScopeKey(ctx, key)))
		}
	}

	return
}

func (h *UsersHandler) mapMetadataToCommands(ctx context.Context, user *ScimUser) ([]*command.AddMetadataEntry, error) {
	md := make([]*command.AddMetadataEntry, 0, len(metadata.ScimUserRelevantMetadataKeys))
	for _, key := range metadata.ScimUserRelevantMetadataKeys {
		value, err := getValueForMetadataKey(user, key)
		if err != nil {
			return nil, err
		}

		if len(value) > 0 {
			md = append(md, &command.AddMetadataEntry{
				Key:   string(metadata.ScopeKey(ctx, key)),
				Value: value,
			})
		}
	}

	return md, nil
}

func getValueForMetadataKey(user *ScimUser, key metadata.Key) ([]byte, error) {
	value := getRawValueForMetadataKey(user, key)
	if value == nil {
		return nil, nil
	}

	switch key {
	// json values
	case metadata.KeyRoles,
		metadata.KeyAddresses,
		metadata.KeyEntitlements,
		metadata.KeyIms,
		metadata.KeyPhotos,
		metadata.KeyEmails:
		val, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		// null is considered no value
		if len(val) == 4 && string(val) == "null" {
			return nil, nil
		}

		return val, nil

	// http url values
	case metadata.KeyProfileUrl:
		return []byte(value.(*schemas.HttpURL).String()), nil

	// raw values
	case metadata.KeyTimezone,
		metadata.KeyLocale,
		metadata.KeyTitle,
		metadata.KeyHonorificPrefix,
		metadata.KeyHonorificSuffix,
		metadata.KeyMiddleName,
		metadata.KeyExternalId,
		metadata.KeyProvisioningDomain:
		valueStr := value.(string)
		if valueStr == "" {
			return nil, nil
		}

		return []byte(valueStr), validateValueForMetadataKey(valueStr, key)
	}

	logging.Panicf("Unknown metadata key %s", key)
	return nil, nil
}

func validateValueForMetadataKey(v string, key metadata.Key) error {
	//nolint:exhaustive
	switch key {
	case metadata.KeyLocale:
		if _, err := language.Parse(v); err != nil {
			return serrors.ThrowInvalidValue(zerrors.ThrowInvalidArgument(err, "SCIM-MD11", "Could not parse locale"))
		}
		return nil
	case metadata.KeyTimezone:
		if _, err := time.LoadLocation(v); err != nil {
			return serrors.ThrowInvalidValue(zerrors.ThrowInvalidArgument(err, "SCIM-MD12", "Could not parse timezone"))
		}

		return nil
	}

	return nil
}

func getRawValueForMetadataKey(user *ScimUser, key metadata.Key) interface{} {
	switch key {
	case metadata.KeyIms:
		return user.Ims
	case metadata.KeyPhotos:
		return user.Photos
	case metadata.KeyAddresses:
		return user.Addresses
	case metadata.KeyEntitlements:
		return user.Entitlements
	case metadata.KeyRoles:
		return user.Roles
	case metadata.KeyMiddleName:
		if user.Name == nil {
			return ""
		}
		return user.Name.MiddleName
	case metadata.KeyHonorificPrefix:
		if user.Name == nil {
			return ""
		}
		return user.Name.HonorificPrefix
	case metadata.KeyHonorificSuffix:
		if user.Name == nil {
			return ""
		}
		return user.Name.HonorificSuffix
	case metadata.KeyExternalId:
		return user.ExternalID
	case metadata.KeyProfileUrl:
		return user.ProfileUrl
	case metadata.KeyTitle:
		return user.Title
	case metadata.KeyLocale:
		return user.Locale
	case metadata.KeyTimezone:
		return user.Timezone
	case metadata.KeyEmails:
		return user.Emails
	case metadata.KeyProvisioningDomain:
		break
	}

	logging.Panicf("Unknown or unsupported metadata key %s", key)
	return nil
}

func extractScalarMetadata(ctx context.Context, md map[metadata.ScopedKey][]byte, key metadata.Key) string {
	val, ok := md[metadata.ScopeKey(ctx, key)]
	if !ok {
		return ""
	}

	return string(val)
}

func extractHttpURLMetadata(ctx context.Context, md map[metadata.ScopedKey][]byte, key metadata.Key) *schemas.HttpURL {
	val, ok := md[metadata.ScopeKey(ctx, key)]
	if !ok {
		return nil
	}

	url, err := schemas.ParseHTTPURL(string(val))
	if err != nil {
		logging.OnError(err).Warn("Failed to parse scim url metadata for " + key)
		return nil
	}

	return url
}

func extractJsonMetadata(ctx context.Context, md map[metadata.ScopedKey][]byte, key metadata.Key, v interface{}) error {
	val, ok := md[metadata.ScopeKey(ctx, key)]
	if !ok {
		return nil
	}

	return json.Unmarshal(val, v)
}
