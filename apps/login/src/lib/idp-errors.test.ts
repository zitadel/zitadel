import { describe, expect, test } from "vitest";
import { isKnownIdpError } from "./idp-errors";

describe("isKnownIdpError", () => {
  test.each([
    ["flow_mismatch", true],
    ["missing_idp_info", true],
    ["session_mismatch", true],
    ["external_idp_taken", true],
    ["Your account is locked, call 555", false],
    ["", false],
    [undefined, false],
  ])("%s -> %s", (code, expected) => {
    expect(isKnownIdpError(code)).toBe(expected);
  });
});
