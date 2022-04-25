package saml

import (
	"testing"
	"time"
)

func TestTime_checkIfRequestTimeIsStillValid(t *testing.T) {
	type args struct {
		notBefore    string
		notOnOrAfter string
	}
	now := time.Now().UTC()

	tests := []struct {
		name string
		args args
		res  bool
	}{
		{
			"check ok 1",
			args{
				notBefore:    now.Add(-1 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: now.Add(1 * time.Minute).Format(defaultTimeLayout),
			},
			false,
		},
		{
			"check ok 2",
			args{
				notBefore:    now.Add(-1 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: now.Add(5 * time.Minute).Format(defaultTimeLayout),
			},
			false,
		},
		{
			"check ok 3",
			args{
				notBefore:    now.Add(-5 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: now.Add(5 * time.Minute).Format(defaultTimeLayout),
			},
			false,
		},
		{
			"check not ok 1",
			args{
				notBefore:    now.Add(1 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: now.Add(5 * time.Minute).Format(defaultTimeLayout),
			},
			true,
		},
		{
			"check not ok 2",
			args{
				notBefore:    now.Add(-5 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: now.Add(-1 * time.Minute).Format(defaultTimeLayout),
			},
			true,
		},
		{
			"check ok no times",
			args{
				notBefore:    "",
				notOnOrAfter: "",
			},
			false,
		},
		{
			"check ok only notOnOrAfter",
			args{
				notBefore:    "",
				notOnOrAfter: now.Add(1 * time.Minute).Format(defaultTimeLayout),
			},
			false,
		},
		{
			"check not ok only notOnOrAfter",
			args{
				notBefore:    "",
				notOnOrAfter: now.Add(-1 * time.Minute).Format(defaultTimeLayout),
			},
			true,
		},
		{
			"check not ok only notBefore",
			args{
				notBefore:    now.Add(1 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: "",
			},
			true,
		},
		{
			"check ok only notBefore",
			args{
				notBefore:    now.Add(-1 * time.Minute).Format(defaultTimeLayout),
				notOnOrAfter: "",
			},
			false,
		},
		{
			"check cant parse notBefore",
			args{
				notBefore:    "what time is it?",
				notOnOrAfter: "",
			},
			true,
		},
		{
			"check cant parse notOnOrAfter",
			args{
				notBefore:    "",
				notOnOrAfter: "what time is it?",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notBeforeF := func() string {
				return tt.args.notBefore
			}
			notOnOrAfterF := func() string {
				return tt.args.notOnOrAfter
			}

			errF := checkIfRequestTimeIsStillValid(notBeforeF, notOnOrAfterF)
			err := errF()
			if (err != nil) != tt.res {
				t.Errorf("ParseCertificates() got = %v, want %v", err != nil, tt.res)
			}
		})
	}
}
