/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "zitadel.user.v2beta";

export enum NotificationType {
  NOTIFICATION_TYPE_Unspecified = 0,
  NOTIFICATION_TYPE_Email = 1,
  NOTIFICATION_TYPE_SMS = 2,
  UNRECOGNIZED = -1,
}

export function notificationTypeFromJSON(object: any): NotificationType {
  switch (object) {
    case 0:
    case "NOTIFICATION_TYPE_Unspecified":
      return NotificationType.NOTIFICATION_TYPE_Unspecified;
    case 1:
    case "NOTIFICATION_TYPE_Email":
      return NotificationType.NOTIFICATION_TYPE_Email;
    case 2:
    case "NOTIFICATION_TYPE_SMS":
      return NotificationType.NOTIFICATION_TYPE_SMS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return NotificationType.UNRECOGNIZED;
  }
}

export function notificationTypeToJSON(object: NotificationType): string {
  switch (object) {
    case NotificationType.NOTIFICATION_TYPE_Unspecified:
      return "NOTIFICATION_TYPE_Unspecified";
    case NotificationType.NOTIFICATION_TYPE_Email:
      return "NOTIFICATION_TYPE_Email";
    case NotificationType.NOTIFICATION_TYPE_SMS:
      return "NOTIFICATION_TYPE_SMS";
    case NotificationType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Password {
  password: string;
  changeRequired: boolean;
}

export interface HashedPassword {
  hash: string;
  changeRequired: boolean;
}

export interface SendPasswordResetLink {
  notificationType: NotificationType;
  urlTemplate?: string | undefined;
}

export interface ReturnPasswordResetCode {
}

export interface SetPassword {
  password?: Password | undefined;
  hashedPassword?: HashedPassword | undefined;
  currentPassword?: string | undefined;
  verificationCode?: string | undefined;
}

function createBasePassword(): Password {
  return { password: "", changeRequired: false };
}

export const Password = {
  encode(message: Password, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.password !== "") {
      writer.uint32(10).string(message.password);
    }
    if (message.changeRequired === true) {
      writer.uint32(16).bool(message.changeRequired);
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

          message.password = reader.string();
          continue;
        case 2:
          if (tag != 16) {
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

  fromJSON(object: any): Password {
    return {
      password: isSet(object.password) ? String(object.password) : "",
      changeRequired: isSet(object.changeRequired) ? Boolean(object.changeRequired) : false,
    };
  },

  toJSON(message: Password): unknown {
    const obj: any = {};
    message.password !== undefined && (obj.password = message.password);
    message.changeRequired !== undefined && (obj.changeRequired = message.changeRequired);
    return obj;
  },

  create(base?: DeepPartial<Password>): Password {
    return Password.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Password>): Password {
    const message = createBasePassword();
    message.password = object.password ?? "";
    message.changeRequired = object.changeRequired ?? false;
    return message;
  },
};

function createBaseHashedPassword(): HashedPassword {
  return { hash: "", changeRequired: false };
}

export const HashedPassword = {
  encode(message: HashedPassword, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.hash !== "") {
      writer.uint32(10).string(message.hash);
    }
    if (message.changeRequired === true) {
      writer.uint32(16).bool(message.changeRequired);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HashedPassword {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHashedPassword();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.hash = reader.string();
          continue;
        case 2:
          if (tag != 16) {
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

  fromJSON(object: any): HashedPassword {
    return {
      hash: isSet(object.hash) ? String(object.hash) : "",
      changeRequired: isSet(object.changeRequired) ? Boolean(object.changeRequired) : false,
    };
  },

  toJSON(message: HashedPassword): unknown {
    const obj: any = {};
    message.hash !== undefined && (obj.hash = message.hash);
    message.changeRequired !== undefined && (obj.changeRequired = message.changeRequired);
    return obj;
  },

  create(base?: DeepPartial<HashedPassword>): HashedPassword {
    return HashedPassword.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<HashedPassword>): HashedPassword {
    const message = createBaseHashedPassword();
    message.hash = object.hash ?? "";
    message.changeRequired = object.changeRequired ?? false;
    return message;
  },
};

function createBaseSendPasswordResetLink(): SendPasswordResetLink {
  return { notificationType: 0, urlTemplate: undefined };
}

export const SendPasswordResetLink = {
  encode(message: SendPasswordResetLink, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.notificationType !== 0) {
      writer.uint32(8).int32(message.notificationType);
    }
    if (message.urlTemplate !== undefined) {
      writer.uint32(18).string(message.urlTemplate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SendPasswordResetLink {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSendPasswordResetLink();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.notificationType = reader.int32() as any;
          continue;
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

  fromJSON(object: any): SendPasswordResetLink {
    return {
      notificationType: isSet(object.notificationType) ? notificationTypeFromJSON(object.notificationType) : 0,
      urlTemplate: isSet(object.urlTemplate) ? String(object.urlTemplate) : undefined,
    };
  },

  toJSON(message: SendPasswordResetLink): unknown {
    const obj: any = {};
    message.notificationType !== undefined && (obj.notificationType = notificationTypeToJSON(message.notificationType));
    message.urlTemplate !== undefined && (obj.urlTemplate = message.urlTemplate);
    return obj;
  },

  create(base?: DeepPartial<SendPasswordResetLink>): SendPasswordResetLink {
    return SendPasswordResetLink.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SendPasswordResetLink>): SendPasswordResetLink {
    const message = createBaseSendPasswordResetLink();
    message.notificationType = object.notificationType ?? 0;
    message.urlTemplate = object.urlTemplate ?? undefined;
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

function createBaseSetPassword(): SetPassword {
  return { password: undefined, hashedPassword: undefined, currentPassword: undefined, verificationCode: undefined };
}

export const SetPassword = {
  encode(message: SetPassword, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.password !== undefined) {
      Password.encode(message.password, writer.uint32(10).fork()).ldelim();
    }
    if (message.hashedPassword !== undefined) {
      HashedPassword.encode(message.hashedPassword, writer.uint32(18).fork()).ldelim();
    }
    if (message.currentPassword !== undefined) {
      writer.uint32(26).string(message.currentPassword);
    }
    if (message.verificationCode !== undefined) {
      writer.uint32(34).string(message.verificationCode);
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

          message.password = Password.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.hashedPassword = HashedPassword.decode(reader, reader.uint32());
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

  fromJSON(object: any): SetPassword {
    return {
      password: isSet(object.password) ? Password.fromJSON(object.password) : undefined,
      hashedPassword: isSet(object.hashedPassword) ? HashedPassword.fromJSON(object.hashedPassword) : undefined,
      currentPassword: isSet(object.currentPassword) ? String(object.currentPassword) : undefined,
      verificationCode: isSet(object.verificationCode) ? String(object.verificationCode) : undefined,
    };
  },

  toJSON(message: SetPassword): unknown {
    const obj: any = {};
    message.password !== undefined && (obj.password = message.password ? Password.toJSON(message.password) : undefined);
    message.hashedPassword !== undefined &&
      (obj.hashedPassword = message.hashedPassword ? HashedPassword.toJSON(message.hashedPassword) : undefined);
    message.currentPassword !== undefined && (obj.currentPassword = message.currentPassword);
    message.verificationCode !== undefined && (obj.verificationCode = message.verificationCode);
    return obj;
  },

  create(base?: DeepPartial<SetPassword>): SetPassword {
    return SetPassword.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetPassword>): SetPassword {
    const message = createBaseSetPassword();
    message.password = (object.password !== undefined && object.password !== null)
      ? Password.fromPartial(object.password)
      : undefined;
    message.hashedPassword = (object.hashedPassword !== undefined && object.hashedPassword !== null)
      ? HashedPassword.fromPartial(object.hashedPassword)
      : undefined;
    message.currentPassword = object.currentPassword ?? undefined;
    message.verificationCode = object.verificationCode ?? undefined;
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
