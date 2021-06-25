package kinds

import (
	"github.com/caos/zitadel/operator/docu"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func GetDocuInfo() []*docu.Type {
	path, orbVersions := orb.GetDocuInfo()

	infos := []*docu.Type{{
		Name: "orb",
		Kinds: []*docu.Info{
			{
				Path:     path,
				Kind:     "orbiter.caos.ch/Orb",
				Versions: orbVersions,
			},
		},
	}}

	infos = append(infos, iam.GetDocuInfo()...)
	return infos
}
