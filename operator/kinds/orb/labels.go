package orb

import "github.com/caos/orbos/pkg/labels"

func mustDatabaseOperator(binaryVersion string) *labels.Operator {
	return labels.MustForOperator("ZITADEL", "zitadel.caos.ch", binaryVersion)
}
