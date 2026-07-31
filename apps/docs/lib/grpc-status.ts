// Mirrors the table in content/apis/statuscodes.mdx. Kept here so
// scripts/generate-error-reference.ts can derive HTTP statuses/example
// responses without hand-duplicating the mapping.

export interface GrpcStatusEntry {
  code: number;
  name: string;
  httpStatus: number;
  httpText: string;
}

export const GRPC_STATUS: Record<string, GrpcStatusEntry> = {
  OK: { code: 0, name: 'OK', httpStatus: 200, httpText: 'OK' },
  CANCELLED: { code: 1, name: 'CANCELLED', httpStatus: 499, httpText: 'Client Closed Request' },
  UNKNOWN: { code: 2, name: 'UNKNOWN', httpStatus: 500, httpText: 'Internal' },
  INVALID_ARGUMENT: { code: 3, name: 'INVALID_ARGUMENT', httpStatus: 400, httpText: 'Bad Request' },
  DEADLINE_EXCEEDED: { code: 4, name: 'DEADLINE_EXCEEDED', httpStatus: 504, httpText: 'Gateway Timeout' },
  NOT_FOUND: { code: 5, name: 'NOT_FOUND', httpStatus: 404, httpText: 'Not found' },
  ALREADY_EXISTS: { code: 6, name: 'ALREADY_EXISTS', httpStatus: 409, httpText: 'Conflict' },
  PERMISSION_DENIED: { code: 7, name: 'PERMISSION_DENIED', httpStatus: 403, httpText: 'Forbidden' },
  RESOURCE_EXHAUSTED: { code: 8, name: 'RESOURCE_EXHAUSTED', httpStatus: 429, httpText: 'Too Many Requests' },
  FAILED_PRECONDITION: { code: 9, name: 'FAILED_PRECONDITION', httpStatus: 400, httpText: 'Bad Request' },
  ABORTED: { code: 10, name: 'ABORTED', httpStatus: 409, httpText: 'Conflict' },
  OUT_OF_RANGE: { code: 11, name: 'OUT_OF_RANGE', httpStatus: 400, httpText: 'Bad Request' },
  UNIMPLEMENTED: { code: 12, name: 'UNIMPLEMENTED', httpStatus: 501, httpText: 'Not Implemented' },
  INTERNAL: { code: 13, name: 'INTERNAL', httpStatus: 500, httpText: 'Internal' },
  UNAVAILABLE: { code: 14, name: 'UNAVAILABLE', httpStatus: 503, httpText: 'Service Unavailable' },
  DATA_LOSS: { code: 15, name: 'DATA_LOSS', httpStatus: 500, httpText: 'Internal' },
  UNAUTHENTICATED: { code: 16, name: 'UNAUTHENTICATED', httpStatus: 401, httpText: 'Unauthorized' },
};

// zerrors.Throw<Kind> function-name suffix -> gRPC status key.
export const THROW_KIND_TO_GRPC_STATUS: Record<string, keyof typeof GRPC_STATUS> = {
  Internal: 'INTERNAL',
  InvalidArgument: 'INVALID_ARGUMENT',
  NotFound: 'NOT_FOUND',
  AlreadyExists: 'ALREADY_EXISTS',
  PermissionDenied: 'PERMISSION_DENIED',
  PreconditionFailed: 'FAILED_PRECONDITION',
  Unauthenticated: 'UNAUTHENTICATED',
  Unimplemented: 'UNIMPLEMENTED',
  ResourceExhausted: 'RESOURCE_EXHAUSTED',
  Unavailable: 'UNAVAILABLE',
  DeadlineExceeded: 'DEADLINE_EXCEEDED',
  Aborted: 'ABORTED',
  OutOfRange: 'OUT_OF_RANGE',
  DataLoss: 'DATA_LOSS',
  Unknown: 'UNKNOWN',
  Canceled: 'CANCELLED',
};
