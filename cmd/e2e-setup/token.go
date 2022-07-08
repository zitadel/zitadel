package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func newToken(cfg E2EConfig) (string, error) {

	keyBytes, err := os.ReadFile(cfg.MachineKeyPath)
	if err != nil {
		return "", err
	}

	key := struct{ UserId, KeyId, Key string }{}
	if err := json.Unmarshal(keyBytes, &key); err != nil {
		return "", err
	}

	if key.KeyId == "" ||
		key.UserId == "" ||
		key.Key == "" {
		return "", fmt.Errorf("key is incomplete: %+v", key)
	}

	now := time.Now()
	iat := now.Unix()
	exp := now.Add(55 * time.Minute).Unix()

	audience := cfg.Audience
	if audience == "" {
		audience = cfg.IssuerURL
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": key.UserId,
		"sub": key.UserId,
		"aud": audience,
		"iat": iat,
		"exp": exp,
	})

	token.Header["alg"] = "RS256"
	token.Header["kid"] = key.KeyId

	rsaKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key.Key))
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(rsaKey)
	if err != nil {
		return "", err
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/oauth/v2/token", cfg.APIURL), map[string][]string{
		"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
		"scope":      {fmt.Sprintf("openid urn:zitadel:iam:org:project:id:%s:aud", strings.TrimPrefix(cfg.ZitadelProjectResourceID, "bignumber-"))},
		"assertion":  {tokenString},
	})
	if err != nil {
		return "", err
	}

	tokenBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return "", fmt.Errorf("getting token returned status code %d and body %s", resp.StatusCode, string(tokenBody))
	}

	tokenResp := struct {
		AccessToken string `json:"access_token"`
	}{}

	return tokenResp.AccessToken, json.Unmarshal(tokenBody, &tokenResp)
}
