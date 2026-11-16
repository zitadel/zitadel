import { Trend } from 'k6/metrics';
import { Org } from './org';
import http, { RefinedResponse } from 'k6/http';
import url from './url';
import { check } from 'k6';
import { setTimeout } from 'k6/timers';
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
// Diagnostics from the last attempt, so an exhausted retry loop can say *why* it gave
// up rather than only that it did. Two failures on 2026-08-27 burned 30 attempts each
// and reported nothing more useful than "not readable", which cost a whole cycle.
export type ReadUserResult = { user?: any; lastStatus: number; lastBody: string; attempts: number };

// Backoff that yields to the event loop. `sleep()` would not: k6 documents that it "blocks VU
// execution and prevents promises from resolving and event handlers from running", and the same
// is true of the synchronous `http.get`. That matters here and nowhere else in this file, because
// this is the one helper called from inside a `.then()` callback while up to MaxVUs sibling
// creations are still in flight -- blocking serialises every one of them behind this loop.
function delay(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function readUserAfterCreate(path: string, userId: string, org: Org, accessToken: string): Promise<ReadUserResult> {
  let lastStatus = -1;
  let lastBody = '';
  let attempt = 0;
  for (attempt = 0; attempt < 60; attempt++) {
    const res = await http.asyncRequest('GET', url(`${path}/${userId}`), null, {
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
      const user = res.json('user') as any;
      // The user record can come back before its login names have been projected.
      // Accepting a merely-truthy user here returns a partial record, and every
      // caller immediately does loginNames[0] on it -- so this must wait for the
      // field itself, not just for the object. Observed at 600 VUs on
      // human_password_login, 2026-08-27: ten retries available, none consumed.
      if (user && Array.isArray(user.loginNames) && user.loginNames.length > 0) {
        // A loop that only reports itself when it gives up says nothing about how close it
        // came, which is how the budget below got guessed at three times before it was
        // measured once. Deliberately a log line and not a Trend: a Trend recorded here
        // reaches output.json for manipulate_user, where this runs inside the measured
        // iteration, and moving a published series to observe it would defeat the point.
        if (attempt > 0) {
          console.warn(
            `read-back retried: user ${userId} readable after ${attempt + 1} attempts, previous status ${lastStatus}`,
          );
        }
        return { user, lastStatus: res.status, lastBody: '', attempts: attempt + 1 };
      }
    }
    lastStatus = res.status;
    lastBody = (res.body ? String(res.body) : '<null body>').slice(0, 200);
    // ~100s of waiting in total (backoff rises to 2s and stays there). Sized against a
    // measured behaviour, not a guess: unloaded, login names are present on the *first*
    // GET, but a burst of MaxVUs concurrent creations queues them behind the projection
    // handler advisory lock -- the largest single database cost in the whole v4.17.1 sweep.
    //
    // Two things measured on 2026-08-31 that this budget cannot fix, recorded so the next
    // person does not widen it a fourth time:
    //
    // Retries are never consumed successfully. Across every run that day the warn above
    // fired zero times -- each user either had its login names on the *first* GET or never
    // had them at all. The read-back is bimodal, so the budget only sets how long we wait
    // before reporting a failure, not how many succeed.
    //
    // And some of those failures are permanent. In one 200-VU run all 200 creates returned
    // 2xx and, ten minutes later, the query side was still missing zitizen-6 entirely and
    // zitizen-7/-24/-25 had no login names. Widening the budget to ~220s changed nothing
    // except how long it took to give up, which is why it was reverted to this.
    //
    // Note also that serialising this loop used to grant an accidental grace period: user N
    // was not asked until users 0..N-1 had finished, so a late user got the queue's time for
    // free. Awaiting concurrently aligns every deadline on the same instant. That is correct
    // and strictly less forgiving, and it is a real effect -- it was simply not the cause of
    // the failures above, which were loss rather than lag.
    await delay(Math.min(100 * (attempt + 1), 2000));
  }
  return { user: undefined, lastStatus, lastBody, attempts: attempt };
}

const createHumanTrend = new Trend('user_create_human_duration', true);
export async function createHuman(username: string, org: Org, accessToken: string): Promise<Human> {
  const res = await http.asyncRequest(
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

  if (
    !check(res, {
      'create user is status ok': (r) => r.status >= 200 && r.status < 300,
    })
  ) {
    throw new Error(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
  }
  createHumanTrend.add(res.timings.duration);

  const userId = res.json('userId') as string;
  const read = await readUserAfterCreate('/v2/users', userId, org, accessToken);
  if (!read.user) {
    throw new Error(
      `user ${userId} not readable with login names after create ` +
        `(username: ${username}, attempts: ${read.attempts}, last status: ${read.lastStatus}, ` +
        `last body: ${read.lastBody})`,
    );
  }
  return read.user as unknown as Human;
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
    let response = http.asyncRequest('PATCH', url(`/v2/users/${userId}`), JSON.stringify(payload), {
      tags: { name: '/v2/users/{userId}' },
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
export async function createMachine(username: string, org: Org, accessToken: string): Promise<Machine> {
  const res = await http.asyncRequest(
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

  if (
    !check(res, {
      'create user is status ok': (r) => r.status >= 200 && r.status < 300,
    })
  ) {
    throw new Error(`unable to create user(username: ${username}) status: ${res.status} body: ${res.body}`);
  }
  createMachineTrend.add(res.timings.duration);

  const userId = res.json('userId') as string;
  const read = await readUserAfterCreate('/v2/users', userId, org, accessToken);
  if (!read.user) {
    throw new Error(
      `user ${userId} not readable with login names after create ` +
        `(username: ${username}, attempts: ${read.attempts}, last status: ${read.lastStatus}, ` +
        `last body: ${read.lastBody})`,
    );
  }
  return read.user as unknown as Machine;
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
    let response = http.asyncRequest('POST', url(`/v2/users/${userId}/lock`), null, {
      tags: { name: '/v2/users/{userId}/lock' },
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
    let response = http.asyncRequest('DELETE', url(`/v2/users/${userId}`), null, {
      tags: { name: '/v2/users/{userId}' },
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
