import { ensureItemDoesntExist, ensureItemExists } from './ensure';
import { API, Entity } from './types';
import { ensureSMTPProviderExists } from './smtp';

export function ensureSMSProviderExists(api: API) {
  // remove and create
  ensureSMSProviderDoesntExist(api);
  return ensureItemExists(
    api,
    `${api.adminBaseURL}/sms/_search`,
    ({ twilio: { sid: foundSid } }: any) => foundSid === 'initial-sid',
    `${api.adminBaseURL}/sms/twilio`,
    {
      sid: 'initial-sid',
      senderNumber: 'initial-senderNumber',
      token: 'initial-token',
    },
  );
}

export function ensureSMSProviderDoesntExist(api: API) {
  return ensureItemDoesntExist(
    api,
    `${api.adminBaseURL}/sms/_search`,
    (provider: any) => !!provider,
    (provider) => `${api.adminBaseURL}/sms/${provider.id}`,
  );
}
