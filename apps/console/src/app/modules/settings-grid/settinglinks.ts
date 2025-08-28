export interface SettingLinks {
  i18nTitle: string;
  i18nDesc: string;
  iamRouterLink: any;
  orgRouterLink?: any;
  queryParams: any;
  iamWithRole?: string[];
  orgWithRole?: string[];
  icon?: string;
  svgIcon?: string;
  color: string;
}

export const LOGIN_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.LOGIN',
  i18nDesc: 'POLICY.LOGIN_POLICY.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'login' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-sign-in-alt',
  color: 'green',
};

export const APPEARANCE_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.APPEARANCE',
  i18nDesc: 'POLICY.PRIVATELABELING.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'branding' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-swatchbook',
  color: 'blue',
};

export const PRIVACY_POLICY: SettingLinks = {
  i18nTitle: 'DESCRIPTIONS.SETTINGS.PRIVACY_POLICY.TITLE',
  i18nDesc: 'POLICY.PRIVACY_POLICY.DESCRIPTION',
  iamRouterLink: ['/settings'],
  orgRouterLink: ['/org-settings'],
  queryParams: { id: 'privacypolicy' },
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  icon: 'las la-file-contract',
  color: 'black',
};

export const NOTIFICATION_GROUP: SettingLinks = {
  i18nTitle: 'SETTINGS.GROUPS.NOTIFICATIONS',
  i18nDesc: 'SETTINGS.LIST.NOTIFICATIONS_DESC',
  iamRouterLink: ['/settings'],
  queryParams: { id: 'smtpprovider' },
  iamWithRole: ['iam.policy.read'],
  icon: 'las la-bell',
  color: 'red',
};

export const SETTINGLINKS: SettingLinks[] = [LOGIN_GROUP, APPEARANCE_GROUP, PRIVACY_POLICY, NOTIFICATION_GROUP];
