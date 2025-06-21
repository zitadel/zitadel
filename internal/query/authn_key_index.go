package query

import "github.com/zitadel/zitadel/internal/cachekey"


type authnKeyIndex = cachekey.AuthnKeyIndex

const (
	InstanceID = cachekey.InstanceID
	UserType   = cachekey.UserType
	KeyID      = cachekey.KeyID
)
