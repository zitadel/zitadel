import * as jspb from 'google-protobuf'

import * as zitadel_app_pb from '../zitadel/app_pb'; // proto import: "zitadel/app.proto"
import * as zitadel_idp_pb from '../zitadel/idp_pb'; // proto import: "zitadel/idp.proto"
import * as zitadel_user_pb from '../zitadel/user_pb'; // proto import: "zitadel/user.proto"
import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_options_pb from '../zitadel/options_pb'; // proto import: "zitadel/options.proto"
import * as zitadel_org_pb from '../zitadel/org_pb'; // proto import: "zitadel/org.proto"
import * as zitadel_member_pb from '../zitadel/member_pb'; // proto import: "zitadel/member.proto"
import * as zitadel_project_pb from '../zitadel/project_pb'; // proto import: "zitadel/project.proto"
import * as zitadel_policy_pb from '../zitadel/policy_pb'; // proto import: "zitadel/policy.proto"
import * as zitadel_text_pb from '../zitadel/text_pb'; // proto import: "zitadel/text.proto"
import * as zitadel_message_pb from '../zitadel/message_pb'; // proto import: "zitadel/message.proto"
import * as zitadel_change_pb from '../zitadel/change_pb'; // proto import: "zitadel/change.proto"
import * as zitadel_auth_n_key_pb from '../zitadel/auth_n_key_pb'; // proto import: "zitadel/auth_n_key.proto"
import * as zitadel_metadata_pb from '../zitadel/metadata_pb'; // proto import: "zitadel/metadata.proto"
import * as zitadel_action_pb from '../zitadel/action_pb'; // proto import: "zitadel/action.proto"
import * as google_api_annotations_pb from '../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as google_api_field_behavior_pb from '../google/api/field_behavior_pb'; // proto import: "google/api/field_behavior.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as google_protobuf_duration_pb from 'google-protobuf/google/protobuf/duration_pb'; // proto import: "google/protobuf/duration.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class HealthzRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthzRequest.AsObject;
  static toObject(includeInstance: boolean, msg: HealthzRequest): HealthzRequest.AsObject;
  static serializeBinaryToWriter(message: HealthzRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthzRequest;
  static deserializeBinaryFromReader(message: HealthzRequest, reader: jspb.BinaryReader): HealthzRequest;
}

export namespace HealthzRequest {
  export type AsObject = {
  }
}

export class HealthzResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HealthzResponse.AsObject;
  static toObject(includeInstance: boolean, msg: HealthzResponse): HealthzResponse.AsObject;
  static serializeBinaryToWriter(message: HealthzResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HealthzResponse;
  static deserializeBinaryFromReader(message: HealthzResponse, reader: jspb.BinaryReader): HealthzResponse;
}

export namespace HealthzResponse {
  export type AsObject = {
  }
}

export class GetOIDCInformationRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCInformationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCInformationRequest): GetOIDCInformationRequest.AsObject;
  static serializeBinaryToWriter(message: GetOIDCInformationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCInformationRequest;
  static deserializeBinaryFromReader(message: GetOIDCInformationRequest, reader: jspb.BinaryReader): GetOIDCInformationRequest;
}

export namespace GetOIDCInformationRequest {
  export type AsObject = {
  }
}

export class GetOIDCInformationResponse extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): GetOIDCInformationResponse;

  getDiscoveryEndpoint(): string;
  setDiscoveryEndpoint(value: string): GetOIDCInformationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOIDCInformationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOIDCInformationResponse): GetOIDCInformationResponse.AsObject;
  static serializeBinaryToWriter(message: GetOIDCInformationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOIDCInformationResponse;
  static deserializeBinaryFromReader(message: GetOIDCInformationResponse, reader: jspb.BinaryReader): GetOIDCInformationResponse;
}

export namespace GetOIDCInformationResponse {
  export type AsObject = {
    issuer: string,
    discoveryEndpoint: string,
  }
}

export class GetIAMRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIAMRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetIAMRequest): GetIAMRequest.AsObject;
  static serializeBinaryToWriter(message: GetIAMRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIAMRequest;
  static deserializeBinaryFromReader(message: GetIAMRequest, reader: jspb.BinaryReader): GetIAMRequest;
}

export namespace GetIAMRequest {
  export type AsObject = {
  }
}

export class GetIAMResponse extends jspb.Message {
  getGlobalOrgId(): string;
  setGlobalOrgId(value: string): GetIAMResponse;

  getIamProjectId(): string;
  setIamProjectId(value: string): GetIAMResponse;

  getDefaultOrgId(): string;
  setDefaultOrgId(value: string): GetIAMResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetIAMResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetIAMResponse): GetIAMResponse.AsObject;
  static serializeBinaryToWriter(message: GetIAMResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetIAMResponse;
  static deserializeBinaryFromReader(message: GetIAMResponse, reader: jspb.BinaryReader): GetIAMResponse;
}

export namespace GetIAMResponse {
  export type AsObject = {
    globalOrgId: string,
    iamProjectId: string,
    defaultOrgId: string,
  }
}

export class GetSupportedLanguagesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupportedLanguagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupportedLanguagesRequest): GetSupportedLanguagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetSupportedLanguagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupportedLanguagesRequest;
  static deserializeBinaryFromReader(message: GetSupportedLanguagesRequest, reader: jspb.BinaryReader): GetSupportedLanguagesRequest;
}

export namespace GetSupportedLanguagesRequest {
  export type AsObject = {
  }
}

export class GetSupportedLanguagesResponse extends jspb.Message {
  getLanguagesList(): Array<string>;
  setLanguagesList(value: Array<string>): GetSupportedLanguagesResponse;
  clearLanguagesList(): GetSupportedLanguagesResponse;
  addLanguages(value: string, index?: number): GetSupportedLanguagesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSupportedLanguagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetSupportedLanguagesResponse): GetSupportedLanguagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetSupportedLanguagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSupportedLanguagesResponse;
  static deserializeBinaryFromReader(message: GetSupportedLanguagesResponse, reader: jspb.BinaryReader): GetSupportedLanguagesResponse;
}

export namespace GetSupportedLanguagesResponse {
  export type AsObject = {
    languagesList: Array<string>,
  }
}

export class GetUserByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetUserByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByIDRequest): GetUserByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByIDRequest;
  static deserializeBinaryFromReader(message: GetUserByIDRequest, reader: jspb.BinaryReader): GetUserByIDRequest;
}

export namespace GetUserByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetUserByIDResponse extends jspb.Message {
  getUser(): zitadel_user_pb.User | undefined;
  setUser(value?: zitadel_user_pb.User): GetUserByIDResponse;
  hasUser(): boolean;
  clearUser(): GetUserByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByIDResponse): GetUserByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByIDResponse;
  static deserializeBinaryFromReader(message: GetUserByIDResponse, reader: jspb.BinaryReader): GetUserByIDResponse;
}

export namespace GetUserByIDResponse {
  export type AsObject = {
    user?: zitadel_user_pb.User.AsObject,
  }
}

export class GetUserByLoginNameGlobalRequest extends jspb.Message {
  getLoginName(): string;
  setLoginName(value: string): GetUserByLoginNameGlobalRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByLoginNameGlobalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByLoginNameGlobalRequest): GetUserByLoginNameGlobalRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserByLoginNameGlobalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByLoginNameGlobalRequest;
  static deserializeBinaryFromReader(message: GetUserByLoginNameGlobalRequest, reader: jspb.BinaryReader): GetUserByLoginNameGlobalRequest;
}

export namespace GetUserByLoginNameGlobalRequest {
  export type AsObject = {
    loginName: string,
  }
}

export class GetUserByLoginNameGlobalResponse extends jspb.Message {
  getUser(): zitadel_user_pb.User | undefined;
  setUser(value?: zitadel_user_pb.User): GetUserByLoginNameGlobalResponse;
  hasUser(): boolean;
  clearUser(): GetUserByLoginNameGlobalResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserByLoginNameGlobalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserByLoginNameGlobalResponse): GetUserByLoginNameGlobalResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserByLoginNameGlobalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserByLoginNameGlobalResponse;
  static deserializeBinaryFromReader(message: GetUserByLoginNameGlobalResponse, reader: jspb.BinaryReader): GetUserByLoginNameGlobalResponse;
}

export namespace GetUserByLoginNameGlobalResponse {
  export type AsObject = {
    user?: zitadel_user_pb.User.AsObject,
  }
}

export class ListUsersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListUsersRequest;
  hasQuery(): boolean;
  clearQuery(): ListUsersRequest;

  getSortingColumn(): zitadel_user_pb.UserFieldName;
  setSortingColumn(value: zitadel_user_pb.UserFieldName): ListUsersRequest;

  getQueriesList(): Array<zitadel_user_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_user_pb.SearchQuery>): ListUsersRequest;
  clearQueriesList(): ListUsersRequest;
  addQueries(value?: zitadel_user_pb.SearchQuery, index?: number): zitadel_user_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersRequest): ListUsersRequest.AsObject;
  static serializeBinaryToWriter(message: ListUsersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersRequest;
  static deserializeBinaryFromReader(message: ListUsersRequest, reader: jspb.BinaryReader): ListUsersRequest;
}

export namespace ListUsersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_user_pb.UserFieldName,
    queriesList: Array<zitadel_user_pb.SearchQuery.AsObject>,
  }
}

export class ListUsersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListUsersResponse;
  hasDetails(): boolean;
  clearDetails(): ListUsersResponse;

  getSortingColumn(): zitadel_user_pb.UserFieldName;
  setSortingColumn(value: zitadel_user_pb.UserFieldName): ListUsersResponse;

  getResultList(): Array<zitadel_user_pb.User>;
  setResultList(value: Array<zitadel_user_pb.User>): ListUsersResponse;
  clearResultList(): ListUsersResponse;
  addResult(value?: zitadel_user_pb.User, index?: number): zitadel_user_pb.User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUsersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUsersResponse): ListUsersResponse.AsObject;
  static serializeBinaryToWriter(message: ListUsersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUsersResponse;
  static deserializeBinaryFromReader(message: ListUsersResponse, reader: jspb.BinaryReader): ListUsersResponse;
}

export namespace ListUsersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_user_pb.UserFieldName,
    resultList: Array<zitadel_user_pb.User.AsObject>,
  }
}

export class ListUserChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListUserChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListUserChangesRequest;

  getUserId(): string;
  setUserId(value: string): ListUserChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserChangesRequest): ListUserChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserChangesRequest;
  static deserializeBinaryFromReader(message: ListUserChangesRequest, reader: jspb.BinaryReader): ListUserChangesRequest;
}

export namespace ListUserChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
    userId: string,
  }
}

export class ListUserChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListUserChangesResponse;
  clearResultList(): ListUserChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserChangesResponse): ListUserChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserChangesResponse;
  static deserializeBinaryFromReader(message: ListUserChangesResponse, reader: jspb.BinaryReader): ListUserChangesResponse;
}

export namespace ListUserChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class IsUserUniqueRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): IsUserUniqueRequest;

  getEmail(): string;
  setEmail(value: string): IsUserUniqueRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsUserUniqueRequest.AsObject;
  static toObject(includeInstance: boolean, msg: IsUserUniqueRequest): IsUserUniqueRequest.AsObject;
  static serializeBinaryToWriter(message: IsUserUniqueRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsUserUniqueRequest;
  static deserializeBinaryFromReader(message: IsUserUniqueRequest, reader: jspb.BinaryReader): IsUserUniqueRequest;
}

export namespace IsUserUniqueRequest {
  export type AsObject = {
    userName: string,
    email: string,
  }
}

export class IsUserUniqueResponse extends jspb.Message {
  getIsUnique(): boolean;
  setIsUnique(value: boolean): IsUserUniqueResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsUserUniqueResponse.AsObject;
  static toObject(includeInstance: boolean, msg: IsUserUniqueResponse): IsUserUniqueResponse.AsObject;
  static serializeBinaryToWriter(message: IsUserUniqueResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsUserUniqueResponse;
  static deserializeBinaryFromReader(message: IsUserUniqueResponse, reader: jspb.BinaryReader): IsUserUniqueResponse;
}

export namespace IsUserUniqueResponse {
  export type AsObject = {
    isUnique: boolean,
  }
}

export class AddHumanUserRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): AddHumanUserRequest;

  getProfile(): AddHumanUserRequest.Profile | undefined;
  setProfile(value?: AddHumanUserRequest.Profile): AddHumanUserRequest;
  hasProfile(): boolean;
  clearProfile(): AddHumanUserRequest;

  getEmail(): AddHumanUserRequest.Email | undefined;
  setEmail(value?: AddHumanUserRequest.Email): AddHumanUserRequest;
  hasEmail(): boolean;
  clearEmail(): AddHumanUserRequest;

  getPhone(): AddHumanUserRequest.Phone | undefined;
  setPhone(value?: AddHumanUserRequest.Phone): AddHumanUserRequest;
  hasPhone(): boolean;
  clearPhone(): AddHumanUserRequest;

  getInitialPassword(): string;
  setInitialPassword(value: string): AddHumanUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddHumanUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddHumanUserRequest): AddHumanUserRequest.AsObject;
  static serializeBinaryToWriter(message: AddHumanUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddHumanUserRequest;
  static deserializeBinaryFromReader(message: AddHumanUserRequest, reader: jspb.BinaryReader): AddHumanUserRequest;
}

export namespace AddHumanUserRequest {
  export type AsObject = {
    userName: string,
    profile?: AddHumanUserRequest.Profile.AsObject,
    email?: AddHumanUserRequest.Email.AsObject,
    phone?: AddHumanUserRequest.Phone.AsObject,
    initialPassword: string,
  }

  export class Profile extends jspb.Message {
    getFirstName(): string;
    setFirstName(value: string): Profile;

    getLastName(): string;
    setLastName(value: string): Profile;

    getNickName(): string;
    setNickName(value: string): Profile;

    getDisplayName(): string;
    setDisplayName(value: string): Profile;

    getPreferredLanguage(): string;
    setPreferredLanguage(value: string): Profile;

    getGender(): zitadel_user_pb.Gender;
    setGender(value: zitadel_user_pb.Gender): Profile;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Profile.AsObject;
    static toObject(includeInstance: boolean, msg: Profile): Profile.AsObject;
    static serializeBinaryToWriter(message: Profile, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Profile;
    static deserializeBinaryFromReader(message: Profile, reader: jspb.BinaryReader): Profile;
  }

  export namespace Profile {
    export type AsObject = {
      firstName: string,
      lastName: string,
      nickName: string,
      displayName: string,
      preferredLanguage: string,
      gender: zitadel_user_pb.Gender,
    }
  }


  export class Email extends jspb.Message {
    getEmail(): string;
    setEmail(value: string): Email;

    getIsEmailVerified(): boolean;
    setIsEmailVerified(value: boolean): Email;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Email.AsObject;
    static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
    static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Email;
    static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
  }

  export namespace Email {
    export type AsObject = {
      email: string,
      isEmailVerified: boolean,
    }
  }


  export class Phone extends jspb.Message {
    getPhone(): string;
    setPhone(value: string): Phone;

    getIsPhoneVerified(): boolean;
    setIsPhoneVerified(value: boolean): Phone;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Phone.AsObject;
    static toObject(includeInstance: boolean, msg: Phone): Phone.AsObject;
    static serializeBinaryToWriter(message: Phone, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Phone;
    static deserializeBinaryFromReader(message: Phone, reader: jspb.BinaryReader): Phone;
  }

  export namespace Phone {
    export type AsObject = {
      phone: string,
      isPhoneVerified: boolean,
    }
  }

}

export class AddHumanUserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddHumanUserResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddHumanUserResponse;
  hasDetails(): boolean;
  clearDetails(): AddHumanUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddHumanUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddHumanUserResponse): AddHumanUserResponse.AsObject;
  static serializeBinaryToWriter(message: AddHumanUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddHumanUserResponse;
  static deserializeBinaryFromReader(message: AddHumanUserResponse, reader: jspb.BinaryReader): AddHumanUserResponse;
}

export namespace AddHumanUserResponse {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ImportHumanUserRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): ImportHumanUserRequest;

  getProfile(): ImportHumanUserRequest.Profile | undefined;
  setProfile(value?: ImportHumanUserRequest.Profile): ImportHumanUserRequest;
  hasProfile(): boolean;
  clearProfile(): ImportHumanUserRequest;

  getEmail(): ImportHumanUserRequest.Email | undefined;
  setEmail(value?: ImportHumanUserRequest.Email): ImportHumanUserRequest;
  hasEmail(): boolean;
  clearEmail(): ImportHumanUserRequest;

  getPhone(): ImportHumanUserRequest.Phone | undefined;
  setPhone(value?: ImportHumanUserRequest.Phone): ImportHumanUserRequest;
  hasPhone(): boolean;
  clearPhone(): ImportHumanUserRequest;

  getPassword(): string;
  setPassword(value: string): ImportHumanUserRequest;

  getHashedPassword(): ImportHumanUserRequest.HashedPassword | undefined;
  setHashedPassword(value?: ImportHumanUserRequest.HashedPassword): ImportHumanUserRequest;
  hasHashedPassword(): boolean;
  clearHashedPassword(): ImportHumanUserRequest;

  getPasswordChangeRequired(): boolean;
  setPasswordChangeRequired(value: boolean): ImportHumanUserRequest;

  getRequestPasswordlessRegistration(): boolean;
  setRequestPasswordlessRegistration(value: boolean): ImportHumanUserRequest;

  getOtpCode(): string;
  setOtpCode(value: string): ImportHumanUserRequest;

  getIdpsList(): Array<ImportHumanUserRequest.IDP>;
  setIdpsList(value: Array<ImportHumanUserRequest.IDP>): ImportHumanUserRequest;
  clearIdpsList(): ImportHumanUserRequest;
  addIdps(value?: ImportHumanUserRequest.IDP, index?: number): ImportHumanUserRequest.IDP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportHumanUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ImportHumanUserRequest): ImportHumanUserRequest.AsObject;
  static serializeBinaryToWriter(message: ImportHumanUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportHumanUserRequest;
  static deserializeBinaryFromReader(message: ImportHumanUserRequest, reader: jspb.BinaryReader): ImportHumanUserRequest;
}

export namespace ImportHumanUserRequest {
  export type AsObject = {
    userName: string,
    profile?: ImportHumanUserRequest.Profile.AsObject,
    email?: ImportHumanUserRequest.Email.AsObject,
    phone?: ImportHumanUserRequest.Phone.AsObject,
    password: string,
    hashedPassword?: ImportHumanUserRequest.HashedPassword.AsObject,
    passwordChangeRequired: boolean,
    requestPasswordlessRegistration: boolean,
    otpCode: string,
    idpsList: Array<ImportHumanUserRequest.IDP.AsObject>,
  }

  export class Profile extends jspb.Message {
    getFirstName(): string;
    setFirstName(value: string): Profile;

    getLastName(): string;
    setLastName(value: string): Profile;

    getNickName(): string;
    setNickName(value: string): Profile;

    getDisplayName(): string;
    setDisplayName(value: string): Profile;

    getPreferredLanguage(): string;
    setPreferredLanguage(value: string): Profile;

    getGender(): zitadel_user_pb.Gender;
    setGender(value: zitadel_user_pb.Gender): Profile;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Profile.AsObject;
    static toObject(includeInstance: boolean, msg: Profile): Profile.AsObject;
    static serializeBinaryToWriter(message: Profile, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Profile;
    static deserializeBinaryFromReader(message: Profile, reader: jspb.BinaryReader): Profile;
  }

  export namespace Profile {
    export type AsObject = {
      firstName: string,
      lastName: string,
      nickName: string,
      displayName: string,
      preferredLanguage: string,
      gender: zitadel_user_pb.Gender,
    }
  }


  export class Email extends jspb.Message {
    getEmail(): string;
    setEmail(value: string): Email;

    getIsEmailVerified(): boolean;
    setIsEmailVerified(value: boolean): Email;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Email.AsObject;
    static toObject(includeInstance: boolean, msg: Email): Email.AsObject;
    static serializeBinaryToWriter(message: Email, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Email;
    static deserializeBinaryFromReader(message: Email, reader: jspb.BinaryReader): Email;
  }

  export namespace Email {
    export type AsObject = {
      email: string,
      isEmailVerified: boolean,
    }
  }


  export class Phone extends jspb.Message {
    getPhone(): string;
    setPhone(value: string): Phone;

    getIsPhoneVerified(): boolean;
    setIsPhoneVerified(value: boolean): Phone;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Phone.AsObject;
    static toObject(includeInstance: boolean, msg: Phone): Phone.AsObject;
    static serializeBinaryToWriter(message: Phone, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Phone;
    static deserializeBinaryFromReader(message: Phone, reader: jspb.BinaryReader): Phone;
  }

  export namespace Phone {
    export type AsObject = {
      phone: string,
      isPhoneVerified: boolean,
    }
  }


  export class HashedPassword extends jspb.Message {
    getValue(): string;
    setValue(value: string): HashedPassword;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HashedPassword.AsObject;
    static toObject(includeInstance: boolean, msg: HashedPassword): HashedPassword.AsObject;
    static serializeBinaryToWriter(message: HashedPassword, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HashedPassword;
    static deserializeBinaryFromReader(message: HashedPassword, reader: jspb.BinaryReader): HashedPassword;
  }

  export namespace HashedPassword {
    export type AsObject = {
      value: string,
    }
  }


  export class IDP extends jspb.Message {
    getConfigId(): string;
    setConfigId(value: string): IDP;

    getExternalUserId(): string;
    setExternalUserId(value: string): IDP;

    getDisplayName(): string;
    setDisplayName(value: string): IDP;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): IDP.AsObject;
    static toObject(includeInstance: boolean, msg: IDP): IDP.AsObject;
    static serializeBinaryToWriter(message: IDP, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): IDP;
    static deserializeBinaryFromReader(message: IDP, reader: jspb.BinaryReader): IDP;
  }

  export namespace IDP {
    export type AsObject = {
      configId: string,
      externalUserId: string,
      displayName: string,
    }
  }

}

export class ImportHumanUserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ImportHumanUserResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ImportHumanUserResponse;
  hasDetails(): boolean;
  clearDetails(): ImportHumanUserResponse;

  getPasswordlessRegistration(): ImportHumanUserResponse.PasswordlessRegistration | undefined;
  setPasswordlessRegistration(value?: ImportHumanUserResponse.PasswordlessRegistration): ImportHumanUserResponse;
  hasPasswordlessRegistration(): boolean;
  clearPasswordlessRegistration(): ImportHumanUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportHumanUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ImportHumanUserResponse): ImportHumanUserResponse.AsObject;
  static serializeBinaryToWriter(message: ImportHumanUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportHumanUserResponse;
  static deserializeBinaryFromReader(message: ImportHumanUserResponse, reader: jspb.BinaryReader): ImportHumanUserResponse;
}

export namespace ImportHumanUserResponse {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    passwordlessRegistration?: ImportHumanUserResponse.PasswordlessRegistration.AsObject,
  }

  export class PasswordlessRegistration extends jspb.Message {
    getLink(): string;
    setLink(value: string): PasswordlessRegistration;

    getLifetime(): google_protobuf_duration_pb.Duration | undefined;
    setLifetime(value?: google_protobuf_duration_pb.Duration): PasswordlessRegistration;
    hasLifetime(): boolean;
    clearLifetime(): PasswordlessRegistration;

    getExpiration(): google_protobuf_duration_pb.Duration | undefined;
    setExpiration(value?: google_protobuf_duration_pb.Duration): PasswordlessRegistration;
    hasExpiration(): boolean;
    clearExpiration(): PasswordlessRegistration;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PasswordlessRegistration.AsObject;
    static toObject(includeInstance: boolean, msg: PasswordlessRegistration): PasswordlessRegistration.AsObject;
    static serializeBinaryToWriter(message: PasswordlessRegistration, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PasswordlessRegistration;
    static deserializeBinaryFromReader(message: PasswordlessRegistration, reader: jspb.BinaryReader): PasswordlessRegistration;
  }

  export namespace PasswordlessRegistration {
    export type AsObject = {
      link: string,
      lifetime?: google_protobuf_duration_pb.Duration.AsObject,
      expiration?: google_protobuf_duration_pb.Duration.AsObject,
    }
  }

}

export class AddMachineUserRequest extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): AddMachineUserRequest;

  getName(): string;
  setName(value: string): AddMachineUserRequest;

  getDescription(): string;
  setDescription(value: string): AddMachineUserRequest;

  getAccessTokenType(): zitadel_user_pb.AccessTokenType;
  setAccessTokenType(value: zitadel_user_pb.AccessTokenType): AddMachineUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMachineUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMachineUserRequest): AddMachineUserRequest.AsObject;
  static serializeBinaryToWriter(message: AddMachineUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMachineUserRequest;
  static deserializeBinaryFromReader(message: AddMachineUserRequest, reader: jspb.BinaryReader): AddMachineUserRequest;
}

export namespace AddMachineUserRequest {
  export type AsObject = {
    userName: string,
    name: string,
    description: string,
    accessTokenType: zitadel_user_pb.AccessTokenType,
  }
}

export class AddMachineUserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddMachineUserResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMachineUserResponse;
  hasDetails(): boolean;
  clearDetails(): AddMachineUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMachineUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMachineUserResponse): AddMachineUserResponse.AsObject;
  static serializeBinaryToWriter(message: AddMachineUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMachineUserResponse;
  static deserializeBinaryFromReader(message: AddMachineUserResponse, reader: jspb.BinaryReader): AddMachineUserResponse;
}

export namespace AddMachineUserResponse {
  export type AsObject = {
    userId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserRequest): DeactivateUserRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserRequest;
  static deserializeBinaryFromReader(message: DeactivateUserRequest, reader: jspb.BinaryReader): DeactivateUserRequest;
}

export namespace DeactivateUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateUserResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserResponse): DeactivateUserResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserResponse;
  static deserializeBinaryFromReader(message: DeactivateUserResponse, reader: jspb.BinaryReader): DeactivateUserResponse;
}

export namespace DeactivateUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ReactivateUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserRequest): ReactivateUserRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserRequest;
  static deserializeBinaryFromReader(message: ReactivateUserRequest, reader: jspb.BinaryReader): ReactivateUserRequest;
}

export namespace ReactivateUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class ReactivateUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateUserResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserResponse): ReactivateUserResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserResponse;
  static deserializeBinaryFromReader(message: ReactivateUserResponse, reader: jspb.BinaryReader): ReactivateUserResponse;
}

export namespace ReactivateUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class LockUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): LockUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LockUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LockUserRequest): LockUserRequest.AsObject;
  static serializeBinaryToWriter(message: LockUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LockUserRequest;
  static deserializeBinaryFromReader(message: LockUserRequest, reader: jspb.BinaryReader): LockUserRequest;
}

export namespace LockUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class LockUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): LockUserResponse;
  hasDetails(): boolean;
  clearDetails(): LockUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LockUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LockUserResponse): LockUserResponse.AsObject;
  static serializeBinaryToWriter(message: LockUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LockUserResponse;
  static deserializeBinaryFromReader(message: LockUserResponse, reader: jspb.BinaryReader): LockUserResponse;
}

export namespace LockUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UnlockUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UnlockUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnlockUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UnlockUserRequest): UnlockUserRequest.AsObject;
  static serializeBinaryToWriter(message: UnlockUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnlockUserRequest;
  static deserializeBinaryFromReader(message: UnlockUserRequest, reader: jspb.BinaryReader): UnlockUserRequest;
}

export namespace UnlockUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class UnlockUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UnlockUserResponse;
  hasDetails(): boolean;
  clearDetails(): UnlockUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UnlockUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UnlockUserResponse): UnlockUserResponse.AsObject;
  static serializeBinaryToWriter(message: UnlockUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UnlockUserResponse;
  static deserializeBinaryFromReader(message: UnlockUserResponse, reader: jspb.BinaryReader): UnlockUserResponse;
}

export namespace UnlockUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveUserRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RemoveUserRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserRequest): RemoveUserRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveUserRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserRequest;
  static deserializeBinaryFromReader(message: RemoveUserRequest, reader: jspb.BinaryReader): RemoveUserRequest;
}

export namespace RemoveUserRequest {
  export type AsObject = {
    id: string,
  }
}

export class RemoveUserResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveUserResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveUserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserResponse): RemoveUserResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveUserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserResponse;
  static deserializeBinaryFromReader(message: RemoveUserResponse, reader: jspb.BinaryReader): RemoveUserResponse;
}

export namespace RemoveUserResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateUserNameRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateUserNameRequest;

  getUserName(): string;
  setUserName(value: string): UpdateUserNameRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserNameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserNameRequest): UpdateUserNameRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserNameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserNameRequest;
  static deserializeBinaryFromReader(message: UpdateUserNameRequest, reader: jspb.BinaryReader): UpdateUserNameRequest;
}

export namespace UpdateUserNameRequest {
  export type AsObject = {
    userId: string,
    userName: string,
  }
}

export class UpdateUserNameResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateUserNameResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateUserNameResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserNameResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserNameResponse): UpdateUserNameResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateUserNameResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserNameResponse;
  static deserializeBinaryFromReader(message: UpdateUserNameResponse, reader: jspb.BinaryReader): UpdateUserNameResponse;
}

export namespace UpdateUserNameResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ListUserMetadataRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListUserMetadataRequest;
  hasQuery(): boolean;
  clearQuery(): ListUserMetadataRequest;

  getQueriesList(): Array<zitadel_metadata_pb.MetadataQuery>;
  setQueriesList(value: Array<zitadel_metadata_pb.MetadataQuery>): ListUserMetadataRequest;
  clearQueriesList(): ListUserMetadataRequest;
  addQueries(value?: zitadel_metadata_pb.MetadataQuery, index?: number): zitadel_metadata_pb.MetadataQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserMetadataRequest): ListUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserMetadataRequest;
  static deserializeBinaryFromReader(message: ListUserMetadataRequest, reader: jspb.BinaryReader): ListUserMetadataRequest;
}

export namespace ListUserMetadataRequest {
  export type AsObject = {
    id: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_metadata_pb.MetadataQuery.AsObject>,
  }
}

export class ListUserMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListUserMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): ListUserMetadataResponse;

  getResultList(): Array<zitadel_metadata_pb.Metadata>;
  setResultList(value: Array<zitadel_metadata_pb.Metadata>): ListUserMetadataResponse;
  clearResultList(): ListUserMetadataResponse;
  addResult(value?: zitadel_metadata_pb.Metadata, index?: number): zitadel_metadata_pb.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserMetadataResponse): ListUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserMetadataResponse;
  static deserializeBinaryFromReader(message: ListUserMetadataResponse, reader: jspb.BinaryReader): ListUserMetadataResponse;
}

export namespace ListUserMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_metadata_pb.Metadata.AsObject>,
  }
}

export class GetUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetUserMetadataRequest;

  getKey(): string;
  setKey(value: string): GetUserMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserMetadataRequest): GetUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserMetadataRequest;
  static deserializeBinaryFromReader(message: GetUserMetadataRequest, reader: jspb.BinaryReader): GetUserMetadataRequest;
}

export namespace GetUserMetadataRequest {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class GetUserMetadataResponse extends jspb.Message {
  getMetadata(): zitadel_metadata_pb.Metadata | undefined;
  setMetadata(value?: zitadel_metadata_pb.Metadata): GetUserMetadataResponse;
  hasMetadata(): boolean;
  clearMetadata(): GetUserMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserMetadataResponse): GetUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserMetadataResponse;
  static deserializeBinaryFromReader(message: GetUserMetadataResponse, reader: jspb.BinaryReader): GetUserMetadataResponse;
}

export namespace GetUserMetadataResponse {
  export type AsObject = {
    metadata?: zitadel_metadata_pb.Metadata.AsObject,
  }
}

export class SetUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): SetUserMetadataRequest;

  getKey(): string;
  setKey(value: string): SetUserMetadataRequest;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): SetUserMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserMetadataRequest): SetUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: SetUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserMetadataRequest;
  static deserializeBinaryFromReader(message: SetUserMetadataRequest, reader: jspb.BinaryReader): SetUserMetadataRequest;
}

export namespace SetUserMetadataRequest {
  export type AsObject = {
    id: string,
    key: string,
    value: Uint8Array | string,
  }
}

export class SetUserMetadataResponse extends jspb.Message {
  getId(): string;
  setId(value: string): SetUserMetadataResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetUserMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): SetUserMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetUserMetadataResponse): SetUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: SetUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetUserMetadataResponse;
  static deserializeBinaryFromReader(message: SetUserMetadataResponse, reader: jspb.BinaryReader): SetUserMetadataResponse;
}

export namespace SetUserMetadataResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkSetUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): BulkSetUserMetadataRequest;

  getMetadataList(): Array<BulkSetUserMetadataRequest.Metadata>;
  setMetadataList(value: Array<BulkSetUserMetadataRequest.Metadata>): BulkSetUserMetadataRequest;
  clearMetadataList(): BulkSetUserMetadataRequest;
  addMetadata(value?: BulkSetUserMetadataRequest.Metadata, index?: number): BulkSetUserMetadataRequest.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetUserMetadataRequest): BulkSetUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkSetUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetUserMetadataRequest;
  static deserializeBinaryFromReader(message: BulkSetUserMetadataRequest, reader: jspb.BinaryReader): BulkSetUserMetadataRequest;
}

export namespace BulkSetUserMetadataRequest {
  export type AsObject = {
    id: string,
    metadataList: Array<BulkSetUserMetadataRequest.Metadata.AsObject>,
  }

  export class Metadata extends jspb.Message {
    getKey(): string;
    setKey(value: string): Metadata;

    getValue(): Uint8Array | string;
    getValue_asU8(): Uint8Array;
    getValue_asB64(): string;
    setValue(value: Uint8Array | string): Metadata;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Metadata.AsObject;
    static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
    static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Metadata;
    static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
  }

  export namespace Metadata {
    export type AsObject = {
      key: string,
      value: Uint8Array | string,
    }
  }

}

export class BulkSetUserMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkSetUserMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkSetUserMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetUserMetadataResponse): BulkSetUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkSetUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetUserMetadataResponse;
  static deserializeBinaryFromReader(message: BulkSetUserMetadataResponse, reader: jspb.BinaryReader): BulkSetUserMetadataResponse;
}

export namespace BulkSetUserMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RemoveUserMetadataRequest;

  getKey(): string;
  setKey(value: string): RemoveUserMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserMetadataRequest): RemoveUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserMetadataRequest;
  static deserializeBinaryFromReader(message: RemoveUserMetadataRequest, reader: jspb.BinaryReader): RemoveUserMetadataRequest;
}

export namespace RemoveUserMetadataRequest {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class RemoveUserMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveUserMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveUserMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserMetadataResponse): RemoveUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserMetadataResponse;
  static deserializeBinaryFromReader(message: RemoveUserMetadataResponse, reader: jspb.BinaryReader): RemoveUserMetadataResponse;
}

export namespace RemoveUserMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkRemoveUserMetadataRequest extends jspb.Message {
  getId(): string;
  setId(value: string): BulkRemoveUserMetadataRequest;

  getKeysList(): Array<string>;
  setKeysList(value: Array<string>): BulkRemoveUserMetadataRequest;
  clearKeysList(): BulkRemoveUserMetadataRequest;
  addKeys(value: string, index?: number): BulkRemoveUserMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveUserMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveUserMetadataRequest): BulkRemoveUserMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveUserMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveUserMetadataRequest;
  static deserializeBinaryFromReader(message: BulkRemoveUserMetadataRequest, reader: jspb.BinaryReader): BulkRemoveUserMetadataRequest;
}

export namespace BulkRemoveUserMetadataRequest {
  export type AsObject = {
    id: string,
    keysList: Array<string>,
  }
}

export class BulkRemoveUserMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkRemoveUserMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkRemoveUserMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveUserMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveUserMetadataResponse): BulkRemoveUserMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveUserMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveUserMetadataResponse;
  static deserializeBinaryFromReader(message: BulkRemoveUserMetadataResponse, reader: jspb.BinaryReader): BulkRemoveUserMetadataResponse;
}

export namespace BulkRemoveUserMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetHumanProfileRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetHumanProfileRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanProfileRequest): GetHumanProfileRequest.AsObject;
  static serializeBinaryToWriter(message: GetHumanProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanProfileRequest;
  static deserializeBinaryFromReader(message: GetHumanProfileRequest, reader: jspb.BinaryReader): GetHumanProfileRequest;
}

export namespace GetHumanProfileRequest {
  export type AsObject = {
    userId: string,
  }
}

export class GetHumanProfileResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetHumanProfileResponse;
  hasDetails(): boolean;
  clearDetails(): GetHumanProfileResponse;

  getProfile(): zitadel_user_pb.Profile | undefined;
  setProfile(value?: zitadel_user_pb.Profile): GetHumanProfileResponse;
  hasProfile(): boolean;
  clearProfile(): GetHumanProfileResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanProfileResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanProfileResponse): GetHumanProfileResponse.AsObject;
  static serializeBinaryToWriter(message: GetHumanProfileResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanProfileResponse;
  static deserializeBinaryFromReader(message: GetHumanProfileResponse, reader: jspb.BinaryReader): GetHumanProfileResponse;
}

export namespace GetHumanProfileResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    profile?: zitadel_user_pb.Profile.AsObject,
  }
}

export class UpdateHumanProfileRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateHumanProfileRequest;

  getFirstName(): string;
  setFirstName(value: string): UpdateHumanProfileRequest;

  getLastName(): string;
  setLastName(value: string): UpdateHumanProfileRequest;

  getNickName(): string;
  setNickName(value: string): UpdateHumanProfileRequest;

  getDisplayName(): string;
  setDisplayName(value: string): UpdateHumanProfileRequest;

  getPreferredLanguage(): string;
  setPreferredLanguage(value: string): UpdateHumanProfileRequest;

  getGender(): zitadel_user_pb.Gender;
  setGender(value: zitadel_user_pb.Gender): UpdateHumanProfileRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanProfileRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanProfileRequest): UpdateHumanProfileRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanProfileRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanProfileRequest;
  static deserializeBinaryFromReader(message: UpdateHumanProfileRequest, reader: jspb.BinaryReader): UpdateHumanProfileRequest;
}

export namespace UpdateHumanProfileRequest {
  export type AsObject = {
    userId: string,
    firstName: string,
    lastName: string,
    nickName: string,
    displayName: string,
    preferredLanguage: string,
    gender: zitadel_user_pb.Gender,
  }
}

export class UpdateHumanProfileResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateHumanProfileResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateHumanProfileResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanProfileResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanProfileResponse): UpdateHumanProfileResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanProfileResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanProfileResponse;
  static deserializeBinaryFromReader(message: UpdateHumanProfileResponse, reader: jspb.BinaryReader): UpdateHumanProfileResponse;
}

export namespace UpdateHumanProfileResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetHumanEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetHumanEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanEmailRequest): GetHumanEmailRequest.AsObject;
  static serializeBinaryToWriter(message: GetHumanEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanEmailRequest;
  static deserializeBinaryFromReader(message: GetHumanEmailRequest, reader: jspb.BinaryReader): GetHumanEmailRequest;
}

export namespace GetHumanEmailRequest {
  export type AsObject = {
    userId: string,
  }
}

export class GetHumanEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetHumanEmailResponse;
  hasDetails(): boolean;
  clearDetails(): GetHumanEmailResponse;

  getEmail(): zitadel_user_pb.Email | undefined;
  setEmail(value?: zitadel_user_pb.Email): GetHumanEmailResponse;
  hasEmail(): boolean;
  clearEmail(): GetHumanEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanEmailResponse): GetHumanEmailResponse.AsObject;
  static serializeBinaryToWriter(message: GetHumanEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanEmailResponse;
  static deserializeBinaryFromReader(message: GetHumanEmailResponse, reader: jspb.BinaryReader): GetHumanEmailResponse;
}

export namespace GetHumanEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    email?: zitadel_user_pb.Email.AsObject,
  }
}

export class UpdateHumanEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateHumanEmailRequest;

  getEmail(): string;
  setEmail(value: string): UpdateHumanEmailRequest;

  getIsEmailVerified(): boolean;
  setIsEmailVerified(value: boolean): UpdateHumanEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanEmailRequest): UpdateHumanEmailRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanEmailRequest;
  static deserializeBinaryFromReader(message: UpdateHumanEmailRequest, reader: jspb.BinaryReader): UpdateHumanEmailRequest;
}

export namespace UpdateHumanEmailRequest {
  export type AsObject = {
    userId: string,
    email: string,
    isEmailVerified: boolean,
  }
}

export class UpdateHumanEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateHumanEmailResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateHumanEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanEmailResponse): UpdateHumanEmailResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanEmailResponse;
  static deserializeBinaryFromReader(message: UpdateHumanEmailResponse, reader: jspb.BinaryReader): UpdateHumanEmailResponse;
}

export namespace UpdateHumanEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResendHumanInitializationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendHumanInitializationRequest;

  getEmail(): string;
  setEmail(value: string): ResendHumanInitializationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanInitializationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanInitializationRequest): ResendHumanInitializationRequest.AsObject;
  static serializeBinaryToWriter(message: ResendHumanInitializationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanInitializationRequest;
  static deserializeBinaryFromReader(message: ResendHumanInitializationRequest, reader: jspb.BinaryReader): ResendHumanInitializationRequest;
}

export namespace ResendHumanInitializationRequest {
  export type AsObject = {
    userId: string,
    email: string,
  }
}

export class ResendHumanInitializationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResendHumanInitializationResponse;
  hasDetails(): boolean;
  clearDetails(): ResendHumanInitializationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanInitializationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanInitializationResponse): ResendHumanInitializationResponse.AsObject;
  static serializeBinaryToWriter(message: ResendHumanInitializationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanInitializationResponse;
  static deserializeBinaryFromReader(message: ResendHumanInitializationResponse, reader: jspb.BinaryReader): ResendHumanInitializationResponse;
}

export namespace ResendHumanInitializationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResendHumanEmailVerificationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendHumanEmailVerificationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanEmailVerificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanEmailVerificationRequest): ResendHumanEmailVerificationRequest.AsObject;
  static serializeBinaryToWriter(message: ResendHumanEmailVerificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanEmailVerificationRequest;
  static deserializeBinaryFromReader(message: ResendHumanEmailVerificationRequest, reader: jspb.BinaryReader): ResendHumanEmailVerificationRequest;
}

export namespace ResendHumanEmailVerificationRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ResendHumanEmailVerificationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResendHumanEmailVerificationResponse;
  hasDetails(): boolean;
  clearDetails(): ResendHumanEmailVerificationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanEmailVerificationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanEmailVerificationResponse): ResendHumanEmailVerificationResponse.AsObject;
  static serializeBinaryToWriter(message: ResendHumanEmailVerificationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanEmailVerificationResponse;
  static deserializeBinaryFromReader(message: ResendHumanEmailVerificationResponse, reader: jspb.BinaryReader): ResendHumanEmailVerificationResponse;
}

export namespace ResendHumanEmailVerificationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetHumanPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetHumanPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanPhoneRequest): GetHumanPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: GetHumanPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanPhoneRequest;
  static deserializeBinaryFromReader(message: GetHumanPhoneRequest, reader: jspb.BinaryReader): GetHumanPhoneRequest;
}

export namespace GetHumanPhoneRequest {
  export type AsObject = {
    userId: string,
  }
}

export class GetHumanPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GetHumanPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): GetHumanPhoneResponse;

  getPhone(): zitadel_user_pb.Phone | undefined;
  setPhone(value?: zitadel_user_pb.Phone): GetHumanPhoneResponse;
  hasPhone(): boolean;
  clearPhone(): GetHumanPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHumanPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetHumanPhoneResponse): GetHumanPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: GetHumanPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHumanPhoneResponse;
  static deserializeBinaryFromReader(message: GetHumanPhoneResponse, reader: jspb.BinaryReader): GetHumanPhoneResponse;
}

export namespace GetHumanPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    phone?: zitadel_user_pb.Phone.AsObject,
  }
}

export class UpdateHumanPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateHumanPhoneRequest;

  getPhone(): string;
  setPhone(value: string): UpdateHumanPhoneRequest;

  getIsPhoneVerified(): boolean;
  setIsPhoneVerified(value: boolean): UpdateHumanPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanPhoneRequest): UpdateHumanPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanPhoneRequest;
  static deserializeBinaryFromReader(message: UpdateHumanPhoneRequest, reader: jspb.BinaryReader): UpdateHumanPhoneRequest;
}

export namespace UpdateHumanPhoneRequest {
  export type AsObject = {
    userId: string,
    phone: string,
    isPhoneVerified: boolean,
  }
}

export class UpdateHumanPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateHumanPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateHumanPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateHumanPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateHumanPhoneResponse): UpdateHumanPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateHumanPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateHumanPhoneResponse;
  static deserializeBinaryFromReader(message: UpdateHumanPhoneResponse, reader: jspb.BinaryReader): UpdateHumanPhoneResponse;
}

export namespace UpdateHumanPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanPhoneRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanPhoneRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanPhoneRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanPhoneRequest): RemoveHumanPhoneRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanPhoneRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanPhoneRequest;
  static deserializeBinaryFromReader(message: RemoveHumanPhoneRequest, reader: jspb.BinaryReader): RemoveHumanPhoneRequest;
}

export namespace RemoveHumanPhoneRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveHumanPhoneResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanPhoneResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanPhoneResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanPhoneResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanPhoneResponse): RemoveHumanPhoneResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanPhoneResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanPhoneResponse;
  static deserializeBinaryFromReader(message: RemoveHumanPhoneResponse, reader: jspb.BinaryReader): RemoveHumanPhoneResponse;
}

export namespace RemoveHumanPhoneResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResendHumanPhoneVerificationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ResendHumanPhoneVerificationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanPhoneVerificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanPhoneVerificationRequest): ResendHumanPhoneVerificationRequest.AsObject;
  static serializeBinaryToWriter(message: ResendHumanPhoneVerificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanPhoneVerificationRequest;
  static deserializeBinaryFromReader(message: ResendHumanPhoneVerificationRequest, reader: jspb.BinaryReader): ResendHumanPhoneVerificationRequest;
}

export namespace ResendHumanPhoneVerificationRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ResendHumanPhoneVerificationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResendHumanPhoneVerificationResponse;
  hasDetails(): boolean;
  clearDetails(): ResendHumanPhoneVerificationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResendHumanPhoneVerificationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResendHumanPhoneVerificationResponse): ResendHumanPhoneVerificationResponse.AsObject;
  static serializeBinaryToWriter(message: ResendHumanPhoneVerificationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResendHumanPhoneVerificationResponse;
  static deserializeBinaryFromReader(message: ResendHumanPhoneVerificationResponse, reader: jspb.BinaryReader): ResendHumanPhoneVerificationResponse;
}

export namespace ResendHumanPhoneVerificationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanAvatarRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanAvatarRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAvatarRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAvatarRequest): RemoveHumanAvatarRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAvatarRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAvatarRequest;
  static deserializeBinaryFromReader(message: RemoveHumanAvatarRequest, reader: jspb.BinaryReader): RemoveHumanAvatarRequest;
}

export namespace RemoveHumanAvatarRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveHumanAvatarResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanAvatarResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanAvatarResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAvatarResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAvatarResponse): RemoveHumanAvatarResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAvatarResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAvatarResponse;
  static deserializeBinaryFromReader(message: RemoveHumanAvatarResponse, reader: jspb.BinaryReader): RemoveHumanAvatarResponse;
}

export namespace RemoveHumanAvatarResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetHumanInitialPasswordRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetHumanInitialPasswordRequest;

  getPassword(): string;
  setPassword(value: string): SetHumanInitialPasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanInitialPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanInitialPasswordRequest): SetHumanInitialPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: SetHumanInitialPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanInitialPasswordRequest;
  static deserializeBinaryFromReader(message: SetHumanInitialPasswordRequest, reader: jspb.BinaryReader): SetHumanInitialPasswordRequest;
}

export namespace SetHumanInitialPasswordRequest {
  export type AsObject = {
    userId: string,
    password: string,
  }
}

export class SetHumanInitialPasswordResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetHumanInitialPasswordResponse;
  hasDetails(): boolean;
  clearDetails(): SetHumanInitialPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanInitialPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanInitialPasswordResponse): SetHumanInitialPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: SetHumanInitialPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanInitialPasswordResponse;
  static deserializeBinaryFromReader(message: SetHumanInitialPasswordResponse, reader: jspb.BinaryReader): SetHumanInitialPasswordResponse;
}

export namespace SetHumanInitialPasswordResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetHumanPasswordRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SetHumanPasswordRequest;

  getPassword(): string;
  setPassword(value: string): SetHumanPasswordRequest;

  getNoChangeRequired(): boolean;
  setNoChangeRequired(value: boolean): SetHumanPasswordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanPasswordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanPasswordRequest): SetHumanPasswordRequest.AsObject;
  static serializeBinaryToWriter(message: SetHumanPasswordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanPasswordRequest;
  static deserializeBinaryFromReader(message: SetHumanPasswordRequest, reader: jspb.BinaryReader): SetHumanPasswordRequest;
}

export namespace SetHumanPasswordRequest {
  export type AsObject = {
    userId: string,
    password: string,
    noChangeRequired: boolean,
  }
}

export class SetHumanPasswordResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetHumanPasswordResponse;
  hasDetails(): boolean;
  clearDetails(): SetHumanPasswordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetHumanPasswordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetHumanPasswordResponse): SetHumanPasswordResponse.AsObject;
  static serializeBinaryToWriter(message: SetHumanPasswordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetHumanPasswordResponse;
  static deserializeBinaryFromReader(message: SetHumanPasswordResponse, reader: jspb.BinaryReader): SetHumanPasswordResponse;
}

export namespace SetHumanPasswordResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SendHumanResetPasswordNotificationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SendHumanResetPasswordNotificationRequest;

  getType(): SendHumanResetPasswordNotificationRequest.Type;
  setType(value: SendHumanResetPasswordNotificationRequest.Type): SendHumanResetPasswordNotificationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendHumanResetPasswordNotificationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendHumanResetPasswordNotificationRequest): SendHumanResetPasswordNotificationRequest.AsObject;
  static serializeBinaryToWriter(message: SendHumanResetPasswordNotificationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendHumanResetPasswordNotificationRequest;
  static deserializeBinaryFromReader(message: SendHumanResetPasswordNotificationRequest, reader: jspb.BinaryReader): SendHumanResetPasswordNotificationRequest;
}

export namespace SendHumanResetPasswordNotificationRequest {
  export type AsObject = {
    userId: string,
    type: SendHumanResetPasswordNotificationRequest.Type,
  }

  export enum Type { 
    TYPE_EMAIL = 0,
    TYPE_SMS = 1,
  }
}

export class SendHumanResetPasswordNotificationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SendHumanResetPasswordNotificationResponse;
  hasDetails(): boolean;
  clearDetails(): SendHumanResetPasswordNotificationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendHumanResetPasswordNotificationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SendHumanResetPasswordNotificationResponse): SendHumanResetPasswordNotificationResponse.AsObject;
  static serializeBinaryToWriter(message: SendHumanResetPasswordNotificationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendHumanResetPasswordNotificationResponse;
  static deserializeBinaryFromReader(message: SendHumanResetPasswordNotificationResponse, reader: jspb.BinaryReader): SendHumanResetPasswordNotificationResponse;
}

export namespace SendHumanResetPasswordNotificationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListHumanAuthFactorsRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListHumanAuthFactorsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanAuthFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanAuthFactorsRequest): ListHumanAuthFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListHumanAuthFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanAuthFactorsRequest;
  static deserializeBinaryFromReader(message: ListHumanAuthFactorsRequest, reader: jspb.BinaryReader): ListHumanAuthFactorsRequest;
}

export namespace ListHumanAuthFactorsRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ListHumanAuthFactorsResponse extends jspb.Message {
  getResultList(): Array<zitadel_user_pb.AuthFactor>;
  setResultList(value: Array<zitadel_user_pb.AuthFactor>): ListHumanAuthFactorsResponse;
  clearResultList(): ListHumanAuthFactorsResponse;
  addResult(value?: zitadel_user_pb.AuthFactor, index?: number): zitadel_user_pb.AuthFactor;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanAuthFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanAuthFactorsResponse): ListHumanAuthFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListHumanAuthFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanAuthFactorsResponse;
  static deserializeBinaryFromReader(message: ListHumanAuthFactorsResponse, reader: jspb.BinaryReader): ListHumanAuthFactorsResponse;
}

export namespace ListHumanAuthFactorsResponse {
  export type AsObject = {
    resultList: Array<zitadel_user_pb.AuthFactor.AsObject>,
  }
}

export class RemoveHumanAuthFactorOTPRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanAuthFactorOTPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPRequest): RemoveHumanAuthFactorOTPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPRequest;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPRequest, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPRequest;
}

export namespace RemoveHumanAuthFactorOTPRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveHumanAuthFactorOTPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanAuthFactorOTPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanAuthFactorOTPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPResponse): RemoveHumanAuthFactorOTPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPResponse;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPResponse, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPResponse;
}

export namespace RemoveHumanAuthFactorOTPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanAuthFactorU2FRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanAuthFactorU2FRequest;

  getTokenId(): string;
  setTokenId(value: string): RemoveHumanAuthFactorU2FRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorU2FRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorU2FRequest): RemoveHumanAuthFactorU2FRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorU2FRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorU2FRequest;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorU2FRequest, reader: jspb.BinaryReader): RemoveHumanAuthFactorU2FRequest;
}

export namespace RemoveHumanAuthFactorU2FRequest {
  export type AsObject = {
    userId: string,
    tokenId: string,
  }
}

export class RemoveHumanAuthFactorU2FResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanAuthFactorU2FResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanAuthFactorU2FResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorU2FResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorU2FResponse): RemoveHumanAuthFactorU2FResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorU2FResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorU2FResponse;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorU2FResponse, reader: jspb.BinaryReader): RemoveHumanAuthFactorU2FResponse;
}

export namespace RemoveHumanAuthFactorU2FResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanAuthFactorOTPSMSRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanAuthFactorOTPSMSRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPSMSRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPSMSRequest): RemoveHumanAuthFactorOTPSMSRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPSMSRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPSMSRequest;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPSMSRequest, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPSMSRequest;
}

export namespace RemoveHumanAuthFactorOTPSMSRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveHumanAuthFactorOTPSMSResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanAuthFactorOTPSMSResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanAuthFactorOTPSMSResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPSMSResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPSMSResponse): RemoveHumanAuthFactorOTPSMSResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPSMSResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPSMSResponse;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPSMSResponse, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPSMSResponse;
}

export namespace RemoveHumanAuthFactorOTPSMSResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanAuthFactorOTPEmailRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanAuthFactorOTPEmailRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPEmailRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPEmailRequest): RemoveHumanAuthFactorOTPEmailRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPEmailRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPEmailRequest;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPEmailRequest, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPEmailRequest;
}

export namespace RemoveHumanAuthFactorOTPEmailRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveHumanAuthFactorOTPEmailResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanAuthFactorOTPEmailResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanAuthFactorOTPEmailResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanAuthFactorOTPEmailResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanAuthFactorOTPEmailResponse): RemoveHumanAuthFactorOTPEmailResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanAuthFactorOTPEmailResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanAuthFactorOTPEmailResponse;
  static deserializeBinaryFromReader(message: RemoveHumanAuthFactorOTPEmailResponse, reader: jspb.BinaryReader): RemoveHumanAuthFactorOTPEmailResponse;
}

export namespace RemoveHumanAuthFactorOTPEmailResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListHumanPasswordlessRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListHumanPasswordlessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanPasswordlessRequest): ListHumanPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: ListHumanPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanPasswordlessRequest;
  static deserializeBinaryFromReader(message: ListHumanPasswordlessRequest, reader: jspb.BinaryReader): ListHumanPasswordlessRequest;
}

export namespace ListHumanPasswordlessRequest {
  export type AsObject = {
    userId: string,
  }
}

export class ListHumanPasswordlessResponse extends jspb.Message {
  getResultList(): Array<zitadel_user_pb.WebAuthNToken>;
  setResultList(value: Array<zitadel_user_pb.WebAuthNToken>): ListHumanPasswordlessResponse;
  clearResultList(): ListHumanPasswordlessResponse;
  addResult(value?: zitadel_user_pb.WebAuthNToken, index?: number): zitadel_user_pb.WebAuthNToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanPasswordlessResponse): ListHumanPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: ListHumanPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanPasswordlessResponse;
  static deserializeBinaryFromReader(message: ListHumanPasswordlessResponse, reader: jspb.BinaryReader): ListHumanPasswordlessResponse;
}

export namespace ListHumanPasswordlessResponse {
  export type AsObject = {
    resultList: Array<zitadel_user_pb.WebAuthNToken.AsObject>,
  }
}

export class AddPasswordlessRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddPasswordlessRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPasswordlessRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddPasswordlessRegistrationRequest): AddPasswordlessRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: AddPasswordlessRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPasswordlessRegistrationRequest;
  static deserializeBinaryFromReader(message: AddPasswordlessRegistrationRequest, reader: jspb.BinaryReader): AddPasswordlessRegistrationRequest;
}

export namespace AddPasswordlessRegistrationRequest {
  export type AsObject = {
    userId: string,
  }
}

export class AddPasswordlessRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddPasswordlessRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): AddPasswordlessRegistrationResponse;

  getLink(): string;
  setLink(value: string): AddPasswordlessRegistrationResponse;

  getExpiration(): google_protobuf_duration_pb.Duration | undefined;
  setExpiration(value?: google_protobuf_duration_pb.Duration): AddPasswordlessRegistrationResponse;
  hasExpiration(): boolean;
  clearExpiration(): AddPasswordlessRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPasswordlessRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddPasswordlessRegistrationResponse): AddPasswordlessRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: AddPasswordlessRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPasswordlessRegistrationResponse;
  static deserializeBinaryFromReader(message: AddPasswordlessRegistrationResponse, reader: jspb.BinaryReader): AddPasswordlessRegistrationResponse;
}

export namespace AddPasswordlessRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    link: string,
    expiration?: google_protobuf_duration_pb.Duration.AsObject,
  }
}

export class SendPasswordlessRegistrationRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): SendPasswordlessRegistrationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasswordlessRegistrationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasswordlessRegistrationRequest): SendPasswordlessRegistrationRequest.AsObject;
  static serializeBinaryToWriter(message: SendPasswordlessRegistrationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasswordlessRegistrationRequest;
  static deserializeBinaryFromReader(message: SendPasswordlessRegistrationRequest, reader: jspb.BinaryReader): SendPasswordlessRegistrationRequest;
}

export namespace SendPasswordlessRegistrationRequest {
  export type AsObject = {
    userId: string,
  }
}

export class SendPasswordlessRegistrationResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SendPasswordlessRegistrationResponse;
  hasDetails(): boolean;
  clearDetails(): SendPasswordlessRegistrationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendPasswordlessRegistrationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SendPasswordlessRegistrationResponse): SendPasswordlessRegistrationResponse.AsObject;
  static serializeBinaryToWriter(message: SendPasswordlessRegistrationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendPasswordlessRegistrationResponse;
  static deserializeBinaryFromReader(message: SendPasswordlessRegistrationResponse, reader: jspb.BinaryReader): SendPasswordlessRegistrationResponse;
}

export namespace SendPasswordlessRegistrationResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveHumanPasswordlessRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanPasswordlessRequest;

  getTokenId(): string;
  setTokenId(value: string): RemoveHumanPasswordlessRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanPasswordlessRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanPasswordlessRequest): RemoveHumanPasswordlessRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanPasswordlessRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanPasswordlessRequest;
  static deserializeBinaryFromReader(message: RemoveHumanPasswordlessRequest, reader: jspb.BinaryReader): RemoveHumanPasswordlessRequest;
}

export namespace RemoveHumanPasswordlessRequest {
  export type AsObject = {
    userId: string,
    tokenId: string,
  }
}

export class RemoveHumanPasswordlessResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanPasswordlessResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanPasswordlessResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanPasswordlessResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanPasswordlessResponse): RemoveHumanPasswordlessResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanPasswordlessResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanPasswordlessResponse;
  static deserializeBinaryFromReader(message: RemoveHumanPasswordlessResponse, reader: jspb.BinaryReader): RemoveHumanPasswordlessResponse;
}

export namespace RemoveHumanPasswordlessResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateMachineRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateMachineRequest;

  getDescription(): string;
  setDescription(value: string): UpdateMachineRequest;

  getName(): string;
  setName(value: string): UpdateMachineRequest;

  getAccessTokenType(): zitadel_user_pb.AccessTokenType;
  setAccessTokenType(value: zitadel_user_pb.AccessTokenType): UpdateMachineRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMachineRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMachineRequest): UpdateMachineRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateMachineRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMachineRequest;
  static deserializeBinaryFromReader(message: UpdateMachineRequest, reader: jspb.BinaryReader): UpdateMachineRequest;
}

export namespace UpdateMachineRequest {
  export type AsObject = {
    userId: string,
    description: string,
    name: string,
    accessTokenType: zitadel_user_pb.AccessTokenType,
  }
}

export class UpdateMachineResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateMachineResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateMachineResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateMachineResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateMachineResponse): UpdateMachineResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateMachineResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateMachineResponse;
  static deserializeBinaryFromReader(message: UpdateMachineResponse, reader: jspb.BinaryReader): UpdateMachineResponse;
}

export namespace UpdateMachineResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GenerateMachineSecretRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GenerateMachineSecretRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateMachineSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateMachineSecretRequest): GenerateMachineSecretRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateMachineSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateMachineSecretRequest;
  static deserializeBinaryFromReader(message: GenerateMachineSecretRequest, reader: jspb.BinaryReader): GenerateMachineSecretRequest;
}

export namespace GenerateMachineSecretRequest {
  export type AsObject = {
    userId: string,
  }
}

export class GenerateMachineSecretResponse extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): GenerateMachineSecretResponse;

  getClientSecret(): string;
  setClientSecret(value: string): GenerateMachineSecretResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GenerateMachineSecretResponse;
  hasDetails(): boolean;
  clearDetails(): GenerateMachineSecretResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateMachineSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateMachineSecretResponse): GenerateMachineSecretResponse.AsObject;
  static serializeBinaryToWriter(message: GenerateMachineSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateMachineSecretResponse;
  static deserializeBinaryFromReader(message: GenerateMachineSecretResponse, reader: jspb.BinaryReader): GenerateMachineSecretResponse;
}

export namespace GenerateMachineSecretResponse {
  export type AsObject = {
    clientId: string,
    clientSecret: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMachineSecretRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveMachineSecretRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMachineSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMachineSecretRequest): RemoveMachineSecretRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMachineSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMachineSecretRequest;
  static deserializeBinaryFromReader(message: RemoveMachineSecretRequest, reader: jspb.BinaryReader): RemoveMachineSecretRequest;
}

export namespace RemoveMachineSecretRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveMachineSecretResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMachineSecretResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMachineSecretResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMachineSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMachineSecretResponse): RemoveMachineSecretResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMachineSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMachineSecretResponse;
  static deserializeBinaryFromReader(message: RemoveMachineSecretResponse, reader: jspb.BinaryReader): RemoveMachineSecretResponse;
}

export namespace RemoveMachineSecretResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetMachineKeyByIDsRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetMachineKeyByIDsRequest;

  getKeyId(): string;
  setKeyId(value: string): GetMachineKeyByIDsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMachineKeyByIDsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMachineKeyByIDsRequest): GetMachineKeyByIDsRequest.AsObject;
  static serializeBinaryToWriter(message: GetMachineKeyByIDsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMachineKeyByIDsRequest;
  static deserializeBinaryFromReader(message: GetMachineKeyByIDsRequest, reader: jspb.BinaryReader): GetMachineKeyByIDsRequest;
}

export namespace GetMachineKeyByIDsRequest {
  export type AsObject = {
    userId: string,
    keyId: string,
  }
}

export class GetMachineKeyByIDsResponse extends jspb.Message {
  getKey(): zitadel_auth_n_key_pb.Key | undefined;
  setKey(value?: zitadel_auth_n_key_pb.Key): GetMachineKeyByIDsResponse;
  hasKey(): boolean;
  clearKey(): GetMachineKeyByIDsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMachineKeyByIDsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMachineKeyByIDsResponse): GetMachineKeyByIDsResponse.AsObject;
  static serializeBinaryToWriter(message: GetMachineKeyByIDsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMachineKeyByIDsResponse;
  static deserializeBinaryFromReader(message: GetMachineKeyByIDsResponse, reader: jspb.BinaryReader): GetMachineKeyByIDsResponse;
}

export namespace GetMachineKeyByIDsResponse {
  export type AsObject = {
    key?: zitadel_auth_n_key_pb.Key.AsObject,
  }
}

export class ListMachineKeysRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListMachineKeysRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListMachineKeysRequest;
  hasQuery(): boolean;
  clearQuery(): ListMachineKeysRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMachineKeysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMachineKeysRequest): ListMachineKeysRequest.AsObject;
  static serializeBinaryToWriter(message: ListMachineKeysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMachineKeysRequest;
  static deserializeBinaryFromReader(message: ListMachineKeysRequest, reader: jspb.BinaryReader): ListMachineKeysRequest;
}

export namespace ListMachineKeysRequest {
  export type AsObject = {
    userId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListMachineKeysResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListMachineKeysResponse;
  hasDetails(): boolean;
  clearDetails(): ListMachineKeysResponse;

  getResultList(): Array<zitadel_auth_n_key_pb.Key>;
  setResultList(value: Array<zitadel_auth_n_key_pb.Key>): ListMachineKeysResponse;
  clearResultList(): ListMachineKeysResponse;
  addResult(value?: zitadel_auth_n_key_pb.Key, index?: number): zitadel_auth_n_key_pb.Key;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMachineKeysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListMachineKeysResponse): ListMachineKeysResponse.AsObject;
  static serializeBinaryToWriter(message: ListMachineKeysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMachineKeysResponse;
  static deserializeBinaryFromReader(message: ListMachineKeysResponse, reader: jspb.BinaryReader): ListMachineKeysResponse;
}

export namespace ListMachineKeysResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_auth_n_key_pb.Key.AsObject>,
  }
}

export class AddMachineKeyRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddMachineKeyRequest;

  getType(): zitadel_auth_n_key_pb.KeyType;
  setType(value: zitadel_auth_n_key_pb.KeyType): AddMachineKeyRequest;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): AddMachineKeyRequest;
  hasExpirationDate(): boolean;
  clearExpirationDate(): AddMachineKeyRequest;

  getPublicKey(): Uint8Array | string;
  getPublicKey_asU8(): Uint8Array;
  getPublicKey_asB64(): string;
  setPublicKey(value: Uint8Array | string): AddMachineKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMachineKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMachineKeyRequest): AddMachineKeyRequest.AsObject;
  static serializeBinaryToWriter(message: AddMachineKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMachineKeyRequest;
  static deserializeBinaryFromReader(message: AddMachineKeyRequest, reader: jspb.BinaryReader): AddMachineKeyRequest;
}

export namespace AddMachineKeyRequest {
  export type AsObject = {
    userId: string,
    type: zitadel_auth_n_key_pb.KeyType,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    publicKey: Uint8Array | string,
  }
}

export class AddMachineKeyResponse extends jspb.Message {
  getKeyId(): string;
  setKeyId(value: string): AddMachineKeyResponse;

  getKeyDetails(): Uint8Array | string;
  getKeyDetails_asU8(): Uint8Array;
  getKeyDetails_asB64(): string;
  setKeyDetails(value: Uint8Array | string): AddMachineKeyResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMachineKeyResponse;
  hasDetails(): boolean;
  clearDetails(): AddMachineKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMachineKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMachineKeyResponse): AddMachineKeyResponse.AsObject;
  static serializeBinaryToWriter(message: AddMachineKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMachineKeyResponse;
  static deserializeBinaryFromReader(message: AddMachineKeyResponse, reader: jspb.BinaryReader): AddMachineKeyResponse;
}

export namespace AddMachineKeyResponse {
  export type AsObject = {
    keyId: string,
    keyDetails: Uint8Array | string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMachineKeyRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveMachineKeyRequest;

  getKeyId(): string;
  setKeyId(value: string): RemoveMachineKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMachineKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMachineKeyRequest): RemoveMachineKeyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMachineKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMachineKeyRequest;
  static deserializeBinaryFromReader(message: RemoveMachineKeyRequest, reader: jspb.BinaryReader): RemoveMachineKeyRequest;
}

export namespace RemoveMachineKeyRequest {
  export type AsObject = {
    userId: string,
    keyId: string,
  }
}

export class RemoveMachineKeyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMachineKeyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMachineKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMachineKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMachineKeyResponse): RemoveMachineKeyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMachineKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMachineKeyResponse;
  static deserializeBinaryFromReader(message: RemoveMachineKeyResponse, reader: jspb.BinaryReader): RemoveMachineKeyResponse;
}

export namespace RemoveMachineKeyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPersonalAccessTokenByIDsRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetPersonalAccessTokenByIDsRequest;

  getTokenId(): string;
  setTokenId(value: string): GetPersonalAccessTokenByIDsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPersonalAccessTokenByIDsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPersonalAccessTokenByIDsRequest): GetPersonalAccessTokenByIDsRequest.AsObject;
  static serializeBinaryToWriter(message: GetPersonalAccessTokenByIDsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPersonalAccessTokenByIDsRequest;
  static deserializeBinaryFromReader(message: GetPersonalAccessTokenByIDsRequest, reader: jspb.BinaryReader): GetPersonalAccessTokenByIDsRequest;
}

export namespace GetPersonalAccessTokenByIDsRequest {
  export type AsObject = {
    userId: string,
    tokenId: string,
  }
}

export class GetPersonalAccessTokenByIDsResponse extends jspb.Message {
  getToken(): zitadel_user_pb.PersonalAccessToken | undefined;
  setToken(value?: zitadel_user_pb.PersonalAccessToken): GetPersonalAccessTokenByIDsResponse;
  hasToken(): boolean;
  clearToken(): GetPersonalAccessTokenByIDsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPersonalAccessTokenByIDsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPersonalAccessTokenByIDsResponse): GetPersonalAccessTokenByIDsResponse.AsObject;
  static serializeBinaryToWriter(message: GetPersonalAccessTokenByIDsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPersonalAccessTokenByIDsResponse;
  static deserializeBinaryFromReader(message: GetPersonalAccessTokenByIDsResponse, reader: jspb.BinaryReader): GetPersonalAccessTokenByIDsResponse;
}

export namespace GetPersonalAccessTokenByIDsResponse {
  export type AsObject = {
    token?: zitadel_user_pb.PersonalAccessToken.AsObject,
  }
}

export class ListPersonalAccessTokensRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListPersonalAccessTokensRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListPersonalAccessTokensRequest;
  hasQuery(): boolean;
  clearQuery(): ListPersonalAccessTokensRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPersonalAccessTokensRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPersonalAccessTokensRequest): ListPersonalAccessTokensRequest.AsObject;
  static serializeBinaryToWriter(message: ListPersonalAccessTokensRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPersonalAccessTokensRequest;
  static deserializeBinaryFromReader(message: ListPersonalAccessTokensRequest, reader: jspb.BinaryReader): ListPersonalAccessTokensRequest;
}

export namespace ListPersonalAccessTokensRequest {
  export type AsObject = {
    userId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListPersonalAccessTokensResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListPersonalAccessTokensResponse;
  hasDetails(): boolean;
  clearDetails(): ListPersonalAccessTokensResponse;

  getResultList(): Array<zitadel_user_pb.PersonalAccessToken>;
  setResultList(value: Array<zitadel_user_pb.PersonalAccessToken>): ListPersonalAccessTokensResponse;
  clearResultList(): ListPersonalAccessTokensResponse;
  addResult(value?: zitadel_user_pb.PersonalAccessToken, index?: number): zitadel_user_pb.PersonalAccessToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPersonalAccessTokensResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPersonalAccessTokensResponse): ListPersonalAccessTokensResponse.AsObject;
  static serializeBinaryToWriter(message: ListPersonalAccessTokensResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPersonalAccessTokensResponse;
  static deserializeBinaryFromReader(message: ListPersonalAccessTokensResponse, reader: jspb.BinaryReader): ListPersonalAccessTokensResponse;
}

export namespace ListPersonalAccessTokensResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_user_pb.PersonalAccessToken.AsObject>,
  }
}

export class AddPersonalAccessTokenRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddPersonalAccessTokenRequest;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): AddPersonalAccessTokenRequest;
  hasExpirationDate(): boolean;
  clearExpirationDate(): AddPersonalAccessTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPersonalAccessTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddPersonalAccessTokenRequest): AddPersonalAccessTokenRequest.AsObject;
  static serializeBinaryToWriter(message: AddPersonalAccessTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPersonalAccessTokenRequest;
  static deserializeBinaryFromReader(message: AddPersonalAccessTokenRequest, reader: jspb.BinaryReader): AddPersonalAccessTokenRequest;
}

export namespace AddPersonalAccessTokenRequest {
  export type AsObject = {
    userId: string,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AddPersonalAccessTokenResponse extends jspb.Message {
  getTokenId(): string;
  setTokenId(value: string): AddPersonalAccessTokenResponse;

  getToken(): string;
  setToken(value: string): AddPersonalAccessTokenResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddPersonalAccessTokenResponse;
  hasDetails(): boolean;
  clearDetails(): AddPersonalAccessTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPersonalAccessTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddPersonalAccessTokenResponse): AddPersonalAccessTokenResponse.AsObject;
  static serializeBinaryToWriter(message: AddPersonalAccessTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPersonalAccessTokenResponse;
  static deserializeBinaryFromReader(message: AddPersonalAccessTokenResponse, reader: jspb.BinaryReader): AddPersonalAccessTokenResponse;
}

export namespace AddPersonalAccessTokenResponse {
  export type AsObject = {
    tokenId: string,
    token: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemovePersonalAccessTokenRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemovePersonalAccessTokenRequest;

  getTokenId(): string;
  setTokenId(value: string): RemovePersonalAccessTokenRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemovePersonalAccessTokenRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemovePersonalAccessTokenRequest): RemovePersonalAccessTokenRequest.AsObject;
  static serializeBinaryToWriter(message: RemovePersonalAccessTokenRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemovePersonalAccessTokenRequest;
  static deserializeBinaryFromReader(message: RemovePersonalAccessTokenRequest, reader: jspb.BinaryReader): RemovePersonalAccessTokenRequest;
}

export namespace RemovePersonalAccessTokenRequest {
  export type AsObject = {
    userId: string,
    tokenId: string,
  }
}

export class RemovePersonalAccessTokenResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemovePersonalAccessTokenResponse;
  hasDetails(): boolean;
  clearDetails(): RemovePersonalAccessTokenResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemovePersonalAccessTokenResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemovePersonalAccessTokenResponse): RemovePersonalAccessTokenResponse.AsObject;
  static serializeBinaryToWriter(message: RemovePersonalAccessTokenResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemovePersonalAccessTokenResponse;
  static deserializeBinaryFromReader(message: RemovePersonalAccessTokenResponse, reader: jspb.BinaryReader): RemovePersonalAccessTokenResponse;
}

export namespace RemovePersonalAccessTokenResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListHumanLinkedIDPsRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListHumanLinkedIDPsRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListHumanLinkedIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListHumanLinkedIDPsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanLinkedIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanLinkedIDPsRequest): ListHumanLinkedIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListHumanLinkedIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanLinkedIDPsRequest;
  static deserializeBinaryFromReader(message: ListHumanLinkedIDPsRequest, reader: jspb.BinaryReader): ListHumanLinkedIDPsRequest;
}

export namespace ListHumanLinkedIDPsRequest {
  export type AsObject = {
    userId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListHumanLinkedIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListHumanLinkedIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListHumanLinkedIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDPUserLink>;
  setResultList(value: Array<zitadel_idp_pb.IDPUserLink>): ListHumanLinkedIDPsResponse;
  clearResultList(): ListHumanLinkedIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDPUserLink, index?: number): zitadel_idp_pb.IDPUserLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListHumanLinkedIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListHumanLinkedIDPsResponse): ListHumanLinkedIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListHumanLinkedIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListHumanLinkedIDPsResponse;
  static deserializeBinaryFromReader(message: ListHumanLinkedIDPsResponse, reader: jspb.BinaryReader): ListHumanLinkedIDPsResponse;
}

export namespace ListHumanLinkedIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.IDPUserLink.AsObject>,
  }
}

export class RemoveHumanLinkedIDPRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveHumanLinkedIDPRequest;

  getIdpId(): string;
  setIdpId(value: string): RemoveHumanLinkedIDPRequest;

  getLinkedUserId(): string;
  setLinkedUserId(value: string): RemoveHumanLinkedIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanLinkedIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanLinkedIDPRequest): RemoveHumanLinkedIDPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanLinkedIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanLinkedIDPRequest;
  static deserializeBinaryFromReader(message: RemoveHumanLinkedIDPRequest, reader: jspb.BinaryReader): RemoveHumanLinkedIDPRequest;
}

export namespace RemoveHumanLinkedIDPRequest {
  export type AsObject = {
    userId: string,
    idpId: string,
    linkedUserId: string,
  }
}

export class RemoveHumanLinkedIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveHumanLinkedIDPResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveHumanLinkedIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveHumanLinkedIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveHumanLinkedIDPResponse): RemoveHumanLinkedIDPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveHumanLinkedIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveHumanLinkedIDPResponse;
  static deserializeBinaryFromReader(message: RemoveHumanLinkedIDPResponse, reader: jspb.BinaryReader): RemoveHumanLinkedIDPResponse;
}

export namespace RemoveHumanLinkedIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListUserMembershipsRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ListUserMembershipsRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListUserMembershipsRequest;
  hasQuery(): boolean;
  clearQuery(): ListUserMembershipsRequest;

  getQueriesList(): Array<zitadel_user_pb.MembershipQuery>;
  setQueriesList(value: Array<zitadel_user_pb.MembershipQuery>): ListUserMembershipsRequest;
  clearQueriesList(): ListUserMembershipsRequest;
  addQueries(value?: zitadel_user_pb.MembershipQuery, index?: number): zitadel_user_pb.MembershipQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserMembershipsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserMembershipsRequest): ListUserMembershipsRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserMembershipsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserMembershipsRequest;
  static deserializeBinaryFromReader(message: ListUserMembershipsRequest, reader: jspb.BinaryReader): ListUserMembershipsRequest;
}

export namespace ListUserMembershipsRequest {
  export type AsObject = {
    userId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_user_pb.MembershipQuery.AsObject>,
  }
}

export class ListUserMembershipsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListUserMembershipsResponse;
  hasDetails(): boolean;
  clearDetails(): ListUserMembershipsResponse;

  getResultList(): Array<zitadel_user_pb.Membership>;
  setResultList(value: Array<zitadel_user_pb.Membership>): ListUserMembershipsResponse;
  clearResultList(): ListUserMembershipsResponse;
  addResult(value?: zitadel_user_pb.Membership, index?: number): zitadel_user_pb.Membership;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserMembershipsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserMembershipsResponse): ListUserMembershipsResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserMembershipsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserMembershipsResponse;
  static deserializeBinaryFromReader(message: ListUserMembershipsResponse, reader: jspb.BinaryReader): ListUserMembershipsResponse;
}

export namespace ListUserMembershipsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_user_pb.Membership.AsObject>,
  }
}

export class GetMyOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyOrgRequest): GetMyOrgRequest.AsObject;
  static serializeBinaryToWriter(message: GetMyOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrgRequest;
  static deserializeBinaryFromReader(message: GetMyOrgRequest, reader: jspb.BinaryReader): GetMyOrgRequest;
}

export namespace GetMyOrgRequest {
  export type AsObject = {
  }
}

export class GetMyOrgResponse extends jspb.Message {
  getOrg(): zitadel_org_pb.Org | undefined;
  setOrg(value?: zitadel_org_pb.Org): GetMyOrgResponse;
  hasOrg(): boolean;
  clearOrg(): GetMyOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMyOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMyOrgResponse): GetMyOrgResponse.AsObject;
  static serializeBinaryToWriter(message: GetMyOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMyOrgResponse;
  static deserializeBinaryFromReader(message: GetMyOrgResponse, reader: jspb.BinaryReader): GetMyOrgResponse;
}

export namespace GetMyOrgResponse {
  export type AsObject = {
    org?: zitadel_org_pb.Org.AsObject,
  }
}

export class GetOrgByDomainGlobalRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): GetOrgByDomainGlobalRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgByDomainGlobalRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgByDomainGlobalRequest): GetOrgByDomainGlobalRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgByDomainGlobalRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgByDomainGlobalRequest;
  static deserializeBinaryFromReader(message: GetOrgByDomainGlobalRequest, reader: jspb.BinaryReader): GetOrgByDomainGlobalRequest;
}

export namespace GetOrgByDomainGlobalRequest {
  export type AsObject = {
    domain: string,
  }
}

export class ListOrgChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListOrgChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgChangesRequest): ListOrgChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgChangesRequest;
  static deserializeBinaryFromReader(message: ListOrgChangesRequest, reader: jspb.BinaryReader): ListOrgChangesRequest;
}

export namespace ListOrgChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
  }
}

export class ListOrgChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListOrgChangesResponse;
  clearResultList(): ListOrgChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgChangesResponse): ListOrgChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgChangesResponse;
  static deserializeBinaryFromReader(message: ListOrgChangesResponse, reader: jspb.BinaryReader): ListOrgChangesResponse;
}

export namespace ListOrgChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class GetOrgByDomainGlobalResponse extends jspb.Message {
  getOrg(): zitadel_org_pb.Org | undefined;
  setOrg(value?: zitadel_org_pb.Org): GetOrgByDomainGlobalResponse;
  hasOrg(): boolean;
  clearOrg(): GetOrgByDomainGlobalResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgByDomainGlobalResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgByDomainGlobalResponse): GetOrgByDomainGlobalResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgByDomainGlobalResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgByDomainGlobalResponse;
  static deserializeBinaryFromReader(message: GetOrgByDomainGlobalResponse, reader: jspb.BinaryReader): GetOrgByDomainGlobalResponse;
}

export namespace GetOrgByDomainGlobalResponse {
  export type AsObject = {
    org?: zitadel_org_pb.Org.AsObject,
  }
}

export class AddOrgRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddOrgRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgRequest): AddOrgRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgRequest;
  static deserializeBinaryFromReader(message: AddOrgRequest, reader: jspb.BinaryReader): AddOrgRequest;
}

export namespace AddOrgRequest {
  export type AsObject = {
    name: string,
  }
}

export class AddOrgResponse extends jspb.Message {
  getId(): string;
  setId(value: string): AddOrgResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOrgResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgResponse): AddOrgResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgResponse;
  static deserializeBinaryFromReader(message: AddOrgResponse, reader: jspb.BinaryReader): AddOrgResponse;
}

export namespace AddOrgResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOrgRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UpdateOrgRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgRequest): UpdateOrgRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgRequest;
  static deserializeBinaryFromReader(message: UpdateOrgRequest, reader: jspb.BinaryReader): UpdateOrgRequest;
}

export namespace UpdateOrgRequest {
  export type AsObject = {
    name: string,
  }
}

export class UpdateOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgResponse): UpdateOrgResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgResponse;
  static deserializeBinaryFromReader(message: UpdateOrgResponse, reader: jspb.BinaryReader): UpdateOrgResponse;
}

export namespace UpdateOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateOrgRequest): DeactivateOrgRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateOrgRequest;
  static deserializeBinaryFromReader(message: DeactivateOrgRequest, reader: jspb.BinaryReader): DeactivateOrgRequest;
}

export namespace DeactivateOrgRequest {
  export type AsObject = {
  }
}

export class DeactivateOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateOrgResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateOrgResponse): DeactivateOrgResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateOrgResponse;
  static deserializeBinaryFromReader(message: DeactivateOrgResponse, reader: jspb.BinaryReader): DeactivateOrgResponse;
}

export namespace DeactivateOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateOrgRequest): ReactivateOrgRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateOrgRequest;
  static deserializeBinaryFromReader(message: ReactivateOrgRequest, reader: jspb.BinaryReader): ReactivateOrgRequest;
}

export namespace ReactivateOrgRequest {
  export type AsObject = {
  }
}

export class ReactivateOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateOrgResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateOrgResponse): ReactivateOrgResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateOrgResponse;
  static deserializeBinaryFromReader(message: ReactivateOrgResponse, reader: jspb.BinaryReader): ReactivateOrgResponse;
}

export namespace ReactivateOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveOrgRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgRequest): RemoveOrgRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgRequest;
  static deserializeBinaryFromReader(message: RemoveOrgRequest, reader: jspb.BinaryReader): RemoveOrgRequest;
}

export namespace RemoveOrgRequest {
  export type AsObject = {
  }
}

export class RemoveOrgResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveOrgResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOrgResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgResponse): RemoveOrgResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgResponse;
  static deserializeBinaryFromReader(message: RemoveOrgResponse, reader: jspb.BinaryReader): RemoveOrgResponse;
}

export namespace RemoveOrgResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListOrgDomainsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListOrgDomainsRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgDomainsRequest;

  getQueriesList(): Array<zitadel_org_pb.DomainSearchQuery>;
  setQueriesList(value: Array<zitadel_org_pb.DomainSearchQuery>): ListOrgDomainsRequest;
  clearQueriesList(): ListOrgDomainsRequest;
  addQueries(value?: zitadel_org_pb.DomainSearchQuery, index?: number): zitadel_org_pb.DomainSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgDomainsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgDomainsRequest): ListOrgDomainsRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgDomainsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgDomainsRequest;
  static deserializeBinaryFromReader(message: ListOrgDomainsRequest, reader: jspb.BinaryReader): ListOrgDomainsRequest;
}

export namespace ListOrgDomainsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_org_pb.DomainSearchQuery.AsObject>,
  }
}

export class ListOrgDomainsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListOrgDomainsResponse;
  hasDetails(): boolean;
  clearDetails(): ListOrgDomainsResponse;

  getResultList(): Array<zitadel_org_pb.Domain>;
  setResultList(value: Array<zitadel_org_pb.Domain>): ListOrgDomainsResponse;
  clearResultList(): ListOrgDomainsResponse;
  addResult(value?: zitadel_org_pb.Domain, index?: number): zitadel_org_pb.Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgDomainsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgDomainsResponse): ListOrgDomainsResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgDomainsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgDomainsResponse;
  static deserializeBinaryFromReader(message: ListOrgDomainsResponse, reader: jspb.BinaryReader): ListOrgDomainsResponse;
}

export namespace ListOrgDomainsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_org_pb.Domain.AsObject>,
  }
}

export class AddOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): AddOrgDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgDomainRequest): AddOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgDomainRequest;
  static deserializeBinaryFromReader(message: AddOrgDomainRequest, reader: jspb.BinaryReader): AddOrgDomainRequest;
}

export namespace AddOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class AddOrgDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOrgDomainResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrgDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgDomainResponse): AddOrgDomainResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrgDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgDomainResponse;
  static deserializeBinaryFromReader(message: AddOrgDomainResponse, reader: jspb.BinaryReader): AddOrgDomainResponse;
}

export namespace AddOrgDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): RemoveOrgDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgDomainRequest): RemoveOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgDomainRequest;
  static deserializeBinaryFromReader(message: RemoveOrgDomainRequest, reader: jspb.BinaryReader): RemoveOrgDomainRequest;
}

export namespace RemoveOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class RemoveOrgDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveOrgDomainResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOrgDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgDomainResponse): RemoveOrgDomainResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgDomainResponse;
  static deserializeBinaryFromReader(message: RemoveOrgDomainResponse, reader: jspb.BinaryReader): RemoveOrgDomainResponse;
}

export namespace RemoveOrgDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GenerateOrgDomainValidationRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): GenerateOrgDomainValidationRequest;

  getType(): zitadel_org_pb.DomainValidationType;
  setType(value: zitadel_org_pb.DomainValidationType): GenerateOrgDomainValidationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateOrgDomainValidationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateOrgDomainValidationRequest): GenerateOrgDomainValidationRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateOrgDomainValidationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateOrgDomainValidationRequest;
  static deserializeBinaryFromReader(message: GenerateOrgDomainValidationRequest, reader: jspb.BinaryReader): GenerateOrgDomainValidationRequest;
}

export namespace GenerateOrgDomainValidationRequest {
  export type AsObject = {
    domain: string,
    type: zitadel_org_pb.DomainValidationType,
  }
}

export class GenerateOrgDomainValidationResponse extends jspb.Message {
  getToken(): string;
  setToken(value: string): GenerateOrgDomainValidationResponse;

  getUrl(): string;
  setUrl(value: string): GenerateOrgDomainValidationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateOrgDomainValidationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateOrgDomainValidationResponse): GenerateOrgDomainValidationResponse.AsObject;
  static serializeBinaryToWriter(message: GenerateOrgDomainValidationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateOrgDomainValidationResponse;
  static deserializeBinaryFromReader(message: GenerateOrgDomainValidationResponse, reader: jspb.BinaryReader): GenerateOrgDomainValidationResponse;
}

export namespace GenerateOrgDomainValidationResponse {
  export type AsObject = {
    token: string,
    url: string,
  }
}

export class ValidateOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): ValidateOrgDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateOrgDomainRequest): ValidateOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: ValidateOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateOrgDomainRequest;
  static deserializeBinaryFromReader(message: ValidateOrgDomainRequest, reader: jspb.BinaryReader): ValidateOrgDomainRequest;
}

export namespace ValidateOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class ValidateOrgDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ValidateOrgDomainResponse;
  hasDetails(): boolean;
  clearDetails(): ValidateOrgDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ValidateOrgDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ValidateOrgDomainResponse): ValidateOrgDomainResponse.AsObject;
  static serializeBinaryToWriter(message: ValidateOrgDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ValidateOrgDomainResponse;
  static deserializeBinaryFromReader(message: ValidateOrgDomainResponse, reader: jspb.BinaryReader): ValidateOrgDomainResponse;
}

export namespace ValidateOrgDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetPrimaryOrgDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): SetPrimaryOrgDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPrimaryOrgDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPrimaryOrgDomainRequest): SetPrimaryOrgDomainRequest.AsObject;
  static serializeBinaryToWriter(message: SetPrimaryOrgDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPrimaryOrgDomainRequest;
  static deserializeBinaryFromReader(message: SetPrimaryOrgDomainRequest, reader: jspb.BinaryReader): SetPrimaryOrgDomainRequest;
}

export namespace SetPrimaryOrgDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class SetPrimaryOrgDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetPrimaryOrgDomainResponse;
  hasDetails(): boolean;
  clearDetails(): SetPrimaryOrgDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPrimaryOrgDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetPrimaryOrgDomainResponse): SetPrimaryOrgDomainResponse.AsObject;
  static serializeBinaryToWriter(message: SetPrimaryOrgDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPrimaryOrgDomainResponse;
  static deserializeBinaryFromReader(message: SetPrimaryOrgDomainResponse, reader: jspb.BinaryReader): SetPrimaryOrgDomainResponse;
}

export namespace SetPrimaryOrgDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListOrgMemberRolesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMemberRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMemberRolesRequest): ListOrgMemberRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgMemberRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMemberRolesRequest;
  static deserializeBinaryFromReader(message: ListOrgMemberRolesRequest, reader: jspb.BinaryReader): ListOrgMemberRolesRequest;
}

export namespace ListOrgMemberRolesRequest {
  export type AsObject = {
  }
}

export class ListOrgMemberRolesResponse extends jspb.Message {
  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListOrgMemberRolesResponse;
  clearResultList(): ListOrgMemberRolesResponse;
  addResult(value: string, index?: number): ListOrgMemberRolesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMemberRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMemberRolesResponse): ListOrgMemberRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgMemberRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMemberRolesResponse;
  static deserializeBinaryFromReader(message: ListOrgMemberRolesResponse, reader: jspb.BinaryReader): ListOrgMemberRolesResponse;
}

export namespace ListOrgMemberRolesResponse {
  export type AsObject = {
    resultList: Array<string>,
  }
}

export class ListOrgMembersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListOrgMembersRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgMembersRequest;

  getQueriesList(): Array<zitadel_member_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_member_pb.SearchQuery>): ListOrgMembersRequest;
  clearQueriesList(): ListOrgMembersRequest;
  addQueries(value?: zitadel_member_pb.SearchQuery, index?: number): zitadel_member_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMembersRequest): ListOrgMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMembersRequest;
  static deserializeBinaryFromReader(message: ListOrgMembersRequest, reader: jspb.BinaryReader): ListOrgMembersRequest;
}

export namespace ListOrgMembersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_member_pb.SearchQuery.AsObject>,
  }
}

export class ListOrgMembersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListOrgMembersResponse;
  hasDetails(): boolean;
  clearDetails(): ListOrgMembersResponse;

  getResultList(): Array<zitadel_member_pb.Member>;
  setResultList(value: Array<zitadel_member_pb.Member>): ListOrgMembersResponse;
  clearResultList(): ListOrgMembersResponse;
  addResult(value?: zitadel_member_pb.Member, index?: number): zitadel_member_pb.Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMembersResponse): ListOrgMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMembersResponse;
  static deserializeBinaryFromReader(message: ListOrgMembersResponse, reader: jspb.BinaryReader): ListOrgMembersResponse;
}

export namespace ListOrgMembersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_member_pb.Member.AsObject>,
  }
}

export class AddOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddOrgMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): AddOrgMemberRequest;
  clearRolesList(): AddOrgMemberRequest;
  addRoles(value: string, index?: number): AddOrgMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgMemberRequest): AddOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgMemberRequest;
  static deserializeBinaryFromReader(message: AddOrgMemberRequest, reader: jspb.BinaryReader): AddOrgMemberRequest;
}

export namespace AddOrgMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class AddOrgMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOrgMemberResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrgMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgMemberResponse): AddOrgMemberResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrgMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgMemberResponse;
  static deserializeBinaryFromReader(message: AddOrgMemberResponse, reader: jspb.BinaryReader): AddOrgMemberResponse;
}

export namespace AddOrgMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateOrgMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): UpdateOrgMemberRequest;
  clearRolesList(): UpdateOrgMemberRequest;
  addRoles(value: string, index?: number): UpdateOrgMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgMemberRequest): UpdateOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgMemberRequest;
  static deserializeBinaryFromReader(message: UpdateOrgMemberRequest, reader: jspb.BinaryReader): UpdateOrgMemberRequest;
}

export namespace UpdateOrgMemberRequest {
  export type AsObject = {
    userId: string,
    rolesList: Array<string>,
  }
}

export class UpdateOrgMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgMemberResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgMemberResponse): UpdateOrgMemberResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgMemberResponse;
  static deserializeBinaryFromReader(message: UpdateOrgMemberResponse, reader: jspb.BinaryReader): UpdateOrgMemberResponse;
}

export namespace UpdateOrgMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveOrgMemberRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveOrgMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgMemberRequest): RemoveOrgMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgMemberRequest;
  static deserializeBinaryFromReader(message: RemoveOrgMemberRequest, reader: jspb.BinaryReader): RemoveOrgMemberRequest;
}

export namespace RemoveOrgMemberRequest {
  export type AsObject = {
    userId: string,
  }
}

export class RemoveOrgMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveOrgMemberResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOrgMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgMemberResponse): RemoveOrgMemberResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgMemberResponse;
  static deserializeBinaryFromReader(message: RemoveOrgMemberResponse, reader: jspb.BinaryReader): RemoveOrgMemberResponse;
}

export namespace RemoveOrgMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListOrgMetadataRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListOrgMetadataRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgMetadataRequest;

  getQueriesList(): Array<zitadel_metadata_pb.MetadataQuery>;
  setQueriesList(value: Array<zitadel_metadata_pb.MetadataQuery>): ListOrgMetadataRequest;
  clearQueriesList(): ListOrgMetadataRequest;
  addQueries(value?: zitadel_metadata_pb.MetadataQuery, index?: number): zitadel_metadata_pb.MetadataQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMetadataRequest): ListOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMetadataRequest;
  static deserializeBinaryFromReader(message: ListOrgMetadataRequest, reader: jspb.BinaryReader): ListOrgMetadataRequest;
}

export namespace ListOrgMetadataRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_metadata_pb.MetadataQuery.AsObject>,
  }
}

export class ListOrgMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListOrgMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): ListOrgMetadataResponse;

  getResultList(): Array<zitadel_metadata_pb.Metadata>;
  setResultList(value: Array<zitadel_metadata_pb.Metadata>): ListOrgMetadataResponse;
  clearResultList(): ListOrgMetadataResponse;
  addResult(value?: zitadel_metadata_pb.Metadata, index?: number): zitadel_metadata_pb.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgMetadataResponse): ListOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgMetadataResponse;
  static deserializeBinaryFromReader(message: ListOrgMetadataResponse, reader: jspb.BinaryReader): ListOrgMetadataResponse;
}

export namespace ListOrgMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_metadata_pb.Metadata.AsObject>,
  }
}

export class GetOrgMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): GetOrgMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgMetadataRequest): GetOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgMetadataRequest;
  static deserializeBinaryFromReader(message: GetOrgMetadataRequest, reader: jspb.BinaryReader): GetOrgMetadataRequest;
}

export namespace GetOrgMetadataRequest {
  export type AsObject = {
    key: string,
  }
}

export class GetOrgMetadataResponse extends jspb.Message {
  getMetadata(): zitadel_metadata_pb.Metadata | undefined;
  setMetadata(value?: zitadel_metadata_pb.Metadata): GetOrgMetadataResponse;
  hasMetadata(): boolean;
  clearMetadata(): GetOrgMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgMetadataResponse): GetOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgMetadataResponse;
  static deserializeBinaryFromReader(message: GetOrgMetadataResponse, reader: jspb.BinaryReader): GetOrgMetadataResponse;
}

export namespace GetOrgMetadataResponse {
  export type AsObject = {
    metadata?: zitadel_metadata_pb.Metadata.AsObject,
  }
}

export class SetOrgMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): SetOrgMetadataRequest;

  getValue(): Uint8Array | string;
  getValue_asU8(): Uint8Array;
  getValue_asB64(): string;
  setValue(value: Uint8Array | string): SetOrgMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetOrgMetadataRequest): SetOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: SetOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetOrgMetadataRequest;
  static deserializeBinaryFromReader(message: SetOrgMetadataRequest, reader: jspb.BinaryReader): SetOrgMetadataRequest;
}

export namespace SetOrgMetadataRequest {
  export type AsObject = {
    key: string,
    value: Uint8Array | string,
  }
}

export class SetOrgMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetOrgMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): SetOrgMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetOrgMetadataResponse): SetOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: SetOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetOrgMetadataResponse;
  static deserializeBinaryFromReader(message: SetOrgMetadataResponse, reader: jspb.BinaryReader): SetOrgMetadataResponse;
}

export namespace SetOrgMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkSetOrgMetadataRequest extends jspb.Message {
  getMetadataList(): Array<BulkSetOrgMetadataRequest.Metadata>;
  setMetadataList(value: Array<BulkSetOrgMetadataRequest.Metadata>): BulkSetOrgMetadataRequest;
  clearMetadataList(): BulkSetOrgMetadataRequest;
  addMetadata(value?: BulkSetOrgMetadataRequest.Metadata, index?: number): BulkSetOrgMetadataRequest.Metadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetOrgMetadataRequest): BulkSetOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkSetOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetOrgMetadataRequest;
  static deserializeBinaryFromReader(message: BulkSetOrgMetadataRequest, reader: jspb.BinaryReader): BulkSetOrgMetadataRequest;
}

export namespace BulkSetOrgMetadataRequest {
  export type AsObject = {
    metadataList: Array<BulkSetOrgMetadataRequest.Metadata.AsObject>,
  }

  export class Metadata extends jspb.Message {
    getKey(): string;
    setKey(value: string): Metadata;

    getValue(): Uint8Array | string;
    getValue_asU8(): Uint8Array;
    getValue_asB64(): string;
    setValue(value: Uint8Array | string): Metadata;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Metadata.AsObject;
    static toObject(includeInstance: boolean, msg: Metadata): Metadata.AsObject;
    static serializeBinaryToWriter(message: Metadata, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Metadata;
    static deserializeBinaryFromReader(message: Metadata, reader: jspb.BinaryReader): Metadata;
  }

  export namespace Metadata {
    export type AsObject = {
      key: string,
      value: Uint8Array | string,
    }
  }

}

export class BulkSetOrgMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkSetOrgMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkSetOrgMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetOrgMetadataResponse): BulkSetOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkSetOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetOrgMetadataResponse;
  static deserializeBinaryFromReader(message: BulkSetOrgMetadataResponse, reader: jspb.BinaryReader): BulkSetOrgMetadataResponse;
}

export namespace BulkSetOrgMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveOrgMetadataRequest extends jspb.Message {
  getKey(): string;
  setKey(value: string): RemoveOrgMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgMetadataRequest): RemoveOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgMetadataRequest;
  static deserializeBinaryFromReader(message: RemoveOrgMetadataRequest, reader: jspb.BinaryReader): RemoveOrgMetadataRequest;
}

export namespace RemoveOrgMetadataRequest {
  export type AsObject = {
    key: string,
  }
}

export class RemoveOrgMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveOrgMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveOrgMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgMetadataResponse): RemoveOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgMetadataResponse;
  static deserializeBinaryFromReader(message: RemoveOrgMetadataResponse, reader: jspb.BinaryReader): RemoveOrgMetadataResponse;
}

export namespace RemoveOrgMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkRemoveOrgMetadataRequest extends jspb.Message {
  getKeysList(): Array<string>;
  setKeysList(value: Array<string>): BulkRemoveOrgMetadataRequest;
  clearKeysList(): BulkRemoveOrgMetadataRequest;
  addKeys(value: string, index?: number): BulkRemoveOrgMetadataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveOrgMetadataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveOrgMetadataRequest): BulkRemoveOrgMetadataRequest.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveOrgMetadataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveOrgMetadataRequest;
  static deserializeBinaryFromReader(message: BulkRemoveOrgMetadataRequest, reader: jspb.BinaryReader): BulkRemoveOrgMetadataRequest;
}

export namespace BulkRemoveOrgMetadataRequest {
  export type AsObject = {
    keysList: Array<string>,
  }
}

export class BulkRemoveOrgMetadataResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkRemoveOrgMetadataResponse;
  hasDetails(): boolean;
  clearDetails(): BulkRemoveOrgMetadataResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveOrgMetadataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveOrgMetadataResponse): BulkRemoveOrgMetadataResponse.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveOrgMetadataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveOrgMetadataResponse;
  static deserializeBinaryFromReader(message: BulkRemoveOrgMetadataResponse, reader: jspb.BinaryReader): BulkRemoveOrgMetadataResponse;
}

export namespace BulkRemoveOrgMetadataResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetProjectByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetProjectByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectByIDRequest): GetProjectByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetProjectByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectByIDRequest;
  static deserializeBinaryFromReader(message: GetProjectByIDRequest, reader: jspb.BinaryReader): GetProjectByIDRequest;
}

export namespace GetProjectByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetProjectByIDResponse extends jspb.Message {
  getProject(): zitadel_project_pb.Project | undefined;
  setProject(value?: zitadel_project_pb.Project): GetProjectByIDResponse;
  hasProject(): boolean;
  clearProject(): GetProjectByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectByIDResponse): GetProjectByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetProjectByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectByIDResponse;
  static deserializeBinaryFromReader(message: GetProjectByIDResponse, reader: jspb.BinaryReader): GetProjectByIDResponse;
}

export namespace GetProjectByIDResponse {
  export type AsObject = {
    project?: zitadel_project_pb.Project.AsObject,
  }
}

export class GetGrantedProjectByIDRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): GetGrantedProjectByIDRequest;

  getGrantId(): string;
  setGrantId(value: string): GetGrantedProjectByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGrantedProjectByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetGrantedProjectByIDRequest): GetGrantedProjectByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetGrantedProjectByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGrantedProjectByIDRequest;
  static deserializeBinaryFromReader(message: GetGrantedProjectByIDRequest, reader: jspb.BinaryReader): GetGrantedProjectByIDRequest;
}

export namespace GetGrantedProjectByIDRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
  }
}

export class GetGrantedProjectByIDResponse extends jspb.Message {
  getGrantedProject(): zitadel_project_pb.GrantedProject | undefined;
  setGrantedProject(value?: zitadel_project_pb.GrantedProject): GetGrantedProjectByIDResponse;
  hasGrantedProject(): boolean;
  clearGrantedProject(): GetGrantedProjectByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetGrantedProjectByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetGrantedProjectByIDResponse): GetGrantedProjectByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetGrantedProjectByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetGrantedProjectByIDResponse;
  static deserializeBinaryFromReader(message: GetGrantedProjectByIDResponse, reader: jspb.BinaryReader): GetGrantedProjectByIDResponse;
}

export namespace GetGrantedProjectByIDResponse {
  export type AsObject = {
    grantedProject?: zitadel_project_pb.GrantedProject.AsObject,
  }
}

export class ListProjectsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectsRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectsRequest;

  getQueriesList(): Array<zitadel_project_pb.ProjectQuery>;
  setQueriesList(value: Array<zitadel_project_pb.ProjectQuery>): ListProjectsRequest;
  clearQueriesList(): ListProjectsRequest;
  addQueries(value?: zitadel_project_pb.ProjectQuery, index?: number): zitadel_project_pb.ProjectQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectsRequest): ListProjectsRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectsRequest;
  static deserializeBinaryFromReader(message: ListProjectsRequest, reader: jspb.BinaryReader): ListProjectsRequest;
}

export namespace ListProjectsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.ProjectQuery.AsObject>,
  }
}

export class ListProjectsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectsResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectsResponse;

  getResultList(): Array<zitadel_project_pb.Project>;
  setResultList(value: Array<zitadel_project_pb.Project>): ListProjectsResponse;
  clearResultList(): ListProjectsResponse;
  addResult(value?: zitadel_project_pb.Project, index?: number): zitadel_project_pb.Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectsResponse): ListProjectsResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectsResponse;
  static deserializeBinaryFromReader(message: ListProjectsResponse, reader: jspb.BinaryReader): ListProjectsResponse;
}

export namespace ListProjectsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.Project.AsObject>,
  }
}

export class ListGrantedProjectsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListGrantedProjectsRequest;
  hasQuery(): boolean;
  clearQuery(): ListGrantedProjectsRequest;

  getQueriesList(): Array<zitadel_project_pb.ProjectQuery>;
  setQueriesList(value: Array<zitadel_project_pb.ProjectQuery>): ListGrantedProjectsRequest;
  clearQueriesList(): ListGrantedProjectsRequest;
  addQueries(value?: zitadel_project_pb.ProjectQuery, index?: number): zitadel_project_pb.ProjectQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGrantedProjectsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGrantedProjectsRequest): ListGrantedProjectsRequest.AsObject;
  static serializeBinaryToWriter(message: ListGrantedProjectsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGrantedProjectsRequest;
  static deserializeBinaryFromReader(message: ListGrantedProjectsRequest, reader: jspb.BinaryReader): ListGrantedProjectsRequest;
}

export namespace ListGrantedProjectsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.ProjectQuery.AsObject>,
  }
}

export class ListGrantedProjectsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListGrantedProjectsResponse;
  hasDetails(): boolean;
  clearDetails(): ListGrantedProjectsResponse;

  getResultList(): Array<zitadel_project_pb.GrantedProject>;
  setResultList(value: Array<zitadel_project_pb.GrantedProject>): ListGrantedProjectsResponse;
  clearResultList(): ListGrantedProjectsResponse;
  addResult(value?: zitadel_project_pb.GrantedProject, index?: number): zitadel_project_pb.GrantedProject;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGrantedProjectsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGrantedProjectsResponse): ListGrantedProjectsResponse.AsObject;
  static serializeBinaryToWriter(message: ListGrantedProjectsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGrantedProjectsResponse;
  static deserializeBinaryFromReader(message: ListGrantedProjectsResponse, reader: jspb.BinaryReader): ListGrantedProjectsResponse;
}

export namespace ListGrantedProjectsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.GrantedProject.AsObject>,
  }
}

export class ListProjectChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListProjectChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectChangesRequest;

  getProjectId(): string;
  setProjectId(value: string): ListProjectChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectChangesRequest): ListProjectChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectChangesRequest;
  static deserializeBinaryFromReader(message: ListProjectChangesRequest, reader: jspb.BinaryReader): ListProjectChangesRequest;
}

export namespace ListProjectChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
    projectId: string,
  }
}

export class ListProjectChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListProjectChangesResponse;
  clearResultList(): ListProjectChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectChangesResponse): ListProjectChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectChangesResponse;
  static deserializeBinaryFromReader(message: ListProjectChangesResponse, reader: jspb.BinaryReader): ListProjectChangesResponse;
}

export namespace ListProjectChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class AddProjectRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddProjectRequest;

  getProjectRoleAssertion(): boolean;
  setProjectRoleAssertion(value: boolean): AddProjectRequest;

  getProjectRoleCheck(): boolean;
  setProjectRoleCheck(value: boolean): AddProjectRequest;

  getHasProjectCheck(): boolean;
  setHasProjectCheck(value: boolean): AddProjectRequest;

  getPrivateLabelingSetting(): zitadel_project_pb.PrivateLabelingSetting;
  setPrivateLabelingSetting(value: zitadel_project_pb.PrivateLabelingSetting): AddProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectRequest): AddProjectRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectRequest;
  static deserializeBinaryFromReader(message: AddProjectRequest, reader: jspb.BinaryReader): AddProjectRequest;
}

export namespace AddProjectRequest {
  export type AsObject = {
    name: string,
    projectRoleAssertion: boolean,
    projectRoleCheck: boolean,
    hasProjectCheck: boolean,
    privateLabelingSetting: zitadel_project_pb.PrivateLabelingSetting,
  }
}

export class AddProjectResponse extends jspb.Message {
  getId(): string;
  setId(value: string): AddProjectResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddProjectResponse;
  hasDetails(): boolean;
  clearDetails(): AddProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectResponse): AddProjectResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectResponse;
  static deserializeBinaryFromReader(message: AddProjectResponse, reader: jspb.BinaryReader): AddProjectResponse;
}

export namespace AddProjectResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateProjectRequest;

  getName(): string;
  setName(value: string): UpdateProjectRequest;

  getProjectRoleAssertion(): boolean;
  setProjectRoleAssertion(value: boolean): UpdateProjectRequest;

  getProjectRoleCheck(): boolean;
  setProjectRoleCheck(value: boolean): UpdateProjectRequest;

  getHasProjectCheck(): boolean;
  setHasProjectCheck(value: boolean): UpdateProjectRequest;

  getPrivateLabelingSetting(): zitadel_project_pb.PrivateLabelingSetting;
  setPrivateLabelingSetting(value: zitadel_project_pb.PrivateLabelingSetting): UpdateProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRequest): UpdateProjectRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRequest;
  static deserializeBinaryFromReader(message: UpdateProjectRequest, reader: jspb.BinaryReader): UpdateProjectRequest;
}

export namespace UpdateProjectRequest {
  export type AsObject = {
    id: string,
    name: string,
    projectRoleAssertion: boolean,
    projectRoleCheck: boolean,
    hasProjectCheck: boolean,
    privateLabelingSetting: zitadel_project_pb.PrivateLabelingSetting,
  }
}

export class UpdateProjectResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateProjectResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectResponse): UpdateProjectResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectResponse;
  static deserializeBinaryFromReader(message: UpdateProjectResponse, reader: jspb.BinaryReader): UpdateProjectResponse;
}

export namespace UpdateProjectResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateProjectRequest): DeactivateProjectRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateProjectRequest;
  static deserializeBinaryFromReader(message: DeactivateProjectRequest, reader: jspb.BinaryReader): DeactivateProjectRequest;
}

export namespace DeactivateProjectRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateProjectResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateProjectResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateProjectResponse): DeactivateProjectResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateProjectResponse;
  static deserializeBinaryFromReader(message: DeactivateProjectResponse, reader: jspb.BinaryReader): DeactivateProjectResponse;
}

export namespace DeactivateProjectResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ReactivateProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateProjectRequest): ReactivateProjectRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateProjectRequest;
  static deserializeBinaryFromReader(message: ReactivateProjectRequest, reader: jspb.BinaryReader): ReactivateProjectRequest;
}

export namespace ReactivateProjectRequest {
  export type AsObject = {
    id: string,
  }
}

export class ReactivateProjectResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateProjectResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateProjectResponse): ReactivateProjectResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateProjectResponse;
  static deserializeBinaryFromReader(message: ReactivateProjectResponse, reader: jspb.BinaryReader): ReactivateProjectResponse;
}

export namespace ReactivateProjectResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveProjectRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RemoveProjectRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectRequest): RemoveProjectRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectRequest;
  static deserializeBinaryFromReader(message: RemoveProjectRequest, reader: jspb.BinaryReader): RemoveProjectRequest;
}

export namespace RemoveProjectRequest {
  export type AsObject = {
    id: string,
  }
}

export class RemoveProjectResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveProjectResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectResponse): RemoveProjectResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectResponse;
  static deserializeBinaryFromReader(message: RemoveProjectResponse, reader: jspb.BinaryReader): RemoveProjectResponse;
}

export namespace RemoveProjectResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProjectMemberRolesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectMemberRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectMemberRolesRequest): ListProjectMemberRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectMemberRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectMemberRolesRequest;
  static deserializeBinaryFromReader(message: ListProjectMemberRolesRequest, reader: jspb.BinaryReader): ListProjectMemberRolesRequest;
}

export namespace ListProjectMemberRolesRequest {
  export type AsObject = {
  }
}

export class ListProjectMemberRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectMemberRolesResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectMemberRolesResponse;

  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListProjectMemberRolesResponse;
  clearResultList(): ListProjectMemberRolesResponse;
  addResult(value: string, index?: number): ListProjectMemberRolesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectMemberRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectMemberRolesResponse): ListProjectMemberRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectMemberRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectMemberRolesResponse;
  static deserializeBinaryFromReader(message: ListProjectMemberRolesResponse, reader: jspb.BinaryReader): ListProjectMemberRolesResponse;
}

export namespace ListProjectMemberRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<string>,
  }
}

export class AddProjectRoleRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddProjectRoleRequest;

  getRoleKey(): string;
  setRoleKey(value: string): AddProjectRoleRequest;

  getDisplayName(): string;
  setDisplayName(value: string): AddProjectRoleRequest;

  getGroup(): string;
  setGroup(value: string): AddProjectRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectRoleRequest): AddProjectRoleRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectRoleRequest;
  static deserializeBinaryFromReader(message: AddProjectRoleRequest, reader: jspb.BinaryReader): AddProjectRoleRequest;
}

export namespace AddProjectRoleRequest {
  export type AsObject = {
    projectId: string,
    roleKey: string,
    displayName: string,
    group: string,
  }
}

export class AddProjectRoleResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddProjectRoleResponse;
  hasDetails(): boolean;
  clearDetails(): AddProjectRoleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectRoleResponse): AddProjectRoleResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectRoleResponse;
  static deserializeBinaryFromReader(message: AddProjectRoleResponse, reader: jspb.BinaryReader): AddProjectRoleResponse;
}

export namespace AddProjectRoleResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkAddProjectRolesRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): BulkAddProjectRolesRequest;

  getRolesList(): Array<BulkAddProjectRolesRequest.Role>;
  setRolesList(value: Array<BulkAddProjectRolesRequest.Role>): BulkAddProjectRolesRequest;
  clearRolesList(): BulkAddProjectRolesRequest;
  addRoles(value?: BulkAddProjectRolesRequest.Role, index?: number): BulkAddProjectRolesRequest.Role;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkAddProjectRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkAddProjectRolesRequest): BulkAddProjectRolesRequest.AsObject;
  static serializeBinaryToWriter(message: BulkAddProjectRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkAddProjectRolesRequest;
  static deserializeBinaryFromReader(message: BulkAddProjectRolesRequest, reader: jspb.BinaryReader): BulkAddProjectRolesRequest;
}

export namespace BulkAddProjectRolesRequest {
  export type AsObject = {
    projectId: string,
    rolesList: Array<BulkAddProjectRolesRequest.Role.AsObject>,
  }

  export class Role extends jspb.Message {
    getKey(): string;
    setKey(value: string): Role;

    getDisplayName(): string;
    setDisplayName(value: string): Role;

    getGroup(): string;
    setGroup(value: string): Role;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Role.AsObject;
    static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
    static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Role;
    static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
  }

  export namespace Role {
    export type AsObject = {
      key: string,
      displayName: string,
      group: string,
    }
  }

}

export class BulkAddProjectRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkAddProjectRolesResponse;
  hasDetails(): boolean;
  clearDetails(): BulkAddProjectRolesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkAddProjectRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkAddProjectRolesResponse): BulkAddProjectRolesResponse.AsObject;
  static serializeBinaryToWriter(message: BulkAddProjectRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkAddProjectRolesResponse;
  static deserializeBinaryFromReader(message: BulkAddProjectRolesResponse, reader: jspb.BinaryReader): BulkAddProjectRolesResponse;
}

export namespace BulkAddProjectRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateProjectRoleRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateProjectRoleRequest;

  getRoleKey(): string;
  setRoleKey(value: string): UpdateProjectRoleRequest;

  getDisplayName(): string;
  setDisplayName(value: string): UpdateProjectRoleRequest;

  getGroup(): string;
  setGroup(value: string): UpdateProjectRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRoleRequest): UpdateProjectRoleRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRoleRequest;
  static deserializeBinaryFromReader(message: UpdateProjectRoleRequest, reader: jspb.BinaryReader): UpdateProjectRoleRequest;
}

export namespace UpdateProjectRoleRequest {
  export type AsObject = {
    projectId: string,
    roleKey: string,
    displayName: string,
    group: string,
  }
}

export class UpdateProjectRoleResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateProjectRoleResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateProjectRoleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRoleResponse): UpdateProjectRoleResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRoleResponse;
  static deserializeBinaryFromReader(message: UpdateProjectRoleResponse, reader: jspb.BinaryReader): UpdateProjectRoleResponse;
}

export namespace UpdateProjectRoleResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveProjectRoleRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveProjectRoleRequest;

  getRoleKey(): string;
  setRoleKey(value: string): RemoveProjectRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectRoleRequest): RemoveProjectRoleRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectRoleRequest;
  static deserializeBinaryFromReader(message: RemoveProjectRoleRequest, reader: jspb.BinaryReader): RemoveProjectRoleRequest;
}

export namespace RemoveProjectRoleRequest {
  export type AsObject = {
    projectId: string,
    roleKey: string,
  }
}

export class RemoveProjectRoleResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveProjectRoleResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveProjectRoleResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectRoleResponse): RemoveProjectRoleResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectRoleResponse;
  static deserializeBinaryFromReader(message: RemoveProjectRoleResponse, reader: jspb.BinaryReader): RemoveProjectRoleResponse;
}

export namespace RemoveProjectRoleResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProjectRolesRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListProjectRolesRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectRolesRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectRolesRequest;

  getQueriesList(): Array<zitadel_project_pb.RoleQuery>;
  setQueriesList(value: Array<zitadel_project_pb.RoleQuery>): ListProjectRolesRequest;
  clearQueriesList(): ListProjectRolesRequest;
  addQueries(value?: zitadel_project_pb.RoleQuery, index?: number): zitadel_project_pb.RoleQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectRolesRequest): ListProjectRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectRolesRequest;
  static deserializeBinaryFromReader(message: ListProjectRolesRequest, reader: jspb.BinaryReader): ListProjectRolesRequest;
}

export namespace ListProjectRolesRequest {
  export type AsObject = {
    projectId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.RoleQuery.AsObject>,
  }
}

export class ListProjectRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectRolesResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectRolesResponse;

  getResultList(): Array<zitadel_project_pb.Role>;
  setResultList(value: Array<zitadel_project_pb.Role>): ListProjectRolesResponse;
  clearResultList(): ListProjectRolesResponse;
  addResult(value?: zitadel_project_pb.Role, index?: number): zitadel_project_pb.Role;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectRolesResponse): ListProjectRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectRolesResponse;
  static deserializeBinaryFromReader(message: ListProjectRolesResponse, reader: jspb.BinaryReader): ListProjectRolesResponse;
}

export namespace ListProjectRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.Role.AsObject>,
  }
}

export class ListGrantedProjectRolesRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListGrantedProjectRolesRequest;

  getGrantId(): string;
  setGrantId(value: string): ListGrantedProjectRolesRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListGrantedProjectRolesRequest;
  hasQuery(): boolean;
  clearQuery(): ListGrantedProjectRolesRequest;

  getQueriesList(): Array<zitadel_project_pb.RoleQuery>;
  setQueriesList(value: Array<zitadel_project_pb.RoleQuery>): ListGrantedProjectRolesRequest;
  clearQueriesList(): ListGrantedProjectRolesRequest;
  addQueries(value?: zitadel_project_pb.RoleQuery, index?: number): zitadel_project_pb.RoleQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGrantedProjectRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListGrantedProjectRolesRequest): ListGrantedProjectRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListGrantedProjectRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGrantedProjectRolesRequest;
  static deserializeBinaryFromReader(message: ListGrantedProjectRolesRequest, reader: jspb.BinaryReader): ListGrantedProjectRolesRequest;
}

export namespace ListGrantedProjectRolesRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.RoleQuery.AsObject>,
  }
}

export class ListGrantedProjectRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListGrantedProjectRolesResponse;
  hasDetails(): boolean;
  clearDetails(): ListGrantedProjectRolesResponse;

  getResultList(): Array<zitadel_project_pb.Role>;
  setResultList(value: Array<zitadel_project_pb.Role>): ListGrantedProjectRolesResponse;
  clearResultList(): ListGrantedProjectRolesResponse;
  addResult(value?: zitadel_project_pb.Role, index?: number): zitadel_project_pb.Role;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListGrantedProjectRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListGrantedProjectRolesResponse): ListGrantedProjectRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListGrantedProjectRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListGrantedProjectRolesResponse;
  static deserializeBinaryFromReader(message: ListGrantedProjectRolesResponse, reader: jspb.BinaryReader): ListGrantedProjectRolesResponse;
}

export namespace ListGrantedProjectRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.Role.AsObject>,
  }
}

export class ListProjectMembersRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListProjectMembersRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectMembersRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectMembersRequest;

  getQueriesList(): Array<zitadel_member_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_member_pb.SearchQuery>): ListProjectMembersRequest;
  clearQueriesList(): ListProjectMembersRequest;
  addQueries(value?: zitadel_member_pb.SearchQuery, index?: number): zitadel_member_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectMembersRequest): ListProjectMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectMembersRequest;
  static deserializeBinaryFromReader(message: ListProjectMembersRequest, reader: jspb.BinaryReader): ListProjectMembersRequest;
}

export namespace ListProjectMembersRequest {
  export type AsObject = {
    projectId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_member_pb.SearchQuery.AsObject>,
  }
}

export class ListProjectMembersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectMembersResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectMembersResponse;

  getResultList(): Array<zitadel_member_pb.Member>;
  setResultList(value: Array<zitadel_member_pb.Member>): ListProjectMembersResponse;
  clearResultList(): ListProjectMembersResponse;
  addResult(value?: zitadel_member_pb.Member, index?: number): zitadel_member_pb.Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectMembersResponse): ListProjectMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectMembersResponse;
  static deserializeBinaryFromReader(message: ListProjectMembersResponse, reader: jspb.BinaryReader): ListProjectMembersResponse;
}

export namespace ListProjectMembersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_member_pb.Member.AsObject>,
  }
}

export class AddProjectMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddProjectMemberRequest;

  getUserId(): string;
  setUserId(value: string): AddProjectMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): AddProjectMemberRequest;
  clearRolesList(): AddProjectMemberRequest;
  addRoles(value: string, index?: number): AddProjectMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectMemberRequest): AddProjectMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectMemberRequest;
  static deserializeBinaryFromReader(message: AddProjectMemberRequest, reader: jspb.BinaryReader): AddProjectMemberRequest;
}

export namespace AddProjectMemberRequest {
  export type AsObject = {
    projectId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class AddProjectMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddProjectMemberResponse;
  hasDetails(): boolean;
  clearDetails(): AddProjectMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectMemberResponse): AddProjectMemberResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectMemberResponse;
  static deserializeBinaryFromReader(message: AddProjectMemberResponse, reader: jspb.BinaryReader): AddProjectMemberResponse;
}

export namespace AddProjectMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateProjectMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateProjectMemberRequest;

  getUserId(): string;
  setUserId(value: string): UpdateProjectMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): UpdateProjectMemberRequest;
  clearRolesList(): UpdateProjectMemberRequest;
  addRoles(value: string, index?: number): UpdateProjectMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectMemberRequest): UpdateProjectMemberRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectMemberRequest;
  static deserializeBinaryFromReader(message: UpdateProjectMemberRequest, reader: jspb.BinaryReader): UpdateProjectMemberRequest;
}

export namespace UpdateProjectMemberRequest {
  export type AsObject = {
    projectId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class UpdateProjectMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateProjectMemberResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateProjectMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectMemberResponse): UpdateProjectMemberResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectMemberResponse;
  static deserializeBinaryFromReader(message: UpdateProjectMemberResponse, reader: jspb.BinaryReader): UpdateProjectMemberResponse;
}

export namespace UpdateProjectMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveProjectMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveProjectMemberRequest;

  getUserId(): string;
  setUserId(value: string): RemoveProjectMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectMemberRequest): RemoveProjectMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectMemberRequest;
  static deserializeBinaryFromReader(message: RemoveProjectMemberRequest, reader: jspb.BinaryReader): RemoveProjectMemberRequest;
}

export namespace RemoveProjectMemberRequest {
  export type AsObject = {
    projectId: string,
    userId: string,
  }
}

export class RemoveProjectMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveProjectMemberResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveProjectMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectMemberResponse): RemoveProjectMemberResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectMemberResponse;
  static deserializeBinaryFromReader(message: RemoveProjectMemberResponse, reader: jspb.BinaryReader): RemoveProjectMemberResponse;
}

export namespace RemoveProjectMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetAppByIDRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): GetAppByIDRequest;

  getAppId(): string;
  setAppId(value: string): GetAppByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAppByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAppByIDRequest): GetAppByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetAppByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAppByIDRequest;
  static deserializeBinaryFromReader(message: GetAppByIDRequest, reader: jspb.BinaryReader): GetAppByIDRequest;
}

export namespace GetAppByIDRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class GetAppByIDResponse extends jspb.Message {
  getApp(): zitadel_app_pb.App | undefined;
  setApp(value?: zitadel_app_pb.App): GetAppByIDResponse;
  hasApp(): boolean;
  clearApp(): GetAppByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAppByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAppByIDResponse): GetAppByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetAppByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAppByIDResponse;
  static deserializeBinaryFromReader(message: GetAppByIDResponse, reader: jspb.BinaryReader): GetAppByIDResponse;
}

export namespace GetAppByIDResponse {
  export type AsObject = {
    app?: zitadel_app_pb.App.AsObject,
  }
}

export class ListAppsRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListAppsRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListAppsRequest;
  hasQuery(): boolean;
  clearQuery(): ListAppsRequest;

  getQueriesList(): Array<zitadel_app_pb.AppQuery>;
  setQueriesList(value: Array<zitadel_app_pb.AppQuery>): ListAppsRequest;
  clearQueriesList(): ListAppsRequest;
  addQueries(value?: zitadel_app_pb.AppQuery, index?: number): zitadel_app_pb.AppQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppsRequest): ListAppsRequest.AsObject;
  static serializeBinaryToWriter(message: ListAppsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppsRequest;
  static deserializeBinaryFromReader(message: ListAppsRequest, reader: jspb.BinaryReader): ListAppsRequest;
}

export namespace ListAppsRequest {
  export type AsObject = {
    projectId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_app_pb.AppQuery.AsObject>,
  }
}

export class ListAppsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListAppsResponse;
  hasDetails(): boolean;
  clearDetails(): ListAppsResponse;

  getResultList(): Array<zitadel_app_pb.App>;
  setResultList(value: Array<zitadel_app_pb.App>): ListAppsResponse;
  clearResultList(): ListAppsResponse;
  addResult(value?: zitadel_app_pb.App, index?: number): zitadel_app_pb.App;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppsResponse): ListAppsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAppsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppsResponse;
  static deserializeBinaryFromReader(message: ListAppsResponse, reader: jspb.BinaryReader): ListAppsResponse;
}

export namespace ListAppsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_app_pb.App.AsObject>,
  }
}

export class ListAppChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListAppChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListAppChangesRequest;

  getProjectId(): string;
  setProjectId(value: string): ListAppChangesRequest;

  getAppId(): string;
  setAppId(value: string): ListAppChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppChangesRequest): ListAppChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListAppChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppChangesRequest;
  static deserializeBinaryFromReader(message: ListAppChangesRequest, reader: jspb.BinaryReader): ListAppChangesRequest;
}

export namespace ListAppChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
    projectId: string,
    appId: string,
  }
}

export class ListAppChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListAppChangesResponse;
  clearResultList(): ListAppChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppChangesResponse): ListAppChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListAppChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppChangesResponse;
  static deserializeBinaryFromReader(message: ListAppChangesResponse, reader: jspb.BinaryReader): ListAppChangesResponse;
}

export namespace ListAppChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class AddOIDCAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddOIDCAppRequest;

  getName(): string;
  setName(value: string): AddOIDCAppRequest;

  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): AddOIDCAppRequest;
  clearRedirectUrisList(): AddOIDCAppRequest;
  addRedirectUris(value: string, index?: number): AddOIDCAppRequest;

  getResponseTypesList(): Array<zitadel_app_pb.OIDCResponseType>;
  setResponseTypesList(value: Array<zitadel_app_pb.OIDCResponseType>): AddOIDCAppRequest;
  clearResponseTypesList(): AddOIDCAppRequest;
  addResponseTypes(value: zitadel_app_pb.OIDCResponseType, index?: number): AddOIDCAppRequest;

  getGrantTypesList(): Array<zitadel_app_pb.OIDCGrantType>;
  setGrantTypesList(value: Array<zitadel_app_pb.OIDCGrantType>): AddOIDCAppRequest;
  clearGrantTypesList(): AddOIDCAppRequest;
  addGrantTypes(value: zitadel_app_pb.OIDCGrantType, index?: number): AddOIDCAppRequest;

  getAppType(): zitadel_app_pb.OIDCAppType;
  setAppType(value: zitadel_app_pb.OIDCAppType): AddOIDCAppRequest;

  getAuthMethodType(): zitadel_app_pb.OIDCAuthMethodType;
  setAuthMethodType(value: zitadel_app_pb.OIDCAuthMethodType): AddOIDCAppRequest;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): AddOIDCAppRequest;
  clearPostLogoutRedirectUrisList(): AddOIDCAppRequest;
  addPostLogoutRedirectUris(value: string, index?: number): AddOIDCAppRequest;

  getVersion(): zitadel_app_pb.OIDCVersion;
  setVersion(value: zitadel_app_pb.OIDCVersion): AddOIDCAppRequest;

  getDevMode(): boolean;
  setDevMode(value: boolean): AddOIDCAppRequest;

  getAccessTokenType(): zitadel_app_pb.OIDCTokenType;
  setAccessTokenType(value: zitadel_app_pb.OIDCTokenType): AddOIDCAppRequest;

  getAccessTokenRoleAssertion(): boolean;
  setAccessTokenRoleAssertion(value: boolean): AddOIDCAppRequest;

  getIdTokenRoleAssertion(): boolean;
  setIdTokenRoleAssertion(value: boolean): AddOIDCAppRequest;

  getIdTokenUserinfoAssertion(): boolean;
  setIdTokenUserinfoAssertion(value: boolean): AddOIDCAppRequest;

  getClockSkew(): google_protobuf_duration_pb.Duration | undefined;
  setClockSkew(value?: google_protobuf_duration_pb.Duration): AddOIDCAppRequest;
  hasClockSkew(): boolean;
  clearClockSkew(): AddOIDCAppRequest;

  getAdditionalOriginsList(): Array<string>;
  setAdditionalOriginsList(value: Array<string>): AddOIDCAppRequest;
  clearAdditionalOriginsList(): AddOIDCAppRequest;
  addAdditionalOrigins(value: string, index?: number): AddOIDCAppRequest;

  getSkipNativeAppSuccessPage(): boolean;
  setSkipNativeAppSuccessPage(value: boolean): AddOIDCAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCAppRequest): AddOIDCAppRequest.AsObject;
  static serializeBinaryToWriter(message: AddOIDCAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCAppRequest;
  static deserializeBinaryFromReader(message: AddOIDCAppRequest, reader: jspb.BinaryReader): AddOIDCAppRequest;
}

export namespace AddOIDCAppRequest {
  export type AsObject = {
    projectId: string,
    name: string,
    redirectUrisList: Array<string>,
    responseTypesList: Array<zitadel_app_pb.OIDCResponseType>,
    grantTypesList: Array<zitadel_app_pb.OIDCGrantType>,
    appType: zitadel_app_pb.OIDCAppType,
    authMethodType: zitadel_app_pb.OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
    version: zitadel_app_pb.OIDCVersion,
    devMode: boolean,
    accessTokenType: zitadel_app_pb.OIDCTokenType,
    accessTokenRoleAssertion: boolean,
    idTokenRoleAssertion: boolean,
    idTokenUserinfoAssertion: boolean,
    clockSkew?: google_protobuf_duration_pb.Duration.AsObject,
    additionalOriginsList: Array<string>,
    skipNativeAppSuccessPage: boolean,
  }
}

export class AddOIDCAppResponse extends jspb.Message {
  getAppId(): string;
  setAppId(value: string): AddOIDCAppResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOIDCAppResponse;
  hasDetails(): boolean;
  clearDetails(): AddOIDCAppResponse;

  getClientId(): string;
  setClientId(value: string): AddOIDCAppResponse;

  getClientSecret(): string;
  setClientSecret(value: string): AddOIDCAppResponse;

  getNoneCompliant(): boolean;
  setNoneCompliant(value: boolean): AddOIDCAppResponse;

  getComplianceProblemsList(): Array<zitadel_message_pb.LocalizedMessage>;
  setComplianceProblemsList(value: Array<zitadel_message_pb.LocalizedMessage>): AddOIDCAppResponse;
  clearComplianceProblemsList(): AddOIDCAppResponse;
  addComplianceProblems(value?: zitadel_message_pb.LocalizedMessage, index?: number): zitadel_message_pb.LocalizedMessage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOIDCAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOIDCAppResponse): AddOIDCAppResponse.AsObject;
  static serializeBinaryToWriter(message: AddOIDCAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOIDCAppResponse;
  static deserializeBinaryFromReader(message: AddOIDCAppResponse, reader: jspb.BinaryReader): AddOIDCAppResponse;
}

export namespace AddOIDCAppResponse {
  export type AsObject = {
    appId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    clientId: string,
    clientSecret: string,
    noneCompliant: boolean,
    complianceProblemsList: Array<zitadel_message_pb.LocalizedMessage.AsObject>,
  }
}

export class AddSAMLAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddSAMLAppRequest;

  getName(): string;
  setName(value: string): AddSAMLAppRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): AddSAMLAppRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): AddSAMLAppRequest;

  getMetadataCase(): AddSAMLAppRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLAppRequest): AddSAMLAppRequest.AsObject;
  static serializeBinaryToWriter(message: AddSAMLAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLAppRequest;
  static deserializeBinaryFromReader(message: AddSAMLAppRequest, reader: jspb.BinaryReader): AddSAMLAppRequest;
}

export namespace AddSAMLAppRequest {
  export type AsObject = {
    projectId: string,
    name: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 3,
    METADATA_URL = 4,
  }
}

export class AddSAMLAppResponse extends jspb.Message {
  getAppId(): string;
  setAppId(value: string): AddSAMLAppResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSAMLAppResponse;
  hasDetails(): boolean;
  clearDetails(): AddSAMLAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLAppResponse): AddSAMLAppResponse.AsObject;
  static serializeBinaryToWriter(message: AddSAMLAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLAppResponse;
  static deserializeBinaryFromReader(message: AddSAMLAppResponse, reader: jspb.BinaryReader): AddSAMLAppResponse;
}

export namespace AddSAMLAppResponse {
  export type AsObject = {
    appId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddAPIAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddAPIAppRequest;

  getName(): string;
  setName(value: string): AddAPIAppRequest;

  getAuthMethodType(): zitadel_app_pb.APIAuthMethodType;
  setAuthMethodType(value: zitadel_app_pb.APIAuthMethodType): AddAPIAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAPIAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAPIAppRequest): AddAPIAppRequest.AsObject;
  static serializeBinaryToWriter(message: AddAPIAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAPIAppRequest;
  static deserializeBinaryFromReader(message: AddAPIAppRequest, reader: jspb.BinaryReader): AddAPIAppRequest;
}

export namespace AddAPIAppRequest {
  export type AsObject = {
    projectId: string,
    name: string,
    authMethodType: zitadel_app_pb.APIAuthMethodType,
  }
}

export class AddAPIAppResponse extends jspb.Message {
  getAppId(): string;
  setAppId(value: string): AddAPIAppResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAPIAppResponse;
  hasDetails(): boolean;
  clearDetails(): AddAPIAppResponse;

  getClientId(): string;
  setClientId(value: string): AddAPIAppResponse;

  getClientSecret(): string;
  setClientSecret(value: string): AddAPIAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAPIAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAPIAppResponse): AddAPIAppResponse.AsObject;
  static serializeBinaryToWriter(message: AddAPIAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAPIAppResponse;
  static deserializeBinaryFromReader(message: AddAPIAppResponse, reader: jspb.BinaryReader): AddAPIAppResponse;
}

export namespace AddAPIAppResponse {
  export type AsObject = {
    appId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    clientId: string,
    clientSecret: string,
  }
}

export class UpdateAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateAppRequest;

  getAppId(): string;
  setAppId(value: string): UpdateAppRequest;

  getName(): string;
  setName(value: string): UpdateAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppRequest): UpdateAppRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppRequest;
  static deserializeBinaryFromReader(message: UpdateAppRequest, reader: jspb.BinaryReader): UpdateAppRequest;
}

export namespace UpdateAppRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    name: string,
  }
}

export class UpdateAppResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAppResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppResponse): UpdateAppResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppResponse;
  static deserializeBinaryFromReader(message: UpdateAppResponse, reader: jspb.BinaryReader): UpdateAppResponse;
}

export namespace UpdateAppResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOIDCAppConfigRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateOIDCAppConfigRequest;

  getAppId(): string;
  setAppId(value: string): UpdateOIDCAppConfigRequest;

  getRedirectUrisList(): Array<string>;
  setRedirectUrisList(value: Array<string>): UpdateOIDCAppConfigRequest;
  clearRedirectUrisList(): UpdateOIDCAppConfigRequest;
  addRedirectUris(value: string, index?: number): UpdateOIDCAppConfigRequest;

  getResponseTypesList(): Array<zitadel_app_pb.OIDCResponseType>;
  setResponseTypesList(value: Array<zitadel_app_pb.OIDCResponseType>): UpdateOIDCAppConfigRequest;
  clearResponseTypesList(): UpdateOIDCAppConfigRequest;
  addResponseTypes(value: zitadel_app_pb.OIDCResponseType, index?: number): UpdateOIDCAppConfigRequest;

  getGrantTypesList(): Array<zitadel_app_pb.OIDCGrantType>;
  setGrantTypesList(value: Array<zitadel_app_pb.OIDCGrantType>): UpdateOIDCAppConfigRequest;
  clearGrantTypesList(): UpdateOIDCAppConfigRequest;
  addGrantTypes(value: zitadel_app_pb.OIDCGrantType, index?: number): UpdateOIDCAppConfigRequest;

  getAppType(): zitadel_app_pb.OIDCAppType;
  setAppType(value: zitadel_app_pb.OIDCAppType): UpdateOIDCAppConfigRequest;

  getAuthMethodType(): zitadel_app_pb.OIDCAuthMethodType;
  setAuthMethodType(value: zitadel_app_pb.OIDCAuthMethodType): UpdateOIDCAppConfigRequest;

  getPostLogoutRedirectUrisList(): Array<string>;
  setPostLogoutRedirectUrisList(value: Array<string>): UpdateOIDCAppConfigRequest;
  clearPostLogoutRedirectUrisList(): UpdateOIDCAppConfigRequest;
  addPostLogoutRedirectUris(value: string, index?: number): UpdateOIDCAppConfigRequest;

  getDevMode(): boolean;
  setDevMode(value: boolean): UpdateOIDCAppConfigRequest;

  getAccessTokenType(): zitadel_app_pb.OIDCTokenType;
  setAccessTokenType(value: zitadel_app_pb.OIDCTokenType): UpdateOIDCAppConfigRequest;

  getAccessTokenRoleAssertion(): boolean;
  setAccessTokenRoleAssertion(value: boolean): UpdateOIDCAppConfigRequest;

  getIdTokenRoleAssertion(): boolean;
  setIdTokenRoleAssertion(value: boolean): UpdateOIDCAppConfigRequest;

  getIdTokenUserinfoAssertion(): boolean;
  setIdTokenUserinfoAssertion(value: boolean): UpdateOIDCAppConfigRequest;

  getClockSkew(): google_protobuf_duration_pb.Duration | undefined;
  setClockSkew(value?: google_protobuf_duration_pb.Duration): UpdateOIDCAppConfigRequest;
  hasClockSkew(): boolean;
  clearClockSkew(): UpdateOIDCAppConfigRequest;

  getAdditionalOriginsList(): Array<string>;
  setAdditionalOriginsList(value: Array<string>): UpdateOIDCAppConfigRequest;
  clearAdditionalOriginsList(): UpdateOIDCAppConfigRequest;
  addAdditionalOrigins(value: string, index?: number): UpdateOIDCAppConfigRequest;

  getSkipNativeAppSuccessPage(): boolean;
  setSkipNativeAppSuccessPage(value: boolean): UpdateOIDCAppConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOIDCAppConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOIDCAppConfigRequest): UpdateOIDCAppConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOIDCAppConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOIDCAppConfigRequest;
  static deserializeBinaryFromReader(message: UpdateOIDCAppConfigRequest, reader: jspb.BinaryReader): UpdateOIDCAppConfigRequest;
}

export namespace UpdateOIDCAppConfigRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    redirectUrisList: Array<string>,
    responseTypesList: Array<zitadel_app_pb.OIDCResponseType>,
    grantTypesList: Array<zitadel_app_pb.OIDCGrantType>,
    appType: zitadel_app_pb.OIDCAppType,
    authMethodType: zitadel_app_pb.OIDCAuthMethodType,
    postLogoutRedirectUrisList: Array<string>,
    devMode: boolean,
    accessTokenType: zitadel_app_pb.OIDCTokenType,
    accessTokenRoleAssertion: boolean,
    idTokenRoleAssertion: boolean,
    idTokenUserinfoAssertion: boolean,
    clockSkew?: google_protobuf_duration_pb.Duration.AsObject,
    additionalOriginsList: Array<string>,
    skipNativeAppSuccessPage: boolean,
  }
}

export class UpdateOIDCAppConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOIDCAppConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOIDCAppConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOIDCAppConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOIDCAppConfigResponse): UpdateOIDCAppConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOIDCAppConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOIDCAppConfigResponse;
  static deserializeBinaryFromReader(message: UpdateOIDCAppConfigResponse, reader: jspb.BinaryReader): UpdateOIDCAppConfigResponse;
}

export namespace UpdateOIDCAppConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateSAMLAppConfigRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateSAMLAppConfigRequest;

  getAppId(): string;
  setAppId(value: string): UpdateSAMLAppConfigRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): UpdateSAMLAppConfigRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): UpdateSAMLAppConfigRequest;

  getMetadataCase(): UpdateSAMLAppConfigRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLAppConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLAppConfigRequest): UpdateSAMLAppConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLAppConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLAppConfigRequest;
  static deserializeBinaryFromReader(message: UpdateSAMLAppConfigRequest, reader: jspb.BinaryReader): UpdateSAMLAppConfigRequest;
}

export namespace UpdateSAMLAppConfigRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 3,
    METADATA_URL = 4,
  }
}

export class UpdateSAMLAppConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSAMLAppConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSAMLAppConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLAppConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLAppConfigResponse): UpdateSAMLAppConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLAppConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLAppConfigResponse;
  static deserializeBinaryFromReader(message: UpdateSAMLAppConfigResponse, reader: jspb.BinaryReader): UpdateSAMLAppConfigResponse;
}

export namespace UpdateSAMLAppConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateAPIAppConfigRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateAPIAppConfigRequest;

  getAppId(): string;
  setAppId(value: string): UpdateAPIAppConfigRequest;

  getAuthMethodType(): zitadel_app_pb.APIAuthMethodType;
  setAuthMethodType(value: zitadel_app_pb.APIAuthMethodType): UpdateAPIAppConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAPIAppConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAPIAppConfigRequest): UpdateAPIAppConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAPIAppConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAPIAppConfigRequest;
  static deserializeBinaryFromReader(message: UpdateAPIAppConfigRequest, reader: jspb.BinaryReader): UpdateAPIAppConfigRequest;
}

export namespace UpdateAPIAppConfigRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    authMethodType: zitadel_app_pb.APIAuthMethodType,
  }
}

export class UpdateAPIAppConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAPIAppConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAPIAppConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAPIAppConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAPIAppConfigResponse): UpdateAPIAppConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAPIAppConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAPIAppConfigResponse;
  static deserializeBinaryFromReader(message: UpdateAPIAppConfigResponse, reader: jspb.BinaryReader): UpdateAPIAppConfigResponse;
}

export namespace UpdateAPIAppConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): DeactivateAppRequest;

  getAppId(): string;
  setAppId(value: string): DeactivateAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateAppRequest): DeactivateAppRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateAppRequest;
  static deserializeBinaryFromReader(message: DeactivateAppRequest, reader: jspb.BinaryReader): DeactivateAppRequest;
}

export namespace DeactivateAppRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class DeactivateAppResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateAppResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateAppResponse): DeactivateAppResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateAppResponse;
  static deserializeBinaryFromReader(message: DeactivateAppResponse, reader: jspb.BinaryReader): DeactivateAppResponse;
}

export namespace DeactivateAppResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ReactivateAppRequest;

  getAppId(): string;
  setAppId(value: string): ReactivateAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateAppRequest): ReactivateAppRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateAppRequest;
  static deserializeBinaryFromReader(message: ReactivateAppRequest, reader: jspb.BinaryReader): ReactivateAppRequest;
}

export namespace ReactivateAppRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class ReactivateAppResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateAppResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateAppResponse): ReactivateAppResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateAppResponse;
  static deserializeBinaryFromReader(message: ReactivateAppResponse, reader: jspb.BinaryReader): ReactivateAppResponse;
}

export namespace ReactivateAppResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveAppRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveAppRequest;

  getAppId(): string;
  setAppId(value: string): RemoveAppRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAppRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAppRequest): RemoveAppRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveAppRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAppRequest;
  static deserializeBinaryFromReader(message: RemoveAppRequest, reader: jspb.BinaryReader): RemoveAppRequest;
}

export namespace RemoveAppRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class RemoveAppResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveAppResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveAppResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAppResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAppResponse): RemoveAppResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveAppResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAppResponse;
  static deserializeBinaryFromReader(message: RemoveAppResponse, reader: jspb.BinaryReader): RemoveAppResponse;
}

export namespace RemoveAppResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RegenerateOIDCClientSecretRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RegenerateOIDCClientSecretRequest;

  getAppId(): string;
  setAppId(value: string): RegenerateOIDCClientSecretRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateOIDCClientSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateOIDCClientSecretRequest): RegenerateOIDCClientSecretRequest.AsObject;
  static serializeBinaryToWriter(message: RegenerateOIDCClientSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateOIDCClientSecretRequest;
  static deserializeBinaryFromReader(message: RegenerateOIDCClientSecretRequest, reader: jspb.BinaryReader): RegenerateOIDCClientSecretRequest;
}

export namespace RegenerateOIDCClientSecretRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class RegenerateOIDCClientSecretResponse extends jspb.Message {
  getClientSecret(): string;
  setClientSecret(value: string): RegenerateOIDCClientSecretResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RegenerateOIDCClientSecretResponse;
  hasDetails(): boolean;
  clearDetails(): RegenerateOIDCClientSecretResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateOIDCClientSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateOIDCClientSecretResponse): RegenerateOIDCClientSecretResponse.AsObject;
  static serializeBinaryToWriter(message: RegenerateOIDCClientSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateOIDCClientSecretResponse;
  static deserializeBinaryFromReader(message: RegenerateOIDCClientSecretResponse, reader: jspb.BinaryReader): RegenerateOIDCClientSecretResponse;
}

export namespace RegenerateOIDCClientSecretResponse {
  export type AsObject = {
    clientSecret: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RegenerateAPIClientSecretRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RegenerateAPIClientSecretRequest;

  getAppId(): string;
  setAppId(value: string): RegenerateAPIClientSecretRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateAPIClientSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateAPIClientSecretRequest): RegenerateAPIClientSecretRequest.AsObject;
  static serializeBinaryToWriter(message: RegenerateAPIClientSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateAPIClientSecretRequest;
  static deserializeBinaryFromReader(message: RegenerateAPIClientSecretRequest, reader: jspb.BinaryReader): RegenerateAPIClientSecretRequest;
}

export namespace RegenerateAPIClientSecretRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
  }
}

export class RegenerateAPIClientSecretResponse extends jspb.Message {
  getClientSecret(): string;
  setClientSecret(value: string): RegenerateAPIClientSecretResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RegenerateAPIClientSecretResponse;
  hasDetails(): boolean;
  clearDetails(): RegenerateAPIClientSecretResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateAPIClientSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateAPIClientSecretResponse): RegenerateAPIClientSecretResponse.AsObject;
  static serializeBinaryToWriter(message: RegenerateAPIClientSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateAPIClientSecretResponse;
  static deserializeBinaryFromReader(message: RegenerateAPIClientSecretResponse, reader: jspb.BinaryReader): RegenerateAPIClientSecretResponse;
}

export namespace RegenerateAPIClientSecretResponse {
  export type AsObject = {
    clientSecret: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetAppKeyRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): GetAppKeyRequest;

  getAppId(): string;
  setAppId(value: string): GetAppKeyRequest;

  getKeyId(): string;
  setKeyId(value: string): GetAppKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAppKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAppKeyRequest): GetAppKeyRequest.AsObject;
  static serializeBinaryToWriter(message: GetAppKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAppKeyRequest;
  static deserializeBinaryFromReader(message: GetAppKeyRequest, reader: jspb.BinaryReader): GetAppKeyRequest;
}

export namespace GetAppKeyRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    keyId: string,
  }
}

export class GetAppKeyResponse extends jspb.Message {
  getKey(): zitadel_auth_n_key_pb.Key | undefined;
  setKey(value?: zitadel_auth_n_key_pb.Key): GetAppKeyResponse;
  hasKey(): boolean;
  clearKey(): GetAppKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAppKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetAppKeyResponse): GetAppKeyResponse.AsObject;
  static serializeBinaryToWriter(message: GetAppKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAppKeyResponse;
  static deserializeBinaryFromReader(message: GetAppKeyResponse, reader: jspb.BinaryReader): GetAppKeyResponse;
}

export namespace GetAppKeyResponse {
  export type AsObject = {
    key?: zitadel_auth_n_key_pb.Key.AsObject,
  }
}

export class ListAppKeysRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListAppKeysRequest;
  hasQuery(): boolean;
  clearQuery(): ListAppKeysRequest;

  getAppId(): string;
  setAppId(value: string): ListAppKeysRequest;

  getProjectId(): string;
  setProjectId(value: string): ListAppKeysRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppKeysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppKeysRequest): ListAppKeysRequest.AsObject;
  static serializeBinaryToWriter(message: ListAppKeysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppKeysRequest;
  static deserializeBinaryFromReader(message: ListAppKeysRequest, reader: jspb.BinaryReader): ListAppKeysRequest;
}

export namespace ListAppKeysRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    appId: string,
    projectId: string,
  }
}

export class ListAppKeysResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListAppKeysResponse;
  hasDetails(): boolean;
  clearDetails(): ListAppKeysResponse;

  getResultList(): Array<zitadel_auth_n_key_pb.Key>;
  setResultList(value: Array<zitadel_auth_n_key_pb.Key>): ListAppKeysResponse;
  clearResultList(): ListAppKeysResponse;
  addResult(value?: zitadel_auth_n_key_pb.Key, index?: number): zitadel_auth_n_key_pb.Key;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAppKeysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAppKeysResponse): ListAppKeysResponse.AsObject;
  static serializeBinaryToWriter(message: ListAppKeysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAppKeysResponse;
  static deserializeBinaryFromReader(message: ListAppKeysResponse, reader: jspb.BinaryReader): ListAppKeysResponse;
}

export namespace ListAppKeysResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_auth_n_key_pb.Key.AsObject>,
  }
}

export class AddAppKeyRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddAppKeyRequest;

  getAppId(): string;
  setAppId(value: string): AddAppKeyRequest;

  getType(): zitadel_auth_n_key_pb.KeyType;
  setType(value: zitadel_auth_n_key_pb.KeyType): AddAppKeyRequest;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): AddAppKeyRequest;
  hasExpirationDate(): boolean;
  clearExpirationDate(): AddAppKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppKeyRequest): AddAppKeyRequest.AsObject;
  static serializeBinaryToWriter(message: AddAppKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppKeyRequest;
  static deserializeBinaryFromReader(message: AddAppKeyRequest, reader: jspb.BinaryReader): AddAppKeyRequest;
}

export namespace AddAppKeyRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    type: zitadel_auth_n_key_pb.KeyType,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class AddAppKeyResponse extends jspb.Message {
  getId(): string;
  setId(value: string): AddAppKeyResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAppKeyResponse;
  hasDetails(): boolean;
  clearDetails(): AddAppKeyResponse;

  getKeyDetails(): Uint8Array | string;
  getKeyDetails_asU8(): Uint8Array;
  getKeyDetails_asB64(): string;
  setKeyDetails(value: Uint8Array | string): AddAppKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppKeyResponse): AddAppKeyResponse.AsObject;
  static serializeBinaryToWriter(message: AddAppKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppKeyResponse;
  static deserializeBinaryFromReader(message: AddAppKeyResponse, reader: jspb.BinaryReader): AddAppKeyResponse;
}

export namespace AddAppKeyResponse {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    keyDetails: Uint8Array | string,
  }
}

export class RemoveAppKeyRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveAppKeyRequest;

  getAppId(): string;
  setAppId(value: string): RemoveAppKeyRequest;

  getKeyId(): string;
  setKeyId(value: string): RemoveAppKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAppKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAppKeyRequest): RemoveAppKeyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveAppKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAppKeyRequest;
  static deserializeBinaryFromReader(message: RemoveAppKeyRequest, reader: jspb.BinaryReader): RemoveAppKeyRequest;
}

export namespace RemoveAppKeyRequest {
  export type AsObject = {
    projectId: string,
    appId: string,
    keyId: string,
  }
}

export class RemoveAppKeyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveAppKeyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveAppKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveAppKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveAppKeyResponse): RemoveAppKeyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveAppKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveAppKeyResponse;
  static deserializeBinaryFromReader(message: RemoveAppKeyResponse, reader: jspb.BinaryReader): RemoveAppKeyResponse;
}

export namespace RemoveAppKeyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProjectGrantChangesRequest extends jspb.Message {
  getQuery(): zitadel_change_pb.ChangeQuery | undefined;
  setQuery(value?: zitadel_change_pb.ChangeQuery): ListProjectGrantChangesRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectGrantChangesRequest;

  getProjectId(): string;
  setProjectId(value: string): ListProjectGrantChangesRequest;

  getGrantId(): string;
  setGrantId(value: string): ListProjectGrantChangesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantChangesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantChangesRequest): ListProjectGrantChangesRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantChangesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantChangesRequest;
  static deserializeBinaryFromReader(message: ListProjectGrantChangesRequest, reader: jspb.BinaryReader): ListProjectGrantChangesRequest;
}

export namespace ListProjectGrantChangesRequest {
  export type AsObject = {
    query?: zitadel_change_pb.ChangeQuery.AsObject,
    projectId: string,
    grantId: string,
  }
}

export class ListProjectGrantChangesResponse extends jspb.Message {
  getResultList(): Array<zitadel_change_pb.Change>;
  setResultList(value: Array<zitadel_change_pb.Change>): ListProjectGrantChangesResponse;
  clearResultList(): ListProjectGrantChangesResponse;
  addResult(value?: zitadel_change_pb.Change, index?: number): zitadel_change_pb.Change;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantChangesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantChangesResponse): ListProjectGrantChangesResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantChangesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantChangesResponse;
  static deserializeBinaryFromReader(message: ListProjectGrantChangesResponse, reader: jspb.BinaryReader): ListProjectGrantChangesResponse;
}

export namespace ListProjectGrantChangesResponse {
  export type AsObject = {
    resultList: Array<zitadel_change_pb.Change.AsObject>,
  }
}

export class GetProjectGrantByIDRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): GetProjectGrantByIDRequest;

  getGrantId(): string;
  setGrantId(value: string): GetProjectGrantByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectGrantByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectGrantByIDRequest): GetProjectGrantByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetProjectGrantByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectGrantByIDRequest;
  static deserializeBinaryFromReader(message: GetProjectGrantByIDRequest, reader: jspb.BinaryReader): GetProjectGrantByIDRequest;
}

export namespace GetProjectGrantByIDRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
  }
}

export class GetProjectGrantByIDResponse extends jspb.Message {
  getProjectGrant(): zitadel_project_pb.GrantedProject | undefined;
  setProjectGrant(value?: zitadel_project_pb.GrantedProject): GetProjectGrantByIDResponse;
  hasProjectGrant(): boolean;
  clearProjectGrant(): GetProjectGrantByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectGrantByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectGrantByIDResponse): GetProjectGrantByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetProjectGrantByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectGrantByIDResponse;
  static deserializeBinaryFromReader(message: GetProjectGrantByIDResponse, reader: jspb.BinaryReader): GetProjectGrantByIDResponse;
}

export namespace GetProjectGrantByIDResponse {
  export type AsObject = {
    projectGrant?: zitadel_project_pb.GrantedProject.AsObject,
  }
}

export class ListProjectGrantsRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListProjectGrantsRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectGrantsRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectGrantsRequest;

  getQueriesList(): Array<zitadel_project_pb.ProjectGrantQuery>;
  setQueriesList(value: Array<zitadel_project_pb.ProjectGrantQuery>): ListProjectGrantsRequest;
  clearQueriesList(): ListProjectGrantsRequest;
  addQueries(value?: zitadel_project_pb.ProjectGrantQuery, index?: number): zitadel_project_pb.ProjectGrantQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantsRequest): ListProjectGrantsRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantsRequest;
  static deserializeBinaryFromReader(message: ListProjectGrantsRequest, reader: jspb.BinaryReader): ListProjectGrantsRequest;
}

export namespace ListProjectGrantsRequest {
  export type AsObject = {
    projectId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.ProjectGrantQuery.AsObject>,
  }
}

export class ListProjectGrantsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectGrantsResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectGrantsResponse;

  getResultList(): Array<zitadel_project_pb.GrantedProject>;
  setResultList(value: Array<zitadel_project_pb.GrantedProject>): ListProjectGrantsResponse;
  clearResultList(): ListProjectGrantsResponse;
  addResult(value?: zitadel_project_pb.GrantedProject, index?: number): zitadel_project_pb.GrantedProject;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantsResponse): ListProjectGrantsResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantsResponse;
  static deserializeBinaryFromReader(message: ListProjectGrantsResponse, reader: jspb.BinaryReader): ListProjectGrantsResponse;
}

export namespace ListProjectGrantsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.GrantedProject.AsObject>,
  }
}

export class ListAllProjectGrantsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListAllProjectGrantsRequest;
  hasQuery(): boolean;
  clearQuery(): ListAllProjectGrantsRequest;

  getQueriesList(): Array<zitadel_project_pb.AllProjectGrantQuery>;
  setQueriesList(value: Array<zitadel_project_pb.AllProjectGrantQuery>): ListAllProjectGrantsRequest;
  clearQueriesList(): ListAllProjectGrantsRequest;
  addQueries(value?: zitadel_project_pb.AllProjectGrantQuery, index?: number): zitadel_project_pb.AllProjectGrantQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAllProjectGrantsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAllProjectGrantsRequest): ListAllProjectGrantsRequest.AsObject;
  static serializeBinaryToWriter(message: ListAllProjectGrantsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAllProjectGrantsRequest;
  static deserializeBinaryFromReader(message: ListAllProjectGrantsRequest, reader: jspb.BinaryReader): ListAllProjectGrantsRequest;
}

export namespace ListAllProjectGrantsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_project_pb.AllProjectGrantQuery.AsObject>,
  }
}

export class ListAllProjectGrantsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListAllProjectGrantsResponse;
  hasDetails(): boolean;
  clearDetails(): ListAllProjectGrantsResponse;

  getResultList(): Array<zitadel_project_pb.GrantedProject>;
  setResultList(value: Array<zitadel_project_pb.GrantedProject>): ListAllProjectGrantsResponse;
  clearResultList(): ListAllProjectGrantsResponse;
  addResult(value?: zitadel_project_pb.GrantedProject, index?: number): zitadel_project_pb.GrantedProject;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAllProjectGrantsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAllProjectGrantsResponse): ListAllProjectGrantsResponse.AsObject;
  static serializeBinaryToWriter(message: ListAllProjectGrantsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAllProjectGrantsResponse;
  static deserializeBinaryFromReader(message: ListAllProjectGrantsResponse, reader: jspb.BinaryReader): ListAllProjectGrantsResponse;
}

export namespace ListAllProjectGrantsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_project_pb.GrantedProject.AsObject>,
  }
}

export class AddProjectGrantRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddProjectGrantRequest;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): AddProjectGrantRequest;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): AddProjectGrantRequest;
  clearRoleKeysList(): AddProjectGrantRequest;
  addRoleKeys(value: string, index?: number): AddProjectGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectGrantRequest): AddProjectGrantRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectGrantRequest;
  static deserializeBinaryFromReader(message: AddProjectGrantRequest, reader: jspb.BinaryReader): AddProjectGrantRequest;
}

export namespace AddProjectGrantRequest {
  export type AsObject = {
    projectId: string,
    grantedOrgId: string,
    roleKeysList: Array<string>,
  }
}

export class AddProjectGrantResponse extends jspb.Message {
  getGrantId(): string;
  setGrantId(value: string): AddProjectGrantResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddProjectGrantResponse;
  hasDetails(): boolean;
  clearDetails(): AddProjectGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectGrantResponse): AddProjectGrantResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectGrantResponse;
  static deserializeBinaryFromReader(message: AddProjectGrantResponse, reader: jspb.BinaryReader): AddProjectGrantResponse;
}

export namespace AddProjectGrantResponse {
  export type AsObject = {
    grantId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateProjectGrantRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateProjectGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): UpdateProjectGrantRequest;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): UpdateProjectGrantRequest;
  clearRoleKeysList(): UpdateProjectGrantRequest;
  addRoleKeys(value: string, index?: number): UpdateProjectGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectGrantRequest): UpdateProjectGrantRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectGrantRequest;
  static deserializeBinaryFromReader(message: UpdateProjectGrantRequest, reader: jspb.BinaryReader): UpdateProjectGrantRequest;
}

export namespace UpdateProjectGrantRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    roleKeysList: Array<string>,
  }
}

export class UpdateProjectGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateProjectGrantResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateProjectGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectGrantResponse): UpdateProjectGrantResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectGrantResponse;
  static deserializeBinaryFromReader(message: UpdateProjectGrantResponse, reader: jspb.BinaryReader): UpdateProjectGrantResponse;
}

export namespace UpdateProjectGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateProjectGrantRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): DeactivateProjectGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): DeactivateProjectGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateProjectGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateProjectGrantRequest): DeactivateProjectGrantRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateProjectGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateProjectGrantRequest;
  static deserializeBinaryFromReader(message: DeactivateProjectGrantRequest, reader: jspb.BinaryReader): DeactivateProjectGrantRequest;
}

export namespace DeactivateProjectGrantRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
  }
}

export class DeactivateProjectGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateProjectGrantResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateProjectGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateProjectGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateProjectGrantResponse): DeactivateProjectGrantResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateProjectGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateProjectGrantResponse;
  static deserializeBinaryFromReader(message: DeactivateProjectGrantResponse, reader: jspb.BinaryReader): DeactivateProjectGrantResponse;
}

export namespace DeactivateProjectGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateProjectGrantRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ReactivateProjectGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): ReactivateProjectGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateProjectGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateProjectGrantRequest): ReactivateProjectGrantRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateProjectGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateProjectGrantRequest;
  static deserializeBinaryFromReader(message: ReactivateProjectGrantRequest, reader: jspb.BinaryReader): ReactivateProjectGrantRequest;
}

export namespace ReactivateProjectGrantRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
  }
}

export class ReactivateProjectGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateProjectGrantResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateProjectGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateProjectGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateProjectGrantResponse): ReactivateProjectGrantResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateProjectGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateProjectGrantResponse;
  static deserializeBinaryFromReader(message: ReactivateProjectGrantResponse, reader: jspb.BinaryReader): ReactivateProjectGrantResponse;
}

export namespace ReactivateProjectGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveProjectGrantRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveProjectGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): RemoveProjectGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectGrantRequest): RemoveProjectGrantRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectGrantRequest;
  static deserializeBinaryFromReader(message: RemoveProjectGrantRequest, reader: jspb.BinaryReader): RemoveProjectGrantRequest;
}

export namespace RemoveProjectGrantRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
  }
}

export class RemoveProjectGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveProjectGrantResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveProjectGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectGrantResponse): RemoveProjectGrantResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectGrantResponse;
  static deserializeBinaryFromReader(message: RemoveProjectGrantResponse, reader: jspb.BinaryReader): RemoveProjectGrantResponse;
}

export namespace RemoveProjectGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProjectGrantMemberRolesRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectGrantMemberRolesRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectGrantMemberRolesRequest;

  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListProjectGrantMemberRolesRequest;
  clearResultList(): ListProjectGrantMemberRolesRequest;
  addResult(value: string, index?: number): ListProjectGrantMemberRolesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantMemberRolesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantMemberRolesRequest): ListProjectGrantMemberRolesRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantMemberRolesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantMemberRolesRequest;
  static deserializeBinaryFromReader(message: ListProjectGrantMemberRolesRequest, reader: jspb.BinaryReader): ListProjectGrantMemberRolesRequest;
}

export namespace ListProjectGrantMemberRolesRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    resultList: Array<string>,
  }
}

export class ListProjectGrantMemberRolesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectGrantMemberRolesResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectGrantMemberRolesResponse;

  getResultList(): Array<string>;
  setResultList(value: Array<string>): ListProjectGrantMemberRolesResponse;
  clearResultList(): ListProjectGrantMemberRolesResponse;
  addResult(value: string, index?: number): ListProjectGrantMemberRolesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantMemberRolesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantMemberRolesResponse): ListProjectGrantMemberRolesResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantMemberRolesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantMemberRolesResponse;
  static deserializeBinaryFromReader(message: ListProjectGrantMemberRolesResponse, reader: jspb.BinaryReader): ListProjectGrantMemberRolesResponse;
}

export namespace ListProjectGrantMemberRolesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<string>,
  }
}

export class ListProjectGrantMembersRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListProjectGrantMembersRequest;

  getGrantId(): string;
  setGrantId(value: string): ListProjectGrantMembersRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProjectGrantMembersRequest;
  hasQuery(): boolean;
  clearQuery(): ListProjectGrantMembersRequest;

  getQueriesList(): Array<zitadel_member_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_member_pb.SearchQuery>): ListProjectGrantMembersRequest;
  clearQueriesList(): ListProjectGrantMembersRequest;
  addQueries(value?: zitadel_member_pb.SearchQuery, index?: number): zitadel_member_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantMembersRequest): ListProjectGrantMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantMembersRequest;
  static deserializeBinaryFromReader(message: ListProjectGrantMembersRequest, reader: jspb.BinaryReader): ListProjectGrantMembersRequest;
}

export namespace ListProjectGrantMembersRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_member_pb.SearchQuery.AsObject>,
  }
}

export class ListProjectGrantMembersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProjectGrantMembersResponse;
  hasDetails(): boolean;
  clearDetails(): ListProjectGrantMembersResponse;

  getResultList(): Array<zitadel_member_pb.Member>;
  setResultList(value: Array<zitadel_member_pb.Member>): ListProjectGrantMembersResponse;
  clearResultList(): ListProjectGrantMembersResponse;
  addResult(value?: zitadel_member_pb.Member, index?: number): zitadel_member_pb.Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProjectGrantMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProjectGrantMembersResponse): ListProjectGrantMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListProjectGrantMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProjectGrantMembersResponse;
  static deserializeBinaryFromReader(message: ListProjectGrantMembersResponse, reader: jspb.BinaryReader): ListProjectGrantMembersResponse;
}

export namespace ListProjectGrantMembersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_member_pb.Member.AsObject>,
  }
}

export class AddProjectGrantMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): AddProjectGrantMemberRequest;

  getGrantId(): string;
  setGrantId(value: string): AddProjectGrantMemberRequest;

  getUserId(): string;
  setUserId(value: string): AddProjectGrantMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): AddProjectGrantMemberRequest;
  clearRolesList(): AddProjectGrantMemberRequest;
  addRoles(value: string, index?: number): AddProjectGrantMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectGrantMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectGrantMemberRequest): AddProjectGrantMemberRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectGrantMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectGrantMemberRequest;
  static deserializeBinaryFromReader(message: AddProjectGrantMemberRequest, reader: jspb.BinaryReader): AddProjectGrantMemberRequest;
}

export namespace AddProjectGrantMemberRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class AddProjectGrantMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddProjectGrantMemberResponse;
  hasDetails(): boolean;
  clearDetails(): AddProjectGrantMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectGrantMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectGrantMemberResponse): AddProjectGrantMemberResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectGrantMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectGrantMemberResponse;
  static deserializeBinaryFromReader(message: AddProjectGrantMemberResponse, reader: jspb.BinaryReader): AddProjectGrantMemberResponse;
}

export namespace AddProjectGrantMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateProjectGrantMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): UpdateProjectGrantMemberRequest;

  getGrantId(): string;
  setGrantId(value: string): UpdateProjectGrantMemberRequest;

  getUserId(): string;
  setUserId(value: string): UpdateProjectGrantMemberRequest;

  getRolesList(): Array<string>;
  setRolesList(value: Array<string>): UpdateProjectGrantMemberRequest;
  clearRolesList(): UpdateProjectGrantMemberRequest;
  addRoles(value: string, index?: number): UpdateProjectGrantMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectGrantMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectGrantMemberRequest): UpdateProjectGrantMemberRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectGrantMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectGrantMemberRequest;
  static deserializeBinaryFromReader(message: UpdateProjectGrantMemberRequest, reader: jspb.BinaryReader): UpdateProjectGrantMemberRequest;
}

export namespace UpdateProjectGrantMemberRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
    rolesList: Array<string>,
  }
}

export class UpdateProjectGrantMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateProjectGrantMemberResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateProjectGrantMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectGrantMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectGrantMemberResponse): UpdateProjectGrantMemberResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectGrantMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectGrantMemberResponse;
  static deserializeBinaryFromReader(message: UpdateProjectGrantMemberResponse, reader: jspb.BinaryReader): UpdateProjectGrantMemberResponse;
}

export namespace UpdateProjectGrantMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveProjectGrantMemberRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): RemoveProjectGrantMemberRequest;

  getGrantId(): string;
  setGrantId(value: string): RemoveProjectGrantMemberRequest;

  getUserId(): string;
  setUserId(value: string): RemoveProjectGrantMemberRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectGrantMemberRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectGrantMemberRequest): RemoveProjectGrantMemberRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectGrantMemberRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectGrantMemberRequest;
  static deserializeBinaryFromReader(message: RemoveProjectGrantMemberRequest, reader: jspb.BinaryReader): RemoveProjectGrantMemberRequest;
}

export namespace RemoveProjectGrantMemberRequest {
  export type AsObject = {
    projectId: string,
    grantId: string,
    userId: string,
  }
}

export class RemoveProjectGrantMemberResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveProjectGrantMemberResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveProjectGrantMemberResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveProjectGrantMemberResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveProjectGrantMemberResponse): RemoveProjectGrantMemberResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveProjectGrantMemberResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveProjectGrantMemberResponse;
  static deserializeBinaryFromReader(message: RemoveProjectGrantMemberResponse, reader: jspb.BinaryReader): RemoveProjectGrantMemberResponse;
}

export namespace RemoveProjectGrantMemberResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetUserGrantByIDRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): GetUserGrantByIDRequest;

  getGrantId(): string;
  setGrantId(value: string): GetUserGrantByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserGrantByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserGrantByIDRequest): GetUserGrantByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetUserGrantByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserGrantByIDRequest;
  static deserializeBinaryFromReader(message: GetUserGrantByIDRequest, reader: jspb.BinaryReader): GetUserGrantByIDRequest;
}

export namespace GetUserGrantByIDRequest {
  export type AsObject = {
    userId: string,
    grantId: string,
  }
}

export class GetUserGrantByIDResponse extends jspb.Message {
  getUserGrant(): zitadel_user_pb.UserGrant | undefined;
  setUserGrant(value?: zitadel_user_pb.UserGrant): GetUserGrantByIDResponse;
  hasUserGrant(): boolean;
  clearUserGrant(): GetUserGrantByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUserGrantByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetUserGrantByIDResponse): GetUserGrantByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetUserGrantByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUserGrantByIDResponse;
  static deserializeBinaryFromReader(message: GetUserGrantByIDResponse, reader: jspb.BinaryReader): GetUserGrantByIDResponse;
}

export namespace GetUserGrantByIDResponse {
  export type AsObject = {
    userGrant?: zitadel_user_pb.UserGrant.AsObject,
  }
}

export class ListUserGrantRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListUserGrantRequest;
  hasQuery(): boolean;
  clearQuery(): ListUserGrantRequest;

  getQueriesList(): Array<zitadel_user_pb.UserGrantQuery>;
  setQueriesList(value: Array<zitadel_user_pb.UserGrantQuery>): ListUserGrantRequest;
  clearQueriesList(): ListUserGrantRequest;
  addQueries(value?: zitadel_user_pb.UserGrantQuery, index?: number): zitadel_user_pb.UserGrantQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserGrantRequest): ListUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: ListUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserGrantRequest;
  static deserializeBinaryFromReader(message: ListUserGrantRequest, reader: jspb.BinaryReader): ListUserGrantRequest;
}

export namespace ListUserGrantRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<zitadel_user_pb.UserGrantQuery.AsObject>,
  }
}

export class ListUserGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): ListUserGrantResponse;

  getResultList(): Array<zitadel_user_pb.UserGrant>;
  setResultList(value: Array<zitadel_user_pb.UserGrant>): ListUserGrantResponse;
  clearResultList(): ListUserGrantResponse;
  addResult(value?: zitadel_user_pb.UserGrant, index?: number): zitadel_user_pb.UserGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUserGrantResponse): ListUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: ListUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUserGrantResponse;
  static deserializeBinaryFromReader(message: ListUserGrantResponse, reader: jspb.BinaryReader): ListUserGrantResponse;
}

export namespace ListUserGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_user_pb.UserGrant.AsObject>,
  }
}

export class AddUserGrantRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): AddUserGrantRequest;

  getProjectId(): string;
  setProjectId(value: string): AddUserGrantRequest;

  getProjectGrantId(): string;
  setProjectGrantId(value: string): AddUserGrantRequest;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): AddUserGrantRequest;
  clearRoleKeysList(): AddUserGrantRequest;
  addRoleKeys(value: string, index?: number): AddUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserGrantRequest): AddUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: AddUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserGrantRequest;
  static deserializeBinaryFromReader(message: AddUserGrantRequest, reader: jspb.BinaryReader): AddUserGrantRequest;
}

export namespace AddUserGrantRequest {
  export type AsObject = {
    userId: string,
    projectId: string,
    projectGrantId: string,
    roleKeysList: Array<string>,
  }
}

export class AddUserGrantResponse extends jspb.Message {
  getUserGrantId(): string;
  setUserGrantId(value: string): AddUserGrantResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): AddUserGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddUserGrantResponse): AddUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: AddUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddUserGrantResponse;
  static deserializeBinaryFromReader(message: AddUserGrantResponse, reader: jspb.BinaryReader): AddUserGrantResponse;
}

export namespace AddUserGrantResponse {
  export type AsObject = {
    userGrantId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateUserGrantRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UpdateUserGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): UpdateUserGrantRequest;

  getRoleKeysList(): Array<string>;
  setRoleKeysList(value: Array<string>): UpdateUserGrantRequest;
  clearRoleKeysList(): UpdateUserGrantRequest;
  addRoleKeys(value: string, index?: number): UpdateUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserGrantRequest): UpdateUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserGrantRequest;
  static deserializeBinaryFromReader(message: UpdateUserGrantRequest, reader: jspb.BinaryReader): UpdateUserGrantRequest;
}

export namespace UpdateUserGrantRequest {
  export type AsObject = {
    userId: string,
    grantId: string,
    roleKeysList: Array<string>,
  }
}

export class UpdateUserGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateUserGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateUserGrantResponse): UpdateUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateUserGrantResponse;
  static deserializeBinaryFromReader(message: UpdateUserGrantResponse, reader: jspb.BinaryReader): UpdateUserGrantResponse;
}

export namespace UpdateUserGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeactivateUserGrantRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): DeactivateUserGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): DeactivateUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserGrantRequest): DeactivateUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserGrantRequest;
  static deserializeBinaryFromReader(message: DeactivateUserGrantRequest, reader: jspb.BinaryReader): DeactivateUserGrantRequest;
}

export namespace DeactivateUserGrantRequest {
  export type AsObject = {
    userId: string,
    grantId: string,
  }
}

export class DeactivateUserGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateUserGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateUserGrantResponse): DeactivateUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateUserGrantResponse;
  static deserializeBinaryFromReader(message: DeactivateUserGrantResponse, reader: jspb.BinaryReader): DeactivateUserGrantResponse;
}

export namespace DeactivateUserGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateUserGrantRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): ReactivateUserGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): ReactivateUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserGrantRequest): ReactivateUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserGrantRequest;
  static deserializeBinaryFromReader(message: ReactivateUserGrantRequest, reader: jspb.BinaryReader): ReactivateUserGrantRequest;
}

export namespace ReactivateUserGrantRequest {
  export type AsObject = {
    userId: string,
    grantId: string,
  }
}

export class ReactivateUserGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateUserGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateUserGrantResponse): ReactivateUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateUserGrantResponse;
  static deserializeBinaryFromReader(message: ReactivateUserGrantResponse, reader: jspb.BinaryReader): ReactivateUserGrantResponse;
}

export namespace ReactivateUserGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveUserGrantRequest extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): RemoveUserGrantRequest;

  getGrantId(): string;
  setGrantId(value: string): RemoveUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserGrantRequest): RemoveUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserGrantRequest;
  static deserializeBinaryFromReader(message: RemoveUserGrantRequest, reader: jspb.BinaryReader): RemoveUserGrantRequest;
}

export namespace RemoveUserGrantRequest {
  export type AsObject = {
    userId: string,
    grantId: string,
  }
}

export class RemoveUserGrantResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveUserGrantResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveUserGrantResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveUserGrantResponse): RemoveUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveUserGrantResponse;
  static deserializeBinaryFromReader(message: RemoveUserGrantResponse, reader: jspb.BinaryReader): RemoveUserGrantResponse;
}

export namespace RemoveUserGrantResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkRemoveUserGrantRequest extends jspb.Message {
  getGrantIdList(): Array<string>;
  setGrantIdList(value: Array<string>): BulkRemoveUserGrantRequest;
  clearGrantIdList(): BulkRemoveUserGrantRequest;
  addGrantId(value: string, index?: number): BulkRemoveUserGrantRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveUserGrantRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveUserGrantRequest): BulkRemoveUserGrantRequest.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveUserGrantRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveUserGrantRequest;
  static deserializeBinaryFromReader(message: BulkRemoveUserGrantRequest, reader: jspb.BinaryReader): BulkRemoveUserGrantRequest;
}

export namespace BulkRemoveUserGrantRequest {
  export type AsObject = {
    grantIdList: Array<string>,
  }
}

export class BulkRemoveUserGrantResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkRemoveUserGrantResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkRemoveUserGrantResponse): BulkRemoveUserGrantResponse.AsObject;
  static serializeBinaryToWriter(message: BulkRemoveUserGrantResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkRemoveUserGrantResponse;
  static deserializeBinaryFromReader(message: BulkRemoveUserGrantResponse, reader: jspb.BinaryReader): BulkRemoveUserGrantResponse;
}

export namespace BulkRemoveUserGrantResponse {
  export type AsObject = {
  }
}

export class GetOrgIAMPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIAMPolicyRequest): GetOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: GetOrgIAMPolicyRequest, reader: jspb.BinaryReader): GetOrgIAMPolicyRequest;
}

export namespace GetOrgIAMPolicyRequest {
  export type AsObject = {
  }
}

export class GetOrgIAMPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.OrgIAMPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.OrgIAMPolicy): GetOrgIAMPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetOrgIAMPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIAMPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIAMPolicyResponse): GetOrgIAMPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgIAMPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIAMPolicyResponse;
  static deserializeBinaryFromReader(message: GetOrgIAMPolicyResponse, reader: jspb.BinaryReader): GetOrgIAMPolicyResponse;
}

export namespace GetOrgIAMPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.OrgIAMPolicy.AsObject,
  }
}

export class GetDomainPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDomainPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDomainPolicyRequest): GetDomainPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDomainPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDomainPolicyRequest;
  static deserializeBinaryFromReader(message: GetDomainPolicyRequest, reader: jspb.BinaryReader): GetDomainPolicyRequest;
}

export namespace GetDomainPolicyRequest {
  export type AsObject = {
  }
}

export class GetDomainPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.DomainPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.DomainPolicy): GetDomainPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDomainPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDomainPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDomainPolicyResponse): GetDomainPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDomainPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDomainPolicyResponse;
  static deserializeBinaryFromReader(message: GetDomainPolicyResponse, reader: jspb.BinaryReader): GetDomainPolicyResponse;
}

export namespace GetDomainPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.DomainPolicy.AsObject,
  }
}

export class GetLoginPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLoginPolicyRequest): GetLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLoginPolicyRequest;
  static deserializeBinaryFromReader(message: GetLoginPolicyRequest, reader: jspb.BinaryReader): GetLoginPolicyRequest;
}

export namespace GetLoginPolicyRequest {
  export type AsObject = {
  }
}

export class GetLoginPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LoginPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LoginPolicy): GetLoginPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLoginPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLoginPolicyResponse): GetLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLoginPolicyResponse;
  static deserializeBinaryFromReader(message: GetLoginPolicyResponse, reader: jspb.BinaryReader): GetLoginPolicyResponse;
}

export namespace GetLoginPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LoginPolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetDefaultLoginPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginPolicyRequest): GetDefaultLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultLoginPolicyRequest, reader: jspb.BinaryReader): GetDefaultLoginPolicyRequest;
}

export namespace GetDefaultLoginPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultLoginPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LoginPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LoginPolicy): GetDefaultLoginPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginPolicyResponse): GetDefaultLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultLoginPolicyResponse, reader: jspb.BinaryReader): GetDefaultLoginPolicyResponse;
}

export namespace GetDefaultLoginPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LoginPolicy.AsObject,
  }
}

export class AddCustomLoginPolicyRequest extends jspb.Message {
  getAllowUsernamePassword(): boolean;
  setAllowUsernamePassword(value: boolean): AddCustomLoginPolicyRequest;

  getAllowRegister(): boolean;
  setAllowRegister(value: boolean): AddCustomLoginPolicyRequest;

  getAllowExternalIdp(): boolean;
  setAllowExternalIdp(value: boolean): AddCustomLoginPolicyRequest;

  getForceMfa(): boolean;
  setForceMfa(value: boolean): AddCustomLoginPolicyRequest;

  getPasswordlessType(): zitadel_policy_pb.PasswordlessType;
  setPasswordlessType(value: zitadel_policy_pb.PasswordlessType): AddCustomLoginPolicyRequest;

  getHidePasswordReset(): boolean;
  setHidePasswordReset(value: boolean): AddCustomLoginPolicyRequest;

  getIgnoreUnknownUsernames(): boolean;
  setIgnoreUnknownUsernames(value: boolean): AddCustomLoginPolicyRequest;

  getDefaultRedirectUri(): string;
  setDefaultRedirectUri(value: string): AddCustomLoginPolicyRequest;

  getPasswordCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setPasswordCheckLifetime(value?: google_protobuf_duration_pb.Duration): AddCustomLoginPolicyRequest;
  hasPasswordCheckLifetime(): boolean;
  clearPasswordCheckLifetime(): AddCustomLoginPolicyRequest;

  getExternalLoginCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setExternalLoginCheckLifetime(value?: google_protobuf_duration_pb.Duration): AddCustomLoginPolicyRequest;
  hasExternalLoginCheckLifetime(): boolean;
  clearExternalLoginCheckLifetime(): AddCustomLoginPolicyRequest;

  getMfaInitSkipLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMfaInitSkipLifetime(value?: google_protobuf_duration_pb.Duration): AddCustomLoginPolicyRequest;
  hasMfaInitSkipLifetime(): boolean;
  clearMfaInitSkipLifetime(): AddCustomLoginPolicyRequest;

  getSecondFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setSecondFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): AddCustomLoginPolicyRequest;
  hasSecondFactorCheckLifetime(): boolean;
  clearSecondFactorCheckLifetime(): AddCustomLoginPolicyRequest;

  getMultiFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMultiFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): AddCustomLoginPolicyRequest;
  hasMultiFactorCheckLifetime(): boolean;
  clearMultiFactorCheckLifetime(): AddCustomLoginPolicyRequest;

  getSecondFactorsList(): Array<zitadel_policy_pb.SecondFactorType>;
  setSecondFactorsList(value: Array<zitadel_policy_pb.SecondFactorType>): AddCustomLoginPolicyRequest;
  clearSecondFactorsList(): AddCustomLoginPolicyRequest;
  addSecondFactors(value: zitadel_policy_pb.SecondFactorType, index?: number): AddCustomLoginPolicyRequest;

  getMultiFactorsList(): Array<zitadel_policy_pb.MultiFactorType>;
  setMultiFactorsList(value: Array<zitadel_policy_pb.MultiFactorType>): AddCustomLoginPolicyRequest;
  clearMultiFactorsList(): AddCustomLoginPolicyRequest;
  addMultiFactors(value: zitadel_policy_pb.MultiFactorType, index?: number): AddCustomLoginPolicyRequest;

  getIdpsList(): Array<AddCustomLoginPolicyRequest.IDP>;
  setIdpsList(value: Array<AddCustomLoginPolicyRequest.IDP>): AddCustomLoginPolicyRequest;
  clearIdpsList(): AddCustomLoginPolicyRequest;
  addIdps(value?: AddCustomLoginPolicyRequest.IDP, index?: number): AddCustomLoginPolicyRequest.IDP;

  getAllowDomainDiscovery(): boolean;
  setAllowDomainDiscovery(value: boolean): AddCustomLoginPolicyRequest;

  getDisableLoginWithEmail(): boolean;
  setDisableLoginWithEmail(value: boolean): AddCustomLoginPolicyRequest;

  getDisableLoginWithPhone(): boolean;
  setDisableLoginWithPhone(value: boolean): AddCustomLoginPolicyRequest;

  getForceMfaLocalOnly(): boolean;
  setForceMfaLocalOnly(value: boolean): AddCustomLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLoginPolicyRequest): AddCustomLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomLoginPolicyRequest, reader: jspb.BinaryReader): AddCustomLoginPolicyRequest;
}

export namespace AddCustomLoginPolicyRequest {
  export type AsObject = {
    allowUsernamePassword: boolean,
    allowRegister: boolean,
    allowExternalIdp: boolean,
    forceMfa: boolean,
    passwordlessType: zitadel_policy_pb.PasswordlessType,
    hidePasswordReset: boolean,
    ignoreUnknownUsernames: boolean,
    defaultRedirectUri: string,
    passwordCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    externalLoginCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    mfaInitSkipLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    multiFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorsList: Array<zitadel_policy_pb.SecondFactorType>,
    multiFactorsList: Array<zitadel_policy_pb.MultiFactorType>,
    idpsList: Array<AddCustomLoginPolicyRequest.IDP.AsObject>,
    allowDomainDiscovery: boolean,
    disableLoginWithEmail: boolean,
    disableLoginWithPhone: boolean,
    forceMfaLocalOnly: boolean,
  }

  export class IDP extends jspb.Message {
    getIdpId(): string;
    setIdpId(value: string): IDP;

    getOwnertype(): zitadel_idp_pb.IDPOwnerType;
    setOwnertype(value: zitadel_idp_pb.IDPOwnerType): IDP;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): IDP.AsObject;
    static toObject(includeInstance: boolean, msg: IDP): IDP.AsObject;
    static serializeBinaryToWriter(message: IDP, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): IDP;
    static deserializeBinaryFromReader(message: IDP, reader: jspb.BinaryReader): IDP;
  }

  export namespace IDP {
    export type AsObject = {
      idpId: string,
      ownertype: zitadel_idp_pb.IDPOwnerType,
    }
  }

}

export class AddCustomLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLoginPolicyResponse): AddCustomLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomLoginPolicyResponse, reader: jspb.BinaryReader): AddCustomLoginPolicyResponse;
}

export namespace AddCustomLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomLoginPolicyRequest extends jspb.Message {
  getAllowUsernamePassword(): boolean;
  setAllowUsernamePassword(value: boolean): UpdateCustomLoginPolicyRequest;

  getAllowRegister(): boolean;
  setAllowRegister(value: boolean): UpdateCustomLoginPolicyRequest;

  getAllowExternalIdp(): boolean;
  setAllowExternalIdp(value: boolean): UpdateCustomLoginPolicyRequest;

  getForceMfa(): boolean;
  setForceMfa(value: boolean): UpdateCustomLoginPolicyRequest;

  getPasswordlessType(): zitadel_policy_pb.PasswordlessType;
  setPasswordlessType(value: zitadel_policy_pb.PasswordlessType): UpdateCustomLoginPolicyRequest;

  getHidePasswordReset(): boolean;
  setHidePasswordReset(value: boolean): UpdateCustomLoginPolicyRequest;

  getIgnoreUnknownUsernames(): boolean;
  setIgnoreUnknownUsernames(value: boolean): UpdateCustomLoginPolicyRequest;

  getDefaultRedirectUri(): string;
  setDefaultRedirectUri(value: string): UpdateCustomLoginPolicyRequest;

  getPasswordCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setPasswordCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateCustomLoginPolicyRequest;
  hasPasswordCheckLifetime(): boolean;
  clearPasswordCheckLifetime(): UpdateCustomLoginPolicyRequest;

  getExternalLoginCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setExternalLoginCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateCustomLoginPolicyRequest;
  hasExternalLoginCheckLifetime(): boolean;
  clearExternalLoginCheckLifetime(): UpdateCustomLoginPolicyRequest;

  getMfaInitSkipLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMfaInitSkipLifetime(value?: google_protobuf_duration_pb.Duration): UpdateCustomLoginPolicyRequest;
  hasMfaInitSkipLifetime(): boolean;
  clearMfaInitSkipLifetime(): UpdateCustomLoginPolicyRequest;

  getSecondFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setSecondFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateCustomLoginPolicyRequest;
  hasSecondFactorCheckLifetime(): boolean;
  clearSecondFactorCheckLifetime(): UpdateCustomLoginPolicyRequest;

  getMultiFactorCheckLifetime(): google_protobuf_duration_pb.Duration | undefined;
  setMultiFactorCheckLifetime(value?: google_protobuf_duration_pb.Duration): UpdateCustomLoginPolicyRequest;
  hasMultiFactorCheckLifetime(): boolean;
  clearMultiFactorCheckLifetime(): UpdateCustomLoginPolicyRequest;

  getAllowDomainDiscovery(): boolean;
  setAllowDomainDiscovery(value: boolean): UpdateCustomLoginPolicyRequest;

  getDisableLoginWithEmail(): boolean;
  setDisableLoginWithEmail(value: boolean): UpdateCustomLoginPolicyRequest;

  getDisableLoginWithPhone(): boolean;
  setDisableLoginWithPhone(value: boolean): UpdateCustomLoginPolicyRequest;

  getForceMfaLocalOnly(): boolean;
  setForceMfaLocalOnly(value: boolean): UpdateCustomLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLoginPolicyRequest): UpdateCustomLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLoginPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomLoginPolicyRequest, reader: jspb.BinaryReader): UpdateCustomLoginPolicyRequest;
}

export namespace UpdateCustomLoginPolicyRequest {
  export type AsObject = {
    allowUsernamePassword: boolean,
    allowRegister: boolean,
    allowExternalIdp: boolean,
    forceMfa: boolean,
    passwordlessType: zitadel_policy_pb.PasswordlessType,
    hidePasswordReset: boolean,
    ignoreUnknownUsernames: boolean,
    defaultRedirectUri: string,
    passwordCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    externalLoginCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    mfaInitSkipLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    secondFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    multiFactorCheckLifetime?: google_protobuf_duration_pb.Duration.AsObject,
    allowDomainDiscovery: boolean,
    disableLoginWithEmail: boolean,
    disableLoginWithPhone: boolean,
    forceMfaLocalOnly: boolean,
  }
}

export class UpdateCustomLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLoginPolicyResponse): UpdateCustomLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLoginPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomLoginPolicyResponse, reader: jspb.BinaryReader): UpdateCustomLoginPolicyResponse;
}

export namespace UpdateCustomLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetLoginPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLoginPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLoginPolicyToDefaultRequest): ResetLoginPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetLoginPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLoginPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetLoginPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetLoginPolicyToDefaultRequest;
}

export namespace ResetLoginPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetLoginPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetLoginPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetLoginPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLoginPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLoginPolicyToDefaultResponse): ResetLoginPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetLoginPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLoginPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetLoginPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetLoginPolicyToDefaultResponse;
}

export namespace ResetLoginPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicyIDPsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListLoginPolicyIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListLoginPolicyIDPsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyIDPsRequest): ListLoginPolicyIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyIDPsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicyIDPsRequest, reader: jspb.BinaryReader): ListLoginPolicyIDPsRequest;
}

export namespace ListLoginPolicyIDPsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
  }
}

export class ListLoginPolicyIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicyIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicyIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDPLoginPolicyLink>;
  setResultList(value: Array<zitadel_idp_pb.IDPLoginPolicyLink>): ListLoginPolicyIDPsResponse;
  clearResultList(): ListLoginPolicyIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDPLoginPolicyLink, index?: number): zitadel_idp_pb.IDPLoginPolicyLink;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyIDPsResponse): ListLoginPolicyIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyIDPsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicyIDPsResponse, reader: jspb.BinaryReader): ListLoginPolicyIDPsResponse;
}

export namespace ListLoginPolicyIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.IDPLoginPolicyLink.AsObject>,
  }
}

export class AddIDPToLoginPolicyRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): AddIDPToLoginPolicyRequest;

  getOwnertype(): zitadel_idp_pb.IDPOwnerType;
  setOwnertype(value: zitadel_idp_pb.IDPOwnerType): AddIDPToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPToLoginPolicyRequest): AddIDPToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddIDPToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddIDPToLoginPolicyRequest, reader: jspb.BinaryReader): AddIDPToLoginPolicyRequest;
}

export namespace AddIDPToLoginPolicyRequest {
  export type AsObject = {
    idpId: string,
    ownertype: zitadel_idp_pb.IDPOwnerType,
  }
}

export class AddIDPToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddIDPToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddIDPToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddIDPToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddIDPToLoginPolicyResponse): AddIDPToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddIDPToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddIDPToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddIDPToLoginPolicyResponse, reader: jspb.BinaryReader): AddIDPToLoginPolicyResponse;
}

export namespace AddIDPToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveIDPFromLoginPolicyRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): RemoveIDPFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPFromLoginPolicyRequest): RemoveIDPFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveIDPFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveIDPFromLoginPolicyRequest;
}

export namespace RemoveIDPFromLoginPolicyRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class RemoveIDPFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveIDPFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveIDPFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveIDPFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveIDPFromLoginPolicyResponse): RemoveIDPFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveIDPFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveIDPFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveIDPFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveIDPFromLoginPolicyResponse;
}

export namespace RemoveIDPFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicySecondFactorsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicySecondFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicySecondFactorsRequest): ListLoginPolicySecondFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicySecondFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicySecondFactorsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicySecondFactorsRequest, reader: jspb.BinaryReader): ListLoginPolicySecondFactorsRequest;
}

export namespace ListLoginPolicySecondFactorsRequest {
  export type AsObject = {
  }
}

export class ListLoginPolicySecondFactorsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicySecondFactorsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicySecondFactorsResponse;

  getResultList(): Array<zitadel_policy_pb.SecondFactorType>;
  setResultList(value: Array<zitadel_policy_pb.SecondFactorType>): ListLoginPolicySecondFactorsResponse;
  clearResultList(): ListLoginPolicySecondFactorsResponse;
  addResult(value: zitadel_policy_pb.SecondFactorType, index?: number): ListLoginPolicySecondFactorsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicySecondFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicySecondFactorsResponse): ListLoginPolicySecondFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicySecondFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicySecondFactorsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicySecondFactorsResponse, reader: jspb.BinaryReader): ListLoginPolicySecondFactorsResponse;
}

export namespace ListLoginPolicySecondFactorsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_policy_pb.SecondFactorType>,
  }
}

export class AddSecondFactorToLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.SecondFactorType;
  setType(value: zitadel_policy_pb.SecondFactorType): AddSecondFactorToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSecondFactorToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSecondFactorToLoginPolicyRequest): AddSecondFactorToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddSecondFactorToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSecondFactorToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddSecondFactorToLoginPolicyRequest, reader: jspb.BinaryReader): AddSecondFactorToLoginPolicyRequest;
}

export namespace AddSecondFactorToLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.SecondFactorType,
  }
}

export class AddSecondFactorToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSecondFactorToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddSecondFactorToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSecondFactorToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSecondFactorToLoginPolicyResponse): AddSecondFactorToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddSecondFactorToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSecondFactorToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddSecondFactorToLoginPolicyResponse, reader: jspb.BinaryReader): AddSecondFactorToLoginPolicyResponse;
}

export namespace AddSecondFactorToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveSecondFactorFromLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.SecondFactorType;
  setType(value: zitadel_policy_pb.SecondFactorType): RemoveSecondFactorFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSecondFactorFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSecondFactorFromLoginPolicyRequest): RemoveSecondFactorFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveSecondFactorFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSecondFactorFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveSecondFactorFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveSecondFactorFromLoginPolicyRequest;
}

export namespace RemoveSecondFactorFromLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.SecondFactorType,
  }
}

export class RemoveSecondFactorFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveSecondFactorFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveSecondFactorFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveSecondFactorFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveSecondFactorFromLoginPolicyResponse): RemoveSecondFactorFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveSecondFactorFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveSecondFactorFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveSecondFactorFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveSecondFactorFromLoginPolicyResponse;
}

export namespace RemoveSecondFactorFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListLoginPolicyMultiFactorsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyMultiFactorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyMultiFactorsRequest): ListLoginPolicyMultiFactorsRequest.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyMultiFactorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyMultiFactorsRequest;
  static deserializeBinaryFromReader(message: ListLoginPolicyMultiFactorsRequest, reader: jspb.BinaryReader): ListLoginPolicyMultiFactorsRequest;
}

export namespace ListLoginPolicyMultiFactorsRequest {
  export type AsObject = {
  }
}

export class ListLoginPolicyMultiFactorsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListLoginPolicyMultiFactorsResponse;
  hasDetails(): boolean;
  clearDetails(): ListLoginPolicyMultiFactorsResponse;

  getResultList(): Array<zitadel_policy_pb.MultiFactorType>;
  setResultList(value: Array<zitadel_policy_pb.MultiFactorType>): ListLoginPolicyMultiFactorsResponse;
  clearResultList(): ListLoginPolicyMultiFactorsResponse;
  addResult(value: zitadel_policy_pb.MultiFactorType, index?: number): ListLoginPolicyMultiFactorsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListLoginPolicyMultiFactorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListLoginPolicyMultiFactorsResponse): ListLoginPolicyMultiFactorsResponse.AsObject;
  static serializeBinaryToWriter(message: ListLoginPolicyMultiFactorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListLoginPolicyMultiFactorsResponse;
  static deserializeBinaryFromReader(message: ListLoginPolicyMultiFactorsResponse, reader: jspb.BinaryReader): ListLoginPolicyMultiFactorsResponse;
}

export namespace ListLoginPolicyMultiFactorsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_policy_pb.MultiFactorType>,
  }
}

export class AddMultiFactorToLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.MultiFactorType;
  setType(value: zitadel_policy_pb.MultiFactorType): AddMultiFactorToLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMultiFactorToLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddMultiFactorToLoginPolicyRequest): AddMultiFactorToLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddMultiFactorToLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMultiFactorToLoginPolicyRequest;
  static deserializeBinaryFromReader(message: AddMultiFactorToLoginPolicyRequest, reader: jspb.BinaryReader): AddMultiFactorToLoginPolicyRequest;
}

export namespace AddMultiFactorToLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.MultiFactorType,
  }
}

export class AddMultiFactorToLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddMultiFactorToLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddMultiFactorToLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddMultiFactorToLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddMultiFactorToLoginPolicyResponse): AddMultiFactorToLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddMultiFactorToLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddMultiFactorToLoginPolicyResponse;
  static deserializeBinaryFromReader(message: AddMultiFactorToLoginPolicyResponse, reader: jspb.BinaryReader): AddMultiFactorToLoginPolicyResponse;
}

export namespace AddMultiFactorToLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveMultiFactorFromLoginPolicyRequest extends jspb.Message {
  getType(): zitadel_policy_pb.MultiFactorType;
  setType(value: zitadel_policy_pb.MultiFactorType): RemoveMultiFactorFromLoginPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMultiFactorFromLoginPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMultiFactorFromLoginPolicyRequest): RemoveMultiFactorFromLoginPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveMultiFactorFromLoginPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMultiFactorFromLoginPolicyRequest;
  static deserializeBinaryFromReader(message: RemoveMultiFactorFromLoginPolicyRequest, reader: jspb.BinaryReader): RemoveMultiFactorFromLoginPolicyRequest;
}

export namespace RemoveMultiFactorFromLoginPolicyRequest {
  export type AsObject = {
    type: zitadel_policy_pb.MultiFactorType,
  }
}

export class RemoveMultiFactorFromLoginPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveMultiFactorFromLoginPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveMultiFactorFromLoginPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveMultiFactorFromLoginPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveMultiFactorFromLoginPolicyResponse): RemoveMultiFactorFromLoginPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveMultiFactorFromLoginPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveMultiFactorFromLoginPolicyResponse;
  static deserializeBinaryFromReader(message: RemoveMultiFactorFromLoginPolicyResponse, reader: jspb.BinaryReader): RemoveMultiFactorFromLoginPolicyResponse;
}

export namespace RemoveMultiFactorFromLoginPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPasswordComplexityPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordComplexityPolicyRequest): GetPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: GetPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): GetPasswordComplexityPolicyRequest;
}

export namespace GetPasswordComplexityPolicyRequest {
  export type AsObject = {
  }
}

export class GetPasswordComplexityPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordComplexityPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordComplexityPolicy): GetPasswordComplexityPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPasswordComplexityPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordComplexityPolicyResponse): GetPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: GetPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): GetPasswordComplexityPolicyResponse;
}

export namespace GetPasswordComplexityPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordComplexityPolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetDefaultPasswordComplexityPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordComplexityPolicyRequest): GetDefaultPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): GetDefaultPasswordComplexityPolicyRequest;
}

export namespace GetDefaultPasswordComplexityPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultPasswordComplexityPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordComplexityPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordComplexityPolicy): GetDefaultPasswordComplexityPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordComplexityPolicyResponse): GetDefaultPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): GetDefaultPasswordComplexityPolicyResponse;
}

export namespace GetDefaultPasswordComplexityPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordComplexityPolicy.AsObject,
  }
}

export class AddCustomPasswordComplexityPolicyRequest extends jspb.Message {
  getMinLength(): number;
  setMinLength(value: number): AddCustomPasswordComplexityPolicyRequest;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): AddCustomPasswordComplexityPolicyRequest;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): AddCustomPasswordComplexityPolicyRequest;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): AddCustomPasswordComplexityPolicyRequest;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): AddCustomPasswordComplexityPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPasswordComplexityPolicyRequest): AddCustomPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): AddCustomPasswordComplexityPolicyRequest;
}

export namespace AddCustomPasswordComplexityPolicyRequest {
  export type AsObject = {
    minLength: number,
    hasUppercase: boolean,
    hasLowercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
  }
}

export class AddCustomPasswordComplexityPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomPasswordComplexityPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPasswordComplexityPolicyResponse): AddCustomPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): AddCustomPasswordComplexityPolicyResponse;
}

export namespace AddCustomPasswordComplexityPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomPasswordComplexityPolicyRequest extends jspb.Message {
  getMinLength(): number;
  setMinLength(value: number): UpdateCustomPasswordComplexityPolicyRequest;

  getHasUppercase(): boolean;
  setHasUppercase(value: boolean): UpdateCustomPasswordComplexityPolicyRequest;

  getHasLowercase(): boolean;
  setHasLowercase(value: boolean): UpdateCustomPasswordComplexityPolicyRequest;

  getHasNumber(): boolean;
  setHasNumber(value: boolean): UpdateCustomPasswordComplexityPolicyRequest;

  getHasSymbol(): boolean;
  setHasSymbol(value: boolean): UpdateCustomPasswordComplexityPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPasswordComplexityPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPasswordComplexityPolicyRequest): UpdateCustomPasswordComplexityPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPasswordComplexityPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPasswordComplexityPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomPasswordComplexityPolicyRequest, reader: jspb.BinaryReader): UpdateCustomPasswordComplexityPolicyRequest;
}

export namespace UpdateCustomPasswordComplexityPolicyRequest {
  export type AsObject = {
    minLength: number,
    hasUppercase: boolean,
    hasLowercase: boolean,
    hasNumber: boolean,
    hasSymbol: boolean,
  }
}

export class UpdateCustomPasswordComplexityPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomPasswordComplexityPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomPasswordComplexityPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPasswordComplexityPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPasswordComplexityPolicyResponse): UpdateCustomPasswordComplexityPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPasswordComplexityPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPasswordComplexityPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomPasswordComplexityPolicyResponse, reader: jspb.BinaryReader): UpdateCustomPasswordComplexityPolicyResponse;
}

export namespace UpdateCustomPasswordComplexityPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetPasswordComplexityPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordComplexityPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordComplexityPolicyToDefaultRequest): ResetPasswordComplexityPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordComplexityPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordComplexityPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetPasswordComplexityPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetPasswordComplexityPolicyToDefaultRequest;
}

export namespace ResetPasswordComplexityPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetPasswordComplexityPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetPasswordComplexityPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetPasswordComplexityPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordComplexityPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordComplexityPolicyToDefaultResponse): ResetPasswordComplexityPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordComplexityPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordComplexityPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetPasswordComplexityPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetPasswordComplexityPolicyToDefaultResponse;
}

export namespace ResetPasswordComplexityPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPasswordAgePolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordAgePolicyRequest): GetPasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: GetPasswordAgePolicyRequest, reader: jspb.BinaryReader): GetPasswordAgePolicyRequest;
}

export namespace GetPasswordAgePolicyRequest {
  export type AsObject = {
  }
}

export class GetPasswordAgePolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordAgePolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordAgePolicy): GetPasswordAgePolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPasswordAgePolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetPasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPasswordAgePolicyResponse): GetPasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: GetPasswordAgePolicyResponse, reader: jspb.BinaryReader): GetPasswordAgePolicyResponse;
}

export namespace GetPasswordAgePolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordAgePolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetDefaultPasswordAgePolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordAgePolicyRequest): GetDefaultPasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordAgePolicyRequest, reader: jspb.BinaryReader): GetDefaultPasswordAgePolicyRequest;
}

export namespace GetDefaultPasswordAgePolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultPasswordAgePolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PasswordAgePolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PasswordAgePolicy): GetDefaultPasswordAgePolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultPasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordAgePolicyResponse): GetDefaultPasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordAgePolicyResponse, reader: jspb.BinaryReader): GetDefaultPasswordAgePolicyResponse;
}

export namespace GetDefaultPasswordAgePolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PasswordAgePolicy.AsObject,
  }
}

export class AddCustomPasswordAgePolicyRequest extends jspb.Message {
  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): AddCustomPasswordAgePolicyRequest;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): AddCustomPasswordAgePolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPasswordAgePolicyRequest): AddCustomPasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomPasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomPasswordAgePolicyRequest, reader: jspb.BinaryReader): AddCustomPasswordAgePolicyRequest;
}

export namespace AddCustomPasswordAgePolicyRequest {
  export type AsObject = {
    maxAgeDays: number,
    expireWarnDays: number,
  }
}

export class AddCustomPasswordAgePolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomPasswordAgePolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomPasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPasswordAgePolicyResponse): AddCustomPasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomPasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomPasswordAgePolicyResponse, reader: jspb.BinaryReader): AddCustomPasswordAgePolicyResponse;
}

export namespace AddCustomPasswordAgePolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomPasswordAgePolicyRequest extends jspb.Message {
  getMaxAgeDays(): number;
  setMaxAgeDays(value: number): UpdateCustomPasswordAgePolicyRequest;

  getExpireWarnDays(): number;
  setExpireWarnDays(value: number): UpdateCustomPasswordAgePolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPasswordAgePolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPasswordAgePolicyRequest): UpdateCustomPasswordAgePolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPasswordAgePolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPasswordAgePolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomPasswordAgePolicyRequest, reader: jspb.BinaryReader): UpdateCustomPasswordAgePolicyRequest;
}

export namespace UpdateCustomPasswordAgePolicyRequest {
  export type AsObject = {
    maxAgeDays: number,
    expireWarnDays: number,
  }
}

export class UpdateCustomPasswordAgePolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomPasswordAgePolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomPasswordAgePolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPasswordAgePolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPasswordAgePolicyResponse): UpdateCustomPasswordAgePolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPasswordAgePolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPasswordAgePolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomPasswordAgePolicyResponse, reader: jspb.BinaryReader): UpdateCustomPasswordAgePolicyResponse;
}

export namespace UpdateCustomPasswordAgePolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetPasswordAgePolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordAgePolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordAgePolicyToDefaultRequest): ResetPasswordAgePolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordAgePolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordAgePolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetPasswordAgePolicyToDefaultRequest, reader: jspb.BinaryReader): ResetPasswordAgePolicyToDefaultRequest;
}

export namespace ResetPasswordAgePolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetPasswordAgePolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetPasswordAgePolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetPasswordAgePolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPasswordAgePolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPasswordAgePolicyToDefaultResponse): ResetPasswordAgePolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetPasswordAgePolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPasswordAgePolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetPasswordAgePolicyToDefaultResponse, reader: jspb.BinaryReader): ResetPasswordAgePolicyToDefaultResponse;
}

export namespace ResetPasswordAgePolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetLockoutPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLockoutPolicyRequest): GetLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: GetLockoutPolicyRequest, reader: jspb.BinaryReader): GetLockoutPolicyRequest;
}

export namespace GetLockoutPolicyRequest {
  export type AsObject = {
  }
}

export class GetLockoutPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LockoutPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LockoutPolicy): GetLockoutPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLockoutPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLockoutPolicyResponse): GetLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: GetLockoutPolicyResponse, reader: jspb.BinaryReader): GetLockoutPolicyResponse;
}

export namespace GetLockoutPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LockoutPolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetDefaultLockoutPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLockoutPolicyRequest): GetDefaultLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultLockoutPolicyRequest, reader: jspb.BinaryReader): GetDefaultLockoutPolicyRequest;
}

export namespace GetDefaultLockoutPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultLockoutPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LockoutPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LockoutPolicy): GetDefaultLockoutPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLockoutPolicyResponse): GetDefaultLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultLockoutPolicyResponse, reader: jspb.BinaryReader): GetDefaultLockoutPolicyResponse;
}

export namespace GetDefaultLockoutPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LockoutPolicy.AsObject,
  }
}

export class AddCustomLockoutPolicyRequest extends jspb.Message {
  getMaxPasswordAttempts(): number;
  setMaxPasswordAttempts(value: number): AddCustomLockoutPolicyRequest;

  getMaxOtpAttempts(): number;
  setMaxOtpAttempts(value: number): AddCustomLockoutPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLockoutPolicyRequest): AddCustomLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomLockoutPolicyRequest, reader: jspb.BinaryReader): AddCustomLockoutPolicyRequest;
}

export namespace AddCustomLockoutPolicyRequest {
  export type AsObject = {
    maxPasswordAttempts: number,
    maxOtpAttempts: number,
  }
}

export class AddCustomLockoutPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomLockoutPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLockoutPolicyResponse): AddCustomLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomLockoutPolicyResponse, reader: jspb.BinaryReader): AddCustomLockoutPolicyResponse;
}

export namespace AddCustomLockoutPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomLockoutPolicyRequest extends jspb.Message {
  getMaxPasswordAttempts(): number;
  setMaxPasswordAttempts(value: number): UpdateCustomLockoutPolicyRequest;

  getMaxOtpAttempts(): number;
  setMaxOtpAttempts(value: number): UpdateCustomLockoutPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLockoutPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLockoutPolicyRequest): UpdateCustomLockoutPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLockoutPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLockoutPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomLockoutPolicyRequest, reader: jspb.BinaryReader): UpdateCustomLockoutPolicyRequest;
}

export namespace UpdateCustomLockoutPolicyRequest {
  export type AsObject = {
    maxPasswordAttempts: number,
    maxOtpAttempts: number,
  }
}

export class UpdateCustomLockoutPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomLockoutPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomLockoutPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLockoutPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLockoutPolicyResponse): UpdateCustomLockoutPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLockoutPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLockoutPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomLockoutPolicyResponse, reader: jspb.BinaryReader): UpdateCustomLockoutPolicyResponse;
}

export namespace UpdateCustomLockoutPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetLockoutPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLockoutPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLockoutPolicyToDefaultRequest): ResetLockoutPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetLockoutPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLockoutPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetLockoutPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetLockoutPolicyToDefaultRequest;
}

export namespace ResetLockoutPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetLockoutPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetLockoutPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetLockoutPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLockoutPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLockoutPolicyToDefaultResponse): ResetLockoutPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetLockoutPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLockoutPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetLockoutPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetLockoutPolicyToDefaultResponse;
}

export namespace ResetLockoutPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetPrivacyPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPrivacyPolicyRequest): GetPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: GetPrivacyPolicyRequest, reader: jspb.BinaryReader): GetPrivacyPolicyRequest;
}

export namespace GetPrivacyPolicyRequest {
  export type AsObject = {
  }
}

export class GetPrivacyPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PrivacyPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PrivacyPolicy): GetPrivacyPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPrivacyPolicyResponse): GetPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: GetPrivacyPolicyResponse, reader: jspb.BinaryReader): GetPrivacyPolicyResponse;
}

export namespace GetPrivacyPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PrivacyPolicy.AsObject,
  }
}

export class GetDefaultPrivacyPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPrivacyPolicyRequest): GetDefaultPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultPrivacyPolicyRequest, reader: jspb.BinaryReader): GetDefaultPrivacyPolicyRequest;
}

export namespace GetDefaultPrivacyPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultPrivacyPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.PrivacyPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.PrivacyPolicy): GetDefaultPrivacyPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPrivacyPolicyResponse): GetDefaultPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultPrivacyPolicyResponse, reader: jspb.BinaryReader): GetDefaultPrivacyPolicyResponse;
}

export namespace GetDefaultPrivacyPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.PrivacyPolicy.AsObject,
  }
}

export class AddCustomPrivacyPolicyRequest extends jspb.Message {
  getTosLink(): string;
  setTosLink(value: string): AddCustomPrivacyPolicyRequest;

  getPrivacyLink(): string;
  setPrivacyLink(value: string): AddCustomPrivacyPolicyRequest;

  getHelpLink(): string;
  setHelpLink(value: string): AddCustomPrivacyPolicyRequest;

  getSupportEmail(): string;
  setSupportEmail(value: string): AddCustomPrivacyPolicyRequest;

  getDocsLink(): string;
  setDocsLink(value: string): AddCustomPrivacyPolicyRequest;

  getCustomLink(): string;
  setCustomLink(value: string): AddCustomPrivacyPolicyRequest;

  getCustomLinkText(): string;
  setCustomLinkText(value: string): AddCustomPrivacyPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPrivacyPolicyRequest): AddCustomPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomPrivacyPolicyRequest, reader: jspb.BinaryReader): AddCustomPrivacyPolicyRequest;
}

export namespace AddCustomPrivacyPolicyRequest {
  export type AsObject = {
    tosLink: string,
    privacyLink: string,
    helpLink: string,
    supportEmail: string,
    docsLink: string,
    customLink: string,
    customLinkText: string,
  }
}

export class AddCustomPrivacyPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomPrivacyPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomPrivacyPolicyResponse): AddCustomPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomPrivacyPolicyResponse, reader: jspb.BinaryReader): AddCustomPrivacyPolicyResponse;
}

export namespace AddCustomPrivacyPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomPrivacyPolicyRequest extends jspb.Message {
  getTosLink(): string;
  setTosLink(value: string): UpdateCustomPrivacyPolicyRequest;

  getPrivacyLink(): string;
  setPrivacyLink(value: string): UpdateCustomPrivacyPolicyRequest;

  getHelpLink(): string;
  setHelpLink(value: string): UpdateCustomPrivacyPolicyRequest;

  getSupportEmail(): string;
  setSupportEmail(value: string): UpdateCustomPrivacyPolicyRequest;

  getDocsLink(): string;
  setDocsLink(value: string): UpdateCustomPrivacyPolicyRequest;

  getCustomLink(): string;
  setCustomLink(value: string): UpdateCustomPrivacyPolicyRequest;

  getCustomLinkText(): string;
  setCustomLinkText(value: string): UpdateCustomPrivacyPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPrivacyPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPrivacyPolicyRequest): UpdateCustomPrivacyPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPrivacyPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPrivacyPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomPrivacyPolicyRequest, reader: jspb.BinaryReader): UpdateCustomPrivacyPolicyRequest;
}

export namespace UpdateCustomPrivacyPolicyRequest {
  export type AsObject = {
    tosLink: string,
    privacyLink: string,
    helpLink: string,
    supportEmail: string,
    docsLink: string,
    customLink: string,
    customLinkText: string,
  }
}

export class UpdateCustomPrivacyPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomPrivacyPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomPrivacyPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomPrivacyPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomPrivacyPolicyResponse): UpdateCustomPrivacyPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomPrivacyPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomPrivacyPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomPrivacyPolicyResponse, reader: jspb.BinaryReader): UpdateCustomPrivacyPolicyResponse;
}

export namespace UpdateCustomPrivacyPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetPrivacyPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPrivacyPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPrivacyPolicyToDefaultRequest): ResetPrivacyPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetPrivacyPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPrivacyPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetPrivacyPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetPrivacyPolicyToDefaultRequest;
}

export namespace ResetPrivacyPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetPrivacyPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetPrivacyPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetPrivacyPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetPrivacyPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetPrivacyPolicyToDefaultResponse): ResetPrivacyPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetPrivacyPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetPrivacyPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetPrivacyPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetPrivacyPolicyToDefaultResponse;
}

export namespace ResetPrivacyPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetNotificationPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationPolicyRequest): GetNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: GetNotificationPolicyRequest, reader: jspb.BinaryReader): GetNotificationPolicyRequest;
}

export namespace GetNotificationPolicyRequest {
  export type AsObject = {
  }
}

export class GetNotificationPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.NotificationPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.NotificationPolicy): GetNotificationPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetNotificationPolicyResponse): GetNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: GetNotificationPolicyResponse, reader: jspb.BinaryReader): GetNotificationPolicyResponse;
}

export namespace GetNotificationPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.NotificationPolicy.AsObject,
  }
}

export class GetDefaultNotificationPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultNotificationPolicyRequest): GetDefaultNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultNotificationPolicyRequest, reader: jspb.BinaryReader): GetDefaultNotificationPolicyRequest;
}

export namespace GetDefaultNotificationPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultNotificationPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.NotificationPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.NotificationPolicy): GetDefaultNotificationPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultNotificationPolicyResponse): GetDefaultNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultNotificationPolicyResponse, reader: jspb.BinaryReader): GetDefaultNotificationPolicyResponse;
}

export namespace GetDefaultNotificationPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.NotificationPolicy.AsObject,
  }
}

export class AddCustomNotificationPolicyRequest extends jspb.Message {
  getPasswordChange(): boolean;
  setPasswordChange(value: boolean): AddCustomNotificationPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomNotificationPolicyRequest): AddCustomNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomNotificationPolicyRequest, reader: jspb.BinaryReader): AddCustomNotificationPolicyRequest;
}

export namespace AddCustomNotificationPolicyRequest {
  export type AsObject = {
    passwordChange: boolean,
  }
}

export class AddCustomNotificationPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomNotificationPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomNotificationPolicyResponse): AddCustomNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomNotificationPolicyResponse, reader: jspb.BinaryReader): AddCustomNotificationPolicyResponse;
}

export namespace AddCustomNotificationPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomNotificationPolicyRequest extends jspb.Message {
  getPasswordChange(): boolean;
  setPasswordChange(value: boolean): UpdateCustomNotificationPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomNotificationPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomNotificationPolicyRequest): UpdateCustomNotificationPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomNotificationPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomNotificationPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomNotificationPolicyRequest, reader: jspb.BinaryReader): UpdateCustomNotificationPolicyRequest;
}

export namespace UpdateCustomNotificationPolicyRequest {
  export type AsObject = {
    passwordChange: boolean,
  }
}

export class UpdateCustomNotificationPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomNotificationPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomNotificationPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomNotificationPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomNotificationPolicyResponse): UpdateCustomNotificationPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomNotificationPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomNotificationPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomNotificationPolicyResponse, reader: jspb.BinaryReader): UpdateCustomNotificationPolicyResponse;
}

export namespace UpdateCustomNotificationPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetNotificationPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetNotificationPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetNotificationPolicyToDefaultRequest): ResetNotificationPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetNotificationPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetNotificationPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetNotificationPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetNotificationPolicyToDefaultRequest;
}

export namespace ResetNotificationPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetNotificationPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetNotificationPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetNotificationPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetNotificationPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetNotificationPolicyToDefaultResponse): ResetNotificationPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetNotificationPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetNotificationPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetNotificationPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetNotificationPolicyToDefaultResponse;
}

export namespace ResetNotificationPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetLabelPolicyRequest): GetLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetLabelPolicyRequest, reader: jspb.BinaryReader): GetLabelPolicyRequest;
}

export namespace GetLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetLabelPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetLabelPolicyResponse): GetLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetLabelPolicyResponse, reader: jspb.BinaryReader): GetLabelPolicyResponse;
}

export namespace GetLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetPreviewLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPreviewLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPreviewLabelPolicyRequest): GetPreviewLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetPreviewLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPreviewLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetPreviewLabelPolicyRequest, reader: jspb.BinaryReader): GetPreviewLabelPolicyRequest;
}

export namespace GetPreviewLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetPreviewLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetPreviewLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetPreviewLabelPolicyResponse;

  getIsDefault(): boolean;
  setIsDefault(value: boolean): GetPreviewLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPreviewLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPreviewLabelPolicyResponse): GetPreviewLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetPreviewLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPreviewLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetPreviewLabelPolicyResponse, reader: jspb.BinaryReader): GetPreviewLabelPolicyResponse;
}

export namespace GetPreviewLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
    isDefault: boolean,
  }
}

export class GetDefaultLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLabelPolicyRequest): GetDefaultLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLabelPolicyRequest;
  static deserializeBinaryFromReader(message: GetDefaultLabelPolicyRequest, reader: jspb.BinaryReader): GetDefaultLabelPolicyRequest;
}

export namespace GetDefaultLabelPolicyRequest {
  export type AsObject = {
  }
}

export class GetDefaultLabelPolicyResponse extends jspb.Message {
  getPolicy(): zitadel_policy_pb.LabelPolicy | undefined;
  setPolicy(value?: zitadel_policy_pb.LabelPolicy): GetDefaultLabelPolicyResponse;
  hasPolicy(): boolean;
  clearPolicy(): GetDefaultLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLabelPolicyResponse): GetDefaultLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLabelPolicyResponse;
  static deserializeBinaryFromReader(message: GetDefaultLabelPolicyResponse, reader: jspb.BinaryReader): GetDefaultLabelPolicyResponse;
}

export namespace GetDefaultLabelPolicyResponse {
  export type AsObject = {
    policy?: zitadel_policy_pb.LabelPolicy.AsObject,
  }
}

export class AddCustomLabelPolicyRequest extends jspb.Message {
  getPrimaryColor(): string;
  setPrimaryColor(value: string): AddCustomLabelPolicyRequest;

  getHideLoginNameSuffix(): boolean;
  setHideLoginNameSuffix(value: boolean): AddCustomLabelPolicyRequest;

  getWarnColor(): string;
  setWarnColor(value: string): AddCustomLabelPolicyRequest;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): AddCustomLabelPolicyRequest;

  getFontColor(): string;
  setFontColor(value: string): AddCustomLabelPolicyRequest;

  getPrimaryColorDark(): string;
  setPrimaryColorDark(value: string): AddCustomLabelPolicyRequest;

  getBackgroundColorDark(): string;
  setBackgroundColorDark(value: string): AddCustomLabelPolicyRequest;

  getWarnColorDark(): string;
  setWarnColorDark(value: string): AddCustomLabelPolicyRequest;

  getFontColorDark(): string;
  setFontColorDark(value: string): AddCustomLabelPolicyRequest;

  getDisableWatermark(): boolean;
  setDisableWatermark(value: boolean): AddCustomLabelPolicyRequest;

  getThemeMode(): zitadel_policy_pb.ThemeMode;
  setThemeMode(value: zitadel_policy_pb.ThemeMode): AddCustomLabelPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLabelPolicyRequest): AddCustomLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLabelPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomLabelPolicyRequest, reader: jspb.BinaryReader): AddCustomLabelPolicyRequest;
}

export namespace AddCustomLabelPolicyRequest {
  export type AsObject = {
    primaryColor: string,
    hideLoginNameSuffix: boolean,
    warnColor: string,
    backgroundColor: string,
    fontColor: string,
    primaryColorDark: string,
    backgroundColorDark: string,
    warnColorDark: string,
    fontColorDark: string,
    disableWatermark: boolean,
    themeMode: zitadel_policy_pb.ThemeMode,
  }
}

export class AddCustomLabelPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddCustomLabelPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): AddCustomLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomLabelPolicyResponse): AddCustomLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: AddCustomLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomLabelPolicyResponse;
  static deserializeBinaryFromReader(message: AddCustomLabelPolicyResponse, reader: jspb.BinaryReader): AddCustomLabelPolicyResponse;
}

export namespace AddCustomLabelPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateCustomLabelPolicyRequest extends jspb.Message {
  getPrimaryColor(): string;
  setPrimaryColor(value: string): UpdateCustomLabelPolicyRequest;

  getHideLoginNameSuffix(): boolean;
  setHideLoginNameSuffix(value: boolean): UpdateCustomLabelPolicyRequest;

  getWarnColor(): string;
  setWarnColor(value: string): UpdateCustomLabelPolicyRequest;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): UpdateCustomLabelPolicyRequest;

  getFontColor(): string;
  setFontColor(value: string): UpdateCustomLabelPolicyRequest;

  getPrimaryColorDark(): string;
  setPrimaryColorDark(value: string): UpdateCustomLabelPolicyRequest;

  getBackgroundColorDark(): string;
  setBackgroundColorDark(value: string): UpdateCustomLabelPolicyRequest;

  getWarnColorDark(): string;
  setWarnColorDark(value: string): UpdateCustomLabelPolicyRequest;

  getFontColorDark(): string;
  setFontColorDark(value: string): UpdateCustomLabelPolicyRequest;

  getDisableWatermark(): boolean;
  setDisableWatermark(value: boolean): UpdateCustomLabelPolicyRequest;

  getThemeMode(): zitadel_policy_pb.ThemeMode;
  setThemeMode(value: zitadel_policy_pb.ThemeMode): UpdateCustomLabelPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLabelPolicyRequest): UpdateCustomLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLabelPolicyRequest;
  static deserializeBinaryFromReader(message: UpdateCustomLabelPolicyRequest, reader: jspb.BinaryReader): UpdateCustomLabelPolicyRequest;
}

export namespace UpdateCustomLabelPolicyRequest {
  export type AsObject = {
    primaryColor: string,
    hideLoginNameSuffix: boolean,
    warnColor: string,
    backgroundColor: string,
    fontColor: string,
    primaryColorDark: string,
    backgroundColorDark: string,
    warnColorDark: string,
    fontColorDark: string,
    disableWatermark: boolean,
    themeMode: zitadel_policy_pb.ThemeMode,
  }
}

export class UpdateCustomLabelPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateCustomLabelPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateCustomLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateCustomLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateCustomLabelPolicyResponse): UpdateCustomLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateCustomLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateCustomLabelPolicyResponse;
  static deserializeBinaryFromReader(message: UpdateCustomLabelPolicyResponse, reader: jspb.BinaryReader): UpdateCustomLabelPolicyResponse;
}

export namespace UpdateCustomLabelPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ActivateCustomLabelPolicyRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateCustomLabelPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateCustomLabelPolicyRequest): ActivateCustomLabelPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: ActivateCustomLabelPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateCustomLabelPolicyRequest;
  static deserializeBinaryFromReader(message: ActivateCustomLabelPolicyRequest, reader: jspb.BinaryReader): ActivateCustomLabelPolicyRequest;
}

export namespace ActivateCustomLabelPolicyRequest {
  export type AsObject = {
  }
}

export class ActivateCustomLabelPolicyResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ActivateCustomLabelPolicyResponse;
  hasDetails(): boolean;
  clearDetails(): ActivateCustomLabelPolicyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActivateCustomLabelPolicyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ActivateCustomLabelPolicyResponse): ActivateCustomLabelPolicyResponse.AsObject;
  static serializeBinaryToWriter(message: ActivateCustomLabelPolicyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActivateCustomLabelPolicyResponse;
  static deserializeBinaryFromReader(message: ActivateCustomLabelPolicyResponse, reader: jspb.BinaryReader): ActivateCustomLabelPolicyResponse;
}

export namespace ActivateCustomLabelPolicyResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveCustomLabelPolicyLogoRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyLogoRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyLogoRequest): RemoveCustomLabelPolicyLogoRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyLogoRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyLogoRequest;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyLogoRequest, reader: jspb.BinaryReader): RemoveCustomLabelPolicyLogoRequest;
}

export namespace RemoveCustomLabelPolicyLogoRequest {
  export type AsObject = {
  }
}

export class RemoveCustomLabelPolicyLogoResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveCustomLabelPolicyLogoResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveCustomLabelPolicyLogoResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyLogoResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyLogoResponse): RemoveCustomLabelPolicyLogoResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyLogoResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyLogoResponse;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyLogoResponse, reader: jspb.BinaryReader): RemoveCustomLabelPolicyLogoResponse;
}

export namespace RemoveCustomLabelPolicyLogoResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveCustomLabelPolicyLogoDarkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyLogoDarkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyLogoDarkRequest): RemoveCustomLabelPolicyLogoDarkRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyLogoDarkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyLogoDarkRequest;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyLogoDarkRequest, reader: jspb.BinaryReader): RemoveCustomLabelPolicyLogoDarkRequest;
}

export namespace RemoveCustomLabelPolicyLogoDarkRequest {
  export type AsObject = {
  }
}

export class RemoveCustomLabelPolicyLogoDarkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveCustomLabelPolicyLogoDarkResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveCustomLabelPolicyLogoDarkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyLogoDarkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyLogoDarkResponse): RemoveCustomLabelPolicyLogoDarkResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyLogoDarkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyLogoDarkResponse;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyLogoDarkResponse, reader: jspb.BinaryReader): RemoveCustomLabelPolicyLogoDarkResponse;
}

export namespace RemoveCustomLabelPolicyLogoDarkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveCustomLabelPolicyIconRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyIconRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyIconRequest): RemoveCustomLabelPolicyIconRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyIconRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyIconRequest;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyIconRequest, reader: jspb.BinaryReader): RemoveCustomLabelPolicyIconRequest;
}

export namespace RemoveCustomLabelPolicyIconRequest {
  export type AsObject = {
  }
}

export class RemoveCustomLabelPolicyIconResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveCustomLabelPolicyIconResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveCustomLabelPolicyIconResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyIconResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyIconResponse): RemoveCustomLabelPolicyIconResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyIconResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyIconResponse;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyIconResponse, reader: jspb.BinaryReader): RemoveCustomLabelPolicyIconResponse;
}

export namespace RemoveCustomLabelPolicyIconResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveCustomLabelPolicyIconDarkRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyIconDarkRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyIconDarkRequest): RemoveCustomLabelPolicyIconDarkRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyIconDarkRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyIconDarkRequest;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyIconDarkRequest, reader: jspb.BinaryReader): RemoveCustomLabelPolicyIconDarkRequest;
}

export namespace RemoveCustomLabelPolicyIconDarkRequest {
  export type AsObject = {
  }
}

export class RemoveCustomLabelPolicyIconDarkResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveCustomLabelPolicyIconDarkResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveCustomLabelPolicyIconDarkResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyIconDarkResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyIconDarkResponse): RemoveCustomLabelPolicyIconDarkResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyIconDarkResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyIconDarkResponse;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyIconDarkResponse, reader: jspb.BinaryReader): RemoveCustomLabelPolicyIconDarkResponse;
}

export namespace RemoveCustomLabelPolicyIconDarkResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveCustomLabelPolicyFontRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyFontRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyFontRequest): RemoveCustomLabelPolicyFontRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyFontRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyFontRequest;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyFontRequest, reader: jspb.BinaryReader): RemoveCustomLabelPolicyFontRequest;
}

export namespace RemoveCustomLabelPolicyFontRequest {
  export type AsObject = {
  }
}

export class RemoveCustomLabelPolicyFontResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveCustomLabelPolicyFontResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveCustomLabelPolicyFontResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveCustomLabelPolicyFontResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveCustomLabelPolicyFontResponse): RemoveCustomLabelPolicyFontResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveCustomLabelPolicyFontResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveCustomLabelPolicyFontResponse;
  static deserializeBinaryFromReader(message: RemoveCustomLabelPolicyFontResponse, reader: jspb.BinaryReader): RemoveCustomLabelPolicyFontResponse;
}

export namespace RemoveCustomLabelPolicyFontResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetLabelPolicyToDefaultRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLabelPolicyToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLabelPolicyToDefaultRequest): ResetLabelPolicyToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetLabelPolicyToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLabelPolicyToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetLabelPolicyToDefaultRequest, reader: jspb.BinaryReader): ResetLabelPolicyToDefaultRequest;
}

export namespace ResetLabelPolicyToDefaultRequest {
  export type AsObject = {
  }
}

export class ResetLabelPolicyToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetLabelPolicyToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetLabelPolicyToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLabelPolicyToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLabelPolicyToDefaultResponse): ResetLabelPolicyToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetLabelPolicyToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLabelPolicyToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetLabelPolicyToDefaultResponse, reader: jspb.BinaryReader): ResetLabelPolicyToDefaultResponse;
}

export namespace ResetLabelPolicyToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomInitMessageTextRequest): GetCustomInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomInitMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomInitMessageTextRequest, reader: jspb.BinaryReader): GetCustomInitMessageTextRequest;
}

export namespace GetCustomInitMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomInitMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomInitMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomInitMessageTextResponse): GetCustomInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomInitMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomInitMessageTextResponse, reader: jspb.BinaryReader): GetCustomInitMessageTextResponse;
}

export namespace GetCustomInitMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultInitMessageTextRequest): GetDefaultInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultInitMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultInitMessageTextRequest, reader: jspb.BinaryReader): GetDefaultInitMessageTextRequest;
}

export namespace GetDefaultInitMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultInitMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultInitMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultInitMessageTextResponse): GetDefaultInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultInitMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultInitMessageTextResponse, reader: jspb.BinaryReader): GetDefaultInitMessageTextResponse;
}

export namespace GetDefaultInitMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomInitMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomInitMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomInitMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomInitMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomInitMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomInitMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomInitMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomInitMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomInitMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomInitMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomInitMessageTextRequest): SetCustomInitMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomInitMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomInitMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomInitMessageTextRequest, reader: jspb.BinaryReader): SetCustomInitMessageTextRequest;
}

export namespace SetCustomInitMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomInitMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomInitMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomInitMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomInitMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomInitMessageTextResponse): SetCustomInitMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomInitMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomInitMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomInitMessageTextResponse, reader: jspb.BinaryReader): SetCustomInitMessageTextResponse;
}

export namespace SetCustomInitMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomInitMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomInitMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomInitMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomInitMessageTextToDefaultRequest): ResetCustomInitMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomInitMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomInitMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomInitMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomInitMessageTextToDefaultRequest;
}

export namespace ResetCustomInitMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomInitMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomInitMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomInitMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomInitMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomInitMessageTextToDefaultResponse): ResetCustomInitMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomInitMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomInitMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomInitMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomInitMessageTextToDefaultResponse;
}

export namespace ResetCustomInitMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetDefaultLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginTextsRequest): GetDefaultLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginTextsRequest;
  static deserializeBinaryFromReader(message: GetDefaultLoginTextsRequest, reader: jspb.BinaryReader): GetDefaultLoginTextsRequest;
}

export namespace GetDefaultLoginTextsRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultLoginTextsResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.LoginCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.LoginCustomText): GetDefaultLoginTextsResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultLoginTextsResponse): GetDefaultLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultLoginTextsResponse;
  static deserializeBinaryFromReader(message: GetDefaultLoginTextsResponse, reader: jspb.BinaryReader): GetDefaultLoginTextsResponse;
}

export namespace GetDefaultLoginTextsResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.LoginCustomText.AsObject,
  }
}

export class GetCustomLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomLoginTextsRequest): GetCustomLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomLoginTextsRequest;
  static deserializeBinaryFromReader(message: GetCustomLoginTextsRequest, reader: jspb.BinaryReader): GetCustomLoginTextsRequest;
}

export namespace GetCustomLoginTextsRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomLoginTextsResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.LoginCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.LoginCustomText): GetCustomLoginTextsResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomLoginTextsResponse): GetCustomLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomLoginTextsResponse;
  static deserializeBinaryFromReader(message: GetCustomLoginTextsResponse, reader: jspb.BinaryReader): GetCustomLoginTextsResponse;
}

export namespace GetCustomLoginTextsResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.LoginCustomText.AsObject,
  }
}

export class SetCustomLoginTextsRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomLoginTextsRequest;

  getSelectAccountText(): zitadel_text_pb.SelectAccountScreenText | undefined;
  setSelectAccountText(value?: zitadel_text_pb.SelectAccountScreenText): SetCustomLoginTextsRequest;
  hasSelectAccountText(): boolean;
  clearSelectAccountText(): SetCustomLoginTextsRequest;

  getLoginText(): zitadel_text_pb.LoginScreenText | undefined;
  setLoginText(value?: zitadel_text_pb.LoginScreenText): SetCustomLoginTextsRequest;
  hasLoginText(): boolean;
  clearLoginText(): SetCustomLoginTextsRequest;

  getPasswordText(): zitadel_text_pb.PasswordScreenText | undefined;
  setPasswordText(value?: zitadel_text_pb.PasswordScreenText): SetCustomLoginTextsRequest;
  hasPasswordText(): boolean;
  clearPasswordText(): SetCustomLoginTextsRequest;

  getUsernameChangeText(): zitadel_text_pb.UsernameChangeScreenText | undefined;
  setUsernameChangeText(value?: zitadel_text_pb.UsernameChangeScreenText): SetCustomLoginTextsRequest;
  hasUsernameChangeText(): boolean;
  clearUsernameChangeText(): SetCustomLoginTextsRequest;

  getUsernameChangeDoneText(): zitadel_text_pb.UsernameChangeDoneScreenText | undefined;
  setUsernameChangeDoneText(value?: zitadel_text_pb.UsernameChangeDoneScreenText): SetCustomLoginTextsRequest;
  hasUsernameChangeDoneText(): boolean;
  clearUsernameChangeDoneText(): SetCustomLoginTextsRequest;

  getInitPasswordText(): zitadel_text_pb.InitPasswordScreenText | undefined;
  setInitPasswordText(value?: zitadel_text_pb.InitPasswordScreenText): SetCustomLoginTextsRequest;
  hasInitPasswordText(): boolean;
  clearInitPasswordText(): SetCustomLoginTextsRequest;

  getInitPasswordDoneText(): zitadel_text_pb.InitPasswordDoneScreenText | undefined;
  setInitPasswordDoneText(value?: zitadel_text_pb.InitPasswordDoneScreenText): SetCustomLoginTextsRequest;
  hasInitPasswordDoneText(): boolean;
  clearInitPasswordDoneText(): SetCustomLoginTextsRequest;

  getEmailVerificationText(): zitadel_text_pb.EmailVerificationScreenText | undefined;
  setEmailVerificationText(value?: zitadel_text_pb.EmailVerificationScreenText): SetCustomLoginTextsRequest;
  hasEmailVerificationText(): boolean;
  clearEmailVerificationText(): SetCustomLoginTextsRequest;

  getEmailVerificationDoneText(): zitadel_text_pb.EmailVerificationDoneScreenText | undefined;
  setEmailVerificationDoneText(value?: zitadel_text_pb.EmailVerificationDoneScreenText): SetCustomLoginTextsRequest;
  hasEmailVerificationDoneText(): boolean;
  clearEmailVerificationDoneText(): SetCustomLoginTextsRequest;

  getInitializeUserText(): zitadel_text_pb.InitializeUserScreenText | undefined;
  setInitializeUserText(value?: zitadel_text_pb.InitializeUserScreenText): SetCustomLoginTextsRequest;
  hasInitializeUserText(): boolean;
  clearInitializeUserText(): SetCustomLoginTextsRequest;

  getInitializeDoneText(): zitadel_text_pb.InitializeUserDoneScreenText | undefined;
  setInitializeDoneText(value?: zitadel_text_pb.InitializeUserDoneScreenText): SetCustomLoginTextsRequest;
  hasInitializeDoneText(): boolean;
  clearInitializeDoneText(): SetCustomLoginTextsRequest;

  getInitMfaPromptText(): zitadel_text_pb.InitMFAPromptScreenText | undefined;
  setInitMfaPromptText(value?: zitadel_text_pb.InitMFAPromptScreenText): SetCustomLoginTextsRequest;
  hasInitMfaPromptText(): boolean;
  clearInitMfaPromptText(): SetCustomLoginTextsRequest;

  getInitMfaOtpText(): zitadel_text_pb.InitMFAOTPScreenText | undefined;
  setInitMfaOtpText(value?: zitadel_text_pb.InitMFAOTPScreenText): SetCustomLoginTextsRequest;
  hasInitMfaOtpText(): boolean;
  clearInitMfaOtpText(): SetCustomLoginTextsRequest;

  getInitMfaU2fText(): zitadel_text_pb.InitMFAU2FScreenText | undefined;
  setInitMfaU2fText(value?: zitadel_text_pb.InitMFAU2FScreenText): SetCustomLoginTextsRequest;
  hasInitMfaU2fText(): boolean;
  clearInitMfaU2fText(): SetCustomLoginTextsRequest;

  getInitMfaDoneText(): zitadel_text_pb.InitMFADoneScreenText | undefined;
  setInitMfaDoneText(value?: zitadel_text_pb.InitMFADoneScreenText): SetCustomLoginTextsRequest;
  hasInitMfaDoneText(): boolean;
  clearInitMfaDoneText(): SetCustomLoginTextsRequest;

  getMfaProvidersText(): zitadel_text_pb.MFAProvidersText | undefined;
  setMfaProvidersText(value?: zitadel_text_pb.MFAProvidersText): SetCustomLoginTextsRequest;
  hasMfaProvidersText(): boolean;
  clearMfaProvidersText(): SetCustomLoginTextsRequest;

  getVerifyMfaOtpText(): zitadel_text_pb.VerifyMFAOTPScreenText | undefined;
  setVerifyMfaOtpText(value?: zitadel_text_pb.VerifyMFAOTPScreenText): SetCustomLoginTextsRequest;
  hasVerifyMfaOtpText(): boolean;
  clearVerifyMfaOtpText(): SetCustomLoginTextsRequest;

  getVerifyMfaU2fText(): zitadel_text_pb.VerifyMFAU2FScreenText | undefined;
  setVerifyMfaU2fText(value?: zitadel_text_pb.VerifyMFAU2FScreenText): SetCustomLoginTextsRequest;
  hasVerifyMfaU2fText(): boolean;
  clearVerifyMfaU2fText(): SetCustomLoginTextsRequest;

  getPasswordlessText(): zitadel_text_pb.PasswordlessScreenText | undefined;
  setPasswordlessText(value?: zitadel_text_pb.PasswordlessScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessText(): boolean;
  clearPasswordlessText(): SetCustomLoginTextsRequest;

  getPasswordChangeText(): zitadel_text_pb.PasswordChangeScreenText | undefined;
  setPasswordChangeText(value?: zitadel_text_pb.PasswordChangeScreenText): SetCustomLoginTextsRequest;
  hasPasswordChangeText(): boolean;
  clearPasswordChangeText(): SetCustomLoginTextsRequest;

  getPasswordChangeDoneText(): zitadel_text_pb.PasswordChangeDoneScreenText | undefined;
  setPasswordChangeDoneText(value?: zitadel_text_pb.PasswordChangeDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordChangeDoneText(): boolean;
  clearPasswordChangeDoneText(): SetCustomLoginTextsRequest;

  getPasswordResetDoneText(): zitadel_text_pb.PasswordResetDoneScreenText | undefined;
  setPasswordResetDoneText(value?: zitadel_text_pb.PasswordResetDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordResetDoneText(): boolean;
  clearPasswordResetDoneText(): SetCustomLoginTextsRequest;

  getRegistrationOptionText(): zitadel_text_pb.RegistrationOptionScreenText | undefined;
  setRegistrationOptionText(value?: zitadel_text_pb.RegistrationOptionScreenText): SetCustomLoginTextsRequest;
  hasRegistrationOptionText(): boolean;
  clearRegistrationOptionText(): SetCustomLoginTextsRequest;

  getRegistrationUserText(): zitadel_text_pb.RegistrationUserScreenText | undefined;
  setRegistrationUserText(value?: zitadel_text_pb.RegistrationUserScreenText): SetCustomLoginTextsRequest;
  hasRegistrationUserText(): boolean;
  clearRegistrationUserText(): SetCustomLoginTextsRequest;

  getRegistrationOrgText(): zitadel_text_pb.RegistrationOrgScreenText | undefined;
  setRegistrationOrgText(value?: zitadel_text_pb.RegistrationOrgScreenText): SetCustomLoginTextsRequest;
  hasRegistrationOrgText(): boolean;
  clearRegistrationOrgText(): SetCustomLoginTextsRequest;

  getLinkingUserDoneText(): zitadel_text_pb.LinkingUserDoneScreenText | undefined;
  setLinkingUserDoneText(value?: zitadel_text_pb.LinkingUserDoneScreenText): SetCustomLoginTextsRequest;
  hasLinkingUserDoneText(): boolean;
  clearLinkingUserDoneText(): SetCustomLoginTextsRequest;

  getExternalUserNotFoundText(): zitadel_text_pb.ExternalUserNotFoundScreenText | undefined;
  setExternalUserNotFoundText(value?: zitadel_text_pb.ExternalUserNotFoundScreenText): SetCustomLoginTextsRequest;
  hasExternalUserNotFoundText(): boolean;
  clearExternalUserNotFoundText(): SetCustomLoginTextsRequest;

  getSuccessLoginText(): zitadel_text_pb.SuccessLoginScreenText | undefined;
  setSuccessLoginText(value?: zitadel_text_pb.SuccessLoginScreenText): SetCustomLoginTextsRequest;
  hasSuccessLoginText(): boolean;
  clearSuccessLoginText(): SetCustomLoginTextsRequest;

  getLogoutText(): zitadel_text_pb.LogoutDoneScreenText | undefined;
  setLogoutText(value?: zitadel_text_pb.LogoutDoneScreenText): SetCustomLoginTextsRequest;
  hasLogoutText(): boolean;
  clearLogoutText(): SetCustomLoginTextsRequest;

  getFooterText(): zitadel_text_pb.FooterText | undefined;
  setFooterText(value?: zitadel_text_pb.FooterText): SetCustomLoginTextsRequest;
  hasFooterText(): boolean;
  clearFooterText(): SetCustomLoginTextsRequest;

  getPasswordlessPromptText(): zitadel_text_pb.PasswordlessPromptScreenText | undefined;
  setPasswordlessPromptText(value?: zitadel_text_pb.PasswordlessPromptScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessPromptText(): boolean;
  clearPasswordlessPromptText(): SetCustomLoginTextsRequest;

  getPasswordlessRegistrationText(): zitadel_text_pb.PasswordlessRegistrationScreenText | undefined;
  setPasswordlessRegistrationText(value?: zitadel_text_pb.PasswordlessRegistrationScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessRegistrationText(): boolean;
  clearPasswordlessRegistrationText(): SetCustomLoginTextsRequest;

  getPasswordlessRegistrationDoneText(): zitadel_text_pb.PasswordlessRegistrationDoneScreenText | undefined;
  setPasswordlessRegistrationDoneText(value?: zitadel_text_pb.PasswordlessRegistrationDoneScreenText): SetCustomLoginTextsRequest;
  hasPasswordlessRegistrationDoneText(): boolean;
  clearPasswordlessRegistrationDoneText(): SetCustomLoginTextsRequest;

  getExternalRegistrationUserOverviewText(): zitadel_text_pb.ExternalRegistrationUserOverviewScreenText | undefined;
  setExternalRegistrationUserOverviewText(value?: zitadel_text_pb.ExternalRegistrationUserOverviewScreenText): SetCustomLoginTextsRequest;
  hasExternalRegistrationUserOverviewText(): boolean;
  clearExternalRegistrationUserOverviewText(): SetCustomLoginTextsRequest;

  getLinkingUserPromptText(): zitadel_text_pb.LinkingUserPromptScreenText | undefined;
  setLinkingUserPromptText(value?: zitadel_text_pb.LinkingUserPromptScreenText): SetCustomLoginTextsRequest;
  hasLinkingUserPromptText(): boolean;
  clearLinkingUserPromptText(): SetCustomLoginTextsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomLoginTextsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomLoginTextsRequest): SetCustomLoginTextsRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomLoginTextsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomLoginTextsRequest;
  static deserializeBinaryFromReader(message: SetCustomLoginTextsRequest, reader: jspb.BinaryReader): SetCustomLoginTextsRequest;
}

export namespace SetCustomLoginTextsRequest {
  export type AsObject = {
    language: string,
    selectAccountText?: zitadel_text_pb.SelectAccountScreenText.AsObject,
    loginText?: zitadel_text_pb.LoginScreenText.AsObject,
    passwordText?: zitadel_text_pb.PasswordScreenText.AsObject,
    usernameChangeText?: zitadel_text_pb.UsernameChangeScreenText.AsObject,
    usernameChangeDoneText?: zitadel_text_pb.UsernameChangeDoneScreenText.AsObject,
    initPasswordText?: zitadel_text_pb.InitPasswordScreenText.AsObject,
    initPasswordDoneText?: zitadel_text_pb.InitPasswordDoneScreenText.AsObject,
    emailVerificationText?: zitadel_text_pb.EmailVerificationScreenText.AsObject,
    emailVerificationDoneText?: zitadel_text_pb.EmailVerificationDoneScreenText.AsObject,
    initializeUserText?: zitadel_text_pb.InitializeUserScreenText.AsObject,
    initializeDoneText?: zitadel_text_pb.InitializeUserDoneScreenText.AsObject,
    initMfaPromptText?: zitadel_text_pb.InitMFAPromptScreenText.AsObject,
    initMfaOtpText?: zitadel_text_pb.InitMFAOTPScreenText.AsObject,
    initMfaU2fText?: zitadel_text_pb.InitMFAU2FScreenText.AsObject,
    initMfaDoneText?: zitadel_text_pb.InitMFADoneScreenText.AsObject,
    mfaProvidersText?: zitadel_text_pb.MFAProvidersText.AsObject,
    verifyMfaOtpText?: zitadel_text_pb.VerifyMFAOTPScreenText.AsObject,
    verifyMfaU2fText?: zitadel_text_pb.VerifyMFAU2FScreenText.AsObject,
    passwordlessText?: zitadel_text_pb.PasswordlessScreenText.AsObject,
    passwordChangeText?: zitadel_text_pb.PasswordChangeScreenText.AsObject,
    passwordChangeDoneText?: zitadel_text_pb.PasswordChangeDoneScreenText.AsObject,
    passwordResetDoneText?: zitadel_text_pb.PasswordResetDoneScreenText.AsObject,
    registrationOptionText?: zitadel_text_pb.RegistrationOptionScreenText.AsObject,
    registrationUserText?: zitadel_text_pb.RegistrationUserScreenText.AsObject,
    registrationOrgText?: zitadel_text_pb.RegistrationOrgScreenText.AsObject,
    linkingUserDoneText?: zitadel_text_pb.LinkingUserDoneScreenText.AsObject,
    externalUserNotFoundText?: zitadel_text_pb.ExternalUserNotFoundScreenText.AsObject,
    successLoginText?: zitadel_text_pb.SuccessLoginScreenText.AsObject,
    logoutText?: zitadel_text_pb.LogoutDoneScreenText.AsObject,
    footerText?: zitadel_text_pb.FooterText.AsObject,
    passwordlessPromptText?: zitadel_text_pb.PasswordlessPromptScreenText.AsObject,
    passwordlessRegistrationText?: zitadel_text_pb.PasswordlessRegistrationScreenText.AsObject,
    passwordlessRegistrationDoneText?: zitadel_text_pb.PasswordlessRegistrationDoneScreenText.AsObject,
    externalRegistrationUserOverviewText?: zitadel_text_pb.ExternalRegistrationUserOverviewScreenText.AsObject,
    linkingUserPromptText?: zitadel_text_pb.LinkingUserPromptScreenText.AsObject,
  }
}

export class SetCustomLoginTextsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomLoginTextsResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomLoginTextsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomLoginTextsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomLoginTextsResponse): SetCustomLoginTextsResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomLoginTextsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomLoginTextsResponse;
  static deserializeBinaryFromReader(message: SetCustomLoginTextsResponse, reader: jspb.BinaryReader): SetCustomLoginTextsResponse;
}

export namespace SetCustomLoginTextsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomLoginTextsToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomLoginTextsToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomLoginTextsToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomLoginTextsToDefaultRequest): ResetCustomLoginTextsToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomLoginTextsToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomLoginTextsToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomLoginTextsToDefaultRequest, reader: jspb.BinaryReader): ResetCustomLoginTextsToDefaultRequest;
}

export namespace ResetCustomLoginTextsToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomLoginTextsToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomLoginTextsToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomLoginTextsToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomLoginTextsToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomLoginTextsToDefaultResponse): ResetCustomLoginTextsToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomLoginTextsToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomLoginTextsToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomLoginTextsToDefaultResponse, reader: jspb.BinaryReader): ResetCustomLoginTextsToDefaultResponse;
}

export namespace ResetCustomLoginTextsToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordResetMessageTextRequest): GetCustomPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordResetMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordResetMessageTextRequest;
}

export namespace GetCustomPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordResetMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordResetMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordResetMessageTextResponse): GetCustomPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordResetMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordResetMessageTextResponse;
}

export namespace GetCustomPasswordResetMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordResetMessageTextRequest): GetDefaultPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordResetMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordResetMessageTextRequest;
}

export namespace GetDefaultPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordResetMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordResetMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordResetMessageTextResponse): GetDefaultPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordResetMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordResetMessageTextResponse;
}

export namespace GetDefaultPasswordResetMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomPasswordResetMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomPasswordResetMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomPasswordResetMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomPasswordResetMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomPasswordResetMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomPasswordResetMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomPasswordResetMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomPasswordResetMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomPasswordResetMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordResetMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordResetMessageTextRequest): SetCustomPasswordResetMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordResetMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordResetMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomPasswordResetMessageTextRequest, reader: jspb.BinaryReader): SetCustomPasswordResetMessageTextRequest;
}

export namespace SetCustomPasswordResetMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomPasswordResetMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomPasswordResetMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomPasswordResetMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordResetMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordResetMessageTextResponse): SetCustomPasswordResetMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordResetMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordResetMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomPasswordResetMessageTextResponse, reader: jspb.BinaryReader): SetCustomPasswordResetMessageTextResponse;
}

export namespace SetCustomPasswordResetMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordResetMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordResetMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordResetMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordResetMessageTextToDefaultRequest): ResetCustomPasswordResetMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordResetMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordResetMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordResetMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordResetMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordResetMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordResetMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordResetMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordResetMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordResetMessageTextToDefaultResponse): ResetCustomPasswordResetMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordResetMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordResetMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordResetMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordResetMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordResetMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailMessageTextRequest): GetCustomVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyEmailMessageTextRequest;
}

export namespace GetCustomVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyEmailMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyEmailMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailMessageTextResponse): GetCustomVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyEmailMessageTextResponse;
}

export namespace GetCustomVerifyEmailMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailMessageTextRequest): GetDefaultVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyEmailMessageTextRequest;
}

export namespace GetDefaultVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyEmailMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyEmailMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailMessageTextResponse): GetDefaultVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyEmailMessageTextResponse;
}

export namespace GetDefaultVerifyEmailMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomVerifyEmailMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomVerifyEmailMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomVerifyEmailMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomVerifyEmailMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomVerifyEmailMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomVerifyEmailMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomVerifyEmailMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomVerifyEmailMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomVerifyEmailMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyEmailMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyEmailMessageTextRequest): SetCustomVerifyEmailMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyEmailMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyEmailMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomVerifyEmailMessageTextRequest, reader: jspb.BinaryReader): SetCustomVerifyEmailMessageTextRequest;
}

export namespace SetCustomVerifyEmailMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomVerifyEmailMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomVerifyEmailMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomVerifyEmailMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyEmailMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyEmailMessageTextResponse): SetCustomVerifyEmailMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyEmailMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyEmailMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomVerifyEmailMessageTextResponse, reader: jspb.BinaryReader): SetCustomVerifyEmailMessageTextResponse;
}

export namespace SetCustomVerifyEmailMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyEmailMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyEmailMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailMessageTextToDefaultRequest): ResetCustomVerifyEmailMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyEmailMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyEmailMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyEmailMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyEmailMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyEmailMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailMessageTextToDefaultResponse): ResetCustomVerifyEmailMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyEmailMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyEmailMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyPhoneMessageTextRequest): GetCustomVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyPhoneMessageTextRequest;
}

export namespace GetCustomVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyPhoneMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyPhoneMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyPhoneMessageTextResponse): GetCustomVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyPhoneMessageTextResponse;
}

export namespace GetCustomVerifyPhoneMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyPhoneMessageTextRequest): GetDefaultVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyPhoneMessageTextRequest;
}

export namespace GetDefaultVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyPhoneMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyPhoneMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyPhoneMessageTextResponse): GetDefaultVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyPhoneMessageTextResponse;
}

export namespace GetDefaultVerifyPhoneMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomVerifyPhoneMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomVerifyPhoneMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomVerifyPhoneMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyPhoneMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyPhoneMessageTextRequest): SetCustomVerifyPhoneMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyPhoneMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyPhoneMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomVerifyPhoneMessageTextRequest, reader: jspb.BinaryReader): SetCustomVerifyPhoneMessageTextRequest;
}

export namespace SetCustomVerifyPhoneMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomVerifyPhoneMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomVerifyPhoneMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomVerifyPhoneMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyPhoneMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyPhoneMessageTextResponse): SetCustomVerifyPhoneMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyPhoneMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyPhoneMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomVerifyPhoneMessageTextResponse, reader: jspb.BinaryReader): SetCustomVerifyPhoneMessageTextResponse;
}

export namespace SetCustomVerifyPhoneMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyPhoneMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyPhoneMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyPhoneMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyPhoneMessageTextToDefaultRequest): ResetCustomVerifyPhoneMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyPhoneMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyPhoneMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyPhoneMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyPhoneMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyPhoneMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyPhoneMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyPhoneMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyPhoneMessageTextToDefaultResponse): ResetCustomVerifyPhoneMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyPhoneMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyPhoneMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyPhoneMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyPhoneMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifySMSOTPMessageTextRequest): GetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifySMSOTPMessageTextRequest;
}

export namespace GetCustomVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifySMSOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifySMSOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifySMSOTPMessageTextResponse): GetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifySMSOTPMessageTextResponse;
}

export namespace GetCustomVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifySMSOTPMessageTextRequest): GetDefaultVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifySMSOTPMessageTextRequest;
}

export namespace GetDefaultVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifySMSOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifySMSOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifySMSOTPMessageTextResponse): GetDefaultVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifySMSOTPMessageTextResponse;
}

export namespace GetDefaultVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomVerifySMSOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomVerifySMSOTPMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomVerifySMSOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifySMSOTPMessageTextRequest): SetCustomVerifySMSOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifySMSOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifySMSOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomVerifySMSOTPMessageTextRequest, reader: jspb.BinaryReader): SetCustomVerifySMSOTPMessageTextRequest;
}

export namespace SetCustomVerifySMSOTPMessageTextRequest {
  export type AsObject = {
    language: string,
    text: string,
  }
}

export class SetCustomVerifySMSOTPMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomVerifySMSOTPMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomVerifySMSOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifySMSOTPMessageTextResponse): SetCustomVerifySMSOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifySMSOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifySMSOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomVerifySMSOTPMessageTextResponse, reader: jspb.BinaryReader): SetCustomVerifySMSOTPMessageTextResponse;
}

export namespace SetCustomVerifySMSOTPMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifySMSOTPMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifySMSOTPMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifySMSOTPMessageTextToDefaultRequest): ResetCustomVerifySMSOTPMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifySMSOTPMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifySMSOTPMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifySMSOTPMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifySMSOTPMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifySMSOTPMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifySMSOTPMessageTextToDefaultResponse): ResetCustomVerifySMSOTPMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifySMSOTPMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifySMSOTPMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifySMSOTPMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifySMSOTPMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailOTPMessageTextRequest): GetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): GetCustomVerifyEmailOTPMessageTextRequest;
}

export namespace GetCustomVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomVerifyEmailOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomVerifyEmailOTPMessageTextResponse): GetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): GetCustomVerifyEmailOTPMessageTextResponse;
}

export namespace GetCustomVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailOTPMessageTextRequest): GetDefaultVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): GetDefaultVerifyEmailOTPMessageTextRequest;
}

export namespace GetDefaultVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultVerifyEmailOTPMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultVerifyEmailOTPMessageTextResponse): GetDefaultVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): GetDefaultVerifyEmailOTPMessageTextResponse;
}

export namespace GetDefaultVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomVerifyEmailOTPMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomVerifyEmailOTPMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyEmailOTPMessageTextRequest): SetCustomVerifyEmailOTPMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyEmailOTPMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyEmailOTPMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomVerifyEmailOTPMessageTextRequest, reader: jspb.BinaryReader): SetCustomVerifyEmailOTPMessageTextRequest;
}

export namespace SetCustomVerifyEmailOTPMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomVerifyEmailOTPMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomVerifyEmailOTPMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomVerifyEmailOTPMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomVerifyEmailOTPMessageTextResponse): SetCustomVerifyEmailOTPMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomVerifyEmailOTPMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomVerifyEmailOTPMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomVerifyEmailOTPMessageTextResponse, reader: jspb.BinaryReader): SetCustomVerifyEmailOTPMessageTextResponse;
}

export namespace SetCustomVerifyEmailOTPMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomVerifyEmailOTPMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailOTPMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomVerifyEmailOTPMessageTextToDefaultRequest;
}

export namespace ResetCustomVerifyEmailOTPMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomVerifyEmailOTPMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomVerifyEmailOTPMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomVerifyEmailOTPMessageTextToDefaultResponse;
}

export namespace ResetCustomVerifyEmailOTPMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainClaimedMessageTextRequest): GetCustomDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): GetCustomDomainClaimedMessageTextRequest;
}

export namespace GetCustomDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomDomainClaimedMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomDomainClaimedMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomDomainClaimedMessageTextResponse): GetCustomDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): GetCustomDomainClaimedMessageTextResponse;
}

export namespace GetCustomDomainClaimedMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultDomainClaimedMessageTextRequest): GetDefaultDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): GetDefaultDomainClaimedMessageTextRequest;
}

export namespace GetDefaultDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultDomainClaimedMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultDomainClaimedMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultDomainClaimedMessageTextResponse): GetDefaultDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): GetDefaultDomainClaimedMessageTextResponse;
}

export namespace GetDefaultDomainClaimedMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomDomainClaimedMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomDomainClaimedMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomDomainClaimedMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomDomainClaimedMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomDomainClaimedMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomDomainClaimedMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomDomainClaimedMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomDomainClaimedMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomDomainClaimedMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomDomainClaimedMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomDomainClaimedMessageTextRequest): SetCustomDomainClaimedMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomDomainClaimedMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomDomainClaimedMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomDomainClaimedMessageTextRequest, reader: jspb.BinaryReader): SetCustomDomainClaimedMessageTextRequest;
}

export namespace SetCustomDomainClaimedMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomDomainClaimedMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomDomainClaimedMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomDomainClaimedMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomDomainClaimedMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomDomainClaimedMessageTextResponse): SetCustomDomainClaimedMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomDomainClaimedMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomDomainClaimedMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomDomainClaimedMessageTextResponse, reader: jspb.BinaryReader): SetCustomDomainClaimedMessageTextResponse;
}

export namespace SetCustomDomainClaimedMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomDomainClaimedMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomDomainClaimedMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainClaimedMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainClaimedMessageTextToDefaultRequest): ResetCustomDomainClaimedMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainClaimedMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainClaimedMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomDomainClaimedMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomDomainClaimedMessageTextToDefaultRequest;
}

export namespace ResetCustomDomainClaimedMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomDomainClaimedMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomDomainClaimedMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomDomainClaimedMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomDomainClaimedMessageTextToDefaultResponse): ResetCustomDomainClaimedMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomDomainClaimedMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomDomainClaimedMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomDomainClaimedMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomDomainClaimedMessageTextToDefaultResponse;
}

export namespace ResetCustomDomainClaimedMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordlessRegistrationMessageTextRequest): GetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordlessRegistrationMessageTextRequest;
}

export namespace GetCustomPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordlessRegistrationMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordlessRegistrationMessageTextResponse): GetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordlessRegistrationMessageTextResponse;
}

export namespace GetCustomPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordlessRegistrationMessageTextRequest): GetDefaultPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordlessRegistrationMessageTextRequest;
}

export namespace GetDefaultPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordlessRegistrationMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordlessRegistrationMessageTextResponse): GetDefaultPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordlessRegistrationMessageTextResponse;
}

export namespace GetDefaultPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomPasswordlessRegistrationMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomPasswordlessRegistrationMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordlessRegistrationMessageTextRequest): SetCustomPasswordlessRegistrationMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordlessRegistrationMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordlessRegistrationMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomPasswordlessRegistrationMessageTextRequest, reader: jspb.BinaryReader): SetCustomPasswordlessRegistrationMessageTextRequest;
}

export namespace SetCustomPasswordlessRegistrationMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomPasswordlessRegistrationMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomPasswordlessRegistrationMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomPasswordlessRegistrationMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordlessRegistrationMessageTextResponse): SetCustomPasswordlessRegistrationMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordlessRegistrationMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordlessRegistrationMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomPasswordlessRegistrationMessageTextResponse, reader: jspb.BinaryReader): SetCustomPasswordlessRegistrationMessageTextResponse;
}

export namespace SetCustomPasswordlessRegistrationMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordlessRegistrationMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordlessRegistrationMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetCustomPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetCustomPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordChangeMessageTextRequest): GetCustomPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: GetCustomPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): GetCustomPasswordChangeMessageTextRequest;
}

export namespace GetCustomPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetCustomPasswordChangeMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetCustomPasswordChangeMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetCustomPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCustomPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCustomPasswordChangeMessageTextResponse): GetCustomPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetCustomPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCustomPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: GetCustomPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): GetCustomPasswordChangeMessageTextResponse;
}

export namespace GetCustomPasswordChangeMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class GetDefaultPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): GetDefaultPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordChangeMessageTextRequest): GetDefaultPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: GetDefaultPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): GetDefaultPasswordChangeMessageTextRequest;
}

export namespace GetDefaultPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
  }
}

export class GetDefaultPasswordChangeMessageTextResponse extends jspb.Message {
  getCustomText(): zitadel_text_pb.MessageCustomText | undefined;
  setCustomText(value?: zitadel_text_pb.MessageCustomText): GetDefaultPasswordChangeMessageTextResponse;
  hasCustomText(): boolean;
  clearCustomText(): GetDefaultPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDefaultPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDefaultPasswordChangeMessageTextResponse): GetDefaultPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: GetDefaultPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDefaultPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: GetDefaultPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): GetDefaultPasswordChangeMessageTextResponse;
}

export namespace GetDefaultPasswordChangeMessageTextResponse {
  export type AsObject = {
    customText?: zitadel_text_pb.MessageCustomText.AsObject,
  }
}

export class SetCustomPasswordChangeMessageTextRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): SetCustomPasswordChangeMessageTextRequest;

  getTitle(): string;
  setTitle(value: string): SetCustomPasswordChangeMessageTextRequest;

  getPreHeader(): string;
  setPreHeader(value: string): SetCustomPasswordChangeMessageTextRequest;

  getSubject(): string;
  setSubject(value: string): SetCustomPasswordChangeMessageTextRequest;

  getGreeting(): string;
  setGreeting(value: string): SetCustomPasswordChangeMessageTextRequest;

  getText(): string;
  setText(value: string): SetCustomPasswordChangeMessageTextRequest;

  getButtonText(): string;
  setButtonText(value: string): SetCustomPasswordChangeMessageTextRequest;

  getFooterText(): string;
  setFooterText(value: string): SetCustomPasswordChangeMessageTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordChangeMessageTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordChangeMessageTextRequest): SetCustomPasswordChangeMessageTextRequest.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordChangeMessageTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordChangeMessageTextRequest;
  static deserializeBinaryFromReader(message: SetCustomPasswordChangeMessageTextRequest, reader: jspb.BinaryReader): SetCustomPasswordChangeMessageTextRequest;
}

export namespace SetCustomPasswordChangeMessageTextRequest {
  export type AsObject = {
    language: string,
    title: string,
    preHeader: string,
    subject: string,
    greeting: string,
    text: string,
    buttonText: string,
    footerText: string,
  }
}

export class SetCustomPasswordChangeMessageTextResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetCustomPasswordChangeMessageTextResponse;
  hasDetails(): boolean;
  clearDetails(): SetCustomPasswordChangeMessageTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetCustomPasswordChangeMessageTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetCustomPasswordChangeMessageTextResponse): SetCustomPasswordChangeMessageTextResponse.AsObject;
  static serializeBinaryToWriter(message: SetCustomPasswordChangeMessageTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetCustomPasswordChangeMessageTextResponse;
  static deserializeBinaryFromReader(message: SetCustomPasswordChangeMessageTextResponse, reader: jspb.BinaryReader): SetCustomPasswordChangeMessageTextResponse;
}

export namespace SetCustomPasswordChangeMessageTextResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ResetCustomPasswordChangeMessageTextToDefaultRequest extends jspb.Message {
  getLanguage(): string;
  setLanguage(value: string): ResetCustomPasswordChangeMessageTextToDefaultRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordChangeMessageTextToDefaultRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordChangeMessageTextToDefaultRequest): ResetCustomPasswordChangeMessageTextToDefaultRequest.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordChangeMessageTextToDefaultRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordChangeMessageTextToDefaultRequest;
  static deserializeBinaryFromReader(message: ResetCustomPasswordChangeMessageTextToDefaultRequest, reader: jspb.BinaryReader): ResetCustomPasswordChangeMessageTextToDefaultRequest;
}

export namespace ResetCustomPasswordChangeMessageTextToDefaultRequest {
  export type AsObject = {
    language: string,
  }
}

export class ResetCustomPasswordChangeMessageTextToDefaultResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetCustomPasswordChangeMessageTextToDefaultResponse;
  hasDetails(): boolean;
  clearDetails(): ResetCustomPasswordChangeMessageTextToDefaultResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetCustomPasswordChangeMessageTextToDefaultResponse): ResetCustomPasswordChangeMessageTextToDefaultResponse.AsObject;
  static serializeBinaryToWriter(message: ResetCustomPasswordChangeMessageTextToDefaultResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetCustomPasswordChangeMessageTextToDefaultResponse;
  static deserializeBinaryFromReader(message: ResetCustomPasswordChangeMessageTextToDefaultResponse, reader: jspb.BinaryReader): ResetCustomPasswordChangeMessageTextToDefaultResponse;
}

export namespace ResetCustomPasswordChangeMessageTextToDefaultResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetOrgIDPByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetOrgIDPByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIDPByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIDPByIDRequest): GetOrgIDPByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetOrgIDPByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIDPByIDRequest;
  static deserializeBinaryFromReader(message: GetOrgIDPByIDRequest, reader: jspb.BinaryReader): GetOrgIDPByIDRequest;
}

export namespace GetOrgIDPByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetOrgIDPByIDResponse extends jspb.Message {
  getIdp(): zitadel_idp_pb.IDP | undefined;
  setIdp(value?: zitadel_idp_pb.IDP): GetOrgIDPByIDResponse;
  hasIdp(): boolean;
  clearIdp(): GetOrgIDPByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetOrgIDPByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetOrgIDPByIDResponse): GetOrgIDPByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetOrgIDPByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetOrgIDPByIDResponse;
  static deserializeBinaryFromReader(message: GetOrgIDPByIDResponse, reader: jspb.BinaryReader): GetOrgIDPByIDResponse;
}

export namespace GetOrgIDPByIDResponse {
  export type AsObject = {
    idp?: zitadel_idp_pb.IDP.AsObject,
  }
}

export class ListOrgIDPsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListOrgIDPsRequest;
  hasQuery(): boolean;
  clearQuery(): ListOrgIDPsRequest;

  getSortingColumn(): zitadel_idp_pb.IDPFieldName;
  setSortingColumn(value: zitadel_idp_pb.IDPFieldName): ListOrgIDPsRequest;

  getQueriesList(): Array<IDPQuery>;
  setQueriesList(value: Array<IDPQuery>): ListOrgIDPsRequest;
  clearQueriesList(): ListOrgIDPsRequest;
  addQueries(value?: IDPQuery, index?: number): IDPQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgIDPsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgIDPsRequest): ListOrgIDPsRequest.AsObject;
  static serializeBinaryToWriter(message: ListOrgIDPsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgIDPsRequest;
  static deserializeBinaryFromReader(message: ListOrgIDPsRequest, reader: jspb.BinaryReader): ListOrgIDPsRequest;
}

export namespace ListOrgIDPsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_idp_pb.IDPFieldName,
    queriesList: Array<IDPQuery.AsObject>,
  }
}

export class IDPQuery extends jspb.Message {
  getIdpIdQuery(): zitadel_idp_pb.IDPIDQuery | undefined;
  setIdpIdQuery(value?: zitadel_idp_pb.IDPIDQuery): IDPQuery;
  hasIdpIdQuery(): boolean;
  clearIdpIdQuery(): IDPQuery;

  getIdpNameQuery(): zitadel_idp_pb.IDPNameQuery | undefined;
  setIdpNameQuery(value?: zitadel_idp_pb.IDPNameQuery): IDPQuery;
  hasIdpNameQuery(): boolean;
  clearIdpNameQuery(): IDPQuery;

  getOwnerTypeQuery(): zitadel_idp_pb.IDPOwnerTypeQuery | undefined;
  setOwnerTypeQuery(value?: zitadel_idp_pb.IDPOwnerTypeQuery): IDPQuery;
  hasOwnerTypeQuery(): boolean;
  clearOwnerTypeQuery(): IDPQuery;

  getQueryCase(): IDPQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IDPQuery.AsObject;
  static toObject(includeInstance: boolean, msg: IDPQuery): IDPQuery.AsObject;
  static serializeBinaryToWriter(message: IDPQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IDPQuery;
  static deserializeBinaryFromReader(message: IDPQuery, reader: jspb.BinaryReader): IDPQuery;
}

export namespace IDPQuery {
  export type AsObject = {
    idpIdQuery?: zitadel_idp_pb.IDPIDQuery.AsObject,
    idpNameQuery?: zitadel_idp_pb.IDPNameQuery.AsObject,
    ownerTypeQuery?: zitadel_idp_pb.IDPOwnerTypeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IDP_ID_QUERY = 1,
    IDP_NAME_QUERY = 2,
    OWNER_TYPE_QUERY = 3,
  }
}

export class ListOrgIDPsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListOrgIDPsResponse;
  hasDetails(): boolean;
  clearDetails(): ListOrgIDPsResponse;

  getSortingColumn(): zitadel_idp_pb.IDPFieldName;
  setSortingColumn(value: zitadel_idp_pb.IDPFieldName): ListOrgIDPsResponse;

  getResultList(): Array<zitadel_idp_pb.IDP>;
  setResultList(value: Array<zitadel_idp_pb.IDP>): ListOrgIDPsResponse;
  clearResultList(): ListOrgIDPsResponse;
  addResult(value?: zitadel_idp_pb.IDP, index?: number): zitadel_idp_pb.IDP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListOrgIDPsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListOrgIDPsResponse): ListOrgIDPsResponse.AsObject;
  static serializeBinaryToWriter(message: ListOrgIDPsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListOrgIDPsResponse;
  static deserializeBinaryFromReader(message: ListOrgIDPsResponse, reader: jspb.BinaryReader): ListOrgIDPsResponse;
}

export namespace ListOrgIDPsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_idp_pb.IDPFieldName,
    resultList: Array<zitadel_idp_pb.IDP.AsObject>,
  }
}

export class AddOrgOIDCIDPRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddOrgOIDCIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): AddOrgOIDCIDPRequest;

  getClientId(): string;
  setClientId(value: string): AddOrgOIDCIDPRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddOrgOIDCIDPRequest;

  getIssuer(): string;
  setIssuer(value: string): AddOrgOIDCIDPRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddOrgOIDCIDPRequest;
  clearScopesList(): AddOrgOIDCIDPRequest;
  addScopes(value: string, index?: number): AddOrgOIDCIDPRequest;

  getDisplayNameMapping(): zitadel_idp_pb.OIDCMappingField;
  setDisplayNameMapping(value: zitadel_idp_pb.OIDCMappingField): AddOrgOIDCIDPRequest;

  getUsernameMapping(): zitadel_idp_pb.OIDCMappingField;
  setUsernameMapping(value: zitadel_idp_pb.OIDCMappingField): AddOrgOIDCIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): AddOrgOIDCIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgOIDCIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgOIDCIDPRequest): AddOrgOIDCIDPRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgOIDCIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgOIDCIDPRequest;
  static deserializeBinaryFromReader(message: AddOrgOIDCIDPRequest, reader: jspb.BinaryReader): AddOrgOIDCIDPRequest;
}

export namespace AddOrgOIDCIDPRequest {
  export type AsObject = {
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    clientId: string,
    clientSecret: string,
    issuer: string,
    scopesList: Array<string>,
    displayNameMapping: zitadel_idp_pb.OIDCMappingField,
    usernameMapping: zitadel_idp_pb.OIDCMappingField,
    autoRegister: boolean,
  }
}

export class AddOrgOIDCIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOrgOIDCIDPResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrgOIDCIDPResponse;

  getIdpId(): string;
  setIdpId(value: string): AddOrgOIDCIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgOIDCIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgOIDCIDPResponse): AddOrgOIDCIDPResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrgOIDCIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgOIDCIDPResponse;
  static deserializeBinaryFromReader(message: AddOrgOIDCIDPResponse, reader: jspb.BinaryReader): AddOrgOIDCIDPResponse;
}

export namespace AddOrgOIDCIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    idpId: string,
  }
}

export class AddOrgJWTIDPRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddOrgJWTIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): AddOrgJWTIDPRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): AddOrgJWTIDPRequest;

  getIssuer(): string;
  setIssuer(value: string): AddOrgJWTIDPRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): AddOrgJWTIDPRequest;

  getHeaderName(): string;
  setHeaderName(value: string): AddOrgJWTIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): AddOrgJWTIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgJWTIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgJWTIDPRequest): AddOrgJWTIDPRequest.AsObject;
  static serializeBinaryToWriter(message: AddOrgJWTIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgJWTIDPRequest;
  static deserializeBinaryFromReader(message: AddOrgJWTIDPRequest, reader: jspb.BinaryReader): AddOrgJWTIDPRequest;
}

export namespace AddOrgJWTIDPRequest {
  export type AsObject = {
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    jwtEndpoint: string,
    issuer: string,
    keysEndpoint: string,
    headerName: string,
    autoRegister: boolean,
  }
}

export class AddOrgJWTIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddOrgJWTIDPResponse;
  hasDetails(): boolean;
  clearDetails(): AddOrgJWTIDPResponse;

  getIdpId(): string;
  setIdpId(value: string): AddOrgJWTIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddOrgJWTIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddOrgJWTIDPResponse): AddOrgJWTIDPResponse.AsObject;
  static serializeBinaryToWriter(message: AddOrgJWTIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddOrgJWTIDPResponse;
  static deserializeBinaryFromReader(message: AddOrgJWTIDPResponse, reader: jspb.BinaryReader): AddOrgJWTIDPResponse;
}

export namespace AddOrgJWTIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    idpId: string,
  }
}

export class DeactivateOrgIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): DeactivateOrgIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateOrgIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateOrgIDPRequest): DeactivateOrgIDPRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateOrgIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateOrgIDPRequest;
  static deserializeBinaryFromReader(message: DeactivateOrgIDPRequest, reader: jspb.BinaryReader): DeactivateOrgIDPRequest;
}

export namespace DeactivateOrgIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class DeactivateOrgIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateOrgIDPResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateOrgIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateOrgIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateOrgIDPResponse): DeactivateOrgIDPResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateOrgIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateOrgIDPResponse;
  static deserializeBinaryFromReader(message: DeactivateOrgIDPResponse, reader: jspb.BinaryReader): DeactivateOrgIDPResponse;
}

export namespace DeactivateOrgIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateOrgIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): ReactivateOrgIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateOrgIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateOrgIDPRequest): ReactivateOrgIDPRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateOrgIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateOrgIDPRequest;
  static deserializeBinaryFromReader(message: ReactivateOrgIDPRequest, reader: jspb.BinaryReader): ReactivateOrgIDPRequest;
}

export namespace ReactivateOrgIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class ReactivateOrgIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateOrgIDPResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateOrgIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateOrgIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateOrgIDPResponse): ReactivateOrgIDPResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateOrgIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateOrgIDPResponse;
  static deserializeBinaryFromReader(message: ReactivateOrgIDPResponse, reader: jspb.BinaryReader): ReactivateOrgIDPResponse;
}

export namespace ReactivateOrgIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveOrgIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): RemoveOrgIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgIDPRequest): RemoveOrgIDPRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgIDPRequest;
  static deserializeBinaryFromReader(message: RemoveOrgIDPRequest, reader: jspb.BinaryReader): RemoveOrgIDPRequest;
}

export namespace RemoveOrgIDPRequest {
  export type AsObject = {
    idpId: string,
  }
}

export class RemoveOrgIDPResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveOrgIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveOrgIDPResponse): RemoveOrgIDPResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveOrgIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveOrgIDPResponse;
  static deserializeBinaryFromReader(message: RemoveOrgIDPResponse, reader: jspb.BinaryReader): RemoveOrgIDPResponse;
}

export namespace RemoveOrgIDPResponse {
  export type AsObject = {
  }
}

export class UpdateOrgIDPRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateOrgIDPRequest;

  getName(): string;
  setName(value: string): UpdateOrgIDPRequest;

  getStylingType(): zitadel_idp_pb.IDPStylingType;
  setStylingType(value: zitadel_idp_pb.IDPStylingType): UpdateOrgIDPRequest;

  getAutoRegister(): boolean;
  setAutoRegister(value: boolean): UpdateOrgIDPRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPRequest): UpdateOrgIDPRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPRequest;
  static deserializeBinaryFromReader(message: UpdateOrgIDPRequest, reader: jspb.BinaryReader): UpdateOrgIDPRequest;
}

export namespace UpdateOrgIDPRequest {
  export type AsObject = {
    idpId: string,
    name: string,
    stylingType: zitadel_idp_pb.IDPStylingType,
    autoRegister: boolean,
  }
}

export class UpdateOrgIDPResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgIDPResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgIDPResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPResponse): UpdateOrgIDPResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPResponse;
  static deserializeBinaryFromReader(message: UpdateOrgIDPResponse, reader: jspb.BinaryReader): UpdateOrgIDPResponse;
}

export namespace UpdateOrgIDPResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOrgIDPOIDCConfigRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateOrgIDPOIDCConfigRequest;

  getClientId(): string;
  setClientId(value: string): UpdateOrgIDPOIDCConfigRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateOrgIDPOIDCConfigRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateOrgIDPOIDCConfigRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateOrgIDPOIDCConfigRequest;
  clearScopesList(): UpdateOrgIDPOIDCConfigRequest;
  addScopes(value: string, index?: number): UpdateOrgIDPOIDCConfigRequest;

  getDisplayNameMapping(): zitadel_idp_pb.OIDCMappingField;
  setDisplayNameMapping(value: zitadel_idp_pb.OIDCMappingField): UpdateOrgIDPOIDCConfigRequest;

  getUsernameMapping(): zitadel_idp_pb.OIDCMappingField;
  setUsernameMapping(value: zitadel_idp_pb.OIDCMappingField): UpdateOrgIDPOIDCConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPOIDCConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPOIDCConfigRequest): UpdateOrgIDPOIDCConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPOIDCConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPOIDCConfigRequest;
  static deserializeBinaryFromReader(message: UpdateOrgIDPOIDCConfigRequest, reader: jspb.BinaryReader): UpdateOrgIDPOIDCConfigRequest;
}

export namespace UpdateOrgIDPOIDCConfigRequest {
  export type AsObject = {
    idpId: string,
    clientId: string,
    clientSecret: string,
    issuer: string,
    scopesList: Array<string>,
    displayNameMapping: zitadel_idp_pb.OIDCMappingField,
    usernameMapping: zitadel_idp_pb.OIDCMappingField,
  }
}

export class UpdateOrgIDPOIDCConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgIDPOIDCConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgIDPOIDCConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPOIDCConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPOIDCConfigResponse): UpdateOrgIDPOIDCConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPOIDCConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPOIDCConfigResponse;
  static deserializeBinaryFromReader(message: UpdateOrgIDPOIDCConfigResponse, reader: jspb.BinaryReader): UpdateOrgIDPOIDCConfigResponse;
}

export namespace UpdateOrgIDPOIDCConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class UpdateOrgIDPJWTConfigRequest extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): UpdateOrgIDPJWTConfigRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): UpdateOrgIDPJWTConfigRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateOrgIDPJWTConfigRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): UpdateOrgIDPJWTConfigRequest;

  getHeaderName(): string;
  setHeaderName(value: string): UpdateOrgIDPJWTConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPJWTConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPJWTConfigRequest): UpdateOrgIDPJWTConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPJWTConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPJWTConfigRequest;
  static deserializeBinaryFromReader(message: UpdateOrgIDPJWTConfigRequest, reader: jspb.BinaryReader): UpdateOrgIDPJWTConfigRequest;
}

export namespace UpdateOrgIDPJWTConfigRequest {
  export type AsObject = {
    idpId: string,
    jwtEndpoint: string,
    issuer: string,
    keysEndpoint: string,
    headerName: string,
  }
}

export class UpdateOrgIDPJWTConfigResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateOrgIDPJWTConfigResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateOrgIDPJWTConfigResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateOrgIDPJWTConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateOrgIDPJWTConfigResponse): UpdateOrgIDPJWTConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateOrgIDPJWTConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateOrgIDPJWTConfigResponse;
  static deserializeBinaryFromReader(message: UpdateOrgIDPJWTConfigResponse, reader: jspb.BinaryReader): UpdateOrgIDPJWTConfigResponse;
}

export namespace UpdateOrgIDPJWTConfigResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListProvidersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListProvidersRequest;
  hasQuery(): boolean;
  clearQuery(): ListProvidersRequest;

  getQueriesList(): Array<ProviderQuery>;
  setQueriesList(value: Array<ProviderQuery>): ListProvidersRequest;
  clearQueriesList(): ListProvidersRequest;
  addQueries(value?: ProviderQuery, index?: number): ProviderQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProvidersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListProvidersRequest): ListProvidersRequest.AsObject;
  static serializeBinaryToWriter(message: ListProvidersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProvidersRequest;
  static deserializeBinaryFromReader(message: ListProvidersRequest, reader: jspb.BinaryReader): ListProvidersRequest;
}

export namespace ListProvidersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    queriesList: Array<ProviderQuery.AsObject>,
  }
}

export class ProviderQuery extends jspb.Message {
  getIdpIdQuery(): zitadel_idp_pb.IDPIDQuery | undefined;
  setIdpIdQuery(value?: zitadel_idp_pb.IDPIDQuery): ProviderQuery;
  hasIdpIdQuery(): boolean;
  clearIdpIdQuery(): ProviderQuery;

  getIdpNameQuery(): zitadel_idp_pb.IDPNameQuery | undefined;
  setIdpNameQuery(value?: zitadel_idp_pb.IDPNameQuery): ProviderQuery;
  hasIdpNameQuery(): boolean;
  clearIdpNameQuery(): ProviderQuery;

  getOwnerTypeQuery(): zitadel_idp_pb.IDPOwnerTypeQuery | undefined;
  setOwnerTypeQuery(value?: zitadel_idp_pb.IDPOwnerTypeQuery): ProviderQuery;
  hasOwnerTypeQuery(): boolean;
  clearOwnerTypeQuery(): ProviderQuery;

  getQueryCase(): ProviderQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProviderQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProviderQuery): ProviderQuery.AsObject;
  static serializeBinaryToWriter(message: ProviderQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProviderQuery;
  static deserializeBinaryFromReader(message: ProviderQuery, reader: jspb.BinaryReader): ProviderQuery;
}

export namespace ProviderQuery {
  export type AsObject = {
    idpIdQuery?: zitadel_idp_pb.IDPIDQuery.AsObject,
    idpNameQuery?: zitadel_idp_pb.IDPNameQuery.AsObject,
    ownerTypeQuery?: zitadel_idp_pb.IDPOwnerTypeQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    IDP_ID_QUERY = 1,
    IDP_NAME_QUERY = 2,
    OWNER_TYPE_QUERY = 3,
  }
}

export class ListProvidersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListProvidersResponse;
  hasDetails(): boolean;
  clearDetails(): ListProvidersResponse;

  getResultList(): Array<zitadel_idp_pb.Provider>;
  setResultList(value: Array<zitadel_idp_pb.Provider>): ListProvidersResponse;
  clearResultList(): ListProvidersResponse;
  addResult(value?: zitadel_idp_pb.Provider, index?: number): zitadel_idp_pb.Provider;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListProvidersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListProvidersResponse): ListProvidersResponse.AsObject;
  static serializeBinaryToWriter(message: ListProvidersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListProvidersResponse;
  static deserializeBinaryFromReader(message: ListProvidersResponse, reader: jspb.BinaryReader): ListProvidersResponse;
}

export namespace ListProvidersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_idp_pb.Provider.AsObject>,
  }
}

export class GetProviderByIDRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetProviderByIDRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProviderByIDRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProviderByIDRequest): GetProviderByIDRequest.AsObject;
  static serializeBinaryToWriter(message: GetProviderByIDRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProviderByIDRequest;
  static deserializeBinaryFromReader(message: GetProviderByIDRequest, reader: jspb.BinaryReader): GetProviderByIDRequest;
}

export namespace GetProviderByIDRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetProviderByIDResponse extends jspb.Message {
  getIdp(): zitadel_idp_pb.Provider | undefined;
  setIdp(value?: zitadel_idp_pb.Provider): GetProviderByIDResponse;
  hasIdp(): boolean;
  clearIdp(): GetProviderByIDResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProviderByIDResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProviderByIDResponse): GetProviderByIDResponse.AsObject;
  static serializeBinaryToWriter(message: GetProviderByIDResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProviderByIDResponse;
  static deserializeBinaryFromReader(message: GetProviderByIDResponse, reader: jspb.BinaryReader): GetProviderByIDResponse;
}

export namespace GetProviderByIDResponse {
  export type AsObject = {
    idp?: zitadel_idp_pb.Provider.AsObject,
  }
}

export class AddGenericOAuthProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGenericOAuthProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGenericOAuthProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGenericOAuthProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): AddGenericOAuthProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): AddGenericOAuthProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): AddGenericOAuthProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGenericOAuthProviderRequest;
  clearScopesList(): AddGenericOAuthProviderRequest;
  addScopes(value: string, index?: number): AddGenericOAuthProviderRequest;

  getIdAttribute(): string;
  setIdAttribute(value: string): AddGenericOAuthProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGenericOAuthProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGenericOAuthProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOAuthProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOAuthProviderRequest): AddGenericOAuthProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGenericOAuthProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOAuthProviderRequest;
  static deserializeBinaryFromReader(message: AddGenericOAuthProviderRequest, reader: jspb.BinaryReader): AddGenericOAuthProviderRequest;
}

export namespace AddGenericOAuthProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    idAttribute: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGenericOAuthProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGenericOAuthProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGenericOAuthProviderResponse;

  getId(): string;
  setId(value: string): AddGenericOAuthProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOAuthProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOAuthProviderResponse): AddGenericOAuthProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGenericOAuthProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOAuthProviderResponse;
  static deserializeBinaryFromReader(message: AddGenericOAuthProviderResponse, reader: jspb.BinaryReader): AddGenericOAuthProviderResponse;
}

export namespace AddGenericOAuthProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGenericOAuthProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGenericOAuthProviderRequest;

  getName(): string;
  setName(value: string): UpdateGenericOAuthProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGenericOAuthProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGenericOAuthProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): UpdateGenericOAuthProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGenericOAuthProviderRequest;
  clearScopesList(): UpdateGenericOAuthProviderRequest;
  addScopes(value: string, index?: number): UpdateGenericOAuthProviderRequest;

  getIdAttribute(): string;
  setIdAttribute(value: string): UpdateGenericOAuthProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGenericOAuthProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGenericOAuthProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOAuthProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOAuthProviderRequest): UpdateGenericOAuthProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOAuthProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOAuthProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGenericOAuthProviderRequest, reader: jspb.BinaryReader): UpdateGenericOAuthProviderRequest;
}

export namespace UpdateGenericOAuthProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    idAttribute: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGenericOAuthProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGenericOAuthProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGenericOAuthProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOAuthProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOAuthProviderResponse): UpdateGenericOAuthProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOAuthProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOAuthProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGenericOAuthProviderResponse, reader: jspb.BinaryReader): UpdateGenericOAuthProviderResponse;
}

export namespace UpdateGenericOAuthProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGenericOIDCProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGenericOIDCProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): AddGenericOIDCProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGenericOIDCProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGenericOIDCProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGenericOIDCProviderRequest;
  clearScopesList(): AddGenericOIDCProviderRequest;
  addScopes(value: string, index?: number): AddGenericOIDCProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGenericOIDCProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGenericOIDCProviderRequest;

  getIsIdTokenMapping(): boolean;
  setIsIdTokenMapping(value: boolean): AddGenericOIDCProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOIDCProviderRequest): AddGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: AddGenericOIDCProviderRequest, reader: jspb.BinaryReader): AddGenericOIDCProviderRequest;
}

export namespace AddGenericOIDCProviderRequest {
  export type AsObject = {
    name: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    isIdTokenMapping: boolean,
  }
}

export class AddGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGenericOIDCProviderResponse;

  getId(): string;
  setId(value: string): AddGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGenericOIDCProviderResponse): AddGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: AddGenericOIDCProviderResponse, reader: jspb.BinaryReader): AddGenericOIDCProviderResponse;
}

export namespace AddGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGenericOIDCProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGenericOIDCProviderRequest;

  getName(): string;
  setName(value: string): UpdateGenericOIDCProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateGenericOIDCProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGenericOIDCProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGenericOIDCProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGenericOIDCProviderRequest;
  clearScopesList(): UpdateGenericOIDCProviderRequest;
  addScopes(value: string, index?: number): UpdateGenericOIDCProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGenericOIDCProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGenericOIDCProviderRequest;

  getIsIdTokenMapping(): boolean;
  setIsIdTokenMapping(value: boolean): UpdateGenericOIDCProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOIDCProviderRequest): UpdateGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGenericOIDCProviderRequest, reader: jspb.BinaryReader): UpdateGenericOIDCProviderRequest;
}

export namespace UpdateGenericOIDCProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    issuer: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    isIdTokenMapping: boolean,
  }
}

export class UpdateGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGenericOIDCProviderResponse): UpdateGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGenericOIDCProviderResponse, reader: jspb.BinaryReader): UpdateGenericOIDCProviderResponse;
}

export namespace UpdateGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class MigrateGenericOIDCProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): MigrateGenericOIDCProviderRequest;

  getAzure(): AddAzureADProviderRequest | undefined;
  setAzure(value?: AddAzureADProviderRequest): MigrateGenericOIDCProviderRequest;
  hasAzure(): boolean;
  clearAzure(): MigrateGenericOIDCProviderRequest;

  getGoogle(): AddGoogleProviderRequest | undefined;
  setGoogle(value?: AddGoogleProviderRequest): MigrateGenericOIDCProviderRequest;
  hasGoogle(): boolean;
  clearGoogle(): MigrateGenericOIDCProviderRequest;

  getTemplateCase(): MigrateGenericOIDCProviderRequest.TemplateCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateGenericOIDCProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateGenericOIDCProviderRequest): MigrateGenericOIDCProviderRequest.AsObject;
  static serializeBinaryToWriter(message: MigrateGenericOIDCProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateGenericOIDCProviderRequest;
  static deserializeBinaryFromReader(message: MigrateGenericOIDCProviderRequest, reader: jspb.BinaryReader): MigrateGenericOIDCProviderRequest;
}

export namespace MigrateGenericOIDCProviderRequest {
  export type AsObject = {
    id: string,
    azure?: AddAzureADProviderRequest.AsObject,
    google?: AddGoogleProviderRequest.AsObject,
  }

  export enum TemplateCase { 
    TEMPLATE_NOT_SET = 0,
    AZURE = 2,
    GOOGLE = 3,
  }
}

export class MigrateGenericOIDCProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): MigrateGenericOIDCProviderResponse;
  hasDetails(): boolean;
  clearDetails(): MigrateGenericOIDCProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): MigrateGenericOIDCProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: MigrateGenericOIDCProviderResponse): MigrateGenericOIDCProviderResponse.AsObject;
  static serializeBinaryToWriter(message: MigrateGenericOIDCProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): MigrateGenericOIDCProviderResponse;
  static deserializeBinaryFromReader(message: MigrateGenericOIDCProviderResponse, reader: jspb.BinaryReader): MigrateGenericOIDCProviderResponse;
}

export namespace MigrateGenericOIDCProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddJWTProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddJWTProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): AddJWTProviderRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): AddJWTProviderRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): AddJWTProviderRequest;

  getHeaderName(): string;
  setHeaderName(value: string): AddJWTProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddJWTProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddJWTProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTProviderRequest): AddJWTProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddJWTProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTProviderRequest;
  static deserializeBinaryFromReader(message: AddJWTProviderRequest, reader: jspb.BinaryReader): AddJWTProviderRequest;
}

export namespace AddJWTProviderRequest {
  export type AsObject = {
    name: string,
    issuer: string,
    jwtEndpoint: string,
    keysEndpoint: string,
    headerName: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddJWTProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddJWTProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddJWTProviderResponse;

  getId(): string;
  setId(value: string): AddJWTProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddJWTProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddJWTProviderResponse): AddJWTProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddJWTProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddJWTProviderResponse;
  static deserializeBinaryFromReader(message: AddJWTProviderResponse, reader: jspb.BinaryReader): AddJWTProviderResponse;
}

export namespace AddJWTProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateJWTProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateJWTProviderRequest;

  getName(): string;
  setName(value: string): UpdateJWTProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateJWTProviderRequest;

  getJwtEndpoint(): string;
  setJwtEndpoint(value: string): UpdateJWTProviderRequest;

  getKeysEndpoint(): string;
  setKeysEndpoint(value: string): UpdateJWTProviderRequest;

  getHeaderName(): string;
  setHeaderName(value: string): UpdateJWTProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateJWTProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateJWTProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateJWTProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateJWTProviderRequest): UpdateJWTProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateJWTProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateJWTProviderRequest;
  static deserializeBinaryFromReader(message: UpdateJWTProviderRequest, reader: jspb.BinaryReader): UpdateJWTProviderRequest;
}

export namespace UpdateJWTProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    issuer: string,
    jwtEndpoint: string,
    keysEndpoint: string,
    headerName: string,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateJWTProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateJWTProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateJWTProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateJWTProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateJWTProviderResponse): UpdateJWTProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateJWTProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateJWTProviderResponse;
  static deserializeBinaryFromReader(message: UpdateJWTProviderResponse, reader: jspb.BinaryReader): UpdateJWTProviderResponse;
}

export namespace UpdateJWTProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddAzureADProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddAzureADProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddAzureADProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddAzureADProviderRequest;

  getTenant(): zitadel_idp_pb.AzureADTenant | undefined;
  setTenant(value?: zitadel_idp_pb.AzureADTenant): AddAzureADProviderRequest;
  hasTenant(): boolean;
  clearTenant(): AddAzureADProviderRequest;

  getEmailVerified(): boolean;
  setEmailVerified(value: boolean): AddAzureADProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddAzureADProviderRequest;
  clearScopesList(): AddAzureADProviderRequest;
  addScopes(value: string, index?: number): AddAzureADProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddAzureADProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddAzureADProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAzureADProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAzureADProviderRequest): AddAzureADProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddAzureADProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAzureADProviderRequest;
  static deserializeBinaryFromReader(message: AddAzureADProviderRequest, reader: jspb.BinaryReader): AddAzureADProviderRequest;
}

export namespace AddAzureADProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    tenant?: zitadel_idp_pb.AzureADTenant.AsObject,
    emailVerified: boolean,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddAzureADProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAzureADProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddAzureADProviderResponse;

  getId(): string;
  setId(value: string): AddAzureADProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAzureADProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAzureADProviderResponse): AddAzureADProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddAzureADProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAzureADProviderResponse;
  static deserializeBinaryFromReader(message: AddAzureADProviderResponse, reader: jspb.BinaryReader): AddAzureADProviderResponse;
}

export namespace AddAzureADProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateAzureADProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateAzureADProviderRequest;

  getName(): string;
  setName(value: string): UpdateAzureADProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateAzureADProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateAzureADProviderRequest;

  getTenant(): zitadel_idp_pb.AzureADTenant | undefined;
  setTenant(value?: zitadel_idp_pb.AzureADTenant): UpdateAzureADProviderRequest;
  hasTenant(): boolean;
  clearTenant(): UpdateAzureADProviderRequest;

  getEmailVerified(): boolean;
  setEmailVerified(value: boolean): UpdateAzureADProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateAzureADProviderRequest;
  clearScopesList(): UpdateAzureADProviderRequest;
  addScopes(value: string, index?: number): UpdateAzureADProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateAzureADProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateAzureADProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAzureADProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAzureADProviderRequest): UpdateAzureADProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAzureADProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAzureADProviderRequest;
  static deserializeBinaryFromReader(message: UpdateAzureADProviderRequest, reader: jspb.BinaryReader): UpdateAzureADProviderRequest;
}

export namespace UpdateAzureADProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    tenant?: zitadel_idp_pb.AzureADTenant.AsObject,
    emailVerified: boolean,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateAzureADProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAzureADProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAzureADProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAzureADProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAzureADProviderResponse): UpdateAzureADProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAzureADProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAzureADProviderResponse;
  static deserializeBinaryFromReader(message: UpdateAzureADProviderResponse, reader: jspb.BinaryReader): UpdateAzureADProviderResponse;
}

export namespace UpdateAzureADProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitHubProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGitHubProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitHubProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitHubProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitHubProviderRequest;
  clearScopesList(): AddGitHubProviderRequest;
  addScopes(value: string, index?: number): AddGitHubProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitHubProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitHubProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubProviderRequest): AddGitHubProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitHubProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubProviderRequest;
  static deserializeBinaryFromReader(message: AddGitHubProviderRequest, reader: jspb.BinaryReader): AddGitHubProviderRequest;
}

export namespace AddGitHubProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitHubProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitHubProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitHubProviderResponse;

  getId(): string;
  setId(value: string): AddGitHubProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubProviderResponse): AddGitHubProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitHubProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubProviderResponse;
  static deserializeBinaryFromReader(message: AddGitHubProviderResponse, reader: jspb.BinaryReader): AddGitHubProviderResponse;
}

export namespace AddGitHubProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitHubProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitHubProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitHubProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitHubProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitHubProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitHubProviderRequest;
  clearScopesList(): UpdateGitHubProviderRequest;
  addScopes(value: string, index?: number): UpdateGitHubProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitHubProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitHubProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubProviderRequest): UpdateGitHubProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitHubProviderRequest, reader: jspb.BinaryReader): UpdateGitHubProviderRequest;
}

export namespace UpdateGitHubProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitHubProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitHubProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitHubProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubProviderResponse): UpdateGitHubProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitHubProviderResponse, reader: jspb.BinaryReader): UpdateGitHubProviderResponse;
}

export namespace UpdateGitHubProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitHubEnterpriseServerProviderRequest extends jspb.Message {
  getClientId(): string;
  setClientId(value: string): AddGitHubEnterpriseServerProviderRequest;

  getName(): string;
  setName(value: string): AddGitHubEnterpriseServerProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitHubEnterpriseServerProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): AddGitHubEnterpriseServerProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitHubEnterpriseServerProviderRequest;
  clearScopesList(): AddGitHubEnterpriseServerProviderRequest;
  addScopes(value: string, index?: number): AddGitHubEnterpriseServerProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitHubEnterpriseServerProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitHubEnterpriseServerProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubEnterpriseServerProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubEnterpriseServerProviderRequest): AddGitHubEnterpriseServerProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitHubEnterpriseServerProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubEnterpriseServerProviderRequest;
  static deserializeBinaryFromReader(message: AddGitHubEnterpriseServerProviderRequest, reader: jspb.BinaryReader): AddGitHubEnterpriseServerProviderRequest;
}

export namespace AddGitHubEnterpriseServerProviderRequest {
  export type AsObject = {
    clientId: string,
    name: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitHubEnterpriseServerProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitHubEnterpriseServerProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitHubEnterpriseServerProviderResponse;

  getId(): string;
  setId(value: string): AddGitHubEnterpriseServerProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitHubEnterpriseServerProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitHubEnterpriseServerProviderResponse): AddGitHubEnterpriseServerProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitHubEnterpriseServerProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitHubEnterpriseServerProviderResponse;
  static deserializeBinaryFromReader(message: AddGitHubEnterpriseServerProviderResponse, reader: jspb.BinaryReader): AddGitHubEnterpriseServerProviderResponse;
}

export namespace AddGitHubEnterpriseServerProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitHubEnterpriseServerProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getAuthorizationEndpoint(): string;
  setAuthorizationEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getTokenEndpoint(): string;
  setTokenEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getUserEndpoint(): string;
  setUserEndpoint(value: string): UpdateGitHubEnterpriseServerProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitHubEnterpriseServerProviderRequest;
  clearScopesList(): UpdateGitHubEnterpriseServerProviderRequest;
  addScopes(value: string, index?: number): UpdateGitHubEnterpriseServerProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitHubEnterpriseServerProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitHubEnterpriseServerProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubEnterpriseServerProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubEnterpriseServerProviderRequest): UpdateGitHubEnterpriseServerProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubEnterpriseServerProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubEnterpriseServerProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitHubEnterpriseServerProviderRequest, reader: jspb.BinaryReader): UpdateGitHubEnterpriseServerProviderRequest;
}

export namespace UpdateGitHubEnterpriseServerProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    authorizationEndpoint: string,
    tokenEndpoint: string,
    userEndpoint: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitHubEnterpriseServerProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitHubEnterpriseServerProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitHubEnterpriseServerProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitHubEnterpriseServerProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitHubEnterpriseServerProviderResponse): UpdateGitHubEnterpriseServerProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitHubEnterpriseServerProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitHubEnterpriseServerProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitHubEnterpriseServerProviderResponse, reader: jspb.BinaryReader): UpdateGitHubEnterpriseServerProviderResponse;
}

export namespace UpdateGitHubEnterpriseServerProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitLabProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGitLabProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitLabProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitLabProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitLabProviderRequest;
  clearScopesList(): AddGitLabProviderRequest;
  addScopes(value: string, index?: number): AddGitLabProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitLabProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitLabProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabProviderRequest): AddGitLabProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitLabProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabProviderRequest;
  static deserializeBinaryFromReader(message: AddGitLabProviderRequest, reader: jspb.BinaryReader): AddGitLabProviderRequest;
}

export namespace AddGitLabProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitLabProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitLabProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitLabProviderResponse;

  getId(): string;
  setId(value: string): AddGitLabProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabProviderResponse): AddGitLabProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitLabProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabProviderResponse;
  static deserializeBinaryFromReader(message: AddGitLabProviderResponse, reader: jspb.BinaryReader): AddGitLabProviderResponse;
}

export namespace AddGitLabProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitLabProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitLabProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitLabProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitLabProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitLabProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitLabProviderRequest;
  clearScopesList(): UpdateGitLabProviderRequest;
  addScopes(value: string, index?: number): UpdateGitLabProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitLabProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitLabProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabProviderRequest): UpdateGitLabProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitLabProviderRequest, reader: jspb.BinaryReader): UpdateGitLabProviderRequest;
}

export namespace UpdateGitLabProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitLabProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitLabProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitLabProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabProviderResponse): UpdateGitLabProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitLabProviderResponse, reader: jspb.BinaryReader): UpdateGitLabProviderResponse;
}

export namespace UpdateGitLabProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGitLabSelfHostedProviderRequest extends jspb.Message {
  getIssuer(): string;
  setIssuer(value: string): AddGitLabSelfHostedProviderRequest;

  getName(): string;
  setName(value: string): AddGitLabSelfHostedProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGitLabSelfHostedProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGitLabSelfHostedProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGitLabSelfHostedProviderRequest;
  clearScopesList(): AddGitLabSelfHostedProviderRequest;
  addScopes(value: string, index?: number): AddGitLabSelfHostedProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGitLabSelfHostedProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGitLabSelfHostedProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabSelfHostedProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabSelfHostedProviderRequest): AddGitLabSelfHostedProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGitLabSelfHostedProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabSelfHostedProviderRequest;
  static deserializeBinaryFromReader(message: AddGitLabSelfHostedProviderRequest, reader: jspb.BinaryReader): AddGitLabSelfHostedProviderRequest;
}

export namespace AddGitLabSelfHostedProviderRequest {
  export type AsObject = {
    issuer: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGitLabSelfHostedProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGitLabSelfHostedProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGitLabSelfHostedProviderResponse;

  getId(): string;
  setId(value: string): AddGitLabSelfHostedProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGitLabSelfHostedProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGitLabSelfHostedProviderResponse): AddGitLabSelfHostedProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGitLabSelfHostedProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGitLabSelfHostedProviderResponse;
  static deserializeBinaryFromReader(message: AddGitLabSelfHostedProviderResponse, reader: jspb.BinaryReader): AddGitLabSelfHostedProviderResponse;
}

export namespace AddGitLabSelfHostedProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGitLabSelfHostedProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGitLabSelfHostedProviderRequest;

  getIssuer(): string;
  setIssuer(value: string): UpdateGitLabSelfHostedProviderRequest;

  getName(): string;
  setName(value: string): UpdateGitLabSelfHostedProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGitLabSelfHostedProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGitLabSelfHostedProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGitLabSelfHostedProviderRequest;
  clearScopesList(): UpdateGitLabSelfHostedProviderRequest;
  addScopes(value: string, index?: number): UpdateGitLabSelfHostedProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGitLabSelfHostedProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGitLabSelfHostedProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabSelfHostedProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabSelfHostedProviderRequest): UpdateGitLabSelfHostedProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabSelfHostedProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabSelfHostedProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGitLabSelfHostedProviderRequest, reader: jspb.BinaryReader): UpdateGitLabSelfHostedProviderRequest;
}

export namespace UpdateGitLabSelfHostedProviderRequest {
  export type AsObject = {
    id: string,
    issuer: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGitLabSelfHostedProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGitLabSelfHostedProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGitLabSelfHostedProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGitLabSelfHostedProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGitLabSelfHostedProviderResponse): UpdateGitLabSelfHostedProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGitLabSelfHostedProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGitLabSelfHostedProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGitLabSelfHostedProviderResponse, reader: jspb.BinaryReader): UpdateGitLabSelfHostedProviderResponse;
}

export namespace UpdateGitLabSelfHostedProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddGoogleProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddGoogleProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddGoogleProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): AddGoogleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddGoogleProviderRequest;
  clearScopesList(): AddGoogleProviderRequest;
  addScopes(value: string, index?: number): AddGoogleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddGoogleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddGoogleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGoogleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddGoogleProviderRequest): AddGoogleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddGoogleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGoogleProviderRequest;
  static deserializeBinaryFromReader(message: AddGoogleProviderRequest, reader: jspb.BinaryReader): AddGoogleProviderRequest;
}

export namespace AddGoogleProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddGoogleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddGoogleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddGoogleProviderResponse;

  getId(): string;
  setId(value: string): AddGoogleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddGoogleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddGoogleProviderResponse): AddGoogleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddGoogleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddGoogleProviderResponse;
  static deserializeBinaryFromReader(message: AddGoogleProviderResponse, reader: jspb.BinaryReader): AddGoogleProviderResponse;
}

export namespace AddGoogleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateGoogleProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateGoogleProviderRequest;

  getName(): string;
  setName(value: string): UpdateGoogleProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateGoogleProviderRequest;

  getClientSecret(): string;
  setClientSecret(value: string): UpdateGoogleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateGoogleProviderRequest;
  clearScopesList(): UpdateGoogleProviderRequest;
  addScopes(value: string, index?: number): UpdateGoogleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateGoogleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateGoogleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoogleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoogleProviderRequest): UpdateGoogleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateGoogleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoogleProviderRequest;
  static deserializeBinaryFromReader(message: UpdateGoogleProviderRequest, reader: jspb.BinaryReader): UpdateGoogleProviderRequest;
}

export namespace UpdateGoogleProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    clientSecret: string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateGoogleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateGoogleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateGoogleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateGoogleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateGoogleProviderResponse): UpdateGoogleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateGoogleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateGoogleProviderResponse;
  static deserializeBinaryFromReader(message: UpdateGoogleProviderResponse, reader: jspb.BinaryReader): UpdateGoogleProviderResponse;
}

export namespace UpdateGoogleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddLDAPProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddLDAPProviderRequest;

  getServersList(): Array<string>;
  setServersList(value: Array<string>): AddLDAPProviderRequest;
  clearServersList(): AddLDAPProviderRequest;
  addServers(value: string, index?: number): AddLDAPProviderRequest;

  getStartTls(): boolean;
  setStartTls(value: boolean): AddLDAPProviderRequest;

  getBaseDn(): string;
  setBaseDn(value: string): AddLDAPProviderRequest;

  getBindDn(): string;
  setBindDn(value: string): AddLDAPProviderRequest;

  getBindPassword(): string;
  setBindPassword(value: string): AddLDAPProviderRequest;

  getUserBase(): string;
  setUserBase(value: string): AddLDAPProviderRequest;

  getUserObjectClassesList(): Array<string>;
  setUserObjectClassesList(value: Array<string>): AddLDAPProviderRequest;
  clearUserObjectClassesList(): AddLDAPProviderRequest;
  addUserObjectClasses(value: string, index?: number): AddLDAPProviderRequest;

  getUserFiltersList(): Array<string>;
  setUserFiltersList(value: Array<string>): AddLDAPProviderRequest;
  clearUserFiltersList(): AddLDAPProviderRequest;
  addUserFilters(value: string, index?: number): AddLDAPProviderRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): AddLDAPProviderRequest;
  hasTimeout(): boolean;
  clearTimeout(): AddLDAPProviderRequest;

  getAttributes(): zitadel_idp_pb.LDAPAttributes | undefined;
  setAttributes(value?: zitadel_idp_pb.LDAPAttributes): AddLDAPProviderRequest;
  hasAttributes(): boolean;
  clearAttributes(): AddLDAPProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddLDAPProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddLDAPProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddLDAPProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddLDAPProviderRequest): AddLDAPProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddLDAPProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddLDAPProviderRequest;
  static deserializeBinaryFromReader(message: AddLDAPProviderRequest, reader: jspb.BinaryReader): AddLDAPProviderRequest;
}

export namespace AddLDAPProviderRequest {
  export type AsObject = {
    name: string,
    serversList: Array<string>,
    startTls: boolean,
    baseDn: string,
    bindDn: string,
    bindPassword: string,
    userBase: string,
    userObjectClassesList: Array<string>,
    userFiltersList: Array<string>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    attributes?: zitadel_idp_pb.LDAPAttributes.AsObject,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddLDAPProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddLDAPProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddLDAPProviderResponse;

  getId(): string;
  setId(value: string): AddLDAPProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddLDAPProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddLDAPProviderResponse): AddLDAPProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddLDAPProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddLDAPProviderResponse;
  static deserializeBinaryFromReader(message: AddLDAPProviderResponse, reader: jspb.BinaryReader): AddLDAPProviderResponse;
}

export namespace AddLDAPProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateLDAPProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateLDAPProviderRequest;

  getName(): string;
  setName(value: string): UpdateLDAPProviderRequest;

  getServersList(): Array<string>;
  setServersList(value: Array<string>): UpdateLDAPProviderRequest;
  clearServersList(): UpdateLDAPProviderRequest;
  addServers(value: string, index?: number): UpdateLDAPProviderRequest;

  getStartTls(): boolean;
  setStartTls(value: boolean): UpdateLDAPProviderRequest;

  getBaseDn(): string;
  setBaseDn(value: string): UpdateLDAPProviderRequest;

  getBindDn(): string;
  setBindDn(value: string): UpdateLDAPProviderRequest;

  getBindPassword(): string;
  setBindPassword(value: string): UpdateLDAPProviderRequest;

  getUserBase(): string;
  setUserBase(value: string): UpdateLDAPProviderRequest;

  getUserObjectClassesList(): Array<string>;
  setUserObjectClassesList(value: Array<string>): UpdateLDAPProviderRequest;
  clearUserObjectClassesList(): UpdateLDAPProviderRequest;
  addUserObjectClasses(value: string, index?: number): UpdateLDAPProviderRequest;

  getUserFiltersList(): Array<string>;
  setUserFiltersList(value: Array<string>): UpdateLDAPProviderRequest;
  clearUserFiltersList(): UpdateLDAPProviderRequest;
  addUserFilters(value: string, index?: number): UpdateLDAPProviderRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): UpdateLDAPProviderRequest;
  hasTimeout(): boolean;
  clearTimeout(): UpdateLDAPProviderRequest;

  getAttributes(): zitadel_idp_pb.LDAPAttributes | undefined;
  setAttributes(value?: zitadel_idp_pb.LDAPAttributes): UpdateLDAPProviderRequest;
  hasAttributes(): boolean;
  clearAttributes(): UpdateLDAPProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateLDAPProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateLDAPProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLDAPProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLDAPProviderRequest): UpdateLDAPProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateLDAPProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLDAPProviderRequest;
  static deserializeBinaryFromReader(message: UpdateLDAPProviderRequest, reader: jspb.BinaryReader): UpdateLDAPProviderRequest;
}

export namespace UpdateLDAPProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    serversList: Array<string>,
    startTls: boolean,
    baseDn: string,
    bindDn: string,
    bindPassword: string,
    userBase: string,
    userObjectClassesList: Array<string>,
    userFiltersList: Array<string>,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    attributes?: zitadel_idp_pb.LDAPAttributes.AsObject,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateLDAPProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateLDAPProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateLDAPProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateLDAPProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateLDAPProviderResponse): UpdateLDAPProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateLDAPProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateLDAPProviderResponse;
  static deserializeBinaryFromReader(message: UpdateLDAPProviderResponse, reader: jspb.BinaryReader): UpdateLDAPProviderResponse;
}

export namespace UpdateLDAPProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddSAMLProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddSAMLProviderRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): AddSAMLProviderRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): AddSAMLProviderRequest;

  getBinding(): zitadel_idp_pb.SAMLBinding;
  setBinding(value: zitadel_idp_pb.SAMLBinding): AddSAMLProviderRequest;

  getWithSignedRequest(): boolean;
  setWithSignedRequest(value: boolean): AddSAMLProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddSAMLProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddSAMLProviderRequest;

  getNameIdFormat(): zitadel_idp_pb.SAMLNameIDFormat;
  setNameIdFormat(value: zitadel_idp_pb.SAMLNameIDFormat): AddSAMLProviderRequest;
  hasNameIdFormat(): boolean;
  clearNameIdFormat(): AddSAMLProviderRequest;

  getTransientMappingAttributeName(): string;
  setTransientMappingAttributeName(value: string): AddSAMLProviderRequest;
  hasTransientMappingAttributeName(): boolean;
  clearTransientMappingAttributeName(): AddSAMLProviderRequest;

  getMetadataCase(): AddSAMLProviderRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLProviderRequest): AddSAMLProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddSAMLProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLProviderRequest;
  static deserializeBinaryFromReader(message: AddSAMLProviderRequest, reader: jspb.BinaryReader): AddSAMLProviderRequest;
}

export namespace AddSAMLProviderRequest {
  export type AsObject = {
    name: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
    binding: zitadel_idp_pb.SAMLBinding,
    withSignedRequest: boolean,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    nameIdFormat?: zitadel_idp_pb.SAMLNameIDFormat,
    transientMappingAttributeName?: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 2,
    METADATA_URL = 3,
  }

  export enum NameIdFormatCase { 
    _NAME_ID_FORMAT_NOT_SET = 0,
    NAME_ID_FORMAT = 7,
  }

  export enum TransientMappingAttributeNameCase { 
    _TRANSIENT_MAPPING_ATTRIBUTE_NAME_NOT_SET = 0,
    TRANSIENT_MAPPING_ATTRIBUTE_NAME = 8,
  }
}

export class AddSAMLProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddSAMLProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddSAMLProviderResponse;

  getId(): string;
  setId(value: string): AddSAMLProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddSAMLProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddSAMLProviderResponse): AddSAMLProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddSAMLProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddSAMLProviderResponse;
  static deserializeBinaryFromReader(message: AddSAMLProviderResponse, reader: jspb.BinaryReader): AddSAMLProviderResponse;
}

export namespace AddSAMLProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateSAMLProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateSAMLProviderRequest;

  getName(): string;
  setName(value: string): UpdateSAMLProviderRequest;

  getMetadataXml(): Uint8Array | string;
  getMetadataXml_asU8(): Uint8Array;
  getMetadataXml_asB64(): string;
  setMetadataXml(value: Uint8Array | string): UpdateSAMLProviderRequest;

  getMetadataUrl(): string;
  setMetadataUrl(value: string): UpdateSAMLProviderRequest;

  getBinding(): zitadel_idp_pb.SAMLBinding;
  setBinding(value: zitadel_idp_pb.SAMLBinding): UpdateSAMLProviderRequest;

  getWithSignedRequest(): boolean;
  setWithSignedRequest(value: boolean): UpdateSAMLProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateSAMLProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateSAMLProviderRequest;

  getNameIdFormat(): zitadel_idp_pb.SAMLNameIDFormat;
  setNameIdFormat(value: zitadel_idp_pb.SAMLNameIDFormat): UpdateSAMLProviderRequest;
  hasNameIdFormat(): boolean;
  clearNameIdFormat(): UpdateSAMLProviderRequest;

  getTransientMappingAttributeName(): string;
  setTransientMappingAttributeName(value: string): UpdateSAMLProviderRequest;
  hasTransientMappingAttributeName(): boolean;
  clearTransientMappingAttributeName(): UpdateSAMLProviderRequest;

  getMetadataCase(): UpdateSAMLProviderRequest.MetadataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLProviderRequest): UpdateSAMLProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLProviderRequest;
  static deserializeBinaryFromReader(message: UpdateSAMLProviderRequest, reader: jspb.BinaryReader): UpdateSAMLProviderRequest;
}

export namespace UpdateSAMLProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    metadataXml: Uint8Array | string,
    metadataUrl: string,
    binding: zitadel_idp_pb.SAMLBinding,
    withSignedRequest: boolean,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
    nameIdFormat?: zitadel_idp_pb.SAMLNameIDFormat,
    transientMappingAttributeName?: string,
  }

  export enum MetadataCase { 
    METADATA_NOT_SET = 0,
    METADATA_XML = 3,
    METADATA_URL = 4,
  }

  export enum NameIdFormatCase { 
    _NAME_ID_FORMAT_NOT_SET = 0,
    NAME_ID_FORMAT = 8,
  }

  export enum TransientMappingAttributeNameCase { 
    _TRANSIENT_MAPPING_ATTRIBUTE_NAME_NOT_SET = 0,
    TRANSIENT_MAPPING_ATTRIBUTE_NAME = 9,
  }
}

export class UpdateSAMLProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateSAMLProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateSAMLProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateSAMLProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateSAMLProviderResponse): UpdateSAMLProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateSAMLProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateSAMLProviderResponse;
  static deserializeBinaryFromReader(message: UpdateSAMLProviderResponse, reader: jspb.BinaryReader): UpdateSAMLProviderResponse;
}

export namespace UpdateSAMLProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RegenerateSAMLProviderCertificateRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RegenerateSAMLProviderCertificateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateSAMLProviderCertificateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateSAMLProviderCertificateRequest): RegenerateSAMLProviderCertificateRequest.AsObject;
  static serializeBinaryToWriter(message: RegenerateSAMLProviderCertificateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateSAMLProviderCertificateRequest;
  static deserializeBinaryFromReader(message: RegenerateSAMLProviderCertificateRequest, reader: jspb.BinaryReader): RegenerateSAMLProviderCertificateRequest;
}

export namespace RegenerateSAMLProviderCertificateRequest {
  export type AsObject = {
    id: string,
  }
}

export class RegenerateSAMLProviderCertificateResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RegenerateSAMLProviderCertificateResponse;
  hasDetails(): boolean;
  clearDetails(): RegenerateSAMLProviderCertificateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegenerateSAMLProviderCertificateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegenerateSAMLProviderCertificateResponse): RegenerateSAMLProviderCertificateResponse.AsObject;
  static serializeBinaryToWriter(message: RegenerateSAMLProviderCertificateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegenerateSAMLProviderCertificateResponse;
  static deserializeBinaryFromReader(message: RegenerateSAMLProviderCertificateResponse, reader: jspb.BinaryReader): RegenerateSAMLProviderCertificateResponse;
}

export namespace RegenerateSAMLProviderCertificateResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class AddAppleProviderRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddAppleProviderRequest;

  getClientId(): string;
  setClientId(value: string): AddAppleProviderRequest;

  getTeamId(): string;
  setTeamId(value: string): AddAppleProviderRequest;

  getKeyId(): string;
  setKeyId(value: string): AddAppleProviderRequest;

  getPrivateKey(): Uint8Array | string;
  getPrivateKey_asU8(): Uint8Array;
  getPrivateKey_asB64(): string;
  setPrivateKey(value: Uint8Array | string): AddAppleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): AddAppleProviderRequest;
  clearScopesList(): AddAppleProviderRequest;
  addScopes(value: string, index?: number): AddAppleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): AddAppleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): AddAppleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppleProviderRequest): AddAppleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: AddAppleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppleProviderRequest;
  static deserializeBinaryFromReader(message: AddAppleProviderRequest, reader: jspb.BinaryReader): AddAppleProviderRequest;
}

export namespace AddAppleProviderRequest {
  export type AsObject = {
    name: string,
    clientId: string,
    teamId: string,
    keyId: string,
    privateKey: Uint8Array | string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class AddAppleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddAppleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): AddAppleProviderResponse;

  getId(): string;
  setId(value: string): AddAppleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddAppleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddAppleProviderResponse): AddAppleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: AddAppleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddAppleProviderResponse;
  static deserializeBinaryFromReader(message: AddAppleProviderResponse, reader: jspb.BinaryReader): AddAppleProviderResponse;
}

export namespace AddAppleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class UpdateAppleProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateAppleProviderRequest;

  getName(): string;
  setName(value: string): UpdateAppleProviderRequest;

  getClientId(): string;
  setClientId(value: string): UpdateAppleProviderRequest;

  getTeamId(): string;
  setTeamId(value: string): UpdateAppleProviderRequest;

  getKeyId(): string;
  setKeyId(value: string): UpdateAppleProviderRequest;

  getPrivateKey(): Uint8Array | string;
  getPrivateKey_asU8(): Uint8Array;
  getPrivateKey_asB64(): string;
  setPrivateKey(value: Uint8Array | string): UpdateAppleProviderRequest;

  getScopesList(): Array<string>;
  setScopesList(value: Array<string>): UpdateAppleProviderRequest;
  clearScopesList(): UpdateAppleProviderRequest;
  addScopes(value: string, index?: number): UpdateAppleProviderRequest;

  getProviderOptions(): zitadel_idp_pb.Options | undefined;
  setProviderOptions(value?: zitadel_idp_pb.Options): UpdateAppleProviderRequest;
  hasProviderOptions(): boolean;
  clearProviderOptions(): UpdateAppleProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppleProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppleProviderRequest): UpdateAppleProviderRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateAppleProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppleProviderRequest;
  static deserializeBinaryFromReader(message: UpdateAppleProviderRequest, reader: jspb.BinaryReader): UpdateAppleProviderRequest;
}

export namespace UpdateAppleProviderRequest {
  export type AsObject = {
    id: string,
    name: string,
    clientId: string,
    teamId: string,
    keyId: string,
    privateKey: Uint8Array | string,
    scopesList: Array<string>,
    providerOptions?: zitadel_idp_pb.Options.AsObject,
  }
}

export class UpdateAppleProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateAppleProviderResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateAppleProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateAppleProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateAppleProviderResponse): UpdateAppleProviderResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateAppleProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateAppleProviderResponse;
  static deserializeBinaryFromReader(message: UpdateAppleProviderResponse, reader: jspb.BinaryReader): UpdateAppleProviderResponse;
}

export namespace UpdateAppleProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeleteProviderRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteProviderRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProviderRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProviderRequest): DeleteProviderRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteProviderRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProviderRequest;
  static deserializeBinaryFromReader(message: DeleteProviderRequest, reader: jspb.BinaryReader): DeleteProviderRequest;
}

export namespace DeleteProviderRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteProviderResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeleteProviderResponse;
  hasDetails(): boolean;
  clearDetails(): DeleteProviderResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProviderResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProviderResponse): DeleteProviderResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteProviderResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProviderResponse;
  static deserializeBinaryFromReader(message: DeleteProviderResponse, reader: jspb.BinaryReader): DeleteProviderResponse;
}

export namespace DeleteProviderResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListActionsRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListActionsRequest;
  hasQuery(): boolean;
  clearQuery(): ListActionsRequest;

  getSortingColumn(): zitadel_action_pb.ActionFieldName;
  setSortingColumn(value: zitadel_action_pb.ActionFieldName): ListActionsRequest;

  getQueriesList(): Array<ActionQuery>;
  setQueriesList(value: Array<ActionQuery>): ListActionsRequest;
  clearQueriesList(): ListActionsRequest;
  addQueries(value?: ActionQuery, index?: number): ActionQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListActionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListActionsRequest): ListActionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListActionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListActionsRequest;
  static deserializeBinaryFromReader(message: ListActionsRequest, reader: jspb.BinaryReader): ListActionsRequest;
}

export namespace ListActionsRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_action_pb.ActionFieldName,
    queriesList: Array<ActionQuery.AsObject>,
  }
}

export class ActionQuery extends jspb.Message {
  getActionIdQuery(): zitadel_action_pb.ActionIDQuery | undefined;
  setActionIdQuery(value?: zitadel_action_pb.ActionIDQuery): ActionQuery;
  hasActionIdQuery(): boolean;
  clearActionIdQuery(): ActionQuery;

  getActionNameQuery(): zitadel_action_pb.ActionNameQuery | undefined;
  setActionNameQuery(value?: zitadel_action_pb.ActionNameQuery): ActionQuery;
  hasActionNameQuery(): boolean;
  clearActionNameQuery(): ActionQuery;

  getActionStateQuery(): zitadel_action_pb.ActionStateQuery | undefined;
  setActionStateQuery(value?: zitadel_action_pb.ActionStateQuery): ActionQuery;
  hasActionStateQuery(): boolean;
  clearActionStateQuery(): ActionQuery;

  getQueryCase(): ActionQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ActionQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ActionQuery): ActionQuery.AsObject;
  static serializeBinaryToWriter(message: ActionQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ActionQuery;
  static deserializeBinaryFromReader(message: ActionQuery, reader: jspb.BinaryReader): ActionQuery;
}

export namespace ActionQuery {
  export type AsObject = {
    actionIdQuery?: zitadel_action_pb.ActionIDQuery.AsObject,
    actionNameQuery?: zitadel_action_pb.ActionNameQuery.AsObject,
    actionStateQuery?: zitadel_action_pb.ActionStateQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    ACTION_ID_QUERY = 1,
    ACTION_NAME_QUERY = 2,
    ACTION_STATE_QUERY = 3,
  }
}

export class ListActionsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListActionsResponse;
  hasDetails(): boolean;
  clearDetails(): ListActionsResponse;

  getSortingColumn(): zitadel_action_pb.ActionFieldName;
  setSortingColumn(value: zitadel_action_pb.ActionFieldName): ListActionsResponse;

  getResultList(): Array<zitadel_action_pb.Action>;
  setResultList(value: Array<zitadel_action_pb.Action>): ListActionsResponse;
  clearResultList(): ListActionsResponse;
  addResult(value?: zitadel_action_pb.Action, index?: number): zitadel_action_pb.Action;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListActionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListActionsResponse): ListActionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListActionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListActionsResponse;
  static deserializeBinaryFromReader(message: ListActionsResponse, reader: jspb.BinaryReader): ListActionsResponse;
}

export namespace ListActionsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_action_pb.ActionFieldName,
    resultList: Array<zitadel_action_pb.Action.AsObject>,
  }
}

export class CreateActionRequest extends jspb.Message {
  getName(): string;
  setName(value: string): CreateActionRequest;

  getScript(): string;
  setScript(value: string): CreateActionRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): CreateActionRequest;
  hasTimeout(): boolean;
  clearTimeout(): CreateActionRequest;

  getAllowedToFail(): boolean;
  setAllowedToFail(value: boolean): CreateActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateActionRequest): CreateActionRequest.AsObject;
  static serializeBinaryToWriter(message: CreateActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateActionRequest;
  static deserializeBinaryFromReader(message: CreateActionRequest, reader: jspb.BinaryReader): CreateActionRequest;
}

export namespace CreateActionRequest {
  export type AsObject = {
    name: string,
    script: string,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    allowedToFail: boolean,
  }
}

export class CreateActionResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): CreateActionResponse;
  hasDetails(): boolean;
  clearDetails(): CreateActionResponse;

  getId(): string;
  setId(value: string): CreateActionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateActionResponse): CreateActionResponse.AsObject;
  static serializeBinaryToWriter(message: CreateActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateActionResponse;
  static deserializeBinaryFromReader(message: CreateActionResponse, reader: jspb.BinaryReader): CreateActionResponse;
}

export namespace CreateActionResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    id: string,
  }
}

export class GetActionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): GetActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetActionRequest): GetActionRequest.AsObject;
  static serializeBinaryToWriter(message: GetActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetActionRequest;
  static deserializeBinaryFromReader(message: GetActionRequest, reader: jspb.BinaryReader): GetActionRequest;
}

export namespace GetActionRequest {
  export type AsObject = {
    id: string,
  }
}

export class GetActionResponse extends jspb.Message {
  getAction(): zitadel_action_pb.Action | undefined;
  setAction(value?: zitadel_action_pb.Action): GetActionResponse;
  hasAction(): boolean;
  clearAction(): GetActionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetActionResponse): GetActionResponse.AsObject;
  static serializeBinaryToWriter(message: GetActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetActionResponse;
  static deserializeBinaryFromReader(message: GetActionResponse, reader: jspb.BinaryReader): GetActionResponse;
}

export namespace GetActionResponse {
  export type AsObject = {
    action?: zitadel_action_pb.Action.AsObject,
  }
}

export class UpdateActionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): UpdateActionRequest;

  getName(): string;
  setName(value: string): UpdateActionRequest;

  getScript(): string;
  setScript(value: string): UpdateActionRequest;

  getTimeout(): google_protobuf_duration_pb.Duration | undefined;
  setTimeout(value?: google_protobuf_duration_pb.Duration): UpdateActionRequest;
  hasTimeout(): boolean;
  clearTimeout(): UpdateActionRequest;

  getAllowedToFail(): boolean;
  setAllowedToFail(value: boolean): UpdateActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateActionRequest): UpdateActionRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateActionRequest;
  static deserializeBinaryFromReader(message: UpdateActionRequest, reader: jspb.BinaryReader): UpdateActionRequest;
}

export namespace UpdateActionRequest {
  export type AsObject = {
    id: string,
    name: string,
    script: string,
    timeout?: google_protobuf_duration_pb.Duration.AsObject,
    allowedToFail: boolean,
  }
}

export class UpdateActionResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateActionResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateActionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateActionResponse): UpdateActionResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateActionResponse;
  static deserializeBinaryFromReader(message: UpdateActionResponse, reader: jspb.BinaryReader): UpdateActionResponse;
}

export namespace UpdateActionResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class DeleteActionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeleteActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteActionRequest): DeleteActionRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteActionRequest;
  static deserializeBinaryFromReader(message: DeleteActionRequest, reader: jspb.BinaryReader): DeleteActionRequest;
}

export namespace DeleteActionRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteActionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteActionResponse): DeleteActionResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteActionResponse;
  static deserializeBinaryFromReader(message: DeleteActionResponse, reader: jspb.BinaryReader): DeleteActionResponse;
}

export namespace DeleteActionResponse {
  export type AsObject = {
  }
}

export class ListFlowTypesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlowTypesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlowTypesRequest): ListFlowTypesRequest.AsObject;
  static serializeBinaryToWriter(message: ListFlowTypesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlowTypesRequest;
  static deserializeBinaryFromReader(message: ListFlowTypesRequest, reader: jspb.BinaryReader): ListFlowTypesRequest;
}

export namespace ListFlowTypesRequest {
  export type AsObject = {
  }
}

export class ListFlowTypesResponse extends jspb.Message {
  getResultList(): Array<zitadel_action_pb.FlowType>;
  setResultList(value: Array<zitadel_action_pb.FlowType>): ListFlowTypesResponse;
  clearResultList(): ListFlowTypesResponse;
  addResult(value?: zitadel_action_pb.FlowType, index?: number): zitadel_action_pb.FlowType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlowTypesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlowTypesResponse): ListFlowTypesResponse.AsObject;
  static serializeBinaryToWriter(message: ListFlowTypesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlowTypesResponse;
  static deserializeBinaryFromReader(message: ListFlowTypesResponse, reader: jspb.BinaryReader): ListFlowTypesResponse;
}

export namespace ListFlowTypesResponse {
  export type AsObject = {
    resultList: Array<zitadel_action_pb.FlowType.AsObject>,
  }
}

export class ListFlowTriggerTypesRequest extends jspb.Message {
  getType(): string;
  setType(value: string): ListFlowTriggerTypesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlowTriggerTypesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlowTriggerTypesRequest): ListFlowTriggerTypesRequest.AsObject;
  static serializeBinaryToWriter(message: ListFlowTriggerTypesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlowTriggerTypesRequest;
  static deserializeBinaryFromReader(message: ListFlowTriggerTypesRequest, reader: jspb.BinaryReader): ListFlowTriggerTypesRequest;
}

export namespace ListFlowTriggerTypesRequest {
  export type AsObject = {
    type: string,
  }
}

export class ListFlowTriggerTypesResponse extends jspb.Message {
  getResultList(): Array<zitadel_action_pb.TriggerType>;
  setResultList(value: Array<zitadel_action_pb.TriggerType>): ListFlowTriggerTypesResponse;
  clearResultList(): ListFlowTriggerTypesResponse;
  addResult(value?: zitadel_action_pb.TriggerType, index?: number): zitadel_action_pb.TriggerType;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFlowTriggerTypesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFlowTriggerTypesResponse): ListFlowTriggerTypesResponse.AsObject;
  static serializeBinaryToWriter(message: ListFlowTriggerTypesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFlowTriggerTypesResponse;
  static deserializeBinaryFromReader(message: ListFlowTriggerTypesResponse, reader: jspb.BinaryReader): ListFlowTriggerTypesResponse;
}

export namespace ListFlowTriggerTypesResponse {
  export type AsObject = {
    resultList: Array<zitadel_action_pb.TriggerType.AsObject>,
  }
}

export class DeactivateActionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DeactivateActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateActionRequest): DeactivateActionRequest.AsObject;
  static serializeBinaryToWriter(message: DeactivateActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateActionRequest;
  static deserializeBinaryFromReader(message: DeactivateActionRequest, reader: jspb.BinaryReader): DeactivateActionRequest;
}

export namespace DeactivateActionRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeactivateActionResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): DeactivateActionResponse;
  hasDetails(): boolean;
  clearDetails(): DeactivateActionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeactivateActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeactivateActionResponse): DeactivateActionResponse.AsObject;
  static serializeBinaryToWriter(message: DeactivateActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeactivateActionResponse;
  static deserializeBinaryFromReader(message: DeactivateActionResponse, reader: jspb.BinaryReader): DeactivateActionResponse;
}

export namespace DeactivateActionResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ReactivateActionRequest extends jspb.Message {
  getId(): string;
  setId(value: string): ReactivateActionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateActionRequest): ReactivateActionRequest.AsObject;
  static serializeBinaryToWriter(message: ReactivateActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateActionRequest;
  static deserializeBinaryFromReader(message: ReactivateActionRequest, reader: jspb.BinaryReader): ReactivateActionRequest;
}

export namespace ReactivateActionRequest {
  export type AsObject = {
    id: string,
  }
}

export class ReactivateActionResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ReactivateActionResponse;
  hasDetails(): boolean;
  clearDetails(): ReactivateActionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReactivateActionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ReactivateActionResponse): ReactivateActionResponse.AsObject;
  static serializeBinaryToWriter(message: ReactivateActionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReactivateActionResponse;
  static deserializeBinaryFromReader(message: ReactivateActionResponse, reader: jspb.BinaryReader): ReactivateActionResponse;
}

export namespace ReactivateActionResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class GetFlowRequest extends jspb.Message {
  getType(): string;
  setType(value: string): GetFlowRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFlowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetFlowRequest): GetFlowRequest.AsObject;
  static serializeBinaryToWriter(message: GetFlowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFlowRequest;
  static deserializeBinaryFromReader(message: GetFlowRequest, reader: jspb.BinaryReader): GetFlowRequest;
}

export namespace GetFlowRequest {
  export type AsObject = {
    type: string,
  }
}

export class GetFlowResponse extends jspb.Message {
  getFlow(): zitadel_action_pb.Flow | undefined;
  setFlow(value?: zitadel_action_pb.Flow): GetFlowResponse;
  hasFlow(): boolean;
  clearFlow(): GetFlowResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetFlowResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetFlowResponse): GetFlowResponse.AsObject;
  static serializeBinaryToWriter(message: GetFlowResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetFlowResponse;
  static deserializeBinaryFromReader(message: GetFlowResponse, reader: jspb.BinaryReader): GetFlowResponse;
}

export namespace GetFlowResponse {
  export type AsObject = {
    flow?: zitadel_action_pb.Flow.AsObject,
  }
}

export class ClearFlowRequest extends jspb.Message {
  getType(): string;
  setType(value: string): ClearFlowRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearFlowRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ClearFlowRequest): ClearFlowRequest.AsObject;
  static serializeBinaryToWriter(message: ClearFlowRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearFlowRequest;
  static deserializeBinaryFromReader(message: ClearFlowRequest, reader: jspb.BinaryReader): ClearFlowRequest;
}

export namespace ClearFlowRequest {
  export type AsObject = {
    type: string,
  }
}

export class ClearFlowResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ClearFlowResponse;
  hasDetails(): boolean;
  clearDetails(): ClearFlowResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearFlowResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ClearFlowResponse): ClearFlowResponse.AsObject;
  static serializeBinaryToWriter(message: ClearFlowResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearFlowResponse;
  static deserializeBinaryFromReader(message: ClearFlowResponse, reader: jspb.BinaryReader): ClearFlowResponse;
}

export namespace ClearFlowResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetTriggerActionsRequest extends jspb.Message {
  getFlowType(): string;
  setFlowType(value: string): SetTriggerActionsRequest;

  getTriggerType(): string;
  setTriggerType(value: string): SetTriggerActionsRequest;

  getActionIdsList(): Array<string>;
  setActionIdsList(value: Array<string>): SetTriggerActionsRequest;
  clearActionIdsList(): SetTriggerActionsRequest;
  addActionIds(value: string, index?: number): SetTriggerActionsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetTriggerActionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetTriggerActionsRequest): SetTriggerActionsRequest.AsObject;
  static serializeBinaryToWriter(message: SetTriggerActionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetTriggerActionsRequest;
  static deserializeBinaryFromReader(message: SetTriggerActionsRequest, reader: jspb.BinaryReader): SetTriggerActionsRequest;
}

export namespace SetTriggerActionsRequest {
  export type AsObject = {
    flowType: string,
    triggerType: string,
    actionIdsList: Array<string>,
  }
}

export class SetTriggerActionsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetTriggerActionsResponse;
  hasDetails(): boolean;
  clearDetails(): SetTriggerActionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetTriggerActionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetTriggerActionsResponse): SetTriggerActionsResponse.AsObject;
  static serializeBinaryToWriter(message: SetTriggerActionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetTriggerActionsResponse;
  static deserializeBinaryFromReader(message: SetTriggerActionsResponse, reader: jspb.BinaryReader): SetTriggerActionsResponse;
}

export namespace SetTriggerActionsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

