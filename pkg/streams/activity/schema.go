package activity

import "github.com/zitadel/zitadel/pkg/streams"

const (
	LogFieldKeyOrgID         streams.LogFieldKey = "orgID"
	LogFieldKeyUserID        streams.LogFieldKey = "userID"
	LogFieldKeyDomain        streams.LogFieldKey = "domain"
	LogFieldKeyTrigger       streams.LogFieldKey = "trigger"
	LogFieldKeyMethod        streams.LogFieldKey = "method"
	LogFieldKeyPath          streams.LogFieldKey = "path"
	LogFieldKeyRequestMethod streams.LogFieldKey = "requestMethod"
	LogFieldKeyIsSystemUser  streams.LogFieldKey = "isSystemUser"
	LogFieldKeyGRPCStatus    streams.LogFieldKey = "grpcStatus"
	LogFieldKeyHTTPStatus    streams.LogFieldKey = "httpStatus"
)

type TriggerMethod int

const (
	Unspecified TriggerMethod = iota
	ResourceAPI
	OIDCAccessToken
	OIDCRefreshToken
	SessionAPI
	SAMLResponse
)

func (t TriggerMethod) String() string {
	switch t {
	case Unspecified:
		return "unspecified"
	case ResourceAPI:
		return "resourceAPI"
	case OIDCRefreshToken:
		return "refreshToken"
	case OIDCAccessToken:
		return "accessToken"
	case SessionAPI:
		return "sessionAPI"
	case SAMLResponse:
		return "samlResponse"
	default:
		return "unknown"
	}
}
