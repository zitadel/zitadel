import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

export const accountLanguages = ['pt', 'en', 'es'] as const;
export type AccountLanguage = (typeof accountLanguages)[number];
export function accountLocale(value: string | undefined): AccountLanguage {
  const base = value?.split('-')[0].toLowerCase();
  return accountLanguages.includes(base as AccountLanguage) ? (base as AccountLanguage) : 'en';
}
@Injectable({ providedIn: 'root' })
export class AccountPreferences {
  readonly languages = accountLanguages;
  constructor(private translate: TranslateService) {}
  initialize(): void {
    this.language(localStorage.getItem('account:language') ?? this.translate.getCurrentLang() ?? navigator.language);
  }
  language(value: string): void {
    const locale = accountLocale(value);
    this.translate.use(locale);
    document.documentElement.lang = locale === 'pt' ? 'pt-BR' : locale;
    localStorage.setItem('account:language', locale);
  }
  languageName(value: string): string {
    if (!value || value === 'und') return this.translate.instant('ACCOUNT.AUTO_LANGUAGE');
    try {
      const locale = accountLocale(this.translate.getCurrentLang());
      const name = new Intl.DisplayNames([locale], { type: 'language' }).of(value);
      return name ? name[0].toLocaleUpperCase(locale) + name.slice(1) : value;
    } catch {
      return value;
    }
  }
}
