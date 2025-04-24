package instance

import (
	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/instance/v2"
)

func ToProtoObject(inst *query.Instance) *instance.Instance {
	return &instance.Instance{
		Id:      inst.ID,
		Name:    inst.Name,
		Domains: DomainsToPb(inst.Domains),
		Version: build.Version(),
		Details: object.ToViewDetailsPb(inst.Sequence, inst.CreationDate, inst.ChangeDate, inst.ID),
	}
}

func DomainsToPb(domains []*query.InstanceDomain) []*instance.Domain {
	d := []*instance.Domain{}
	for _, dm := range domains {
		pbDomain := DomainToPb(dm)
		d = append(d, pbDomain)
	}
	return d
}

func DomainToPb(d *query.InstanceDomain) *instance.Domain {
	return &instance.Domain{
		Domain:    d.Domain,
		Primary:   d.IsPrimary,
		Generated: d.IsGenerated,
		Details: object.ToViewDetailsPb(
			d.Sequence,
			d.CreationDate,
			d.ChangeDate,
			d.InstanceID,
		),
	}
}
