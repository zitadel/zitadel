package v4

type Org struct {
	InstanceID string
	ID         string
	Name       string
}

type GetOrg struct{}

type ListOrgs struct{}

type CreateOrg struct{}

type UpdateOrg struct{}

type DeleteOrg struct{}
