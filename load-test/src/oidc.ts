import { JSONObject, check, fail } from 'k6';
import encoding from 'k6/encoding';
import http from 'k6/http';
import { Trend } from 'k6/metrics';
import url from './url';

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
    'openid configuration': (r) => r.status == 200 || fail('unable to load openid configuration'),
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
    'userinfo status ok': (r) => r.status === 200,
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
    'introspect status ok': (r) => r.status === 200,
  });

  introspectTrend.add(res.timings.duration);
}

const clientCredentialsTrend = new Trend('oidc_client_credentials_duration', true);
export function clientCredentials(clientId: string, clientSecret: string): Promise<Tokens> {
  return new Promise((resolve, reject) => {
  const response = http.asyncRequest('POST', configuration().token_endpoint,
      {
        grant_type: "client_credentials",
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
        'client credentials status ok': (r) => r.status === 200,
      }) || reject(`client credentials request failed (client id: ${clientId}) status: ${res.status} body: ${res.body}`);

      clientCredentialsTrend.add(res.timings.duration);
      const tokens = new Tokens(res.json() as JSONObject)
      check(tokens, {
        'client credentials token ok': (t) => t.accessToken !== undefined,
      }) || reject(`client credentials access token missing (client id: ${clientId}`);

      resolve(tokens)
    });
  });
}