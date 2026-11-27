import { RedirectPipe } from './redirect.pipe';

describe('RedirectPipe', () => {
  const pipe = new RedirectPipe();

  describe('native applications', () => {
    it('accepts https redirect URIs', () => {
      expect(pipe.transform('https://example.com/callback', true)).toBeTrue();
    });

    it('accepts custom scheme redirect URIs', () => {
      expect(pipe.transform('myapp://callback', true)).toBeTrue();
      expect(pipe.transform('com.example.app:/callback', true)).toBeTrue();
    });

    it('accepts http loopback redirect URIs', () => {
      expect(pipe.transform('http://localhost:3000/callback', true)).toBeTrue();
      expect(pipe.transform('http://127.0.0.1/callback', true)).toBeTrue();
      expect(pipe.transform('http://[::1]:8080/callback', true)).toBeTrue();
    });

    it('rejects http redirect URIs that are not loopback', () => {
      expect(pipe.transform('http://example.com/callback', true)).toBeFalse();
    });
  });

  describe('web and user agent applications', () => {
    it('accepts https redirect URIs', () => {
      expect(pipe.transform('https://example.com/callback', false)).toBeTrue();
    });

    it('rejects http and custom scheme redirect URIs', () => {
      expect(pipe.transform('http://localhost:3000/callback', false)).toBeFalse();
      expect(pipe.transform('myapp://callback', false)).toBeFalse();
    });
  });

  it('rejects invalid URIs', () => {
    expect(pipe.transform('not a uri', true)).toBeFalse();
    expect(pipe.transform('not a uri', false)).toBeFalse();
  });
});
