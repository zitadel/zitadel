package command

import "github.com/zitadel/zitadel/internal/repository/idp"

type GoogleProvider struct {
	Name         string
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type LDAPProvider struct {
	Name                string
	Host                string
	Port                string
	TLS                 bool
	BaseDN              string
	UserObjectClass     string
	UserUniqueAttribute string
	Admin               string
	Password            string
	LDAPAttributes      idp.LDAPAttributes
	IDPOptions          idp.Options
}
