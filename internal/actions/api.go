package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/text/language"

	"github.com/dop251/goja"
	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"github.com/zitadel/zitadel/internal/domain"
)

type apiParam struct {
	runtime *goja.Runtime
	parameter
}

type APIOption func(*apiParam)

func SetHuman(human *domain.Human) APIOption {
	return func(a *apiParam) {
		a.set("setFirstName", func(firstName string) {
			human.FirstName = firstName
		})
		a.set("setLastName", func(lastName string) {
			human.LastName = lastName
		})
		a.set("setNickName", func(nickName string) {
			human.NickName = nickName
		})
		a.set("setDisplayName", func(displayName string) {
			human.DisplayName = displayName
		})
		a.set("setPreferredLanguage", func(preferredLanguage string) {
			human.PreferredLanguage = language.Make(preferredLanguage)
		})
		a.set("setGender", func(gender domain.Gender) {
			human.Gender = gender
		})
		a.set("setUsername", func(username string) {
			human.Username = username
		})
		a.set("setEmail", func(email string) {
			if human.Email == nil {
				human.Email = &domain.Email{}
			}
			human.Email.EmailAddress = email
		})
		a.set("setEmailVerified", func(verified bool) {
			if human.Email == nil {
				return
			}
			human.Email.IsEmailVerified = verified
		})
		a.set("setPhone", func(email string) {
			if human.Phone == nil {
				human.Phone = &domain.Phone{}
			}
			human.Phone.PhoneNumber = email
		})
		a.set("setPhoneVerified", func(verified bool) {
			if human.Phone == nil {
				return
			}
			human.Phone.IsPhoneVerified = verified
		})
	}
}

func SetExternalUser(user *domain.ExternalUser) APIOption {
	return func(a *apiParam) {
		a.set("setFirstName", func(firstName string) {
			user.FirstName = firstName
		})
		a.set("setLastName", func(lastName string) {
			user.LastName = lastName
		})
		a.set("setNickName", func(nickName string) {
			user.NickName = nickName
		})
		a.set("setDisplayName", func(displayName string) {
			user.DisplayName = displayName
		})
		a.set("setPreferredLanguage", func(preferredLanguage string) {
			user.PreferredLanguage = language.Make(preferredLanguage)
		})
		a.set("setPreferredUsername", func(username string) {
			user.PreferredUsername = username
		})
		a.set("setEmail", func(email string) {
			user.Email = email
		})
		a.set("setEmailVerified", func(verified bool) {
			user.IsEmailVerified = verified
		})
		a.set("setPhone", func(phone string) {
			user.Phone = phone
		})
		a.set("setPhoneVerified", func(verified bool) {
			user.IsPhoneVerified = verified
		})
	}
}

func SetMetadata(metadata *[]*domain.Metadata) APIOption {
	return func(a *apiParam) {
		a.set("metadata", metadata)
	}
}

func SetUserGrants(usergrants *[]UserGrant) APIOption {
	return func(a *apiParam) {
		a.set("userGrants", usergrants)
	}
}

func SetClaims(claims *map[string]interface{}, logs *[]string) APIOption {
	return func(a *apiParam) {
		if *claims == nil {
			*claims = make(map[string]interface{})
		}
		a.setPath([]string{"v1", "claims", "set"}, func(key string, value interface{}) {
			if _, ok := (*claims)[key]; !ok {
				(*claims)[key] = value
				return
			}
			*logs = append(*logs, fmt.Sprintf("key %q already exists", key))
		})
		a.setPath([]string{"v1", "claims", "appendLogIntoClaims"}, func(entry string) {
			*logs = append(*logs, entry)
		})
	}
}

func SetUserinfo(userInfo oidc.UserInfoSetter, logs *[]string) APIOption {
	return func(a *apiParam) {
		a.setPath([]string{"v1", "userinfo", "setClaim"}, func(key string, value interface{}) {
			if userInfo.GetClaim(key) == nil {
				userInfo.AppendClaims(key, value)
				return
			}
			*logs = append(*logs, fmt.Sprintf("key %q already exists", key))
		})
		a.setPath([]string{"v1", "userinfo", "appendLogIntoClaims"}, func(entry string) {
			*logs = append(*logs, entry)
		})
	}
}

type userMetadataSetter interface {
	SetUserMetadata(ctx context.Context, metadata *domain.Metadata, userID, resourceOwner string) (*domain.Metadata, error)
}

func SetUserMetadataSetter(ctx context.Context, setter userMetadataSetter, userID, resourceOwner string) APIOption {
	return func(c *apiParam) {
		c.setPath([]string{"v1", "user", "setMetadata"}, func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) != 2 {
				panic("exactly 2 (key, value) arguments expected")
			}
			key := call.Arguments[0].Export().(string)
			val := call.Arguments[1].Export()

			value, err := json.Marshal(val)
			if err != nil {
				logging.WithError(err).Debug("unable to marshal")
				panic(err)
			}

			metadata := &domain.Metadata{
				Key:   key,
				Value: value,
			}
			if _, err = setter.SetUserMetadata(ctx, metadata, userID, resourceOwner); err != nil {
				logging.WithError(err).Info("unable to set md in action")
				panic(err)
			}
			return nil
		})
	}
}
