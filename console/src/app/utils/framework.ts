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
    .setRedirectUrisList(['http://localhost:3000/auth/callback/zitadel'])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ]),
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
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ]),
  // web applications (SSR)
  ['next']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['astro']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['hono']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['nestjs']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['nuxtjs']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['qwik']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['solidstart']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callack'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/api/auth/logout/callback'])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['sveltekit']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['spring-boot']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp'])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['java']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp'])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['symfony']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8000/login_check'])
    .setPostLogoutRedirectUrisList(['http://localhost:8000/logout'])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['django']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8000/oidc/callback/'])
    .setPostLogoutRedirectUrisList(['http://localhost:8000/oidc/logout/ '])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['express']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['fastify']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList([
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ])
    .setPostLogoutRedirectUrisList([
      'http://localhost:3000/api/auth/logout/callback',
      'http://localhost:3000/auth/logout/callback',
    ])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  ['flask']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://127.0.0.1:5000/callback'])
    .setPostLogoutRedirectUrisList(['http://127.0.0.1:5000'])
    .setIdTokenRoleAssertion(true)
    .setIdTokenUserinfoAssertion(true),
  // native
  ['flutter']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_NATIVE)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:4444/auth.html', 'com.example.zitadelflutter'])
    .setPostLogoutRedirectUrisList(['http://localhost:4444', 'com.example.zitadelflutter']),
};
