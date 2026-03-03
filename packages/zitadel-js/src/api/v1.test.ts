import { describe, expect, test } from "vitest";
import * as v1Exports from "./v1.js";

describe("api/v1 exports", () => {
  test("expose legacy v1 API client factories", () => {
    expect(Object.keys(v1Exports).sort()).toEqual([
      "createAdminServiceClient",
      "createAuthServiceClient",
      "createManagementServiceClient",
      "createSystemServiceClient",
    ]);
  });
});
