// only these codes are shown on the failure and linking-failed pages, anything else in the URL is ignored
const KNOWN_IDP_ERRORS = new Set([
  "flow_mismatch",
  "no_organization_context",
  "missing_idp_user_info",
  "user_creation_failed",
  "missing_idp_info",
  "session_invalid",
  "session_mismatch",
  "linking_not_allowed",
  "user_not_found",
  "validation_failed",
  "missing_organization",
  "external_idp_taken",
]);

export function isKnownIdpError(code: string | undefined): code is string {
  return !!code && KNOWN_IDP_ERRORS.has(code);
}
