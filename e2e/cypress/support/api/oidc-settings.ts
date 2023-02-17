import { ensureSetting } from './ensure';
import { API } from './types';

export function ensureOIDCSettingsSet(
  api: API,
  accessTokenLifetime: number,
  idTokenLifetime: number,
  refreshTokenExpiration: number,
  refreshTokenIdleExpiration: number,
) {
  return ensureSetting(
    api,
    `${api.adminBaseURL}/settings/oidc`,
    (body: any) => {
      const result = {
        sequence: body.settings?.details?.sequence,
        id: body.settings.id,
        entity: null,
      };

      if (
        body.settings &&
        body.settings.accessTokenLifetime === hoursToDuration(accessTokenLifetime) &&
        body.settings.idTokenLifetime === hoursToDuration(idTokenLifetime) &&
        body.settings.refreshTokenExpiration === daysToDuration(refreshTokenExpiration) &&
        body.settings.refreshTokenIdleExpiration === daysToDuration(refreshTokenIdleExpiration)
      ) {
        return { ...result, entity: body.settings };
      }
      return result;
    },
    `${api.adminBaseURL}/settings/oidc`,
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
