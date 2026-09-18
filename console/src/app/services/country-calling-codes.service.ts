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
    const countryName = this.countryNameResolver(currentLang);
    const countryPhoneCodes = getCountries()
      .map((code: CountryCode) => ({ code, name: countryName(code) }))
      .filter((country): country is { code: CountryCode; name: string } => !!country.name)
      .map(({ code, name }) => {
        return <CountryPhoneCode>{
          countryCode: code,
          countryName: name,
          countryCallingCode: getCountryCallingCode(code),
        };
      })
      .sort((a, b) => a.countryName.localeCompare(b.countryName));
    return countryPhoneCodes;
  }

  /**
   * i18n-iso-countries only ships country names for a subset of the languages we
   * support, for example it has no Traditional Chinese data. Fall back to the
   * browsers own CLDR data in that case so the list stays localized instead of
   * ending up empty.
   */
  private countryNameResolver(lang: string): (code: CountryCode) => string | undefined {
    if (i18nIsoCountries.getName('CH', lang)) {
      return (code) => i18nIsoCountries.getName(code.toString(), lang);
    }
    try {
      const regionNames = new Intl.DisplayNames([lang], { type: 'region' });
      return (code) => regionNames.of(code.toString());
    } catch {
      return (code) => i18nIsoCountries.getName(code.toString(), 'en');
    }
  }
}
