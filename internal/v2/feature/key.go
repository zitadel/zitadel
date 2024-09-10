package feature

import (
	feature "github.com/zitadel/zitadel/internal/feature"
)

type Key = feature.Key

var Keys = feature.KeyValues()

func keyFromString(text string) Key {
	key, _ := feature.KeyString(text)
	return key
}
