import { ensureSetting } from './ensure';
import { API } from './types';

export function ensureOIDCSettingsSet(
  api: API,
  accessTokenLifetime: number,
  idTokenLifetime: number,
  refreshTokenExpiration: number,
  refreshTokenIdleExpiration: number,
): Cypress.Chainable<number> {
  return ensureSetting(
    api,
    `${api.adminBaseURL}settings/oidc`,
    (body: any) => {
      const result = {
        entity: body.settings,
        sequence: body.settings?.details?.sequence,
        id: body.settings.id,
      };

      if (
        body.settings?.accessTokenLifetime != hoursToDuration(accessTokenLifetime) ||
        body.settings?.idTokenLifetime != hoursToDuration(idTokenLifetime) ||
        body.settings?.refreshTokenExpiration != daysToDuration(refreshTokenExpiration) ||
        body.settings?.refreshTokenIdleExpiration != daysToDuration(refreshTokenIdleExpiration)
      ) {
        result.entity = null;
      }
      return result;
    },
    `${api.adminBaseURL}settings/oidc`,
    {
      accessTokenLifetime: hoursToDuration(accessTokenLifetime),
      idTokenLifetime: hoursToDuration(idTokenLifetime),
      refreshTokenExpiration: daysToDuration(refreshTokenExpiration),
      refreshTokenIdleExpiration: daysToDuration(refreshTokenIdleExpiration),
    },
  );
}

function hoursToDuration(hours: number): string {
  return (hours * 3600).toString() + 's';
}
function daysToDuration(days: number): string {
  return hoursToDuration(24 * days);
}
