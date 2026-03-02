import { describe, expect, test } from "vitest";
import {
  isSessionExpired as isSessionExpiredFromAuthLane,
  isSessionValid as isSessionValidFromAuthLane,
} from "./session.js";
import { isSessionExpired, isSessionValid } from "../session.js";

describe("auth/session exports", () => {
  test("re-export session helpers", () => {
    expect(isSessionExpiredFromAuthLane).toBe(isSessionExpired);
    expect(isSessionValidFromAuthLane).toBe(isSessionValid);
  });
});
