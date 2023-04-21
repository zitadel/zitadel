package crdb

import "testing"

func Test_defaultValue(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				value: "asdf",
			},
			want: "'asdf'",
		},
		{
			name: "primitive non string",
			args: args{
				value: 1,
			},
			want: "1",
		},
		{
			name: "stringer",
			args: args{
				value: testStringer(0),
			},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultValue(tt.args.value); got != tt.want {
				t.Errorf("defaultValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testStringer int

func (t testStringer) String() string {
	return "0529958243"
}
