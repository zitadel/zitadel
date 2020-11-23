package configuration

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
)

func getHash(dataMap map[string]string) string {
	data, err := json.Marshal(dataMap)
	if err != nil {
		return ""
	}
	h := sha512.New()
	return base64.URLEncoding.EncodeToString(h.Sum(data))
}
