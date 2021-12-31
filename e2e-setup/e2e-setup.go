package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func main() {
	registerOrg(getCsrfTokenAndCookies())
}

func registerOrg(csrfToken string, cookies []*http.Cookie) {

	req, err := http.NewRequest("POST", "http://localhost:50003/login/register/org", strings.NewReader(url.Values(map[string][]string{
		"gorilla.csrf.Token":             {csrfToken},
		"orgname":                        {"e2e-org"},
		"firstname":                      {"e2e-firstname"},
		"lastname":                       {"e2e-lastname"},
		"username":                       {"e2e-username"},
		"email":                          {"e2e@example.com"},
		"register-password":              {"e2e-Passw0rd"},
		"register-password-confirmation": {"e2e-Passw0rd"},
		"register-term-confirmation":     {"on"},
	}).Encode()))

	for i := range cookies {
		req.AddCookie(cookies[i])
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	panicOnErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	panicOnErr(err)
	fmt.Println(string(body))
}

func getCsrfTokenAndCookies() (string, []*http.Cookie) {
	resp, err := http.Get("http://localhost:50003/login/register/org")
	defer resp.Body.Close()
	panicOnErr(err)

	node, err := html.Parse(resp.Body)
	panicOnErr(err)

	selector := `[name="gorilla.csrf.Token"]`
	csrfTokenInput := cascadia.Query(node, cascadia.MustCompile(selector))
	for i := range csrfTokenInput.Attr {
		attr := csrfTokenInput.Attr[i]
		if attr.Key == "value" {
			return attr.Val, resp.Cookies()
		}
	}
	panic("node not found using selector " + selector)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
