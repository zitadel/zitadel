package zitadel

import (
	"github.com/caos/zitadel/operator/docu"
	"path"
	"runtime"
)

func GetDocuInfo() (string, []*docu.Version) {
	_, filename, _, ok := runtime.Caller(0)

	if ok {
		return path.Dir(filename), []*docu.Version{
			{
				Struct:   "DesiredV0",
				Version:  "V0",
				SubKinds: map[string]string{},
			},
		}
	}
	return "", nil
}
