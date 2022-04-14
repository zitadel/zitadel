package checker_test

import (
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CheckerEmpty(t *testing.T) {
	checkerInstance := &checker.Checker{}
	assert.False(t, checkerInstance.CheckFailed())
	assert.Zero(t, checkerInstance.StepCount())
}

func Test_CheckerWithValueNotEmptyCheck(t *testing.T) {
	tests := []struct {
		name   string
		arg    string
		failed bool
	}{
		{
			"value with text",
			"test",
			false,
		},
		{
			"value empty",
			"",
			true,
		},
		{
			"value space",
			" ",
			false,
		},
		{
			"value long",
			"iwnakspemwebfuiasdvniaoenfkalnewrjaksduvfhypcvojamlkefnalewkrjbaiuhdsvhnypcxovnalkjwejfbnazfvghoiyxcpvjyomfklaenfqoefha<opdvnyxjkcvn aoifhjaöosdifnkjaewbfaisdvhypfovnalJDSJFBNAODVHANOFBVNaefpasdivnajksef aps9djcadslfknaeiofihvloyuixcvnljknaWEFDÖOAISVHNÖOawejkfhankl",
			false,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithValueNotEmptyCheck(
				"test",
				func() string { return tt.arg },
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithValueNotEmptyCheck got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithValuesNotEmptyCheck(t *testing.T) {
	tests := []struct {
		name   string
		arg    []string
		failed bool
	}{
		{
			"values single not empty",
			[]string{"test"},
			false,
		},
		{
			"values 3 not empty",
			[]string{"test1", "test2", "test3"},
			false,
		},
		{
			"values 5 not empty",
			[]string{"test1", "test2", "test3", "test4", "test5"},
			false,
		},
		{
			"values first empty",
			[]string{"", "test2", "test3", "test4", "test5"},
			true,
		},
		{
			"values last empty",
			[]string{"test1", "test2", "test3", "test4", ""},
			true,
		}, {
			"values mid empty",
			[]string{"test1", "test2", "", "test4", "test5"},
			true,
		},
		{
			"values all empty",
			[]string{"", "", "", "", ""},
			true,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithValuesNotEmptyCheck(
				func() []string { return tt.arg },
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithValuesNotEmptyCheck got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithValueLengthCheck(t *testing.T) {
	type args struct {
		value     string
		minlength int
		maxlength int
	}

	tests := []struct {
		name   string
		args   args
		failed bool
	}{
		{
			"value empty no fail",
			args{
				"", 0, 0,
			},
			false,
		},
		{
			"value empty fail",
			args{
				"", 1, 0,
			},
			true,
		},
		{
			"value len=4 min=0 max=4",
			args{
				"test", 0, 4,
			},
			false,
		},
		{
			"value len=4 min=4 max=10",
			args{
				"test", 4, 10,
			},
			false,
		},
		{
			"value len=4 min=4 max=10",
			args{
				"test", 4, 10,
			},
			false,
		},
		{
			"value len=4 min=0 max=10",
			args{
				"test", 0, 10,
			},
			false,
		},
		{
			"value len=5 min=0 max=4",
			args{
				"test1", 0, 4,
			},
			true,
		},
		{
			"value len=4 min=5 max=10",
			args{
				"test", 5, 10,
			},
			true,
		},
		{
			"value len=4 min=6 max=4",
			args{
				"test1", 6, 4,
			},
			true,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithValueLengthCheck(
				"test",
				func() string { return tt.args.value },
				tt.args.minlength,
				tt.args.maxlength,
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithValueLengthCheck got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithValueEqualsCheck(t *testing.T) {
	type args struct {
		value  string
		equals string
	}

	tests := []struct {
		name   string
		args   args
		failed bool
	}{
		{
			"value empty equals",
			args{
				"", "",
			},
			false,
		},
		{
			"value empty not equals",
			args{
				"", "1",
			},
			true,
		},
		{
			"value equals",
			args{
				"test", "test",
			},
			false,
		},
		{
			"value not equals",
			args{
				"test", "nononono",
			},
			true,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithValueEqualsCheck(
				"test",
				func() string { return tt.args.value },
				func() string { return tt.args.equals },
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithValueEqualsCheck got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithConditionalValueNotEmpty(t *testing.T) {
	type args struct {
		condition bool
		value     string
	}

	tests := []struct {
		name   string
		args   args
		failed bool
	}{
		{
			"value with text",
			args{
				condition: false,
				value:     "test",
			},
			false,
		},
		{
			"value with text",
			args{
				condition: true,
				value:     "test",
			},
			false,
		},
		{
			"value empty",
			args{
				condition: false,
				value:     "",
			},
			false,
		},
		{
			"value empty",
			args{
				condition: true,
				value:     "",
			},
			true,
		},
		{
			"value space",
			args{
				condition: false,
				value:     " ",
			},
			false,
		},
		{
			"value space",
			args{
				condition: true,
				value:     " ",
			},
			false,
		},
		{
			"value long",
			args{
				condition: true,
				value:     "iwnakspemwebfuiasdvniaoenfkalnewrjaksduvfhypcvojamlkefnalewkrjbaiuhdsvhnypcxovnalkjwejfbnazfvghoiyxcpvjyomfklaenfqoefha<opdvnyxjkcvn aoifhjaöosdifnkjaewbfaisdvhypfovnalJDSJFBNAODVHANOFBVNaefpasdivnajksef aps9djcadslfknaeiofihvloyuixcvnljknaWEFDÖOAISVHNÖOawejkfhankl",
			},
			false,
		},
		{
			"value long",
			args{
				condition: false,
				value:     "iwnakspemwebfuiasdvniaoenfkalnewrjaksduvfhypcvojamlkefnalewkrjbaiuhdsvhnypcxovnalkjwejfbnazfvghoiyxcpvjyomfklaenfqoefha<opdvnyxjkcvn aoifhjaöosdifnkjaewbfaisdvhypfovnalJDSJFBNAODVHANOFBVNaefpasdivnajksef aps9djcadslfknaeiofihvloyuixcvnljknaWEFDÖOAISVHNÖOawejkfhankl",
			},
			false,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithConditionalValueNotEmpty(
				func() bool { return tt.args.condition },
				"test",
				func() string { return tt.args.value },
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithConditionalValueNotEmpty got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithLogicStep(t *testing.T) {
	tests := []struct {
		name   string
		arg    func() error
		failed bool
	}{
		{
			"logic no error",
			func() error { return nil },
			false,
		},
		{
			"logic error",
			func() error { return fmt.Errorf("failed") },
			true,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithLogicStep(
				tt.arg,
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithLogicStep got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithConditionalLogicStep(t *testing.T) {
	type args struct {
		cond func() bool
		f    func() error
	}
	tests := []struct {
		name   string
		args   args
		failed bool
	}{
		{
			"logic no error condition true",
			args{
				cond: func() bool { return true },
				f:    func() error { return nil },
			},
			false,
		},
		{
			"logic no error condition false",
			args{
				cond: func() bool { return false },
				f:    func() error { return nil },
			},
			false,
		},
		{
			"logic error condition true",
			args{
				cond: func() bool { return true },
				f:    func() error { return fmt.Errorf("failed") },
			},
			true,
		},
		{
			"logic error condition false",
			args{
				cond: func() bool { return false },
				f:    func() error { return fmt.Errorf("failed") },
			},
			false,
		},
	}

	errorFunc := func() {}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkerInstance := &checker.Checker{}
			checkerInstance.WithConditionalLogicStep(
				tt.args.cond,
				tt.args.f,
				"test",
				errorFunc,
			)

			got := checkerInstance.CheckFailed()
			if got != tt.failed {
				t.Errorf("WithConditionalLogicStep got = %v, want = %v", got, tt.failed)
			}
		})
	}
}

func Test_CheckerWithValueStep(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		done := false
		checkerInstance := &checker.Checker{}
		checkerInstance.WithValueStep(
			func() { done = true },
		)

		got := checkerInstance.CheckFailed()
		if got != false {
			t.Errorf("WithValueStep got = %v, want = %v", got, false)
		}
		if !done {
			t.Errorf("WithValueStep not done")
		}
	})
}
