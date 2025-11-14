import { MessageInitShape } from '@bufbuild/protobuf';
import { CreateApplicationRequestSchema } from '@zitadel/proto/zitadel/app/v2beta/app_service_pb';
import { OIDCAppType, OIDCAuthMethodType, OIDCGrantType, OIDCResponseType } from '@zitadel/proto/zitadel/app/v2beta/oidc_pb';
import frameworkDefinition from '../../../../docs/frameworks.json';

type OIDCConfiguration = Extract<
  MessageInitShape<typeof CreateApplicationRequestSchema>['creationRequestType'],
  { case: 'oidcRequest' }
>['value'];

export const OIDC_CONFIGURATIONS = {
  // user agent applications (SPA)
  angular: {
    appType: OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:3000/auth/callback/zitadel'],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
  },
  react: {
    appType: OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:3000/callback'],
    postLogoutRedirectUris: ['http://localhost:3000'],
  },
  vue: {
    appType: OIDCAppType.OIDC_APP_TYPE_USER_AGENT,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
  },
  // web applications (SSR)
  next: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  astro: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  hono: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  nestjs: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  nuxtjs: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  qwik: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  solidstart: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  sveltekit: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  'spring-boot': {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:8080/login/oauth2/code/zitadel'],
    postLogoutRedirectUris: ['http://localhost:8080'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  java: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:18080/webapp/login/oauth2/code/zitadel'],
    postLogoutRedirectUris: ['http://localhost:18080/webapp'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  symfony: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:8000/login_check'],
    postLogoutRedirectUris: ['http://localhost:8000/logout'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  django: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://localhost:8000/oidc/callback/'],
    postLogoutRedirectUris: ['http://localhost:8000/oidc/logout/'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  express: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  fastify: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: [
      'http://localhost:3000/auth/callback',
      'http://localhost:3000/api/auth/callback',
      'http://localhost:3000/api/auth/callback/zitadel',
      'http://localhost:3000/auth/callback/zitadel',
    ],
    postLogoutRedirectUris: ['http://localhost:3000/api/auth/logout/callback', 'http://localhost:3000/auth/logout/callback'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  flask: {
    appType: OIDCAppType.OIDC_APP_TYPE_WEB,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_BASIC,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://127.0.0.1:5000/callback'],
    postLogoutRedirectUris: ['http://127.0.0.1:5000'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
  // native
  flutter: {
    appType: OIDCAppType.OIDC_APP_TYPE_NATIVE,
    authMethodType: OIDCAuthMethodType.OIDC_AUTH_METHOD_TYPE_NONE,
    responseTypes: [OIDCResponseType.OIDC_RESPONSE_TYPE_CODE],
    grantTypes: [OIDCGrantType.OIDC_GRANT_TYPE_AUTHORIZATION_CODE],
    redirectUris: ['http://127.0.0.1:5000/callback'],
    postLogoutRedirectUris: ['http://127.0.0.1:5000'],
    idTokenRoleAssertion: true,
    idTokenUserinfoAssertion: true,
  },
} satisfies Record<string, OIDCConfiguration>;

export const frameworks = frameworkDefinition.map((f) => ({
  ...f,
  imgSrcDark: `assets${f.imgSrcDark}`,
  imgSrcLight: `assets${f.imgSrcLight ?? f.imgSrcDark}`,
}));

export const frameworksWithOidcConfiguration = frameworks
  .map((f) => {
    const id = f.id as unknown as keyof typeof OIDC_CONFIGURATIONS;
    const oidcConfiguration: (typeof OIDC_CONFIGURATIONS)[typeof id] | undefined = OIDC_CONFIGURATIONS[id];
    if (oidcConfiguration) {
      return {
        ...f,
        id,
      };
    }

    return undefined;
  })
  .filter((f) => !!f);
