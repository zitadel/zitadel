package saml

import "github.com/caos/zitadel/internal/api/saml/xml/protocol/saml"

type AttributeSetter interface {
	SetEmail(string)
	SetFullName(string)
	SetGivenName(string)
	SetSurname(string)
	SetUserID(string)
	SetUsername(string)
	SetApplicationID(string)
}

const (
	AttributeEmail int = iota
	AttributeFullName
	AttributeGivenName
	AttributeSurname
	AttributeUsername
	AttributeUserID
	AttributeApplicationID
)

type Attributes struct {
	email         string
	fullName      string
	givenName     string
	surname       string
	userID        string
	username      string
	applicationID string
}

var _ AttributeSetter = &Attributes{}

func (a *Attributes) GetNameID() *saml.NameIDType {
	return &saml.NameIDType{
		Format: "urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress",
		Text:   a.username,
	}
}

func (a *Attributes) SetEmail(value string) {
	a.email = value
}

func (a *Attributes) SetFullName(value string) {
	a.fullName = value
}

func (a *Attributes) SetGivenName(value string) {
	a.givenName = value
}

func (a *Attributes) SetSurname(value string) {
	a.surname = value
}

func (a *Attributes) SetUsername(value string) {
	a.username = value
}

func (a *Attributes) SetUserID(value string) {
	a.userID = value
}

func (a *Attributes) SetApplicationID(value string) {
	a.applicationID = value
}

func (a *Attributes) GetSAML() []*saml.AttributeType {
	attrs := make([]*saml.AttributeType, 0)
	if a.email != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "Email",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.email},
		})
	}
	if a.surname != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "SurName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.surname},
		})
	}
	if a.givenName != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "FirstName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.givenName},
		})
	}
	if a.fullName != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "FullName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.fullName},
		})
	}
	if a.username != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "UserName",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.username},
		})
	}
	if a.userID != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "UserID",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.userID},
		})
	}
	if a.applicationID != "" {
		attrs = append(attrs, &saml.AttributeType{
			Name:           "ApplicationID",
			NameFormat:     "urn:oasis:names:tc:SAML:2.0:attrname-format:basic",
			AttributeValue: []string{a.applicationID},
		})
	}
	return attrs
}
