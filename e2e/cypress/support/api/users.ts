import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API } from './types';

export function ensureHumanUserExists(api: API, username: string){
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    `${api.mgmtBaseURL}/users/human`,
    {
      user_name: username,
      profile: {
        first_name: 'e2efirstName',
        last_name: 'e2elastName',
      },
      email: {
        email: 'e2e@email.ch',
      },
      phone: {
        phone: '+41 123456789',
      },
    },
    undefined,
    'userId',
  );
}

export function ensureMachineUserExists(api: API, username: string) {
  return ensureItemExists(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    `${api.mgmtBaseURL}/users/machine`,
    {
      user_name: username,
      name: 'e2emachinename',
      description: 'e2emachinedescription',
    },
    undefined,
    'userId',
  );
}

export function ensureUserDoesntExist(api: API, username: string) {
  return ensureItemDoesntExist(
    api,
    `${api.mgmtBaseURL}/users/_search`,
    (user: any) => user.userName === username,
    (user) => `${api.mgmtBaseURL}/users/${user.id}`,
  );
}
