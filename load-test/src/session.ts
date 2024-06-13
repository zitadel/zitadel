import { Trend } from 'k6/metrics';
import { Org } from './org';
// import http from 'k6/http';
// import url from './url';
import { check } from 'k6';
import { User } from './user';
import grpc from 'k6/net/grpc';

export type Session = {
  id: string;
};

export const SessionClient = new grpc.Client();
SessionClient.load(['../../proto'], 'zitadel/session/v2beta/session_service.proto', 'zitadel/session/v2beta/session.proto', 'zitadel/object/v2beta/object.proto');

const addSessionTrend = new Trend('session_add_session_duration', true);
export function createSession(user: User, org: Org, accessToken: string): Promise<Session> {
  return new Promise((resolve, reject) => {
    const start = new Date();

    const response = SessionClient.invoke(
      'zitadel.session.v2beta.SessionService/CreateSession', 
      {
        checks: {
          user: {
              userId: user.userId,
          }
        }
      },
      {
        metadata: {
          authorization: `Bearer ${accessToken}`,
          'x-zitadel-orgid': org.organizationId
        }
      }
    );
    addSessionTrend.add(new Date().getTime() - start.getTime());

    check(response, {
      'add Session status ok': (r) => r.status === grpc.StatusOK,
    }) || reject(`unable to add Session status: ${response.status} body: ${JSON.stringify(response.message)} error: ${JSON.stringify(response.error)}`);
    resolve(response.message as Session);
    });
}
