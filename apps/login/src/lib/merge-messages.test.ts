import { describe, expect, it } from "vitest";
import { mergeMessages } from "./merge-messages";

describe("mergeMessages", () => {
  it("lets later sources override earlier ones", () => {
    expect(mergeMessages({ title: "Login" }, { title: "Anmelden" }, { title: "Custom" })).toEqual({ title: "Custom" });
  });

  it("merges nested keys instead of replacing the whole object", () => {
    const fallback = { loginname: { title: "Welcome back!", description: "Enter your login data." } };
    const locale = { loginname: { title: "Willkommen zurück!" } };

    expect(mergeMessages(fallback, locale)).toEqual({
      loginname: { title: "Willkommen zurück!", description: "Enter your login data." },
    });
  });

  it("keeps fallback messages for keys the locale is missing", () => {
    const fallback = { password: { verify: { title: "Password" } }, register: { title: "Register" } };

    expect(mergeMessages(fallback, { password: {} }, {})).toEqual(fallback);
  });

  it("replaces arrays and primitives instead of merging them", () => {
    expect(mergeMessages({ a: ["x"], b: { c: "d" } }, { a: ["y"], b: "flat" })).toEqual({ a: ["y"], b: "flat" });
  });

  it("does not modify its inputs", () => {
    const fallback = { loginname: { title: "Welcome back!" } };
    mergeMessages(fallback, { loginname: { title: "Custom" } });

    expect(fallback).toEqual({ loginname: { title: "Welcome back!" } });
  });

  it("ignores missing sources", () => {
    expect(mergeMessages({ title: "Login" }, undefined, null)).toEqual({ title: "Login" });
  });

  it("returns a plain object", () => {
    expect(Object.getPrototypeOf(mergeMessages({ title: "Login" }))).toBe(Object.prototype);
  });

  it("skips keys that would change prototypes", () => {
    const custom = JSON.parse(
      '{"__proto__": {"polluted": "yes"}, "constructor": {"prototype": {"polluted": "yes"}}, "loginname": {"__proto__": {"polluted": "yes"}, "title": "Custom"}}',
    );

    const merged = mergeMessages({ loginname: { title: "Welcome back!" } }, custom);

    expect(merged).toEqual({ loginname: { title: "Custom" } });
    expect(Object.getPrototypeOf(merged)).toBe(Object.prototype);
    expect(Object.getPrototypeOf(merged.loginname)).toBe(Object.prototype);
    expect(({} as Record<string, unknown>).polluted).toBeUndefined();
  });
});
