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
  ['java']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp']),
  ['client-java']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:18080/webapp/login/oauth2/code/zitadel'])
    .setPostLogoutRedirectUrisList(['http://localhost:18080/webapp']),
  ['spring']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000/auth/logout/callback']),
  ['client-php']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8080/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:8080']),
  ['client-python']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:8000']),
  ['client-node']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['client-go']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:8089/auth/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:8089']),
  ['client-ruby']: new AddOIDCAppRequest()
    .setAppType(OIDCAppType.OIDC_APP_TYPE_WEB)
    .setAuthMethodType(OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC)
    .setResponseTypesList([OIDCResponseType.OIDC_RESPONSE_TYPE_CODE])
    .setGrantTypesList([OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE])
    .setRedirectUrisList(['http://localhost:3000/callback'])
    .setPostLogoutRedirectUrisList(['http://localhost:3000']),
  ['symfony']: new AddOIDCAppRequest()
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
