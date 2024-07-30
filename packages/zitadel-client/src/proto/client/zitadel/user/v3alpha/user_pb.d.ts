import * as jspb from 'google-protobuf'

import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_user_v3alpha_authenticator_pb from '../../../zitadel/user/v3alpha/authenticator_pb'; // proto import: "zitadel/user/v3alpha/authenticator.proto"
import * as zitadel_user_v3alpha_communication_pb from '../../../zitadel/user/v3alpha/communication_pb'; // proto import: "zitadel/user/v3alpha/communication.proto"


export class User extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): User;

  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): User;
  hasDetails(): boolean;
  clearDetails(): User;

  getAuthenticators(): zitadel_user_v3alpha_authenticator_pb.Authenticators | undefined;
  setAuthenticators(value?: zitadel_user_v3alpha_authenticator_pb.Authenticators): User;
  hasAuthenticators(): boolean;
  clearAuthenticators(): User;

  getContact(): zitadel_user_v3alpha_communication_pb.Contact | undefined;
  setContact(value?: zitadel_user_v3alpha_communication_pb.Contact): User;
  hasContact(): boolean;
  clearContact(): User;

  getState(): State;
  setState(value: State): User;

  getSchema(): Schema | undefined;
  setSchema(value?: Schema): User;
  hasSchema(): boolean;
  clearSchema(): User;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): User;
  hasData(): boolean;
  clearData(): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    authenticators?: zitadel_user_v3alpha_authenticator_pb.Authenticators.AsObject,
    contact?: zitadel_user_v3alpha_communication_pb.Contact.AsObject,
    state: State,
    schema?: Schema.AsObject,
    data?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class Schema extends jspb.Message {
  getId(): string;
  setId(value: string): Schema;

  getType(): string;
  setType(value: string): Schema;

  getRevision(): number;
  setRevision(value: number): Schema;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Schema.AsObject;
  static toObject(includeInstance: boolean, msg: Schema): Schema.AsObject;
  static serializeBinaryToWriter(message: Schema, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Schema;
  static deserializeBinaryFromReader(message: Schema, reader: jspb.BinaryReader): Schema;
}

export namespace Schema {
  export type AsObject = {
    id: string,
    type: string,
    revision: number,
  }
}

export enum State { 
  USER_STATE_UNSPECIFIED = 0,
  USER_STATE_ACTIVE = 1,
  USER_STATE_INACTIVE = 2,
  USER_STATE_DELETED = 3,
  USER_STATE_LOCKED = 4,
}
