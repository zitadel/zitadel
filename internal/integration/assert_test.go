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
					ChangeDate:    timestamppb.Now(),
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

func TestPartiallyDeepEqual(t *testing.T) {
	type SecondaryNestedType struct {
		Value int
	}
	type NestedType struct {
		Value         int
		ValueSlice    []int
		Nested        SecondaryNestedType
		NestedPointer *SecondaryNestedType
	}

	type args struct {
		expected interface{}
		actual   interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil",
			args: args{
				expected: nil,
				actual:   nil,
			},
			want: true,
		},
		{
			name: "scalar value",
			args: args{
				expected: 10,
				actual:   10,
			},
			want: true,
		},
		{
			name: "different scalar value",
			args: args{
				expected: 11,
				actual:   10,
			},
			want: false,
		},
		{
			name: "string value",
			args: args{
				expected: "foo",
				actual:   "foo",
			},
			want: true,
		},
		{
			name: "different string value",
			args: args{
				expected: "foo2",
				actual:   "foo",
			},
			want: false,
		},
		{
			name: "scalar only set in actual",
			args: args{
				expected: &SecondaryNestedType{},
				actual:   &SecondaryNestedType{Value: 10},
			},
			want: true,
		},
		{
			name: "scalar equal",
			args: args{
				expected: &SecondaryNestedType{Value: 10},
				actual:   &SecondaryNestedType{Value: 10},
			},
			want: true,
		},
		{
			name: "scalar only set in expected",
			args: args{
				expected: &SecondaryNestedType{Value: 10},
				actual:   &SecondaryNestedType{},
			},
			want: false,
		},
		{
			name: "ptr only set in expected",
			args: args{
				expected: &NestedType{NestedPointer: &SecondaryNestedType{Value: 10}},
				actual:   &NestedType{},
			},
			want: false,
		},
		{
			name: "ptr only set in actual",
			args: args{
				expected: &NestedType{},
				actual:   &NestedType{NestedPointer: &SecondaryNestedType{Value: 10}},
			},
			want: true,
		},
		{
			name: "ptr equal",
			args: args{
				expected: &NestedType{NestedPointer: &SecondaryNestedType{Value: 10}},
				actual:   &NestedType{NestedPointer: &SecondaryNestedType{Value: 10}},
			},
			want: true,
		},
		{
			name: "nested equal",
			args: args{
				expected: &NestedType{Nested: SecondaryNestedType{Value: 10}},
				actual:   &NestedType{Nested: SecondaryNestedType{Value: 10}},
			},
			want: true,
		},
		{
			name: "slice equal",
			args: args{
				expected: &NestedType{ValueSlice: []int{10, 20}},
				actual:   &NestedType{ValueSlice: []int{10, 20}},
			},
			want: true,
		},
		{
			name: "slice additional in expected",
			args: args{
				expected: &NestedType{ValueSlice: []int{10, 20, 30}},
				actual:   &NestedType{ValueSlice: []int{10, 20}},
			},
			want: false,
		},
		{
			name: "slice additional in actual",
			args: args{
				expected: &NestedType{ValueSlice: []int{10, 20}},
				actual:   &NestedType{ValueSlice: []int{10, 20, 30}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PartiallyDeepEqual(tt.args.expected, tt.args.actual); got != tt.want {
				t.Errorf("PartiallyDeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
