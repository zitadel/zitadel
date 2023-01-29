import { CountryCode, parsePhoneNumber } from 'libphonenumber-js';

export function formatPhone(phone: string): { phone: string; country: CountryCode } {
  // Format phone before save (add +)
  try {
    const phoneNumber = parsePhoneNumber(phone ?? '', 'CH');
    const country = phoneNumber.country ?? 'CH';
    if (phoneNumber) {
      return { phone: phoneNumber.formatInternational(), country };
    }
  } catch (error) {
    console.error(error);
  }
  return { phone, country: 'CH' };
}
