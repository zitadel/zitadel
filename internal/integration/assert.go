package integration

import (
	"testing"
	"time"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type ListDetails interface {
	comparable
	GetTotalResult() uint64
	GetTimestamp() *timestamppb.Timestamp
}

type ListDetailsMsg[L ListDetails] interface {
	GetDetails() L
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

func AssertListDetails[L ListDetails, D ListDetailsMsg[L]](t testing.TB, expected, actual D) {
	wantDetails, gotDetails := expected.GetDetails(), actual.GetDetails()
	var nilDetails L
	if wantDetails == nilDetails {
		assert.Nil(t, gotDetails)
		return
	}
	assert.Equal(t, wantDetails.GetTotalResult(), gotDetails.GetTotalResult())

	if wantDetails.GetTimestamp() != nil {
		gotCD := gotDetails.GetTimestamp().AsTime()
		wantCD := time.Now()
		assert.WithinRange(t, gotCD, wantCD.Add(-time.Minute), wantCD.Add(time.Minute))
	}
}

// EqualProto is inspired by [assert.Equal], only that it tests equality of a proto message.
// A message diff is printed on the error test log if the messages are not equal.
//
// As [assert.Equal] is based on reflection, comparing 2 proto messages sometimes fails,
// due to their internal state.
// Expected messages are usually with a vanilla state, eg only exported fields contain data.
// Actual messages obtained from the gRPC client had unexported fields with data.
// This makes them hard to compare.
func EqualProto(t testing.TB, expected, actual proto.Message) bool {
	t.Helper()
	if proto.Equal(expected, actual) {
		return true
	}
	t.Errorf("Proto messages not equal: %s", diffProto(expected, actual))
	return false
}

func diffProto(expected, actual proto.Message) string {
	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(protojson.Format(expected)),
		B:        difflib.SplitLines(protojson.Format(actual)),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})
	if err != nil {
		panic(err)
	}
	return "\n\nDiff:\n" + diff
}
