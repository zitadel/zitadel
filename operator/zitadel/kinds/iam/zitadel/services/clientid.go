package services

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetClientIDFunc(
	namespace string,
	httpServiceName string,
	httpPort int,
) func() string {
	return func() string {
		resp, err := http.Get("http://" + httpServiceName + "." + namespace + ":" + strconv.Itoa(httpPort) + "/clientID")
		if err != nil || resp.StatusCode >= 400 {
			return ""
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		return strings.TrimSuffix(strings.TrimPrefix(string(body), "\""), "\"")
	}
}
