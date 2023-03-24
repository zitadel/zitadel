import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import * as i18nIsoCountries from 'i18n-iso-countries';
import { CountryCode, getCountries, getCountryCallingCode } from 'libphonenumber-js';

export interface CountryPhoneCode {
  countryCode: string;
  countryName: string;
  countryCallingCode: string;
}

@Injectable()
export class CountryCallingCodesService {
  constructor(private translateService: TranslateService) {}

  public getCountryCallingCodes(): CountryPhoneCode[] {
    const currentLang = this.translateService.currentLang ?? 'en';
    const countryPhoneCodes = getCountries()
      .filter((code: CountryCode) => i18nIsoCountries.getName(code.toString(), currentLang))
      .map((code: CountryCode) => {
        return <CountryPhoneCode>{
          countryCode: code,
          countryName: i18nIsoCountries.getName(code.toString(), currentLang),
          countryCallingCode: getCountryCallingCode(code),
        };
      })
      .sort((a, b) => a.countryName.localeCompare(b.countryName));
    return countryPhoneCodes;
  }
}
