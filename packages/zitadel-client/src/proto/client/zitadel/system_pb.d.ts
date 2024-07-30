import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as zitadel_options_pb from '../zitadel/options_pb'; // proto import: "zitadel/options.proto"
import * as zitadel_instance_pb from '../zitadel/instance_pb'; // proto import: "zitadel/instance.proto"
import * as zitadel_member_pb from '../zitadel/member_pb'; // proto import: "zitadel/member.proto"
import * as zitadel_quota_pb from '../zitadel/quota_pb'; // proto import: "zitadel/quota.proto"
import * as zitadel_auth_n_key_pb from '../zitadel/auth_n_key_pb'; // proto import: "zitadel/auth_n_key.proto"
import * as zitadel_feature_pb from '../zitadel/feature_pb'; // proto import: "zitadel/feature.proto"
import * as google_api_annotations_pb from '../google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
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

export class ListInstancesRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListInstancesRequest;
  hasQuery(): boolean;
  clearQuery(): ListInstancesRequest;

  getSortingColumn(): zitadel_instance_pb.FieldName;
  setSortingColumn(value: zitadel_instance_pb.FieldName): ListInstancesRequest;

  getQueriesList(): Array<zitadel_instance_pb.Query>;
  setQueriesList(value: Array<zitadel_instance_pb.Query>): ListInstancesRequest;
  clearQueriesList(): ListInstancesRequest;
  addQueries(value?: zitadel_instance_pb.Query, index?: number): zitadel_instance_pb.Query;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstancesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstancesRequest): ListInstancesRequest.AsObject;
  static serializeBinaryToWriter(message: ListInstancesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstancesRequest;
  static deserializeBinaryFromReader(message: ListInstancesRequest, reader: jspb.BinaryReader): ListInstancesRequest;
}

export namespace ListInstancesRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_instance_pb.FieldName,
    queriesList: Array<zitadel_instance_pb.Query.AsObject>,
  }
}

export class ListInstancesResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListInstancesResponse;
  hasDetails(): boolean;
  clearDetails(): ListInstancesResponse;

  getSortingColumn(): zitadel_instance_pb.FieldName;
  setSortingColumn(value: zitadel_instance_pb.FieldName): ListInstancesResponse;

  getResultList(): Array<zitadel_instance_pb.Instance>;
  setResultList(value: Array<zitadel_instance_pb.Instance>): ListInstancesResponse;
  clearResultList(): ListInstancesResponse;
  addResult(value?: zitadel_instance_pb.Instance, index?: number): zitadel_instance_pb.Instance;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListInstancesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListInstancesResponse): ListInstancesResponse.AsObject;
  static serializeBinaryToWriter(message: ListInstancesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListInstancesResponse;
  static deserializeBinaryFromReader(message: ListInstancesResponse, reader: jspb.BinaryReader): ListInstancesResponse;
}

export namespace ListInstancesResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_instance_pb.FieldName,
    resultList: Array<zitadel_instance_pb.Instance.AsObject>,
  }
}

export class GetInstanceRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): GetInstanceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetInstanceRequest): GetInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: GetInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInstanceRequest;
  static deserializeBinaryFromReader(message: GetInstanceRequest, reader: jspb.BinaryReader): GetInstanceRequest;
}

export namespace GetInstanceRequest {
  export type AsObject = {
    instanceId: string,
  }
}

export class GetInstanceResponse extends jspb.Message {
  getInstance(): zitadel_instance_pb.InstanceDetail | undefined;
  setInstance(value?: zitadel_instance_pb.InstanceDetail): GetInstanceResponse;
  hasInstance(): boolean;
  clearInstance(): GetInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetInstanceResponse): GetInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: GetInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInstanceResponse;
  static deserializeBinaryFromReader(message: GetInstanceResponse, reader: jspb.BinaryReader): GetInstanceResponse;
}

export namespace GetInstanceResponse {
  export type AsObject = {
    instance?: zitadel_instance_pb.InstanceDetail.AsObject,
  }
}

export class AddInstanceRequest extends jspb.Message {
  getInstanceName(): string;
  setInstanceName(value: string): AddInstanceRequest;

  getFirstOrgName(): string;
  setFirstOrgName(value: string): AddInstanceRequest;

  getCustomDomain(): string;
  setCustomDomain(value: string): AddInstanceRequest;

  getOwnerUserName(): string;
  setOwnerUserName(value: string): AddInstanceRequest;

  getOwnerEmail(): AddInstanceRequest.Email | undefined;
  setOwnerEmail(value?: AddInstanceRequest.Email): AddInstanceRequest;
  hasOwnerEmail(): boolean;
  clearOwnerEmail(): AddInstanceRequest;

  getOwnerProfile(): AddInstanceRequest.Profile | undefined;
  setOwnerProfile(value?: AddInstanceRequest.Profile): AddInstanceRequest;
  hasOwnerProfile(): boolean;
  clearOwnerProfile(): AddInstanceRequest;

  getOwnerPassword(): AddInstanceRequest.Password | undefined;
  setOwnerPassword(value?: AddInstanceRequest.Password): AddInstanceRequest;
  hasOwnerPassword(): boolean;
  clearOwnerPassword(): AddInstanceRequest;

  getDefaultLanguage(): string;
  setDefaultLanguage(value: string): AddInstanceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddInstanceRequest): AddInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: AddInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddInstanceRequest;
  static deserializeBinaryFromReader(message: AddInstanceRequest, reader: jspb.BinaryReader): AddInstanceRequest;
}

export namespace AddInstanceRequest {
  export type AsObject = {
    instanceName: string,
    firstOrgName: string,
    customDomain: string,
    ownerUserName: string,
    ownerEmail?: AddInstanceRequest.Email.AsObject,
    ownerProfile?: AddInstanceRequest.Profile.AsObject,
    ownerPassword?: AddInstanceRequest.Password.AsObject,
    defaultLanguage: string,
  }

  export class Profile extends jspb.Message {
    getFirstName(): string;
    setFirstName(value: string): Profile;

    getLastName(): string;
    setLastName(value: string): Profile;

    getPreferredLanguage(): string;
    setPreferredLanguage(value: string): Profile;

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
      preferredLanguage: string,
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


  export class Password extends jspb.Message {
    getPassword(): string;
    setPassword(value: string): Password;

    getPasswordChangeRequired(): boolean;
    setPasswordChangeRequired(value: boolean): Password;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Password.AsObject;
    static toObject(includeInstance: boolean, msg: Password): Password.AsObject;
    static serializeBinaryToWriter(message: Password, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Password;
    static deserializeBinaryFromReader(message: Password, reader: jspb.BinaryReader): Password;
  }

  export namespace Password {
    export type AsObject = {
      password: string,
      passwordChangeRequired: boolean,
    }
  }

}

export class AddInstanceResponse extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): AddInstanceResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddInstanceResponse;
  hasDetails(): boolean;
  clearDetails(): AddInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddInstanceResponse): AddInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: AddInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddInstanceResponse;
  static deserializeBinaryFromReader(message: AddInstanceResponse, reader: jspb.BinaryReader): AddInstanceResponse;
}

export namespace AddInstanceResponse {
  export type AsObject = {
    instanceId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class CreateInstanceRequest extends jspb.Message {
  getInstanceName(): string;
  setInstanceName(value: string): CreateInstanceRequest;

  getFirstOrgName(): string;
  setFirstOrgName(value: string): CreateInstanceRequest;

  getCustomDomain(): string;
  setCustomDomain(value: string): CreateInstanceRequest;

  getHuman(): CreateInstanceRequest.Human | undefined;
  setHuman(value?: CreateInstanceRequest.Human): CreateInstanceRequest;
  hasHuman(): boolean;
  clearHuman(): CreateInstanceRequest;

  getMachine(): CreateInstanceRequest.Machine | undefined;
  setMachine(value?: CreateInstanceRequest.Machine): CreateInstanceRequest;
  hasMachine(): boolean;
  clearMachine(): CreateInstanceRequest;

  getDefaultLanguage(): string;
  setDefaultLanguage(value: string): CreateInstanceRequest;

  getOwnerCase(): CreateInstanceRequest.OwnerCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CreateInstanceRequest): CreateInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: CreateInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateInstanceRequest;
  static deserializeBinaryFromReader(message: CreateInstanceRequest, reader: jspb.BinaryReader): CreateInstanceRequest;
}

export namespace CreateInstanceRequest {
  export type AsObject = {
    instanceName: string,
    firstOrgName: string,
    customDomain: string,
    human?: CreateInstanceRequest.Human.AsObject,
    machine?: CreateInstanceRequest.Machine.AsObject,
    defaultLanguage: string,
  }

  export class Profile extends jspb.Message {
    getFirstName(): string;
    setFirstName(value: string): Profile;

    getLastName(): string;
    setLastName(value: string): Profile;

    getPreferredLanguage(): string;
    setPreferredLanguage(value: string): Profile;

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
      preferredLanguage: string,
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


  export class Password extends jspb.Message {
    getPassword(): string;
    setPassword(value: string): Password;

    getPasswordChangeRequired(): boolean;
    setPasswordChangeRequired(value: boolean): Password;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Password.AsObject;
    static toObject(includeInstance: boolean, msg: Password): Password.AsObject;
    static serializeBinaryToWriter(message: Password, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Password;
    static deserializeBinaryFromReader(message: Password, reader: jspb.BinaryReader): Password;
  }

  export namespace Password {
    export type AsObject = {
      password: string,
      passwordChangeRequired: boolean,
    }
  }


  export class Human extends jspb.Message {
    getUserName(): string;
    setUserName(value: string): Human;

    getEmail(): CreateInstanceRequest.Email | undefined;
    setEmail(value?: CreateInstanceRequest.Email): Human;
    hasEmail(): boolean;
    clearEmail(): Human;

    getProfile(): CreateInstanceRequest.Profile | undefined;
    setProfile(value?: CreateInstanceRequest.Profile): Human;
    hasProfile(): boolean;
    clearProfile(): Human;

    getPassword(): CreateInstanceRequest.Password | undefined;
    setPassword(value?: CreateInstanceRequest.Password): Human;
    hasPassword(): boolean;
    clearPassword(): Human;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Human.AsObject;
    static toObject(includeInstance: boolean, msg: Human): Human.AsObject;
    static serializeBinaryToWriter(message: Human, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Human;
    static deserializeBinaryFromReader(message: Human, reader: jspb.BinaryReader): Human;
  }

  export namespace Human {
    export type AsObject = {
      userName: string,
      email?: CreateInstanceRequest.Email.AsObject,
      profile?: CreateInstanceRequest.Profile.AsObject,
      password?: CreateInstanceRequest.Password.AsObject,
    }
  }


  export class PersonalAccessToken extends jspb.Message {
    getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): PersonalAccessToken;
    hasExpirationDate(): boolean;
    clearExpirationDate(): PersonalAccessToken;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PersonalAccessToken.AsObject;
    static toObject(includeInstance: boolean, msg: PersonalAccessToken): PersonalAccessToken.AsObject;
    static serializeBinaryToWriter(message: PersonalAccessToken, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PersonalAccessToken;
    static deserializeBinaryFromReader(message: PersonalAccessToken, reader: jspb.BinaryReader): PersonalAccessToken;
  }

  export namespace PersonalAccessToken {
    export type AsObject = {
      expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }


  export class MachineKey extends jspb.Message {
    getType(): zitadel_auth_n_key_pb.KeyType;
    setType(value: zitadel_auth_n_key_pb.KeyType): MachineKey;

    getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): MachineKey;
    hasExpirationDate(): boolean;
    clearExpirationDate(): MachineKey;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): MachineKey.AsObject;
    static toObject(includeInstance: boolean, msg: MachineKey): MachineKey.AsObject;
    static serializeBinaryToWriter(message: MachineKey, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): MachineKey;
    static deserializeBinaryFromReader(message: MachineKey, reader: jspb.BinaryReader): MachineKey;
  }

  export namespace MachineKey {
    export type AsObject = {
      type: zitadel_auth_n_key_pb.KeyType,
      expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
  }


  export class Machine extends jspb.Message {
    getUserName(): string;
    setUserName(value: string): Machine;

    getName(): string;
    setName(value: string): Machine;

    getPersonalAccessToken(): CreateInstanceRequest.PersonalAccessToken | undefined;
    setPersonalAccessToken(value?: CreateInstanceRequest.PersonalAccessToken): Machine;
    hasPersonalAccessToken(): boolean;
    clearPersonalAccessToken(): Machine;

    getMachineKey(): CreateInstanceRequest.MachineKey | undefined;
    setMachineKey(value?: CreateInstanceRequest.MachineKey): Machine;
    hasMachineKey(): boolean;
    clearMachineKey(): Machine;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Machine.AsObject;
    static toObject(includeInstance: boolean, msg: Machine): Machine.AsObject;
    static serializeBinaryToWriter(message: Machine, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Machine;
    static deserializeBinaryFromReader(message: Machine, reader: jspb.BinaryReader): Machine;
  }

  export namespace Machine {
    export type AsObject = {
      userName: string,
      name: string,
      personalAccessToken?: CreateInstanceRequest.PersonalAccessToken.AsObject,
      machineKey?: CreateInstanceRequest.MachineKey.AsObject,
    }
  }


  export enum OwnerCase { 
    OWNER_NOT_SET = 0,
    HUMAN = 4,
    MACHINE = 5,
  }
}

export class CreateInstanceResponse extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): CreateInstanceResponse;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): CreateInstanceResponse;
  hasDetails(): boolean;
  clearDetails(): CreateInstanceResponse;

  getPat(): string;
  setPat(value: string): CreateInstanceResponse;

  getMachineKey(): Uint8Array | string;
  getMachineKey_asU8(): Uint8Array;
  getMachineKey_asB64(): string;
  setMachineKey(value: Uint8Array | string): CreateInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CreateInstanceResponse): CreateInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: CreateInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateInstanceResponse;
  static deserializeBinaryFromReader(message: CreateInstanceResponse, reader: jspb.BinaryReader): CreateInstanceResponse;
}

export namespace CreateInstanceResponse {
  export type AsObject = {
    instanceId: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    pat: string,
    machineKey: Uint8Array | string,
  }
}

export class UpdateInstanceRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): UpdateInstanceRequest;

  getInstanceName(): string;
  setInstanceName(value: string): UpdateInstanceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateInstanceRequest): UpdateInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateInstanceRequest;
  static deserializeBinaryFromReader(message: UpdateInstanceRequest, reader: jspb.BinaryReader): UpdateInstanceRequest;
}

export namespace UpdateInstanceRequest {
  export type AsObject = {
    instanceId: string,
    instanceName: string,
  }
}

export class UpdateInstanceResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): UpdateInstanceResponse;
  hasDetails(): boolean;
  clearDetails(): UpdateInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateInstanceResponse): UpdateInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateInstanceResponse;
  static deserializeBinaryFromReader(message: UpdateInstanceResponse, reader: jspb.BinaryReader): UpdateInstanceResponse;
}

export namespace UpdateInstanceResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveInstanceRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): RemoveInstanceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveInstanceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveInstanceRequest): RemoveInstanceRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveInstanceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveInstanceRequest;
  static deserializeBinaryFromReader(message: RemoveInstanceRequest, reader: jspb.BinaryReader): RemoveInstanceRequest;
}

export namespace RemoveInstanceRequest {
  export type AsObject = {
    instanceId: string,
  }
}

export class RemoveInstanceResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveInstanceResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveInstanceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveInstanceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveInstanceResponse): RemoveInstanceResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveInstanceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveInstanceResponse;
  static deserializeBinaryFromReader(message: RemoveInstanceResponse, reader: jspb.BinaryReader): RemoveInstanceResponse;
}

export namespace RemoveInstanceResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListIAMMembersRequest extends jspb.Message {
  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListIAMMembersRequest;
  hasQuery(): boolean;
  clearQuery(): ListIAMMembersRequest;

  getInstanceId(): string;
  setInstanceId(value: string): ListIAMMembersRequest;

  getQueriesList(): Array<zitadel_member_pb.SearchQuery>;
  setQueriesList(value: Array<zitadel_member_pb.SearchQuery>): ListIAMMembersRequest;
  clearQueriesList(): ListIAMMembersRequest;
  addQueries(value?: zitadel_member_pb.SearchQuery, index?: number): zitadel_member_pb.SearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMembersRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMembersRequest): ListIAMMembersRequest.AsObject;
  static serializeBinaryToWriter(message: ListIAMMembersRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMembersRequest;
  static deserializeBinaryFromReader(message: ListIAMMembersRequest, reader: jspb.BinaryReader): ListIAMMembersRequest;
}

export namespace ListIAMMembersRequest {
  export type AsObject = {
    query?: zitadel_object_pb.ListQuery.AsObject,
    instanceId: string,
    queriesList: Array<zitadel_member_pb.SearchQuery.AsObject>,
  }
}

export class ListIAMMembersResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListIAMMembersResponse;
  hasDetails(): boolean;
  clearDetails(): ListIAMMembersResponse;

  getResultList(): Array<zitadel_member_pb.Member>;
  setResultList(value: Array<zitadel_member_pb.Member>): ListIAMMembersResponse;
  clearResultList(): ListIAMMembersResponse;
  addResult(value?: zitadel_member_pb.Member, index?: number): zitadel_member_pb.Member;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListIAMMembersResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListIAMMembersResponse): ListIAMMembersResponse.AsObject;
  static serializeBinaryToWriter(message: ListIAMMembersResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListIAMMembersResponse;
  static deserializeBinaryFromReader(message: ListIAMMembersResponse, reader: jspb.BinaryReader): ListIAMMembersResponse;
}

export namespace ListIAMMembersResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    resultList: Array<zitadel_member_pb.Member.AsObject>,
  }
}

export class GetUsageRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): GetUsageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetUsageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetUsageRequest): GetUsageRequest.AsObject;
  static serializeBinaryToWriter(message: GetUsageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetUsageRequest;
  static deserializeBinaryFromReader(message: GetUsageRequest, reader: jspb.BinaryReader): GetUsageRequest;
}

export namespace GetUsageRequest {
  export type AsObject = {
    instanceId: string,
  }
}

export class AddQuotaRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): AddQuotaRequest;

  getUnit(): zitadel_quota_pb.Unit;
  setUnit(value: zitadel_quota_pb.Unit): AddQuotaRequest;

  getFrom(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFrom(value?: google_protobuf_timestamp_pb.Timestamp): AddQuotaRequest;
  hasFrom(): boolean;
  clearFrom(): AddQuotaRequest;

  getResetInterval(): google_protobuf_duration_pb.Duration | undefined;
  setResetInterval(value?: google_protobuf_duration_pb.Duration): AddQuotaRequest;
  hasResetInterval(): boolean;
  clearResetInterval(): AddQuotaRequest;

  getAmount(): number;
  setAmount(value: number): AddQuotaRequest;

  getLimit(): boolean;
  setLimit(value: boolean): AddQuotaRequest;

  getNotificationsList(): Array<zitadel_quota_pb.Notification>;
  setNotificationsList(value: Array<zitadel_quota_pb.Notification>): AddQuotaRequest;
  clearNotificationsList(): AddQuotaRequest;
  addNotifications(value?: zitadel_quota_pb.Notification, index?: number): zitadel_quota_pb.Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddQuotaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddQuotaRequest): AddQuotaRequest.AsObject;
  static serializeBinaryToWriter(message: AddQuotaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddQuotaRequest;
  static deserializeBinaryFromReader(message: AddQuotaRequest, reader: jspb.BinaryReader): AddQuotaRequest;
}

export namespace AddQuotaRequest {
  export type AsObject = {
    instanceId: string,
    unit: zitadel_quota_pb.Unit,
    from?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resetInterval?: google_protobuf_duration_pb.Duration.AsObject,
    amount: number,
    limit: boolean,
    notificationsList: Array<zitadel_quota_pb.Notification.AsObject>,
  }
}

export class AddQuotaResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddQuotaResponse;
  hasDetails(): boolean;
  clearDetails(): AddQuotaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddQuotaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddQuotaResponse): AddQuotaResponse.AsObject;
  static serializeBinaryToWriter(message: AddQuotaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddQuotaResponse;
  static deserializeBinaryFromReader(message: AddQuotaResponse, reader: jspb.BinaryReader): AddQuotaResponse;
}

export namespace AddQuotaResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetQuotaRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): SetQuotaRequest;

  getUnit(): zitadel_quota_pb.Unit;
  setUnit(value: zitadel_quota_pb.Unit): SetQuotaRequest;

  getFrom(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setFrom(value?: google_protobuf_timestamp_pb.Timestamp): SetQuotaRequest;
  hasFrom(): boolean;
  clearFrom(): SetQuotaRequest;

  getResetInterval(): google_protobuf_duration_pb.Duration | undefined;
  setResetInterval(value?: google_protobuf_duration_pb.Duration): SetQuotaRequest;
  hasResetInterval(): boolean;
  clearResetInterval(): SetQuotaRequest;

  getAmount(): number;
  setAmount(value: number): SetQuotaRequest;

  getLimit(): boolean;
  setLimit(value: boolean): SetQuotaRequest;

  getNotificationsList(): Array<zitadel_quota_pb.Notification>;
  setNotificationsList(value: Array<zitadel_quota_pb.Notification>): SetQuotaRequest;
  clearNotificationsList(): SetQuotaRequest;
  addNotifications(value?: zitadel_quota_pb.Notification, index?: number): zitadel_quota_pb.Notification;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetQuotaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetQuotaRequest): SetQuotaRequest.AsObject;
  static serializeBinaryToWriter(message: SetQuotaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetQuotaRequest;
  static deserializeBinaryFromReader(message: SetQuotaRequest, reader: jspb.BinaryReader): SetQuotaRequest;
}

export namespace SetQuotaRequest {
  export type AsObject = {
    instanceId: string,
    unit: zitadel_quota_pb.Unit,
    from?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    resetInterval?: google_protobuf_duration_pb.Duration.AsObject,
    amount: number,
    limit: boolean,
    notificationsList: Array<zitadel_quota_pb.Notification.AsObject>,
  }
}

export class SetQuotaResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetQuotaResponse;
  hasDetails(): boolean;
  clearDetails(): SetQuotaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetQuotaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetQuotaResponse): SetQuotaResponse.AsObject;
  static serializeBinaryToWriter(message: SetQuotaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetQuotaResponse;
  static deserializeBinaryFromReader(message: SetQuotaResponse, reader: jspb.BinaryReader): SetQuotaResponse;
}

export namespace SetQuotaResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveQuotaRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): RemoveQuotaRequest;

  getUnit(): zitadel_quota_pb.Unit;
  setUnit(value: zitadel_quota_pb.Unit): RemoveQuotaRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveQuotaRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveQuotaRequest): RemoveQuotaRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveQuotaRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveQuotaRequest;
  static deserializeBinaryFromReader(message: RemoveQuotaRequest, reader: jspb.BinaryReader): RemoveQuotaRequest;
}

export namespace RemoveQuotaRequest {
  export type AsObject = {
    instanceId: string,
    unit: zitadel_quota_pb.Unit,
  }
}

export class RemoveQuotaResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveQuotaResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveQuotaResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveQuotaResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveQuotaResponse): RemoveQuotaResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveQuotaResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveQuotaResponse;
  static deserializeBinaryFromReader(message: RemoveQuotaResponse, reader: jspb.BinaryReader): RemoveQuotaResponse;
}

export namespace RemoveQuotaResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetLimitsRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): SetLimitsRequest;

  getAuditLogRetention(): google_protobuf_duration_pb.Duration | undefined;
  setAuditLogRetention(value?: google_protobuf_duration_pb.Duration): SetLimitsRequest;
  hasAuditLogRetention(): boolean;
  clearAuditLogRetention(): SetLimitsRequest;

  getBlock(): boolean;
  setBlock(value: boolean): SetLimitsRequest;
  hasBlock(): boolean;
  clearBlock(): SetLimitsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetLimitsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetLimitsRequest): SetLimitsRequest.AsObject;
  static serializeBinaryToWriter(message: SetLimitsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetLimitsRequest;
  static deserializeBinaryFromReader(message: SetLimitsRequest, reader: jspb.BinaryReader): SetLimitsRequest;
}

export namespace SetLimitsRequest {
  export type AsObject = {
    instanceId: string,
    auditLogRetention?: google_protobuf_duration_pb.Duration.AsObject,
    block?: boolean,
  }

  export enum BlockCase { 
    _BLOCK_NOT_SET = 0,
    BLOCK = 3,
  }
}

export class SetLimitsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetLimitsResponse;
  hasDetails(): boolean;
  clearDetails(): SetLimitsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetLimitsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetLimitsResponse): SetLimitsResponse.AsObject;
  static serializeBinaryToWriter(message: SetLimitsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetLimitsResponse;
  static deserializeBinaryFromReader(message: SetLimitsResponse, reader: jspb.BinaryReader): SetLimitsResponse;
}

export namespace SetLimitsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class BulkSetLimitsRequest extends jspb.Message {
  getLimitsList(): Array<SetLimitsRequest>;
  setLimitsList(value: Array<SetLimitsRequest>): BulkSetLimitsRequest;
  clearLimitsList(): BulkSetLimitsRequest;
  addLimits(value?: SetLimitsRequest, index?: number): SetLimitsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetLimitsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetLimitsRequest): BulkSetLimitsRequest.AsObject;
  static serializeBinaryToWriter(message: BulkSetLimitsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetLimitsRequest;
  static deserializeBinaryFromReader(message: BulkSetLimitsRequest, reader: jspb.BinaryReader): BulkSetLimitsRequest;
}

export namespace BulkSetLimitsRequest {
  export type AsObject = {
    limitsList: Array<SetLimitsRequest.AsObject>,
  }
}

export class BulkSetLimitsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): BulkSetLimitsResponse;
  hasDetails(): boolean;
  clearDetails(): BulkSetLimitsResponse;

  getTargetDetailsList(): Array<zitadel_object_pb.ObjectDetails>;
  setTargetDetailsList(value: Array<zitadel_object_pb.ObjectDetails>): BulkSetLimitsResponse;
  clearTargetDetailsList(): BulkSetLimitsResponse;
  addTargetDetails(value?: zitadel_object_pb.ObjectDetails, index?: number): zitadel_object_pb.ObjectDetails;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BulkSetLimitsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BulkSetLimitsResponse): BulkSetLimitsResponse.AsObject;
  static serializeBinaryToWriter(message: BulkSetLimitsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BulkSetLimitsResponse;
  static deserializeBinaryFromReader(message: BulkSetLimitsResponse, reader: jspb.BinaryReader): BulkSetLimitsResponse;
}

export namespace BulkSetLimitsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    targetDetailsList: Array<zitadel_object_pb.ObjectDetails.AsObject>,
  }
}

export class ResetLimitsRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): ResetLimitsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLimitsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLimitsRequest): ResetLimitsRequest.AsObject;
  static serializeBinaryToWriter(message: ResetLimitsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLimitsRequest;
  static deserializeBinaryFromReader(message: ResetLimitsRequest, reader: jspb.BinaryReader): ResetLimitsRequest;
}

export namespace ResetLimitsRequest {
  export type AsObject = {
    instanceId: string,
  }
}

export class ResetLimitsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ResetLimitsResponse;
  hasDetails(): boolean;
  clearDetails(): ResetLimitsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResetLimitsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ResetLimitsResponse): ResetLimitsResponse.AsObject;
  static serializeBinaryToWriter(message: ResetLimitsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResetLimitsResponse;
  static deserializeBinaryFromReader(message: ResetLimitsResponse, reader: jspb.BinaryReader): ResetLimitsResponse;
}

export namespace ResetLimitsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ExistsDomainRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): ExistsDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExistsDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ExistsDomainRequest): ExistsDomainRequest.AsObject;
  static serializeBinaryToWriter(message: ExistsDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExistsDomainRequest;
  static deserializeBinaryFromReader(message: ExistsDomainRequest, reader: jspb.BinaryReader): ExistsDomainRequest;
}

export namespace ExistsDomainRequest {
  export type AsObject = {
    domain: string,
  }
}

export class ExistsDomainResponse extends jspb.Message {
  getExists(): boolean;
  setExists(value: boolean): ExistsDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExistsDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ExistsDomainResponse): ExistsDomainResponse.AsObject;
  static serializeBinaryToWriter(message: ExistsDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExistsDomainResponse;
  static deserializeBinaryFromReader(message: ExistsDomainResponse, reader: jspb.BinaryReader): ExistsDomainResponse;
}

export namespace ExistsDomainResponse {
  export type AsObject = {
    exists: boolean,
  }
}

export class ListDomainsRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): ListDomainsRequest;

  getQuery(): zitadel_object_pb.ListQuery | undefined;
  setQuery(value?: zitadel_object_pb.ListQuery): ListDomainsRequest;
  hasQuery(): boolean;
  clearQuery(): ListDomainsRequest;

  getSortingColumn(): zitadel_instance_pb.DomainFieldName;
  setSortingColumn(value: zitadel_instance_pb.DomainFieldName): ListDomainsRequest;

  getQueriesList(): Array<zitadel_instance_pb.DomainSearchQuery>;
  setQueriesList(value: Array<zitadel_instance_pb.DomainSearchQuery>): ListDomainsRequest;
  clearQueriesList(): ListDomainsRequest;
  addQueries(value?: zitadel_instance_pb.DomainSearchQuery, index?: number): zitadel_instance_pb.DomainSearchQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDomainsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListDomainsRequest): ListDomainsRequest.AsObject;
  static serializeBinaryToWriter(message: ListDomainsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDomainsRequest;
  static deserializeBinaryFromReader(message: ListDomainsRequest, reader: jspb.BinaryReader): ListDomainsRequest;
}

export namespace ListDomainsRequest {
  export type AsObject = {
    instanceId: string,
    query?: zitadel_object_pb.ListQuery.AsObject,
    sortingColumn: zitadel_instance_pb.DomainFieldName,
    queriesList: Array<zitadel_instance_pb.DomainSearchQuery.AsObject>,
  }
}

export class ListDomainsResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ListDetails | undefined;
  setDetails(value?: zitadel_object_pb.ListDetails): ListDomainsResponse;
  hasDetails(): boolean;
  clearDetails(): ListDomainsResponse;

  getSortingColumn(): zitadel_instance_pb.DomainFieldName;
  setSortingColumn(value: zitadel_instance_pb.DomainFieldName): ListDomainsResponse;

  getResultList(): Array<zitadel_instance_pb.Domain>;
  setResultList(value: Array<zitadel_instance_pb.Domain>): ListDomainsResponse;
  clearResultList(): ListDomainsResponse;
  addResult(value?: zitadel_instance_pb.Domain, index?: number): zitadel_instance_pb.Domain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDomainsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDomainsResponse): ListDomainsResponse.AsObject;
  static serializeBinaryToWriter(message: ListDomainsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDomainsResponse;
  static deserializeBinaryFromReader(message: ListDomainsResponse, reader: jspb.BinaryReader): ListDomainsResponse;
}

export namespace ListDomainsResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ListDetails.AsObject,
    sortingColumn: zitadel_instance_pb.DomainFieldName,
    resultList: Array<zitadel_instance_pb.Domain.AsObject>,
  }
}

export class AddDomainRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): AddDomainRequest;

  getDomain(): string;
  setDomain(value: string): AddDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddDomainRequest): AddDomainRequest.AsObject;
  static serializeBinaryToWriter(message: AddDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDomainRequest;
  static deserializeBinaryFromReader(message: AddDomainRequest, reader: jspb.BinaryReader): AddDomainRequest;
}

export namespace AddDomainRequest {
  export type AsObject = {
    instanceId: string,
    domain: string,
  }
}

export class AddDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): AddDomainResponse;
  hasDetails(): boolean;
  clearDetails(): AddDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddDomainResponse): AddDomainResponse.AsObject;
  static serializeBinaryToWriter(message: AddDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddDomainResponse;
  static deserializeBinaryFromReader(message: AddDomainResponse, reader: jspb.BinaryReader): AddDomainResponse;
}

export namespace AddDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class RemoveDomainRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): RemoveDomainRequest;

  getDomain(): string;
  setDomain(value: string): RemoveDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDomainRequest): RemoveDomainRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDomainRequest;
  static deserializeBinaryFromReader(message: RemoveDomainRequest, reader: jspb.BinaryReader): RemoveDomainRequest;
}

export namespace RemoveDomainRequest {
  export type AsObject = {
    instanceId: string,
    domain: string,
  }
}

export class RemoveDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): RemoveDomainResponse;
  hasDetails(): boolean;
  clearDetails(): RemoveDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveDomainResponse): RemoveDomainResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveDomainResponse;
  static deserializeBinaryFromReader(message: RemoveDomainResponse, reader: jspb.BinaryReader): RemoveDomainResponse;
}

export namespace RemoveDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class SetPrimaryDomainRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): SetPrimaryDomainRequest;

  getDomain(): string;
  setDomain(value: string): SetPrimaryDomainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPrimaryDomainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetPrimaryDomainRequest): SetPrimaryDomainRequest.AsObject;
  static serializeBinaryToWriter(message: SetPrimaryDomainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPrimaryDomainRequest;
  static deserializeBinaryFromReader(message: SetPrimaryDomainRequest, reader: jspb.BinaryReader): SetPrimaryDomainRequest;
}

export namespace SetPrimaryDomainRequest {
  export type AsObject = {
    instanceId: string,
    domain: string,
  }
}

export class SetPrimaryDomainResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetPrimaryDomainResponse;
  hasDetails(): boolean;
  clearDetails(): SetPrimaryDomainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetPrimaryDomainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetPrimaryDomainResponse): SetPrimaryDomainResponse.AsObject;
  static serializeBinaryToWriter(message: SetPrimaryDomainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetPrimaryDomainResponse;
  static deserializeBinaryFromReader(message: SetPrimaryDomainResponse, reader: jspb.BinaryReader): SetPrimaryDomainResponse;
}

export namespace SetPrimaryDomainResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ChangeSubscriptionRequest extends jspb.Message {
  getDomain(): string;
  setDomain(value: string): ChangeSubscriptionRequest;

  getSubscriptionName(): string;
  setSubscriptionName(value: string): ChangeSubscriptionRequest;

  getRequestLimit(): number;
  setRequestLimit(value: number): ChangeSubscriptionRequest;

  getActionMinsLimit(): number;
  setActionMinsLimit(value: number): ChangeSubscriptionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeSubscriptionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeSubscriptionRequest): ChangeSubscriptionRequest.AsObject;
  static serializeBinaryToWriter(message: ChangeSubscriptionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeSubscriptionRequest;
  static deserializeBinaryFromReader(message: ChangeSubscriptionRequest, reader: jspb.BinaryReader): ChangeSubscriptionRequest;
}

export namespace ChangeSubscriptionRequest {
  export type AsObject = {
    domain: string,
    subscriptionName: string,
    requestLimit: number,
    actionMinsLimit: number,
  }
}

export class ChangeSubscriptionResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): ChangeSubscriptionResponse;
  hasDetails(): boolean;
  clearDetails(): ChangeSubscriptionResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChangeSubscriptionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ChangeSubscriptionResponse): ChangeSubscriptionResponse.AsObject;
  static serializeBinaryToWriter(message: ChangeSubscriptionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChangeSubscriptionResponse;
  static deserializeBinaryFromReader(message: ChangeSubscriptionResponse, reader: jspb.BinaryReader): ChangeSubscriptionResponse;
}

export namespace ChangeSubscriptionResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ListViewsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListViewsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListViewsRequest): ListViewsRequest.AsObject;
  static serializeBinaryToWriter(message: ListViewsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListViewsRequest;
  static deserializeBinaryFromReader(message: ListViewsRequest, reader: jspb.BinaryReader): ListViewsRequest;
}

export namespace ListViewsRequest {
  export type AsObject = {
  }
}

export class ListViewsResponse extends jspb.Message {
  getResultList(): Array<View>;
  setResultList(value: Array<View>): ListViewsResponse;
  clearResultList(): ListViewsResponse;
  addResult(value?: View, index?: number): View;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListViewsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListViewsResponse): ListViewsResponse.AsObject;
  static serializeBinaryToWriter(message: ListViewsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListViewsResponse;
  static deserializeBinaryFromReader(message: ListViewsResponse, reader: jspb.BinaryReader): ListViewsResponse;
}

export namespace ListViewsResponse {
  export type AsObject = {
    resultList: Array<View.AsObject>,
  }
}

export class ClearViewRequest extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): ClearViewRequest;

  getViewName(): string;
  setViewName(value: string): ClearViewRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearViewRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ClearViewRequest): ClearViewRequest.AsObject;
  static serializeBinaryToWriter(message: ClearViewRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearViewRequest;
  static deserializeBinaryFromReader(message: ClearViewRequest, reader: jspb.BinaryReader): ClearViewRequest;
}

export namespace ClearViewRequest {
  export type AsObject = {
    database: string,
    viewName: string,
  }
}

export class ClearViewResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ClearViewResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ClearViewResponse): ClearViewResponse.AsObject;
  static serializeBinaryToWriter(message: ClearViewResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ClearViewResponse;
  static deserializeBinaryFromReader(message: ClearViewResponse, reader: jspb.BinaryReader): ClearViewResponse;
}

export namespace ClearViewResponse {
  export type AsObject = {
  }
}

export class ListFailedEventsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFailedEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListFailedEventsRequest): ListFailedEventsRequest.AsObject;
  static serializeBinaryToWriter(message: ListFailedEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFailedEventsRequest;
  static deserializeBinaryFromReader(message: ListFailedEventsRequest, reader: jspb.BinaryReader): ListFailedEventsRequest;
}

export namespace ListFailedEventsRequest {
  export type AsObject = {
  }
}

export class ListFailedEventsResponse extends jspb.Message {
  getResultList(): Array<FailedEvent>;
  setResultList(value: Array<FailedEvent>): ListFailedEventsResponse;
  clearResultList(): ListFailedEventsResponse;
  addResult(value?: FailedEvent, index?: number): FailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListFailedEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListFailedEventsResponse): ListFailedEventsResponse.AsObject;
  static serializeBinaryToWriter(message: ListFailedEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListFailedEventsResponse;
  static deserializeBinaryFromReader(message: ListFailedEventsResponse, reader: jspb.BinaryReader): ListFailedEventsResponse;
}

export namespace ListFailedEventsResponse {
  export type AsObject = {
    resultList: Array<FailedEvent.AsObject>,
  }
}

export class RemoveFailedEventRequest extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): RemoveFailedEventRequest;

  getViewName(): string;
  setViewName(value: string): RemoveFailedEventRequest;

  getFailedSequence(): number;
  setFailedSequence(value: number): RemoveFailedEventRequest;

  getInstanceId(): string;
  setInstanceId(value: string): RemoveFailedEventRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveFailedEventRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveFailedEventRequest): RemoveFailedEventRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveFailedEventRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveFailedEventRequest;
  static deserializeBinaryFromReader(message: RemoveFailedEventRequest, reader: jspb.BinaryReader): RemoveFailedEventRequest;
}

export namespace RemoveFailedEventRequest {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
    instanceId: string,
  }
}

export class RemoveFailedEventResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveFailedEventResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveFailedEventResponse): RemoveFailedEventResponse.AsObject;
  static serializeBinaryToWriter(message: RemoveFailedEventResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveFailedEventResponse;
  static deserializeBinaryFromReader(message: RemoveFailedEventResponse, reader: jspb.BinaryReader): RemoveFailedEventResponse;
}

export namespace RemoveFailedEventResponse {
  export type AsObject = {
  }
}

export class View extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): View;

  getViewName(): string;
  setViewName(value: string): View;

  getProcessedSequence(): number;
  setProcessedSequence(value: number): View;

  getEventTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setEventTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): View;
  hasEventTimestamp(): boolean;
  clearEventTimestamp(): View;

  getLastSuccessfulSpoolerRun(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastSuccessfulSpoolerRun(value?: google_protobuf_timestamp_pb.Timestamp): View;
  hasLastSuccessfulSpoolerRun(): boolean;
  clearLastSuccessfulSpoolerRun(): View;

  getInstance(): string;
  setInstance(value: string): View;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): View.AsObject;
  static toObject(includeInstance: boolean, msg: View): View.AsObject;
  static serializeBinaryToWriter(message: View, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): View;
  static deserializeBinaryFromReader(message: View, reader: jspb.BinaryReader): View;
}

export namespace View {
  export type AsObject = {
    database: string,
    viewName: string,
    processedSequence: number,
    eventTimestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    lastSuccessfulSpoolerRun?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    instance: string,
  }
}

export class FailedEvent extends jspb.Message {
  getDatabase(): string;
  setDatabase(value: string): FailedEvent;

  getViewName(): string;
  setViewName(value: string): FailedEvent;

  getFailedSequence(): number;
  setFailedSequence(value: number): FailedEvent;

  getFailureCount(): number;
  setFailureCount(value: number): FailedEvent;

  getErrorMessage(): string;
  setErrorMessage(value: string): FailedEvent;

  getLastFailed(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setLastFailed(value?: google_protobuf_timestamp_pb.Timestamp): FailedEvent;
  hasLastFailed(): boolean;
  clearLastFailed(): FailedEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FailedEvent.AsObject;
  static toObject(includeInstance: boolean, msg: FailedEvent): FailedEvent.AsObject;
  static serializeBinaryToWriter(message: FailedEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FailedEvent;
  static deserializeBinaryFromReader(message: FailedEvent, reader: jspb.BinaryReader): FailedEvent;
}

export namespace FailedEvent {
  export type AsObject = {
    database: string,
    viewName: string,
    failedSequence: number,
    failureCount: number,
    errorMessage: string,
    lastFailed?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class SetInstanceFeatureRequest extends jspb.Message {
  getInstanceId(): string;
  setInstanceId(value: string): SetInstanceFeatureRequest;

  getFeatureId(): zitadel_feature_pb.InstanceFeature;
  setFeatureId(value: zitadel_feature_pb.InstanceFeature): SetInstanceFeatureRequest;

  getBool(): boolean;
  setBool(value: boolean): SetInstanceFeatureRequest;

  getValueCase(): SetInstanceFeatureRequest.ValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetInstanceFeatureRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SetInstanceFeatureRequest): SetInstanceFeatureRequest.AsObject;
  static serializeBinaryToWriter(message: SetInstanceFeatureRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetInstanceFeatureRequest;
  static deserializeBinaryFromReader(message: SetInstanceFeatureRequest, reader: jspb.BinaryReader): SetInstanceFeatureRequest;
}

export namespace SetInstanceFeatureRequest {
  export type AsObject = {
    instanceId: string,
    featureId: zitadel_feature_pb.InstanceFeature,
    bool: boolean,
  }

  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    BOOL = 3,
  }
}

export class SetInstanceFeatureResponse extends jspb.Message {
  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): SetInstanceFeatureResponse;
  hasDetails(): boolean;
  clearDetails(): SetInstanceFeatureResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SetInstanceFeatureResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SetInstanceFeatureResponse): SetInstanceFeatureResponse.AsObject;
  static serializeBinaryToWriter(message: SetInstanceFeatureResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SetInstanceFeatureResponse;
  static deserializeBinaryFromReader(message: SetInstanceFeatureResponse, reader: jspb.BinaryReader): SetInstanceFeatureResponse;
}

export namespace SetInstanceFeatureResponse {
  export type AsObject = {
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

