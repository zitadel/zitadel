package apple

import (
	"context"
	"encoding/json"

	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

// Session extends the [oidc.Session] with the formValues returned from the callback.
// This enables to parse the user (name and email), which Apple only returns as form params on registration
type Session struct {
	*oidc.Session
	UserFormValue string
}

type userFormValue struct {
	Name userNamesFormValue `json:"name,omitempty" schema:"name"`
}

type userNamesFormValue struct {
	FirstName string `json:"firstName,omitempty" schema:"firstName"`
	LastName  string `json:"lastName,omitempty" schema:"lastName"`
}

// FetchUser implements the [idp.Session] interface.
// It will execute an OIDC code exchange if needed to retrieve the tokens,
// extract the information from the id_token and if available also from the `user` form value.
// The information will be mapped into an [idp.User].
func (s *Session) FetchUser(ctx context.Context) (user idp.User, err error) {
	if s.Tokens == nil {
		if err = s.Authorize(ctx); err != nil {
			return nil, err
		}
	}
	info := s.Tokens.IDTokenClaims.GetUserInfo()
	userName := userFormValue{}
	if s.UserFormValue != "" {
		if err = json.Unmarshal([]byte(s.UserFormValue), &userName); err != nil {
			return nil, err
		}
	}

	return NewUser(info, userName.Name), nil
}

func NewUser(info *openid.UserInfo, names userNamesFormValue) *User {
	user := oidc.NewUser(info)
	user.GivenName = names.FirstName
	user.FamilyName = names.LastName
	return &User{User: user}
}

// User extends the [oidc.User] by returning the email as preferred_username, since Apple does not return the latter.
type User struct {
	*oidc.User
}

func (u *User) GetPreferredUsername() string {
	return u.Email
}
