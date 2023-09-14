package integration

import (
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

type myMsg struct {
	details *object.Details
}

func (m myMsg) GetDetails() *object.Details {
	return m.details
}

func TestAssertDetails(t *testing.T) {
	tests := []struct {
		name      string
		exptected myMsg
		actual    myMsg
	}{
		{
			name:      "nil",
			exptected: myMsg{},
			actual:    myMsg{},
		},
		{
			name: "values",
			exptected: myMsg{
				details: &object.Details{
					ResourceOwner: "me",
				},
			},
			actual: myMsg{
				details: &object.Details{
					Sequence:      123,
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: "me",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertDetails(t, tt.exptected, tt.actual)
		})
	}
}
