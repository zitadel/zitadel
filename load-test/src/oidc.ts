import { JSONObject, check, fail } from 'k6';
import encoding from 'k6/encoding';
import http, { RequestBody, Response } from 'k6/http';
import { Trend } from 'k6/metrics';
import url from './url';
import { Client, Config } from './config';
// @ts-ignore Import module
import zitadel from 'k6/x/zitadel';

export class Tokens {
  idToken?: string;
  accessToken?: string;
  info?: any;

  constructor(res: JSONObject) {
    this.idToken = res.id_token ? res.id_token!.toString() : undefined;
    this.accessToken = res.access_token ? res.access_token!.toString() : undefined;
    this.info = this.idToken
      ? JSON.parse(encoding.b64decode(this.idToken?.split('.')[1].toString(), 'rawstd', 's'))
      : undefined;
  }
}

let oidcConfig: any | undefined;

function configuration() {
  if (oidcConfig !== undefined) {
    return oidcConfig;
  }

  const res = http.get(url('/.well-known/openid-configuration'));
  check(res, {
    'openid configuration': (r) => r.status >= 200 && r.status < 300 || fail('unable to load openid configuration'),
  });

  oidcConfig = res.json();
  return oidcConfig;
}

const userinfoTrend = new Trend('oidc_user_info_duration', true);
export function userinfo(token: string) {
  const userinfo = http.get(configuration().userinfo_endpoint, {
    headers: {
      authorization: 'Bearer ' + token,
      'Content-Type': 'application/json',
    },
  });

  check(userinfo, {
    'userinfo status ok': (r) => r.status >= 200 && r.status < 300,
  });

  userinfoTrend.add(userinfo.timings.duration);
}

const introspectTrend = new Trend('oidc_introspect_duration', true);
export function introspect(jwt: string, token: string) {
  const res = http.post(
    configuration().introspection_endpoint,
    {
      client_assertion: jwt,
      token: token,
      client_assertion_type: 'urn:ietf:params:oauth:client-assertion-type:jwt-bearer',
    },
    {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        alg: 'RS256',
      },
    },
  );
  check(res, {
    'introspect status ok': (r) => r.status >= 200 && r.status < 300,
  });

  introspectTrend.add(res.timings.duration);
}

const clientCredentialsTrend = new Trend('oidc_client_credentials_duration', true);
export function clientCredentials(clientId: string, clientSecret: string): Promise<Tokens> {
  return new Promise((resolve, reject) => {
    const response = http.asyncRequest(
      'POST',
      configuration().token_endpoint,
      {
        grant_type: 'client_credentials',
        scope: 'openid profile urn:zitadel:iam:org:project:id:zitadel:aud',
        client_id: clientId,
        client_secret: clientSecret,
      },
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
      },
    );
    response.then((res) => {
      check(res, {
        'client credentials status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`client credentials request failed (client id: ${clientId}) status: ${res.status} body: ${res.body}`);

      clientCredentialsTrend.add(res.timings.duration);
      const tokens = new Tokens(res.json() as JSONObject);
      check(tokens, {
        'client credentials token ok': (t) => t.accessToken !== undefined,
      }) || reject(`client credentials access token missing (client id: ${clientId}`);

      resolve(tokens);
    });
  });
}

export interface TokenRequest {
  payload(): RequestBody;
  headers(): { [name: string]: string };
}

const privateKey = open('../.keys/key.pem');

export class JWTProfileRequest implements TokenRequest {
  keyPayload!: {
    userId: string;
    expiration: number;
    keyId: string;
  };

  constructor(userId: string, keyId: string) {
    this.keyPayload = {
      userId: userId,
      // 1 minute
      expiration: 60 * 1_000_000_000,
      keyId: keyId,
    };
  }

  payload(): RequestBody {
    const assertion = zitadel.signJWTProfileAssertion(this.keyPayload.userId, this.keyPayload.keyId, {
      audience: [Config.host],
      expiration: this.keyPayload.expiration,
      key: privateKey,
    });
    return {
      grant_type: 'urn:ietf:params:oauth:grant-type:jwt-bearer',
      scope: 'openid urn:zitadel:iam:org:project:id:zitadel:aud',
      assertion: `${assertion}`,
    };
  }
  public headers(): { [name: string]: string } {
    return {
      'Content-Type': 'application/x-www-form-urlencoded',
    };
  }
}

const tokenDurationTrend = new Trend('oidc_token_duration', true);
export async function token(request: TokenRequest): Promise<Tokens> {
  return http
    .asyncRequest('POST', configuration().token_endpoint, request.payload(), {
      headers: request.headers(),
    })
    .then((res) => {
      tokenDurationTrend.add(res.timings.duration);
      check(res, {
        'token status ok': (r) => r.status >= 200 && r.status < 300,
        'access token returned': (r) => r.json('access_token')! != undefined && r.json('access_token')! != '',
      });
      return new Tokens(res.json() as JSONObject);
    });
}

const authRequestByIDTrend = new Trend('oidc_auth_request_by_id_duration', true);
export async function authRequestByID(id: string, tokens: any): Promise<Response> {
  const response = http.get(url(`/v2/oidc/auth_requests/${id}`), {
    headers: {
      Authorization: `Bearer ${tokens.accessToken}`,
    },
  });
  check(response, {
    'authorize status ok': (r) => r.status >= 200 && r.status < 300 || fail(`auth request by failed: ${JSON.stringify(r)}`),
  });
  authRequestByIDTrend.add(response.timings.duration);
  return response;
}

const finalizeAuthRequestTrend = new Trend('oidc_auth_request_finalize', true);
export async function finalizeAuthRequest(id: string, session: any, tokens: any): Promise<Response> {
  const res = await http.post(
    url(`/v2/oidc/auth_requests/${id}`),
    JSON.stringify({
      session: {
        sessionId: session.sessionId,
        sessionToken: session.sessionToken,
      },
    }),
    {
      headers: {
        Authorization: `Bearer ${tokens.accessToken}`,
        'Content-Type': 'application/json',
        // 'Accept': 'application/json',
        'x-zitadel-login-client': tokens.info.client_id,
      },
    },
  );
  check(res, {
    'finalize auth request status ok': (r) => r.status >= 200 && r.status < 300 || fail(`finalize auth request failed: ${JSON.stringify(r)}`),
  });
  finalizeAuthRequestTrend.add(res.timings.duration);

  return res;
}
