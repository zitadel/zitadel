import { Trend } from 'k6/metrics';
import { Org } from './org';
import http from 'k6/http';
import url from './url';
import { check } from 'k6';

export type API = {
  appId: string;
};

const addAPITrend = new Trend('app_add_app_duration', true);
export function createAPI(name: string, projectId: string, org: Org, accessToken: string): Promise<API> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url(`/management/v1/projects/${projectId}/apps/api`),
      JSON.stringify({
        name: name,
        authMethodType: 'API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT',
      }),
      {
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );
    response.then((res) => {
      check(res, {
        'add api status ok': (r) => r.status === 200,
      }) || reject(`unable to add api project: ${projectId} status: ${res.status} body: ${res.body}`);
      resolve(res.json() as API);

      addAPITrend.add(res.timings.duration);
    });
  });
}

export type AppKey = {
  keyDetails: string;
};

const addAppKeyTrend = new Trend('app_add_app_key_duration', true);
export function createAppKey(appId: string, projectId: string, org: Org, accessToken: string): Promise<AppKey> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url(`/management/v1/projects/${projectId}/apps/${appId}/keys`),
      JSON.stringify({
        type: 'KEY_TYPE_JSON',
      }),
      {
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );
    response.then((res) => {
      check(res, {
        'add app key status ok': (r) => r.status === 200,
      }) || reject(`unable to add app key project: ${projectId} app: ${appId} status: ${res.status} body: ${res.body}`);
      resolve(res.json() as AppKey);

      addAppKeyTrend.add(res.timings.duration);
    });
  });
}
