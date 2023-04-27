package object

func (o *Organisation) AuthContext() string {
	return o.GetOrgId()
}
