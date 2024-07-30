/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.user.v3alpha";

export interface Contact {
  /** Email contact information of the user. */
  email:
    | Email
    | undefined;
  /** Phone contact information of the user. */
  phone: Phone | undefined;
}

export interface Email {
  /** Email address of the user. */
  address: string;
  /** IsVerified states if the email address has been verified to belong to the user. */
  isVerified: boolean;
}

export interface Phone {
  /** Phone number of the user. */
  number: string;
  /** IsVerified states if the phone number has been verified to belong to the user. */
  isVerified: boolean;
}

export interface SetContact {
  email?: SetEmail | undefined;
  phone?: SetPhone | undefined;
}

export interface SetEmail {
  /** Set the email address. */
  address: string;
  /** Let ZITADEL send the link to the user via email. */
  sendCode?:
    | SendEmailVerificationCode
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?:
    | ReturnEmailVerificationCode
    | undefined;
  /** Set the email as already verified. */
  isVerified?: boolean | undefined;
}

export interface SendEmailVerificationCode {
  /**
   * Optionally set a url_template, which will be used in the verification mail sent by ZITADEL
   * to guide the user to your verification page.
   * If no template is set, the default ZITADEL url will be used.
   */
  urlTemplate?: string | undefined;
}

export interface ReturnEmailVerificationCode {
}

export interface SetPhone {
  /** Set the user's phone number. */
  number: string;
  /** Let ZITADEL send the link to the user via SMS. */
  sendCode?:
    | SendPhoneVerificationCode
    | undefined;
  /** Get the code back to provide it to the user in your preferred mechanism. */
  returnCode?:
    | ReturnPhoneVerificationCode
    | undefined;
  /** Set the phone as already verified. */
  isVerified?: boolean | undefined;
}

export interface SendPhoneVerificationCode {
}

export interface ReturnPhoneVerificationCode {
}

function createBaseContact(): Contact {
  return { email: undefined, phone: undefined };
}

export const Contact = {
  encode(message: Contact, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== undefined) {
      Email.encode(message.email, writer.uint32(10).fork()).ldelim();
    }
    if (message.phone !== undefined) {
      Phone.encode(message.phone, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Contact {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContact();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.email = Email.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.phone = Phone.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Contact {
    return {
      email: isSet(object.email) ? Email.fromJSON(object.email) : undefined,
      phone: isSet(object.phone) ? Phone.fromJSON(object.phone) : undefined,
    };
  },

  toJSON(message: Contact): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email ? Email.toJSON(message.email) : undefined);
    message.phone !== undefined && (obj.phone = message.phone ? Phone.toJSON(message.phone) : undefined);
    return obj;
  },

  create(base?: DeepPartial<Contact>): Contact {
    return Contact.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Contact>): Contact {
    const message = createBaseContact();
    message.email = (object.email !== undefined && object.email !== null) ? Email.fromPartial(object.email) : undefined;
    message.phone = (object.phone !== undefined && object.phone !== null) ? Phone.fromPartial(object.phone) : undefined;
    return message;
  },
};

function createBaseEmail(): Email {
  return { address: "", isVerified: false };
}

export const Email = {
  encode(message: Email, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.isVerified === true) {
      writer.uint32(16).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Email {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.address = reader.string();
          continue;
        case 2:
          if (tag != 16) {
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

  fromJSON(object: any): Email {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: Email): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<Email>): Email {
    return Email.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Email>): Email {
    const message = createBaseEmail();
    message.address = object.address ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBasePhone(): Phone {
  return { number: "", isVerified: false };
}

export const Phone = {
  encode(message: Phone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.number !== "") {
      writer.uint32(10).string(message.number);
    }
    if (message.isVerified === true) {
      writer.uint32(16).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Phone {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePhone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.number = reader.string();
          continue;
        case 2:
          if (tag != 16) {
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

  fromJSON(object: any): Phone {
    return {
      number: isSet(object.number) ? String(object.number) : "",
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : false,
    };
  },

  toJSON(message: Phone): unknown {
    const obj: any = {};
    message.number !== undefined && (obj.number = message.number);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<Phone>): Phone {
    return Phone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Phone>): Phone {
    const message = createBasePhone();
    message.number = object.number ?? "";
    message.isVerified = object.isVerified ?? false;
    return message;
  },
};

function createBaseSetContact(): SetContact {
  return { email: undefined, phone: undefined };
}

export const SetContact = {
  encode(message: SetContact, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== undefined) {
      SetEmail.encode(message.email, writer.uint32(10).fork()).ldelim();
    }
    if (message.phone !== undefined) {
      SetPhone.encode(message.phone, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetContact {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetContact();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.email = SetEmail.decode(reader, reader.uint32());
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

  fromJSON(object: any): SetContact {
    return {
      email: isSet(object.email) ? SetEmail.fromJSON(object.email) : undefined,
      phone: isSet(object.phone) ? SetPhone.fromJSON(object.phone) : undefined,
    };
  },

  toJSON(message: SetContact): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email ? SetEmail.toJSON(message.email) : undefined);
    message.phone !== undefined && (obj.phone = message.phone ? SetPhone.toJSON(message.phone) : undefined);
    return obj;
  },

  create(base?: DeepPartial<SetContact>): SetContact {
    return SetContact.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetContact>): SetContact {
    const message = createBaseSetContact();
    message.email = (object.email !== undefined && object.email !== null)
      ? SetEmail.fromPartial(object.email)
      : undefined;
    message.phone = (object.phone !== undefined && object.phone !== null)
      ? SetPhone.fromPartial(object.phone)
      : undefined;
    return message;
  },
};

function createBaseSetEmail(): SetEmail {
  return { address: "", sendCode: undefined, returnCode: undefined, isVerified: undefined };
}

export const SetEmail = {
  encode(message: SetEmail, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    if (message.sendCode !== undefined) {
      SendEmailVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnEmailVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    if (message.isVerified !== undefined) {
      writer.uint32(32).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetEmail {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetEmail();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.address = reader.string();
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
        case 4:
          if (tag != 32) {
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

  fromJSON(object: any): SetEmail {
    return {
      address: isSet(object.address) ? String(object.address) : "",
      sendCode: isSet(object.sendCode) ? SendEmailVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnEmailVerificationCode.fromJSON(object.returnCode) : undefined,
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : undefined,
    };
  },

  toJSON(message: SetEmail): unknown {
    const obj: any = {};
    message.address !== undefined && (obj.address = message.address);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendEmailVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnEmailVerificationCode.toJSON(message.returnCode) : undefined);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<SetEmail>): SetEmail {
    return SetEmail.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetEmail>): SetEmail {
    const message = createBaseSetEmail();
    message.address = object.address ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendEmailVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnEmailVerificationCode.fromPartial(object.returnCode)
      : undefined;
    message.isVerified = object.isVerified ?? undefined;
    return message;
  },
};

function createBaseSendEmailVerificationCode(): SendEmailVerificationCode {
  return { urlTemplate: undefined };
}

export const SendEmailVerificationCode = {
  encode(message: SendEmailVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.urlTemplate !== undefined) {
      writer.uint32(10).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendEmailVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendEmailVerificationCode();
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

  fromJSON(object: any): SendEmailVerificationCode {
    return { urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined };
  },

  toJSON(message: SendEmailVerificationCode): unknown {
    const obj: any = {};
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<SendEmailVerificationCode>): SendEmailVerificationCode {
    return SendEmailVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendEmailVerificationCode>): SendEmailVerificationCode {
    const message = createBaseSendEmailVerificationCode();
    message.urlTemplate = object.urlTemplate ?? undefined;
    return message;
  },
};

function createBaseReturnEmailVerificationCode(): ReturnEmailVerificationCode {
  return {};
}

export const ReturnEmailVerificationCode = {
  encode(_: ReturnEmailVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnEmailVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnEmailVerificationCode();
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

  fromJSON(_: any): ReturnEmailVerificationCode {
    return {};
  },

  toJSON(_: ReturnEmailVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnEmailVerificationCode>): ReturnEmailVerificationCode {
    return ReturnEmailVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnEmailVerificationCode>): ReturnEmailVerificationCode {
    const message = createBaseReturnEmailVerificationCode();
    return message;
  },
};

function createBaseSetPhone(): SetPhone {
  return { number: "", sendCode: undefined, returnCode: undefined, isVerified: undefined };
}

export const SetPhone = {
  encode(message: SetPhone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.number !== "") {
      writer.uint32(10).string(message.number);
    }
    if (message.sendCode !== undefined) {
      SendPhoneVerificationCode.encode(message.sendCode, writer.uint32(18).fork()).ldelim();
    }
    if (message.returnCode !== undefined) {
      ReturnPhoneVerificationCode.encode(message.returnCode, writer.uint32(26).fork()).ldelim();
    }
    if (message.isVerified !== undefined) {
      writer.uint32(32).bool(message.isVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetPhone {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetPhone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.number = reader.string();
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
        case 4:
          if (tag != 32) {
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

  fromJSON(object: any): SetPhone {
    return {
      number: isSet(object.number) ? String(object.number) : "",
      sendCode: isSet(object.sendCode) ? SendPhoneVerificationCode.fromJSON(object.sendCode) : undefined,
      returnCode: isSet(object.returnCode) ? ReturnPhoneVerificationCode.fromJSON(object.returnCode) : undefined,
      isVerified: isSet(object.isVerified) ? Boolean(object.isVerified) : undefined,
    };
  },

  toJSON(message: SetPhone): unknown {
    const obj: any = {};
    message.number !== undefined && (obj.number = message.number);
    message.sendCode !== undefined &&
      (obj.sendCode = message.sendCode ? SendPhoneVerificationCode.toJSON(message.sendCode) : undefined);
    message.returnCode !== undefined &&
      (obj.returnCode = message.returnCode ? ReturnPhoneVerificationCode.toJSON(message.returnCode) : undefined);
    message.isVerified !== undefined && (obj.isVerified = message.isVerified);
    return obj;
  },

  create(base?: DeepPartial<SetPhone>): SetPhone {
    return SetPhone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetPhone>): SetPhone {
    const message = createBaseSetPhone();
    message.number = object.number ?? "";
    message.sendCode = (object.sendCode !== undefined && object.sendCode !== null)
      ? SendPhoneVerificationCode.fromPartial(object.sendCode)
      : undefined;
    message.returnCode = (object.returnCode !== undefined && object.returnCode !== null)
      ? ReturnPhoneVerificationCode.fromPartial(object.returnCode)
      : undefined;
    message.isVerified = object.isVerified ?? undefined;
    return message;
  },
};

function createBaseSendPhoneVerificationCode(): SendPhoneVerificationCode {
  return {};
}

export const SendPhoneVerificationCode = {
  encode(_: SendPhoneVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendPhoneVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendPhoneVerificationCode();
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

  fromJSON(_: any): SendPhoneVerificationCode {
    return {};
  },

  toJSON(_: SendPhoneVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<SendPhoneVerificationCode>): SendPhoneVerificationCode {
    return SendPhoneVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<SendPhoneVerificationCode>): SendPhoneVerificationCode {
    const message = createBaseSendPhoneVerificationCode();
    return message;
  },
};

function createBaseReturnPhoneVerificationCode(): ReturnPhoneVerificationCode {
  return {};
}

export const ReturnPhoneVerificationCode = {
  encode(_: ReturnPhoneVerificationCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReturnPhoneVerificationCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReturnPhoneVerificationCode();
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

  fromJSON(_: any): ReturnPhoneVerificationCode {
    return {};
  },

  toJSON(_: ReturnPhoneVerificationCode): unknown {
    const obj: any = {};
    return obj;
  },

  create(base?: DeepPartial<ReturnPhoneVerificationCode>): ReturnPhoneVerificationCode {
    return ReturnPhoneVerificationCode.fromPartial(base ?? {});
  },

  fromPartial(_: DeepPartial<ReturnPhoneVerificationCode>): ReturnPhoneVerificationCode {
    const message = createBaseReturnPhoneVerificationCode();
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
