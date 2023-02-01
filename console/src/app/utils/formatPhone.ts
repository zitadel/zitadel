import { CountryCode, parsePhoneNumber } from 'libphonenumber-js';

export function formatPhone(phone: string): { phone: string; country: CountryCode } {
  const defaultCountry = 'CH';

  if (phone) {
    try {
      const phoneNumber = parsePhoneNumber(phone, defaultCountry);
      const country = phoneNumber.country ?? defaultCountry;
      if (phoneNumber) {
        return { phone: phoneNumber.formatInternational(), country };
      }
    } catch (error) {
      console.error(error);
    }
  }

  return { phone, country: defaultCountry };
}
