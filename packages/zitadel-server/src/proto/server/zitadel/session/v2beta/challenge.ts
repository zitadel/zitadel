/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Struct } from "../../../google/protobuf/struct";

export const protobufPackage = "zitadel.session.v2beta";

export enum UserVerificationRequirement {
  USER_VERIFICATION_REQUIREMENT_UNSPECIFIED = 0,
  USER_VERIFICATION_REQUIREMENT_REQUIRED = 1,
  USER_VERIFICATION_REQUIREMENT_PREFERRED = 2,
  USER_VERIFICATION_REQUIREMENT_DISCOURAGED = 3,
  UNRECOGNIZED = -1,
}

export function userVerificationRequirementFromJSON(object: any): UserVerificationRequirement {
  switch (object) {
    case 0:
    case "USER_VERIFICATION_REQUIREMENT_UNSPECIFIED":
      return UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_UNSPECIFIED;
    case 1:
    case "USER_VERIFICATION_REQUIREMENT_REQUIRED":
      return UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_REQUIRED;
    case 2:
    case "USER_VERIFICATION_REQUIREMENT_PREFERRED":
      return UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_PREFERRED;
    case 3:
    case "USER_VERIFICATION_REQUIREMENT_DISCOURAGED":
      return UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_DISCOURAGED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return UserVerificationRequirement.UNRECOGNIZED;
  }
}

export function userVerificationRequirementToJSON(object: UserVerificationRequirement): string {
  switch (object) {
    case UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_UNSPECIFIED:
      return "USER_VERIFICATION_REQUIREMENT_UNSPECIFIED";
    case UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_REQUIRED:
      return "USER_VERIFICATION_REQUIREMENT_REQUIRED";
    case UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_PREFERRED:
      return "USER_VERIFICATION_REQUIREMENT_PREFERRED";
    case UserVerificationRequirement.USER_VERIFICATION_REQUIREMENT_DISCOURAGED:
      return "USER_VERIFICATION_REQUIREMENT_DISCOURAGED";
    case UserVerificationRequirement.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface RequestChallenges {
  webAuthN?: RequestChallenges_WebAuthN | undefined;
  otpSms?: RequestChallenges_OTPSMS | undefined;
  otpEmail?: RequestChallenges_OTPEmail | undefined;
}

export interface RequestChallenges_WebAuthN {
  domain: string;
  userVerificationRequirement: UserVerificationRequirement;
}

export interface RequestChallenges_OTPSMS {
  returnCode: boolean;
}

export interface RequestChallenges_OTPEmail {
  sendCode?: RequestChallenges_OTPEmail_SendCode | undefined;
  returnCode?: RequestChallenges_OTPEmail_ReturnCode | undefined;
}

export interface RequestChallenges_OTPEmail_SendCode {
  urlTemplate?: string | undefined;
}

export interface RequestChallenges_OTPEmail_ReturnCode {
}

export interface Challenges {
  webAuthN?: Challenges_WebAuthN | undefined;
  otpSms?: string | undefined;
  otpEmail?: string | undefined;
}

export interface Challenges_WebAuthN {
  publicKeyCredentialRequestOptions: { [key: string]: any } | undefined;
}

function createBaseRequestChallenges(): RequestChallenges {
  return { webAuthN: undefined, otpSms: undefined, otpEmail: undefined };
}

export const RequestChallenges = {
  encode(message: RequestChallenges, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.webAuthN !== undefined) {
      RequestChallenges_WebAuthN.encode(message.webAuthN, writer.uint32(10).fork()).ldelim();
    }
    if (message.otpSms !== undefined) {
      RequestChallenges_OTPSMS.encode(message.otpSms, writer.uint32(18).fork()).ldelim();
    }
    if (message.otpEmail !== undefined) {
      RequestChallenges_OTPEmail.encode(message.otpEmail, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.webAuthN = RequestChallenges_WebAuthN.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpSms = RequestChallenges_OTPSMS.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.otpEmail = RequestChallenges_OTPEmail.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestChallenges {
    return {
      webAuthN: isSet(object.webAuthN) ? RequestChallenges_WebAuthN.fromJSON(object.webAuthN) : undefined,
      otpSms: isSet(object.otpSms) ? RequestChallenges_OTPSMS.fromJSON(object.otpSms) : undefined,
      otpEmail: isSet(object.otpEmail) ? RequestChallenges_OTPEmail.fromJSON(object.otpEmail) : undefined,
    };
  },

  toJSON(message: RequestChallenges): unknown {
    const obj: any = {};
    message.webAuthN !== undefined &&
      (obj.webAuthN = message.webAuthN ? RequestChallenges_WebAuthN.toJSON(message.webAuthN) : undefined);
    message.otpSms !== undefined &&
      (obj.otpSms = message.otpSms ? RequestChallenges_OTPSMS.toJSON(message.otpSms) : undefined);
    message.otpEmail !== undefined &&
      (obj.otpEmail = message.otpEmail ? RequestChallenges_OTPEmail.toJSON(message.otpEmail) : undefined);
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges>): RequestChallenges {
    return RequestChallenges.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestChallenges>): RequestChallenges {
    const message = createBaseRequestChallenges();
    message.webAuthN = (object.webAuthN !== undefined && object.webAuthN !== null)
      ? RequestChallenges_WebAuthN.fromPartial(object.webAuthN)
      : undefined;
    message.otpSms = (object.otpSms !== undefined && object.otpSms !== null)
      ? RequestChallenges_OTPSMS.fromPartial(object.otpSms)
      : undefined;
    message.otpEmail = (object.otpEmail !== undefined && object.otpEmail !== null)
      ? RequestChallenges_OTPEmail.fromPartial(object.otpEmail)
      : undefined;
    return message;
  },
};

function createBaseRequestChallenges_WebAuthN(): RequestChallenges_WebAuthN {
  return { domain: "", userVerificationRequirement: 0 };
}

export const RequestChallenges_WebAuthN = {
  encode(message: RequestChallenges_WebAuthN, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.domain !== "") {
      writer.uint32(10).string(message.domain);
    }
    if (message.userVerificationRequirement !== 0) {
      writer.uint32(16).int32(message.userVerificationRequirement);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges_WebAuthN {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges_WebAuthN();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.domain = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.userVerificationRequirement = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestChallenges_WebAuthN {
    return {
      domain: isSet(object.domain) ? String(object.domain) : "",
      userVerificationRequirement: isSet(object.userVerificationRequirement)
        ? userVerificationRequirementFromJSON(object.userVerificationRequirement)
        : 0,
    };
  },

  toJSON(message: RequestChallenges_WebAuthN): unknown {
    const obj: any = {};
    message.domain !== undefined && (obj.domain = message.domain);
    message.userVerificationRequirement !== undefined &&
      (obj.userVerificationRequirement = userVerificationRequirementToJSON(message.userVerificationRequirement));
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges_WebAuthN>): RequestChallenges_WebAuthN {
    return RequestChallenges_WebAuthN.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestChallenges_WebAuthN>): RequestChallenges_WebAuthN {
    const message = createBaseRequestChallenges_WebAuthN();
    message.domain = object.domain ?? "";
    message.userVerificationRequirement = object.userVerificationRequirement ?? 0;
    return message;
  },
};

function createBaseRequestChallenges_OTPSMS(): RequestChallenges_OTPSMS {
  return { returnCode: false };
}

export const RequestChallenges_OTPSMS = {
  encode(message: RequestChallenges_OTPSMS, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.returnCode === true) {
      writer.uint32(8).bool(message.returnCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges_OTPSMS {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges_OTPSMS();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.returnCode = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestChallenges_OTPSMS {
    return { returnCode: isSet(object.returnCode) ? Boolean(object.returnCode) : false };
  },

  toJSON(message: RequestChallenges_OTPSMS): unknown {
    const obj: any = {};
    message.returnCode !== undefined && (obj.returnCode = message.returnCode);
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges_OTPSMS>): RequestChallenges_OTPSMS {
    return RequestChallenges_OTPSMS.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestChallenges_OTPSMS>): RequestChallenges_OTPSMS {
    const message = createBaseRequestChallenges_OTPSMS();
    message.returnCode = object.returnCode ?? false;
    return message;
  },
};

function createBaseRequestChallenges_OTPEmail(): RequestChallenges_OTPEmail {
  return { sendCode: undefined, returnCode: undefined };
}

export const RequestChallenges_OTPEmail = {
  encode(message: RequestChallenges_OTPEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sendCode !== undefined) {
      RequestChallenges_OTPEmail_SendCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      RequestChallenges_OTPEmail_ReturnCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges_OTPEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges_OTPEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          if (tag != 18) {
            break;
          }

          message.sendCode = RequestChallenges_OTPEmail_SendCode.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.returnCode = RequestChallenges_OTPEmail_ReturnCode.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequestChallenges_OTPEmail {
    return {
      sendCode: isSet(object.sendCode) ? RequestChallenges_OTPEmail_SendCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode)
        ? RequestChallenges_OTPEmail_ReturnCode.fromJSON(object.returnCode)
        : undefined,
    };
  },

  toJSON(message: RequestChallenges_OTPEmail): unknown {
    const obj: any = {};
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? RequestChallenges_OTPEmail_SendCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined && (obj.returnCode = message.returnCode
      ? RequestChallenges_OTPEmail_ReturnCode.toJSON(message.returnCode)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges_OTPEmail>): RequestChallenges_OTPEmail {
    return RequestChallenges_OTPEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestChallenges_OTPEmail>): RequestChallenges_OTPEmail {
    const message = createBaseRequestChallenges_OTPEmail();
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? RequestChallenges_OTPEmail_SendCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? RequestChallenges_OTPEmail_ReturnCode.fromPartial(object.returnCode)
      : undefined;
    return message;
  },
};

function createBaseRequestChallenges_OTPEmail_SendCode(): RequestChallenges_OTPEmail_SendCode {
  return { urlTemplate: undefined };
}

export const RequestChallenges_OTPEmail_SendCode = {
  encode(message: RequestChallenges_OTPEmail_SendCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.urlTemplate !== undefined) {
      writer.uint32(10).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges_OTPEmail_SendCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges_OTPEmail_SendCode();
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

  fromJSON(object: any): RequestChallenges_OTPEmail_SendCode {
    return { urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined };
  },

  toJSON(message: RequestChallenges_OTPEmail_SendCode): unknown {
    const obj: any = {};
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges_OTPEmail_SendCode>): RequestChallenges_OTPEmail_SendCode {
    return RequestChallenges_OTPEmail_SendCode.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RequestChallenges_OTPEmail_SendCode>): RequestChallenges_OTPEmail_SendCode {
    const message = createBaseRequestChallenges_OTPEmail_SendCode();
    message.urlTemplate = object.urlTemplate ?? undefined;
    return message;
  },
};

function createBaseRequestChallenges_OTPEmail_ReturnCode(): RequestChallenges_OTPEmail_ReturnCode {
  return {};
}

export const RequestChallenges_OTPEmail_ReturnCode = {
  encode(_: RequestChallenges_OTPEmail_ReturnCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequestChallenges_OTPEmail_ReturnCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequestChallenges_OTPEmail_ReturnCode();
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

  fromJSON(_: any): RequestChallenges_OTPEmail_ReturnCode {
    return {};
  },

  toJSON(_: RequestChallenges_OTPEmail_ReturnCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<RequestChallenges_OTPEmail_ReturnCode>): RequestChallenges_OTPEmail_ReturnCode {
    return RequestChallenges_OTPEmail_ReturnCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<RequestChallenges_OTPEmail_ReturnCode>): RequestChallenges_OTPEmail_ReturnCode {
    const message = createBaseRequestChallenges_OTPEmail_ReturnCode();
    return message;
  },
};

function createBaseChallenges(): Challenges {
  return { webAuthN: undefined, otpSms: undefined, otpEmail: undefined };
}

export const Challenges = {
  encode(message: Challenges, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.webAuthN !== undefined) {
      Challenges_WebAuthN.encode(message.webAuthN, writer.uint32(10).fork()).ldelim();
    }
    if (message.otpSms !== undefined) {
      writer.uint32(18).string(message.otpSms);
    }
    if (message.otpEmail !== undefined) {
      writer.uint32(26).string(message.otpEmail);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Challenges {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChallenges();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.webAuthN = Challenges_WebAuthN.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.otpSms = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.otpEmail = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Challenges {
    return {
      webAuthN: isSet(object.webAuthN) ? Challenges_WebAuthN.fromJSON(object.webAuthN) : undefined,
      otpSms: isSet(object.otpSms) ? String(object.otpSms) : undefined,
      otpEmail: isSet(object.otpEmail) ? String(object.otpEmail) : undefined,
    };
  },

  toJSON(message: Challenges): unknown {
    const obj: any = {};
    message.webAuthN !== undefined &&
      (obj.webAuthN = message.webAuthN ? Challenges_WebAuthN.toJSON(message.webAuthN) : undefined);
    message.otpSms !== undefined && (obj.otpSms = message.otpSms);
    message.otpEmail !== undefined && (obj.otpEmail = message.otpEmail);
    return obj;
  },

  create(base?: DeepPartial<Challenges>): Challenges {
    return Challenges.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Challenges>): Challenges {
    const message = createBaseChallenges();
    message.webAuthN = (object.webAuthN !== undefined && object.webAuthN !== null)
      ? Challenges_WebAuthN.fromPartial(object.webAuthN)
      : undefined;
    message.otpSms = object.otpSms ?? undefined;
    message.otpEmail = object.otpEmail ?? undefined;
    return message;
  },
};

function createBaseChallenges_WebAuthN(): Challenges_WebAuthN {
  return { publicKeyCredentialRequestOptions: undefined };
}

export const Challenges_WebAuthN = {
  encode(message: Challenges_WebAuthN, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.publicKeyCredentialRequestOptions !== undefined) {
      Struct.encode(Struct.wrap(message.publicKeyCredentialRequestOptions), writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Challenges_WebAuthN {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseChallenges_WebAuthN();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.publicKeyCredentialRequestOptions = Struct.unwrap(Struct.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Challenges_WebAuthN {
    return {
      publicKeyCredentialRequestOptions: isObject(object.publicKeyCredentialRequestOptions)
        ? object.publicKeyCredentialRequestOptions
        : undefined,
    };
  },

  toJSON(message: Challenges_WebAuthN): unknown {
    const obj: any = {};
    message.publicKeyCredentialRequestOptions !== undefined &&
      (obj.publicKeyCredentialRequestOptions = message.publicKeyCredentialRequestOptions);
    return obj;
  },

  create(base?: DeepPartial<Challenges_WebAuthN>): Challenges_WebAuthN {
    return Challenges_WebAuthN.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Challenges_WebAuthN>): Challenges_WebAuthN {
    const message = createBaseChallenges_WebAuthN();
    message.publicKeyCredentialRequestOptions = object.publicKeyCredentialRequestOptions ?? undefined;
    return message;
  },
};

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
