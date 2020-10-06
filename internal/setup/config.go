package setup

type IAMSetUp struct {
	GlobalOrg          string
	IAMProject         string
	DefaultLoginPolicy LoginPolicy
	DefaultLabelPolicy LabelPolicy
	Orgs               []Org
	Owners             []string
}

type LoginPolicy struct {
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIdp      bool
}
type LabelPolicy struct {
	PrimaryColor   string
	SecundaryColor string
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
	Users        []User
	Owners       []string
	Projects     []Project
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
	DevMode                bool
}
