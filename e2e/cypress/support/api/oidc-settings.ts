import { ZITADELTarget } from 'support/commands';

export function ensureOIDCSettings(
  target: ZITADELTarget,
  accessTokenLifetime: number,
  idTokenLifetime: number,
  refreshTokenExpiration: number,
  refreshTokenIdleExpiration: number,
) {
  updateOIDCSettings(target, accessTokenLifetime, idTokenLifetime, refreshTokenExpiration, refreshTokenIdleExpiration);
  return getOIDCSettings(target).should(
    (res) =>
      res.body.accessTokenLifetime == accessTokenLifetime &&
      res.body.idTokenLifetime == idTokenLifetime &&
      res.body.refreshTokenExpiration == refreshTokenExpiration &&
      res.body.refreshTokenIdleExpiration == refreshTokenIdleExpiration,
  );
}

function getOIDCSettings(target: ZITADELTarget) {
  return cy.request({
    method: 'GET',
    url: `${target.adminBaseURL}/settings/oidc`,
    headers: target.headers,
  });
}

function updateOIDCSettings(
  target: ZITADELTarget,
  accessTokenLifetime: number,
  idTokenLifetime: number,
  refreshTokenExpiration: number,
  refreshTokenIdleExpiration: number,
) {
  return cy
    .request({
      method: 'PUT',
      url: `${target.adminBaseURL}/settings/oidc`,
      body: {
        accessTokenLifetime: hoursToDuration(accessTokenLifetime),
        idTokenLifetime: hoursToDuration(idTokenLifetime),
        refreshTokenExpiration: daysToDuration(refreshTokenExpiration),
        refreshTokenIdleExpiration: daysToDuration(refreshTokenIdleExpiration),
      },
      headers: target.headers,
      failOnStatusCode: false,
    })
    .then((res) => {
      if (!res.isOkStatusCode) {
        expect(res.status).to.equal(400);
        expect(res.body.message).to.contain('No changes');
      }
      return res;
    });
}

function hoursToDuration(hours: number): string {
  return (hours * 3600).toString() + 's';
}
function daysToDuration(days: number): string {
  return hoursToDuration(24 * days);
}
