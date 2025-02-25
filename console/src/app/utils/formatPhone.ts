import { CountryCode, parsePhoneNumber } from 'libphonenumber-js';

export function formatPhone(phone?: string): { phone: string; country: CountryCode } | null {
  const defaultCountry = 'US';

  if (phone) {
    try {
      const phoneNumber = parsePhoneNumber(phone, defaultCountry);
      const country = phoneNumber.country ?? defaultCountry;
      if (phoneNumber) {
        return { phone: phoneNumber.formatInternational(), country };
      }
    } catch (e) {
      return null;
    }
  }

  return null;
}
