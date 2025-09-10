import { API } from './types';

export function getInstanceFeatures(api: API) {
  return cy.request({
    method: 'GET',
    url: `${api.featuresBaseURL}/instance`,
    headers: {
      authorization: `Bearer ${api.token}`,
    },
    body: {},
  });
}

export function setInstanceFeature(api: API, feature: string, enabled: boolean, additionalConfig?: Record<string, any>) {
  const body: Record<string, any> = {
    [feature]: enabled,
    ...additionalConfig,
  };

  return cy.request({
    method: 'PUT',
    url: `${api.featuresBaseURL}/instance`,
    headers: {
      authorization: `Bearer ${api.token}`,
    },
    body,
  });
}

export function resetInstanceFeatures(api: API) {
  return cy.request({
    method: 'DELETE',
    url: `${api.featuresBaseURL}/instance`,
    headers: {
      authorization: `Bearer ${api.token}`,
    },
  });
}

export function ensureFeatureState(api: API, feature: string, enabled: boolean, additionalConfig?: Record<string, any>) {
  return getInstanceFeatures(api).then((response) => {
    const currentState = response.body?.[feature]?.enabled;

    if (currentState !== enabled) {
      return setInstanceFeature(api, feature, enabled, additionalConfig);
    }

    return cy.wrap(response);
  });
}

export function ensureLoginV2FeatureState(api: API, required: boolean, baseUri?: string) {
  return setInstanceFeature(api, 'loginV2', required, {
    loginV2: {
      required,
      baseUri: baseUri || '',
    },
  });
}
