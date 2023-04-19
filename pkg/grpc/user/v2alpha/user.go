package user

func (r *AddUserRequest) AuthContext() string {
	return r.GetOrganisation().GetOrgId()
}
