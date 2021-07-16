import { PolicyComponentType } from '../policies/policy-component-types.enum';

export interface GridPolicy {
  i18nTitle: string;
  i18nDesc: string;
  iamRouterLink: any;
  orgRouterLink: any;
  iamWithRole: string[];
  orgWithRole: string[];
  tags: string[];
  icon?: string;
  svgIcon?: string;
  color: string;
}

export const COMPLEXITY_POLICY: GridPolicy = {
  i18nTitle: 'POLICY.PWD_COMPLEXITY.TITLE',
  i18nDesc: 'POLICY.PWD_COMPLEXITY.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.COMPLEXITY],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.COMPLEXITY],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['login', 'security'],
  svgIcon: 'mdi_textbox_password',
  color: 'yellow',
};

export const IAM_POLICY = {
  i18nTitle: 'POLICY.IAM_POLICY.TITLE',
  i18nDesc: 'POLICY.IAM_POLICY.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.IAM],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.IAM],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['iam.policy.read'],
  tags: ['login'],
  icon: 'las la-gem',
  color: 'purple',
};

export const LOGIN_POLICY = {
  i18nTitle: 'POLICY.LOGIN_POLICY.TITLE',
  i18nDesc: 'POLICY.LOGIN_POLICY.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.LOGIN],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.LOGIN],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['login', 'security'],
  icon: 'las la-sign-in-alt',
  color: 'green',
};

export const PRIVATELABEL_POLICY = {
  i18nTitle: 'POLICY.PRIVATELABELING.TITLE',
  i18nDesc: 'POLICY.PRIVATELABELING.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.PRIVATELABEL],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.PRIVATELABEL],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['login', 'appearance'],
  icon: 'las la-sign-in-alt',
  color: 'blue',
};

export const PRIVACY_POLICY = {
  i18nTitle: 'POLICY.PRIVACY_POLICY.TITLE',
  i18nDesc: 'POLICY.PRIVACY_POLICY.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.PRIVACYPOLICY],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.PRIVACYPOLICY],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['documents', 'text'],
  icon: 'las la-file-contract',
  color: 'black',
};

export const MESSAGE_TEXTS_POLICY = {
  i18nTitle: 'POLICY.MESSAGE_TEXTS.TITLE',
  i18nDesc: 'POLICY.MESSAGE_TEXTS.DESCRIPTION',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.MESSAGETEXTS],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.MESSAGETEXTS],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['appearance', 'text'],
  icon: 'las la-paragraph',
  color: 'red',
};

export const LOGIN_TEXTS_POLICY = {
  i18nTitle: 'POLICY.LOGIN_TEXTS.TITLE',
  i18nDesc: 'POLICY.LOGIN_TEXTS.DESCRIPTION_SHORT',
  iamRouterLink: ['/iam', 'policy', PolicyComponentType.LOGINTEXTS],
  orgRouterLink: ['/org', 'policy', PolicyComponentType.LOGINTEXTS],
  iamWithRole: ['iam.policy.read'],
  orgWithRole: ['policy.read'],
  tags: ['appearance', 'text'],
  icon: 'las la-paragraph',
  color: 'red',
};

export const POLICIES: GridPolicy[] = [
  COMPLEXITY_POLICY,
  IAM_POLICY,
  LOGIN_POLICY,
  PRIVATELABEL_POLICY,
  PRIVACY_POLICY,
  MESSAGE_TEXTS_POLICY,
  LOGIN_TEXTS_POLICY,
];
