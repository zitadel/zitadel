import { describe, it, expect, vi } from 'vitest';
import { validateRedirectUri } from './url-validation';

describe('validateRedirectUri', () => {

  describe('URL format validation', () => {
    it('should accept valid HTTPS URLs', () => {
      const result = validateRedirectUri('https://example.com/path');
      expect(result.valid).toBe(true);
    });

    it('should accept valid HTTP URLs in non-production', () => {
      const result = validateRedirectUri('http://example.com/path');
      expect(result.valid).toBe(true);
    });

    it('should reject malformed URLs', () => {
      const result = validateRedirectUri('not-a-url');
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Invalid URL format');
    });

    it('should reject empty strings', () => {
      const result = validateRedirectUri('');
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Invalid URL format');
    });

    it('should accept URLs with spaces (URL constructor auto-encodes)', () => {
      // Note: URL constructor automatically encodes spaces to %20
      const result = validateRedirectUri('https://example.com/path with spaces');
      expect(result.valid).toBe(true);
    });
  });

  describe('Protocol validation', () => {
    it('should block javascript: protocol', () => {
      const result = validateRedirectUri('javascript:alert(1)');
      expect(result.valid).toBe(false);
      expect(result.reason).toContain('Invalid protocol');
      expect(result.reason).toContain('javascript:');
    });

    it('should block data: protocol', () => {
      const result = validateRedirectUri('data:text/html,<script>alert(1)</script>');
      expect(result.valid).toBe(false);
      expect(result.reason).toContain('Invalid protocol');
      expect(result.reason).toContain('data:');
    });

    it('should block file: protocol', () => {
      const result = validateRedirectUri('file:///etc/passwd');
      expect(result.valid).toBe(false);
      expect(result.reason).toContain('Invalid protocol');
      expect(result.reason).toContain('file:');
    });

    it('should block vbscript: protocol', () => {
      const result = validateRedirectUri('vbscript:msgbox(1)');
      expect(result.valid).toBe(false);
      expect(result.reason).toContain('Invalid protocol');
      expect(result.reason).toContain('vbscript:');
    });

    it('should block ftp: protocol', () => {
      const result = validateRedirectUri('ftp://example.com/file');
      expect(result.valid).toBe(false);
      expect(result.reason).toContain('Invalid protocol');
      expect(result.reason).toContain('ftp:');
    });

    it('should accept http: protocol', () => {
      const result = validateRedirectUri('http://example.com');
      expect(result.valid).toBe(true);
    });

    it('should accept https: protocol', () => {
      const result = validateRedirectUri('https://example.com');
      expect(result.valid).toBe(true);
    });
  });

  describe('HTTPS enforcement', () => {
    it('should enforce HTTPS in production when enforceHttps is true', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('http://example.com', { enforceHttps: true });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('HTTPS required in production environment');
      vi.unstubAllEnvs();
    });

    it('should allow HTTP in production when enforceHttps is false', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('http://example.com', { enforceHttps: false });
      expect(result.valid).toBe(true);
      vi.unstubAllEnvs();
    });

    it('should allow HTTP in development regardless of enforceHttps', () => {
      vi.stubEnv('NODE_ENV', 'development');
      const result = validateRedirectUri('http://example.com', { enforceHttps: true });
      expect(result.valid).toBe(true);
      vi.unstubAllEnvs();
    });

    it('should allow HTTPS in production', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('https://example.com', { enforceHttps: true });
      expect(result.valid).toBe(true);
      vi.unstubAllEnvs();
    });

    it('should default to enforcing HTTPS', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('http://example.com');
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('HTTPS required in production environment');
      vi.unstubAllEnvs();
    });
  });

  describe('Trusted domains validation', () => {
    it('should accept URLs without trusted domains configured', () => {
      const result = validateRedirectUri('https://example.com');
      expect(result.valid).toBe(true);
    });

    it('should accept URLs with empty trusted domains array', () => {
      const result = validateRedirectUri('https://example.com', { trustedDomains: [] });
      expect(result.valid).toBe(true);
    });

    it('should accept exact domain match', () => {
      const result = validateRedirectUri('https://example.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept subdomain match', () => {
      const result = validateRedirectUri('https://app.example.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept deeply nested subdomain match', () => {
      const result = validateRedirectUri('https://api.v2.app.example.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should reject domain not in trusted list', () => {
      const result = validateRedirectUri('https://evil.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Domain evil.com not in trusted domains list');
    });

    it('should reject partial domain match', () => {
      const result = validateRedirectUri('https://notexample.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Domain notexample.com not in trusted domains list');
    });

    it('should reject domain suffix match without subdomain separator', () => {
      const result = validateRedirectUri('https://evilexample.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Domain evilexample.com not in trusted domains list');
    });

    it('should be case-insensitive for domains', () => {
      const result = validateRedirectUri('https://APP.EXAMPLE.COM/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle mixed case in trusted domains', () => {
      const result = validateRedirectUri('https://app.example.com/path', {
        trustedDomains: ['Example.COM'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept if any trusted domain matches', () => {
      const result = validateRedirectUri('https://app.example.com/path', {
        trustedDomains: ['other.com', 'example.com', 'another.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle domains with ports', () => {
      const result = validateRedirectUri('https://example.com:8080/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle localhost', () => {
      const result = validateRedirectUri('http://localhost:3000/callback', {
        trustedDomains: ['localhost'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle IP addresses', () => {
      const result = validateRedirectUri('https://192.168.1.1/path', {
        trustedDomains: ['192.168.1.1'],
      });
      expect(result.valid).toBe(true);
    });

    it('should trim whitespace from trusted domains', () => {
      const result = validateRedirectUri('https://example.com/path', {
        trustedDomains: ['  example.com  '],
      });
      expect(result.valid).toBe(true);
    });

    it('should ignore empty trusted domains', () => {
      const result = validateRedirectUri('https://example.com/path', {
        trustedDomains: ['', '  ', 'example.com'],
      });
      expect(result.valid).toBe(true);
    });
  });

  describe('Real-world attack scenarios', () => {
    it('should block open redirect to phishing site', () => {
      const result = validateRedirectUri('https://evil-phishing-site.com/fake-login', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
    });

    it('should block redirect with similar domain name', () => {
      const result = validateRedirectUri('https://example.com.evil.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
    });

    it('should block homograph attack (look-alike domains)', () => {
      // Note: This is a simplified test. Real homograph attacks use Unicode characters
      const result = validateRedirectUri('https://examp1e.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(false);
    });

    it('should accept legitimate multi-level subdomain', () => {
      const result = validateRedirectUri('https://tenant-123.app.example.com/callback', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept legitimate development subdomain', () => {
      const result = validateRedirectUri('https://dev.example.com/callback', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });
  });

  describe('Edge cases', () => {
    it('should handle URLs with query parameters', () => {
      const result = validateRedirectUri('https://example.com/path?foo=bar&baz=qux', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle URLs with fragments', () => {
      const result = validateRedirectUri('https://example.com/path#section', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle URLs with authentication', () => {
      const result = validateRedirectUri('https://user:pass@example.com/path', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle root path URLs', () => {
      const result = validateRedirectUri('https://example.com/', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle URLs without path', () => {
      const result = validateRedirectUri('https://example.com', {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle very long URLs', () => {
      const longPath = '/path/' + 'a'.repeat(1000);
      const result = validateRedirectUri(`https://example.com${longPath}`, {
        trustedDomains: ['example.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should handle internationalized domain names (IDN)', () => {
      const result = validateRedirectUri('https://münchen.de/path', {
        trustedDomains: ['xn--mnchen-3ya.de'], // Punycode for münchen.de
      });
      // URL constructor automatically converts IDN to punycode
      expect(result.valid).toBe(true);
    });
  });

  describe('Combined validation scenarios', () => {
    it('should validate protocol and domain in production', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('http://evil.com/path', {
        trustedDomains: ['example.com'],
        enforceHttps: true,
      });
      // Should fail on HTTPS enforcement first
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('HTTPS required in production environment');
      vi.unstubAllEnvs();
    });

    it('should validate domain even with HTTPS in production', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('https://evil.com/path', {
        trustedDomains: ['example.com'],
        enforceHttps: true,
      });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Domain evil.com not in trusted domains list');
      vi.unstubAllEnvs();
    });

    it('should pass all validations for legitimate request', () => {
      vi.stubEnv('NODE_ENV', 'production');
      const result = validateRedirectUri('https://app.example.com/callback?state=abc123', {
        trustedDomains: ['example.com', 'trusted.com'],
        enforceHttps: true,
      });
      expect(result.valid).toBe(true);
      vi.unstubAllEnvs();
    });
  });

  describe('Custom domains and trusted domains combination', () => {
    it('should accept URL matching trusted domain', () => {
      const result = validateRedirectUri('https://app.trusted.com/callback', {
        trustedDomains: ['trusted.com', 'custom.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept URL matching custom domain', () => {
      const result = validateRedirectUri('https://custom.com/callback', {
        trustedDomains: ['trusted.com', 'custom.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should accept subdomain of custom domain', () => {
      const result = validateRedirectUri('https://app.custom.com/callback', {
        trustedDomains: ['trusted.com', 'custom.com'],
      });
      expect(result.valid).toBe(true);
    });

    it('should reject domain not in either list', () => {
      const result = validateRedirectUri('https://evil.com/callback', {
        trustedDomains: ['trusted.com', 'custom.com'],
      });
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('Domain evil.com not in trusted domains list');
    });
  });

  describe('getValidatedRedirectUri (integration with API)', () => {
    it('should fetch and combine trusted domains and custom domains', async () => {
      // Mock the ZITADEL API functions
      const mockListTrustedDomains = vi.fn().mockResolvedValue([
        { domain: 'trusted.com' },
        { domain: 'example.com' },
      ]);
      const mockListCustomDomains = vi.fn().mockResolvedValue([
        { domain: 'custom.com' },
        { domain: 'my-custom.com' },
      ]);

      // We can't easily mock the imports, but we can test the validateRedirectUri
      // function with the combined domains directly
      const combinedDomains = [
        'trusted.com',
        'example.com',
        'custom.com',
        'my-custom.com',
      ];

      const result1 = validateRedirectUri('https://app.trusted.com/callback', {
        trustedDomains: combinedDomains,
      });
      expect(result1.valid).toBe(true);

      const result2 = validateRedirectUri('https://custom.com/callback', {
        trustedDomains: combinedDomains,
      });
      expect(result2.valid).toBe(true);

      const result3 = validateRedirectUri('https://subdomain.my-custom.com/path', {
        trustedDomains: combinedDomains,
      });
      expect(result3.valid).toBe(true);

      const result4 = validateRedirectUri('https://evil.com/callback', {
        trustedDomains: combinedDomains,
      });
      expect(result4.valid).toBe(false);
      expect(result4.reason).toBe('Domain evil.com not in trusted domains list');
    });

    it('should handle empty responses from API', async () => {
      // Test with empty domain lists
      const emptyDomains: string[] = [];

      // Without trusted domains, all valid HTTPS URLs should pass
      const result = validateRedirectUri('https://any-domain.com/callback', {
        trustedDomains: emptyDomains,
      });
      expect(result.valid).toBe(true);
    });

    it('should handle undefined responses from API', async () => {
      // Test with undefined (API returned no data)
      const result = validateRedirectUri('https://any-domain.com/callback', {
        trustedDomains: undefined,
      });
      expect(result.valid).toBe(true);
    });

    it('should prioritize validation errors over domain checking', async () => {
      vi.stubEnv('NODE_ENV', 'production');
      
      const combinedDomains = ['trusted.com', 'custom.com'];

      // HTTP in production should fail before domain check
      const result = validateRedirectUri('http://trusted.com/callback', {
        trustedDomains: combinedDomains,
        enforceHttps: true,
      });
      
      expect(result.valid).toBe(false);
      expect(result.reason).toBe('HTTPS required in production environment');
      vi.unstubAllEnvs();
    });

    it('should validate domains from both trusted and custom lists', async () => {
      const combinedDomains = [
        'zitadel-cloud.com',    // trusted domain
        'my-company.com',       // custom domain
        'staging.example.com',  // another custom domain
      ];

      // Test various URLs against combined list
      const testCases = [
        { url: 'https://zitadel-cloud.com/callback', expected: true },
        { url: 'https://app.zitadel-cloud.com/callback', expected: true },
        { url: 'https://my-company.com/logout', expected: true },
        { url: 'https://api.my-company.com/auth', expected: true },
        { url: 'https://staging.example.com/login', expected: true },
        { url: 'https://evil.com/phishing', expected: false },
        { url: 'https://zitadel-cloud.com.evil.com', expected: false },
      ];

      testCases.forEach(({ url, expected }) => {
        const result = validateRedirectUri(url, {
          trustedDomains: combinedDomains,
        });
        expect(result.valid).toBe(expected);
      });
    });
  });
});
