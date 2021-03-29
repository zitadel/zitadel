import { OIDCAppType } from 'src/app/proto/generated/zitadel/app_pb';

// export enum AppType {
//     "WEB",
//     "USER_AGENT",
//     "NATIVE",
//     "API"
// }

export enum AppCreateType {
    API = "API",
    OIDC = "OIDC"
}

export interface RadioItemAppType {
    // key: string;
    createType: AppCreateType;
    oidcAppType?: OIDCAppType;
    titleI18nKey: string;
    descI18nKey: string;
    prefix: string;
    background: string;
}

export const WEB_TYPE = {
    // key: AppType.WEB,
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.DESCRIPTION',
    createType: AppCreateType.OIDC,
    oidcAppType: OIDCAppType.OIDC_APP_TYPE_WEB,
    prefix: 'WEB',
    background: 'rgb(80, 110, 110)',
};

export const USER_AGENT_TYPE = {
    // key: AppType.USER_AGENT,
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.DESCRIPTION',
    createType: AppCreateType.OIDC,
    oidcAppType: OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
    prefix: 'UA',
    background: '#6a506e',
};

export const NATIVE_TYPE = {
    // key: AppType.NATIVE,
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.DESCRIPTION',
    createType: AppCreateType.OIDC,
    oidcAppType: OIDCAppType.OIDC_APP_TYPE_NATIVE,
    prefix: 'N',
    background: '#595d80',
};

export const API_TYPE = {
    // key: AppType.API,
    titleI18nKey: 'APP.API.SELECTION.TITLE',
    descI18nKey: 'APP.API.SELECTION.DESCRIPTION',
    createType: AppCreateType.API,
    prefix: 'API',
    background: 'rgb(73,73,73)',
};
