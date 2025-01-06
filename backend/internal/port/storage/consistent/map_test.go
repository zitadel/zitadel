package consistent

import "github.com/zitadel/zitadel/backend/internal/port"

var _ port.Object = (*testInstance)(nil)

type testInstance struct {
	ID      string        `consistent:"id,pk"`
	Name    string        `consistent:"name,pk"`
	Domains []*testDomain `consistent:"domains"`
}

func (i *testInstance) Columns() []*port.Column {
	return []*port.Column{
		{Name: "id", Value: i.ID},
		{Name: "name", Value: i.Name},
		{Name: "domains", Value: i.Domains},
	}
}

var _ port.Object = (*testDomain)(nil)

type testDomain struct {
	Name       string `consistent:"name,pk"`
	InstanceID string `consistent:"instance_id,pk,fk"`
	IsVerified bool   `consistent:"is_verified"`
}

func (d *testDomain) Columns() []*port.Column {
	return []*port.Column{
		{Name: "name", Value: d.Name},
		{Name: "instance_id", Value: d.InstanceID},
		{Name: "is_verified", Value: d.IsVerified},
	}
}
