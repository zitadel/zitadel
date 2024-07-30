import * as jspb from 'google-protobuf'

import * as zitadel_object_pb from '../zitadel/object_pb'; // proto import: "zitadel/object.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class Project extends jspb.Message {
  getId(): string;
  setId(value: string): Project;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Project;
  hasDetails(): boolean;
  clearDetails(): Project;

  getName(): string;
  setName(value: string): Project;

  getState(): ProjectState;
  setState(value: ProjectState): Project;

  getProjectRoleAssertion(): boolean;
  setProjectRoleAssertion(value: boolean): Project;

  getProjectRoleCheck(): boolean;
  setProjectRoleCheck(value: boolean): Project;

  getHasProjectCheck(): boolean;
  setHasProjectCheck(value: boolean): Project;

  getPrivateLabelingSetting(): PrivateLabelingSetting;
  setPrivateLabelingSetting(value: PrivateLabelingSetting): Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Project.AsObject;
  static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
  static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Project;
  static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
}

export namespace Project {
  export type AsObject = {
    id: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    name: string,
    state: ProjectState,
    projectRoleAssertion: boolean,
    projectRoleCheck: boolean,
    hasProjectCheck: boolean,
    privateLabelingSetting: PrivateLabelingSetting,
  }
}

export class GrantedProject extends jspb.Message {
  getGrantId(): string;
  setGrantId(value: string): GrantedProject;

  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): GrantedProject;

  getGrantedOrgName(): string;
  setGrantedOrgName(value: string): GrantedProject;

  getGrantedRoleKeysList(): Array<string>;
  setGrantedRoleKeysList(value: Array<string>): GrantedProject;
  clearGrantedRoleKeysList(): GrantedProject;
  addGrantedRoleKeys(value: string, index?: number): GrantedProject;

  getState(): ProjectGrantState;
  setState(value: ProjectGrantState): GrantedProject;

  getProjectId(): string;
  setProjectId(value: string): GrantedProject;

  getProjectName(): string;
  setProjectName(value: string): GrantedProject;

  getProjectOwnerId(): string;
  setProjectOwnerId(value: string): GrantedProject;

  getProjectOwnerName(): string;
  setProjectOwnerName(value: string): GrantedProject;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): GrantedProject;
  hasDetails(): boolean;
  clearDetails(): GrantedProject;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantedProject.AsObject;
  static toObject(includeInstance: boolean, msg: GrantedProject): GrantedProject.AsObject;
  static serializeBinaryToWriter(message: GrantedProject, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantedProject;
  static deserializeBinaryFromReader(message: GrantedProject, reader: jspb.BinaryReader): GrantedProject;
}

export namespace GrantedProject {
  export type AsObject = {
    grantId: string,
    grantedOrgId: string,
    grantedOrgName: string,
    grantedRoleKeysList: Array<string>,
    state: ProjectGrantState,
    projectId: string,
    projectName: string,
    projectOwnerId: string,
    projectOwnerName: string,
    details?: zitadel_object_pb.ObjectDetails.AsObject,
  }
}

export class ProjectQuery extends jspb.Message {
  getNameQuery(): ProjectNameQuery | undefined;
  setNameQuery(value?: ProjectNameQuery): ProjectQuery;
  hasNameQuery(): boolean;
  clearNameQuery(): ProjectQuery;

  getProjectResourceOwnerQuery(): ProjectResourceOwnerQuery | undefined;
  setProjectResourceOwnerQuery(value?: ProjectResourceOwnerQuery): ProjectQuery;
  hasProjectResourceOwnerQuery(): boolean;
  clearProjectResourceOwnerQuery(): ProjectQuery;

  getQueryCase(): ProjectQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectQuery): ProjectQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectQuery;
  static deserializeBinaryFromReader(message: ProjectQuery, reader: jspb.BinaryReader): ProjectQuery;
}

export namespace ProjectQuery {
  export type AsObject = {
    nameQuery?: ProjectNameQuery.AsObject,
    projectResourceOwnerQuery?: ProjectResourceOwnerQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    NAME_QUERY = 1,
    PROJECT_RESOURCE_OWNER_QUERY = 2,
  }
}

export class ProjectNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): ProjectNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): ProjectNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectNameQuery): ProjectNameQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectNameQuery;
  static deserializeBinaryFromReader(message: ProjectNameQuery, reader: jspb.BinaryReader): ProjectNameQuery;
}

export namespace ProjectNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class ProjectResourceOwnerQuery extends jspb.Message {
  getResourceOwner(): string;
  setResourceOwner(value: string): ProjectResourceOwnerQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectResourceOwnerQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectResourceOwnerQuery): ProjectResourceOwnerQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectResourceOwnerQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectResourceOwnerQuery;
  static deserializeBinaryFromReader(message: ProjectResourceOwnerQuery, reader: jspb.BinaryReader): ProjectResourceOwnerQuery;
}

export namespace ProjectResourceOwnerQuery {
  export type AsObject = {
    resourceOwner: string,
  }
}

export class Role extends jspb.Message {
  getKey(): string;
  setKey(value: string): Role;

  getDetails(): zitadel_object_pb.ObjectDetails | undefined;
  setDetails(value?: zitadel_object_pb.ObjectDetails): Role;
  hasDetails(): boolean;
  clearDetails(): Role;

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
    details?: zitadel_object_pb.ObjectDetails.AsObject,
    displayName: string,
    group: string,
  }
}

export class RoleQuery extends jspb.Message {
  getKeyQuery(): RoleKeyQuery | undefined;
  setKeyQuery(value?: RoleKeyQuery): RoleQuery;
  hasKeyQuery(): boolean;
  clearKeyQuery(): RoleQuery;

  getDisplayNameQuery(): RoleDisplayNameQuery | undefined;
  setDisplayNameQuery(value?: RoleDisplayNameQuery): RoleQuery;
  hasDisplayNameQuery(): boolean;
  clearDisplayNameQuery(): RoleQuery;

  getQueryCase(): RoleQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoleQuery.AsObject;
  static toObject(includeInstance: boolean, msg: RoleQuery): RoleQuery.AsObject;
  static serializeBinaryToWriter(message: RoleQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoleQuery;
  static deserializeBinaryFromReader(message: RoleQuery, reader: jspb.BinaryReader): RoleQuery;
}

export namespace RoleQuery {
  export type AsObject = {
    keyQuery?: RoleKeyQuery.AsObject,
    displayNameQuery?: RoleDisplayNameQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    KEY_QUERY = 1,
    DISPLAY_NAME_QUERY = 2,
  }
}

export class RoleKeyQuery extends jspb.Message {
  getKey(): string;
  setKey(value: string): RoleKeyQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): RoleKeyQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoleKeyQuery.AsObject;
  static toObject(includeInstance: boolean, msg: RoleKeyQuery): RoleKeyQuery.AsObject;
  static serializeBinaryToWriter(message: RoleKeyQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoleKeyQuery;
  static deserializeBinaryFromReader(message: RoleKeyQuery, reader: jspb.BinaryReader): RoleKeyQuery;
}

export namespace RoleKeyQuery {
  export type AsObject = {
    key: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class RoleDisplayNameQuery extends jspb.Message {
  getDisplayName(): string;
  setDisplayName(value: string): RoleDisplayNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): RoleDisplayNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RoleDisplayNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: RoleDisplayNameQuery): RoleDisplayNameQuery.AsObject;
  static serializeBinaryToWriter(message: RoleDisplayNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RoleDisplayNameQuery;
  static deserializeBinaryFromReader(message: RoleDisplayNameQuery, reader: jspb.BinaryReader): RoleDisplayNameQuery;
}

export namespace RoleDisplayNameQuery {
  export type AsObject = {
    displayName: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class ProjectGrantQuery extends jspb.Message {
  getProjectNameQuery(): GrantProjectNameQuery | undefined;
  setProjectNameQuery(value?: GrantProjectNameQuery): ProjectGrantQuery;
  hasProjectNameQuery(): boolean;
  clearProjectNameQuery(): ProjectGrantQuery;

  getRoleKeyQuery(): GrantRoleKeyQuery | undefined;
  setRoleKeyQuery(value?: GrantRoleKeyQuery): ProjectGrantQuery;
  hasRoleKeyQuery(): boolean;
  clearRoleKeyQuery(): ProjectGrantQuery;

  getQueryCase(): ProjectGrantQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectGrantQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectGrantQuery): ProjectGrantQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectGrantQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectGrantQuery;
  static deserializeBinaryFromReader(message: ProjectGrantQuery, reader: jspb.BinaryReader): ProjectGrantQuery;
}

export namespace ProjectGrantQuery {
  export type AsObject = {
    projectNameQuery?: GrantProjectNameQuery.AsObject,
    roleKeyQuery?: GrantRoleKeyQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    PROJECT_NAME_QUERY = 1,
    ROLE_KEY_QUERY = 2,
  }
}

export class AllProjectGrantQuery extends jspb.Message {
  getProjectNameQuery(): GrantProjectNameQuery | undefined;
  setProjectNameQuery(value?: GrantProjectNameQuery): AllProjectGrantQuery;
  hasProjectNameQuery(): boolean;
  clearProjectNameQuery(): AllProjectGrantQuery;

  getRoleKeyQuery(): GrantRoleKeyQuery | undefined;
  setRoleKeyQuery(value?: GrantRoleKeyQuery): AllProjectGrantQuery;
  hasRoleKeyQuery(): boolean;
  clearRoleKeyQuery(): AllProjectGrantQuery;

  getProjectIdQuery(): ProjectIDQuery | undefined;
  setProjectIdQuery(value?: ProjectIDQuery): AllProjectGrantQuery;
  hasProjectIdQuery(): boolean;
  clearProjectIdQuery(): AllProjectGrantQuery;

  getGrantedOrgIdQuery(): GrantedOrgIDQuery | undefined;
  setGrantedOrgIdQuery(value?: GrantedOrgIDQuery): AllProjectGrantQuery;
  hasGrantedOrgIdQuery(): boolean;
  clearGrantedOrgIdQuery(): AllProjectGrantQuery;

  getQueryCase(): AllProjectGrantQuery.QueryCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AllProjectGrantQuery.AsObject;
  static toObject(includeInstance: boolean, msg: AllProjectGrantQuery): AllProjectGrantQuery.AsObject;
  static serializeBinaryToWriter(message: AllProjectGrantQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AllProjectGrantQuery;
  static deserializeBinaryFromReader(message: AllProjectGrantQuery, reader: jspb.BinaryReader): AllProjectGrantQuery;
}

export namespace AllProjectGrantQuery {
  export type AsObject = {
    projectNameQuery?: GrantProjectNameQuery.AsObject,
    roleKeyQuery?: GrantRoleKeyQuery.AsObject,
    projectIdQuery?: ProjectIDQuery.AsObject,
    grantedOrgIdQuery?: GrantedOrgIDQuery.AsObject,
  }

  export enum QueryCase { 
    QUERY_NOT_SET = 0,
    PROJECT_NAME_QUERY = 1,
    ROLE_KEY_QUERY = 2,
    PROJECT_ID_QUERY = 3,
    GRANTED_ORG_ID_QUERY = 4,
  }
}

export class GrantProjectNameQuery extends jspb.Message {
  getName(): string;
  setName(value: string): GrantProjectNameQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): GrantProjectNameQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantProjectNameQuery.AsObject;
  static toObject(includeInstance: boolean, msg: GrantProjectNameQuery): GrantProjectNameQuery.AsObject;
  static serializeBinaryToWriter(message: GrantProjectNameQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantProjectNameQuery;
  static deserializeBinaryFromReader(message: GrantProjectNameQuery, reader: jspb.BinaryReader): GrantProjectNameQuery;
}

export namespace GrantProjectNameQuery {
  export type AsObject = {
    name: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class GrantRoleKeyQuery extends jspb.Message {
  getRoleKey(): string;
  setRoleKey(value: string): GrantRoleKeyQuery;

  getMethod(): zitadel_object_pb.TextQueryMethod;
  setMethod(value: zitadel_object_pb.TextQueryMethod): GrantRoleKeyQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantRoleKeyQuery.AsObject;
  static toObject(includeInstance: boolean, msg: GrantRoleKeyQuery): GrantRoleKeyQuery.AsObject;
  static serializeBinaryToWriter(message: GrantRoleKeyQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantRoleKeyQuery;
  static deserializeBinaryFromReader(message: GrantRoleKeyQuery, reader: jspb.BinaryReader): GrantRoleKeyQuery;
}

export namespace GrantRoleKeyQuery {
  export type AsObject = {
    roleKey: string,
    method: zitadel_object_pb.TextQueryMethod,
  }
}

export class ProjectIDQuery extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ProjectIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectIDQuery): ProjectIDQuery.AsObject;
  static serializeBinaryToWriter(message: ProjectIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectIDQuery;
  static deserializeBinaryFromReader(message: ProjectIDQuery, reader: jspb.BinaryReader): ProjectIDQuery;
}

export namespace ProjectIDQuery {
  export type AsObject = {
    projectId: string,
  }
}

export class GrantedOrgIDQuery extends jspb.Message {
  getGrantedOrgId(): string;
  setGrantedOrgId(value: string): GrantedOrgIDQuery;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GrantedOrgIDQuery.AsObject;
  static toObject(includeInstance: boolean, msg: GrantedOrgIDQuery): GrantedOrgIDQuery.AsObject;
  static serializeBinaryToWriter(message: GrantedOrgIDQuery, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GrantedOrgIDQuery;
  static deserializeBinaryFromReader(message: GrantedOrgIDQuery, reader: jspb.BinaryReader): GrantedOrgIDQuery;
}

export namespace GrantedOrgIDQuery {
  export type AsObject = {
    grantedOrgId: string,
  }
}

export enum ProjectState { 
  PROJECT_STATE_UNSPECIFIED = 0,
  PROJECT_STATE_ACTIVE = 1,
  PROJECT_STATE_INACTIVE = 2,
}
export enum PrivateLabelingSetting { 
  PRIVATE_LABELING_SETTING_UNSPECIFIED = 0,
  PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY = 1,
  PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY = 2,
}
export enum ProjectGrantState { 
  PROJECT_GRANT_STATE_UNSPECIFIED = 0,
  PROJECT_GRANT_STATE_ACTIVE = 1,
  PROJECT_GRANT_STATE_INACTIVE = 2,
}
