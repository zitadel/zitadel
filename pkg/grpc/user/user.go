package user

type SearchQuery_ResourceOwner struct {
	ResourceOwner *ResourceOwnerQuery
}

func (SearchQuery_ResourceOwner) isSearchQuery_Query() {}

type ResourceOwnerQuery struct {
	OrgID string
}

type MembershipType = isMembership_Type
