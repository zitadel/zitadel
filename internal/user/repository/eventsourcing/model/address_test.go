package model

import "testing"

func TestAddressChanges(t *testing.T) {
	type args struct {
		existing *Address
		new      *Address
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all fields changed",
			args: args{
				existing: &Address{Country: "Country", Locality: "Locality", PostalCode: "PostalCode", Region: "Region", StreetAddress: "StreetAddress"},
				new:      &Address{Country: "CountryChanged", Locality: "LocalityChanged", PostalCode: "PostalCodeChanged", Region: "RegionChanged", StreetAddress: "StreetAddressChanged"},
			},
			res: res{
				changesLen: 5,
			},
		},
		{
			name: "no fields changed",
			args: args{
				existing: &Address{Country: "Country", Locality: "Locality", PostalCode: "PostalCode", Region: "Region", StreetAddress: "StreetAddress"},
				new:      &Address{Country: "Country", Locality: "Locality", PostalCode: "PostalCode", Region: "Region", StreetAddress: "StreetAddress"},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
