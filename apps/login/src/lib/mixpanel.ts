import mixpanel from "mixpanel-browser";

let initialized = false;

export function initMixpanel() {
  if (typeof window === "undefined") return;
  if (initialized) return;

  const token = process.env.NEXT_PUBLIC_MIXPANEL_TOKEN;
  if (!token) return;

  mixpanel.init(token, {
    track_pageview: false,
    ip: false,
    persistence: "localStorage",
  });

  initialized = true;
}

export const MixpanelEvents = {
  // Login
  username_submitted: "username_submitted",
  password_submitted: "password_submitted",
  login_success: "login_success",
  login_failure: "login_failure",
  password_reset_requested: "password_reset_requested",

  // Registration
  register_submitted: "register_submitted",
  register_method_selected: "register_method_selected",
  register_password_submitted: "register_password_submitted",
  register_success: "register_success",
  register_failure: "register_failure",

  // IDP
  idp_button_clicked: "idp_button_clicked",
  idp_register_submitted: "idp_register_submitted",
  idp_callback_started: "idp_callback_started",
  idp_callback_success: "idp_callback_success",
  idp_callback_failure: "idp_callback_failure",

  // Passkey
  passkey_login_started: "passkey_login_started",
  passkey_login_success: "passkey_login_success",
  passkey_login_failure: "passkey_login_failure",
  passkey_register_started: "passkey_register_started",
  passkey_register_success: "passkey_register_success",
  passkey_register_failure: "passkey_register_failure",
  passkey_register_skipped: "passkey_register_skipped",

  // U2F
  u2f_register_started: "u2f_register_started",
  u2f_register_success: "u2f_register_success",
  u2f_register_failure: "u2f_register_failure",

  // OTP
  otp_code_submitted: "otp_code_submitted",
  otp_code_resent: "otp_code_resent",
  otp_success: "otp_success",
  otp_failure: "otp_failure",

  // TOTP
  totp_setup_code_submitted: "totp_setup_code_submitted",
  totp_setup_success: "totp_setup_success",
  totp_setup_failure: "totp_setup_failure",

  // Verify
  verify_code_submitted: "verify_code_submitted",
  verify_code_resent: "verify_code_resent",
  verify_success: "verify_success",
  verify_failure: "verify_failure",

  // Password management
  password_set_submitted: "password_set_submitted",
  password_set_success: "password_set_success",
  password_set_failure: "password_set_failure",
  password_change_submitted: "password_change_submitted",
  password_change_success: "password_change_success",
  password_change_failure: "password_change_failure",

  // LDAP
  ldap_login_submitted: "ldap_login_submitted",
  ldap_login_success: "ldap_login_success",
  ldap_login_failure: "ldap_login_failure",

  // Device flow
  device_code_submitted: "device_code_submitted",
  device_code_success: "device_code_success",
  device_code_failure: "device_code_failure",
  device_consent_approved: "device_consent_approved",
  device_consent_denied: "device_consent_denied",

  // Session
  session_selected: "session_selected",
  session_cleared: "session_cleared",

  // MFA
  mfa_method_selected: "mfa_method_selected",
  mfa_setup_method_selected: "mfa_setup_method_selected",
  mfa_setup_skipped: "mfa_setup_skipped",

  // Page view
  page_view: "page_view",
} as const;

export type MixpanelEvent =
  (typeof MixpanelEvents)[keyof typeof MixpanelEvents];

export function trackEvent(
  event: MixpanelEvent,
  properties?: Record<string, string | number | boolean>,
) {
  if (!initialized) return;
  mixpanel.track(event, properties);
}

export function trackPageView(pathname: string) {
  trackEvent(MixpanelEvents.page_view, { pathname });
}
