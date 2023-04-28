package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
)

type DetailsMsg interface {
	GetDetails() *object.Details
}

func AssertDetails[D DetailsMsg](t testing.TB, exptected, actual D) {
	wantDetails, gotDetails := exptected.GetDetails(), actual.GetDetails()

	if wantDetails != nil {
		assert.NotZero(t, gotDetails.GetSequence())
	}
	wantCD, gotCD := wantDetails.GetChangeDate().AsTime(), gotDetails.GetChangeDate().AsTime()
	assert.WithinRange(t, gotCD, wantCD, wantCD.Add(time.Minute))
	assert.Equal(t, wantDetails.GetResourceOwner(), gotDetails.GetResourceOwner())
}
