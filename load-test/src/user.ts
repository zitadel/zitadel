import { Trend } from 'k6/metrics';
import { Org } from './org';
import http, { RefinedResponse } from 'k6/http';
import url from './url';
import { check } from 'k6';

export type User = {
  userId: string;
  loginName: string;
  password: string;
};

export interface Human extends User {
  loginNames: string[];
}

const createHumanTrend = new Trend('user_create_human_duration', true);
export function createHuman(username: string, org: Org, accessToken: string): Promise<Human> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/v2beta/users/human'),
      JSON.stringify({
        username: username,
        organization: {
          orgId: org.organizationId,
        },
        profile: {
          givenName: 'Gigi',
          familyName: 'Zitizen',
        },
        email: {
          email: `zitizen-@caos.ch`,
          isVerified: true,
        },
        password: {
          password: 'Password1!',
          changeRequired: false,
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

    response
      .then((res) => {
        check(res, {
          'create user is status ok': (r) => r.status === 201,
        }) || reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
        createHumanTrend.add(res.timings.duration);

        const user = http.get(url(`/v2beta/users/${res.json('userId')!}`), {
          headers: {
            authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
            'x-zitadel-orgid': org.organizationId,
          },
        });
        resolve(user.json('user')! as unknown as Human);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

const updateHumanTrend = new Trend('update_human_duration', true);
export function updateHuman(
  payload: any = {},
  userId: string,
  org: Org,
  accessToken: string,
): Promise<RefinedResponse<any>> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('PUT', url(`/v2beta/users/${userId}`), JSON.stringify(payload), {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'update user is status ok': (r) => r.status === 201,
        });
        updateHumanTrend.add(res.timings.duration);
        resolve(res);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

export interface Machine extends User {
  loginNames: string[];
}

const createMachineTrend = new Trend('user_create_machine_duration', true);
export function createMachine(username: string, org: Org, accessToken: string): Promise<Machine> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/management/v1/users/machine'),
      JSON.stringify({
        userName: username,
        name: username,
        // bearer
        access_token_type: 0,
      }),
      {
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );

    response
      .then((res) => {
        check(res, {
          'create user is status ok': (r) => r.status === 200,
        }) || reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
        createMachineTrend.add(res.timings.duration);

        const user = http.get(url(`/v2beta/users/${res.json('userId')!}`), {
          headers: {
            authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
            'x-zitadel-orgid': org.organizationId,
          },
        });
        resolve(user.json('user')! as unknown as Machine);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

export type MachinePat = {
  token: string;
};

const addMachinePatTrend = new Trend('user_add_machine_pat_duration', true);
export function addMachinePat(userId: string, org: Org, accessToken: string): Promise<MachinePat> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url(`/management/v1/users/${userId}/pats`), null, {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });
    response.then((res) => {
      check(res, {
        'add pat status ok': (r) => r.status === 200,
      }) || reject(`unable to add pat (user id: ${userId}) status: ${res.status} body: ${res.body}`);

      addMachinePatTrend.add(res.timings.duration);
      resolve(res.json()! as MachinePat);
    });
  });
}

const lockUserTrend = new Trend('lock_user_duration', true);
export function lockUser(userId: string, org: Org, accessToken: string): Promise<RefinedResponse<any>> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url(`/v2beta/users/${userId}/lock`), null, {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'update user is status ok': (r) => r.status === 201,
        });
        lockUserTrend.add(res.timings.duration);
        resolve(res);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

const deleteUserTrend = new Trend('delete_user_duration', true);
export function deleteUser(userId: string, org: Org, accessToken: string): Promise<RefinedResponse<any>> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('DELETE', url(`/v2beta/users/${userId}`), null, {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'update user is status ok': (r) => r.status === 201,
        });
        deleteUserTrend.add(res.timings.duration);
        resolve(res);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}
