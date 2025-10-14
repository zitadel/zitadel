import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API, Entity } from './types';

export function ensureSMTPProviderExists(api: API, providerDescription: string) {
  return ensureItemExists(
    api,
    `${api.adminBaseURL}/smtp/_search`,
    (provider: any) => {
      return provider.description === providerDescription;
    },
    `${api.adminBaseURL}/smtp`,
    {
      name: providerDescription,
      description: providerDescription,
      senderAddress: 'a@sender.com',
      senderName: 'A Sender',
      host: 'smtp.host.com:587',
      user: 'smtpuser',
    },
  );
}

export function activateSMTPProvider(api: API, providerId: string) {
  return cy.request({
    method: 'POST',
    url: `${api.adminBaseURL}/smtp/${providerId}/_activate`,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${api.token}`,
    },
  });
}
