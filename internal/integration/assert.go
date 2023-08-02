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

// AssertDetails asserts values in a message's object Details,
// if the object Details in expected is a non-nil value.
// It targets API v2 messages that have the `GetDetails()` method.
//
// Dynamically generated values are not compared with expected.
// Instead a sanity check is performed.
// For the sequence a non-zero value is expected.
// The change date has to be now, with a tollerance of 1 second.
//
// The resource owner is compared with expected and is
// therefore the only value that has to be set.
func AssertDetails[D DetailsMsg](t testing.TB, exptected, actual D) {
	wantDetails, gotDetails := exptected.GetDetails(), actual.GetDetails()
	if wantDetails == nil {
		assert.Nil(t, gotDetails)
		return
	}

	assert.NotZero(t, gotDetails.GetSequence())

	gotCD := gotDetails.GetChangeDate().AsTime()
	now := time.Now()
	assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))

	assert.Equal(t, wantDetails.GetResourceOwner(), gotDetails.GetResourceOwner())
}
