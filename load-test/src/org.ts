import http from 'k6/http';
import { Trend } from 'k6/metrics';
import url from './url';
import { Config } from './config';
import { check } from 'k6';

export type Org = {
  organizationId: string;
};

const createOrgTrend = new Trend('org_create_org_duration', true);
export function createOrg(accessToken: string): Promise<Org> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/v2beta/organizations'),
      JSON.stringify({
        name: `load-test-${new Date(Date.now()).toISOString()}`,
      }),
      {
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': Config.orgId,
        },
      },
    );

    response.then((res) => {
      check(res, {
        'org created': (r) => {
          return r !== undefined && r.status === 201;
        },
      }) || reject(`unable to create org status: ${res.status} || body: ${res.body}`);

      createOrgTrend.add(res.timings.duration);

      resolve(res.json() as Org);
    });
  });
}

export function removeOrg(org: Org, accessToken: string) {
  const response = http.del(url('/management/v1/orgs/me'), null, {
    headers: {
      authorization: `Bearer ${accessToken}`,
      'x-zitadel-orgid': org.organizationId,
    },
  });

  check(response, {
    'org removed': (r) => r.status === 200,
  }) || console.log(`status: ${response.status} || body: ${response.body}|| org: ${JSON.stringify(org)}`);

  return response.json();
}
