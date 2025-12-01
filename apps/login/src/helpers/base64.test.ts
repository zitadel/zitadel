import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { coerceToBase64Url, coerceToArrayBuffer } from "./base64";

const originalBtoa = global.btoa;
const originalAtob = global.atob;

beforeEach(() => {
  global.btoa = (str: string) => Buffer.from(str, "binary").toString("base64");
  global.atob = (b64: string) => Buffer.from(b64, "base64").toString("binary");
  global.window = global.window || ({} as any);
  (global.window as any).btoa = global.btoa;
  (global.window as any).atob = global.atob;
});

afterEach(() => {
  global.btoa = originalBtoa;
  global.atob = originalAtob;
});

describe("base64 utilities", () => {
  describe("coerceToBase64Url", () => {
    it("should convert Uint8Array to base64url", () => {
      const uint8Array = new Uint8Array([72, 101, 108, 108, 111]);
      const result = coerceToBase64Url(uint8Array, "test");

      expect(result).not.toContain("+");
      expect(result).not.toContain("/");
      expect(result).not.toContain("=");
      expect(typeof result).toBe("string");
    });

    it("should convert Array to base64url", () => {
      const array = [72, 101, 108, 108, 111];
      const result = coerceToBase64Url(array, "test");

      expect(result).not.toContain("+");
      expect(result).not.toContain("/");
      expect(result).not.toContain("=");
      expect(typeof result).toBe("string");
    });

    it("should convert ArrayBuffer to base64url", () => {
      const buffer = new Uint8Array([72, 101, 108, 108, 111]).buffer;
      const result = coerceToBase64Url(buffer, "test");

      expect(result).not.toContain("+");
      expect(result).not.toContain("/");
      expect(result).not.toContain("=");
      expect(typeof result).toBe("string");
    });

    it("should convert standard base64 string to base64url", () => {
      const base64 = "Hello+World/Test=";
      const result = coerceToBase64Url(base64, "test");

      expect(result).toBe("Hello-World_Test");
      expect(result).not.toContain("+");
      expect(result).not.toContain("/");
      expect(result).not.toContain("=");
    });

    it("should handle already URL-safe base64 strings", () => {
      const base64url = "Hello-World_Test";
      const result = coerceToBase64Url(base64url, "test");

      expect(result).toBe("Hello-World_Test");
    });

    it("should remove padding from base64 strings", () => {
      const base64WithPadding = "Test==";
      const result = coerceToBase64Url(base64WithPadding, "test");

      expect(result).toBe("Test");
      expect(result).not.toContain("=");
    });

    it("should handle empty Uint8Array", () => {
      const uint8Array = new Uint8Array([]);
      const result = coerceToBase64Url(uint8Array, "test");

      expect(typeof result).toBe("string");
      expect(result).toBe("");
    });

    it("should handle empty string", () => {
      const result = coerceToBase64Url("", "test");

      expect(result).toBe("");
    });

    it("should throw error for non-coercible types", () => {
      expect(() => coerceToBase64Url(123, "number")).toThrow("could not coerce 'number' to string");

      expect(() => coerceToBase64Url(null, "null")).toThrow("could not coerce 'null' to string");

      expect(() => coerceToBase64Url(undefined, "undefined")).toThrow("could not coerce 'undefined' to string");

      expect(() => coerceToBase64Url({}, "object")).toThrow("could not coerce 'object' to string");
    });

    it("should handle binary data correctly", () => {
      const binaryData = new Uint8Array([0, 1, 2, 255, 254, 253]);
      const result = coerceToBase64Url(binaryData, "binary");

      expect(typeof result).toBe("string");
      expect(result.length).toBeGreaterThan(0);
    });

    it("should handle large Uint8Arrays", () => {
      const largeArray = new Uint8Array(1000);
      for (let i = 0; i < 1000; i++) {
        largeArray[i] = i % 256;
      }

      const result = coerceToBase64Url(largeArray, "large");

      expect(typeof result).toBe("string");
      expect(result.length).toBeGreaterThan(0);
    });

    it("should produce consistent results for same input", () => {
      const input = new Uint8Array([1, 2, 3, 4, 5]);
      const result1 = coerceToBase64Url(input, "test");
      const result2 = coerceToBase64Url(input, "test");

      expect(result1).toBe(result2);
    });
  });

  describe("coerceToArrayBuffer", () => {
    it("should convert base64url string to ArrayBuffer", () => {
      const base64url = "SGVsbG8";
      const result = coerceToArrayBuffer(base64url, "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result.byteLength).toBeGreaterThan(0);
    });

    it("should convert base64 string with URL-safe characters to ArrayBuffer", () => {
      const base64url = "Hello-World_Test";
      const result = coerceToArrayBuffer(base64url, "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
    });

    it("should convert Array to ArrayBuffer", () => {
      const array = [72, 101, 108, 108, 111];
      const result = coerceToArrayBuffer(array, "test");

      expect(result).toBeInstanceOf(ArrayBuffer);

      const view = new Uint8Array(result);
      expect(Array.from(view)).toEqual(array);
    });

    it("should convert Uint8Array to ArrayBuffer", () => {
      const uint8Array = new Uint8Array([72, 101, 108, 108, 111]);
      const result = coerceToArrayBuffer(uint8Array, "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result).toBe(uint8Array.buffer);
    });

    it("should handle ArrayBuffer input (passthrough)", () => {
      const buffer = new Uint8Array([1, 2, 3]).buffer;
      const result = coerceToArrayBuffer(buffer, "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result).toBe(buffer);
    });

    it("should handle empty string", () => {
      const result = coerceToArrayBuffer("", "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result.byteLength).toBe(0);
    });

    it("should handle empty array", () => {
      const result = coerceToArrayBuffer([], "test");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result.byteLength).toBe(0);
    });

    it("should throw TypeError for non-coercible types", () => {
      expect(() => coerceToArrayBuffer(123, "number")).toThrow(TypeError);
      expect(() => coerceToArrayBuffer(123, "number")).toThrow("could not coerce 'number' to ArrayBuffer");

      expect(() => coerceToArrayBuffer(null, "null")).toThrow(TypeError);
      expect(() => coerceToArrayBuffer(undefined, "undefined")).toThrow(TypeError);
      expect(() => coerceToArrayBuffer({}, "object")).toThrow(TypeError);
    });

    it("should handle binary data correctly", () => {
      const binaryArray = [0, 1, 2, 255, 254, 253];
      const result = coerceToArrayBuffer(binaryArray, "binary");

      expect(result).toBeInstanceOf(ArrayBuffer);

      const view = new Uint8Array(result);
      expect(Array.from(view)).toEqual(binaryArray);
    });

    it("should handle large arrays", () => {
      const largeArray = new Array(1000).fill(0).map((_, i) => i % 256);
      const result = coerceToArrayBuffer(largeArray, "large");

      expect(result).toBeInstanceOf(ArrayBuffer);
      expect(result.byteLength).toBe(1000);
    });

    it("should produce consistent results for same input", () => {
      const input = [1, 2, 3, 4, 5];
      const result1 = coerceToArrayBuffer(input, "test");
      const result2 = coerceToArrayBuffer(input, "test");

      const view1 = new Uint8Array(result1);
      const view2 = new Uint8Array(result2);

      expect(Array.from(view1)).toEqual(Array.from(view2));
    });
  });

  describe("round-trip conversions", () => {
    it("should successfully round-trip Uint8Array -> base64url -> ArrayBuffer", () => {
      const original = new Uint8Array([1, 2, 3, 4, 5, 255, 254, 253]);

      const base64url = coerceToBase64Url(original, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");

      const recoveredArray = new Uint8Array(recovered);
      expect(Array.from(recoveredArray)).toEqual(Array.from(original));
    });

    it("should successfully round-trip Array -> base64url -> ArrayBuffer -> Array", () => {
      const original = [10, 20, 30, 40, 50];

      const base64url = coerceToBase64Url(original, "test");
      const buffer = coerceToArrayBuffer(base64url, "test");
      const recovered = Array.from(new Uint8Array(buffer));

      expect(recovered).toEqual(original);
    });

    it("should handle empty data in round-trip", () => {
      const original = new Uint8Array([]);

      const base64url = coerceToBase64Url(original, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");

      expect(recovered.byteLength).toBe(0);
    });

    it("should handle special characters in round-trip", () => {
      const original = new Uint8Array([0xff, 0xff, 0xff, 0xfe, 0x00, 0x01]);

      const base64url = coerceToBase64Url(original, "test");

      expect(base64url).not.toContain("+");
      expect(base64url).not.toContain("/");
      expect(base64url).not.toContain("=");

      const recovered = coerceToArrayBuffer(base64url, "test");
      const recoveredArray = new Uint8Array(recovered);

      expect(Array.from(recoveredArray)).toEqual(Array.from(original));
    });

    it("should handle large data in round-trip", () => {
      const original = new Uint8Array(1024);
      for (let i = 0; i < 1024; i++) {
        original[i] = Math.floor(Math.random() * 256);
      }

      const base64url = coerceToBase64Url(original, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");
      const recoveredArray = new Uint8Array(recovered);

      expect(Array.from(recoveredArray)).toEqual(Array.from(original));
    });
  });

  describe("edge cases", () => {
    it("should handle all zero bytes", () => {
      const zeros = new Uint8Array([0, 0, 0, 0]);
      const base64url = coerceToBase64Url(zeros, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");

      expect(Array.from(new Uint8Array(recovered))).toEqual([0, 0, 0, 0]);
    });

    it("should handle all max bytes", () => {
      const maxBytes = new Uint8Array([255, 255, 255, 255]);
      const base64url = coerceToBase64Url(maxBytes, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");

      expect(Array.from(new Uint8Array(recovered))).toEqual([255, 255, 255, 255]);
    });

    it("should handle single byte", () => {
      const singleByte = new Uint8Array([42]);
      const base64url = coerceToBase64Url(singleByte, "test");
      const recovered = coerceToArrayBuffer(base64url, "test");

      expect(Array.from(new Uint8Array(recovered))).toEqual([42]);
    });

    it("should provide meaningful error messages with parameter names", () => {
      expect(() => coerceToBase64Url(123, "myParameter")).toThrow("could not coerce 'myParameter' to string");

      expect(() => coerceToArrayBuffer(123, "anotherParam")).toThrow("could not coerce 'anotherParam' to ArrayBuffer");
    });
  });
});
