package auth

import (
	"golang.org/x/oauth2"
)

// PATTokenSource returns an oauth2.TokenSource that uses a Personal Access Token as a Bearer token.
func PATTokenSource(pat string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: pat,
		TokenType:   "Bearer",
	})
}
