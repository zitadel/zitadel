import { OIDCApplicationType } from 'src/app/proto/generated/management_pb';

export const WEB_TYPE = {
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.WEB.DESCRIPTION',
    type: OIDCApplicationType.OIDCAPPLICATIONTYPE_WEB,
    prefix: 'WEB',
    background: 'rgb(80, 110, 110)',
};

export const USER_AGENT_TYPE = {
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.USERAGENT.DESCRIPTION',
    type: OIDCApplicationType.OIDCAPPLICATIONTYPE_USER_AGENT,
    prefix: 'UA',
    background: '#6a506e',
};

export const NATIVE_TYPE = {
    titleI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.TITLE',
    descI18nKey: 'APP.OIDC.SELECTION.APPTYPE.NATIVE.DESCRIPTION',
    type: OIDCApplicationType.OIDCAPPLICATIONTYPE_NATIVE,
    prefix: 'N',
    background: '#595d80',
};
