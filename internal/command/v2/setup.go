package command

type SetupDefaults struct {
	Org         Org
	Admin       User
	LoginPolicy LoginPolicy
}

type LoginPolicy struct {
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIdp      bool
}

type User struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	Password  string
}

type Org struct {
	Name         string
	Domain       string
	OrgIamPolicy bool
	Owner        User
	Projects     []Project
}

type Project struct {
	Name     string
	Users    []User
	Members  []string
	OIDCApps []OIDCApp
	APIs     []API
}

type OIDCApp struct {
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	ApplicationType        string
	AuthMethodType         string
	PostLogoutRedirectUris []string
	DevMode                bool
}

type API struct {
	Name           string
	AuthMethodType string
}
