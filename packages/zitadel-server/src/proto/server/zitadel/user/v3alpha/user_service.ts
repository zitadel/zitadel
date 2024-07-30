/* eslint-disable */
import type { CallContext, CallOptions } from "nice-grpc-common";
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../google/protobuf/struct";
import { Details, ListDetails, ListQuery, Organization } from "../../object/v2beta/object";
import {
  AuthenticatorRegistrationCode,
  IdentityProviderIntent,
  IDPAuthenticator,
  IDPInformation,
  LDAPCredentials,
  RedirectURLs,
  ReturnPasswordResetCode,
  ReturnWebAuthNRegistrationCode,
  SendPasswordResetEmail,
  SendPasswordResetSMS,
  SendWebAuthNRegistrationLink,
  SetAuthenticators,
  SetPassword,
  SetUsername,
  WebAuthNAuthenticatorType,
  webAuthNAuthenticatorTypeFromJSON,
  webAuthNAuthenticatorTypeToJSON,
} from "./authenticator";
import {
  ReturnEmailVerificationCode,
  ReturnPhoneVerificationCode,
  SendEmailVerificationCode,
  SendPhoneVerificationCode,
  SetContact,
  SetEmail,
  SetPhone,
} from "./communication";
import { FieldName, fieldNameFromJSON, fieldNameToJSON, SearchQuery } from "./query";
import { User } from "./user";

export const protobufPackage = "zitadel.user.v3alpha";

export interface ListUsersRequest {
  /** list limitations and ordering. */
  query:
    | ListQuery
    | undefined;
  /** the field the result is sorted. */
  sortingColumn: FieldName;
  /** Define the criteria to query for. */
  queries: SearchQuery[];
}

export interface ListUsersResponse {
  /** Details provides information about the returned result including total amount found. */
  details:
    | ListDetails
    | undefined;
  /** States by which field the results are sorted. */
  sortingColumn: FieldName;
  /** The result contains the user schemas, which matched the queries. */
  result: User[];
}

export interface GetUserByIDRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface GetUserByIDResponse {
  user: User | undefined;
}

export interface CreateUserRequest {
  /** Optionally set a unique identifier of the user. If unset, ZITADEL will take care of it. */
  userId?:
    | string
    | undefined;
  /** Set the organization the user belongs to. */
  organization:
    | Organization
    | undefined;
  /** Set the initial authenticators of the user. */
  authenticators:
    | SetAuthenticators
    | undefined;
  /** Set the contact information (email, phone) for the user. */
  contact:
    | SetContact
    | undefined;
  /** Define the schema the user's data schema by providing it's ID. */
  schemaId: string;
  /** Provide data about the user. It will be validated based on the specified schema. */
  data: { [key: string]: any } | undefined;
}

export interface CreateUserResponse {
  userId: string;
  details:
    | Details
    | undefined;
  /** The email code will be set if a contact email was set with a return_code verification option. */
  emailCode?:
    | string
    | undefined;
  /** The phone code will be set if a contact phone was set with a return_code verification option. */
  phoneCode?: string | undefined;
}

export interface UpdateUserRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Add or update the contact information (email, phone) for the user if needed. */
  contact?:
    | SetContact
    | undefined;
  /** Change the schema the user's data schema by providing it's ID if needed. */
  schemaId?:
    | string
    | undefined;
  /** Update the user data if needed. It will be validated based on the specified schema. */
  data?: { [key: string]: any } | undefined;
}

export interface UpdateUserResponse {
  details:
    | Details
    | undefined;
  /** The email code will be set if a contact email was set with a return_code verification option. */
  emailCode?:
    | string
    | undefined;
  /** The phone code will be set if a contact phone was set with a return_code verification option. */
  phoneCode?: string | undefined;
}

export interface DeactivateUserRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface DeactivateUserResponse {
  details: Details | undefined;
}

export interface ReactivateUserRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface ReactivateUserResponse {
  details: Details | undefined;
}

export interface LockUserRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface LockUserResponse {
  details: Details | undefined;
}

export interface UnlockUserRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface UnlockUserResponse {
  details: Details | undefined;
}

export interface DeleteUserRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface DeleteUserResponse {
  details: Details | undefined;
}

export interface SetContactEmailRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the user's contact email and it's verification state. */
  email: SetEmail | undefined;
}

export interface SetContactEmailResponse {
  details:
    | Details
    | undefined;
  /** The verification code will be set if a contact email was set with a return_code verification option. */
  verificationCode?: string | undefined;
}

export interface VerifyContactEmailRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the verification code generated during the set contact email request. */
  verificationCode: string;
}

export interface VerifyContactEmailResponse {
  details: Details | undefined;
}

export interface ResendContactEmailCodeRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Let ZITADEL send the link to the user via email. */
  sendCode?:
    | SendEmailVerificationCode
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?: ReturnEmailVerificationCode | undefined;
}

export interface ResendContactEmailCodeResponse {
  details:
    | Details
    | undefined;
  /** in case the verification was set to return_code, the code will be returned. */
  verificationCode?: string | undefined;
}

export interface SetContactPhoneRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the user's contact phone and it's verification state. */
  phone: SetPhone | undefined;
}

export interface SetContactPhoneResponse {
  details:
    | Details
    | undefined;
  /** The phone verification code will be set if a contact phone was set with a return_code verification option. */
  emailCode?: string | undefined;
}

export interface VerifyContactPhoneRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the verification code generated during the set contact phone request. */
  verificationCode: string;
}

export interface VerifyContactPhoneResponse {
  details: Details | undefined;
}

export interface ResendContactPhoneCodeRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Let ZITADEL send the link to the user via SMS. */
  sendCode?:
    | SendPhoneVerificationCode
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?: ReturnPhoneVerificationCode | undefined;
}

export interface ResendContactPhoneCodeResponse {
  details:
    | Details
    | undefined;
  /** in case the verification was set to return_code, the code will be returned. */
  verificationCode?: string | undefined;
}

export interface AddUsernameRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the user's new username. */
  username: SetUsername | undefined;
}

export interface AddUsernameResponse {
  details:
    | Details
    | undefined;
  /** unique identifier of the username. */
  usernameId: string;
}

export interface RemoveUsernameRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the username. */
  usernameId: string;
}

export interface RemoveUsernameResponse {
  details: Details | undefined;
}

export interface SetPasswordRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Provide the new password (in plain text or as hash). */
  newPassword:
    | SetPassword
    | undefined;
  /** Provide the current password to verify you're allowed to change the password. */
  currentPassword?:
    | string
    | undefined;
  /** Or provider the verification code generated during password reset request. */
  verificationCode?: string | undefined;
}

export interface SetPasswordResponse {
  details: Details | undefined;
}

export interface RequestPasswordResetRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Let ZITADEL send the link to the user via email. */
  sendEmail?:
    | SendPasswordResetEmail
    | undefined;
  /** Let ZITADEL send the link to the user via SMS. */
  sendSms?:
    | SendPasswordResetSMS
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?: ReturnPasswordResetCode | undefined;
}

export interface RequestPasswordResetResponse {
  details:
    | Details
    | undefined;
  /** In case the medium was set to return_code, the code will be returned. */
  verificationCode?: string | undefined;
}

export interface StartWebAuthNRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Domain on which the user currently is or will be authenticated. */
  domain: string;
  /**
   * Optionally specify the authenticator type of the passkey device (platform or cross-platform).
   * If none is provided, both values are allowed.
   */
  authenticatorType: WebAuthNAuthenticatorType;
  /**
   * Optionally provide a one time code generated by ZITADEL.
   * This is required to start the passkey registration without user authentication.
   */
  code?: AuthenticatorRegistrationCode | undefined;
}

export interface StartWebAuthNRegistrationResponse {
  details:
    | Details
    | undefined;
  /** unique identifier of the WebAuthN registration. */
  webAuthNId: string;
  /**
   * Options for Credential Creation (dictionary PublicKeyCredentialCreationOptions).
   * Generated helper methods transform the field to JSON, for use in a WebauthN client.
   * See also:  https://www.w3.org/TR/webauthn/#dictdef-publickeycredentialcreationoptions
   */
  publicKeyCredentialCreationOptions: { [key: string]: any } | undefined;
}

export interface VerifyWebAuthNRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the WebAuthN registration, which was returned in the start webauthn registration. */
  webAuthNId: string;
  /**
   * PublicKeyCredential Interface.
   * Generated helper methods populate the field from JSON created by a WebAuthN client.
   * See also:  https://www.w3.org/TR/webauthn/#publickeycredential
   */
  publicKeyCredential:
    | { [key: string]: any }
    | undefined;
  /** Provide a name for the WebAuthN device. This will help identify it in the future. */
  webAuthNName: string;
}

export interface VerifyWebAuthNRegistrationResponse {
  details: Details | undefined;
}

export interface CreateWebAuthNRegistrationLinkRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Let ZITADEL send the link to the user via email. */
  sendLink?:
    | SendWebAuthNRegistrationLink
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?: ReturnWebAuthNRegistrationCode | undefined;
}

export interface CreateWebAuthNRegistrationLinkResponse {
  details:
    | Details
    | undefined;
  /** In case the medium was set to return_code, the code will be returned. */
  code?: AuthenticatorRegistrationCode | undefined;
}

export interface RemoveWebAuthNAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the WebAuthN authenticator. */
  webAuthNId: string;
}

export interface RemoveWebAuthNAuthenticatorResponse {
  details: Details | undefined;
}

export interface StartTOTPRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
}

export interface StartTOTPRegistrationResponse {
  details:
    | Details
    | undefined;
  /** unique identifier of the TOTP registration. */
  totpId: string;
  /** The TOTP URI, which can be used to create a QR Code for scanning with an authenticator app. */
  uri: string;
  /** The TOTP secret, which can be used for manually adding in an authenticator app. */
  secret: string;
}

export interface VerifyTOTPRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the TOTP registration, which was returned in the start TOTP registration. */
  totpId: string;
  /** Code generated by TOTP app or device. */
  code: string;
}

export interface VerifyTOTPRegistrationResponse {
  details: Details | undefined;
}

export interface RemoveTOTPAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the TOTP authenticator. */
  totpId: string;
}

export interface RemoveTOTPAuthenticatorResponse {
  details: Details | undefined;
}

export interface AddOTPSMSAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the user's phone for the OTP SMS authenticator and it's verification state. */
  phone: SetPhone | undefined;
}

export interface AddOTPSMSAuthenticatorResponse {
  details:
    | Details
    | undefined;
  /** unique identifier of the OTP SMS registration. */
  otpSmsId: string;
  /** The OTP verification code will be set if a phone was set with a return_code verification option. */
  verificationCode?: string | undefined;
}

export interface VerifyOTPSMSRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the OTP SMS registration, which was returned in the add OTP SMS. */
  otpSmsId: string;
  /** Set the verification code generated during the add OTP SMS request. */
  code: string;
}

export interface VerifyOTPSMSRegistrationResponse {
  details: Details | undefined;
}

export interface RemoveOTPSMSAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the OTP SMS authenticator. */
  otpSmsId: string;
}

export interface RemoveOTPSMSAuthenticatorResponse {
  details: Details | undefined;
}

export interface AddOTPEmailAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** Set the user's email for the OTP Email authenticator and it's verification state. */
  email: SetEmail | undefined;
}

export interface AddOTPEmailAuthenticatorResponse {
  details:
    | Details
    | undefined;
  /** unique identifier of the OTP Email registration. */
  otpEmailId: string;
  /** The OTP verification code will be set if a email was set with a return_code verification option. */
  verificationCode?: string | undefined;
}

export interface VerifyOTPEmailRegistrationRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the OTP Email registration, which was returned in the add OTP Email. */
  otpEmailId: string;
  /** Set the verification code generated during the add OTP Email request. */
  code: string;
}

export interface VerifyOTPEmailRegistrationResponse {
  details: Details | undefined;
}

export interface RemoveOTPEmailAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the OTP Email authenticator. */
  otpEmailId: string;
}

export interface RemoveOTPEmailAuthenticatorResponse {
  details: Details | undefined;
}

export interface StartIdentityProviderIntentRequest {
  /** ID of an existing identity provider (IDP). */
  idpId: string;
  urls?: RedirectURLs | undefined;
  ldap?: LDAPCredentials | undefined;
}

export interface StartIdentityProviderIntentResponse {
  details:
    | Details
    | undefined;
  /** The authentication URL to which the client should redirect. */
  authUrl?:
    | string
    | undefined;
  /**
   * The Start Intent directly succeeded and returned the IDP Intent.
   * Further information can be retrieved by using the retrieve identity provider intent request.
   */
  idpIntent?:
    | IdentityProviderIntent
    | undefined;
  /** The HTML form with the embedded POST call information to render and execute. */
  postForm?: Buffer | undefined;
}

export interface RetrieveIdentityProviderIntentRequest {
  /** ID of the identity provider (IDP) intent, previously returned on the success response of the start identity provider intent. */
  idpIntentId: string;
  /** Token of the identity provider (IDP) intent, previously returned on the success response of the start identity provider intent. */
  idpIntentToken: string;
}

export interface RetrieveIdentityProviderIntentResponse {
  details:
    | Details
    | undefined;
  /**
   * Information returned by the identity provider (IDP) such as the identification of the user
   * and detailed / profile information.
   */
  idpInformation:
    | IDPInformation
    | undefined;
  /** If the user was already federated and linked to a ZITADEL user, it's id will be returned. */
  userId?: string | undefined;
}

export interface AddIDPAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  idpAuthenticator: IDPAuthenticator | undefined;
}

export interface AddIDPAuthenticatorResponse {
  details: Details | undefined;
}

export interface RemoveIDPAuthenticatorRequest {
  /** unique identifier of the user. */
  userId: string;
  /** unique identifier of the identity provider (IDP) authenticator. */
  idpId: string;
}

export interface RemoveIDPAuthenticatorResponse {
  details: Details | undefined;
}

function createBaseListUsersRequest(): ListUsersRequest {
  return { query: undefined, sortingColumn: 0, queries: [] };
}

export const ListUsersRequest = {
  encode(message: ListUsersRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.query !== undefined) {
      ListQuery.encode(message.query, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.queries) {
      SearchQuery.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListUsersRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListUsersRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.query = ListQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.sortingColumn = reader.int32() as any;
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.queries.push(SearchQuery.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListUsersRequest {
    return {
      query: isSet(object.query) ? ListQuery.fromJSON(object.query) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? fieldNameFromJSON(object.sortingColumn) : 0,
      queries: Array.isArray(object?.queries) ? object.queries.map((e: any) => SearchQuery.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListUsersRequest): unknown {
    const obj: any = {};
    message.query !== undefined && (obj.query = message.query ? ListQuery.toJSON(message.query) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = fieldNameToJSON(message.sortingColumn));
    if (message.queries) {
      obj.queries = message.queries.map((e) => e ? SearchQuery.toJSON(e) : undefined);
    } else {
      obj.queries = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListUsersRequest>): ListUsersRequest {
    return ListUsersRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListUsersRequest>): ListUsersRequest {
    const message = createBaseListUsersRequest();
    message.query = (object.query !== undefined && object.query !== null)
      ? ListQuery.fromPartial(object.query)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.queries = object.queries?.map((e) => SearchQuery.fromPartial(e)) || [];
    return message;
  },
};

function createBaseListUsersResponse(): ListUsersResponse {
  return { details: undefined, sortingColumn: 0, result: [] };
}

export const ListUsersResponse = {
  encode(message: ListUsersResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      ListDetails.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.sortingColumn !== 0) {
      writer.uint32(16).int32(message.sortingColumn);
    }
    for (const v of message.result) {
      User.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListUsersResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListUsersResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = ListDetails.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.sortingColumn = reader.int32() as any;
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.result.push(User.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListUsersResponse {
    return {
      details: isSet(object.details) ? ListDetails.fromJSON(object.details) : undefined,
      sortingColumn: isSet(object.sortingColumn) ? fieldNameFromJSON(object.sortingColumn) : 0,
      result: Array.isArray(object?.result) ? object.result.map((e: any) => User.fromJSON(e)) : [],
    };
  },

  toJSON(message: ListUsersResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? ListDetails.toJSON(message.details) : undefined);
    message.sortingColumn !== undefined && (obj.sortingColumn = fieldNameToJSON(message.sortingColumn));
    if (message.result) {
      obj.result = message.result.map((e) => e ? User.toJSON(e) : undefined);
    } else {
      obj.result = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ListUsersResponse>): ListUsersResponse {
    return ListUsersResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ListUsersResponse>): ListUsersResponse {
    const message = createBaseListUsersResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? ListDetails.fromPartial(object.details)
      : undefined;
    message.sortingColumn = object.sortingColumn ?? 0;
    message.result = object.result?.map((e) => User.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetUserByIDRequest(): GetUserByIDRequest {
  return { userId: "" };
}

export const GetUserByIDRequest = {
  encode(message: GetUserByIDRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserByIDRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserByIDRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): GetUserByIDRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: GetUserByIDRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<GetUserByIDRequest>): GetUserByIDRequest {
    return GetUserByIDRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserByIDRequest>): GetUserByIDRequest {
    const message = createBaseGetUserByIDRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseGetUserByIDResponse(): GetUserByIDResponse {
  return { user: undefined };
}

export const GetUserByIDResponse = {
  encode(message: GetUserByIDResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetUserByIDResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetUserByIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetUserByIDResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: GetUserByIDResponse): unknown {
    const obj: any = {};
    message.user !== undefined && (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GetUserByIDResponse>): GetUserByIDResponse {
    return GetUserByIDResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GetUserByIDResponse>): GetUserByIDResponse {
    const message = createBaseGetUserByIDResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseCreateUserRequest(): CreateUserRequest {
  return {
    userId: undefined,
    organization: undefined,
    authenticators: undefined,
    contact: undefined,
    schemaId: "",
    data: undefined,
  };
}

export const CreateUserRequest = {
  encode(message: CreateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== undefined) {
      writer.uint32(10).string(message.userId);
    }
    if (message.organization !== undefined) {
      Organization.encode(message.organization, writer.uint32(18).fork()).ldelim();
    }
    if (message.authenticators !== undefined) {
      SetAuthenticators.encode(message.authenticators, writer.uint32(26).fork()).ldelim();
    }
    if (message.contact !== undefined) {
      SetContact.encode(message.contact, writer.uint32(34).fork()).ldelim();
    }
    if (message.schemaId !== "") {
      writer.uint32(42).string(message.schemaId);
    }
    if (message.data !== undefined) {
      Struct.encode(Struct.wrap(message.data), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.organization = Organization.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.authenticators = SetAuthenticators.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.contact = SetContact.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.schemaId = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.data = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateUserRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : undefined,
      organization: isSet(object.organization) ? Organization.fromJSON(object.organization) : undefined,
      authenticators: isSet(object.authenticators) ? SetAuthenticators.fromJSON(object.authenticators) : undefined,
      contact: isSet(object.contact) ? SetContact.fromJSON(object.contact) : undefined,
      schemaId: isSet(object.schemaId) ? String(object.schemaId) : "",
      data: isObject(object.data) ? object.data : undefined,
    };
  },

  toJSON(message: CreateUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.organization !== undefined &&
      (obj.organization = message.organization ? Organization.toJSON(message.organization) : undefined);
    message.authenticators !== undefined &&
      (obj.authenticators = message.authenticators ? SetAuthenticators.toJSON(message.authenticators) : undefined);
    message.contact !== undefined && (obj.contact = message.contact ? SetContact.toJSON(message.contact) : undefined);
    message.schemaId !== undefined && (obj.schemaId = message.schemaId);
    message.data !== undefined && (obj.data = message.data);
    return obj;
  },

  create(base?: DeepPartial<CreateUserRequest>): CreateUserRequest {
    return CreateUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateUserRequest>): CreateUserRequest {
    const message = createBaseCreateUserRequest();
    message.userId = object.userId ?? undefined;
    message.organization = (object.organization !== undefined && object.organization !== null)
      ? Organization.fromPartial(object.organization)
      : undefined;
    message.authenticators = (object.authenticators !== undefined && object.authenticators !== null)
      ? SetAuthenticators.fromPartial(object.authenticators)
      : undefined;
    message.contact = (object.contact !== undefined && object.contact !== null)
      ? SetContact.fromPartial(object.contact)
      : undefined;
    message.schemaId = object.schemaId ?? "";
    message.data = object.data ?? undefined;
    return message;
  },
};

function createBaseCreateUserResponse(): CreateUserResponse {
  return { userId: "", details: undefined, emailCode: undefined, phoneCode: undefined };
}

export const CreateUserResponse = {
  encode(message: CreateUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.emailCode !== undefined) {
      writer.uint32(26).string(message.emailCode);
    }
    if (message.phoneCode !== undefined) {
      writer.uint32(34).string(message.phoneCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.emailCode = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.phoneCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateUserResponse {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      emailCode: isSet(object.emailCode) ? String(object.emailCode) : undefined,
      phoneCode: isSet(object.phoneCode) ? String(object.phoneCode) : undefined,
    };
  },

  toJSON(message: CreateUserResponse): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.emailCode !== undefined && (obj.emailCode = message.emailCode);
    message.phoneCode !== undefined && (obj.phoneCode = message.phoneCode);
    return obj;
  },

  create(base?: DeepPartial<CreateUserResponse>): CreateUserResponse {
    return CreateUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateUserResponse>): CreateUserResponse {
    const message = createBaseCreateUserResponse();
    message.userId = object.userId ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.emailCode = object.emailCode ?? undefined;
    message.phoneCode = object.phoneCode ?? undefined;
    return message;
  },
};

function createBaseUpdateUserRequest(): UpdateUserRequest {
  return { userId: "", contact: undefined, schemaId: undefined, data: undefined };
}

export const UpdateUserRequest = {
  encode(message: UpdateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.contact !== undefined) {
      SetContact.encode(message.contact, writer.uint32(34).fork()).ldelim();
    }
    if (message.schemaId !== undefined) {
      writer.uint32(42).string(message.schemaId);
    }
    if (message.data !== undefined) {
      Struct.encode(Struct.wrap(message.data), writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.contact = SetContact.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.schemaId = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.data = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateUserRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      contact: isSet(object.contact) ? SetContact.fromJSON(object.contact) : undefined,
      schemaId: isSet(object.schemaId) ? String(object.schemaId) : undefined,
      data: isObject(object.data) ? object.data : undefined,
    };
  },

  toJSON(message: UpdateUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.contact !== undefined && (obj.contact = message.contact ? SetContact.toJSON(message.contact) : undefined);
    message.schemaId !== undefined && (obj.schemaId = message.schemaId);
    message.data !== undefined && (obj.data = message.data);
    return obj;
  },

  create(base?: DeepPartial<UpdateUserRequest>): UpdateUserRequest {
    return UpdateUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateUserRequest>): UpdateUserRequest {
    const message = createBaseUpdateUserRequest();
    message.userId = object.userId ?? "";
    message.contact = (object.contact !== undefined && object.contact !== null)
      ? SetContact.fromPartial(object.contact)
      : undefined;
    message.schemaId = object.schemaId ?? undefined;
    message.data = object.data ?? undefined;
    return message;
  },
};

function createBaseUpdateUserResponse(): UpdateUserResponse {
  return { details: undefined, emailCode: undefined, phoneCode: undefined };
}

export const UpdateUserResponse = {
  encode(message: UpdateUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.emailCode !== undefined) {
      writer.uint32(26).string(message.emailCode);
    }
    if (message.phoneCode !== undefined) {
      writer.uint32(34).string(message.phoneCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.emailCode = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.phoneCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateUserResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      emailCode: isSet(object.emailCode) ? String(object.emailCode) : undefined,
      phoneCode: isSet(object.phoneCode) ? String(object.phoneCode) : undefined,
    };
  },

  toJSON(message: UpdateUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.emailCode !== undefined && (obj.emailCode = message.emailCode);
    message.phoneCode !== undefined && (obj.phoneCode = message.phoneCode);
    return obj;
  },

  create(base?: DeepPartial<UpdateUserResponse>): UpdateUserResponse {
    return UpdateUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UpdateUserResponse>): UpdateUserResponse {
    const message = createBaseUpdateUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.emailCode = object.emailCode ?? undefined;
    message.phoneCode = object.phoneCode ?? undefined;
    return message;
  },
};

function createBaseDeactivateUserRequest(): DeactivateUserRequest {
  return { userId: "" };
}

export const DeactivateUserRequest = {
  encode(message: DeactivateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): DeactivateUserRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: DeactivateUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<DeactivateUserRequest>): DeactivateUserRequest {
    return DeactivateUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateUserRequest>): DeactivateUserRequest {
    const message = createBaseDeactivateUserRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseDeactivateUserResponse(): DeactivateUserResponse {
  return { details: undefined };
}

export const DeactivateUserResponse = {
  encode(message: DeactivateUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeactivateUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeactivateUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeactivateUserResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeactivateUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeactivateUserResponse>): DeactivateUserResponse {
    return DeactivateUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeactivateUserResponse>): DeactivateUserResponse {
    const message = createBaseDeactivateUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseReactivateUserRequest(): ReactivateUserRequest {
  return { userId: "" };
}

export const ReactivateUserRequest = {
  encode(message: ReactivateUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): ReactivateUserRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: ReactivateUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<ReactivateUserRequest>): ReactivateUserRequest {
    return ReactivateUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateUserRequest>): ReactivateUserRequest {
    const message = createBaseReactivateUserRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseReactivateUserResponse(): ReactivateUserResponse {
  return { details: undefined };
}

export const ReactivateUserResponse = {
  encode(message: ReactivateUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReactivateUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReactivateUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReactivateUserResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: ReactivateUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ReactivateUserResponse>): ReactivateUserResponse {
    return ReactivateUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ReactivateUserResponse>): ReactivateUserResponse {
    const message = createBaseReactivateUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseLockUserRequest(): LockUserRequest {
  return { userId: "" };
}

export const LockUserRequest = {
  encode(message: LockUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LockUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLockUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): LockUserRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: LockUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<LockUserRequest>): LockUserRequest {
    return LockUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LockUserRequest>): LockUserRequest {
    const message = createBaseLockUserRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseLockUserResponse(): LockUserResponse {
  return { details: undefined };
}

export const LockUserResponse = {
  encode(message: LockUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LockUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLockUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): LockUserResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: LockUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<LockUserResponse>): LockUserResponse {
    return LockUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<LockUserResponse>): LockUserResponse {
    const message = createBaseLockUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseUnlockUserRequest(): UnlockUserRequest {
  return { userId: "" };
}

export const UnlockUserRequest = {
  encode(message: UnlockUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnlockUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnlockUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): UnlockUserRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: UnlockUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<UnlockUserRequest>): UnlockUserRequest {
    return UnlockUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UnlockUserRequest>): UnlockUserRequest {
    const message = createBaseUnlockUserRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseUnlockUserResponse(): UnlockUserResponse {
  return { details: undefined };
}

export const UnlockUserResponse = {
  encode(message: UnlockUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnlockUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnlockUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UnlockUserResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: UnlockUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<UnlockUserResponse>): UnlockUserResponse {
    return UnlockUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<UnlockUserResponse>): UnlockUserResponse {
    const message = createBaseUnlockUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseDeleteUserRequest(): DeleteUserRequest {
  return { userId: "" };
}

export const DeleteUserRequest = {
  encode(message: DeleteUserRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteUserRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): DeleteUserRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: DeleteUserRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<DeleteUserRequest>): DeleteUserRequest {
    return DeleteUserRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteUserRequest>): DeleteUserRequest {
    const message = createBaseDeleteUserRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseDeleteUserResponse(): DeleteUserResponse {
  return { details: undefined };
}

export const DeleteUserResponse = {
  encode(message: DeleteUserResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteUserResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteUserResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteUserResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: DeleteUserResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DeleteUserResponse>): DeleteUserResponse {
    return DeleteUserResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DeleteUserResponse>): DeleteUserResponse {
    const message = createBaseDeleteUserResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseSetContactEmailRequest(): SetContactEmailRequest {
  return { userId: "", email: undefined };
}

export const SetContactEmailRequest = {
  encode(message: SetContactEmailRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.email !== undefined) {
      SetEmail.encode(message.email, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetContactEmailRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetContactEmailRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.email = SetEmail.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetContactEmailRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      email: isSet(object.email) ? SetEmail.fromJSON(object.email) : undefined,
    };
  },

  toJSON(message: SetContactEmailRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.email !== undefined && (obj.email = message.email ? SetEmail.toJSON(message.email) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetContactEmailRequest>): SetContactEmailRequest {
    return SetContactEmailRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetContactEmailRequest>): SetContactEmailRequest {
    const message = createBaseSetContactEmailRequest();
    message.userId = object.userId ?? "";
    message.email = (object.email !== undefined && object.email !== null)
      ? SetEmail.fromPartial(object.email)
      : undefined;
    return message;
  },
};

function createBaseSetContactEmailResponse(): SetContactEmailResponse {
  return { details: undefined, verificationCode: undefined };
}

export const SetContactEmailResponse = {
  encode(message: SetContactEmailResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(26).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetContactEmailResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetContactEmailResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetContactEmailResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: SetContactEmailResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<SetContactEmailResponse>): SetContactEmailResponse {
    return SetContactEmailResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetContactEmailResponse>): SetContactEmailResponse {
    const message = createBaseSetContactEmailResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseVerifyContactEmailRequest(): VerifyContactEmailRequest {
  return { userId: "", verificationCode: "" };
}

export const VerifyContactEmailRequest = {
  encode(message: VerifyContactEmailRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.verificationCode !== "") {
      writer.uint32(18).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyContactEmailRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyContactEmailRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyContactEmailRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : "",
    };
  },

  toJSON(message: VerifyContactEmailRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<VerifyContactEmailRequest>): VerifyContactEmailRequest {
    return VerifyContactEmailRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyContactEmailRequest>): VerifyContactEmailRequest {
    const message = createBaseVerifyContactEmailRequest();
    message.userId = object.userId ?? "";
    message.verificationCode = object.verificationCode ?? "";
    return message;
  },
};

function createBaseVerifyContactEmailResponse(): VerifyContactEmailResponse {
  return { details: undefined };
}

export const VerifyContactEmailResponse = {
  encode(message: VerifyContactEmailResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyContactEmailResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyContactEmailResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyContactEmailResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyContactEmailResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyContactEmailResponse>): VerifyContactEmailResponse {
    return VerifyContactEmailResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyContactEmailResponse>): VerifyContactEmailResponse {
    const message = createBaseVerifyContactEmailResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResendContactEmailCodeRequest(): ResendContactEmailCodeRequest {
  return { userId: "", sendCode: undefined, returnCode: undefined };
}

export const ResendContactEmailCodeRequest = {
  encode(message: ResendContactEmailCodeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.sendCode !== undefined) {
      SendEmailVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnEmailVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendContactEmailCodeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendContactEmailCodeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendCode = SendEmailVerificationCode.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = ReturnEmailVerificationCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResendContactEmailCodeRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      sendCode: isSet(object.sendCode) ? SendEmailVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnEmailVerificationCode.fromJSON(object.returnCode) : undefined,
    };
  },

  toJSON(message: ResendContactEmailCodeRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendEmailVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnEmailVerificationCode.toJSON(message.returnCode) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResendContactEmailCodeRequest>): ResendContactEmailCodeRequest {
    return ResendContactEmailCodeRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendContactEmailCodeRequest>): ResendContactEmailCodeRequest {
    const message = createBaseResendContactEmailCodeRequest();
    message.userId = object.userId ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendEmailVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnEmailVerificationCode.fromPartial(object.returnCode)
      : undefined;
    return message;
  },
};

function createBaseResendContactEmailCodeResponse(): ResendContactEmailCodeResponse {
  return { details: undefined, verificationCode: undefined };
}

export const ResendContactEmailCodeResponse = {
  encode(message: ResendContactEmailCodeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(18).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendContactEmailCodeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendContactEmailCodeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResendContactEmailCodeResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: ResendContactEmailCodeResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<ResendContactEmailCodeResponse>): ResendContactEmailCodeResponse {
    return ResendContactEmailCodeResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendContactEmailCodeResponse>): ResendContactEmailCodeResponse {
    const message = createBaseResendContactEmailCodeResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseSetContactPhoneRequest(): SetContactPhoneRequest {
  return { userId: "", phone: undefined };
}

export const SetContactPhoneRequest = {
  encode(message: SetContactPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.phone !== undefined) {
      SetPhone.encode(message.phone, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetContactPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetContactPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.phone = SetPhone.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetContactPhoneRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      phone: isSet(object.phone) ? SetPhone.fromJSON(object.phone) : undefined,
    };
  },

  toJSON(message: SetContactPhoneRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.phone !== undefined && (obj.phone = message.phone ? SetPhone.toJSON(message.phone) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetContactPhoneRequest>): SetContactPhoneRequest {
    return SetContactPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetContactPhoneRequest>): SetContactPhoneRequest {
    const message = createBaseSetContactPhoneRequest();
    message.userId = object.userId ?? "";
    message.phone = (object.phone !== undefined && object.phone !== null)
      ? SetPhone.fromPartial(object.phone)
      : undefined;
    return message;
  },
};

function createBaseSetContactPhoneResponse(): SetContactPhoneResponse {
  return { details: undefined, emailCode: undefined };
}

export const SetContactPhoneResponse = {
  encode(message: SetContactPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.emailCode !== undefined) {
      writer.uint32(26).string(message.emailCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetContactPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetContactPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.emailCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetContactPhoneResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      emailCode: isSet(object.emailCode) ? String(object.emailCode) : undefined,
    };
  },

  toJSON(message: SetContactPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.emailCode !== undefined && (obj.emailCode = message.emailCode);
    return obj;
  },

  create(base?: DeepPartial<SetContactPhoneResponse>): SetContactPhoneResponse {
    return SetContactPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetContactPhoneResponse>): SetContactPhoneResponse {
    const message = createBaseSetContactPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.emailCode = object.emailCode ?? undefined;
    return message;
  },
};

function createBaseVerifyContactPhoneRequest(): VerifyContactPhoneRequest {
  return { userId: "", verificationCode: "" };
}

export const VerifyContactPhoneRequest = {
  encode(message: VerifyContactPhoneRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.verificationCode !== "") {
      writer.uint32(18).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyContactPhoneRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyContactPhoneRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyContactPhoneRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : "",
    };
  },

  toJSON(message: VerifyContactPhoneRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<VerifyContactPhoneRequest>): VerifyContactPhoneRequest {
    return VerifyContactPhoneRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyContactPhoneRequest>): VerifyContactPhoneRequest {
    const message = createBaseVerifyContactPhoneRequest();
    message.userId = object.userId ?? "";
    message.verificationCode = object.verificationCode ?? "";
    return message;
  },
};

function createBaseVerifyContactPhoneResponse(): VerifyContactPhoneResponse {
  return { details: undefined };
}

export const VerifyContactPhoneResponse = {
  encode(message: VerifyContactPhoneResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyContactPhoneResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyContactPhoneResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyContactPhoneResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyContactPhoneResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyContactPhoneResponse>): VerifyContactPhoneResponse {
    return VerifyContactPhoneResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyContactPhoneResponse>): VerifyContactPhoneResponse {
    const message = createBaseVerifyContactPhoneResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseResendContactPhoneCodeRequest(): ResendContactPhoneCodeRequest {
  return { userId: "", sendCode: undefined, returnCode: undefined };
}

export const ResendContactPhoneCodeRequest = {
  encode(message: ResendContactPhoneCodeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.sendCode !== undefined) {
      SendPhoneVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnPhoneVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendContactPhoneCodeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendContactPhoneCodeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendCode = SendPhoneVerificationCode.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = ReturnPhoneVerificationCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResendContactPhoneCodeRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      sendCode: isSet(object.sendCode) ? SendPhoneVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnPhoneVerificationCode.fromJSON(object.returnCode) : undefined,
    };
  },

  toJSON(message: ResendContactPhoneCodeRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendPhoneVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnPhoneVerificationCode.toJSON(message.returnCode) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ResendContactPhoneCodeRequest>): ResendContactPhoneCodeRequest {
    return ResendContactPhoneCodeRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendContactPhoneCodeRequest>): ResendContactPhoneCodeRequest {
    const message = createBaseResendContactPhoneCodeRequest();
    message.userId = object.userId ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendPhoneVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnPhoneVerificationCode.fromPartial(object.returnCode)
      : undefined;
    return message;
  },
};

function createBaseResendContactPhoneCodeResponse(): ResendContactPhoneCodeResponse {
  return { details: undefined, verificationCode: undefined };
}

export const ResendContactPhoneCodeResponse = {
  encode(message: ResendContactPhoneCodeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(18).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResendContactPhoneCodeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResendContactPhoneCodeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResendContactPhoneCodeResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: ResendContactPhoneCodeResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<ResendContactPhoneCodeResponse>): ResendContactPhoneCodeResponse {
    return ResendContactPhoneCodeResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ResendContactPhoneCodeResponse>): ResendContactPhoneCodeResponse {
    const message = createBaseResendContactPhoneCodeResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseAddUsernameRequest(): AddUsernameRequest {
  return { userId: "", username: undefined };
}

export const AddUsernameRequest = {
  encode(message: AddUsernameRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.username !== undefined) {
      SetUsername.encode(message.username, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddUsernameRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddUsernameRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.username = SetUsername.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddUsernameRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      username: isSet(object.username) ? SetUsername.fromJSON(object.username) : undefined,
    };
  },

  toJSON(message: AddUsernameRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.username !== undefined &&
      (obj.username = message.username ? SetUsername.toJSON(message.username) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddUsernameRequest>): AddUsernameRequest {
    return AddUsernameRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddUsernameRequest>): AddUsernameRequest {
    const message = createBaseAddUsernameRequest();
    message.userId = object.userId ?? "";
    message.username = (object.username !== undefined && object.username !== null)
      ? SetUsername.fromPartial(object.username)
      : undefined;
    return message;
  },
};

function createBaseAddUsernameResponse(): AddUsernameResponse {
  return { details: undefined, usernameId: "" };
}

export const AddUsernameResponse = {
  encode(message: AddUsernameResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.usernameId !== "") {
      writer.uint32(18).string(message.usernameId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddUsernameResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddUsernameResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.usernameId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddUsernameResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      usernameId: isSet(object.usernameId) ? String(object.usernameId) : "",
    };
  },

  toJSON(message: AddUsernameResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.usernameId !== undefined && (obj.usernameId = message.usernameId);
    return obj;
  },

  create(base?: DeepPartial<AddUsernameResponse>): AddUsernameResponse {
    return AddUsernameResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddUsernameResponse>): AddUsernameResponse {
    const message = createBaseAddUsernameResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.usernameId = object.usernameId ?? "";
    return message;
  },
};

function createBaseRemoveUsernameRequest(): RemoveUsernameRequest {
  return { userId: "", usernameId: "" };
}

export const RemoveUsernameRequest = {
  encode(message: RemoveUsernameRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.usernameId !== "") {
      writer.uint32(18).string(message.usernameId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveUsernameRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveUsernameRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.usernameId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveUsernameRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      usernameId: isSet(object.usernameId) ? String(object.usernameId) : "",
    };
  },

  toJSON(message: RemoveUsernameRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.usernameId !== undefined && (obj.usernameId = message.usernameId);
    return obj;
  },

  create(base?: DeepPartial<RemoveUsernameRequest>): RemoveUsernameRequest {
    return RemoveUsernameRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveUsernameRequest>): RemoveUsernameRequest {
    const message = createBaseRemoveUsernameRequest();
    message.userId = object.userId ?? "";
    message.usernameId = object.usernameId ?? "";
    return message;
  },
};

function createBaseRemoveUsernameResponse(): RemoveUsernameResponse {
  return { details: undefined };
}

export const RemoveUsernameResponse = {
  encode(message: RemoveUsernameResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveUsernameResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveUsernameResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveUsernameResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveUsernameResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveUsernameResponse>): RemoveUsernameResponse {
    return RemoveUsernameResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveUsernameResponse>): RemoveUsernameResponse {
    const message = createBaseRemoveUsernameResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseSetPasswordRequest(): SetPasswordRequest {
  return { userId: "", newPassword: undefined, currentPassword: undefined, verificationCode: undefined };
}

export const SetPasswordRequest = {
  encode(message: SetPasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.newPassword !== undefined) {
      SetPassword.encode(message.newPassword, writer.uint32(18).fork()).ldelim();
    }
    if (message.currentPassword !== undefined) {
      writer.uint32(26).string(message.currentPassword);
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(34).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetPasswordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.newPassword = SetPassword.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.currentPassword = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetPasswordRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      newPassword: isSet(object.newPassword) ? SetPassword.fromJSON(object.newPassword) : undefined,
      currentPassword: isSet(object.currentPassword) ? String(object.currentPassword) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: SetPasswordRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.newPassword !== undefined &&
      (obj.newPassword = message.newPassword ? SetPassword.toJSON(message.newPassword) : undefined);
    message.currentPassword !== undefined && (obj.currentPassword = message.currentPassword);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<SetPasswordRequest>): SetPasswordRequest {
    return SetPasswordRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetPasswordRequest>): SetPasswordRequest {
    const message = createBaseSetPasswordRequest();
    message.userId = object.userId ?? "";
    message.newPassword = (object.newPassword !== undefined && object.newPassword !== null)
      ? SetPassword.fromPartial(object.newPassword)
      : undefined;
    message.currentPassword = object.currentPassword ?? undefined;
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseSetPasswordResponse(): SetPasswordResponse {
  return { details: undefined };
}

export const SetPasswordResponse = {
  encode(message: SetPasswordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetPasswordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetPasswordResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: SetPasswordResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetPasswordResponse>): SetPasswordResponse {
    return SetPasswordResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetPasswordResponse>): SetPasswordResponse {
    const message = createBaseSetPasswordResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRequestPasswordResetRequest(): RequestPasswordResetRequest {
  return { userId: "", sendEmail: undefined, sendSms: undefined, returnCode: undefined };
}

export const RequestPasswordResetRequest = {
  encode(message: RequestPasswordResetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.sendEmail !== undefined) {
      SendPasswordResetEmail.encode(message.sendEmail, writer.uint32(18).fork()).ldelim();
    }
    if (message.sendSms !== undefined) {
      SendPasswordResetSMS.encode(message.sendSms, writer.uint32(26).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnPasswordResetCode.encode(message.returnCode, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestPasswordResetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestPasswordResetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendEmail = SendPasswordResetEmail.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.sendSms = SendPasswordResetSMS.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.returnCode = ReturnPasswordResetCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestPasswordResetRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      sendEmail: isSet(object.sendEmail) ? SendPasswordResetEmail.fromJSON(object.sendEmail) : undefined,
      sendSms: isSet(object.sendSms) ? SendPasswordResetSMS.fromJSON(object.sendSms) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnPasswordResetCode.fromJSON(object.returnCode) : undefined,
    };
  },

  toJSON(message: RequestPasswordResetRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.sendEmail !== undefined &&
      (obj.sendEmail = message.sendEmail ? SendPasswordResetEmail.toJSON(message.sendEmail) : undefined);
    message.sendSms !== undefined &&
      (obj.sendSms = message.sendSms ? SendPasswordResetSMS.toJSON(message.sendSms) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnPasswordResetCode.toJSON(message.returnCode) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RequestPasswordResetRequest>): RequestPasswordResetRequest {
    return RequestPasswordResetRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestPasswordResetRequest>): RequestPasswordResetRequest {
    const message = createBaseRequestPasswordResetRequest();
    message.userId = object.userId ?? "";
    message.sendEmail = (object.sendEmail !== undefined && object.sendEmail !== null)
      ? SendPasswordResetEmail.fromPartial(object.sendEmail)
      : undefined;
    message.sendSms = (object.sendSms !== undefined && object.sendSms !== null)
      ? SendPasswordResetSMS.fromPartial(object.sendSms)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnPasswordResetCode.fromPartial(object.returnCode)
      : undefined;
    return message;
  },
};

function createBaseRequestPasswordResetResponse(): RequestPasswordResetResponse {
  return { details: undefined, verificationCode: undefined };
}

export const RequestPasswordResetResponse = {
  encode(message: RequestPasswordResetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(18).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestPasswordResetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestPasswordResetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestPasswordResetResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: RequestPasswordResetResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<RequestPasswordResetResponse>): RequestPasswordResetResponse {
    return RequestPasswordResetResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestPasswordResetResponse>): RequestPasswordResetResponse {
    const message = createBaseRequestPasswordResetResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseStartWebAuthNRegistrationRequest(): StartWebAuthNRegistrationRequest {
  return { userId: "", domain: "", authenticatorType: 0, code: undefined };
}

export const StartWebAuthNRegistrationRequest = {
  encode(message: StartWebAuthNRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.domain !== "") {
      writer.uint32(34).string(message.domain);
    }
    if (message.authenticatorType !== 0) {
      writer.uint32(24).int32(message.authenticatorType);
    }
    if (message.code !== undefined) {
      AuthenticatorRegistrationCode.encode(message.code, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartWebAuthNRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartWebAuthNRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.domain = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.authenticatorType = reader.int32() as any;
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.code = AuthenticatorRegistrationCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartWebAuthNRegistrationRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      domain: isSet(object.domain) ? String(object.domain) : "",
      authenticatorType: isSet(object.authenticatorType)
        ? webAuthNAuthenticatorTypeFromJSON(object.authenticatorType)
        : 0,
      code: isSet(object.code) ? AuthenticatorRegistrationCode.fromJSON(object.code) : undefined,
    };
  },

  toJSON(message: StartWebAuthNRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.domain !== undefined && (obj.domain = message.domain);
    message.authenticatorType !== undefined &&
      (obj.authenticatorType = webAuthNAuthenticatorTypeToJSON(message.authenticatorType));
    message.code !== undefined &&
      (obj.code = message.code ? AuthenticatorRegistrationCode.toJSON(message.code) : undefined);
    return obj;
  },

  create(base?: DeepPartial<StartWebAuthNRegistrationRequest>): StartWebAuthNRegistrationRequest {
    return StartWebAuthNRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartWebAuthNRegistrationRequest>): StartWebAuthNRegistrationRequest {
    const message = createBaseStartWebAuthNRegistrationRequest();
    message.userId = object.userId ?? "";
    message.domain = object.domain ?? "";
    message.authenticatorType = object.authenticatorType ?? 0;
    message.code = (object.code !== undefined && object.code !== null)
      ? AuthenticatorRegistrationCode.fromPartial(object.code)
      : undefined;
    return message;
  },
};

function createBaseStartWebAuthNRegistrationResponse(): StartWebAuthNRegistrationResponse {
  return { details: undefined, webAuthNId: "", publicKeyCredentialCreationOptions: undefined };
}

export const StartWebAuthNRegistrationResponse = {
  encode(message: StartWebAuthNRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.webAuthNId !== "") {
      writer.uint32(18).string(message.webAuthNId);
    }
    if (message.publicKeyCredentialCreationOptions !== undefined) {
      Struct.encode(Struct.wrap(message.publicKeyCredentialCreationOptions), writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartWebAuthNRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartWebAuthNRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.webAuthNId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.publicKeyCredentialCreationOptions = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartWebAuthNRegistrationResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      webAuthNId: isSet(object.webAuthNId) ? String(object.webAuthNId) : "",
      publicKeyCredentialCreationOptions: isObject(object.publicKeyCredentialCreationOptions)
        ? object.publicKeyCredentialCreationOptions
        : undefined,
    };
  },

  toJSON(message: StartWebAuthNRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.webAuthNId !== undefined && (obj.webAuthNId = message.webAuthNId);
    message.publicKeyCredentialCreationOptions !== undefined &&
      (obj.publicKeyCredentialCreationOptions = message.publicKeyCredentialCreationOptions);
    return obj;
  },

  create(base?: DeepPartial<StartWebAuthNRegistrationResponse>): StartWebAuthNRegistrationResponse {
    return StartWebAuthNRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartWebAuthNRegistrationResponse>): StartWebAuthNRegistrationResponse {
    const message = createBaseStartWebAuthNRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.webAuthNId = object.webAuthNId ?? "";
    message.publicKeyCredentialCreationOptions = object.publicKeyCredentialCreationOptions ?? undefined;
    return message;
  },
};

function createBaseVerifyWebAuthNRegistrationRequest(): VerifyWebAuthNRegistrationRequest {
  return { userId: "", webAuthNId: "", publicKeyCredential: undefined, webAuthNName: "" };
}

export const VerifyWebAuthNRegistrationRequest = {
  encode(message: VerifyWebAuthNRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.webAuthNId !== "") {
      writer.uint32(18).string(message.webAuthNId);
    }
    if (message.publicKeyCredential !== undefined) {
      Struct.encode(Struct.wrap(message.publicKeyCredential), writer.uint32(26).fork()).ldelim();
    }
    if (message.webAuthNName !== "") {
      writer.uint32(34).string(message.webAuthNName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyWebAuthNRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyWebAuthNRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.webAuthNId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.publicKeyCredential = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.webAuthNName = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyWebAuthNRegistrationRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      webAuthNId: isSet(object.webAuthNId) ? String(object.webAuthNId) : "",
      publicKeyCredential: isObject(object.publicKeyCredential) ? object.publicKeyCredential : undefined,
      webAuthNName: isSet(object.webAuthNName) ? String(object.webAuthNName) : "",
    };
  },

  toJSON(message: VerifyWebAuthNRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.webAuthNId !== undefined && (obj.webAuthNId = message.webAuthNId);
    message.publicKeyCredential !== undefined && (obj.publicKeyCredential = message.publicKeyCredential);
    message.webAuthNName !== undefined && (obj.webAuthNName = message.webAuthNName);
    return obj;
  },

  create(base?: DeepPartial<VerifyWebAuthNRegistrationRequest>): VerifyWebAuthNRegistrationRequest {
    return VerifyWebAuthNRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyWebAuthNRegistrationRequest>): VerifyWebAuthNRegistrationRequest {
    const message = createBaseVerifyWebAuthNRegistrationRequest();
    message.userId = object.userId ?? "";
    message.webAuthNId = object.webAuthNId ?? "";
    message.publicKeyCredential = object.publicKeyCredential ?? undefined;
    message.webAuthNName = object.webAuthNName ?? "";
    return message;
  },
};

function createBaseVerifyWebAuthNRegistrationResponse(): VerifyWebAuthNRegistrationResponse {
  return { details: undefined };
}

export const VerifyWebAuthNRegistrationResponse = {
  encode(message: VerifyWebAuthNRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyWebAuthNRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyWebAuthNRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyWebAuthNRegistrationResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyWebAuthNRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyWebAuthNRegistrationResponse>): VerifyWebAuthNRegistrationResponse {
    return VerifyWebAuthNRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyWebAuthNRegistrationResponse>): VerifyWebAuthNRegistrationResponse {
    const message = createBaseVerifyWebAuthNRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseCreateWebAuthNRegistrationLinkRequest(): CreateWebAuthNRegistrationLinkRequest {
  return { userId: "", sendLink: undefined, returnCode: undefined };
}

export const CreateWebAuthNRegistrationLinkRequest = {
  encode(message: CreateWebAuthNRegistrationLinkRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.sendLink !== undefined) {
      SendWebAuthNRegistrationLink.encode(message.sendLink, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnWebAuthNRegistrationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateWebAuthNRegistrationLinkRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateWebAuthNRegistrationLinkRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendLink = SendWebAuthNRegistrationLink.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = ReturnWebAuthNRegistrationCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateWebAuthNRegistrationLinkRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      sendLink: isSet(object.sendLink) ? SendWebAuthNRegistrationLink.fromJSON(object.sendLink) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnWebAuthNRegistrationCode.fromJSON(object.returnCode) : undefined,
    };
  },

  toJSON(message: CreateWebAuthNRegistrationLinkRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.sendLink !== undefined &&
      (obj.sendLink = message.sendLink ? SendWebAuthNRegistrationLink.toJSON(message.sendLink) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnWebAuthNRegistrationCode.toJSON(message.returnCode) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateWebAuthNRegistrationLinkRequest>): CreateWebAuthNRegistrationLinkRequest {
    return CreateWebAuthNRegistrationLinkRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateWebAuthNRegistrationLinkRequest>): CreateWebAuthNRegistrationLinkRequest {
    const message = createBaseCreateWebAuthNRegistrationLinkRequest();
    message.userId = object.userId ?? "";
    message.sendLink = (object.sendLink !== undefined && object.sendLink !== null)
      ? SendWebAuthNRegistrationLink.fromPartial(object.sendLink)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnWebAuthNRegistrationCode.fromPartial(object.returnCode)
      : undefined;
    return message;
  },
};

function createBaseCreateWebAuthNRegistrationLinkResponse(): CreateWebAuthNRegistrationLinkResponse {
  return { details: undefined, code: undefined };
}

export const CreateWebAuthNRegistrationLinkResponse = {
  encode(message: CreateWebAuthNRegistrationLinkResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.code !== undefined) {
      AuthenticatorRegistrationCode.encode(message.code, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateWebAuthNRegistrationLinkResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateWebAuthNRegistrationLinkResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.code = AuthenticatorRegistrationCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateWebAuthNRegistrationLinkResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      code: isSet(object.code) ? AuthenticatorRegistrationCode.fromJSON(object.code) : undefined,
    };
  },

  toJSON(message: CreateWebAuthNRegistrationLinkResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.code !== undefined &&
      (obj.code = message.code ? AuthenticatorRegistrationCode.toJSON(message.code) : undefined);
    return obj;
  },

  create(base?: DeepPartial<CreateWebAuthNRegistrationLinkResponse>): CreateWebAuthNRegistrationLinkResponse {
    return CreateWebAuthNRegistrationLinkResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<CreateWebAuthNRegistrationLinkResponse>): CreateWebAuthNRegistrationLinkResponse {
    const message = createBaseCreateWebAuthNRegistrationLinkResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.code = (object.code !== undefined && object.code !== null)
      ? AuthenticatorRegistrationCode.fromPartial(object.code)
      : undefined;
    return message;
  },
};

function createBaseRemoveWebAuthNAuthenticatorRequest(): RemoveWebAuthNAuthenticatorRequest {
  return { userId: "", webAuthNId: "" };
}

export const RemoveWebAuthNAuthenticatorRequest = {
  encode(message: RemoveWebAuthNAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.webAuthNId !== "") {
      writer.uint32(18).string(message.webAuthNId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveWebAuthNAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveWebAuthNAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.webAuthNId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveWebAuthNAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      webAuthNId: isSet(object.webAuthNId) ? String(object.webAuthNId) : "",
    };
  },

  toJSON(message: RemoveWebAuthNAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.webAuthNId !== undefined && (obj.webAuthNId = message.webAuthNId);
    return obj;
  },

  create(base?: DeepPartial<RemoveWebAuthNAuthenticatorRequest>): RemoveWebAuthNAuthenticatorRequest {
    return RemoveWebAuthNAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveWebAuthNAuthenticatorRequest>): RemoveWebAuthNAuthenticatorRequest {
    const message = createBaseRemoveWebAuthNAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.webAuthNId = object.webAuthNId ?? "";
    return message;
  },
};

function createBaseRemoveWebAuthNAuthenticatorResponse(): RemoveWebAuthNAuthenticatorResponse {
  return { details: undefined };
}

export const RemoveWebAuthNAuthenticatorResponse = {
  encode(message: RemoveWebAuthNAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveWebAuthNAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveWebAuthNAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveWebAuthNAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveWebAuthNAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveWebAuthNAuthenticatorResponse>): RemoveWebAuthNAuthenticatorResponse {
    return RemoveWebAuthNAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveWebAuthNAuthenticatorResponse>): RemoveWebAuthNAuthenticatorResponse {
    const message = createBaseRemoveWebAuthNAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseStartTOTPRegistrationRequest(): StartTOTPRegistrationRequest {
  return { userId: "" };
}

export const StartTOTPRegistrationRequest = {
  encode(message: StartTOTPRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartTOTPRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartTOTPRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
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

  fromJSON(object: any): StartTOTPRegistrationRequest {
    return { userId: isSet(object.userId) ? String(object.userId) : "" };
  },

  toJSON(message: StartTOTPRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<StartTOTPRegistrationRequest>): StartTOTPRegistrationRequest {
    return StartTOTPRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartTOTPRegistrationRequest>): StartTOTPRegistrationRequest {
    const message = createBaseStartTOTPRegistrationRequest();
    message.userId = object.userId ?? "";
    return message;
  },
};

function createBaseStartTOTPRegistrationResponse(): StartTOTPRegistrationResponse {
  return { details: undefined, totpId: "", uri: "", secret: "" };
}

export const StartTOTPRegistrationResponse = {
  encode(message: StartTOTPRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.totpId !== "") {
      writer.uint32(18).string(message.totpId);
    }
    if (message.uri !== "") {
      writer.uint32(26).string(message.uri);
    }
    if (message.secret !== "") {
      writer.uint32(34).string(message.secret);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartTOTPRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartTOTPRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.totpId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.uri = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.secret = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartTOTPRegistrationResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      totpId: isSet(object.totpId) ? String(object.totpId) : "",
      uri: isSet(object.uri) ? String(object.uri) : "",
      secret: isSet(object.secret) ? String(object.secret) : "",
    };
  },

  toJSON(message: StartTOTPRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.totpId !== undefined && (obj.totpId = message.totpId);
    message.uri !== undefined && (obj.uri = message.uri);
    message.secret !== undefined && (obj.secret = message.secret);
    return obj;
  },

  create(base?: DeepPartial<StartTOTPRegistrationResponse>): StartTOTPRegistrationResponse {
    return StartTOTPRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartTOTPRegistrationResponse>): StartTOTPRegistrationResponse {
    const message = createBaseStartTOTPRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.totpId = object.totpId ?? "";
    message.uri = object.uri ?? "";
    message.secret = object.secret ?? "";
    return message;
  },
};

function createBaseVerifyTOTPRegistrationRequest(): VerifyTOTPRegistrationRequest {
  return { userId: "", totpId: "", code: "" };
}

export const VerifyTOTPRegistrationRequest = {
  encode(message: VerifyTOTPRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.totpId !== "") {
      writer.uint32(18).string(message.totpId);
    }
    if (message.code !== "") {
      writer.uint32(26).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyTOTPRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyTOTPRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.totpId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
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

  fromJSON(object: any): VerifyTOTPRegistrationRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      totpId: isSet(object.totpId) ? String(object.totpId) : "",
      code: isSet(object.code) ? String(object.code) : "",
    };
  },

  toJSON(message: VerifyTOTPRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.totpId !== undefined && (obj.totpId = message.totpId);
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyTOTPRegistrationRequest>): VerifyTOTPRegistrationRequest {
    return VerifyTOTPRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyTOTPRegistrationRequest>): VerifyTOTPRegistrationRequest {
    const message = createBaseVerifyTOTPRegistrationRequest();
    message.userId = object.userId ?? "";
    message.totpId = object.totpId ?? "";
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyTOTPRegistrationResponse(): VerifyTOTPRegistrationResponse {
  return { details: undefined };
}

export const VerifyTOTPRegistrationResponse = {
  encode(message: VerifyTOTPRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyTOTPRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyTOTPRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyTOTPRegistrationResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyTOTPRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyTOTPRegistrationResponse>): VerifyTOTPRegistrationResponse {
    return VerifyTOTPRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyTOTPRegistrationResponse>): VerifyTOTPRegistrationResponse {
    const message = createBaseVerifyTOTPRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveTOTPAuthenticatorRequest(): RemoveTOTPAuthenticatorRequest {
  return { userId: "", totpId: "" };
}

export const RemoveTOTPAuthenticatorRequest = {
  encode(message: RemoveTOTPAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.totpId !== "") {
      writer.uint32(18).string(message.totpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveTOTPAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveTOTPAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.totpId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveTOTPAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      totpId: isSet(object.totpId) ? String(object.totpId) : "",
    };
  },

  toJSON(message: RemoveTOTPAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.totpId !== undefined && (obj.totpId = message.totpId);
    return obj;
  },

  create(base?: DeepPartial<RemoveTOTPAuthenticatorRequest>): RemoveTOTPAuthenticatorRequest {
    return RemoveTOTPAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveTOTPAuthenticatorRequest>): RemoveTOTPAuthenticatorRequest {
    const message = createBaseRemoveTOTPAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.totpId = object.totpId ?? "";
    return message;
  },
};

function createBaseRemoveTOTPAuthenticatorResponse(): RemoveTOTPAuthenticatorResponse {
  return { details: undefined };
}

export const RemoveTOTPAuthenticatorResponse = {
  encode(message: RemoveTOTPAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveTOTPAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveTOTPAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveTOTPAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveTOTPAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveTOTPAuthenticatorResponse>): RemoveTOTPAuthenticatorResponse {
    return RemoveTOTPAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveTOTPAuthenticatorResponse>): RemoveTOTPAuthenticatorResponse {
    const message = createBaseRemoveTOTPAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddOTPSMSAuthenticatorRequest(): AddOTPSMSAuthenticatorRequest {
  return { userId: "", phone: undefined };
}

export const AddOTPSMSAuthenticatorRequest = {
  encode(message: AddOTPSMSAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.phone !== undefined) {
      SetPhone.encode(message.phone, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOTPSMSAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOTPSMSAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.phone = SetPhone.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOTPSMSAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      phone: isSet(object.phone) ? SetPhone.fromJSON(object.phone) : undefined,
    };
  },

  toJSON(message: AddOTPSMSAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.phone !== undefined && (obj.phone = message.phone ? SetPhone.toJSON(message.phone) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddOTPSMSAuthenticatorRequest>): AddOTPSMSAuthenticatorRequest {
    return AddOTPSMSAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOTPSMSAuthenticatorRequest>): AddOTPSMSAuthenticatorRequest {
    const message = createBaseAddOTPSMSAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.phone = (object.phone !== undefined && object.phone !== null)
      ? SetPhone.fromPartial(object.phone)
      : undefined;
    return message;
  },
};

function createBaseAddOTPSMSAuthenticatorResponse(): AddOTPSMSAuthenticatorResponse {
  return { details: undefined, otpSmsId: "", verificationCode: undefined };
}

export const AddOTPSMSAuthenticatorResponse = {
  encode(message: AddOTPSMSAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.otpSmsId !== "") {
      writer.uint32(18).string(message.otpSmsId);
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(26).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOTPSMSAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOTPSMSAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpSmsId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOTPSMSAuthenticatorResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      otpSmsId: isSet(object.otpSmsId) ? String(object.otpSmsId) : "",
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: AddOTPSMSAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.otpSmsId !== undefined && (obj.otpSmsId = message.otpSmsId);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<AddOTPSMSAuthenticatorResponse>): AddOTPSMSAuthenticatorResponse {
    return AddOTPSMSAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOTPSMSAuthenticatorResponse>): AddOTPSMSAuthenticatorResponse {
    const message = createBaseAddOTPSMSAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.otpSmsId = object.otpSmsId ?? "";
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseVerifyOTPSMSRegistrationRequest(): VerifyOTPSMSRegistrationRequest {
  return { userId: "", otpSmsId: "", code: "" };
}

export const VerifyOTPSMSRegistrationRequest = {
  encode(message: VerifyOTPSMSRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.otpSmsId !== "") {
      writer.uint32(18).string(message.otpSmsId);
    }
    if (message.code !== "") {
      writer.uint32(26).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyOTPSMSRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyOTPSMSRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpSmsId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
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

  fromJSON(object: any): VerifyOTPSMSRegistrationRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      otpSmsId: isSet(object.otpSmsId) ? String(object.otpSmsId) : "",
      code: isSet(object.code) ? String(object.code) : "",
    };
  },

  toJSON(message: VerifyOTPSMSRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.otpSmsId !== undefined && (obj.otpSmsId = message.otpSmsId);
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyOTPSMSRegistrationRequest>): VerifyOTPSMSRegistrationRequest {
    return VerifyOTPSMSRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyOTPSMSRegistrationRequest>): VerifyOTPSMSRegistrationRequest {
    const message = createBaseVerifyOTPSMSRegistrationRequest();
    message.userId = object.userId ?? "";
    message.otpSmsId = object.otpSmsId ?? "";
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyOTPSMSRegistrationResponse(): VerifyOTPSMSRegistrationResponse {
  return { details: undefined };
}

export const VerifyOTPSMSRegistrationResponse = {
  encode(message: VerifyOTPSMSRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyOTPSMSRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyOTPSMSRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyOTPSMSRegistrationResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyOTPSMSRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyOTPSMSRegistrationResponse>): VerifyOTPSMSRegistrationResponse {
    return VerifyOTPSMSRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyOTPSMSRegistrationResponse>): VerifyOTPSMSRegistrationResponse {
    const message = createBaseVerifyOTPSMSRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveOTPSMSAuthenticatorRequest(): RemoveOTPSMSAuthenticatorRequest {
  return { userId: "", otpSmsId: "" };
}

export const RemoveOTPSMSAuthenticatorRequest = {
  encode(message: RemoveOTPSMSAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.otpSmsId !== "") {
      writer.uint32(18).string(message.otpSmsId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOTPSMSAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOTPSMSAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpSmsId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveOTPSMSAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      otpSmsId: isSet(object.otpSmsId) ? String(object.otpSmsId) : "",
    };
  },

  toJSON(message: RemoveOTPSMSAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.otpSmsId !== undefined && (obj.otpSmsId = message.otpSmsId);
    return obj;
  },

  create(base?: DeepPartial<RemoveOTPSMSAuthenticatorRequest>): RemoveOTPSMSAuthenticatorRequest {
    return RemoveOTPSMSAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOTPSMSAuthenticatorRequest>): RemoveOTPSMSAuthenticatorRequest {
    const message = createBaseRemoveOTPSMSAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.otpSmsId = object.otpSmsId ?? "";
    return message;
  },
};

function createBaseRemoveOTPSMSAuthenticatorResponse(): RemoveOTPSMSAuthenticatorResponse {
  return { details: undefined };
}

export const RemoveOTPSMSAuthenticatorResponse = {
  encode(message: RemoveOTPSMSAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOTPSMSAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOTPSMSAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveOTPSMSAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveOTPSMSAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveOTPSMSAuthenticatorResponse>): RemoveOTPSMSAuthenticatorResponse {
    return RemoveOTPSMSAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOTPSMSAuthenticatorResponse>): RemoveOTPSMSAuthenticatorResponse {
    const message = createBaseRemoveOTPSMSAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseAddOTPEmailAuthenticatorRequest(): AddOTPEmailAuthenticatorRequest {
  return { userId: "", email: undefined };
}

export const AddOTPEmailAuthenticatorRequest = {
  encode(message: AddOTPEmailAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.email !== undefined) {
      SetEmail.encode(message.email, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOTPEmailAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOTPEmailAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.email = SetEmail.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOTPEmailAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      email: isSet(object.email) ? SetEmail.fromJSON(object.email) : undefined,
    };
  },

  toJSON(message: AddOTPEmailAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.email !== undefined && (obj.email = message.email ? SetEmail.toJSON(message.email) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddOTPEmailAuthenticatorRequest>): AddOTPEmailAuthenticatorRequest {
    return AddOTPEmailAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOTPEmailAuthenticatorRequest>): AddOTPEmailAuthenticatorRequest {
    const message = createBaseAddOTPEmailAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.email = (object.email !== undefined && object.email !== null)
      ? SetEmail.fromPartial(object.email)
      : undefined;
    return message;
  },
};

function createBaseAddOTPEmailAuthenticatorResponse(): AddOTPEmailAuthenticatorResponse {
  return { details: undefined, otpEmailId: "", verificationCode: undefined };
}

export const AddOTPEmailAuthenticatorResponse = {
  encode(message: AddOTPEmailAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.otpEmailId !== "") {
      writer.uint32(18).string(message.otpEmailId);
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(26).string(message.verificationCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddOTPEmailAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddOTPEmailAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpEmailId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.verificationCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddOTPEmailAuthenticatorResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      otpEmailId: isSet(object.otpEmailId) ? String(object.otpEmailId) : "",
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: AddOTPEmailAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.otpEmailId !== undefined && (obj.otpEmailId = message.otpEmailId);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<AddOTPEmailAuthenticatorResponse>): AddOTPEmailAuthenticatorResponse {
    return AddOTPEmailAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddOTPEmailAuthenticatorResponse>): AddOTPEmailAuthenticatorResponse {
    const message = createBaseAddOTPEmailAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.otpEmailId = object.otpEmailId ?? "";
    message.verificationCode = object.verificationCode ?? undefined;
    return message;
  },
};

function createBaseVerifyOTPEmailRegistrationRequest(): VerifyOTPEmailRegistrationRequest {
  return { userId: "", otpEmailId: "", code: "" };
}

export const VerifyOTPEmailRegistrationRequest = {
  encode(message: VerifyOTPEmailRegistrationRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.otpEmailId !== "") {
      writer.uint32(18).string(message.otpEmailId);
    }
    if (message.code !== "") {
      writer.uint32(26).string(message.code);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyOTPEmailRegistrationRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyOTPEmailRegistrationRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpEmailId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
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

  fromJSON(object: any): VerifyOTPEmailRegistrationRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      otpEmailId: isSet(object.otpEmailId) ? String(object.otpEmailId) : "",
      code: isSet(object.code) ? String(object.code) : "",
    };
  },

  toJSON(message: VerifyOTPEmailRegistrationRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.otpEmailId !== undefined && (obj.otpEmailId = message.otpEmailId);
    message.code !== undefined && (obj.code = message.code);
    return obj;
  },

  create(base?: DeepPartial<VerifyOTPEmailRegistrationRequest>): VerifyOTPEmailRegistrationRequest {
    return VerifyOTPEmailRegistrationRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyOTPEmailRegistrationRequest>): VerifyOTPEmailRegistrationRequest {
    const message = createBaseVerifyOTPEmailRegistrationRequest();
    message.userId = object.userId ?? "";
    message.otpEmailId = object.otpEmailId ?? "";
    message.code = object.code ?? "";
    return message;
  },
};

function createBaseVerifyOTPEmailRegistrationResponse(): VerifyOTPEmailRegistrationResponse {
  return { details: undefined };
}

export const VerifyOTPEmailRegistrationResponse = {
  encode(message: VerifyOTPEmailRegistrationResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyOTPEmailRegistrationResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyOTPEmailRegistrationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): VerifyOTPEmailRegistrationResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: VerifyOTPEmailRegistrationResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<VerifyOTPEmailRegistrationResponse>): VerifyOTPEmailRegistrationResponse {
    return VerifyOTPEmailRegistrationResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<VerifyOTPEmailRegistrationResponse>): VerifyOTPEmailRegistrationResponse {
    const message = createBaseVerifyOTPEmailRegistrationResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveOTPEmailAuthenticatorRequest(): RemoveOTPEmailAuthenticatorRequest {
  return { userId: "", otpEmailId: "" };
}

export const RemoveOTPEmailAuthenticatorRequest = {
  encode(message: RemoveOTPEmailAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.otpEmailId !== "") {
      writer.uint32(18).string(message.otpEmailId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOTPEmailAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOTPEmailAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpEmailId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveOTPEmailAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      otpEmailId: isSet(object.otpEmailId) ? String(object.otpEmailId) : "",
    };
  },

  toJSON(message: RemoveOTPEmailAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.otpEmailId !== undefined && (obj.otpEmailId = message.otpEmailId);
    return obj;
  },

  create(base?: DeepPartial<RemoveOTPEmailAuthenticatorRequest>): RemoveOTPEmailAuthenticatorRequest {
    return RemoveOTPEmailAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOTPEmailAuthenticatorRequest>): RemoveOTPEmailAuthenticatorRequest {
    const message = createBaseRemoveOTPEmailAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.otpEmailId = object.otpEmailId ?? "";
    return message;
  },
};

function createBaseRemoveOTPEmailAuthenticatorResponse(): RemoveOTPEmailAuthenticatorResponse {
  return { details: undefined };
}

export const RemoveOTPEmailAuthenticatorResponse = {
  encode(message: RemoveOTPEmailAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveOTPEmailAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveOTPEmailAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveOTPEmailAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveOTPEmailAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveOTPEmailAuthenticatorResponse>): RemoveOTPEmailAuthenticatorResponse {
    return RemoveOTPEmailAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveOTPEmailAuthenticatorResponse>): RemoveOTPEmailAuthenticatorResponse {
    const message = createBaseRemoveOTPEmailAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseStartIdentityProviderIntentRequest(): StartIdentityProviderIntentRequest {
  return { idpId: "", urls: undefined, ldap: undefined };
}

export const StartIdentityProviderIntentRequest = {
  encode(message: StartIdentityProviderIntentRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.urls !== undefined) {
      RedirectURLs.encode(message.urls, writer.uint32(18).fork()).ldelim();
    }
    if (message.ldap !== undefined) {
      LDAPCredentials.encode(message.ldap, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartIdentityProviderIntentRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartIdentityProviderIntentRequest();
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

          message.urls = RedirectURLs.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.ldap = LDAPCredentials.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartIdentityProviderIntentRequest {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      urls: isSet(object.urls) ? RedirectURLs.fromJSON(object.urls) : undefined,
      ldap: isSet(object.ldap) ? LDAPCredentials.fromJSON(object.ldap) : undefined,
    };
  },

  toJSON(message: StartIdentityProviderIntentRequest): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.urls !== undefined && (obj.urls = message.urls ? RedirectURLs.toJSON(message.urls) : undefined);
    message.ldap !== undefined && (obj.ldap = message.ldap ? LDAPCredentials.toJSON(message.ldap) : undefined);
    return obj;
  },

  create(base?: DeepPartial<StartIdentityProviderIntentRequest>): StartIdentityProviderIntentRequest {
    return StartIdentityProviderIntentRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartIdentityProviderIntentRequest>): StartIdentityProviderIntentRequest {
    const message = createBaseStartIdentityProviderIntentRequest();
    message.idpId = object.idpId ?? "";
    message.urls = (object.urls !== undefined && object.urls !== null)
      ? RedirectURLs.fromPartial(object.urls)
      : undefined;
    message.ldap = (object.ldap !== undefined && object.ldap !== null)
      ? LDAPCredentials.fromPartial(object.ldap)
      : undefined;
    return message;
  },
};

function createBaseStartIdentityProviderIntentResponse(): StartIdentityProviderIntentResponse {
  return { details: undefined, authUrl: undefined, idpIntent: undefined, postForm: undefined };
}

export const StartIdentityProviderIntentResponse = {
  encode(message: StartIdentityProviderIntentResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.authUrl !== undefined) {
      writer.uint32(18).string(message.authUrl);
    }
    if (message.idpIntent !== undefined) {
      IdentityProviderIntent.encode(message.idpIntent, writer.uint32(26).fork()).ldelim();
    }
    if (message.postForm !== undefined) {
      writer.uint32(34).bytes(message.postForm);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StartIdentityProviderIntentResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStartIdentityProviderIntentResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.authUrl = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.idpIntent = IdentityProviderIntent.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.postForm = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StartIdentityProviderIntentResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      authUrl: isSet(object.authUrl) ? String(object.authUrl) : undefined,
      idpIntent: isSet(object.idpIntent) ? IdentityProviderIntent.fromJSON(object.idpIntent) : undefined,
      postForm: isSet(object.postForm) ? Buffer.from(bytesFromBase64(object.postForm)) : undefined,
    };
  },

  toJSON(message: StartIdentityProviderIntentResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.authUrl !== undefined && (obj.authUrl = message.authUrl);
    message.idpIntent !== undefined &&
      (obj.idpIntent = message.idpIntent ? IdentityProviderIntent.toJSON(message.idpIntent) : undefined);
    message.postForm !== undefined &&
      (obj.postForm = message.postForm !== undefined ? base64FromBytes(message.postForm) : undefined);
    return obj;
  },

  create(base?: DeepPartial<StartIdentityProviderIntentResponse>): StartIdentityProviderIntentResponse {
    return StartIdentityProviderIntentResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<StartIdentityProviderIntentResponse>): StartIdentityProviderIntentResponse {
    const message = createBaseStartIdentityProviderIntentResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.authUrl = object.authUrl ?? undefined;
    message.idpIntent = (object.idpIntent !== undefined && object.idpIntent !== null)
      ? IdentityProviderIntent.fromPartial(object.idpIntent)
      : undefined;
    message.postForm = object.postForm ?? undefined;
    return message;
  },
};

function createBaseRetrieveIdentityProviderIntentRequest(): RetrieveIdentityProviderIntentRequest {
  return { idpIntentId: "", idpIntentToken: "" };
}

export const RetrieveIdentityProviderIntentRequest = {
  encode(message: RetrieveIdentityProviderIntentRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpIntentId !== "") {
      writer.uint32(10).string(message.idpIntentId);
    }
    if (message.idpIntentToken !== "") {
      writer.uint32(18).string(message.idpIntentToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RetrieveIdentityProviderIntentRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRetrieveIdentityProviderIntentRequest();
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
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RetrieveIdentityProviderIntentRequest {
    return {
      idpIntentId: isSet(object.idpIntentId) ? String(object.idpIntentId) : "",
      idpIntentToken: isSet(object.idpIntentToken) ? String(object.idpIntentToken) : "",
    };
  },

  toJSON(message: RetrieveIdentityProviderIntentRequest): unknown {
    const obj: any = {};
    message.idpIntentId !== undefined && (obj.idpIntentId = message.idpIntentId);
    message.idpIntentToken !== undefined && (obj.idpIntentToken = message.idpIntentToken);
    return obj;
  },

  create(base?: DeepPartial<RetrieveIdentityProviderIntentRequest>): RetrieveIdentityProviderIntentRequest {
    return RetrieveIdentityProviderIntentRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RetrieveIdentityProviderIntentRequest>): RetrieveIdentityProviderIntentRequest {
    const message = createBaseRetrieveIdentityProviderIntentRequest();
    message.idpIntentId = object.idpIntentId ?? "";
    message.idpIntentToken = object.idpIntentToken ?? "";
    return message;
  },
};

function createBaseRetrieveIdentityProviderIntentResponse(): RetrieveIdentityProviderIntentResponse {
  return { details: undefined, idpInformation: undefined, userId: undefined };
}

export const RetrieveIdentityProviderIntentResponse = {
  encode(message: RetrieveIdentityProviderIntentResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    if (message.idpInformation !== undefined) {
      IDPInformation.encode(message.idpInformation, writer.uint32(18).fork()).ldelim();
    }
    if (message.userId !== undefined) {
      writer.uint32(26).string(message.userId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RetrieveIdentityProviderIntentResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRetrieveIdentityProviderIntentResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpInformation = IDPInformation.decode(reader, reader.uint32());
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

  fromJSON(object: any): RetrieveIdentityProviderIntentResponse {
    return {
      details: isSet(object.details) ? Details.fromJSON(object.details) : undefined,
      idpInformation: isSet(object.idpInformation) ? IDPInformation.fromJSON(object.idpInformation) : undefined,
      userId: isSet(object.userId) ? String(object.userId) : undefined,
    };
  },

  toJSON(message: RetrieveIdentityProviderIntentResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    message.idpInformation !== undefined &&
      (obj.idpInformation = message.idpInformation ? IDPInformation.toJSON(message.idpInformation) : undefined);
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  create(base?: DeepPartial<RetrieveIdentityProviderIntentResponse>): RetrieveIdentityProviderIntentResponse {
    return RetrieveIdentityProviderIntentResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RetrieveIdentityProviderIntentResponse>): RetrieveIdentityProviderIntentResponse {
    const message = createBaseRetrieveIdentityProviderIntentResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    message.idpInformation = (object.idpInformation !== undefined && object.idpInformation !== null)
      ? IDPInformation.fromPartial(object.idpInformation)
      : undefined;
    message.userId = object.userId ?? undefined;
    return message;
  },
};

function createBaseAddIDPAuthenticatorRequest(): AddIDPAuthenticatorRequest {
  return { userId: "", idpAuthenticator: undefined };
}

export const AddIDPAuthenticatorRequest = {
  encode(message: AddIDPAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.idpAuthenticator !== undefined) {
      IDPAuthenticator.encode(message.idpAuthenticator, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIDPAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIDPAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpAuthenticator = IDPAuthenticator.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddIDPAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      idpAuthenticator: isSet(object.idpAuthenticator) ? IDPAuthenticator.fromJSON(object.idpAuthenticator) : undefined,
    };
  },

  toJSON(message: AddIDPAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.idpAuthenticator !== undefined &&
      (obj.idpAuthenticator = message.idpAuthenticator ? IDPAuthenticator.toJSON(message.idpAuthenticator) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddIDPAuthenticatorRequest>): AddIDPAuthenticatorRequest {
    return AddIDPAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIDPAuthenticatorRequest>): AddIDPAuthenticatorRequest {
    const message = createBaseAddIDPAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.idpAuthenticator = (object.idpAuthenticator !== undefined && object.idpAuthenticator !== null)
      ? IDPAuthenticator.fromPartial(object.idpAuthenticator)
      : undefined;
    return message;
  },
};

function createBaseAddIDPAuthenticatorResponse(): AddIDPAuthenticatorResponse {
  return { details: undefined };
}

export const AddIDPAuthenticatorResponse = {
  encode(message: AddIDPAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddIDPAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddIDPAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddIDPAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: AddIDPAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<AddIDPAuthenticatorResponse>): AddIDPAuthenticatorResponse {
    return AddIDPAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddIDPAuthenticatorResponse>): AddIDPAuthenticatorResponse {
    const message = createBaseAddIDPAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseRemoveIDPAuthenticatorRequest(): RemoveIDPAuthenticatorRequest {
  return { userId: "", idpId: "" };
}

export const RemoveIDPAuthenticatorRequest = {
  encode(message: RemoveIDPAuthenticatorRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.idpId !== "") {
      writer.uint32(18).string(message.idpId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPAuthenticatorRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPAuthenticatorRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idpId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPAuthenticatorRequest {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
    };
  },

  toJSON(message: RemoveIDPAuthenticatorRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.idpId !== undefined && (obj.idpId = message.idpId);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPAuthenticatorRequest>): RemoveIDPAuthenticatorRequest {
    return RemoveIDPAuthenticatorRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPAuthenticatorRequest>): RemoveIDPAuthenticatorRequest {
    const message = createBaseRemoveIDPAuthenticatorRequest();
    message.userId = object.userId ?? "";
    message.idpId = object.idpId ?? "";
    return message;
  },
};

function createBaseRemoveIDPAuthenticatorResponse(): RemoveIDPAuthenticatorResponse {
  return { details: undefined };
}

export const RemoveIDPAuthenticatorResponse = {
  encode(message: RemoveIDPAuthenticatorResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.details !== undefined) {
      Details.encode(message.details, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveIDPAuthenticatorResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveIDPAuthenticatorResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.details = Details.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveIDPAuthenticatorResponse {
    return { details: isSet(object.details) ? Details.fromJSON(object.details) : undefined };
  },

  toJSON(message: RemoveIDPAuthenticatorResponse): unknown {
    const obj: any = {};
    message.details !== undefined && (obj.details = message.details ? Details.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RemoveIDPAuthenticatorResponse>): RemoveIDPAuthenticatorResponse {
    return RemoveIDPAuthenticatorResponse.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RemoveIDPAuthenticatorResponse>): RemoveIDPAuthenticatorResponse {
    const message = createBaseRemoveIDPAuthenticatorResponse();
    message.details = (object.details !== undefined && object.details !== null)
      ? Details.fromPartial(object.details)
      : undefined;
    return message;
  },
};

export type UserServiceDefinition = typeof UserServiceDefinition;
export const UserServiceDefinition = {
  name: "UserService",
  fullName: "zitadel.user.v3alpha.UserService",
  methods: {
    /**
     * List users
     *
     * List all matching users. By default, we will return all users of your instance.
     * Make sure to include a limit and sorting for pagination.
     */
    listUsers: {
      name: "ListUsers",
      requestType: ListUsersRequest,
      requestStream: false,
      responseType: ListUsersResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              107,
              74,
              47,
              10,
              3,
              50,
              48,
              48,
              18,
              40,
              10,
              38,
              65,
              32,
              108,
              105,
              115,
              116,
              32,
              111,
              102,
              32,
              97,
              108,
              108,
              32,
              117,
              115,
              101,
              114,
              115,
              32,
              109,
              97,
              116,
              99,
              104,
              105,
              110,
              103,
              32,
              116,
              104,
              101,
              32,
              113,
              117,
              101,
              114,
              121,
              74,
              56,
              10,
              3,
              52,
              48,
              48,
              18,
              49,
              10,
              18,
              105,
              110,
              118,
              97,
              108,
              105,
              100,
              32,
              108,
              105,
              115,
              116,
              32,
              113,
              117,
              101,
              114,
              121,
              18,
              27,
              10,
              25,
              26,
              23,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              114,
              112,
              99,
              83,
              116,
              97,
              116,
              117,
              115,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              26,
              58,
              1,
              42,
              34,
              21,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              115,
              101,
              97,
              114,
              99,
              104,
            ]),
          ],
        },
      },
    },
    /**
     * User by ID
     *
     * Returns the user identified by the requested ID.
     */
    getUserByID: {
      name: "GetUserByID",
      requestType: GetUserByIDRequest,
      requestStream: false,
      responseType: GetUserByIDResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              38,
              74,
              36,
              10,
              3,
              50,
              48,
              48,
              18,
              29,
              10,
              27,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              116,
              114,
              105,
              101,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              26,
              18,
              24,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Create a user
     *
     * Create a new user with an optional data schema.
     */
    createUser: {
      name: "CreateUser",
      requestType: CreateUserRequest,
      requestStream: false,
      responseType: CreateUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              81,
              74,
              79,
              10,
              3,
              50,
              48,
              49,
              18,
              72,
              10,
              25,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              99,
              114,
              101,
              97,
              116,
              101,
              100,
              18,
              43,
              10,
              41,
              26,
              39,
              35,
              47,
              100,
              101,
              102,
              105,
              110,
              105,
              116,
              105,
              111,
              110,
              115,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              67,
              114,
              101,
              97,
              116,
              101,
              85,
              115,
              101,
              114,
              82,
              101,
              115,
              112,
              111,
              110,
              115,
              101,
            ]),
          ],
          400010: [
            Buffer.from([
              33,
              10,
              26,
              10,
              10,
              117,
              115,
              101,
              114,
              46,
              119,
              114,
              105,
              116,
              101,
              26,
              12,
              111,
              114,
              103,
              97,
              110,
              105,
              122,
              97,
              116,
              105,
              111,
              110,
              18,
              3,
              8,
              201,
              1,
            ]),
          ],
          578365826: [
            Buffer.from([19, 58, 1, 42, 34, 14, 47, 118, 51, 97, 108, 112, 104, 97, 47, 117, 115, 101, 114, 115]),
          ],
        },
      },
    },
    /**
     * Update a user
     *
     * Update an existing user with data based on a user schema.
     */
    updateUser: {
      name: "UpdateUser",
      requestType: UpdateUserRequest,
      requestStream: false,
      responseType: UpdateUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              36,
              74,
              34,
              10,
              3,
              50,
              48,
              48,
              18,
              27,
              10,
              25,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              117,
              112,
              100,
              97,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              29,
              58,
              1,
              42,
              26,
              24,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Deactivate a user
     *
     * Deactivate an existing user and change the state 'deactivated'.
     * The user will not be able to log in anymore.
     * Use deactivate user when the user should not be able to use the account anymore,
     * but you still need access to the user data.
     *
     * The endpoint returns an error if the user is already in the state 'deactivated'.
     */
    deactivateUser: {
      name: "DeactivateUser",
      requestType: DeactivateUserRequest,
      requestStream: false,
      responseType: DeactivateUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              40,
              74,
              38,
              10,
              3,
              50,
              48,
              48,
              18,
              31,
              10,
              29,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              100,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              37,
              34,
              35,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              100,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Reactivate a user
     *
     * Reactivate a previously deactivated user and change the state to 'active'.
     * The user will be able to log in again.
     *
     * The endpoint returns an error if the user is not in the state 'deactivated'.
     */
    reactivateUser: {
      name: "ReactivateUser",
      requestType: ReactivateUserRequest,
      requestStream: false,
      responseType: ReactivateUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              40,
              74,
              38,
              10,
              3,
              50,
              48,
              48,
              18,
              31,
              10,
              29,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              37,
              34,
              35,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              114,
              101,
              97,
              99,
              116,
              105,
              118,
              97,
              116,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Lock a user
     *
     * Lock an existing user and change the state 'locked'.
     * The user will not be able to log in anymore.
     * Use lock user when the user should temporarily not be able to log in
     * because of an event that happened (wrong password, etc.)
     *
     * The endpoint returns an error if the user is already in the state 'locked'.
     */
    lockUser: {
      name: "LockUser",
      requestType: LockUserRequest,
      requestStream: false,
      responseType: LockUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              35,
              74,
              33,
              10,
              3,
              50,
              48,
              48,
              18,
              26,
              10,
              24,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              108,
              111,
              99,
              107,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              31,
              34,
              29,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              108,
              111,
              99,
              107,
            ]),
          ],
        },
      },
    },
    /**
     * Unlock a user
     *
     * Unlock a previously locked user and change the state to 'active'.
     * The user will be able to log in again.
     *
     * The endpoint returns an error if the user is not in the state 'locked'.
     */
    unlockUser: {
      name: "UnlockUser",
      requestType: UnlockUserRequest,
      requestStream: false,
      responseType: UnlockUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              37,
              74,
              35,
              10,
              3,
              50,
              48,
              48,
              18,
              28,
              10,
              26,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              117,
              110,
              108,
              111,
              99,
              107,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              33,
              34,
              31,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              117,
              110,
              108,
              111,
              99,
              107,
            ]),
          ],
        },
      },
    },
    /**
     * Delete a user
     *
     * Delete an existing user and change the state to 'deleted'.
     * The user will be able to log in anymore.
     */
    deleteUser: {
      name: "DeleteUser",
      requestType: DeleteUserRequest,
      requestStream: false,
      responseType: DeleteUserResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              36,
              74,
              34,
              10,
              3,
              50,
              48,
              48,
              18,
              27,
              10,
              25,
              85,
              115,
              101,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              100,
              101,
              108,
              101,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              26,
              42,
              24,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Set contact email
     *
     * Add or update the contact email address of a user.
     * If the email is not passed as verified, a verification code will be generated,
     * which can be either returned or will be sent to the user by email.
     */
    setContactEmail: {
      name: "SetContactEmail",
      requestType: SetContactEmailRequest,
      requestStream: false,
      responseType: SetContactEmailResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              33,
              74,
              31,
              10,
              3,
              50,
              48,
              48,
              18,
              24,
              10,
              22,
              69,
              109,
              97,
              105,
              108,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              101,
              116,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              35,
              58,
              1,
              42,
              26,
              30,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              101,
              109,
              97,
              105,
              108,
            ]),
          ],
        },
      },
    },
    /**
     * Verify the contact email
     *
     * Verify the contact email with the provided code.
     */
    verifyContactEmail: {
      name: "VerifyContactEmail",
      requestType: VerifyContactEmailRequest,
      requestStream: false,
      responseType: VerifyContactEmailResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              38,
              74,
              36,
              10,
              3,
              50,
              48,
              48,
              18,
              29,
              10,
              27,
              69,
              109,
              97,
              105,
              108,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              42,
              58,
              1,
              42,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              101,
              109,
              97,
              105,
              108,
              47,
              118,
              101,
              114,
              105,
              102,
              121,
            ]),
          ],
        },
      },
    },
    /**
     * Resend the contact email code
     *
     * Resend the email with the verification code for the contact email address.
     */
    resendContactEmailCode: {
      name: "ResendContactEmailCode",
      requestType: ResendContactEmailCodeRequest,
      requestStream: false,
      responseType: ResendContactEmailCodeResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              45,
              74,
              43,
              10,
              3,
              50,
              48,
              48,
              18,
              36,
              10,
              34,
              67,
              111,
              100,
              101,
              32,
              114,
              101,
              115,
              101,
              110,
              100,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              42,
              58,
              1,
              42,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              101,
              109,
              97,
              105,
              108,
              47,
              114,
              101,
              115,
              101,
              110,
              100,
            ]),
          ],
        },
      },
    },
    /**
     * Set contact phone
     *
     * Add or update the contact phone number of a user.
     * If the phone is not passed as verified, a verification code will be generated,
     * which can be either returned or will be sent to the user by SMS.
     */
    setContactPhone: {
      name: "SetContactPhone",
      requestType: SetContactPhoneRequest,
      requestStream: false,
      responseType: SetContactPhoneResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              41,
              74,
              39,
              10,
              3,
              50,
              48,
              48,
              18,
              32,
              10,
              30,
              67,
              111,
              110,
              116,
              97,
              99,
              116,
              32,
              112,
              104,
              111,
              110,
              101,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              101,
              116,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              35,
              58,
              1,
              42,
              26,
              30,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              112,
              104,
              111,
              110,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Verify the contact phone
     *
     * Verify the contact phone with the provided code.
     */
    verifyContactPhone: {
      name: "VerifyContactPhone",
      requestType: VerifyContactPhoneRequest,
      requestStream: false,
      responseType: VerifyContactPhoneResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              46,
              74,
              44,
              10,
              3,
              50,
              48,
              48,
              18,
              37,
              10,
              35,
              67,
              111,
              110,
              116,
              97,
              99,
              116,
              32,
              112,
              104,
              111,
              110,
              101,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              42,
              58,
              1,
              42,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              112,
              104,
              111,
              110,
              101,
              47,
              118,
              101,
              114,
              105,
              102,
              121,
            ]),
          ],
        },
      },
    },
    /**
     * Resend the contact phone code
     *
     * Resend the phone with the verification code for the contact phone number.
     */
    resendContactPhoneCode: {
      name: "ResendContactPhoneCode",
      requestType: ResendContactPhoneCodeRequest,
      requestStream: false,
      responseType: ResendContactPhoneCodeResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              45,
              74,
              43,
              10,
              3,
              50,
              48,
              48,
              18,
              36,
              10,
              34,
              67,
              111,
              100,
              101,
              32,
              114,
              101,
              115,
              101,
              110,
              100,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              42,
              58,
              1,
              42,
              34,
              37,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              112,
              104,
              111,
              110,
              101,
              47,
              114,
              101,
              115,
              101,
              110,
              100,
            ]),
          ],
        },
      },
    },
    /**
     * Add a username
     *
     * Add a new unique username to a user. The username will be used to identify the user on authentication.
     */
    addUsername: {
      name: "AddUsername",
      requestType: AddUsernameRequest,
      requestStream: false,
      responseType: AddUsernameResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              38,
              74,
              36,
              10,
              3,
              50,
              48,
              48,
              18,
              29,
              10,
              27,
              85,
              115,
              101,
              114,
              110,
              97,
              109,
              101,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              97,
              100,
              100,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              38,
              58,
              1,
              42,
              34,
              33,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              117,
              115,
              101,
              114,
              110,
              97,
              109,
              101,
            ]),
          ],
        },
      },
    },
    /**
     * Remove a username
     *
     * Remove an existing username of a user, so it cannot be used for authentication anymore.
     */
    removeUsername: {
      name: "RemoveUsername",
      requestType: RemoveUsernameRequest,
      requestStream: false,
      responseType: RemoveUsernameResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              40,
              74,
              38,
              10,
              3,
              50,
              48,
              48,
              18,
              31,
              10,
              29,
              85,
              115,
              101,
              114,
              110,
              97,
              109,
              101,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              49,
              42,
              47,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              117,
              115,
              101,
              114,
              110,
              97,
              109,
              101,
              47,
              123,
              117,
              115,
              101,
              114,
              110,
              97,
              109,
              101,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Set a password
     *
     * Add, update or reset a user's password with either a verification code or the current password.
     */
    setPassword: {
      name: "SetPassword",
      requestType: SetPasswordRequest,
      requestStream: false,
      responseType: SetPasswordResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              36,
              74,
              34,
              10,
              3,
              50,
              48,
              48,
              18,
              27,
              10,
              25,
              80,
              97,
              115,
              115,
              119,
              111,
              114,
              100,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              101,
              116,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              37,
              58,
              1,
              42,
              34,
              32,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              112,
              97,
              115,
              115,
              119,
              111,
              114,
              100,
            ]),
          ],
        },
      },
    },
    /**
     * Request password reset
     *
     * Request a code to be able to set a new password.
     */
    requestPasswordReset: {
      name: "RequestPasswordReset",
      requestType: RequestPasswordResetRequest,
      requestStream: false,
      responseType: RequestPasswordResetResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              48,
              74,
              46,
              10,
              3,
              50,
              48,
              48,
              18,
              39,
              10,
              37,
              80,
              97,
              115,
              115,
              119,
              111,
              114,
              100,
              32,
              114,
              101,
              115,
              101,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              113,
              117,
              101,
              115,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              43,
              58,
              1,
              42,
              34,
              38,
              47,
              118,
              50,
              98,
              101,
              116,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              112,
              97,
              115,
              115,
              119,
              111,
              114,
              100,
              47,
              114,
              101,
              115,
              101,
              116,
            ]),
          ],
        },
      },
    },
    /**
     * Start a WebAuthN registration
     *
     * Start the registration of a new WebAuthN device (e.g. Passkeys) for a user.
     * As a response the public key credential creation options are returned,
     * which are used to verify the device.
     */
    startWebAuthNRegistration: {
      name: "StartWebAuthNRegistration",
      requestType: StartWebAuthNRegistrationRequest,
      requestStream: false,
      responseType: StartWebAuthNRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              53,
              74,
              51,
              10,
              3,
              50,
              48,
              48,
              18,
              44,
              10,
              42,
              87,
              101,
              98,
              65,
              117,
              116,
              104,
              78,
              32,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              116,
              97,
              114,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              38,
              58,
              1,
              42,
              34,
              33,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              119,
              101,
              98,
              97,
              117,
              116,
              104,
              110,
            ]),
          ],
        },
      },
    },
    /**
     * Verify a WebAuthN registration
     *
     * Verify the WebAuthN registration started by StartWebAuthNRegistration with the public key credential.
     */
    verifyWebAuthNRegistration: {
      name: "VerifyWebAuthNRegistration",
      requestType: VerifyWebAuthNRegistrationRequest,
      requestStream: false,
      responseType: VerifyWebAuthNRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              54,
              74,
              52,
              10,
              3,
              50,
              48,
              48,
              18,
              45,
              10,
              43,
              87,
              101,
              98,
              65,
              117,
              116,
              104,
              78,
              32,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              54,
              58,
              1,
              42,
              34,
              49,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              119,
              101,
              98,
              97,
              117,
              116,
              104,
              110,
              47,
              123,
              119,
              101,
              98,
              95,
              97,
              117,
              116,
              104,
              95,
              110,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Create a WebAuthN registration link
     *
     * Create a link, which includes a code, that can either be returned or directly sent to the user.
     * The code will allow the user to start a new WebAuthN registration.
     */
    createWebAuthNRegistrationLink: {
      name: "CreateWebAuthNRegistrationLink",
      requestType: CreateWebAuthNRegistrationLinkRequest,
      requestStream: false,
      responseType: CreateWebAuthNRegistrationLinkResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              53,
              74,
              51,
              10,
              3,
              50,
              48,
              48,
              18,
              44,
              10,
              42,
              87,
              101,
              98,
              65,
              117,
              116,
              104,
              78,
              32,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              99,
              114,
              101,
              97,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              56,
              58,
              1,
              42,
              34,
              51,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              119,
              101,
              98,
              97,
              117,
              116,
              104,
              110,
              47,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              95,
              108,
              105,
              110,
              107,
            ]),
          ],
        },
      },
    },
    /**
     * Remove a WebAuthN authenticator
     *
     * Remove an existing WebAuthN authenticator from a user, so it cannot be used for authentication anymore.
     */
    removeWebAuthNAuthenticator: {
      name: "RemoveWebAuthNAuthenticator",
      requestType: RemoveWebAuthNAuthenticatorRequest,
      requestStream: false,
      responseType: RemoveWebAuthNAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              54,
              74,
              52,
              10,
              3,
              50,
              48,
              48,
              18,
              45,
              10,
              43,
              87,
              101,
              98,
              65,
              117,
              116,
              104,
              78,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              51,
              42,
              49,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              119,
              101,
              98,
              97,
              117,
              116,
              104,
              110,
              47,
              123,
              119,
              101,
              98,
              95,
              97,
              117,
              116,
              104,
              95,
              110,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Start a TOTP registration
     *
     * Start the registration of a new time-based one-time-password (TOTP) generator for a user.
     * As a response a secret is returned, which is used to initialize a TOTP app or device.
     */
    startTOTPRegistration: {
      name: "StartTOTPRegistration",
      requestType: StartTOTPRegistrationRequest,
      requestStream: false,
      responseType: StartTOTPRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              49,
              74,
              47,
              10,
              3,
              50,
              48,
              48,
              18,
              40,
              10,
              38,
              84,
              79,
              84,
              80,
              32,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              116,
              97,
              114,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              31,
              34,
              29,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              116,
              111,
              116,
              112,
            ]),
          ],
        },
      },
    },
    /**
     * Verify a TOTP registration
     *
     * Verify the time-based one-time-password (TOTP) registration with the generated code.
     */
    verifyTOTPRegistration: {
      name: "VerifyTOTPRegistration",
      requestType: VerifyTOTPRegistrationRequest,
      requestStream: false,
      responseType: VerifyTOTPRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              50,
              74,
              48,
              10,
              3,
              50,
              48,
              48,
              18,
              41,
              10,
              39,
              84,
              79,
              84,
              80,
              32,
              114,
              101,
              103,
              105,
              115,
              116,
              114,
              97,
              116,
              105,
              111,
              110,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              51,
              58,
              1,
              42,
              34,
              46,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              116,
              111,
              116,
              112,
              47,
              123,
              116,
              111,
              116,
              112,
              95,
              105,
              100,
              125,
              47,
              118,
              101,
              114,
              105,
              102,
              121,
            ]),
          ],
        },
      },
    },
    /**
     * Remove a TOTP authenticator
     *
     * Remove an existing time-based one-time-password (TOTP) authenticator from a user, so it cannot be used for authentication anymore.
     */
    removeTOTPAuthenticator: {
      name: "RemoveTOTPAuthenticator",
      requestType: RemoveTOTPAuthenticatorRequest,
      requestStream: false,
      responseType: RemoveTOTPAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              50,
              74,
              48,
              10,
              3,
              50,
              48,
              48,
              18,
              41,
              10,
              39,
              84,
              79,
              84,
              80,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              41,
              42,
              39,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              116,
              111,
              116,
              112,
              47,
              123,
              116,
              111,
              116,
              112,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Add a OTP SMS authenticator
     *
     * Add a new one-time-password (OTP) SMS authenticator to a user.
     * If the phone is not passed as verified, a verification code will be generated,
     * which can be either returned or will be sent to the user by SMS.
     */
    addOTPSMSAuthenticator: {
      name: "AddOTPSMSAuthenticator",
      requestType: AddOTPSMSAuthenticatorRequest,
      requestStream: false,
      responseType: AddOTPSMSAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              51,
              74,
              49,
              10,
              3,
              50,
              48,
              48,
              18,
              42,
              10,
              40,
              79,
              84,
              80,
              32,
              83,
              77,
              83,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              97,
              100,
              100,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              37,
              58,
              1,
              42,
              34,
              32,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              115,
              109,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Verify OTP SMS registration
     *
     * Verify the OTP SMS registration with the provided code.
     */
    verifyOTPSMSRegistration: {
      name: "VerifyOTPSMSRegistration",
      requestType: VerifyOTPSMSRegistrationRequest,
      requestStream: false,
      responseType: VerifyOTPSMSRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              54,
              74,
              52,
              10,
              3,
              50,
              48,
              48,
              18,
              45,
              10,
              43,
              79,
              84,
              80,
              32,
              83,
              77,
              83,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              57,
              58,
              1,
              42,
              34,
              52,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              115,
              109,
              115,
              47,
              123,
              111,
              116,
              112,
              95,
              115,
              109,
              115,
              95,
              105,
              100,
              125,
              47,
              118,
              101,
              114,
              105,
              102,
              121,
            ]),
          ],
        },
      },
    },
    /**
     * Remove a OTP SMS authenticator
     *
     * Remove an existing one-time-password (OTP) SMS authenticator from a user, so it cannot be used for authentication anymore.
     */
    removeOTPSMSAuthenticator: {
      name: "RemoveOTPSMSAuthenticator",
      requestType: RemoveOTPSMSAuthenticatorRequest,
      requestStream: false,
      responseType: RemoveOTPSMSAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              53,
              74,
              51,
              10,
              3,
              50,
              48,
              48,
              18,
              44,
              10,
              42,
              79,
              84,
              80,
              32,
              83,
              77,
              83,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              47,
              42,
              45,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              115,
              109,
              115,
              47,
              123,
              111,
              116,
              112,
              95,
              115,
              109,
              115,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Add a OTP Email authenticator
     *
     * Add a new one-time-password (OTP) Email authenticator to a user.
     * If the email is not passed as verified, a verification code will be generated,
     * which can be either returned or will be sent to the user by email.
     */
    addOTPEmailAuthenticator: {
      name: "AddOTPEmailAuthenticator",
      requestType: AddOTPEmailAuthenticatorRequest,
      requestStream: false,
      responseType: AddOTPEmailAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              53,
              74,
              51,
              10,
              3,
              50,
              48,
              48,
              18,
              44,
              10,
              42,
              79,
              84,
              80,
              32,
              69,
              109,
              97,
              105,
              108,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              97,
              100,
              100,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              39,
              58,
              1,
              42,
              34,
              34,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              101,
              109,
              97,
              105,
              108,
            ]),
          ],
        },
      },
    },
    /**
     * Verify OTP Email registration
     *
     * Verify the OTP Email registration with the provided code.
     */
    verifyOTPEmailRegistration: {
      name: "VerifyOTPEmailRegistration",
      requestType: VerifyOTPEmailRegistrationRequest,
      requestStream: false,
      responseType: VerifyOTPEmailRegistrationResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              56,
              74,
              54,
              10,
              3,
              50,
              48,
              48,
              18,
              47,
              10,
              45,
              79,
              84,
              80,
              32,
              69,
              109,
              97,
              105,
              108,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              118,
              101,
              114,
              105,
              102,
              105,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              61,
              58,
              1,
              42,
              34,
              56,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              101,
              109,
              97,
              105,
              108,
              47,
              123,
              111,
              116,
              112,
              95,
              101,
              109,
              97,
              105,
              108,
              95,
              105,
              100,
              125,
              47,
              118,
              101,
              114,
              105,
              102,
              121,
            ]),
          ],
        },
      },
    },
    /**
     * Remove a OTP Email authenticator
     *
     * Remove an existing one-time-password (OTP) Email authenticator from a user, so it cannot be used for authentication anymore.
     */
    removeOTPEmailAuthenticator: {
      name: "RemoveOTPEmailAuthenticator",
      requestType: RemoveOTPEmailAuthenticatorRequest,
      requestStream: false,
      responseType: RemoveOTPEmailAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              55,
              74,
              53,
              10,
              3,
              50,
              48,
              48,
              18,
              46,
              10,
              44,
              79,
              84,
              80,
              32,
              69,
              109,
              97,
              105,
              108,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              51,
              42,
              49,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              111,
              116,
              112,
              95,
              101,
              109,
              97,
              105,
              108,
              47,
              123,
              111,
              116,
              112,
              95,
              101,
              109,
              97,
              105,
              108,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Start an IDP authentication intent
     *
     * Start a new authentication intent on configured identity provider (IDP) for external login, registration or linking.
     */
    startIdentityProviderIntent: {
      name: "StartIdentityProviderIntent",
      requestType: StartIdentityProviderIntentRequest,
      requestStream: false,
      responseType: StartIdentityProviderIntentResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              42,
              74,
              40,
              10,
              3,
              50,
              48,
              48,
              18,
              33,
              10,
              31,
              73,
              68,
              80,
              32,
              105,
              110,
              116,
              101,
              110,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              115,
              116,
              97,
              114,
              116,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              25,
              58,
              1,
              42,
              34,
              20,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              105,
              100,
              112,
              95,
              105,
              110,
              116,
              101,
              110,
              116,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Retrieve the information of the IDP authentication intent
     *
     * Retrieve the information returned by the identity provider (IDP) for registration or updating an existing user with new information.
     */
    retrieveIdentityProviderIntent: {
      name: "RetrieveIdentityProviderIntent",
      requestType: RetrieveIdentityProviderIntentRequest,
      requestStream: false,
      responseType: RetrieveIdentityProviderIntentResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              44,
              74,
              42,
              10,
              3,
              50,
              48,
              48,
              18,
              35,
              10,
              33,
              73,
              68,
              80,
              32,
              105,
              110,
              116,
              101,
              110,
              116,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              116,
              114,
              105,
              101,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              41,
              58,
              1,
              42,
              34,
              36,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              105,
              100,
              112,
              95,
              105,
              110,
              116,
              101,
              110,
              116,
              115,
              47,
              123,
              105,
              100,
              112,
              95,
              105,
              110,
              116,
              101,
              110,
              116,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
    /**
     * Add an IDP authenticator to a user
     *
     * Add a new identity provider (IDP) authenticator to an existing user.
     * This will allow the user to authenticate with the provided IDP.
     */
    addIDPAuthenticator: {
      name: "AddIDPAuthenticator",
      requestType: AddIDPAuthenticatorRequest,
      requestStream: false,
      responseType: AddIDPAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              47,
              74,
              45,
              10,
              3,
              50,
              48,
              48,
              18,
              38,
              10,
              36,
              73,
              68,
              80,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              97,
              100,
              100,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              34,
              58,
              1,
              42,
              34,
              29,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              105,
              100,
              112,
              115,
            ]),
          ],
        },
      },
    },
    /**
     * Remove an IDP authenticator
     *
     * Remove an existing identity provider (IDP) authenticator from a user, so it cannot be used for authentication anymore.
     */
    removeIDPAuthenticator: {
      name: "RemoveIDPAuthenticator",
      requestType: RemoveIDPAuthenticatorRequest,
      requestStream: false,
      responseType: RemoveIDPAuthenticatorResponse,
      responseStream: false,
      options: {
        _unknownFields: {
          8338: [
            Buffer.from([
              49,
              74,
              47,
              10,
              3,
              50,
              48,
              48,
              18,
              40,
              10,
              38,
              73,
              68,
              80,
              32,
              97,
              117,
              116,
              104,
              101,
              110,
              116,
              105,
              99,
              97,
              116,
              111,
              114,
              32,
              115,
              117,
              99,
              99,
              101,
              115,
              115,
              102,
              117,
              108,
              108,
              121,
              32,
              114,
              101,
              109,
              111,
              118,
              101,
              100,
            ]),
          ],
          400010: [Buffer.from([17, 10, 15, 10, 13, 97, 117, 116, 104, 101, 110, 116, 105, 99, 97, 116, 101, 100])],
          578365826: [
            Buffer.from([
              40,
              42,
              38,
              47,
              118,
              51,
              97,
              108,
              112,
              104,
              97,
              47,
              117,
              115,
              101,
              114,
              115,
              47,
              123,
              117,
              115,
              101,
              114,
              95,
              105,
              100,
              125,
              47,
              105,
              100,
              112,
              115,
              47,
              123,
              105,
              100,
              112,
              95,
              105,
              100,
              125,
            ]),
          ],
        },
      },
    },
  },
} as const;

export interface UserServiceImplementation<CallContextExt = {}> {
  /**
   * List users
   *
   * List all matching users. By default, we will return all users of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listUsers(request: ListUsersRequest, context: CallContext & CallContextExt): Promise<DeepPartial<ListUsersResponse>>;
  /**
   * User by ID
   *
   * Returns the user identified by the requested ID.
   */
  getUserByID(
    request: GetUserByIDRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<GetUserByIDResponse>>;
  /**
   * Create a user
   *
   * Create a new user with an optional data schema.
   */
  createUser(
    request: CreateUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateUserResponse>>;
  /**
   * Update a user
   *
   * Update an existing user with data based on a user schema.
   */
  updateUser(
    request: UpdateUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UpdateUserResponse>>;
  /**
   * Deactivate a user
   *
   * Deactivate an existing user and change the state 'deactivated'.
   * The user will not be able to log in anymore.
   * Use deactivate user when the user should not be able to use the account anymore,
   * but you still need access to the user data.
   *
   * The endpoint returns an error if the user is already in the state 'deactivated'.
   */
  deactivateUser(
    request: DeactivateUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeactivateUserResponse>>;
  /**
   * Reactivate a user
   *
   * Reactivate a previously deactivated user and change the state to 'active'.
   * The user will be able to log in again.
   *
   * The endpoint returns an error if the user is not in the state 'deactivated'.
   */
  reactivateUser(
    request: ReactivateUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ReactivateUserResponse>>;
  /**
   * Lock a user
   *
   * Lock an existing user and change the state 'locked'.
   * The user will not be able to log in anymore.
   * Use lock user when the user should temporarily not be able to log in
   * because of an event that happened (wrong password, etc.)
   *
   * The endpoint returns an error if the user is already in the state 'locked'.
   */
  lockUser(request: LockUserRequest, context: CallContext & CallContextExt): Promise<DeepPartial<LockUserResponse>>;
  /**
   * Unlock a user
   *
   * Unlock a previously locked user and change the state to 'active'.
   * The user will be able to log in again.
   *
   * The endpoint returns an error if the user is not in the state 'locked'.
   */
  unlockUser(
    request: UnlockUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<UnlockUserResponse>>;
  /**
   * Delete a user
   *
   * Delete an existing user and change the state to 'deleted'.
   * The user will be able to log in anymore.
   */
  deleteUser(
    request: DeleteUserRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<DeleteUserResponse>>;
  /**
   * Set contact email
   *
   * Add or update the contact email address of a user.
   * If the email is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by email.
   */
  setContactEmail(
    request: SetContactEmailRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetContactEmailResponse>>;
  /**
   * Verify the contact email
   *
   * Verify the contact email with the provided code.
   */
  verifyContactEmail(
    request: VerifyContactEmailRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyContactEmailResponse>>;
  /**
   * Resend the contact email code
   *
   * Resend the email with the verification code for the contact email address.
   */
  resendContactEmailCode(
    request: ResendContactEmailCodeRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResendContactEmailCodeResponse>>;
  /**
   * Set contact phone
   *
   * Add or update the contact phone number of a user.
   * If the phone is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by SMS.
   */
  setContactPhone(
    request: SetContactPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetContactPhoneResponse>>;
  /**
   * Verify the contact phone
   *
   * Verify the contact phone with the provided code.
   */
  verifyContactPhone(
    request: VerifyContactPhoneRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyContactPhoneResponse>>;
  /**
   * Resend the contact phone code
   *
   * Resend the phone with the verification code for the contact phone number.
   */
  resendContactPhoneCode(
    request: ResendContactPhoneCodeRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<ResendContactPhoneCodeResponse>>;
  /**
   * Add a username
   *
   * Add a new unique username to a user. The username will be used to identify the user on authentication.
   */
  addUsername(
    request: AddUsernameRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddUsernameResponse>>;
  /**
   * Remove a username
   *
   * Remove an existing username of a user, so it cannot be used for authentication anymore.
   */
  removeUsername(
    request: RemoveUsernameRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveUsernameResponse>>;
  /**
   * Set a password
   *
   * Add, update or reset a user's password with either a verification code or the current password.
   */
  setPassword(
    request: SetPasswordRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<SetPasswordResponse>>;
  /**
   * Request password reset
   *
   * Request a code to be able to set a new password.
   */
  requestPasswordReset(
    request: RequestPasswordResetRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RequestPasswordResetResponse>>;
  /**
   * Start a WebAuthN registration
   *
   * Start the registration of a new WebAuthN device (e.g. Passkeys) for a user.
   * As a response the public key credential creation options are returned,
   * which are used to verify the device.
   */
  startWebAuthNRegistration(
    request: StartWebAuthNRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<StartWebAuthNRegistrationResponse>>;
  /**
   * Verify a WebAuthN registration
   *
   * Verify the WebAuthN registration started by StartWebAuthNRegistration with the public key credential.
   */
  verifyWebAuthNRegistration(
    request: VerifyWebAuthNRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyWebAuthNRegistrationResponse>>;
  /**
   * Create a WebAuthN registration link
   *
   * Create a link, which includes a code, that can either be returned or directly sent to the user.
   * The code will allow the user to start a new WebAuthN registration.
   */
  createWebAuthNRegistrationLink(
    request: CreateWebAuthNRegistrationLinkRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<CreateWebAuthNRegistrationLinkResponse>>;
  /**
   * Remove a WebAuthN authenticator
   *
   * Remove an existing WebAuthN authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeWebAuthNAuthenticator(
    request: RemoveWebAuthNAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveWebAuthNAuthenticatorResponse>>;
  /**
   * Start a TOTP registration
   *
   * Start the registration of a new time-based one-time-password (TOTP) generator for a user.
   * As a response a secret is returned, which is used to initialize a TOTP app or device.
   */
  startTOTPRegistration(
    request: StartTOTPRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<StartTOTPRegistrationResponse>>;
  /**
   * Verify a TOTP registration
   *
   * Verify the time-based one-time-password (TOTP) registration with the generated code.
   */
  verifyTOTPRegistration(
    request: VerifyTOTPRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyTOTPRegistrationResponse>>;
  /**
   * Remove a TOTP authenticator
   *
   * Remove an existing time-based one-time-password (TOTP) authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeTOTPAuthenticator(
    request: RemoveTOTPAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveTOTPAuthenticatorResponse>>;
  /**
   * Add a OTP SMS authenticator
   *
   * Add a new one-time-password (OTP) SMS authenticator to a user.
   * If the phone is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by SMS.
   */
  addOTPSMSAuthenticator(
    request: AddOTPSMSAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddOTPSMSAuthenticatorResponse>>;
  /**
   * Verify OTP SMS registration
   *
   * Verify the OTP SMS registration with the provided code.
   */
  verifyOTPSMSRegistration(
    request: VerifyOTPSMSRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyOTPSMSRegistrationResponse>>;
  /**
   * Remove a OTP SMS authenticator
   *
   * Remove an existing one-time-password (OTP) SMS authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeOTPSMSAuthenticator(
    request: RemoveOTPSMSAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveOTPSMSAuthenticatorResponse>>;
  /**
   * Add a OTP Email authenticator
   *
   * Add a new one-time-password (OTP) Email authenticator to a user.
   * If the email is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by email.
   */
  addOTPEmailAuthenticator(
    request: AddOTPEmailAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddOTPEmailAuthenticatorResponse>>;
  /**
   * Verify OTP Email registration
   *
   * Verify the OTP Email registration with the provided code.
   */
  verifyOTPEmailRegistration(
    request: VerifyOTPEmailRegistrationRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<VerifyOTPEmailRegistrationResponse>>;
  /**
   * Remove a OTP Email authenticator
   *
   * Remove an existing one-time-password (OTP) Email authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeOTPEmailAuthenticator(
    request: RemoveOTPEmailAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveOTPEmailAuthenticatorResponse>>;
  /**
   * Start an IDP authentication intent
   *
   * Start a new authentication intent on configured identity provider (IDP) for external login, registration or linking.
   */
  startIdentityProviderIntent(
    request: StartIdentityProviderIntentRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<StartIdentityProviderIntentResponse>>;
  /**
   * Retrieve the information of the IDP authentication intent
   *
   * Retrieve the information returned by the identity provider (IDP) for registration or updating an existing user with new information.
   */
  retrieveIdentityProviderIntent(
    request: RetrieveIdentityProviderIntentRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RetrieveIdentityProviderIntentResponse>>;
  /**
   * Add an IDP authenticator to a user
   *
   * Add a new identity provider (IDP) authenticator to an existing user.
   * This will allow the user to authenticate with the provided IDP.
   */
  addIDPAuthenticator(
    request: AddIDPAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<AddIDPAuthenticatorResponse>>;
  /**
   * Remove an IDP authenticator
   *
   * Remove an existing identity provider (IDP) authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeIDPAuthenticator(
    request: RemoveIDPAuthenticatorRequest,
    context: CallContext & CallContextExt,
  ): Promise<DeepPartial<RemoveIDPAuthenticatorResponse>>;
}

export interface UserServiceClient<CallOptionsExt = {}> {
  /**
   * List users
   *
   * List all matching users. By default, we will return all users of your instance.
   * Make sure to include a limit and sorting for pagination.
   */
  listUsers(request: DeepPartial<ListUsersRequest>, options?: CallOptions & CallOptionsExt): Promise<ListUsersResponse>;
  /**
   * User by ID
   *
   * Returns the user identified by the requested ID.
   */
  getUserByID(
    request: DeepPartial<GetUserByIDRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<GetUserByIDResponse>;
  /**
   * Create a user
   *
   * Create a new user with an optional data schema.
   */
  createUser(
    request: DeepPartial<CreateUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateUserResponse>;
  /**
   * Update a user
   *
   * Update an existing user with data based on a user schema.
   */
  updateUser(
    request: DeepPartial<UpdateUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UpdateUserResponse>;
  /**
   * Deactivate a user
   *
   * Deactivate an existing user and change the state 'deactivated'.
   * The user will not be able to log in anymore.
   * Use deactivate user when the user should not be able to use the account anymore,
   * but you still need access to the user data.
   *
   * The endpoint returns an error if the user is already in the state 'deactivated'.
   */
  deactivateUser(
    request: DeepPartial<DeactivateUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeactivateUserResponse>;
  /**
   * Reactivate a user
   *
   * Reactivate a previously deactivated user and change the state to 'active'.
   * The user will be able to log in again.
   *
   * The endpoint returns an error if the user is not in the state 'deactivated'.
   */
  reactivateUser(
    request: DeepPartial<ReactivateUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ReactivateUserResponse>;
  /**
   * Lock a user
   *
   * Lock an existing user and change the state 'locked'.
   * The user will not be able to log in anymore.
   * Use lock user when the user should temporarily not be able to log in
   * because of an event that happened (wrong password, etc.)
   *
   * The endpoint returns an error if the user is already in the state 'locked'.
   */
  lockUser(request: DeepPartial<LockUserRequest>, options?: CallOptions & CallOptionsExt): Promise<LockUserResponse>;
  /**
   * Unlock a user
   *
   * Unlock a previously locked user and change the state to 'active'.
   * The user will be able to log in again.
   *
   * The endpoint returns an error if the user is not in the state 'locked'.
   */
  unlockUser(
    request: DeepPartial<UnlockUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<UnlockUserResponse>;
  /**
   * Delete a user
   *
   * Delete an existing user and change the state to 'deleted'.
   * The user will be able to log in anymore.
   */
  deleteUser(
    request: DeepPartial<DeleteUserRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<DeleteUserResponse>;
  /**
   * Set contact email
   *
   * Add or update the contact email address of a user.
   * If the email is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by email.
   */
  setContactEmail(
    request: DeepPartial<SetContactEmailRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetContactEmailResponse>;
  /**
   * Verify the contact email
   *
   * Verify the contact email with the provided code.
   */
  verifyContactEmail(
    request: DeepPartial<VerifyContactEmailRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyContactEmailResponse>;
  /**
   * Resend the contact email code
   *
   * Resend the email with the verification code for the contact email address.
   */
  resendContactEmailCode(
    request: DeepPartial<ResendContactEmailCodeRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResendContactEmailCodeResponse>;
  /**
   * Set contact phone
   *
   * Add or update the contact phone number of a user.
   * If the phone is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by SMS.
   */
  setContactPhone(
    request: DeepPartial<SetContactPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetContactPhoneResponse>;
  /**
   * Verify the contact phone
   *
   * Verify the contact phone with the provided code.
   */
  verifyContactPhone(
    request: DeepPartial<VerifyContactPhoneRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyContactPhoneResponse>;
  /**
   * Resend the contact phone code
   *
   * Resend the phone with the verification code for the contact phone number.
   */
  resendContactPhoneCode(
    request: DeepPartial<ResendContactPhoneCodeRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<ResendContactPhoneCodeResponse>;
  /**
   * Add a username
   *
   * Add a new unique username to a user. The username will be used to identify the user on authentication.
   */
  addUsername(
    request: DeepPartial<AddUsernameRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddUsernameResponse>;
  /**
   * Remove a username
   *
   * Remove an existing username of a user, so it cannot be used for authentication anymore.
   */
  removeUsername(
    request: DeepPartial<RemoveUsernameRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveUsernameResponse>;
  /**
   * Set a password
   *
   * Add, update or reset a user's password with either a verification code or the current password.
   */
  setPassword(
    request: DeepPartial<SetPasswordRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<SetPasswordResponse>;
  /**
   * Request password reset
   *
   * Request a code to be able to set a new password.
   */
  requestPasswordReset(
    request: DeepPartial<RequestPasswordResetRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RequestPasswordResetResponse>;
  /**
   * Start a WebAuthN registration
   *
   * Start the registration of a new WebAuthN device (e.g. Passkeys) for a user.
   * As a response the public key credential creation options are returned,
   * which are used to verify the device.
   */
  startWebAuthNRegistration(
    request: DeepPartial<StartWebAuthNRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<StartWebAuthNRegistrationResponse>;
  /**
   * Verify a WebAuthN registration
   *
   * Verify the WebAuthN registration started by StartWebAuthNRegistration with the public key credential.
   */
  verifyWebAuthNRegistration(
    request: DeepPartial<VerifyWebAuthNRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyWebAuthNRegistrationResponse>;
  /**
   * Create a WebAuthN registration link
   *
   * Create a link, which includes a code, that can either be returned or directly sent to the user.
   * The code will allow the user to start a new WebAuthN registration.
   */
  createWebAuthNRegistrationLink(
    request: DeepPartial<CreateWebAuthNRegistrationLinkRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<CreateWebAuthNRegistrationLinkResponse>;
  /**
   * Remove a WebAuthN authenticator
   *
   * Remove an existing WebAuthN authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeWebAuthNAuthenticator(
    request: DeepPartial<RemoveWebAuthNAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveWebAuthNAuthenticatorResponse>;
  /**
   * Start a TOTP registration
   *
   * Start the registration of a new time-based one-time-password (TOTP) generator for a user.
   * As a response a secret is returned, which is used to initialize a TOTP app or device.
   */
  startTOTPRegistration(
    request: DeepPartial<StartTOTPRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<StartTOTPRegistrationResponse>;
  /**
   * Verify a TOTP registration
   *
   * Verify the time-based one-time-password (TOTP) registration with the generated code.
   */
  verifyTOTPRegistration(
    request: DeepPartial<VerifyTOTPRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyTOTPRegistrationResponse>;
  /**
   * Remove a TOTP authenticator
   *
   * Remove an existing time-based one-time-password (TOTP) authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeTOTPAuthenticator(
    request: DeepPartial<RemoveTOTPAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveTOTPAuthenticatorResponse>;
  /**
   * Add a OTP SMS authenticator
   *
   * Add a new one-time-password (OTP) SMS authenticator to a user.
   * If the phone is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by SMS.
   */
  addOTPSMSAuthenticator(
    request: DeepPartial<AddOTPSMSAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddOTPSMSAuthenticatorResponse>;
  /**
   * Verify OTP SMS registration
   *
   * Verify the OTP SMS registration with the provided code.
   */
  verifyOTPSMSRegistration(
    request: DeepPartial<VerifyOTPSMSRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyOTPSMSRegistrationResponse>;
  /**
   * Remove a OTP SMS authenticator
   *
   * Remove an existing one-time-password (OTP) SMS authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeOTPSMSAuthenticator(
    request: DeepPartial<RemoveOTPSMSAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveOTPSMSAuthenticatorResponse>;
  /**
   * Add a OTP Email authenticator
   *
   * Add a new one-time-password (OTP) Email authenticator to a user.
   * If the email is not passed as verified, a verification code will be generated,
   * which can be either returned or will be sent to the user by email.
   */
  addOTPEmailAuthenticator(
    request: DeepPartial<AddOTPEmailAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddOTPEmailAuthenticatorResponse>;
  /**
   * Verify OTP Email registration
   *
   * Verify the OTP Email registration with the provided code.
   */
  verifyOTPEmailRegistration(
    request: DeepPartial<VerifyOTPEmailRegistrationRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<VerifyOTPEmailRegistrationResponse>;
  /**
   * Remove a OTP Email authenticator
   *
   * Remove an existing one-time-password (OTP) Email authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeOTPEmailAuthenticator(
    request: DeepPartial<RemoveOTPEmailAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveOTPEmailAuthenticatorResponse>;
  /**
   * Start an IDP authentication intent
   *
   * Start a new authentication intent on configured identity provider (IDP) for external login, registration or linking.
   */
  startIdentityProviderIntent(
    request: DeepPartial<StartIdentityProviderIntentRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<StartIdentityProviderIntentResponse>;
  /**
   * Retrieve the information of the IDP authentication intent
   *
   * Retrieve the information returned by the identity provider (IDP) for registration or updating an existing user with new information.
   */
  retrieveIdentityProviderIntent(
    request: DeepPartial<RetrieveIdentityProviderIntentRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RetrieveIdentityProviderIntentResponse>;
  /**
   * Add an IDP authenticator to a user
   *
   * Add a new identity provider (IDP) authenticator to an existing user.
   * This will allow the user to authenticate with the provided IDP.
   */
  addIDPAuthenticator(
    request: DeepPartial<AddIDPAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<AddIDPAuthenticatorResponse>;
  /**
   * Remove an IDP authenticator
   *
   * Remove an existing identity provider (IDP) authenticator from a user, so it cannot be used for authentication anymore.
   */
  removeIDPAuthenticator(
    request: DeepPartial<RemoveIDPAuthenticatorRequest>,
    options?: CallOptions & CallOptionsExt,
  ): Promise<RemoveIDPAuthenticatorResponse>;
}

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

function isObject(value: any): boolean {
  return typeof value === "object" && value !== null;
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
