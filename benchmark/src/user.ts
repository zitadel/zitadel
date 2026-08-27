import { Trend } from 'k6/metrics';
import { Org } from './org';
import http, { RefinedResponse } from 'k6/http';
import url from './url';
import { check, sleep } from 'k6';
import { Metadata } from './metadata';

export type User = {
  userId: string;
  loginName: string;
  password: string;
};

export interface Human extends User {
  loginNames: string[];
}

// A user is created by a write, but its login names can only be read back from the
// query side. Under load that projection can lag behind the write, and the read
// then returns no user at all. Retry briefly rather than resolving undefined, which
// surfaces much later as "Cannot read property '0' of undefined" while mapping
// login names, hiding the real cause.
function readUserAfterCreate(path: string, userId: string, org: Org, accessToken: string): any {
  for (let attempt = 0; attempt < 10; attempt++) {
    const res = http.get(url(`${path}/${userId}`), {
      tags: { name: `${path}/{userId}` },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });
    // A timed-out request has status 0 and a null body. Calling res.json() on that
    // throws "GoError: the body is null" and kills the iteration instead of consuming
    // a retry, which is exactly the failure this retry loop exists to absorb.
    if (res.status >= 200 && res.status < 300 && res.body) {
      const user = res.json('user');
      if (user) {
        return user;
      }
    }
    sleep(0.1 * (attempt + 1));
  }
  return undefined;
}

const createHumanTrend = new Trend('user_create_human_duration', true);
export function createHuman(username: string, org: Org, accessToken: string): Promise<Human> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url('/v2/users/human'),
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
          email: `${username}@zitadel.com`,
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
        if (
          !check(res, {
            'create user is status ok': (r) => r.status >= 200 && r.status < 300,
          })
        ) {
          reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
          return;
        }
        createHumanTrend.add(res.timings.duration);

        const user = readUserAfterCreate('/v2/users', res.json('userId') as string, org, accessToken);
        if (!user) {
          reject(`user ${res.json('userId')} not readable after create (username: ${username})`);
          return;
        }
        resolve(user as unknown as Human);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

const setEmailOTPOnHumanTrend = new Trend('set_human_email_otp_duration', true);
export async function setEmailOTPOnHuman(user: User, org: Org, accessToken: string): Promise<void> {
  const response = await http.asyncRequest('POST', url(`/v2/users/${user.userId}/otp_email`), null, {
    tags: { name: '/v2/users/{userId}/otp_email' },
    headers: {
      authorization: `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
      'x-zitadel-orgid': org.organizationId,
    },
  });
  check(response, {
    'set email otp status ok': (r) => r.status >= 200 && r.status < 300,
  });
  setEmailOTPOnHumanTrend.add(response.timings.duration);

  return;
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
      tags: { name: '/v2beta/users/{userId}' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'update user is status ok': (r) => r.status >= 200 && r.status < 300,
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
        if (
          !check(res, {
            'create user is status ok': (r) => r.status >= 200 && r.status < 300,
          })
        ) {
          reject(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
          return;
        }
        createMachineTrend.add(res.timings.duration);

        const user = readUserAfterCreate('/v2beta/users', res.json('userId') as string, org, accessToken);
        if (!user) {
          reject(`user ${res.json('userId')} not readable after create (username: ${username})`);
          return;
        }
        resolve(user as unknown as Machine);
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
      tags: { name: '/management/v1/users/{userId}/pats' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });
    response.then((res) => {
      check(res, {
        'add pat status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to add pat (user id: ${userId}) status: ${res.status} body: ${res.body}`);

      addMachinePatTrend.add(res.timings.duration);
      resolve(res.json()! as MachinePat);
    });
  });
}

export type MachineSecret = {
  clientId: string;
  clientSecret: string;
};

const addMachineSecretTrend = new Trend('user_add_machine_secret_duration', true);
export function addMachineSecret(userId: string, org: Org, accessToken: string): Promise<MachineSecret> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('PUT', url(`/management/v1/users/${userId}/secret`), null, {
      tags: { name: '/management/v1/users/{userId}/secret' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });
    response.then((res) => {
      check(res, {
        'generate machine secret status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to generate machine secret (user id: ${userId}) status: ${res.status} body: ${res.body}`);

      addMachineSecretTrend.add(res.timings.duration);
      resolve(res.json()! as MachineSecret);
    });
  });
}

export type MachineKey = {
  keyId: string;
};

const addMachineKeyTrend = new Trend('user_add_machine_key_duration', true);
export function addMachineKey(userId: string, org: Org, accessToken: string, publicKey?: string): Promise<MachineKey> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest(
      'POST',
      url(`/management/v1/users/${userId}/keys`),
      JSON.stringify({
        type: 'KEY_TYPE_JSON',
        userId: userId,
        // base64 encoded public key
        publicKey: publicKey,
      }),
      {
        tags: { name: '/management/v1/users/{userId}/keys' },
        headers: {
          authorization: `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
          'x-zitadel-orgid': org.organizationId,
        },
      },
    );
    response.then((res) => {
      check(res, {
        'generate machine key status ok': (r) => r.status >= 200 && r.status < 300,
      }) || reject(`unable to generate machine Key (user id: ${userId}) status: ${res.status} body: ${res.body}`);

      addMachineKeyTrend.add(res.timings.duration);
      resolve(res.json()! as MachineKey);
    });
  });
}

const lockUserTrend = new Trend('lock_user_duration', true);
export function lockUser(userId: string, org: Org, accessToken: string): Promise<RefinedResponse<any>> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url(`/v2beta/users/${userId}/lock`), null, {
      tags: { name: '/v2beta/users/{userId}/lock' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'lock user is status ok': (r) => r.status >= 200 && r.status < 300,
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
      tags: { name: '/v2beta/users/{userId}' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'x-zitadel-orgid': org.organizationId,
      },
    });

    response
      .then((res) => {
        check(res, {
          'delete user is status ok': (r) => r.status >= 200 && r.status < 300,
        });
        deleteUserTrend.add(res.timings.duration);
        resolve(res);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

const setUserMetadataTrend = new Trend('set_user_metadata_duration', true);
export function setUserMetadata(metadata: Metadata[], userId: string, accessToken: string): Promise<RefinedResponse<any>> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url(`/v2/users/${userId}/metadata`), JSON.stringify({ metadata: metadata }), {
      tags: { name: '/v2/users/{userId}/metadata' },
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
    });

    response
      .then((res) => {
        check(res, {
          'set user metadata is status ok': (r) => r.status >= 200 && r.status < 300,
        }) || console.log(`unable to set user metadata (user id: ${userId}) status: ${res.status} body: ${res.body}`);
        setUserMetadataTrend.add(res.timings.duration);
        resolve(res);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}

export type ListUsersRequest = {
  query?: {
    limit?: number;
    offset?: number;
  };
  queries: {
    loginNameQuery?: {
      loginName: string;
      method:
        | 'TEXT_QUERY_METHOD_EQUALS'
        | 'TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE'
        | 'TEXT_QUERY_METHOD_STARTS_WITH'
        | 'TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE'
        | 'TEXT_QUERY_METHOD_CONTAINS'
        | 'TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE'
        | 'TEXT_QUERY_METHOD_ENDS_WITH'
        | 'TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE';
    };
    organizationIdQuery?: {
      organizationId: string;
    };
    metadataKeyFilter?: {
      key: string;
      method: 'TEXT_FILTER_METHOD_EQUALS' | 'TEXT_FILTER_METHOD_CONTAINS' | 'TEXT_FILTER_METHOD_CONTAINS_IGNORE_CASE';
    };
    metadataValueFilter?: {
      value: string;
      method: 'BYTE_FILTER_METHOD_EQUALS' | 'BYTE_FILTER_METHOD_NOT_EQUALS';
    };
  }[];
};

export type ListUsersResult = {
  details: {
    totalResult: number;
  };
  result?: unknown[];
};

const listUsersTrend = new Trend('list_users_duration', true);
export function listUsers(body: ListUsersRequest, accessToken: string): Promise<ListUsersResult> {
  return new Promise((resolve, reject) => {
    let response = http.asyncRequest('POST', url(`/v2/users`), JSON.stringify(body), {
      headers: {
        authorization: `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
      },
    });

    response
      .then((res) => {
        check(res, {
          'list users is status ok': (r) => r.status >= 200 && r.status < 300,
        }) || console.log(`unable to list users status: ${res.status} body: ${res.body}`);
        listUsersTrend.add(res.timings.duration);
        resolve(res.json()! as ListUsersResult);
      })
      .catch((reason) => {
        reject(reason);
      });
  });
}
