import { apiCallProperties } from './apiauth';
import { ensureSomethingIsSet } from './ensure';

export function ensureOIDCSettingsSet(
  api: apiCallProperties,
  accessTokenLifetime: number,
  idTokenLifetime: number,
  refreshTokenExpiration: number,
  refreshTokenIdleExpiration: number,
): Cypress.Chainable<number> {
  return ensureSomethingIsSet(
    api,
    `${api.adminBaseURL}settings/oidc`,
    (settings: any) => {
      let entity = null;
      if (
        settings.settings?.accessTokenLifetime === hoursToDuration(accessTokenLifetime) &&
        settings.settings?.idTokenLifetime === hoursToDuration(idTokenLifetime) &&
        settings.settings?.refreshTokenExpiration === daysToDuration(refreshTokenExpiration) &&
        settings.settings?.refreshTokenIdleExpiration === daysToDuration(refreshTokenIdleExpiration)
      ) {
        entity = settings.settings;
      }
      return {
        entity: entity,
        sequence: settings.settings?.details?.sequence,
      };
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
