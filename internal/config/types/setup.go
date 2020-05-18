package types

type IAMSetUp struct {
	GlobalOrg  string
	IAMProject string
	Orgs       []Org
	Owners     []string
}

type User struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	Password  string
}

type Org struct {
	Name     string
	Domain   string
	Users    []User
	Owners   []string
	Projects []Project
}

type Project struct {
	Name     string
	Users    []User
	Members  []string
	OIDCApps []OIDCApp
}

type OIDCApp struct {
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	ApplicationType        string
	AuthMethodType         string
	PostLogoutRedirectUris []string
}
