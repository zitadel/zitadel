import { describe, expect, test } from "vitest";
import { makeReqCtx } from "./context.js";

describe("makeReqCtx", () => {
  test("returns orgId when provided", () => {
    const ctx = makeReqCtx("org-123");
    expect(ctx).toEqual({ orgId: "org-123" });
  });

  test("returns empty object when no orgId", () => {
    const ctx = makeReqCtx();
    expect(ctx).toEqual({});
  });
});
