import { describe, expect, test } from "vitest";
import * as rootExports from "./index.js";

describe("root exports", () => {
  test("only expose shared transport/client primitives", () => {
    expect(Object.keys(rootExports).sort()).toEqual([
      "createClientFor",
      "createConnectTransport",
      "createGrpcTransport",
    ]);
  });
});
