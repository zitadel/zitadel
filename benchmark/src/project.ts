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
      addProjectTrend.add(res.timings.duration);
      // Return after rejecting. Carrying on to res.json() threw on a non-JSON body and the
      // uncaught parse error replaced this message: on 2026-09-28 ten HTTP 500s with an HTML
      // body aborted introspect's setup, reported only as "invalid character '<'".
      if (
        !check(res, {
          'add project status ok': (r) => r.status >= 200 && r.status < 300,
        })
      ) {
        reject(`unable to add project status: ${res.status} body: ${String(res.body).slice(0, 200)}`);
        return;
      }
      resolve(res.json() as Project);
    });
  });
}
