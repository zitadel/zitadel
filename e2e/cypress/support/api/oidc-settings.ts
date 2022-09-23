import {apiCallProperties} from './apiauth';
import {ensureSomethingIsSet} from './ensure';

export function ensureOIDCSettingsSet(api: apiCallProperties, accessTokenLifetime, idTokenLifetime, refreshTokenIdleExpiration, refreshTokenExpiration: number): Cypress.Chainable<number> {
    return ensureSomethingIsSet(api, `${api.adminBaseURL}settings/oidc`,
        (settings: any) => {
            let entity = null;
            if (settings.settings?.accessTokenLifetime === durationString(accessTokenLifetime) &&
                settings.settings?.idTokenLifetime === durationString(idTokenLifetime) &&
                settings.settings?.refreshTokenIdleExpiration === durationString(refreshTokenIdleExpiration) &&
                settings.settings?.refreshTokenExpiration === durationString(refreshTokenExpiration)) {
                entity = settings.settings
            }
            return {
                entity: entity,
                sequence: settings.settings?.details?.sequence,
            };
        },
        `${api.adminBaseURL}settings/oidc`,
        {
            accessTokenLifetime: durationString(accessTokenLifetime),
            idTokenLifetime: durationString(idTokenLifetime),
            refreshTokenIdleExpiration: durationString(refreshTokenIdleExpiration),
            refreshTokenExpiration: durationString(refreshTokenExpiration),
        });
}

function durationString(durationInSeconds: number): string {
    return durationInSeconds.toString() + "s"
}
