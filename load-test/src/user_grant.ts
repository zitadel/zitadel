import http from 'k6/http';
import { Trend } from 'k6/metrics';
import url from './url';
import { check, fail } from 'k6';
import { Project } from './project';
import { Org } from './org';

export type UserGrant = {
  userGrantId: string;
};

const addUserGrantTrend = new Trend('user_grant_add', true);
export async function addUserGrant(
  org: Org,
  userId: string,
  project: Project,
  roles: string[],
  accessToken: string,
): Promise<UserGrant> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url(`/management/v1/users/${userId}/grants`),
      JSON.stringify({
        projectId: project.id,
        roleKeys: roles,
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
        'add User Grant status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to add User Grant status: ${res.status} body: ${res.body}`);

      addUserGrantTrend.add(res.timings.duration);
      resolve(res.json() as UserGrant);
    });
  });
}