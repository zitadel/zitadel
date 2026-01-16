package domain

import (
	"fmt"
	"sort"
	"time"

	"github.com/ttacon/libphonenumber"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const defaultRegion = "CH"

type PhoneNumber string

func (p PhoneNumber) Normalize() (PhoneNumber, error) {
	return p.NormalizeWithRegion(defaultRegion)
}

func (p PhoneNumber) NormalizeWithRegion(region string) (PhoneNumber, error) {
	if p == "" {
		return p, zerrors.ThrowInvalidArgument(nil, "PHONE-Zt0NV", "Errors.User.Phone.Empty")
	}
	if region == "" {
		region = defaultRegion
	}
	phoneNr, err := libphonenumber.Parse(string(p), region)
	if err != nil {
		return p, zerrors.ThrowInvalidArgument(err, "PHONE-so0wa", "Errors.User.Phone.Invalid")
	}
	return PhoneNumber(libphonenumber.Format(phoneNr, libphonenumber.E164)), nil
}

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     PhoneNumber
	IsPhoneVerified bool
	// PlainCode is set by the command and can be used to return it to the caller (API)
	PlainCode *string
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue
	Expiry time.Duration
}

func (p *Phone) Normalize() error {
	if p == nil {
		return zerrors.ThrowInvalidArgument(nil, "PHONE-YlbwO", "Errors.User.Phone.Empty")
	}
	normalizedNumber, err := p.PhoneNumber.Normalize()
	if err != nil {
		return err
	}
	// Issue for avoiding mutating state: https://github.com/zitadel/zitadel/issues/5412
	p.PhoneNumber = normalizedNumber
	return nil
}

type PhoneState int32

const (
	PhoneStateUnspecified PhoneState = iota
	PhoneStateActive
	PhoneStateRemoved

	phoneStateCount
)

func (s PhoneState) Valid() bool {
	return s >= 0 && s < phoneStateCount
}

func (s PhoneState) Exists() bool {
	return s == PhoneStateActive
}

// CountryCode represents a country with its ISO code and calling code
type CountryCode struct {
	ISOCode     string
	CallingCode string
	Name        string
}

// GetCountryCodes returns a list of country codes for phone number input using libphonenumber
func GetCountryCodes() []CountryCode {
	regions := libphonenumber.GetSupportedRegions()
	namer := display.English.Regions()

	codes := make([]CountryCode, 0, len(regions))
	for region := range regions {
		countryCode := libphonenumber.GetCountryCodeForRegion(region)
		if countryCode == 0 {
			continue
		}

		// Get the display name for the region
		regionTag, err := language.ParseRegion(region)
		name := region
		if err == nil {
			name = namer.Name(regionTag)
		}

		codes = append(codes, CountryCode{
			ISOCode:     region,
			CallingCode: fmt.Sprintf("+%d", countryCode),
			Name:        name,
		})
	}

	// Sort by country name
	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Name < codes[j].Name
	})

	return codes
}
