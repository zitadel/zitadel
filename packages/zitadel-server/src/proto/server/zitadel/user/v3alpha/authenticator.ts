/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../google/protobuf/struct";
import { Timestamp } from "../../../google/protobuf/timestamp";
import { Details } from "../../object/v2beta/object";

export const protobufPackage = "zitadel.user.v3alpha";

export enum AuthNKeyType {
  AUTHN_KEY_TYPE_UNSPECIFIED = 0,
  AUTHN_KEY_TYPE_JSON = 1,
  UNRECOGNIZED = -1,
}

export function authNKeyTypeFromJSON(object: any): AuthNKeyType {
  switch (object) {
    case 0:
    case "AUTHN_KEY_TYPE_UNSPECIFIED":
      return AuthNKeyType.AUTHN_KEY_TYPE_UNSPECIFIED;
    case 1:
    case "AUTHN_KEY_TYPE_JSON":
      return AuthNKeyType.AUTHN_KEY_TYPE_JSON;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AuthNKeyType.UNRECOGNIZED;
  }
}

export function authNKeyTypeToJSON(object: AuthNKeyType): string {
  switch (object) {
    case AuthNKeyType.AUTHN_KEY_TYPE_UNSPECIFIED:
      return "AUTHN_KEY_TYPE_UNSPECIFIED";
    case AuthNKeyType.AUTHN_KEY_TYPE_JSON:
      return "AUTHN_KEY_TYPE_JSON";
    case AuthNKeyType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum WebAuthNAuthenticatorType {
  WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED = 0,
  WEB_AUTH_N_AUTHENTICATOR_PLATFORM = 1,
  WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM = 2,
  UNRECOGNIZED = -1,
}

export function webAuthNAuthenticatorTypeFromJSON(object: any): WebAuthNAuthenticatorType {
  switch (object) {
    case 0:
    case "WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED":
      return WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED;
    case 1:
    case "WEB_AUTH_N_AUTHENTICATOR_PLATFORM":
      return WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_PLATFORM;
    case 2:
    case "WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM":
      return WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM;
    case -1:
    case "UNRECOGNIZED":
    default:
      return WebAuthNAuthenticatorType.UNRECOGNIZED;
  }
}

export function webAuthNAuthenticatorTypeToJSON(object: WebAuthNAuthenticatorType): string {
  switch (object) {
    case WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED:
      return "WEB_AUTH_N_AUTHENTICATOR_UNSPECIFIED";
    case WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_PLATFORM:
      return "WEB_AUTH_N_AUTHENTICATOR_PLATFORM";
    case WebAuthNAuthenticatorType.WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM:
      return "WEB_AUTH_N_AUTHENTICATOR_CROSS_PLATFORM";
    case WebAuthNAuthenticatorType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Authenticators {
  /** All of the user's usernames, which will be used for identification during authentication. */
  usernames: Username[];
  /** If the user has set a password, the time it was last changed will be returned. */
  password:
    | Password
    | undefined;
  /** Meta information about the user's WebAuthN authenticators. */
  webAuthN: WebAuthN[];
  /**
   * A list of the user's time-based one-time-password (TOTP) authenticators,
   * incl. the name for identification.
   */
  totps: TOTP[];
  /** A list of the user's one-time-password (OTP) SMS authenticators. */
  otpSms: OTPSMS[];
  /** A list of the user's one-time-password (OTP) Email authenticators. */
  otpEmail: OTPEmail[];
  /** A list of the user's authentication keys. They can be used to authenticate e.g. by JWT Profile. */
  authenticationKeys: AuthenticationKey[];
  /** A list of the user's linked identity providers (IDPs). */
  identityProviders: IdentityProvider[];
}

export interface Username {
  /** unique identifier of the username. */
  usernameId: string;
  /** The user's unique username. It is used for identification during authentication. */
  username: string;
  /**
   * By default usernames must be unique across all organizations in an instance.
   * This option allow to restrict the uniqueness to the user's own organization.
   * As a result, this username can only be used if the authentication is limited
   * to the corresponding organization.
   *
   * This can be useful if you provide multiple usernames for a single user, where one
   * if specific to your organization, e.g.:
   * - gigi-giraffe@zitadel.com (unique across organizations)
   * - gigi-giraffe (unique only inside the ZITADEL organization)
   */
  isOrganizationSpecific: boolean;
}

export interface SetUsername {
  /** Set the user's username. This will be used for identification during authentication. */
  username: string;
  /**
   * By default username must be unique across all organizations in an instance.
   * This option allow to restrict the uniqueness to the user's own organization.
   * As a result, this username can only be used if the authentication is limited
   * to the corresponding organization.
   *
   * This can be useful if you provide multiple usernames for a single user, where one
   * if specific to your organization, e.g.:
   * - gigi-giraffe@zitadel.com (unique across organizations)
   * - gigi-giraffe (unique only inside the ZITADEL organization)
   */
  isOrganizationSpecific: boolean;
}

export interface Password {
  /** States the time the password was last changed. */
  lastChanged: Date | undefined;
}

export interface WebAuthN {
  /** unique identifier of the WebAuthN authenticator. */
  webAuthNId: string;
  /** Name of the WebAuthN authenticator. This is used for easier identification. */
  name: string;
  /** State whether the WebAuthN registration has been completed. */
  isVerified: boolean;
  /**
   * States if the user has been verified during the registration. Authentication with this device
   * will be considered as multi factor authentication (MFA) without the need to check a password
   * (typically known as Passkeys).
   * Without user verification it will be a second factor authentication (2FA), typically done
   * after a password check.
   *
   * More on WebAuthN User Verification: https://www.w3.org/TR/webauthn/#user-verification
   */
  userVerified: boolean;
}

export interface OTPSMS {
  /** unique identifier of the one-time-password (OTP) SMS authenticator. */
  otpSmsId: string;
  /** The phone number used for the OTP SMS authenticator. */
  phone: string;
  /** State whether the OTP SMS registration has been completed. */
  isVerified: boolean;
}

export interface OTPEmail {
  /** unique identifier of the one-time-password (OTP) Email authenticator. */
  otpEmailId: string;
  /** The email address used for the OTP Email authenticator. */
  address: string;
  /** State whether the OTP Email registration has been completed. */
  isVerified: boolean;
}

export interface TOTP {
  /** unique identifier of the time-based one-time-password (TOTP) authenticator. */
  totpId: string;
  /** The name provided during registration. This is used for easier identification. */
  name: string;
  /** State whether the TOTP registration has been completed. */
  isVerified: boolean;
}

export interface AuthenticationKey {
  /** ID is the read-only unique identifier of the authentication key. */
  authenticationKeyId: string;
  details:
    | Details
    | undefined;
  /** the file type of the key */
  type: AuthNKeyType;
  /** After the expiration date, the key will no longer be usable for authentication. */
  expirationDate: Date | undefined;
}

export interface IdentityProvider {
  /** IDP ID is the read-only unique identifier of the identity provider in ZITADEL. */
  idpId: string;
  /** IDP name is the name of the identity provider in ZITADEL. */
  idpName: string;
  /**
   * The user ID represents the ID provided by the identity provider.
   * This ID is used to link the user in ZITADEL with the identity provider.
   */
  userId: string;
  /** The username represents the username provided by the identity provider. */
  username: string;
}

export interface SetAuthenticators {
  usernames: SetUsername[];
  password: SetPassword | undefined;
}

export interface SetPassword {
  /** Provide the plain text password. ZITADEL will take care to store it in a secure way (hash). */
  password?:
    | string
    | undefined;
  /**
   * Encoded hash of a password in Modular Crypt Format:
   * https://zitadel.com/docs/concepts/architecture/secrets#hashed-secrets.
   */
  hash?:
    | string
    | undefined;
  /** Provide if the user needs to change the password on the next use. */
  changeRequired: boolean;
}

export interface SendPasswordResetEmail {
  /**
   * Optionally set a url_template, which will be used in the password reset mail
   * sent by ZITADEL to guide the user to your password change page.
   * If no template is set, the default ZITADEL url will be used.
   */
  urlTemplate?: string | undefined;
}

export interface SendPasswordResetSMS {
}

export interface ReturnPasswordResetCode {
}

export interface AuthenticatorRegistrationCode {
  /** ID to the one time code generated by ZITADEL. */
  id: string;
  /** one time code generated by ZITADEL. */
  code: string;
}

export interface SendWebAuthNRegistrationLink {
  /**
   * Optionally set a url_template, which will be used in the mail sent by ZITADEL
   * to guide the user to your passkey registration page.
   * If no template is set, the default ZITADEL url will be used.
   */
  urlTemplate?: string | undefined;
}

export interface ReturnWebAuthNRegistrationCode {
}

export interface RedirectURLs {
  /** URL to which the user will be redirected after a successful login. */
  successUrl: string;
  /** URL to which the user will be redirected after a failed login. */
  failureUrl: string;
}

export interface LDAPCredentials {
  /** Username used to login through LDAP. */
  username: string;
  /** Password used to login through LDAP. */
  password: string;
}

export interface IdentityProviderIntent {
  /** ID of the identity provider (IDP) intent. */
  idpIntentId: string;
  /** Token of the identity provider (IDP) intent. */
  idpIntentToken: string;
  /** If the user was already federated and linked to a ZITADEL user, it's id will be returned. */
  userId?: string | undefined;
}

export interface IDPInformation {
  /** ID of the identity provider. */
  idpId: string;
  /** ID of the user provided by the identity provider. */
  userId: string;
  /** Username of the user provided by the identity provider. */
  userName: string;
  /** Complete information returned by the identity provider. */
  rawInformation:
    | { [key: string]: any }
    | undefined;
  /** OAuth/OIDC access (and id_token) returned by the identity provider. */
  oauth?:
    | IDPOAuthAccessInformation
    | undefined;
  /** LDAP entity attributes returned by the identity provider */
  ldap?:
    | IDPLDAPAccessInformation
    | undefined;
  /** SAMLResponse returned by the identity provider */
  saml?: IDPSAMLAccessInformation | undefined;
}

export interface IDPOAuthAccessInformation {
  /** The access_token returned by the identity provider. */
  accessToken: string;
  /** In case the provider returned an id_token. */
  idToken?: string | undefined;
}

export interface IDPLDAPAccessInformation {
  /** The attributes of the user returned by the identity provider. */
  attributes: { [key: string]: any } | undefined;
}

export interface IDPSAMLAccessInformation {
  /** The SAML assertion returned by the identity provider. */
  assertion: Buffer;
}

export interface IDPAuthenticator {
  /** ID of the identity provider */
  idpId: string;
  /** ID of the user provided by the identity provider */
  userId: string;
  /** Username of the user provided by the identity provider. */
  userName: string;
}

function createBaseAuthenticators(): Authenticators {
  return {
    usernames: [],
    password: undefined,
    webAuthN: [],
    totps: [],
    otpSms: [],
    otpEmail: [],
    authenticationKeys: [],
    identityProviders: [],
  };
}

export const Authenticators = {
  encode(message: Authenticators, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.usernames) {
      Username.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.password !== undefined) {
      Password.encode(message.password, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.webAuthN) {
      WebAuthN.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.totps) {
      TOTP.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.otpSms) {
      OTPSMS.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.otpEmail) {
      OTPEmail.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    for (const v of message.authenticationKeys) {
      AuthenticationKey.encode(v!, writer.uint32(58).fork()).ldelim();
    }
    for (const v of message.identityProviders) {
      IdentityProvider.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Authenticators {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticators();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.usernames.push(Username.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.password = Password.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.webAuthN.push(WebAuthN.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.totps.push(TOTP.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.otpSms.push(OTPSMS.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.otpEmail.push(OTPEmail.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.authenticationKeys.push(AuthenticationKey.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.identityProviders.push(IdentityProvider.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Authenticators {
    return {
      usernames: Array.isArray(object?.usernames) ? object.usernames.map((e: any) => Username.fromJSON(e)) : [],
      password: isSet(object.password) ? Password.fromJSON(object.password) : undefined,
      webAuthN: Array.isArray(object?.webAuthN) ? object.webAuthN.map((e: any) => WebAuthN.fromJSON(e)) : [],
      totps: Array.isArray(object?.totps) ? object.totps.map((e: any) => TOTP.fromJSON(e)) : [],
      otpSms: Array.isArray(object?.otpSms) ? object.otpSms.map((e: any) => OTPSMS.fromJSON(e)) : [],
      otpEmail: Array.isArray(object?.otpEmail) ? object.otpEmail.map((e: any) => OTPEmail.fromJSON(e)) : [],
      authenticationKeys: Array.isArray(object?.authenticationKeys)
        ? object.authenticationKeys.map((e: any) => AuthenticationKey.fromJSON(e))
        : [],
      identityProviders: Array.isArray(object?.identityProviders)
        ? object.identityProviders.map((e: any) => IdentityProvider.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Authenticators): unknown {
    const obj: any = {};
    if (message.usernames) {
      obj.usernames = message.usernames.map((e) => e ? Username.toJSON(e) : undefined);
    } else {
      obj.usernames = [];
    }
    message.password !== undefined && (obj.password = message.password ? Password.toJSON(message.password) : undefined);
    if (message.webAuthN) {
      obj.webAuthN = message.webAuthN.map((e) => e ? WebAuthN.toJSON(e) : undefined);
    } else {
      obj.webAuthN = [];
    }
    if (message.totps) {
      obj.totps = message.totps.map((e) => e ? TOTP.toJSON(e) : undefined);
    } else {
      obj.totps = [];
    }
    if (message.otpSms) {
      obj.otpSms = message.otpSms.map((e) => e ? OTPSMS.toJSON(e) : undefined);
    } else {
      obj.otpSms = [];
    }
    if (message.otpEmail) {
      obj.otpEmail = message.otpEmail.map((e) => e ? OTPEmail.toJSON(e) : undefined);
    } else {
      obj.otpEmail = [];
    }
    if (message.authenticationKeys) {
      obj.authenticationKeys = message.authenticationKeys.map((e) => e ? AuthenticationKey.toJSON(e) : undefined);
    } else {
      obj.authenticationKeys = [];
    }
    if (message.identityProviders) {
      obj.identityProviders = message.identityProviders.map((e) => e ? IdentityProvider.toJSON(e) : undefined);
    } else {
      obj.identityProviders = [];
    }
    return obj;
  },

  create(base?: DeepPartial<Authenticators>): Authenticators {
    return Authenticators.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Authenticators>): Authenticators {
    const message = createBaseAuthenticators();
    message.usernames = object.usernames?.map((e) => Username.fromPartial(e)) || [];
    message.password = (object.password !== undefined && object.password !== null)
      ? Password.fromPartial(object.password)
      : undefined;
    message.webAuthN = object.webAuthN?.map((e) => WebAuthN.fromPartial(e)) || [];
    message.totps = object.totps?.map((e) => TOTP.fromPartial(e)) || [];
    message.otpSms = object.otpSms?.map((e) => OTPSMS.fromPartial(e)) || [];
    message.otpEmail = object.otpEmail?.map((e) => OTPEmail.fromPartial(e)) || [];
    message.authenticationKeys = object.authenticationKeys?.map((e) => AuthenticationKey.fromPartial(e)) || [];
    message.identityProviders = object.identityProviders?.map((e) => IdentityProvider.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUsername(): Username {
  return { usernameId: "", username: "", isOrganizationSpecific: false };
}

export const Username = {
  encode(message: Username, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.usernameId !== "") {
      writer.uint32(10).string(message.usernameId);
    }
    if (message.username !== "") {
      writer.uint32(18).string(message.username);
    }
    if (message.isOrganizationSpecific === true) {
      writer.uint32(24).bool(message.isOrganizationSpecific);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Username {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUsername();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.usernameId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.username = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isOrganizationSpecific = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Username {
    return {
      usernameId: isSet(object.usernameId) ? String(object.usernameId) : "",
      username: isSet(object.username) ? String(object.username) : "",
      isOrganizationSpecific: isSet(object.isOrganizationSpecific) ? Boolean(object.isOrganizationSpecific) : false,
    };
  },

  toJSON(message: Username): unknown {
    const obj: any = {};
    message.usernameId !== undefined && (obj.usernameId = message.usernameId);
    message.username !== undefined && (obj.username = message.username);
    message.isOrganizationSpecific !== undefined && (obj.isOrganizationSpecific = message.isOrganizationSpecific);
    return obj;
  },

  create(base?: DeepPartial<Username>): Username {
    return Username.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Username>): Username {
    const message = createBaseUsername();
    message.usernameId = object.usernameId ?? "";
    message.username = object.username ?? "";
    message.isOrganizationSpecific = object.isOrganizationSpecific ?? false;
    return message;
  },
};

function createBaseSetUsername(): SetUsername {
  return { username: "", isOrganizationSpecific: false };
}

export const SetUsername = {
  encode(message: SetUsername, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.isOrganizationSpecific === true) {
      writer.uint32(16).bool(message.isOrganizationSpecific);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetUsername {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetUsername();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.username = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.isOrganizationSpecific = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetUsername {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      isOrganizationSpecific: isSet(object.isOrganizationSpecific) ? Boolean(object.isOrganizationSpecific) : false,
    };
  },

  toJSON(message: SetUsername): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.isOrganizationSpecific !== undefined && (obj.isOrganizationSpecific = message.isOrganizationSpecific);
    return obj;
  },

  create(base?: DeepPartial<SetUsername>): SetUsername {
    return SetUsername.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetUsername>): SetUsername {
    const message = createBaseSetUsername();
    message.username = object.username ?? "";
    message.isOrganizationSpecific = object.isOrganizationSpecific ?? false;
    return message;
  },
};

function createBasePassword(): Password {
  return { lastChanged: undefined };
}

export const Password = {
  encode(message: Password, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.lastChanged !== undefined) {
      Timestamp.encode(toTimestamp(message.lastChanged), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Password {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePassword();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.lastChanged = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Password {
    return { lastChanged: isSet(object.lastChanged) ? fromJsonTimestamp(object.lastChanged) : undefined };
  },

  toJSON(message: Password): unknown {
    const obj: any = {};
    message.lastChanged !== undefined && (obj.lastChanged = message.lastChanged.toISOString());
    return obj;
  },

  create(base?: DeepPartial<Password>): Password {
    return Password.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Password>): Password {
    const message = createBasePassword();
    message.lastChanged = object.lastChanged ?? undefined;
    return message;
  },
};

function createBaseWebAuthN(): WebAuthN {
  return { webAuthNId: "", name: "", isVerified: false, userVerified: false };
}

export const WebAuthN = {
  encode(message: WebAuthN, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.webAuthNId !== "") {
      writer.uint32(10).string(message.webAuthNId);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.isVerified === true) {
      writer.uint32(24).bool(message.isVerified);
    }
    if (message.userVerified === true) {
      writer.uint32(32).bool(message.userVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): WebAuthN {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseWebAuthN();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.webAuthNId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.userVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): WebAuthN {
    return {
      webAuthNId: isSet(object.webAuthNId) ? String(object.webAuthNId) : "",
      name: isSet(object.name) ? String(object.name) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
      userVerified: isSet(object.userVerified) ? Boolean(object.userVerified) : false,
    };
  },

  toJSON(message: WebAuthN): unknown {
    const obj: any = {};
    message.webAuthNId !== undefined && (obj.webAuthNId = message.webAuthNId);
    message.name !== undefined && (obj.name = message.name);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    message.userVerified !== undefined && (obj.userVerified = message.userVerified);
    return obj;
  },

  create(base?: DeepPartial<WebAuthN>): WebAuthN {
    return WebAuthN.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<WebAuthN>): WebAuthN {
    const message = createBaseWebAuthN();
    message.webAuthNId = object.webAuthNId ?? "";
    message.name = object.name ?? "";
    message.isVerified = object.isVerified ?? false;
    message.userVerified = object.userVerified ?? false;
    return message;
  },
};

function createBaseOTPSMS(): OTPSMS {
  return { otpSmsId: "", phone: "", isVerified: false };
}

export const OTPSMS = {
  encode(message: OTPSMS, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.otpSmsId !== "") {
      writer.uint32(10).string(message.otpSmsId);
    }
    if (message.phone !== "") {
      writer.uint32(18).string(message.phone);
    }
    if (message.isVerified === true) {
      writer.uint32(24).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OTPSMS {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOTPSMS();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.otpSmsId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.phone = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OTPSMS {
    return {
      otpSmsId: isSet(object.otpSmsId) ? String(object.otpSmsId) : "",
      phone: isSet(object.phone) ? String(object.phone) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: OTPSMS): unknown {
    const obj: any = {};
    message.otpSmsId !== undefined && (obj.otpSmsId = message.otpSmsId);
    message.phone !== undefined && (obj.phone = message.phone);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<OTPSMS>): OTPSMS {
    return OTPSMS.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OTPSMS>): OTPSMS {
    const message = createBaseOTPSMS();
    message.otpSmsId = object.otpSmsId ?? "";
    message.phone = object.phone ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseOTPEmail(): OTPEmail {
  return { otpEmailId: "", address: "", isVerified: false };
}

export const OTPEmail = {
  encode(message: OTPEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.otpEmailId !== "") {
      writer.uint32(10).string(message.otpEmailId);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    if (message.isVerified === true) {
      writer.uint32(24).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OTPEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOTPEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.otpEmailId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.address = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OTPEmail {
    return {
      otpEmailId: isSet(object.otpEmailId) ? String(object.otpEmailId) : "",
      address: isSet(object.address) ? String(object.address) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: OTPEmail): unknown {
    const obj: any = {};
    message.otpEmailId !== undefined && (obj.otpEmailId = message.otpEmailId);
    message.address !== undefined && (obj.address = message.address);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<OTPEmail>): OTPEmail {
    return OTPEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<OTPEmail>): OTPEmail {
    const message = createBaseOTPEmail();
    message.otpEmailId = object.otpEmailId ?? "";
    message.address = object.address ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseTOTP(): TOTP {
  return { totpId: "", name: "", isVerified: false };
}

export const TOTP = {
  encode(message: TOTP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.totpId !== "") {
      writer.uint32(10).string(message.totpId);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.isVerified === true) {
      writer.uint32(24).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TOTP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTOTP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.totpId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.name = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.isVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TOTP {
    return {
      totpId: isSet(object.totpId) ? String(object.totpId) : "",
      name: isSet(object.name) ? String(object.name) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: TOTP): unknown {
    const obj: any = {};
    message.totpId !== undefined && (obj.totpId = message.totpId);
    message.name !== undefined && (obj.name = message.name);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<TOTP>): TOTP {
    return TOTP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<TOTP>): TOTP {
    const message = createBaseTOTP();
    message.totpId = object.totpId ?? "";
    message.name = object.name ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseAuthenticationKey(): AuthenticationKey {
  return { authenticationKeyId: "", details: undefined, type: 0, expirationDate: undefined };
}

export const AuthenticationKey = {
  encode(message: AuthenticationKey, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authenticationKeyId !== "") {
      writer.uint32(10).string(message.authenticationKeyId);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    if (message.expirationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.expirationDate), writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticationKey {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticationKey();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.authenticationKeyId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.expirationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthenticationKey {
    return {
      authenticationKeyId: isSet(object.authenticationKeyId) ? String(object.authenticationKeyId) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      type: isSet(object.type) ? authNKeyTypeFromJSON(object.type) : 0,
      expirationDate: isSet(object.expirationDate) ? fromJsonTimestamp(object.expirationDate) : undefined,
    };
  },

  toJSON(message: AuthenticationKey): unknown {
    const obj: any = {};
    message.authenticationKeyId !== undefined && (obj.authenticationKeyId = message.authenticationKeyId);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.type !== undefined && (obj.type = authNKeyTypeToJSON(message.type));
    message.expirationDate !== undefined && (obj.expirationDate = message.expirationDate.toISOString());
    return obj;
  },

  create(base?: DeepPartial<AuthenticationKey>): AuthenticationKey {
    return AuthenticationKey.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AuthenticationKey>): AuthenticationKey {
    const message = createBaseAuthenticationKey();
    message.authenticationKeyId = object.authenticationKeyId ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.type = object.type ?? 0;
    message.expirationDate = object.expirationDate ?? undefined;
    return message;
  },
};

function createBaseIdentityProvider(): IdentityProvider {
  return { idpId: "", idpName: "", userId: "", username: "" };
}

export const IdentityProvider = {
  encode(message: IdentityProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.idpName !== "") {
      writer.uint32(26).string(message.idpName);
    }
    if (message.userId !== "") {
      writer.uint32(34).string(message.userId);
    }
    if (message.username !== "") {
      writer.uint32(42).string(message.username);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdentityProvider {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentityProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.idpName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.username = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IdentityProvider {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      idpName: isSet(object.idpName) ? String(object.idpName) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      username: isSet(object.username) ? String(object.username) : "",
    };
  },

  toJSON(message: IdentityProvider): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.idpName !== undefined && (obj.idpName = message.idpName);
    message.userId !== undefined && (obj.userId = message.userId);
    message.username !== undefined && (obj.username = message.username);
    return obj;
  },

  create(base?: DeepPartial<IdentityProvider>): IdentityProvider {
    return IdentityProvider.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IdentityProvider>): IdentityProvider {
    const message = createBaseIdentityProvider();
    message.idpId = object.idpId ?? "";
    message.idpName = object.idpName ?? "";
    message.userId = object.userId ?? "";
    message.username = object.username ?? "";
    return message;
  },
};

function createBaseSetAuthenticators(): SetAuthenticators {
  return { usernames: [], password: undefined };
}

export const SetAuthenticators = {
  encode(message: SetAuthenticators, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.usernames) {
      SetUsername.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.password !== undefined) {
      SetPassword.encode(message.password, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetAuthenticators {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetAuthenticators();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.usernames.push(SetUsername.decode(reader, reader.uint32()));
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.password = SetPassword.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetAuthenticators {
    return {
      usernames: Array.isArray(object?.usernames) ? object.usernames.map((e: any) => SetUsername.fromJSON(e)) : [],
      password: isSet(object.password) ? SetPassword.fromJSON(object.password) : undefined,
    };
  },

  toJSON(message: SetAuthenticators): unknown {
    const obj: any = {};
    if (message.usernames) {
      obj.usernames = message.usernames.map((e) => e ? SetUsername.toJSON(e) : undefined);
    } else {
      obj.usernames = [];
    }
    message.password !== undefined &&
      (obj.password = message.password ? SetPassword.toJSON(message.password) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetAuthenticators>): SetAuthenticators {
    return SetAuthenticators.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetAuthenticators>): SetAuthenticators {
    const message = createBaseSetAuthenticators();
    message.usernames = object.usernames?.map((e) => SetUsername.fromPartial(e)) || [];
    message.password = (object.password !== undefined && object.password !== null)
      ? SetPassword.fromPartial(object.password)
      : undefined;
    return message;
  },
};

function createBaseSetPassword(): SetPassword {
  return { password: undefined, hash: undefined, changeRequired: false };
}

export const SetPassword = {
  encode(message: SetPassword, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.password !== undefined) {
      writer.uint32(10).string(message.password);
    }
    if (message.hash !== undefined) {
      writer.uint32(18).string(message.hash);
    }
    if (message.changeRequired === true) {
      writer.uint32(24).bool(message.changeRequired);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetPassword {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetPassword();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.password = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.hash = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.changeRequired = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetPassword {
    return {
      password: isSet(object.password) ? String(object.password) : undefined,
      hash: isSet(object.hash) ? String(object.hash) : undefined,
      changeRequired: isSet(object.changeRequired) ? Boolean(object.changeRequired) : false,
    };
  },

  toJSON(message: SetPassword): unknown {
    const obj: any = {};
    message.password !== undefined && (obj.password = message.password);
    message.hash !== undefined && (obj.hash = message.hash);
    message.changeRequired !== undefined && (obj.changeRequired = message.changeRequired);
    return obj;
  },

  create(base?: DeepPartial<SetPassword>): SetPassword {
    return SetPassword.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetPassword>): SetPassword {
    const message = createBaseSetPassword();
    message.password = object.password ?? undefined;
    message.hash = object.hash ?? undefined;
    message.changeRequired = object.changeRequired ?? false;
    return message;
  },
};

function createBaseSendPasswordResetEmail(): SendPasswordResetEmail {
  return { urlTemplate: undefined };
}

export const SendPasswordResetEmail = {
  encode(message: SendPasswordResetEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.urlTemplate !== undefined) {
      writer.uint32(18).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendPasswordResetEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendPasswordResetEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag != 18) {
            break;
          }

          message.urlTemplate = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SendPasswordResetEmail {
    return { urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined };
  },

  toJSON(message: SendPasswordResetEmail): unknown {
    const obj: any = {};
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<SendPasswordResetEmail>): SendPasswordResetEmail {
    return SendPasswordResetEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendPasswordResetEmail>): SendPasswordResetEmail {
    const message = createBaseSendPasswordResetEmail();
    message.urlTemplate = object.urlTemplate ?? undefined;
    return message;
  },
};

function createBaseSendPasswordResetSMS(): SendPasswordResetSMS {
  return {};
}

export const SendPasswordResetSMS = {
  encode(_: SendPasswordResetSMS, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendPasswordResetSMS {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendPasswordResetSMS();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): SendPasswordResetSMS {
    return {};
  },

  toJSON(_: SendPasswordResetSMS): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<SendPasswordResetSMS>): SendPasswordResetSMS {
    return SendPasswordResetSMS.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<SendPasswordResetSMS>): SendPasswordResetSMS {
    const message = createBaseSendPasswordResetSMS();
    return message;
  },
};

function createBaseReturnPasswordResetCode(): ReturnPasswordResetCode {
  return {};
}

export const ReturnPasswordResetCode = {
  encode(_: ReturnPasswordResetCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnPasswordResetCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnPasswordResetCode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ReturnPasswordResetCode {
    return {};
  },

  toJSON(_: ReturnPasswordResetCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnPasswordResetCode>): ReturnPasswordResetCode {
    return ReturnPasswordResetCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnPasswordResetCode>): ReturnPasswordResetCode {
    const message = createBaseReturnPasswordResetCode();
    return message;
  },
};

function createBaseAuthenticatorRegistrationCode(): AuthenticatorRegistrationCode {
  return { id: "", code: "" };
}

export const AuthenticatorRegistrationCode = {
  encode(message: AuthenticatorRegistrationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.code !== "") {
      writer.uint32(18).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticatorRegistrationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticatorRegistrationCode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.code = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthenticatorRegistrationCode {
    return { id: isSet(object.id) ? String(object.id) : "", code: isSet(object.code) ? String(object.code) : "" };
  },

  toJSON(message: AuthenticatorRegistrationCode): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<AuthenticatorRegistrationCode>): AuthenticatorRegistrationCode {
    return AuthenticatorRegistrationCode.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AuthenticatorRegistrationCode>): AuthenticatorRegistrationCode {
    const message = createBaseAuthenticatorRegistrationCode();
    message.id = object.id ?? "";
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseSendWebAuthNRegistrationLink(): SendWebAuthNRegistrationLink {
  return { urlTemplate: undefined };
}

export const SendWebAuthNRegistrationLink = {
  encode(message: SendWebAuthNRegistrationLink, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.urlTemplate !== undefined) {
      writer.uint32(10).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendWebAuthNRegistrationLink {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendWebAuthNRegistrationLink();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.urlTemplate = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SendWebAuthNRegistrationLink {
    return { urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined };
  },

  toJSON(message: SendWebAuthNRegistrationLink): unknown {
    const obj: any = {};
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<SendWebAuthNRegistrationLink>): SendWebAuthNRegistrationLink {
    return SendWebAuthNRegistrationLink.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendWebAuthNRegistrationLink>): SendWebAuthNRegistrationLink {
    const message = createBaseSendWebAuthNRegistrationLink();
    message.urlTemplate = object.urlTemplate ?? undefined;
    return message;
  },
};

function createBaseReturnWebAuthNRegistrationCode(): ReturnWebAuthNRegistrationCode {
  return {};
}

export const ReturnWebAuthNRegistrationCode = {
  encode(_: ReturnWebAuthNRegistrationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnWebAuthNRegistrationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnWebAuthNRegistrationCode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): ReturnWebAuthNRegistrationCode {
    return {};
  },

  toJSON(_: ReturnWebAuthNRegistrationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnWebAuthNRegistrationCode>): ReturnWebAuthNRegistrationCode {
    return ReturnWebAuthNRegistrationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnWebAuthNRegistrationCode>): ReturnWebAuthNRegistrationCode {
    const message = createBaseReturnWebAuthNRegistrationCode();
    return message;
  },
};

function createBaseRedirectURLs(): RedirectURLs {
  return { successUrl: "", failureUrl: "" };
}

export const RedirectURLs = {
  encode(message: RedirectURLs, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.successUrl !== "") {
      writer.uint32(10).string(message.successUrl);
    }
    if (message.failureUrl !== "") {
      writer.uint32(18).string(message.failureUrl);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RedirectURLs {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRedirectURLs();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.successUrl = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.failureUrl = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RedirectURLs {
    return {
      successUrl: isSet(object.successUrl) ? String(object.successUrl) : "",
      failureUrl: isSet(object.failureUrl) ? String(object.failureUrl) : "",
    };
  },

  toJSON(message: RedirectURLs): unknown {
    const obj: any = {};
    message.successUrl !== undefined && (obj.successUrl = message.successUrl);
    message.failureUrl !== undefined && (obj.failureUrl = message.failureUrl);
    return obj;
  },

  create(base?: DeepPartial<RedirectURLs>): RedirectURLs {
    return RedirectURLs.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RedirectURLs>): RedirectURLs {
    const message = createBaseRedirectURLs();
    message.successUrl = object.successUrl ?? "";
    message.failureUrl = object.failureUrl ?? "";
    return message;
  },
};

function createBaseLDAPCredentials(): LDAPCredentials {
  return { username: "", password: "" };
}

export const LDAPCredentials = {
  encode(message: LDAPCredentials, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.username !== "") {
      writer.uint32(10).string(message.username);
    }
    if (message.password !== "") {
      writer.uint32(18).string(message.password);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LDAPCredentials {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLDAPCredentials();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.username = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.password = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LDAPCredentials {
    return {
      username: isSet(object.username) ? String(object.username) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: LDAPCredentials): unknown {
    const obj: any = {};
    message.username !== undefined && (obj.username = message.username);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  create(base?: DeepPartial<LDAPCredentials>): LDAPCredentials {
    return LDAPCredentials.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LDAPCredentials>): LDAPCredentials {
    const message = createBaseLDAPCredentials();
    message.username = object.username ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseIdentityProviderIntent(): IdentityProviderIntent {
  return { idpIntentId: "", idpIntentToken: "", userId: undefined };
}

export const IdentityProviderIntent = {
  encode(message: IdentityProviderIntent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpIntentId !== "") {
      writer.uint32(10).string(message.idpIntentId);
    }
    if (message.idpIntentToken !== "") {
      writer.uint32(18).string(message.idpIntentToken);
    }
    if (message.userId !== undefined) {
      writer.uint32(26).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdentityProviderIntent {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentityProviderIntent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpIntentId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpIntentToken = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.userId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IdentityProviderIntent {
    return {
      idpIntentId: isSet(object.idpIntentId) ? String(object.idpIntentId) : "",
      idpIntentToken: isSet(object.idpIntentToken) ? String(object.idpIntentToken) : "",
      userId: isSet(object.userId) ? String(object.userId) : undefined,
    };
  },

  toJSON(message: IdentityProviderIntent): unknown {
    const obj: any = {};
    message.idpIntentId !== undefined && (obj.idpIntentId = message.idpIntentId);
    message.idpIntentToken !== undefined && (obj.idpIntentToken = message.idpIntentToken);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<IdentityProviderIntent>): IdentityProviderIntent {
    return IdentityProviderIntent.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IdentityProviderIntent>): IdentityProviderIntent {
    const message = createBaseIdentityProviderIntent();
    message.idpIntentId = object.idpIntentId ?? "";
    message.idpIntentToken = object.idpIntentToken ?? "";
    message.userId = object.userId ?? undefined;
    return message;
  },
};

function createBaseIDPInformation(): IDPInformation {
  return {
    idpId: "",
    userId: "",
    userName: "",
    rawInformation: undefined,
    oauth: undefined,
    ldap: undefined,
    saml: undefined,
  };
}

export const IDPInformation = {
  encode(message: IDPInformation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(26).string(message.userName);
    }
    if (message.rawInformation !== undefined) {
      Struct.encode(Struct.wrap(message.rawInformation), writer.uint32(34).fork()).ldelim();
    }
    if (message.oauth !== undefined) {
      IDPOAuthAccessInformation.encode(message.oauth, writer.uint32(42).fork()).ldelim();
    }
    if (message.ldap !== undefined) {
      IDPLDAPAccessInformation.encode(message.ldap, writer.uint32(50).fork()).ldelim();
    }
    if (message.saml !== undefined) {
      IDPSAMLAccessInformation.encode(message.saml, writer.uint32(58).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPInformation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPInformation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.userName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.rawInformation = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.oauth = IDPOAuthAccessInformation.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.ldap = IDPLDAPAccessInformation.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.saml = IDPSAMLAccessInformation.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPInformation {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
      rawInformation: isObject(object.rawInformation) ? object.rawInformation : undefined,
      oauth: isSet(object.oauth) ? IDPOAuthAccessInformation.fromJSON(object.oauth) : undefined,
      ldap: isSet(object.ldap) ? IDPLDAPAccessInformation.fromJSON(object.ldap) : undefined,
      saml: isSet(object.saml) ? IDPSAMLAccessInformation.fromJSON(object.saml) : undefined,
    };
  },

  toJSON(message: IDPInformation): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userName !== undefined && (obj.userName = message.userName);
    message.rawInformation !== undefined && (obj.rawInformation = message.rawInformation);
    message.oauth !== undefined &&
      (obj.oauth = message.oauth ? IDPOAuthAccessInformation.toJSON(message.oauth) : undefined);
    message.ldap !== undefined && (obj.ldap = message.ldap ? IDPLDAPAccessInformation.toJSON(message.ldap) : undefined);
    message.saml !== undefined && (obj.saml = message.saml ? IDPSAMLAccessInformation.toJSON(message.saml) : undefined);
    return obj;
  },

  create(base?: DeepPartial<IDPInformation>): IDPInformation {
    return IDPInformation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPInformation>): IDPInformation {
    const message = createBaseIDPInformation();
    message.idpId = object.idpId ?? "";
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    message.rawInformation = object.rawInformation ?? undefined;
    message.oauth = (object.oauth !== undefined && object.oauth !== null)
      ? IDPOAuthAccessInformation.fromPartial(object.oauth)
      : undefined;
    message.ldap = (object.ldap !== undefined && object.ldap !== null)
      ? IDPLDAPAccessInformation.fromPartial(object.ldap)
      : undefined;
    message.saml = (object.saml !== undefined && object.saml !== null)
      ? IDPSAMLAccessInformation.fromPartial(object.saml)
      : undefined;
    return message;
  },
};

function createBaseIDPOAuthAccessInformation(): IDPOAuthAccessInformation {
  return { accessToken: "", idToken: undefined };
}

export const IDPOAuthAccessInformation = {
  encode(message: IDPOAuthAccessInformation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.idToken !== undefined) {
      writer.uint32(18).string(message.idToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPOAuthAccessInformation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPOAuthAccessInformation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.accessToken = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idToken = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPOAuthAccessInformation {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      idToken: isSet(object.idToken) ? String(object.idToken) : undefined,
    };
  },

  toJSON(message: IDPOAuthAccessInformation): unknown {
    const obj: any = {};
    message.accessToken !== undefined && (obj.accessToken = message.accessToken);
    message.idToken !== undefined && (obj.idToken = message.idToken);
    return obj;
  },

  create(base?: DeepPartial<IDPOAuthAccessInformation>): IDPOAuthAccessInformation {
    return IDPOAuthAccessInformation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPOAuthAccessInformation>): IDPOAuthAccessInformation {
    const message = createBaseIDPOAuthAccessInformation();
    message.accessToken = object.accessToken ?? "";
    message.idToken = object.idToken ?? undefined;
    return message;
  },
};

function createBaseIDPLDAPAccessInformation(): IDPLDAPAccessInformation {
  return { attributes: undefined };
}

export const IDPLDAPAccessInformation = {
  encode(message: IDPLDAPAccessInformation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.attributes !== undefined) {
      Struct.encode(Struct.wrap(message.attributes), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPLDAPAccessInformation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPLDAPAccessInformation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.attributes = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPLDAPAccessInformation {
    return { attributes: isObject(object.attributes) ? object.attributes : undefined };
  },

  toJSON(message: IDPLDAPAccessInformation): unknown {
    const obj: any = {};
    message.attributes !== undefined && (obj.attributes = message.attributes);
    return obj;
  },

  create(base?: DeepPartial<IDPLDAPAccessInformation>): IDPLDAPAccessInformation {
    return IDPLDAPAccessInformation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPLDAPAccessInformation>): IDPLDAPAccessInformation {
    const message = createBaseIDPLDAPAccessInformation();
    message.attributes = object.attributes ?? undefined;
    return message;
  },
};

function createBaseIDPSAMLAccessInformation(): IDPSAMLAccessInformation {
  return { assertion: Buffer.alloc(0) };
}

export const IDPSAMLAccessInformation = {
  encode(message: IDPSAMLAccessInformation, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.assertion.length !== 0) {
      writer.uint32(10).bytes(message.assertion);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPSAMLAccessInformation {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPSAMLAccessInformation();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.assertion = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPSAMLAccessInformation {
    return { assertion: isSet(object.assertion) ? Buffer.from(bytesFromBase64(object.assertion)) : Buffer.alloc(0) };
  },

  toJSON(message: IDPSAMLAccessInformation): unknown {
    const obj: any = {};
    message.assertion !== undefined &&
      (obj.assertion = base64FromBytes(message.assertion !== undefined ? message.assertion : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<IDPSAMLAccessInformation>): IDPSAMLAccessInformation {
    return IDPSAMLAccessInformation.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPSAMLAccessInformation>): IDPSAMLAccessInformation {
    const message = createBaseIDPSAMLAccessInformation();
    message.assertion = object.assertion ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseIDPAuthenticator(): IDPAuthenticator {
  return { idpId: "", userId: "", userName: "" };
}

export const IDPAuthenticator = {
  encode(message: IDPAuthenticator, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.userName !== "") {
      writer.uint32(26).string(message.userName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IDPAuthenticator {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIDPAuthenticator();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.userName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IDPAuthenticator {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      userName: isSet(object.userName) ? String(object.userName) : "",
    };
  },

  toJSON(message: IDPAuthenticator): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userName !== undefined && (obj.userName = message.userName);
    return obj;
  },

  create(base?: DeepPartial<IDPAuthenticator>): IDPAuthenticator {
    return IDPAuthenticator.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<IDPAuthenticator>): IDPAuthenticator {
    const message = createBaseIDPAuthenticator();
    message.idpId = object.idpId ?? "";
    message.userId = object.userId ?? "";
    message.userName = object.userName ?? "";
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
