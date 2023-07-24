package domain

import (
	"testing"
)

func TestNewIAMDomainName(t *testing.T) {
	type args struct {
		orgName   string
		iamDomain string
	}
	tests := []struct {
		name   string
		args   args
		result string
	}{
		{
			name: "Single word domain is already valid",
			args: args{
				orgName:   "single-word-domain",
				iamDomain: "localhost",
			},
			result: "single-word-domain.localhost",
		},
		{
			name: "resulting domain should be in lowercase",
			args: args{
				orgName:   "Uppercase org Name",
				iamDomain: "localhost",
			},
			result: "uppercase-org-name.localhost",
		},
		{
			name: "replace spaces with hyphens",
			args: args{
				orgName:   "my org name",
				iamDomain: "localhost",
			},
			result: "my-org-name.localhost",
		},
		{
			name: "replace invalid characters [^a-zA-Z0-9-] with empty spaces",
			args: args{
				orgName:   "mí Örg name?",
				iamDomain: "localhost",
			},
			result: "m-rg-name.localhost",
		},
		{
			name: "label created from org name size is not greater than 63 chars",
			args: args{
				orgName:   "my organization name must not exceed sixty-three characters 1234",
				iamDomain: "localhost",
			},
			result: "my-organization-name-must-not-exceed-sixty-three-characters-123.localhost",
		},
		{
			name: "resulting domain cannot exceed 253 chars",
			args: args{
				orgName:   "Lorem ipsum dolor sit amet",
				iamDomain: "llgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.co.uk",
			},
			result: "lorem-ipsum-dolor-sit.llgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch.co.uk",
		},
		{
			name: "label based on org name should not end with a hyphen",
			args: args{
				orgName:   "my super long organization name with many many many characters ",
				iamDomain: "localhost",
			},
			result: "my-super-long-organization-name-with-many-many-many-characters.localhost",
		},
		{
			name: "label based on org name should not start with a hyphen",
			args: args{
				orgName:   " my super long organization name with many many many characters",
				iamDomain: "localhost",
			},
			result: "my-super-long-organization-name-with-many-many-many-characters.localhost",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain := NewIAMDomainName(tt.args.orgName, tt.args.iamDomain)
			if tt.result != domain {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, domain)
			}
		})
	}
}
