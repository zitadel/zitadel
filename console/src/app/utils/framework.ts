import { AddOIDCAppRequest } from '../proto/generated/zitadel/management_pb';
import { OIDCAppType, OIDCAuthMethodType, OIDCGrantType, OIDCResponseType } from '../proto/generated/zitadel/app_pb';

type OidcAppConfigurations = {
  [framework: string]: AddOIDCAppRequest;
};

export const OIDC_CONFIGURATIONS: OidcAppConfigurations = {
  // user agent applications (SPA)
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
    .setPostLogoutRedirectUrisList(['http://localhost:5173/']),
  // web applications (SSR)
  ['next']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['astro']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['hono']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['nestjs']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['nuxtjs']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['solidstart']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/api/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['sveltekit']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['spring-boot']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp']),
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
  ['flask']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://127.0.0.1:5000/callback'])
    .setPostLogoutRedirectUrisList(['http://127.0.0.1:5000']),
  // native
  ['flutter']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_NATIVE)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:4444/auth.html', 'com.example.zitadelflutter'])
    .setPostLogoutRedirectUrisList(['http://localhost:4444', 'com.example.zitadelflutter']),
};
