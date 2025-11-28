import { describe, it, expect } from "vitest";
import { symbolValidator, numberValidator, upperCaseValidator, lowerCaseValidator } from "./validators";

describe("validators", () => {
  describe("symbolValidator", () => {
    it("should return true for strings with symbols", () => {
      expect(symbolValidator("!")).toBe(true);
      expect(symbolValidator("@")).toBe(true);
      expect(symbolValidator("#")).toBe(true);
      expect(symbolValidator("$")).toBe(true);
      expect(symbolValidator("%")).toBe(true);
      expect(symbolValidator("^")).toBe(true);
      expect(symbolValidator("&")).toBe(true);
      expect(symbolValidator("*")).toBe(true);
      expect(symbolValidator("(")).toBe(true);
      expect(symbolValidator(")")).toBe(true);
    });

    it("should return true for strings with punctuation", () => {
      expect(symbolValidator(".")).toBe(true);
      expect(symbolValidator(",")).toBe(true);
      expect(symbolValidator(";")).toBe(true);
      expect(symbolValidator(":")).toBe(true);
      expect(symbolValidator("?")).toBe(true);
      expect(symbolValidator("'")).toBe(true);
      expect(symbolValidator('"')).toBe(true);
    });

    it("should return true for strings with special characters", () => {
      expect(symbolValidator("-")).toBe(true);
      expect(symbolValidator("_")).toBe(true);
      expect(symbolValidator("=")).toBe(true);
      expect(symbolValidator("+")).toBe(true);
      expect(symbolValidator("[")).toBe(true);
      expect(symbolValidator("]")).toBe(true);
      expect(symbolValidator("{")).toBe(true);
      expect(symbolValidator("}")).toBe(true);
      expect(symbolValidator("|")).toBe(true);
      expect(symbolValidator("\\")).toBe(true);
      expect(symbolValidator("/")).toBe(true);
    });

    it("should return true for strings with whitespace", () => {
      expect(symbolValidator(" ")).toBe(true);
      expect(symbolValidator("\t")).toBe(true);
      expect(symbolValidator("\n")).toBe(true);
    });

    it("should return false for strings with only letters", () => {
      expect(symbolValidator("abc")).toBe(false);
      expect(symbolValidator("ABC")).toBe(false);
      expect(symbolValidator("aBc")).toBe(false);
    });

    it("should return false for strings with only numbers", () => {
      expect(symbolValidator("123")).toBe(false);
      expect(symbolValidator("0")).toBe(false);
    });

    it("should return false for alphanumeric strings", () => {
      expect(symbolValidator("abc123")).toBe(false);
      expect(symbolValidator("ABC123")).toBe(false);
      expect(symbolValidator("Test123")).toBe(false);
    });

    it("should return true for mixed strings with symbols", () => {
      expect(symbolValidator("Password!")).toBe(true);
      expect(symbolValidator("Test@123")).toBe(true);
      expect(symbolValidator("hello-world")).toBe(true);
      expect(symbolValidator("user_name")).toBe(true);
    });

    it("should return false for empty strings", () => {
      expect(symbolValidator("")).toBe(false);
    });

    it("should return true for Unicode symbols", () => {
      expect(symbolValidator("€")).toBe(true);
      expect(symbolValidator("©")).toBe(true);
      expect(symbolValidator("™")).toBe(true);
    });
  });

  describe("numberValidator", () => {
    it("should return true for strings with single digits", () => {
      expect(numberValidator("0")).toBe(true);
      expect(numberValidator("1")).toBe(true);
      expect(numberValidator("5")).toBe(true);
      expect(numberValidator("9")).toBe(true);
    });

    it("should return true for strings with multiple digits", () => {
      expect(numberValidator("123")).toBe(true);
      expect(numberValidator("000")).toBe(true);
      expect(numberValidator("999")).toBe(true);
    });

    it("should return false for strings without numbers", () => {
      expect(numberValidator("abc")).toBe(false);
      expect(numberValidator("ABC")).toBe(false);
      expect(numberValidator("test")).toBe(false);
    });

    it("should return true for alphanumeric strings", () => {
      expect(numberValidator("abc123")).toBe(true);
      expect(numberValidator("Test1")).toBe(true);
      expect(numberValidator("1test")).toBe(true);
    });

    it("should return false for strings with only symbols", () => {
      expect(numberValidator("!@#$%")).toBe(false);
      expect(numberValidator("***")).toBe(false);
    });

    it("should return true for mixed content with numbers", () => {
      expect(numberValidator("Password1")).toBe(true);
      expect(numberValidator("Test@123")).toBe(true);
      expect(numberValidator("hello-2-world")).toBe(true);
    });

    it("should return false for empty strings", () => {
      expect(numberValidator("")).toBe(false);
    });

    it("should return false for whitespace only", () => {
      expect(numberValidator("   ")).toBe(false);
      expect(numberValidator("\t\n")).toBe(false);
    });

    it("should handle numbers anywhere in the string", () => {
      expect(numberValidator("a1")).toBe(true);
      expect(numberValidator("1a")).toBe(true);
      expect(numberValidator("abc1def")).toBe(true);
    });
  });

  describe("upperCaseValidator", () => {
    it("should return true for strings with single uppercase letters", () => {
      expect(upperCaseValidator("A")).toBe(true);
      expect(upperCaseValidator("Z")).toBe(true);
      expect(upperCaseValidator("M")).toBe(true);
    });

    it("should return true for strings with multiple uppercase letters", () => {
      expect(upperCaseValidator("ABC")).toBe(true);
      expect(upperCaseValidator("TEST")).toBe(true);
      expect(upperCaseValidator("HELLO")).toBe(true);
    });

    it("should return false for strings with only lowercase letters", () => {
      expect(upperCaseValidator("abc")).toBe(false);
      expect(upperCaseValidator("test")).toBe(false);
      expect(upperCaseValidator("hello")).toBe(false);
    });

    it("should return true for mixed case strings", () => {
      expect(upperCaseValidator("Test")).toBe(true);
      expect(upperCaseValidator("heLLo")).toBe(true);
      expect(upperCaseValidator("aBcDeF")).toBe(true);
    });

    it("should return false for strings with only numbers", () => {
      expect(upperCaseValidator("123")).toBe(false);
      expect(upperCaseValidator("000")).toBe(false);
    });

    it("should return false for strings with only symbols", () => {
      expect(upperCaseValidator("!@#$%")).toBe(false);
      expect(upperCaseValidator("***")).toBe(false);
    });

    it("should return true for complex passwords with uppercase", () => {
      expect(upperCaseValidator("Password1!")).toBe(true);
      expect(upperCaseValidator("test@Test123")).toBe(true);
      expect(upperCaseValidator("hello-World")).toBe(true);
    });

    it("should return false for empty strings", () => {
      expect(upperCaseValidator("")).toBe(false);
    });

    it("should handle uppercase anywhere in the string", () => {
      expect(upperCaseValidator("Aaaa")).toBe(true);
      expect(upperCaseValidator("aaaA")).toBe(true);
      expect(upperCaseValidator("aaAaa")).toBe(true);
    });

    it("should work with numbers and symbols present", () => {
      expect(upperCaseValidator("Test123!@#")).toBe(true);
      expect(upperCaseValidator("123A456")).toBe(true);
      expect(upperCaseValidator("!@#A$%^")).toBe(true);
    });
  });

  describe("lowerCaseValidator", () => {
    it("should return true for strings with single lowercase letters", () => {
      expect(lowerCaseValidator("a")).toBe(true);
      expect(lowerCaseValidator("z")).toBe(true);
      expect(lowerCaseValidator("m")).toBe(true);
    });

    it("should return true for strings with multiple lowercase letters", () => {
      expect(lowerCaseValidator("abc")).toBe(true);
      expect(lowerCaseValidator("test")).toBe(true);
      expect(lowerCaseValidator("hello")).toBe(true);
    });

    it("should return false for strings with only uppercase letters", () => {
      expect(lowerCaseValidator("ABC")).toBe(false);
      expect(lowerCaseValidator("TEST")).toBe(false);
      expect(lowerCaseValidator("HELLO")).toBe(false);
    });

    it("should return true for mixed case strings", () => {
      expect(lowerCaseValidator("Test")).toBe(true);
      expect(lowerCaseValidator("heLLo")).toBe(true);
      expect(lowerCaseValidator("aBcDeF")).toBe(true);
    });

    it("should return false for strings with only numbers", () => {
      expect(lowerCaseValidator("123")).toBe(false);
      expect(lowerCaseValidator("000")).toBe(false);
    });

    it("should return false for strings with only symbols", () => {
      expect(lowerCaseValidator("!@#$%")).toBe(false);
      expect(lowerCaseValidator("***")).toBe(false);
    });

    it("should return true for complex passwords with lowercase", () => {
      expect(lowerCaseValidator("Password1!")).toBe(true);
      expect(lowerCaseValidator("TEST@test123")).toBe(true);
      expect(lowerCaseValidator("HELLO-world")).toBe(true);
    });

    it("should return false for empty strings", () => {
      expect(lowerCaseValidator("")).toBe(false);
    });

    it("should handle lowercase anywhere in the string", () => {
      expect(lowerCaseValidator("aAAA")).toBe(true);
      expect(lowerCaseValidator("AAAa")).toBe(true);
      expect(lowerCaseValidator("AAaAA")).toBe(true);
    });

    it("should work with numbers and symbols present", () => {
      expect(lowerCaseValidator("test123!@#")).toBe(true);
      expect(lowerCaseValidator("123a456")).toBe(true);
      expect(lowerCaseValidator("!@#a$%^")).toBe(true);
    });
  });

  describe("password complexity validation scenarios", () => {
    it("should validate a strong password with all requirements", () => {
      const password = "MyP@ssw0rd!";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(true);
    });

    it("should fail weak password without symbols", () => {
      const password = "MyPassword123";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(false);
    });

    it("should fail password without uppercase", () => {
      const password = "myp@ssw0rd!";

      expect(upperCaseValidator(password)).toBe(false);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(true);
    });

    it("should fail password without lowercase", () => {
      const password = "MYP@SSW0RD!";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(false);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(true);
    });

    it("should fail password without numbers", () => {
      const password = "MyP@ssword!";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(false);
      expect(symbolValidator(password)).toBe(true);
    });

    it("should validate very long complex password", () => {
      const password = "ThisIsAVeryLongP@ssw0rd!WithManyCharacters123";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(true);
    });

    it("should validate password with multiple symbols", () => {
      const password = "P@$$w0rd!#%";

      expect(upperCaseValidator(password)).toBe(true);
      expect(lowerCaseValidator(password)).toBe(true);
      expect(numberValidator(password)).toBe(true);
      expect(symbolValidator(password)).toBe(true);
    });
  });
});
