export interface SettingLinks {
  i18nTitle: string;
  i18nDesc: string;
  iamRouterLink: any;
  orgRouterLink: any;
  iamWithRole: string[];
  orgWithRole: string[];
  icon?: string;
  svgIcon?: string;
  color: string;
}

export const LOGIN_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.LOGIN',
  i18nDesc: '',
  iamRouterLink: ['/settings?id=login'],
  orgRouterLink: ['/org-settings?id=login'],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-sign-in-alt',
  color: 'green',
};

export const APPEARANCE_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.APPEARANCE',
  i18nDesc: '',
  iamRouterLink: ['/settings?id=branding'],
  orgRouterLink: ['/org-settings?id=branding'],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-swatchbook',
  color: 'blue',
};

export const PRIVACY_POLICY: SettingLinks = {
  i18nTitle: 'SETTINGS.LIST.PRIVACYPOLICY',
  i18nDesc: '',
  iamRouterLink: ['/settings?id=privacypolicy'],
  orgRouterLink: ['/org-settings?id=privacypolicy'],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-file-contract',
  color: 'black',
};

export const NOTIFICATION_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.NOTIFICATIONS',
  i18nDesc: '',
  iamRouterLink: ['/settings?id=notifications'],
  orgRouterLink: ['/org-settings?id=notifications'],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-bell',
  color: 'red',
};

export const SETTINGLINKS: SettingLinks[] = [LOGIN_GROUP, APPEARANCE_GROUP, PRIVACY_POLICY, NOTIFICATION_GROUP];
