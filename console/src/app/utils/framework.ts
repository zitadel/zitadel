import { AddOIDCAppRequest } from '../proto/generated/zitadel/management_pb';
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
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/signedout']),
  ['react']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_USER_AGENT)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['vue']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_USER_AGENT)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/']),
  // web applications
  ['next']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['spring']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['go']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8089/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:8089']),
  ['symfony']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['laravel']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback/zitadel']),
  ['fastapi']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['flask']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['django']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['astro']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['hono']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['solidstart']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  // native
  ['flutter']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_NATIVE)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:4444/auth.html', 'com.example.zitadelflutter'])
    .setPostLogoutRedirectUrisList(['http://localhost:4444', 'com.example.zitadelflutter']),
};
