import sk from "../../locales/sk.json";
import en from "../../locales/en.json";
import { describe, expect, it } from "vitest";
import { LANGS, getLanguage, normalizeLanguageCode } from "./i18n";

function flatten(value: unknown, path = ""): Map<string, string> {
  if (typeof value === "string") return new Map([[path, value]]);
  if (!value || typeof value !== "object" || Array.isArray(value)) throw new Error("Expected a message object at " + (path || "root"));
  return Object.entries(value).reduce((messages, [key, nestedValue]) => {
    for (const [nestedPath, message] of flatten(nestedValue, path ? path + "." + key : key)) messages.set(nestedPath, message);
    return messages;
  }, new Map<string, string>());
}

function placeholders(message: string): string[] { return [...message.matchAll(/\\{([^}]+)\\}/g)].map((match) => match[1]).sort(); }

describe("Login V2 Slovak locale", () => {
  it("registers sk with its Slovak display name", () => {
    expect(LANGS).toContainEqual({ code: "sk", name: "Slovenčina" });
    expect(getLanguage("sk")).toEqual({ code: "sk", name: "Slovenčina" });
  });
  it("normalizes Slovak regional tags to sk", () => {
    expect(normalizeLanguageCode("sk-SK")).toBe("sk");
    expect(normalizeLanguageCode("SK-sk")).toBe("sk");
  });
  it("has the complete English key set, non-empty messages, and matching placeholders", () => {
    const english = flatten(en);
    const slovak = flatten(sk);
    expect([...slovak.keys()].sort()).toEqual([...english.keys()].sort());
    for (const [key, englishMessage] of english) {
      const slovakMessage = slovak.get(key);
      expect(slovakMessage, key + " must not be empty").toBeTruthy();
      expect(placeholders(slovakMessage!)).toEqual(placeholders(englishMessage));
    }
  });
});
