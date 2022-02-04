package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func awaitConsistency(ctx context.Context, cfg E2EConfig, expectUsers []user) (err error) {

	retry := make(chan struct{})
	go func() {
		// trigger first check
		retry <- struct{}{}
	}()
	for {
		select {
		case <-retry:
			err = checkCondition(ctx, cfg, expectUsers)
			if err == nil {
				fmt.Println("setup is consistent")
				return nil
			}
			fmt.Printf("setup is not consistent yet, retrying in a second: %s\n", err)
			time.Sleep(time.Second)
			go func() {
				retry <- struct{}{}
			}()
		case <-ctx.Done():
			return fmt.Errorf("setup failed to come to a consistent state: %s: %w", ctx.Err(), err)
		}
	}
}

func checkCondition(ctx context.Context, cfg E2EConfig, expectUsers []user) error {
	token, err := newToken(cfg)
	if err != nil {
		return err
	}

	foundUsers, err := listUsers(ctx, cfg.APIURL, token)
	if err != nil {
		return err
	}

	var awaitingUsers []string
expectLoop:
	for i := range expectUsers {
		expectUser := expectUsers[i]
		for j := range foundUsers {
			foundUser := foundUsers[j]
			if expectUser.desc+"_user_name" == foundUser {
				continue expectLoop
			}
		}
		awaitingUsers = append(awaitingUsers, expectUser.desc)
	}

	if len(awaitingUsers) > 0 {
		return fmt.Errorf("users %v are not consistent yet", awaitingUsers)
	}
	return nil
}

func listUsers(ctx context.Context, apiUrl, token string) ([]string, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl+"/management/v1/users/_search", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	unmarshalledResp := struct {
		Result []struct {
			UserName string `json:"userName"`
		}
	}{}

	if err = json.Unmarshal(bodyBytes, &unmarshalledResp); err != nil {
		return nil, err
	}

	users := make([]string, len(unmarshalledResp.Result))
	for i := range unmarshalledResp.Result {
		users[i] = unmarshalledResp.Result[i].UserName
	}

	return users, nil
}

func newToken(cfg E2EConfig) (string, error) {

	keyBytes, err := os.ReadFile(cfg.MachineKeyPath)
	if err != nil {
		return "", err
	}

	key := struct {
		UserId, KeyId string
		Key           string
	}{}
	if err := json.Unmarshal(keyBytes, &key); err != nil {
		return "", err
	}

	now := time.Now()
	iat := now.Unix()
	exp := now.Add(55 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": key.UserId,
		"sub": key.UserId,
		"aud": cfg.IssuerURL,
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

	tokenResp := struct {
		AccessToken string `json:"access_token"`
	}{}

	return tokenResp.AccessToken, json.Unmarshal(tokenBody, &tokenResp)
}
