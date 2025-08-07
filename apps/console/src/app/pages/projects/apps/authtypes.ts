import { OIDCAppType } from 'src/app/proto/generated/zitadel/app_pb';

// export enum AppType {
//     "WEB",
//     "USER_AGENT",
//     "NATIVE",
//     "API"
// }

export enum AppCreateType {
  API = 'API',
  OIDC = 'OIDC',
  SAML = 'SAML',
}

export interface RadioItemAppType {
  // key: string;
  createType: AppCreateType;
  oidcAppType?: OIDCAppType;
  titleI18nKey: string;
  descI18nKey: string;
  prefix: string;
  background: string;
  protocol: 'OIDC' | 'SAML';
}

export const WEB_TYPE: RadioItemAppType = {
  // key: AppType.WEB,
  titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.TITLE',
  descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.DESCRIPTION',
  createType: AppCreateType.OIDC,
  oidcAppType: OIDCAppType.OIDC_APP_TYPE_WEB,
  prefix: 'WEB',
  background: 'linear-gradient(40deg, #059669 30%, #047857)',
  protocol: 'OIDC',
};

export const USER_AGENT_TYPE: RadioItemAppType = {
  // key: AppType.USER_AGENT,
  titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.TITLE',
  descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.DESCRIPTION',
  createType: AppCreateType.OIDC,
  oidcAppType: OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
  prefix: 'UA',
  background: 'linear-gradient(40deg, #dc2626 30%, #db2777)',
  protocol: 'OIDC',
};

export const NATIVE_TYPE: RadioItemAppType = {
  // key: AppType.NATIVE,
  titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.TITLE',
  descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.DESCRIPTION',
  createType: AppCreateType.OIDC,
  oidcAppType: OIDCAppType.OIDC_APP_TYPE_NATIVE,
  prefix: 'N',
  background: 'linear-gradient(40deg, #306ccc 30%, #4f46e5)',
  protocol: 'OIDC',
};

export const API_TYPE: RadioItemAppType = {
  // key: AppType.API,
  titleI18nKey: 'APP.API.SELECTION.TITLE',
  descI18nKey: 'APP.API.SELECTION.DESCRIPTION',
  createType: AppCreateType.API,
  prefix: 'API',
  background: 'linear-gradient(40deg, #1f2937, #111827)',
  protocol: 'OIDC',
};

export const SAML_TYPE: RadioItemAppType = {
  titleI18nKey: 'APP.SAML.SELECTION.TITLE',
  descI18nKey: 'APP.SAML.SELECTION.DESCRIPTION',
  createType: AppCreateType.SAML,
  prefix: 'SAML',
  background: 'linear-gradient(40deg,rgb(110, 56, 124), rgb(88, 37, 103))',
  protocol: 'SAML',
};
