package user

func (r *AddHumanUserRequest) AuthContext() string {
	return r.GetOrganisation().GetOrgId()
}
