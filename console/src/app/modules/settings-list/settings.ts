import { PolicyComponentServiceType } from '../policies/policy-component-types.enum';
import { SidenavSetting } from '../sidenav/sidenav.component';

export const ORGANIZATIONS: SidenavSetting = {
  id: 'organizations',
  i18nKey: 'SETTINGS.LIST.ORGS',
  groupI18nKey: 'SETTINGS.GROUPS.GENERAL',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.read'],
  },
};

export const FEATURESETTINGS: SidenavSetting = {
  id: 'features',
  i18nKey: 'SETTINGS.LIST.FEATURESETTINGS',
  groupI18nKey: 'SETTINGS.GROUPS.GENERAL',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.restrictions.read'],
  },
};

export const LANGUAGES: SidenavSetting = {
  id: 'languages',
  i18nKey: 'SETTINGS.LIST.LANGUAGES',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const OIDC: SidenavSetting = {
  id: 'oidc',
  i18nKey: 'SETTINGS.LIST.OIDC',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const WEBKEYS: SidenavSetting = {
  id: 'webkeys',
  i18nKey: 'SETTINGS.LIST.WEB_KEYS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const SECRETS: SidenavSetting = {
  id: 'secrets',
  i18nKey: 'SETTINGS.LIST.SECRETS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const SECURITY: SidenavSetting = {
  id: 'security',
  i18nKey: 'SETTINGS.LIST.SECURITY',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const VIEWS: SidenavSetting = {
  id: 'views',
  i18nKey: 'SETTINGS.LIST.VIEWS',
  groupI18nKey: 'SETTINGS.GROUPS.STORAGE',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.read'],
  },
};

export const FAILEDEVENTS: SidenavSetting = {
  id: 'failedevents',
  i18nKey: 'SETTINGS.LIST.FAILEDEVENTS',
  groupI18nKey: 'SETTINGS.GROUPS.STORAGE',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.read'],
  },
};

export const EVENTS: SidenavSetting = {
  id: 'events',
  i18nKey: 'SETTINGS.LIST.EVENTS',
  groupI18nKey: 'SETTINGS.GROUPS.STORAGE',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['events.read'],
  },
};

export const LOGIN: SidenavSetting = {
  id: 'login',
  i18nKey: 'SETTINGS.LIST.LOGIN',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const VERIFIED_DOMAINS: SidenavSetting = {
  id: 'verified_domains',
  i18nKey: 'SETTINGS.LIST.VERIFIED_DOMAINS',
  groupI18nKey: 'SETTINGS.GROUPS.DOMAIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['org.read'],
  },
};

export const DOMAIN: SidenavSetting = {
  id: 'domain',
  i18nKey: 'SETTINGS.LIST.DOMAIN',
  groupI18nKey: 'SETTINGS.GROUPS.DOMAIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['iam.policy.write'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.write'],
  },
};

export const LOCKOUT: SidenavSetting = {
  id: 'lockout',
  i18nKey: 'SETTINGS.LIST.LOCKOUT',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const AGE: SidenavSetting = {
  id: 'age',
  i18nKey: 'SETTINGS.LIST.AGE',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const COMPLEXITY: SidenavSetting = {
  id: 'complexity',
  i18nKey: 'SETTINGS.LIST.COMPLEXITY',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const IDP: SidenavSetting = {
  id: 'idp',
  i18nKey: 'SETTINGS.LIST.IDP',
  groupI18nKey: 'SETTINGS.GROUPS.LOGIN',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read', 'org.idp.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read', 'org.idp.read'],
  },
};

export const NOTIFICATIONS: SidenavSetting = {
  id: 'notifications',
  i18nKey: 'SETTINGS.LIST.NOTIFICATIONS',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
  },
};

export const SMTP_PROVIDER: SidenavSetting = {
  id: 'smtpprovider',
  i18nKey: 'SETTINGS.LIST.SMTP_PROVIDER',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const SMS_PROVIDER: SidenavSetting = {
  id: 'smsprovider',
  i18nKey: 'SETTINGS.LIST.SMS_PROVIDER',
  groupI18nKey: 'SETTINGS.GROUPS.NOTIFICATIONS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const MESSAGETEXTS: SidenavSetting = {
  id: 'messagetexts',
  i18nKey: 'SETTINGS.LIST.MESSAGETEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const LOGINTEXTS: SidenavSetting = {
  id: 'logintexts',
  i18nKey: 'SETTINGS.LIST.LOGINTEXTS',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const PRIVACYPOLICY: SidenavSetting = {
  id: 'privacypolicy',
  i18nKey: 'DESCRIPTIONS.SETTINGS.PRIVACY_POLICY.TITLE',
  groupI18nKey: 'SETTINGS.GROUPS.OTHER',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const BRANDING: SidenavSetting = {
  id: 'branding',
  i18nKey: 'SETTINGS.LIST.BRANDING',
  groupI18nKey: 'SETTINGS.GROUPS.APPEARANCE',
  requiredRoles: {
    [PolicyComponentServiceType.MGMT]: ['policy.read'],
    [PolicyComponentServiceType.ADMIN]: ['iam.policy.read'],
  },
};

export const ACTIONS: SidenavSetting = {
  id: 'actions',
  i18nKey: 'SETTINGS.LIST.ACTIONS',
  groupI18nKey: 'SETTINGS.GROUPS.ACTIONS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['action.execution.write', 'action.target.write'],
  },
  beta: true,
};

export const ACTIONS_TARGETS: SidenavSetting = {
  id: 'actions_targets',
  i18nKey: 'SETTINGS.LIST.TARGETS',
  groupI18nKey: 'SETTINGS.GROUPS.ACTIONS',
  requiredRoles: {
    [PolicyComponentServiceType.ADMIN]: ['action.execution.write', 'action.target.write'],
  },
  beta: true,
};
