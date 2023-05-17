import * as settings from "./v2/settings";
import * as session from "./v2/session";
import * as user from "./v2/user";

import * as login from "./proto/server/zitadel/settings/v2alpha/login_settings";
import * as password from "./proto/server/zitadel/settings/v2alpha/password_settings";
import * as legal from "./proto/server/zitadel/settings/v2alpha/legal_settings";

export {
  BrandingSettings,
  Theme,
} from "./proto/server/zitadel/settings/v2alpha/branding_settings";

export { type LegalAndSupportSettings } from "./proto/server/zitadel/settings/v2alpha/legal_settings";
export { type PasswordComplexitySettings } from "./proto/server/zitadel/settings/v2alpha/password_settings";

import {
  getServers,
  initializeServer,
  ZitadelServer,
  ZitadelServerOptions,
} from "./server";
export * from "./middleware";

export {
  getServers,
  ZitadelServer,
  type ZitadelServerOptions,
  initializeServer,
  user,
  session,
  settings,
  login,
  password,
  legal,
};
