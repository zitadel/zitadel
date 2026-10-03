export const VERSION_PREFERENCE_COOKIE = 'zitadel-docs-version';
export const LATEST_VERSION = 'latest';
export const LATEST_VERSION_LABEL = 'Latest (unreleased)';

export function getVersionPreferenceCookie(version) {
  return `${VERSION_PREFERENCE_COOKIE}=${encodeURIComponent(version)}; Path=/docs; Max-Age=31536000; SameSite=Lax`;
}

export function getLatestReleasedVersion(versions) {
  return versions.find((version) => version.param !== LATEST_VERSION)?.param;
}

export function resolveVersion({ requestedVersion, storedVersion, versions }) {
  const availableVersions = new Set(versions.map((version) => version.param));

  if (requestedVersion && availableVersions.has(requestedVersion)) {
    return requestedVersion;
  }

  if (storedVersion && availableVersions.has(storedVersion)) {
    return storedVersion;
  }

  return getLatestReleasedVersion(versions) ?? LATEST_VERSION;
}
