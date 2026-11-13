import { setLanguageCookie } from "@/lib/actions/language";
import { beforeEach, describe, expect, test, vi } from "vitest";

vi.mock("@/lib/server/cookies", () => ({
  setLanguageCookie: vi.fn(),
}));

describe("setLanguageCookie", () => {
  let setCookie: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    setCookie = vi.mocked((await import("@/lib/server/cookies")).setLanguageCookie);
  });

  test("sets a supported language", async () => {
    await setLanguageCookie("de");

    expect(setCookie).toHaveBeenCalledWith("de");
  });

  test("ignores a value that is not a supported language", async () => {
    await setLanguageCookie("de; Path=/evil");

    expect(setCookie).not.toHaveBeenCalled();
  });
});
