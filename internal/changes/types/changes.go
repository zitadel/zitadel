package types

type User struct {
	FirstName    string `json:"firstName,omitempty"`
	LastName     string `json:"lastName,omitempty"`
	EMailAddress string `json:"email,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Language     string `json:"preferredLanguage,omitempty"`
	UserName     string `json:"userName,omitempty"`
}

type Project struct {
	Name string `json:"name,omitempty"`
}

type omitempty struct {
	ClientId string `json:"clientId,omitempty"`
}

type App struct {
	AppId          string `json:"AppId,omitempty"`
	AppType        int    `json:"AppType,omitempty"`
	AuthMethodType int    `json:"authMethodType,omitempty"`
	ClientId       string `json:"clientId,omitempty"`
	Name           string `json:"name,omitempty"`
}

type Org struct {
	Domain string   `json:"domain,omitempty"`
	Name   string   `json:"name,omitempty"`
	Roles  []string `json:"roles,omitempty"`
	UserId string   `json:"userId,omitempty"`
}
