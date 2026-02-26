export interface RequestContext {
  orgId?: string;
  instanceId?: string;
}

/**
 * Creates a request context metadata object for ZITADEL API calls.
 */
export function makeReqCtx(orgId?: string): RequestContext {
  return orgId ? { orgId } : {};
}
