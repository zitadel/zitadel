import { PolicyComponentType } from '../policies/policy-component-types.enum';

export const IAM_COMPLEXITY_LINK = {
  i18nTitle: 'POLICY.PWD_COMPLEXITY.TITLE',
  i18nDesc: 'POLICY.PWD_COMPLEXITY.DESCRIPTION',
  routerLink: ['/iam', 'policy', PolicyComponentType.COMPLEXITY],
  withRole: ['iam.policy.read'],
};

export const IAM_POLICY_LINK = {
  i18nTitle: 'POLICY.IAM_POLICY.TITLE',
  i18nDesc: 'POLICY.IAM_POLICY.DESCRIPTION',
  routerLink: ['/iam', 'policy', PolicyComponentType.IAM],
  withRole: ['iam.policy.read'],
};

export const IAM_LOGIN_POLICY_LINK = {
  i18nTitle: 'POLICY.LOGIN_POLICY.TITLE',
  i18nDesc: 'POLICY.LOGIN_POLICY.DESCRIPTION',
  routerLink: ['/iam', 'policy', PolicyComponentType.LOGIN],
  withRole: ['iam.policy.read'],
};

export const IAM_PRIVATELABEL_LINK = {
  i18nTitle: 'POLICY.PRIVATELABELING.TITLE',
  i18nDesc: 'POLICY.PRIVATELABELING.DESCRIPTION',
  routerLink: ['/iam', 'policy', PolicyComponentType.PRIVATELABEL],
  withRole: ['iam.policy.read'],
};

export const ORG_COMPLEXITY_LINK = {
  i18nTitle: 'POLICY.PWD_COMPLEXITY.TITLE',
  i18nDesc: 'POLICY.PWD_COMPLEXITY.DESCRIPTION',
  routerLink: ['/org', 'policy', PolicyComponentType.COMPLEXITY],
  withRole: ['policy.read'],
};

export const ORG_IAM_POLICY_LINK = {
  i18nTitle: 'POLICY.IAM_POLICY.TITLE',
  i18nDesc: 'POLICY.IAM_POLICY.DESCRIPTION',
  routerLink: ['/org', 'policy', PolicyComponentType.IAM],
  withRole: ['iam.policy.read'],
};

export const ORG_LOGIN_POLICY_LINK = {
  i18nTitle: 'POLICY.LOGIN_POLICY.TITLE',
  i18nDesc: 'POLICY.LOGIN_POLICY.DESCRIPTION',
  routerLink: ['/org', 'policy', PolicyComponentType.LOGIN],
  withRole: ['policy.read'],
};


export const ORG_PRIVATELABEL_LINK = {
  i18nTitle: 'POLICY.PRIVATELABELING.TITLE',
  i18nDesc: 'POLICY.PRIVATELABELING.DESCRIPTION',
  routerLink: ['/org', 'policy', PolicyComponentType.PRIVATELABEL],
  withRole: ['policy.read'],
};
