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
        tags: { name: '/management/v1/projects/{projectId}/apps/api' },
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );
    response.then((res) => {
      check(res, {
        'add api status ok': (r) => r.status >= 200 && r.status < 300,
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
        tags: { name: '/management/v1/projects/{projectId}/apps/{appId}/keys' },
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );
    response.then((res) => {
      check(res, {
        'add app key status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to add app key project: ${projectId} app: ${appId} status: ${res.status} body: ${res.body}`);
      resolve(res.json() as AppKey);

      addAppKeyTrend.add(res.timings.duration);
    });
  });
}

// setMinimalIntrospection enables the minimal_introspection setting on an existing API app via the
// v2 ApplicationService (Connect protocol, plain JSON POST). With it enabled, the introspection
// endpoint skips the userinfo/claims lookup for tokens issued to this app.
const setMinimalIntrospectionTrend = new Trend('app_set_minimal_introspection_duration', true);
export function setMinimalIntrospection(appId: string, projectId: string, org: Org, accessToken: string): Promise<void> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/zitadel.application.v2.ApplicationService/UpdateApplication'),
      JSON.stringify({
        projectId: projectId,
        applicationId: appId,
        apiConfiguration: {
          authMethodType: 'API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT',
          minimalIntrospection: true,
        },
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
        'set minimal introspection status ok': (r) => r.status >= 200 && r.status < 300,
      }) ||
        reject(`unable to set minimal introspection project: ${projectId} app: ${appId} status: ${res.status} body: ${res.body}`);
      resolve();

      setMinimalIntrospectionTrend.add(res.timings.duration);
    });
  });
}
