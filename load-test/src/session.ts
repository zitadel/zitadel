import { Trend } from 'k6/metrics';
import { Org } from './org';
import http from 'k6/http';
import url from './url';
import { check } from 'k6';
import { User } from './user';
import { b64encode } from 'k6/encoding';

export type Session = {
  id: string;
};

const addSessionTrend = new Trend('session_add_session_duration', true);
export function createSession(user: User, org: Org, accessToken: string): Promise<Session> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/v2beta/sessions'),
      JSON.stringify({
        checks: {
            user: {
                userId: user.userId,
            }
        }
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
        'add Session status ok': (r) => r.status === 201,
      }) || reject(`unable to add Session status: ${res.status} body: ${res.body}`);

      addSessionTrend.add(res.timings.duration);
      resolve(res.json() as Session);
    });
  });
}

const addEmptySessionTrend = new Trend('session_add_empty_session_duration', true);
export function createEmptySession(org: Org, accessToken: string): Promise<Session> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/v2beta/sessions'),
      null,
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
        'add Session status ok': (r) => r.status === 201,
      }) || reject(`unable to add Session status: ${res.status} body: ${res.body}`);

      addEmptySessionTrend.add(res.timings.duration);
      resolve({id: res.json('sessionId')} as Session);
    });
  });
}

const updateSessionTrend = new Trend('session_update_session_duration', true);
export function updateSession(session: Session, org: Org, accessToken: string): Promise<void> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'PATCH',
      url(`/v2beta/sessions/${session.id}`),
      JSON.stringify({
        metadata: {
          time: b64encode(new Date().toISOString(), "url"),
        }
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
        'add Session status ok': (r) => r.status === 200,
      }) || reject(`unable to update Session status: ${res.status} body: ${res.body}`);

      updateSessionTrend.add(res.timings.duration);
      resolve();
    });
  });
}