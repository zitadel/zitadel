import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class SecretGenerator extends jspb.Message {
  getGeneratorType(): SecretGeneratorType;
  setGeneratorType(value: SecretGeneratorType): SecretGenerator;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SecretGenerator;
  hasDetails(): boolean;
  clearDetails(): SecretGenerator;

  getLength(): number;
  setLength(value: number): SecretGenerator;

  getExpiry(): google_protobuf_duration_pb.Duration | undefined;
  setExpiry(value?: google_protobuf_duration_pb.Duration): SecretGenerator;
  hasExpiry(): boolean;
  clearExpiry(): SecretGenerator;

  getIncludeLowerLetters(): boolean;
  setIncludeLowerLetters(value: boolean): SecretGenerator;

  getIncludeUpperLetters(): boolean;
  setIncludeUpperLetters(value: boolean): SecretGenerator;

  getIncludeDigits(): boolean;
  setIncludeDigits(value: boolean): SecretGenerator;

  getIncludeSymbols(): boolean;
  setIncludeSymbols(value: boolean): SecretGenerator;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecretGenerator.AsObject;
  static toObject(includeInstance: boolean, msg: SecretGenerator): SecretGenerator.AsObject;
  static serializeBinaryToWriter(message: SecretGenerator, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecretGenerator;
  static deserializeBinaryFromReader(message: SecretGenerator, reader: jspb.BinaryReader): SecretGenerator;
}

export namespace SecretGenerator {
  export type AsObject = {
    generatorType: SecretGeneratorType,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    length: number,
    expiry?: google_protobuf_duration_pb.Duration.AsObject,
    includeLowerLetters: boolean,
    includeUpperLetters: boolean,
    includeDigits: boolean,
    includeSymbols: boolean,
  }
}

export class SecretGeneratorQuery extends jspb.Message {
  getTypeQuery(): SecretGeneratorTypeQuery | undefined;
  setTypeQuery(value?: SecretGeneratorTypeQuery): SecretGeneratorQuery;
  hasTypeQuery(): boolean;
  clearTypeQuery(): SecretGeneratorQuery;

  getQueryCase(): SecretGeneratorQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecretGeneratorQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SecretGeneratorQuery): SecretGeneratorQuery.AsObject;
  static serializeBinaryToWriter(message: SecretGeneratorQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecretGeneratorQuery;
  static deserializeBinaryFromReader(message: SecretGeneratorQuery, reader: jspb.BinaryReader): SecretGeneratorQuery;
}

export namespace SecretGeneratorQuery {
  export type AsObject = {
    typeQuery?: SecretGeneratorTypeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    TYPE_QUERY = 1,
  }
}

export class SecretGeneratorTypeQuery extends jspb.Message {
  getGeneratorType(): SecretGeneratorType;
  setGeneratorType(value: SecretGeneratorType): SecretGeneratorTypeQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecretGeneratorTypeQuery.AsObject;
  static toObject(includeInstance: boolean, msg: SecretGeneratorTypeQuery): SecretGeneratorTypeQuery.AsObject;
  static serializeBinaryToWriter(message: SecretGeneratorTypeQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecretGeneratorTypeQuery;
  static deserializeBinaryFromReader(message: SecretGeneratorTypeQuery, reader: jspb.BinaryReader): SecretGeneratorTypeQuery;
}

export namespace SecretGeneratorTypeQuery {
  export type AsObject = {
    generatorType: SecretGeneratorType,
  }
}

export class SMTPConfig extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SMTPConfig;
  hasDetails(): boolean;
  clearDetails(): SMTPConfig;

  getSenderAddress(): string;
  setSenderAddress(value: string): SMTPConfig;

  getSenderName(): string;
  setSenderName(value: string): SMTPConfig;

  getTls(): boolean;
  setTls(value: boolean): SMTPConfig;

  getHost(): string;
  setHost(value: string): SMTPConfig;

  getUser(): string;
  setUser(value: string): SMTPConfig;

  getReplyToAddress(): string;
  setReplyToAddress(value: string): SMTPConfig;

  getState(): SMTPConfigState;
  setState(value: SMTPConfigState): SMTPConfig;

  getDescription(): string;
  setDescription(value: string): SMTPConfig;

  getId(): string;
  setId(value: string): SMTPConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SMTPConfig.AsObject;
  static toObject(includeInstance: boolean, msg: SMTPConfig): SMTPConfig.AsObject;
  static serializeBinaryToWriter(message: SMTPConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SMTPConfig;
  static deserializeBinaryFromReader(message: SMTPConfig, reader: jspb.BinaryReader): SMTPConfig;
}

export namespace SMTPConfig {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    senderAddress: string,
    senderName: string,
    tls: boolean,
    host: string,
    user: string,
    replyToAddress: string,
    state: SMTPConfigState,
    description: string,
    id: string,
  }
}

export class SMSProvider extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SMSProvider;
  hasDetails(): boolean;
  clearDetails(): SMSProvider;

  getId(): string;
  setId(value: string): SMSProvider;

  getState(): SMSProviderConfigState;
  setState(value: SMSProviderConfigState): SMSProvider;

  getTwilio(): TwilioConfig | undefined;
  setTwilio(value?: TwilioConfig): SMSProvider;
  hasTwilio(): boolean;
  clearTwilio(): SMSProvider;

  getConfigCase(): SMSProvider.ConfigCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SMSProvider.AsObject;
  static toObject(includeInstance: boolean, msg: SMSProvider): SMSProvider.AsObject;
  static serializeBinaryToWriter(message: SMSProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SMSProvider;
  static deserializeBinaryFromReader(message: SMSProvider, reader: jspb.BinaryReader): SMSProvider;
}

export namespace SMSProvider {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
    state: SMSProviderConfigState,
    twilio?: TwilioConfig.AsObject,
  }

  export enum ConfigCase { 
    CONFIG_NOT_SET = 0,
    TWILIO = 4,
  }
}

export class TwilioConfig extends jspb.Message {
  getSid(): string;
  setSid(value: string): TwilioConfig;

  getSenderNumber(): string;
  setSenderNumber(value: string): TwilioConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TwilioConfig.AsObject;
  static toObject(includeInstance: boolean, msg: TwilioConfig): TwilioConfig.AsObject;
  static serializeBinaryToWriter(message: TwilioConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TwilioConfig;
  static deserializeBinaryFromReader(message: TwilioConfig, reader: jspb.BinaryReader): TwilioConfig;
}

export namespace TwilioConfig {
  export type AsObject = {
    sid: string,
    senderNumber: string,
  }
}

export class DebugNotificationProvider extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DebugNotificationProvider;
  hasDetails(): boolean;
  clearDetails(): DebugNotificationProvider;

  getCompact(): boolean;
  setCompact(value: boolean): DebugNotificationProvider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DebugNotificationProvider.AsObject;
  static toObject(includeInstance: boolean, msg: DebugNotificationProvider): DebugNotificationProvider.AsObject;
  static serializeBinaryToWriter(message: DebugNotificationProvider, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DebugNotificationProvider;
  static deserializeBinaryFromReader(message: DebugNotificationProvider, reader: jspb.BinaryReader): DebugNotificationProvider;
}

export namespace DebugNotificationProvider {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    compact: boolean,
  }
}

export class OIDCSettings extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): OIDCSettings;
  hasDetails(): boolean;
  clearDetails(): OIDCSettings;

  getAccessTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setAccessTokenLifetime(value?: google_protobuf_duration_pb.Duration): OIDCSettings;
  hasAccessTokenLifetime(): boolean;
  clearAccessTokenLifetime(): OIDCSettings;

  getIdTokenLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setIdTokenLifetime(value?: google_protobuf_duration_pb.Duration): OIDCSettings;
  hasIdTokenLifetime(): boolean;
  clearIdTokenLifetime(): OIDCSettings;

  getRefreshTokenIdleExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenIdleExpiration(value?: google_protobuf_duration_pb.Duration): OIDCSettings;
  hasRefreshTokenIdleExpiration(): boolean;
  clearRefreshTokenIdleExpiration(): OIDCSettings;

  getRefreshTokenExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setRefreshTokenExpiration(value?: google_protobuf_duration_pb.Duration): OIDCSettings;
  hasRefreshTokenExpiration(): boolean;
  clearRefreshTokenExpiration(): OIDCSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OIDCSettings.AsObject;
  static toObject(includeInstance: boolean, msg: OIDCSettings): OIDCSettings.AsObject;
  static serializeBinaryToWriter(message: OIDCSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OIDCSettings;
  static deserializeBinaryFromReader(message: OIDCSettings, reader: jspb.BinaryReader): OIDCSettings;
}

export namespace OIDCSettings {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    accessTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    idTokenLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenIdleExpiration?: google_protobuf_duration_pb.Duration.AsObject,
    refreshTokenExpiration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class SecurityPolicy extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SecurityPolicy;
  hasDetails(): boolean;
  clearDetails(): SecurityPolicy;

  getEnableIframeEmbedding(): boolean;
  setEnableIframeEmbedding(value: boolean): SecurityPolicy;

  getAllowedOriginsList(): Array<string>;
  setAllowedOriginsList(value: Array<string>): SecurityPolicy;
  clearAllowedOriginsList(): SecurityPolicy;
  addAllowedOrigins(value: string, index?: number): SecurityPolicy;

  getEnableImpersonation(): boolean;
  setEnableImpersonation(value: boolean): SecurityPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SecurityPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: SecurityPolicy): SecurityPolicy.AsObject;
  static serializeBinaryToWriter(message: SecurityPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SecurityPolicy;
  static deserializeBinaryFromReader(message: SecurityPolicy, reader: jspb.BinaryReader): SecurityPolicy;
}

export namespace SecurityPolicy {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    enableIframeEmbedding: boolean,
    allowedOriginsList: Array<string>,
    enableImpersonation: boolean,
  }
}

export enum SMTPConfigState { 
  SMTP_CONFIG_STATE_UNSPECIFIED = 0,
  SMTP_CONFIG_ACTIVE = 1,
  SMTP_CONFIG_INACTIVE = 2,
}
export enum SecretGeneratorType { 
  SECRET_GENERATOR_TYPE_UNSPECIFIED = 0,
  SECRET_GENERATOR_TYPE_INIT_CODE = 1,
  SECRET_GENERATOR_TYPE_VERIFY_EMAIL_CODE = 2,
  SECRET_GENERATOR_TYPE_VERIFY_PHONE_CODE = 3,
  SECRET_GENERATOR_TYPE_PASSWORD_RESET_CODE = 4,
  SECRET_GENERATOR_TYPE_PASSWORDLESS_INIT_CODE = 5,
  SECRET_GENERATOR_TYPE_APP_SECRET = 6,
  SECRET_GENERATOR_TYPE_OTP_SMS = 7,
  SECRET_GENERATOR_TYPE_OTP_EMAIL = 8,
}
export enum SMSProviderConfigState { 
  SMS_PROVIDER_CONFIG_STATE_UNSPECIFIED = 0,
  SMS_PROVIDER_CONFIG_ACTIVE = 1,
  SMS_PROVIDER_CONFIG_INACTIVE = 2,
}
