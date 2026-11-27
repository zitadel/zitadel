import { hasZitadelErrorId, isRejectedByUser, zitadelErrorId, ZitadelErrorIds } from "@/lib/server/api-error";
import { ClassifiedConnectError } from "@/lib/server/grpc/interceptors/error-classification";
import { Code, ConnectError } from "@connectrpc/connect";
import { ErrorDetailSchema } from "@zitadel/proto/zitadel/message_pb";
import { describe, expect, test } from "vitest";

const apiError = (code: Code, message: string, detailId?: string) =>
  new ClassifiedConnectError(
    new ConnectError(
      message,
      code,
      undefined,
      detailId ? [{ desc: ErrorDetailSchema, value: { id: detailId, message } }] : undefined,
    ),
  );

describe("isRejectedByUser", () => {
  test.each([Code.InvalidArgument, Code.FailedPrecondition, Code.NotFound, Code.ResourceExhausted])(
    "counts code %s as a rejection caused by the user",
    (code) => {
      expect(isRejectedByUser(apiError(code, "rejected"))).toBe(true);
    },
  );

  test.each([Code.Unauthenticated, Code.PermissionDenied])(
    "does not count code %s: the API rejected the login's own credentials",
    (code) => {
      expect(isRejectedByUser(apiError(code, "rejected"))).toBe(false);
    },
  );

  test("does not count server faults or errors that did not come from the API", () => {
    expect(isRejectedByUser(apiError(Code.Internal, "boom"))).toBe(false);
    expect(isRejectedByUser(new Error("boom"))).toBe(false);
  });
});

describe("zitadelErrorId", () => {
  test("reads the id from the error detail", () => {
    expect(zitadelErrorId(apiError(Code.FailedPrecondition, "User is locked", "COMMAND-SF3fg"))).toBe("COMMAND-SF3fg");
  });

  test("falls back to the id that ends the translated message", () => {
    expect(zitadelErrorId(apiError(Code.FailedPrecondition, "Benutzer ist gesperrt (COMMAND-SF3fg)"))).toBe("COMMAND-SF3fg");
  });

  test("is undefined without an id", () => {
    expect(zitadelErrorId(apiError(Code.InvalidArgument, "rejected"))).toBeUndefined();
  });

  test("matches a list of ids", () => {
    const locked = apiError(Code.FailedPrecondition, "User is locked (COMMAND-S6h4R)");

    expect(hasZitadelErrorId(locked, ZitadelErrorIds.userLocked)).toBe(true);
    expect(hasZitadelErrorId(locked, ZitadelErrorIds.codeExpired)).toBe(false);
  });
});
