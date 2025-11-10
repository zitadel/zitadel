/**
 * URL validation utilities for security-critical redirect parameters
 * Provides defense-in-depth protection against open redirect attacks
 */

import { listTrustedDomains, listCustomDomains } from './zitadel';

export interface ValidationResult {
  valid: boolean;
  reason?: string;
}

/**
 * Validate redirect URI to prevent open redirect attacks
 * This is a defense-in-depth measure - backend validation is the primary control
 * 
 * @param uri - The URI to validate
 * @param options - Validation options
 * @returns Validation result with reason for failure
 */
export function validateRedirectUri(
  uri: string,
  options: {
    /**
     * Array of trusted domain strings
     * If provided, only these domains will be allowed
     * Supports subdomain matching (e.g., "example.com" allows "app.example.com")
     */
    trustedDomains?: string[];
    /**
     * Whether to enforce HTTPS in production
     * Defaults to true
     */
    enforceHttps?: boolean;
  } = {}
): ValidationResult {
  const { trustedDomains, enforceHttps = true } = options;

  // 1. Validate URL format and scheme
  let parsed: URL;
  try {
    parsed = new URL(uri);
  } catch {
    return { valid: false, reason: 'Invalid URL format' };
  }

  // 2. Only allow HTTP/HTTPS protocols
  // This blocks javascript:, data:, file:, vbscript:, etc.
  if (!['http:', 'https:'].includes(parsed.protocol)) {
    return { 
      valid: false, 
      reason: `Invalid protocol: ${parsed.protocol}. Only http: and https: are allowed` 
    };
  }

  // 3. Enforce HTTPS in production
  if (
    enforceHttps &&
    process.env.NODE_ENV === 'production' && 
    parsed.protocol !== 'https:'
  ) {
    return { 
      valid: false, 
      reason: 'HTTPS required in production environment' 
    };
  }

  // 4. Validate against trusted domains if configured
  if (trustedDomains && trustedDomains.length > 0) {
    const hostname = parsed.hostname.toLowerCase();
    
    const isAllowed = trustedDomains.some((domain: string) => {
      const normalizedDomain = domain.trim().toLowerCase();
      if (!normalizedDomain) return false;
      
      // Exact match
      if (hostname === normalizedDomain) {
        return true;
      }
      
      // Subdomain match (e.g., "app.example.com" matches "example.com")
      if (hostname.endsWith(`.${normalizedDomain}`)) {
        return true;
      }
      
      return false;
    });

    if (!isAllowed) {
      return {
        valid: false,
        reason: `Domain ${parsed.hostname} not in trusted domains list`
      };
    }
  }

  return { valid: true };
}

/**
 * Helper function to safely redirect with validation
 * Logs security events and redirects to safe default if validation fails
 * 
 * @param uri - The URI to redirect to
 * @param fallbackPath - Path to redirect to if validation fails
 * @param options - Validation options
 * @returns The validated URI or fallback path
 */
export async function getValidatedRedirectUri(
  uri: string,
  fallbackPath: string,
  options: {
    serviceUrl: string;
    enforceHttps?: boolean;
    /**
     * Additional parameters to add to fallback path
     */
    fallbackParams?: Record<string, string>;
  }
): Promise<string> {
  const { serviceUrl, enforceHttps, fallbackParams } = options;

  // Fetch trusted domains and custom domains from ZITADEL API
  const [trustedDomainsResponse, customDomainsResponse] = await Promise.all([
    listTrustedDomains({ serviceUrl }),
    listCustomDomains({ serviceUrl }),
  ]);

  const trustedDomains = [
    ...(trustedDomainsResponse?.map(td => td.domain) || []),
    ...(customDomainsResponse?.map(cd => cd.domain) || []),
  ];

  const validation = validateRedirectUri(uri, { 
    trustedDomains, 
    enforceHttps 
  });

  if (!validation.valid) {
    console.warn('[Security] Blocked invalid redirect URI:', {
      uri,
      reason: validation.reason,
      timestamp: new Date().toISOString(),
    });

    // Build fallback URL with parameters
    if (fallbackParams && Object.keys(fallbackParams).length > 0) {
      const params = new URLSearchParams(fallbackParams);
      return `${fallbackPath}?${params}`;
    }

    return fallbackPath;
  }

  return uri;
}
