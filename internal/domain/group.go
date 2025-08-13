package domain

// Group represents a user group in an organization
type Group struct {
	Name           string
	Description    string
	OrganizationID string

	ID string
}
