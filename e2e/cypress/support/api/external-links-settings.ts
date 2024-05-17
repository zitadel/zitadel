import { ensureSetting } from './ensure';
import { API } from './types';

export function ensureExternalLinksSettingsSet(api: API, tosLink: string, privacyPolicyLink: string, docsLink: string) {
  return ensureSetting(
    api,
    `${api.adminBaseURL}/policies/privacy`,
    (body: any) => {
      const result = {
        sequence: body.policy?.details?.sequence,
        id: body.policy.id,
        entity: null,
      };

      if (
        body.policy &&
        body.policy.tosLink === tosLink &&
        body.policy.privacyLink === privacyPolicyLink &&
        body.policy.docsLink === docsLink
      ) {
        return { ...result, entity: body.policy };
      }
      return result;
    },
    `${api.adminBaseURL}/policies/privacy`,
    {
      tosLink,
      privacyLink: privacyPolicyLink,
      docsLink,
    },
  );
}
