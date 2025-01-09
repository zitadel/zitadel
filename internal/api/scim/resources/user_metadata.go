package resources

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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
	case metadata.KeyEntitlements:
		fallthrough
	case metadata.KeyIms:
		fallthrough
	case metadata.KeyPhotos:
		fallthrough
	case metadata.KeyAddresses:
		fallthrough
	case metadata.KeyRoles:
		return json.Marshal(value)

	// http url values
	case metadata.KeyProfileUrl:
		return []byte(value.(*schemas.HttpURL).String()), nil

	// raw values
	case metadata.KeyProvisioningDomain:
		fallthrough
	case metadata.KeyExternalId:
		fallthrough
	case metadata.KeyMiddleName:
		fallthrough
	case metadata.KeyHonorificSuffix:
		fallthrough
	case metadata.KeyHonorificPrefix:
		fallthrough
	case metadata.KeyTitle:
		fallthrough
	case metadata.KeyLocale:
		fallthrough
	case metadata.KeyTimezone:
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
	case metadata.KeyProvisioningDomain:
		break
	}

	logging.Panicf("Unknown or unsupported metadata key %s", key)
	return nil
}
