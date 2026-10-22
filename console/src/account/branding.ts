import { ThemeService } from 'src/app/services/theme.service';
import { Injectable } from '@angular/core';
import { GrpcAuthService } from 'src/app/services/grpc-auth.service';
import { LabelPolicy } from 'src/app/proto/generated/zitadel/policy_pb';
import { AccountAuthService } from './auth.service';
import { OAuthService } from 'angular-oauth2-oidc';

@Injectable({ providedIn: 'root' })
export class AccountBranding {
  policy?: LabelPolicy.AsObject;
  unavailable = false;
  dark = matchMedia('(prefers-color-scheme: dark)').matches;
  name = '';
  logo = '';
  constructor(
    private native: GrpcAuthService,
    private auth: AccountAuthService,
    private oauth: OAuthService,
    private nativeTheme: ThemeService,
  ) {
    this.apply();
  }
  async load(): Promise<void> {
    try {
      this.policy = await this.native.getMyLabelPolicy(this.auth.environment.loginOrgId);
      this.name = String(this.oauth.getIdentityClaims()?.['urn:zitadel:iam:user:resourceowner:name'] ?? '');
      this.apply();
    } catch {
      this.unavailable = true;
    }
  }
  toggle(): void {
    this.dark = !this.dark;
    this.apply();
  }
  private asset(value: string | undefined): string {
    if (!value) return '';
    const url = new URL(value, this.auth.environment.api);
    return url.origin === new URL(this.auth.environment.api).origin && url.protocol === 'https:' ? url.href : '';
  }
  private apply(): void {
    const p = this.policy;
    document.documentElement.classList.toggle('dark-theme', this.dark);
    document.documentElement.classList.toggle('light-theme', !this.dark);
    document.documentElement.style.colorScheme = this.dark ? 'dark' : 'light';
    this.nativeTheme.setDarkTheme(this.dark);
    if (!p) return;
    const colors = {
      '--account-brand': this.dark ? p.primaryColorDark : p.primaryColor,
      '--account-background': this.dark ? p.backgroundColorDark : p.backgroundColor,
      '--account-text': this.dark ? p.fontColorDark : p.fontColor,
      '--account-warn': this.dark ? p.warnColorDark : p.warnColor,
    };
    for (const [key, value] of Object.entries(colors)) {
      document.documentElement.style.removeProperty(key);
      if (value && CSS.supports('color', value)) document.documentElement.style.setProperty(key, value);
    }
    this.logo = this.asset(this.dark ? p.logoUrlDark || p.logoUrl : p.logoUrl || p.logoUrlDark);
    const icon = this.asset(this.dark ? p.iconUrlDark || p.iconUrl : p.iconUrl || p.iconUrlDark);
    document.querySelector<HTMLLinkElement>('#account-icon')?.setAttribute('href', icon || 'data:,');
    const font = this.asset(p.fontUrl);
    if (font) {
      const face = new FontFace('AccountOrganization', `url(${JSON.stringify(font)})`);
      face
        .load()
        .then((loaded) => {
          document.fonts.add(loaded);
          document.documentElement.style.setProperty('--account-font', 'AccountOrganization, sans-serif');
        })
        .catch(() => {
          this.unavailable = true;
        });
    }
    document.title = this.name || 'Account';
  }
}
