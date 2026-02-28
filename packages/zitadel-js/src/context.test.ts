import { describe, expect, test } from "vitest";
import { makeReqCtx } from "./context.js";

describe("makeReqCtx", () => {
  test("returns orgId when provided", () => {
    const ctx = makeReqCtx("org-123");
    expect(ctx.resourceOwner).toEqual({
      case: "orgId",
      value: "org-123",
    });
  });

  test("returns instance scope when no orgId", () => {
    const ctx = makeReqCtx();
    expect(ctx.resourceOwner).toEqual({
      case: "instance",
      value: true,
    });
  });
});
