import { Code, ConnectError } from "@connectrpc/connect";
import { describe, expect, it } from "vitest";
import { ClassifiedConnectError, grpcCodeToHttpStatus, isClassifiedError } from "./error-classification";

describe("grpcCodeToHttpStatus", () => {
  const cases: [Code, number][] = [
    [Code.InvalidArgument, 400],
    [Code.FailedPrecondition, 400],
    [Code.OutOfRange, 400],
    [Code.Unauthenticated, 401],
    [Code.PermissionDenied, 403],
    [Code.NotFound, 404],
    [Code.AlreadyExists, 409],
    [Code.Aborted, 409],
    [Code.ResourceExhausted, 429],
    [Code.Canceled, 499],
    [Code.Unimplemented, 501],
    [Code.Unavailable, 503],
    [Code.DeadlineExceeded, 504],
    [Code.Internal, 500],
    [Code.DataLoss, 500],
    [Code.Unknown, 500],
  ];

  it.each(cases)("maps gRPC code %i to HTTP status %i", (grpcCode, expectedHttp) => {
    expect(grpcCodeToHttpStatus(grpcCode)).toBe(expectedHttp);
  });

  it("defaults to 500 for unmapped codes", () => {
    expect(grpcCodeToHttpStatus(999 as Code)).toBe(500);
  });
});

describe("ClassifiedConnectError", () => {
  it("preserves original error properties", () => {
    const source = new ConnectError("not found", Code.NotFound);
    const classified = new ClassifiedConnectError(source);

    expect(classified.message).toContain("not found");
    expect(classified.code).toBe(Code.NotFound);
    // name must remain "ConnectError" so ConnectError's custom Symbol.hasInstance passes
    expect(classified.name).toBe("ConnectError");
  });

  it("is instanceof ConnectError (required for ConnectError.from compatibility)", () => {
    const source = new ConnectError("test", Code.FailedPrecondition);
    const classified = new ClassifiedConnectError(source);

    expect(classified instanceof ConnectError).toBe(true);
  });

  it("sets httpStatus from gRPC code", () => {
    const source = new ConnectError("permission denied", Code.PermissionDenied);
    const classified = new ClassifiedConnectError(source);

    expect(classified.httpStatus).toBe(403);
  });

  it("marks client errors correctly", () => {
    const clientError = new ClassifiedConnectError(new ConnectError("bad input", Code.InvalidArgument));
    const serverError = new ClassifiedConnectError(new ConnectError("internal", Code.Internal));

    expect(clientError.isUserError).toBe(true);
    expect(serverError.isUserError).toBe(false);
  });

  it("is detectable via isClassifiedError type guard", () => {
    const classified = new ClassifiedConnectError(new ConnectError("test", Code.NotFound));

    expect(isClassifiedError(classified)).toBe(true);
  });

  it("marks FailedPrecondition as client error", () => {
    const classified = new ClassifiedConnectError(new ConnectError("precondition", Code.FailedPrecondition));

    expect(classified.isUserError).toBe(true);
    expect(classified.httpStatus).toBe(400);
  });

  it("marks Unavailable as server error", () => {
    const classified = new ClassifiedConnectError(new ConnectError("unavailable", Code.Unavailable));

    expect(classified.isUserError).toBe(false);
    expect(classified.httpStatus).toBe(503);
  });
});

describe("isClassifiedError", () => {
  it("returns true for ClassifiedConnectError", () => {
    const classified = new ClassifiedConnectError(new ConnectError("test", Code.NotFound));
    expect(isClassifiedError(classified)).toBe(true);
  });

  it("returns false for plain ConnectError", () => {
    const plain = new ConnectError("test", Code.NotFound);
    expect(isClassifiedError(plain)).toBe(false);
  });

  it("returns false for plain Error", () => {
    expect(isClassifiedError(new Error("test"))).toBe(false);
  });

  it("returns false for non-error values", () => {
    expect(isClassifiedError(null)).toBe(false);
    expect(isClassifiedError(undefined)).toBe(false);
    expect(isClassifiedError("string")).toBe(false);
    expect(isClassifiedError({ code: 5 })).toBe(false);
  });
});
