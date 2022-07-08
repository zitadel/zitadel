package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, fmt.Errorf("listing users returned status code %d and body %s", resp.StatusCode, string(bodyBytes))
	}

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
