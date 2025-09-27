import { Trend } from 'k6/metrics';
import { Org } from './org';
import http from 'k6/http';
import url from './url';
import { check } from 'k6';

export type Session = {
  challenges: any;
  id: string;
  token: string;
};

const addSessionTrend = new Trend('session_add_session_duration', true);
export function createSession(org: Org, accessToken: string, checks?: any): Promise<Session> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url('/v2/sessions'), checks ? JSON.stringify({ checks: checks }) : null, {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });
    response.then((res) => {
      check(res, {
        'add Session status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to add Session status: ${res.status} body: ${res.body}`);

      addSessionTrend.add(res.timings.duration);
      resolve(res.json() as Session);
    });
  });
}

const setSessionTrend = new Trend('session_set_session_duration', true);
export function setSession(id: string, session: any, accessToken: string, challenges?: any, checks?: any): Promise<Session> {
  const body = {
    sessionToken: session.sessionToken,
    checks: checks ? checks : null,
    challenges: challenges ? challenges : null,
  };
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('PATCH', url(`/v2/sessions/${id}`), JSON.stringify(body), {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        // 'x-zitadel-orgid': org.organizationId,
      },
    });
    response.then((res) => {
      check(res, {
        'set Session status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to set Session status: ${res.status} body: ${res.body}`);

      setSessionTrend.add(res.timings.duration);
      resolve(res.json() as Session);
    });
  });
}
