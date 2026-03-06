import { create } from "@bufbuild/protobuf";
import {
  RequestContext,
  RequestContextSchema,
} from "./generated/zitadel/object/v2/object_pb.js";

export type { RequestContext };

/**
 * Creates a request context metadata object for ZITADEL API calls.
 */
export function makeReqCtx(orgId?: string): RequestContext {
  return create(RequestContextSchema, {
    resourceOwner: orgId
      ? { case: "orgId", value: orgId }
      : { case: "instance", value: true },
  });
}
