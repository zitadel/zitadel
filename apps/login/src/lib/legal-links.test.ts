import { describe, expect, it } from "vitest";
import { resolveLocalizedLegalLink } from "./legal-links";

describe("resolveLocalizedLegalLink", () => {
  it("replaces the language placeholder with the active locale", () => {
    expect(resolveLocalizedLegalLink("https://example.com/{{.Lang}}/terms", "de")).toBe("https://example.com/de/terms");
  });

  it("replaces every language placeholder in the link", () => {
    expect(resolveLocalizedLegalLink("https://example.com/{{.Lang}}/docs?lang={{.Lang}}", "fr")).toBe(
      "https://example.com/fr/docs?lang=fr",
    );
  });

  it("returns the original link when no placeholder is present", () => {
    expect(resolveLocalizedLegalLink("https://example.com/terms", "en")).toBe("https://example.com/terms");
  });

  it("returns the original link when locale is missing", () => {
    expect(resolveLocalizedLegalLink("https://example.com/{{.Lang}}/terms", undefined)).toBe(
      "https://example.com/{{.Lang}}/terms",
    );
  });

  it("returns undefined when the link is missing", () => {
    expect(resolveLocalizedLegalLink(undefined, "en")).toBeUndefined();
  });
});
