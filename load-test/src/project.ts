import { Trend } from 'k6/metrics';
import { Org } from './org';
import http from 'k6/http';
import url from './url';
import { check } from 'k6';

export type Project = {
  id: string;
};

const addProjectTrend = new Trend('project_add_project_duration', true);
export function createProject(name: string, org: Org, accessToken: string): Promise<Project> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/management/v1/projects'),
      JSON.stringify({
        name: name,
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
        'add project status ok': (r) => r.status === 200,
      }) || reject(`unable to add project status: ${res.status} body: ${res.body}`);

      addProjectTrend.add(res.timings.duration);
      resolve(res.json() as Project);
    });
  });
}
