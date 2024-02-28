import { Framework } from '@netlify/framework-info/lib/types';
import { AddOIDCAppRequest } from '../proto/generated/zitadel/management_pb';
import { FrameworkName } from '@netlify/framework-info/lib/generated/frameworkNames';
import { OIDCAppType, OIDCAuthMethodType, OIDCGrantType, OIDCResponseType } from '../proto/generated/zitadel/app_pb';

type OidcAppConfigurations = {
  [framework: string]: AddOIDCAppRequest;
};

export const OIDC_CONFIGURATIONS: OidcAppConfigurations = {
  // user agent applications
  ['angular']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_USER_AGENT)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:4200/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:4200']),
  ['react']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_USER_AGENT)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['vue']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_USER_AGENT)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:5173/auth/signinwin/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:5173']),
  // web applications
  ['next']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['java']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp']),
  ['symfony']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8000/login_check'])
    .setPostLogoutRedirectUrisList(['http://localhost:8000/logout']),
  ['django']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8000/oidc/callback/'])
    .setPostLogoutRedirectUrisList(['http://localhost:8000/oidc/logout/ ']),
  // native
  ['flutter']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_NATIVE)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:4444/auth.html', 'com.example.zitadelflutter'])
    .setPostLogoutRedirectUrisList(['http://localhost:4444', 'com.example.zitadelflutter']),
};
