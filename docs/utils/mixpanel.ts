import type { OverridedMixpanel } from 'mixpanel-browser';

let mixpanel: OverridedMixpanel | undefined;

export async function initMixpanel() {
  if (typeof window === 'undefined') return;

  const token = process.env.NEXT_PUBLIC_MIXPANEL_TOKEN;
  if (!token) return;

  const module = await import('mixpanel-browser');
  mixpanel = module.default;

  mixpanel.init(token, {
    property_blacklist: ['$referrer', 'referrer', '$current_url_query_params', '$initial_referrer'],
    debug: process.env.NODE_ENV === 'development',
    track_pageview: 'url-with-path',
    persistence: 'localStorage',
    api_host: 'https://api-eu.mixpanel.com',
    record_sessions_percent: 50,
  });
}

export function optInTracking() {
  if (isMixpanelInitialized()) {
    mixpanel?.opt_in_tracking();
  }
}

export function optOutTracking() {
  if (isMixpanelInitialized()) {
    mixpanel?.opt_out_tracking();
  }
}

function isMixpanelInitialized(): boolean {
  return typeof window !== 'undefined' && typeof mixpanel !== 'undefined';
}

export const mixpanelClient = {
  track: (eventName: string, properties?: Record<string, any>) => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.track(eventName, properties);
      } catch (e) {
        console.warn('Mixpanel tracking error:', e);
      }
    } else {
      if (process.env.NODE_ENV === 'development') {
        console.warn('[Mixpanel] Not initialized yet, skipping event:', eventName);
      }
    }
  },

  identify: (userId: string) => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.identify(userId);
      } catch (e) {
        console.warn('Mixpanel identify error:', e);
      }
    }
  },

  setPeople: (properties: Record<string, any>) => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.people.set(properties);
      } catch (e) {
        console.warn('Mixpanel people error:', e);
      }
    }
  },

  setUserProperties: (properties: Record<string, any>) => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.people.set(properties);
      } catch (e) {
        console.warn('Mixpanel set properties error:', e);
      }
    }
  },

  alias: (newId: string) => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.alias(newId);
      } catch (e) {
        console.warn('Mixpanel alias error:', e);
      }
    }
  },

  reset: () => {
    if (isMixpanelInitialized()) {
      try {
        mixpanel?.reset();
      } catch (e) {
        console.warn('Mixpanel reset error:', e);
      }
    }
  },
};
