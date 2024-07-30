import * as jspb from 'google-protobuf'

import * as zitadel_object_v2beta_object_pb from '../../../zitadel/object/v2beta/object_pb'; // proto import: "zitadel/object/v2beta/object.proto"
import * as zitadel_protoc_gen_zitadel_v2_options_pb from '../../../zitadel/protoc_gen_zitadel/v2/options_pb'; // proto import: "zitadel/protoc_gen_zitadel/v2/options.proto"
import * as zitadel_user_v2beta_auth_pb from '../../../zitadel/user/v2beta/auth_pb'; // proto import: "zitadel/user/v2beta/auth.proto"
import * as zitadel_user_v2beta_email_pb from '../../../zitadel/user/v2beta/email_pb'; // proto import: "zitadel/user/v2beta/email.proto"
import * as zitadel_user_v2beta_phone_pb from '../../../zitadel/user/v2beta/phone_pb'; // proto import: "zitadel/user/v2beta/phone.proto"
import * as zitadel_user_v2beta_idp_pb from '../../../zitadel/user/v2beta/idp_pb'; // proto import: "zitadel/user/v2beta/idp.proto"
import * as zitadel_user_v2beta_password_pb from '../../../zitadel/user/v2beta/password_pb'; // proto import: "zitadel/user/v2beta/password.proto"
import * as zitadel_user_v2beta_user_pb from '../../../zitadel/user/v2beta/user_pb'; // proto import: "zitadel/user/v2beta/user.proto"
import * as zitadel_user_v2beta_user_service_pb from '../../../zitadel/user/v2beta/user_service_pb'; // proto import: "zitadel/user/v2beta/user_service.proto"
import * as google_api_annotations_pb from '../../../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../../../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb'; // proto import: "google/protobuf/struct.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../../../validate/validate_pb'; // proto import: "validate/validate.proto"


export class AddOrganizationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddOrganizationRequest;

  getAdminsList(): Array<AddOrganizationRequest.Admin>;
  setAdminsList(value: Array<AddOrganizationRequest.Admin>): AddOrganizationRequest;
  clearAdminsList(): AddOrganizationRequest;
  addAdmins(value?: AddOrganizationRequest.Admin, index?: number): AddOrganizationRequest.Admin;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrganizationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrganizationRequest): AddOrganizationRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrganizationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrganizationRequest;
  static deserializeBinaryFromReader(message: AddOrganizationRequest, reader: jspb.BinaryReader): AddOrganizationRequest;
}

export namespace AddOrganizationRequest {
  export type AsObject = {
    name: string,
    adminsList: Array<AddOrganizationRequest.Admin.AsObject>,
  }

  export class Admin extends jspb.Message {
    getUserId(): string;
    setUserId(value: string): Admin;

    getHuman(): zitadel_user_v2beta_user_service_pb.AddHumanUserRequest | undefined;
    setHuman(value?: zitadel_user_v2beta_user_service_pb.AddHumanUserRequest): Admin;
    hasHuman(): boolean;
    clearHuman(): Admin;

    getRolesList(): Array<string>;
    setRolesList(value: Array<string>): Admin;
    clearRolesList(): Admin;
    addRoles(value: string, index?: number): Admin;

    getUserTypeCase(): Admin.UserTypeCase;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Admin.AsObject;
    static toObject(includeInstance: boolean, msg: Admin): Admin.AsObject;
    static serializeBinaryToWriter(message: Admin, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Admin;
    static deserializeBinaryFromReader(message: Admin, reader: jspb.BinaryReader): Admin;
  }

  export namespace Admin {
    export type AsObject = {
      userId: string,
      human?: zitadel_user_v2beta_user_service_pb.AddHumanUserRequest.AsObject,
      rolesList: Array<string>,
    }

    export enum UserTypeCase { 
      USER_TYPE_NOT_SET = 0,
      USER_ID = 1,
      HUMAN = 2,
    }
  }

}

export class AddOrganizationResponse extends jspb.Message {
  getDetails(): zitadel_object_v2beta_object_pb.Details | undefined;
  setDetails(value?: zitadel_object_v2beta_object_pb.Details): AddOrganizationResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrganizationResponse;

  getOrganizationId(): string;
  setOrganizationId(value: string): AddOrganizationResponse;

  getCreatedAdminsList(): Array<AddOrganizationResponse.CreatedAdmin>;
  setCreatedAdminsList(value: Array<AddOrganizationResponse.CreatedAdmin>): AddOrganizationResponse;
  clearCreatedAdminsList(): AddOrganizationResponse;
  addCreatedAdmins(value?: AddOrganizationResponse.CreatedAdmin, index?: number): AddOrganizationResponse.CreatedAdmin;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrganizationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrganizationResponse): AddOrganizationResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrganizationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrganizationResponse;
  static deserializeBinaryFromReader(message: AddOrganizationResponse, reader: jspb.BinaryReader): AddOrganizationResponse;
}

export namespace AddOrganizationResponse {
  export type AsObject = {
    details?: zitadel_object_v2beta_object_pb.Details.AsObject,
    organizationId: string,
    createdAdminsList: Array<AddOrganizationResponse.CreatedAdmin.AsObject>,
  }

  export class CreatedAdmin extends jspb.Message {
    getUserId(): string;
    setUserId(value: string): CreatedAdmin;

    getEmailCode(): string;
    setEmailCode(value: string): CreatedAdmin;
    hasEmailCode(): boolean;
    clearEmailCode(): CreatedAdmin;

    getPhoneCode(): string;
    setPhoneCode(value: string): CreatedAdmin;
    hasPhoneCode(): boolean;
    clearPhoneCode(): CreatedAdmin;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreatedAdmin.AsObject;
    static toObject(includeInstance: boolean, msg: CreatedAdmin): CreatedAdmin.AsObject;
    static serializeBinaryToWriter(message: CreatedAdmin, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreatedAdmin;
    static deserializeBinaryFromReader(message: CreatedAdmin, reader: jspb.BinaryReader): CreatedAdmin;
  }

  export namespace CreatedAdmin {
    export type AsObject = {
      userId: string,
      emailCode?: string,
      phoneCode?: string,
    }

    export enum EmailCodeCase { 
      _EMAIL_CODE_NOT_SET = 0,
      EMAIL_CODE = 2,
    }

    export enum PhoneCodeCase { 
      _PHONE_CODE_NOT_SET = 0,
      PHONE_CODE = 3,
    }
  }

}

