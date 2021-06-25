package iam

import (
	"github.com/caos/zitadel/operator/docu"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel"
)

func GetDocuInfo() []*docu.Type {
	path, versions := zitadel.GetDocuInfo()
	types := []*docu.Type{{
		Name: "iam",
		Kinds: []*docu.Info{{
			Path:     path,
			Kind:     zitadelKind,
			Versions: versions,
		}},
	}}
	return types
}
