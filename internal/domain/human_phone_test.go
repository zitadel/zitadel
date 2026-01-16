package domain

import (
	"sort"
	"strings"
	"testing"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestFormatPhoneNumber(t *testing.T) {
	type args struct {
		phone *Phone
	}
	tests := []struct {
		name    string
		args    args
		result  *Phone
		errFunc func(err error) bool
	}{
		{
			name: "invalid phone number",
			args: args{
				phone: &Phone{
					PhoneNumber: "PhoneNumber",
				},
			},
			errFunc: zerrors.IsErrorInvalidArgument,
		},
		{
			name: "format phone 071...",
			args: args{
				phone: &Phone{
					PhoneNumber: "0711234567",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 0041...",
			args: args{
				phone: &Phone{
					PhoneNumber: "0041711234567",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 071 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "071 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone +4171 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "+4171 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format phone 004171 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "004171 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+41711234567",
			},
		},
		{
			name: "format non swiss phone 004371 xxx xx xx",
			args: args{
				phone: &Phone{
					PhoneNumber: "004371 123 45 67",
				},
			},
			result: &Phone{
				PhoneNumber: "+43711234567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized, err := tt.args.phone.PhoneNumber.Normalize()
			if tt.errFunc == nil && tt.result.PhoneNumber != normalized {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result.PhoneNumber, normalized)
			}
			if tt.errFunc != nil && !tt.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestGetCountryCodes(t *testing.T) {
	codes := GetCountryCodes()

	t.Run("returns non-empty list", func(t *testing.T) {
		if len(codes) == 0 {
			t.Error("GetCountryCodes() returned empty list")
		}
	})

	t.Run("returns more than 200 countries", func(t *testing.T) {
		// libphonenumber supports 240+ regions
		if len(codes) < 200 {
			t.Errorf("GetCountryCodes() returned %d countries, expected at least 200", len(codes))
		}
	})

	t.Run("list is sorted by name", func(t *testing.T) {
		sorted := sort.SliceIsSorted(codes, func(i, j int) bool {
			return codes[i].Name < codes[j].Name
		})
		if !sorted {
			t.Error("GetCountryCodes() did not return a list sorted by name")
		}
	})

	t.Run("all entries have valid format", func(t *testing.T) {
		for _, code := range codes {
			if code.ISOCode == "" {
				t.Errorf("country code has empty ISOCode: %+v", code)
			}
			if len(code.ISOCode) < 2 {
				t.Errorf("country code has invalid ISOCode length: %+v", code)
			}
			if !strings.HasPrefix(code.CallingCode, "+") {
				t.Errorf("country code CallingCode does not start with '+': %+v", code)
			}
			if code.Name == "" {
				t.Errorf("country code has empty Name: %+v", code)
			}
		}
	})
}

func TestGetCountryCodes_ContainsExpectedCountries(t *testing.T) {
	codes := GetCountryCodes()

	// Build a map for easier lookup
	codeMap := make(map[string]CountryCode)
	for _, code := range codes {
		codeMap[code.ISOCode] = code
	}

	tests := []struct {
		name        string
		isoCode     string
		callingCode string
		countryName string
	}{
		{
			name:        "United States",
			isoCode:     "US",
			callingCode: "+1",
			countryName: "United States",
		},
		{
			name:        "Switzerland",
			isoCode:     "CH",
			callingCode: "+41",
			countryName: "Switzerland",
		},
		{
			name:        "United Kingdom",
			isoCode:     "GB",
			callingCode: "+44",
			countryName: "United Kingdom",
		},
		{
			name:        "Germany",
			isoCode:     "DE",
			callingCode: "+49",
			countryName: "Germany",
		},
		{
			name:        "Australia",
			isoCode:     "AU",
			callingCode: "+61",
			countryName: "Australia",
		},
		{
			name:        "Japan",
			isoCode:     "JP",
			callingCode: "+81",
			countryName: "Japan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, found := codeMap[tt.isoCode]
			if !found {
				t.Errorf("GetCountryCodes() does not contain %s (%s)", tt.name, tt.isoCode)
				return
			}
			if code.CallingCode != tt.callingCode {
				t.Errorf("GetCountryCodes() %s CallingCode = %s, want %s", tt.isoCode, code.CallingCode, tt.callingCode)
			}
			if code.Name != tt.countryName {
				t.Errorf("GetCountryCodes() %s Name = %s, want %s", tt.isoCode, code.Name, tt.countryName)
			}
		})
	}
}
