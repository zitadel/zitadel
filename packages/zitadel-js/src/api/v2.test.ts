import { describe, expect, test } from "vitest";
import * as legacyV2 from "../v2.js";
import * as apiV2 from "./v2.js";

describe("api/v2 exports", () => {
  test("mirror legacy v2 exports", () => {
    expect(Object.keys(apiV2).sort()).toEqual(Object.keys(legacyV2).sort());
    expect(apiV2.createUserServiceClient).toBe(legacyV2.createUserServiceClient);
  });
});
