/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { ObjectDetails, TextQueryMethod, textQueryMethodFromJSON, textQueryMethodToJSON } from "./object";

export const protobufPackage = "zitadel.project.v1";

export enum ProjectState {
  PROJECT_STATE_UNSPECIFIED = 0,
  PROJECT_STATE_ACTIVE = 1,
  PROJECT_STATE_INACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function projectStateFromJSON(object: any): ProjectState {
  switch (object) {
    case 0:
    case "PROJECT_STATE_UNSPECIFIED":
      return ProjectState.PROJECT_STATE_UNSPECIFIED;
    case 1:
    case "PROJECT_STATE_ACTIVE":
      return ProjectState.PROJECT_STATE_ACTIVE;
    case 2:
    case "PROJECT_STATE_INACTIVE":
      return ProjectState.PROJECT_STATE_INACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProjectState.UNRECOGNIZED;
  }
}

export function projectStateToJSON(object: ProjectState): string {
  switch (object) {
    case ProjectState.PROJECT_STATE_UNSPECIFIED:
      return "PROJECT_STATE_UNSPECIFIED";
    case ProjectState.PROJECT_STATE_ACTIVE:
      return "PROJECT_STATE_ACTIVE";
    case ProjectState.PROJECT_STATE_INACTIVE:
      return "PROJECT_STATE_INACTIVE";
    case ProjectState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum PrivateLabelingSetting {
  PRIVATE_LABELING_SETTING_UNSPECIFIED = 0,
  PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY = 1,
  PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY = 2,
  UNRECOGNIZED = -1,
}

export function privateLabelingSettingFromJSON(object: any): PrivateLabelingSetting {
  switch (object) {
    case 0:
    case "PRIVATE_LABELING_SETTING_UNSPECIFIED":
      return PrivateLabelingSetting.PRIVATE_LABELING_SETTING_UNSPECIFIED;
    case 1:
    case "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY":
      return PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY;
    case 2:
    case "PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY":
      return PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PrivateLabelingSetting.UNRECOGNIZED;
  }
}

export function privateLabelingSettingToJSON(object: PrivateLabelingSetting): string {
  switch (object) {
    case PrivateLabelingSetting.PRIVATE_LABELING_SETTING_UNSPECIFIED:
      return "PRIVATE_LABELING_SETTING_UNSPECIFIED";
    case PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY:
      return "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY";
    case PrivateLabelingSetting.PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY:
      return "PRIVATE_LABELING_SETTING_ALLOW_LOGIN_USER_RESOURCE_OWNER_POLICY";
    case PrivateLabelingSetting.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ProjectGrantState {
  PROJECT_GRANT_STATE_UNSPECIFIED = 0,
  PROJECT_GRANT_STATE_ACTIVE = 1,
  PROJECT_GRANT_STATE_INACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function projectGrantStateFromJSON(object: any): ProjectGrantState {
  switch (object) {
    case 0:
    case "PROJECT_GRANT_STATE_UNSPECIFIED":
      return ProjectGrantState.PROJECT_GRANT_STATE_UNSPECIFIED;
    case 1:
    case "PROJECT_GRANT_STATE_ACTIVE":
      return ProjectGrantState.PROJECT_GRANT_STATE_ACTIVE;
    case 2:
    case "PROJECT_GRANT_STATE_INACTIVE":
      return ProjectGrantState.PROJECT_GRANT_STATE_INACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProjectGrantState.UNRECOGNIZED;
  }
}

export function projectGrantStateToJSON(object: ProjectGrantState): string {
  switch (object) {
    case ProjectGrantState.PROJECT_GRANT_STATE_UNSPECIFIED:
      return "PROJECT_GRANT_STATE_UNSPECIFIED";
    case ProjectGrantState.PROJECT_GRANT_STATE_ACTIVE:
      return "PROJECT_GRANT_STATE_ACTIVE";
    case ProjectGrantState.PROJECT_GRANT_STATE_INACTIVE:
      return "PROJECT_GRANT_STATE_INACTIVE";
    case ProjectGrantState.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Project {
  id: string;
  details: ObjectDetails | undefined;
  name: string;
  state: ProjectState;
  /** describes if the roles of the user should be added to the token */
  projectRoleAssertion: boolean;
  /** ZITADEL checks if the user has at least one on this project */
  projectRoleCheck: boolean;
  /** ZITADEL checks if the org of the user has permission for this project */
  hasProjectCheck: boolean;
  /** Defines from where the private labeling should be triggered */
  privateLabelingSetting: PrivateLabelingSetting;
}

export interface GrantedProject {
  grantId: string;
  grantedOrgId: string;
  grantedOrgName: string;
  grantedRoleKeys: string[];
  state: ProjectGrantState;
  projectId: string;
  projectName: string;
  projectOwnerId: string;
  projectOwnerName: string;
  details: ObjectDetails | undefined;
}

export interface ProjectQuery {
  nameQuery?: ProjectNameQuery | undefined;
  projectResourceOwnerQuery?: ProjectResourceOwnerQuery | undefined;
}

export interface ProjectNameQuery {
  name: string;
  method: TextQueryMethod;
}

export interface ProjectResourceOwnerQuery {
  resourceOwner: string;
}

export interface Role {
  key: string;
  details: ObjectDetails | undefined;
  displayName: string;
  group: string;
}

export interface RoleQuery {
  keyQuery?: RoleKeyQuery | undefined;
  displayNameQuery?: RoleDisplayNameQuery | undefined;
}

export interface RoleKeyQuery {
  key: string;
  method: TextQueryMethod;
}

export interface RoleDisplayNameQuery {
  displayName: string;
  method: TextQueryMethod;
}

export interface ProjectGrantQuery {
  projectNameQuery?: GrantProjectNameQuery | undefined;
  roleKeyQuery?: GrantRoleKeyQuery | undefined;
}

export interface AllProjectGrantQuery {
  projectNameQuery?: GrantProjectNameQuery | undefined;
  roleKeyQuery?: GrantRoleKeyQuery | undefined;
  projectIdQuery?: ProjectIDQuery | undefined;
  grantedOrgIdQuery?: GrantedOrgIDQuery | undefined;
}

export interface GrantProjectNameQuery {
  name: string;
  method: TextQueryMethod;
}

export interface GrantRoleKeyQuery {
  roleKey: string;
  method: TextQueryMethod;
}

export interface ProjectIDQuery {
  projectId: string;
}

export interface GrantedOrgIDQuery {
  grantedOrgId: string;
}

function createBaseProject(): Project {
  return {
    id: "",
    details: undefined,
    name: "",
    state: 0,
    projectRoleAssertion: false,
    projectRoleCheck: false,
    hasProjectCheck: false,
    privateLabelingSetting: 0,
  };
}

export const Project = {
  encode(message: Project, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.state !== 0) {
      writer.uint32(32).int32(message.state);
    }
    if (message.projectRoleAssertion === true) {
      writer.uint32(40).bool(message.projectRoleAssertion);
    }
    if (message.projectRoleCheck === true) {
      writer.uint32(48).bool(message.projectRoleCheck);
    }
    if (message.hasProjectCheck === true) {
      writer.uint32(56).bool(message.hasProjectCheck);
    }
    if (message.privateLabelingSetting !== 0) {
      writer.uint32(64).int32(message.privateLabelingSetting);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Project {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.id = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.name = reader.string();
          continue;
        case 4:
          if (tag != 32) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.projectRoleAssertion = reader.bool();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.projectRoleCheck = reader.bool();
          continue;
        case 7:
          if (tag != 56) {
            break;
          }

          message.hasProjectCheck = reader.bool();
          continue;
        case 8:
          if (tag != 64) {
            break;
          }

          message.privateLabelingSetting = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Project {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      name: isSet(object.name) ? String(object.name) : "",
      state: isSet(object.state) ? projectStateFromJSON(object.state) : 0,
      projectRoleAssertion: isSet(object.projectRoleAssertion) ? Boolean(object.projectRoleAssertion) : false,
      projectRoleCheck: isSet(object.projectRoleCheck) ? Boolean(object.projectRoleCheck) : false,
      hasProjectCheck: isSet(object.hasProjectCheck) ? Boolean(object.hasProjectCheck) : false,
      privateLabelingSetting: isSet(object.privateLabelingSetting)
        ? privateLabelingSettingFromJSON(object.privateLabelingSetting)
        : 0,
    };
  },

  toJSON(message: Project): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.name !== undefined && (obj.name = message.name);
    message.state !== undefined && (obj.state = projectStateToJSON(message.state));
    message.projectRoleAssertion !== undefined && (obj.projectRoleAssertion = message.projectRoleAssertion);
    message.projectRoleCheck !== undefined && (obj.projectRoleCheck = message.projectRoleCheck);
    message.hasProjectCheck !== undefined && (obj.hasProjectCheck = message.hasProjectCheck);
    message.privateLabelingSetting !== undefined &&
      (obj.privateLabelingSetting = privateLabelingSettingToJSON(message.privateLabelingSetting));
    return obj;
  },

  create(base?: DeepPartial<Project>): Project {
    return Project.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Project>): Project {
    const message = createBaseProject();
    message.id = object.id ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.name = object.name ?? "";
    message.state = object.state ?? 0;
    message.projectRoleAssertion = object.projectRoleAssertion ?? false;
    message.projectRoleCheck = object.projectRoleCheck ?? false;
    message.hasProjectCheck = object.hasProjectCheck ?? false;
    message.privateLabelingSetting = object.privateLabelingSetting ?? 0;
    return message;
  },
};

function createBaseGrantedProject(): GrantedProject {
  return {
    grantId: "",
    grantedOrgId: "",
    grantedOrgName: "",
    grantedRoleKeys: [],
    state: 0,
    projectId: "",
    projectName: "",
    projectOwnerId: "",
    projectOwnerName: "",
    details: undefined,
  };
}

export const GrantedProject = {
  encode(message: GrantedProject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.grantId !== "") {
      writer.uint32(10).string(message.grantId);
    }
    if (message.grantedOrgId !== "") {
      writer.uint32(18).string(message.grantedOrgId);
    }
    if (message.grantedOrgName !== "") {
      writer.uint32(26).string(message.grantedOrgName);
    }
    for (const v of message.grantedRoleKeys) {
      writer.uint32(34).string(v!);
    }
    if (message.state !== 0) {
      writer.uint32(40).int32(message.state);
    }
    if (message.projectId !== "") {
      writer.uint32(50).string(message.projectId);
    }
    if (message.projectName !== "") {
      writer.uint32(58).string(message.projectName);
    }
    if (message.projectOwnerId !== "") {
      writer.uint32(66).string(message.projectOwnerId);
    }
    if (message.projectOwnerName !== "") {
      writer.uint32(74).string(message.projectOwnerName);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(82).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GrantedProject {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGrantedProject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.grantId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.grantedOrgId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.grantedOrgName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.grantedRoleKeys.push(reader.string());
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.state = reader.int32() as any;
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.projectId = reader.string();
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.projectName = reader.string();
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.projectOwnerId = reader.string();
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.projectOwnerName = reader.string();
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GrantedProject {
    return {
      grantId: isSet(object.grantId) ? String(object.grantId) : "",
      grantedOrgId: isSet(object.grantedOrgId) ? String(object.grantedOrgId) : "",
      grantedOrgName: isSet(object.grantedOrgName) ? String(object.grantedOrgName) : "",
      grantedRoleKeys: Array.isArray(object?.grantedRoleKeys) ? object.grantedRoleKeys.map((e: any) => String(e)) : [],
      state: isSet(object.state) ? projectGrantStateFromJSON(object.state) : 0,
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      projectName: isSet(object.projectName) ? String(object.projectName) : "",
      projectOwnerId: isSet(object.projectOwnerId) ? String(object.projectOwnerId) : "",
      projectOwnerName: isSet(object.projectOwnerName) ? String(object.projectOwnerName) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
    };
  },

  toJSON(message: GrantedProject): unknown {
    const obj: any = {};
    message.grantId !== undefined && (obj.grantId = message.grantId);
    message.grantedOrgId !== undefined && (obj.grantedOrgId = message.grantedOrgId);
    message.grantedOrgName !== undefined && (obj.grantedOrgName = message.grantedOrgName);
    if (message.grantedRoleKeys) {
      obj.grantedRoleKeys = message.grantedRoleKeys.map((e) => e);
    } else {
      obj.grantedRoleKeys = [];
    }
    message.state !== undefined && (obj.state = projectGrantStateToJSON(message.state));
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.projectName !== undefined && (obj.projectName = message.projectName);
    message.projectOwnerId !== undefined && (obj.projectOwnerId = message.projectOwnerId);
    message.projectOwnerName !== undefined && (obj.projectOwnerName = message.projectOwnerName);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    return obj;
  },

  create(base?: DeepPartial<GrantedProject>): GrantedProject {
    return GrantedProject.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GrantedProject>): GrantedProject {
    const message = createBaseGrantedProject();
    message.grantId = object.grantId ?? "";
    message.grantedOrgId = object.grantedOrgId ?? "";
    message.grantedOrgName = object.grantedOrgName ?? "";
    message.grantedRoleKeys = object.grantedRoleKeys?.map((e) => e) || [];
    message.state = object.state ?? 0;
    message.projectId = object.projectId ?? "";
    message.projectName = object.projectName ?? "";
    message.projectOwnerId = object.projectOwnerId ?? "";
    message.projectOwnerName = object.projectOwnerName ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    return message;
  },
};

function createBaseProjectQuery(): ProjectQuery {
  return { nameQuery: undefined, projectResourceOwnerQuery: undefined };
}

export const ProjectQuery = {
  encode(message: ProjectQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nameQuery !== undefined) {
      ProjectNameQuery.encode(message.nameQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.projectResourceOwnerQuery !== undefined) {
      ProjectResourceOwnerQuery.encode(message.projectResourceOwnerQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProjectQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProjectQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.nameQuery = ProjectNameQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.projectResourceOwnerQuery = ProjectResourceOwnerQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProjectQuery {
    return {
      nameQuery: isSet(object.nameQuery) ? ProjectNameQuery.fromJSON(object.nameQuery) : undefined,
      projectResourceOwnerQuery: isSet(object.projectResourceOwnerQuery)
        ? ProjectResourceOwnerQuery.fromJSON(object.projectResourceOwnerQuery)
        : undefined,
    };
  },

  toJSON(message: ProjectQuery): unknown {
    const obj: any = {};
    message.nameQuery !== undefined &&
      (obj.nameQuery = message.nameQuery ? ProjectNameQuery.toJSON(message.nameQuery) : undefined);
    message.projectResourceOwnerQuery !== undefined &&
      (obj.projectResourceOwnerQuery = message.projectResourceOwnerQuery
        ? ProjectResourceOwnerQuery.toJSON(message.projectResourceOwnerQuery)
        : undefined);
    return obj;
  },

  create(base?: DeepPartial<ProjectQuery>): ProjectQuery {
    return ProjectQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProjectQuery>): ProjectQuery {
    const message = createBaseProjectQuery();
    message.nameQuery = (object.nameQuery !== undefined && object.nameQuery !== null)
      ? ProjectNameQuery.fromPartial(object.nameQuery)
      : undefined;
    message.projectResourceOwnerQuery =
      (object.projectResourceOwnerQuery !== undefined && object.projectResourceOwnerQuery !== null)
        ? ProjectResourceOwnerQuery.fromPartial(object.projectResourceOwnerQuery)
        : undefined;
    return message;
  },
};

function createBaseProjectNameQuery(): ProjectNameQuery {
  return { name: "", method: 0 };
}

export const ProjectNameQuery = {
  encode(message: ProjectNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProjectNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProjectNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProjectNameQuery {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: ProjectNameQuery): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<ProjectNameQuery>): ProjectNameQuery {
    return ProjectNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProjectNameQuery>): ProjectNameQuery {
    const message = createBaseProjectNameQuery();
    message.name = object.name ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseProjectResourceOwnerQuery(): ProjectResourceOwnerQuery {
  return { resourceOwner: "" };
}

export const ProjectResourceOwnerQuery = {
  encode(message: ProjectResourceOwnerQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.resourceOwner !== "") {
      writer.uint32(10).string(message.resourceOwner);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProjectResourceOwnerQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProjectResourceOwnerQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.resourceOwner = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProjectResourceOwnerQuery {
    return { resourceOwner: isSet(object.resourceOwner) ? String(object.resourceOwner) : "" };
  },

  toJSON(message: ProjectResourceOwnerQuery): unknown {
    const obj: any = {};
    message.resourceOwner !== undefined && (obj.resourceOwner = message.resourceOwner);
    return obj;
  },

  create(base?: DeepPartial<ProjectResourceOwnerQuery>): ProjectResourceOwnerQuery {
    return ProjectResourceOwnerQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProjectResourceOwnerQuery>): ProjectResourceOwnerQuery {
    const message = createBaseProjectResourceOwnerQuery();
    message.resourceOwner = object.resourceOwner ?? "";
    return message;
  },
};

function createBaseRole(): Role {
  return { key: "", details: undefined, displayName: "", group: "" };
}

export const Role = {
  encode(message: Role, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.details !== undefined) {
      ObjectDetails.encode(message.details, writer.uint32(18).fork()).ldelim();
    }
    if (message.displayName !== "") {
      writer.uint32(26).string(message.displayName);
    }
    if (message.group !== "") {
      writer.uint32(34).string(message.group);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Role {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRole();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.details = ObjectDetails.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.displayName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.group = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Role {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      details: isSet(object.details) ? ObjectDetails.fromJSON(object.details) : undefined,
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      group: isSet(object.group) ? String(object.group) : "",
    };
  },

  toJSON(message: Role): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.details !== undefined &&
      (obj.details = message.details ? ObjectDetails.toJSON(message.details) : undefined);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.group !== undefined && (obj.group = message.group);
    return obj;
  },

  create(base?: DeepPartial<Role>): Role {
    return Role.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<Role>): Role {
    const message = createBaseRole();
    message.key = object.key ?? "";
    message.details = (object.details !== undefined && object.details !== null)
      ? ObjectDetails.fromPartial(object.details)
      : undefined;
    message.displayName = object.displayName ?? "";
    message.group = object.group ?? "";
    return message;
  },
};

function createBaseRoleQuery(): RoleQuery {
  return { keyQuery: undefined, displayNameQuery: undefined };
}

export const RoleQuery = {
  encode(message: RoleQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyQuery !== undefined) {
      RoleKeyQuery.encode(message.keyQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.displayNameQuery !== undefined) {
      RoleDisplayNameQuery.encode(message.displayNameQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RoleQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRoleQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.keyQuery = RoleKeyQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.displayNameQuery = RoleDisplayNameQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RoleQuery {
    return {
      keyQuery: isSet(object.keyQuery) ? RoleKeyQuery.fromJSON(object.keyQuery) : undefined,
      displayNameQuery: isSet(object.displayNameQuery)
        ? RoleDisplayNameQuery.fromJSON(object.displayNameQuery)
        : undefined,
    };
  },

  toJSON(message: RoleQuery): unknown {
    const obj: any = {};
    message.keyQuery !== undefined &&
      (obj.keyQuery = message.keyQuery ? RoleKeyQuery.toJSON(message.keyQuery) : undefined);
    message.displayNameQuery !== undefined && (obj.displayNameQuery = message.displayNameQuery
      ? RoleDisplayNameQuery.toJSON(message.displayNameQuery)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<RoleQuery>): RoleQuery {
    return RoleQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RoleQuery>): RoleQuery {
    const message = createBaseRoleQuery();
    message.keyQuery = (object.keyQuery !== undefined && object.keyQuery !== null)
      ? RoleKeyQuery.fromPartial(object.keyQuery)
      : undefined;
    message.displayNameQuery = (object.displayNameQuery !== undefined && object.displayNameQuery !== null)
      ? RoleDisplayNameQuery.fromPartial(object.displayNameQuery)
      : undefined;
    return message;
  },
};

function createBaseRoleKeyQuery(): RoleKeyQuery {
  return { key: "", method: 0 };
}

export const RoleKeyQuery = {
  encode(message: RoleKeyQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.key !== "") {
      writer.uint32(10).string(message.key);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RoleKeyQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRoleKeyQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.key = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RoleKeyQuery {
    return {
      key: isSet(object.key) ? String(object.key) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: RoleKeyQuery): unknown {
    const obj: any = {};
    message.key !== undefined && (obj.key = message.key);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<RoleKeyQuery>): RoleKeyQuery {
    return RoleKeyQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RoleKeyQuery>): RoleKeyQuery {
    const message = createBaseRoleKeyQuery();
    message.key = object.key ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseRoleDisplayNameQuery(): RoleDisplayNameQuery {
  return { displayName: "", method: 0 };
}

export const RoleDisplayNameQuery = {
  encode(message: RoleDisplayNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.displayName !== "") {
      writer.uint32(10).string(message.displayName);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RoleDisplayNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRoleDisplayNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.displayName = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RoleDisplayNameQuery {
    return {
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: RoleDisplayNameQuery): unknown {
    const obj: any = {};
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<RoleDisplayNameQuery>): RoleDisplayNameQuery {
    return RoleDisplayNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<RoleDisplayNameQuery>): RoleDisplayNameQuery {
    const message = createBaseRoleDisplayNameQuery();
    message.displayName = object.displayName ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseProjectGrantQuery(): ProjectGrantQuery {
  return { projectNameQuery: undefined, roleKeyQuery: undefined };
}

export const ProjectGrantQuery = {
  encode(message: ProjectGrantQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectNameQuery !== undefined) {
      GrantProjectNameQuery.encode(message.projectNameQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.roleKeyQuery !== undefined) {
      GrantRoleKeyQuery.encode(message.roleKeyQuery, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProjectGrantQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProjectGrantQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.projectNameQuery = GrantProjectNameQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.roleKeyQuery = GrantRoleKeyQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProjectGrantQuery {
    return {
      projectNameQuery: isSet(object.projectNameQuery)
        ? GrantProjectNameQuery.fromJSON(object.projectNameQuery)
        : undefined,
      roleKeyQuery: isSet(object.roleKeyQuery) ? GrantRoleKeyQuery.fromJSON(object.roleKeyQuery) : undefined,
    };
  },

  toJSON(message: ProjectGrantQuery): unknown {
    const obj: any = {};
    message.projectNameQuery !== undefined && (obj.projectNameQuery = message.projectNameQuery
      ? GrantProjectNameQuery.toJSON(message.projectNameQuery)
      : undefined);
    message.roleKeyQuery !== undefined &&
      (obj.roleKeyQuery = message.roleKeyQuery ? GrantRoleKeyQuery.toJSON(message.roleKeyQuery) : undefined);
    return obj;
  },

  create(base?: DeepPartial<ProjectGrantQuery>): ProjectGrantQuery {
    return ProjectGrantQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProjectGrantQuery>): ProjectGrantQuery {
    const message = createBaseProjectGrantQuery();
    message.projectNameQuery = (object.projectNameQuery !== undefined && object.projectNameQuery !== null)
      ? GrantProjectNameQuery.fromPartial(object.projectNameQuery)
      : undefined;
    message.roleKeyQuery = (object.roleKeyQuery !== undefined && object.roleKeyQuery !== null)
      ? GrantRoleKeyQuery.fromPartial(object.roleKeyQuery)
      : undefined;
    return message;
  },
};

function createBaseAllProjectGrantQuery(): AllProjectGrantQuery {
  return {
    projectNameQuery: undefined,
    roleKeyQuery: undefined,
    projectIdQuery: undefined,
    grantedOrgIdQuery: undefined,
  };
}

export const AllProjectGrantQuery = {
  encode(message: AllProjectGrantQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectNameQuery !== undefined) {
      GrantProjectNameQuery.encode(message.projectNameQuery, writer.uint32(10).fork()).ldelim();
    }
    if (message.roleKeyQuery !== undefined) {
      GrantRoleKeyQuery.encode(message.roleKeyQuery, writer.uint32(18).fork()).ldelim();
    }
    if (message.projectIdQuery !== undefined) {
      ProjectIDQuery.encode(message.projectIdQuery, writer.uint32(26).fork()).ldelim();
    }
    if (message.grantedOrgIdQuery !== undefined) {
      GrantedOrgIDQuery.encode(message.grantedOrgIdQuery, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AllProjectGrantQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAllProjectGrantQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.projectNameQuery = GrantProjectNameQuery.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.roleKeyQuery = GrantRoleKeyQuery.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.projectIdQuery = ProjectIDQuery.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.grantedOrgIdQuery = GrantedOrgIDQuery.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AllProjectGrantQuery {
    return {
      projectNameQuery: isSet(object.projectNameQuery)
        ? GrantProjectNameQuery.fromJSON(object.projectNameQuery)
        : undefined,
      roleKeyQuery: isSet(object.roleKeyQuery) ? GrantRoleKeyQuery.fromJSON(object.roleKeyQuery) : undefined,
      projectIdQuery: isSet(object.projectIdQuery) ? ProjectIDQuery.fromJSON(object.projectIdQuery) : undefined,
      grantedOrgIdQuery: isSet(object.grantedOrgIdQuery)
        ? GrantedOrgIDQuery.fromJSON(object.grantedOrgIdQuery)
        : undefined,
    };
  },

  toJSON(message: AllProjectGrantQuery): unknown {
    const obj: any = {};
    message.projectNameQuery !== undefined && (obj.projectNameQuery = message.projectNameQuery
      ? GrantProjectNameQuery.toJSON(message.projectNameQuery)
      : undefined);
    message.roleKeyQuery !== undefined &&
      (obj.roleKeyQuery = message.roleKeyQuery ? GrantRoleKeyQuery.toJSON(message.roleKeyQuery) : undefined);
    message.projectIdQuery !== undefined &&
      (obj.projectIdQuery = message.projectIdQuery ? ProjectIDQuery.toJSON(message.projectIdQuery) : undefined);
    message.grantedOrgIdQuery !== undefined && (obj.grantedOrgIdQuery = message.grantedOrgIdQuery
      ? GrantedOrgIDQuery.toJSON(message.grantedOrgIdQuery)
      : undefined);
    return obj;
  },

  create(base?: DeepPartial<AllProjectGrantQuery>): AllProjectGrantQuery {
    return AllProjectGrantQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AllProjectGrantQuery>): AllProjectGrantQuery {
    const message = createBaseAllProjectGrantQuery();
    message.projectNameQuery = (object.projectNameQuery !== undefined && object.projectNameQuery !== null)
      ? GrantProjectNameQuery.fromPartial(object.projectNameQuery)
      : undefined;
    message.roleKeyQuery = (object.roleKeyQuery !== undefined && object.roleKeyQuery !== null)
      ? GrantRoleKeyQuery.fromPartial(object.roleKeyQuery)
      : undefined;
    message.projectIdQuery = (object.projectIdQuery !== undefined && object.projectIdQuery !== null)
      ? ProjectIDQuery.fromPartial(object.projectIdQuery)
      : undefined;
    message.grantedOrgIdQuery = (object.grantedOrgIdQuery !== undefined && object.grantedOrgIdQuery !== null)
      ? GrantedOrgIDQuery.fromPartial(object.grantedOrgIdQuery)
      : undefined;
    return message;
  },
};

function createBaseGrantProjectNameQuery(): GrantProjectNameQuery {
  return { name: "", method: 0 };
}

export const GrantProjectNameQuery = {
  encode(message: GrantProjectNameQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GrantProjectNameQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGrantProjectNameQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GrantProjectNameQuery {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: GrantProjectNameQuery): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<GrantProjectNameQuery>): GrantProjectNameQuery {
    return GrantProjectNameQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GrantProjectNameQuery>): GrantProjectNameQuery {
    const message = createBaseGrantProjectNameQuery();
    message.name = object.name ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseGrantRoleKeyQuery(): GrantRoleKeyQuery {
  return { roleKey: "", method: 0 };
}

export const GrantRoleKeyQuery = {
  encode(message: GrantRoleKeyQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.roleKey !== "") {
      writer.uint32(10).string(message.roleKey);
    }
    if (message.method !== 0) {
      writer.uint32(16).int32(message.method);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GrantRoleKeyQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGrantRoleKeyQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.roleKey = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.method = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GrantRoleKeyQuery {
    return {
      roleKey: isSet(object.roleKey) ? String(object.roleKey) : "",
      method: isSet(object.method) ? textQueryMethodFromJSON(object.method) : 0,
    };
  },

  toJSON(message: GrantRoleKeyQuery): unknown {
    const obj: any = {};
    message.roleKey !== undefined && (obj.roleKey = message.roleKey);
    message.method !== undefined && (obj.method = textQueryMethodToJSON(message.method));
    return obj;
  },

  create(base?: DeepPartial<GrantRoleKeyQuery>): GrantRoleKeyQuery {
    return GrantRoleKeyQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GrantRoleKeyQuery>): GrantRoleKeyQuery {
    const message = createBaseGrantRoleKeyQuery();
    message.roleKey = object.roleKey ?? "";
    message.method = object.method ?? 0;
    return message;
  },
};

function createBaseProjectIDQuery(): ProjectIDQuery {
  return { projectId: "" };
}

export const ProjectIDQuery = {
  encode(message: ProjectIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectId !== "") {
      writer.uint32(10).string(message.projectId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ProjectIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProjectIDQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.projectId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ProjectIDQuery {
    return { projectId: isSet(object.projectId) ? String(object.projectId) : "" };
  },

  toJSON(message: ProjectIDQuery): unknown {
    const obj: any = {};
    message.projectId !== undefined && (obj.projectId = message.projectId);
    return obj;
  },

  create(base?: DeepPartial<ProjectIDQuery>): ProjectIDQuery {
    return ProjectIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ProjectIDQuery>): ProjectIDQuery {
    const message = createBaseProjectIDQuery();
    message.projectId = object.projectId ?? "";
    return message;
  },
};

function createBaseGrantedOrgIDQuery(): GrantedOrgIDQuery {
  return { grantedOrgId: "" };
}

export const GrantedOrgIDQuery = {
  encode(message: GrantedOrgIDQuery, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.grantedOrgId !== "") {
      writer.uint32(10).string(message.grantedOrgId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GrantedOrgIDQuery {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGrantedOrgIDQuery();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.grantedOrgId = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GrantedOrgIDQuery {
    return { grantedOrgId: isSet(object.grantedOrgId) ? String(object.grantedOrgId) : "" };
  },

  toJSON(message: GrantedOrgIDQuery): unknown {
    const obj: any = {};
    message.grantedOrgId !== undefined && (obj.grantedOrgId = message.grantedOrgId);
    return obj;
  },

  create(base?: DeepPartial<GrantedOrgIDQuery>): GrantedOrgIDQuery {
    return GrantedOrgIDQuery.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<GrantedOrgIDQuery>): GrantedOrgIDQuery {
    const message = createBaseGrantedOrgIDQuery();
    message.grantedOrgId = object.grantedOrgId ?? "";
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
