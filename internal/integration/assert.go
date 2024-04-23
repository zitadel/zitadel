package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

// Details is the interface that covers both v1 and v2 proto generated object details.
type Details interface {
	comparable
	GetSequence() uint64
	GetChangeDate() *timestamppb.Timestamp
	GetResourceOwner() string
}

// DetailsMsg is the interface that covers all proto messages which contain v1 or v2 object details.
type DetailsMsg[D Details] interface {
	GetDetails() D
}

type ListDetailsMsg interface {
	GetDetails() *object.ListDetails
}

// AssertDetails asserts values in a message's object Details,
// if the object Details in expected is a non-nil value.
// It targets API v2 messages that have the `GetDetails()` method.
//
// Dynamically generated values are not compared with expected.
// Instead a sanity check is performed.
// For the sequence a non-zero value is expected.
// If the change date is populated, it is checked with a tolerance of 1 minute around Now.
//
// The resource owner is compared with expected.
func AssertDetails[D Details, M DetailsMsg[D]](t testing.TB, expected, actual M) {
	wantDetails, gotDetails := expected.GetDetails(), actual.GetDetails()
	var nilDetails D
	if wantDetails == nilDetails {
		assert.Nil(t, gotDetails)
		return
	}

	assert.NotZero(t, gotDetails.GetSequence())

	if wantDetails.GetChangeDate() != nil {
		wantChangeDate := time.Now()
		gotChangeDate := gotDetails.GetChangeDate().AsTime()
		assert.WithinRange(t, gotChangeDate, wantChangeDate.Add(-time.Minute), wantChangeDate.Add(time.Minute))
	}

	assert.Equal(t, wantDetails.GetResourceOwner(), gotDetails.GetResourceOwner())
}

func AssertListDetails[D ListDetailsMsg](t testing.TB, expected, actual D) {
	wantDetails, gotDetails := expected.GetDetails(), actual.GetDetails()
	if wantDetails == nil {
		assert.Nil(t, gotDetails)
		return
	}

	assert.Equal(t, wantDetails.GetTotalResult(), gotDetails.GetTotalResult())

	gotCD := gotDetails.GetTimestamp().AsTime()
	wantCD := time.Now()
	assert.WithinRange(t, gotCD, wantCD.Add(-time.Minute), wantCD.Add(time.Minute))
}
