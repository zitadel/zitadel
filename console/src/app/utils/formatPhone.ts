import { CountryCode, parsePhoneNumberWithError } from 'libphonenumber-js';

export function formatPhone(phone?: string): { phone: string; country: CountryCode } | null {
  const defaultCountry = 'US';

  if (!phone) {
    return null;
  }

  try {
    const phoneNumber = parsePhoneNumberWithError(phone, defaultCountry);
    const country = phoneNumber.country ?? defaultCountry;
    return { phone: phoneNumber.formatInternational(), country };
  } catch {
    return null;
  }
}
