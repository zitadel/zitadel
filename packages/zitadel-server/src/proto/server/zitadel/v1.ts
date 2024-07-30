/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Timestamp } from "../google/protobuf/timestamp";
import { KeyType, keyTypeFromJSON, keyTypeToJSON } from "./auth_n_key";
import { IDPUserLink } from "./idp";
import {
  AddAPIAppRequest,
  AddCustomLabelPolicyRequest,
  AddCustomLockoutPolicyRequest,
  AddCustomLoginPolicyRequest,
  AddCustomPasswordComplexityPolicyRequest,
  AddCustomPrivacyPolicyRequest,
  AddIDPToLoginPolicyRequest,
  AddMachineUserRequest,
  AddMultiFactorToLoginPolicyRequest,
  AddOIDCAppRequest,
  AddOrgJWTIDPRequest,
  AddOrgMemberRequest,
  AddOrgOIDCIDPRequest,
  AddOrgRequest,
  AddProjectGrantMemberRequest,
  AddProjectGrantRequest,
  AddProjectMemberRequest,
  AddProjectRequest,
  AddProjectRoleRequest,
  AddSecondFactorToLoginPolicyRequest,
  AddUserGrantRequest,
  CreateActionRequest,
  ImportHumanUserRequest,
  SetCustomDomainClaimedMessageTextRequest,
  SetCustomInitMessageTextRequest,
  SetCustomLoginTextsRequest,
  SetCustomPasswordlessRegistrationMessageTextRequest,
  SetCustomPasswordResetMessageTextRequest,
  SetCustomVerifyEmailMessageTextRequest,
  SetCustomVerifyPhoneMessageTextRequest,
  SetUserMetadataRequest,
} from "./management";
import { Domain } from "./org";
import { Gender, genderFromJSON, genderToJSON } from "./user";

export const protobufPackage = "zitadel.v1.v1";

export enum FlowType {
  FLOW_TYPE_UNSPECIFIED = 0,
  FLOW_TYPE_EXTERNAL_AUTHENTICATION = 1,
  UNRECOGNIZED = -1,
}

export function flowTypeFromJSON(object: any): FlowType {
  switch (object) {
    case 0:
    case "FLOW_TYPE_UNSPECIFIED":
      return FlowType.FLOW_TYPE_UNSPECIFIED;
    case 1:
    case "FLOW_TYPE_EXTERNAL_AUTHENTICATION":
      return FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FlowType.UNRECOGNIZED;
  }
}

export function flowTypeToJSON(object: FlowType): string {
  switch (object) {
    case FlowType.FLOW_TYPE_UNSPECIFIED:
      return "FLOW_TYPE_UNSPECIFIED";
    case FlowType.FLOW_TYPE_EXTERNAL_AUTHENTICATION:
      return "FLOW_TYPE_EXTERNAL_AUTHENTICATION";
    case FlowType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum TriggerType {
  TRIGGER_TYPE_UNSPECIFIED = 0,
  TRIGGER_TYPE_POST_AUTHENTICATION = 1,
  TRIGGER_TYPE_PRE_CREATION = 2,
  TRIGGER_TYPE_POST_CREATION = 3,
  UNRECOGNIZED = -1,
}

export function triggerTypeFromJSON(object: any): TriggerType {
  switch (object) {
    case 0:
    case "TRIGGER_TYPE_UNSPECIFIED":
      return TriggerType.TRIGGER_TYPE_UNSPECIFIED;
    case 1:
    case "TRIGGER_TYPE_POST_AUTHENTICATION":
      return TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION;
    case 2:
    case "TRIGGER_TYPE_PRE_CREATION":
      return TriggerType.TRIGGER_TYPE_PRE_CREATION;
    case 3:
    case "TRIGGER_TYPE_POST_CREATION":
      return TriggerType.TRIGGER_TYPE_POST_CREATION;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TriggerType.UNRECOGNIZED;
  }
}

export function triggerTypeToJSON(object: TriggerType): string {
  switch (object) {
    case TriggerType.TRIGGER_TYPE_UNSPECIFIED:
      return "TRIGGER_TYPE_UNSPECIFIED";
    case TriggerType.TRIGGER_TYPE_POST_AUTHENTICATION:
      return "TRIGGER_TYPE_POST_AUTHENTICATION";
    case TriggerType.TRIGGER_TYPE_PRE_CREATION:
      return "TRIGGER_TYPE_PRE_CREATION";
    case TriggerType.TRIGGER_TYPE_POST_CREATION:
      return "TRIGGER_TYPE_POST_CREATION";
    case TriggerType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface AddCustomOrgIAMPolicyRequest {
  orgId: string;
  /** the username has to end with the domain of its organization (uniqueness is organization based) */
  userLoginMustBeDomain: boolean;
}

export interface ImportDataOrg {
  orgs: DataOrg[];
}

export interface DataOrg {
  orgId: string;
  org: AddOrgRequest | undefined;
  iamPolicy: AddCustomOrgIAMPolicyRequest | undefined;
  labelPolicy: AddCustomLabelPolicyRequest | undefined;
  lockoutPolicy: AddCustomLockoutPolicyRequest | undefined;
  loginPolicy: AddCustomLoginPolicyRequest | undefined;
  passwordComplexityPolicy: AddCustomPasswordComplexityPolicyRequest | undefined;
  privacyPolicy: AddCustomPrivacyPolicyRequest | undefined;
  projects: DataProject[];
  projectRoles: AddProjectRoleRequest[];
  apiApps: DataAPIApplication[];
  oidcApps: DataOIDCApplication[];
  humanUsers: DataHumanUser[];
  machineUsers: DataMachineUser[];
  triggerActions: SetTriggerActionsRequest[];
  actions: DataAction[];
  projectGrants: DataProjectGrant[];
  userGrants: AddUserGrantRequest[];
  orgMembers: AddOrgMemberRequest[];
  projectMembers: AddProjectMemberRequest[];
  projectGrantMembers: AddProjectGrantMemberRequest[];
  userMetadata: SetUserMetadataRequest[];
  loginTexts: SetCustomLoginTextsRequest[];
  initMessages: SetCustomInitMessageTextRequest[];
  passwordResetMessages: SetCustomPasswordResetMessageTextRequest[];
  verifyEmailMessages: SetCustomVerifyEmailMessageTextRequest[];
  verifyPhoneMessages: SetCustomVerifyPhoneMessageTextRequest[];
  domainClaimedMessages: SetCustomDomainClaimedMessageTextRequest[];
  passwordlessRegistrationMessages: SetCustomPasswordlessRegistrationMessageTextRequest[];
  oidcIdps: DataOIDCIDP[];
  jwtIdps: DataJWTIDP[];
  secondFactors: AddSecondFactorToLoginPolicyRequest[];
  multiFactors: AddMultiFactorToLoginPolicyRequest[];
  idps: AddIDPToLoginPolicyRequest[];
  userLinks: IDPUserLink[];
  domains: Domain[];
  appKeys: DataAppKey[];
  machineKeys: DataMachineKey[];
}

export interface DataOIDCIDP {
  idpId: string;
  idp: AddOrgOIDCIDPRequest | undefined;
}

export interface DataJWTIDP {
  idpId: string;
  idp: AddOrgJWTIDPRequest | undefined;
}

export interface ExportHumanUser {
  userName: string;
  profile: ExportHumanUser_Profile | undefined;
  email: ExportHumanUser_Email | undefined;
  phone: ExportHumanUser_Phone | undefined;
  password: string;
  hashedPassword: ExportHumanUser_HashedPassword | undefined;
  passwordChangeRequired: boolean;
  requestPasswordlessRegistration: boolean;
  otpCode: string;
}

export interface ExportHumanUser_Profile {
  firstName: string;
  lastName: string;
  nickName: string;
  displayName: string;
  preferredLanguage: string;
  gender: Gender;
}

export interface ExportHumanUser_Email {
  /** TODO: check if no value is allowed */
  email: string;
  isEmailVerified: boolean;
}

export interface ExportHumanUser_Phone {
  /** has to be a global number */
  phone: string;
  isPhoneVerified: boolean;
}

export interface ExportHumanUser_HashedPassword {
  value: string;
  algorithm: string;
}

export interface DataAppKey {
  id: string;
  projectId: string;
  appId: string;
  clientId: string;
  type: KeyType;
  expirationDate: Date | undefined;
  publicKey: Buffer;
}

export interface DataMachineKey {
  keyId: string;
  userId: string;
  type: KeyType;
  expirationDate: Date | undefined;
  publicKey: Buffer;
}

export interface DataProject {
  projectId: string;
  project: AddProjectRequest | undefined;
}

export interface DataAPIApplication {
  appId: string;
  app: AddAPIAppRequest | undefined;
}

export interface DataOIDCApplication {
  appId: string;
  app: AddOIDCAppRequest | undefined;
}

export interface DataHumanUser {
  userId: string;
  user: ImportHumanUserRequest | undefined;
}

export interface DataMachineUser {
  userId: string;
  user: AddMachineUserRequest | undefined;
}

export interface DataAction {
  actionId: string;
  action: CreateActionRequest | undefined;
}

export interface DataProjectGrant {
  grantId: string;
  projectGrant: AddProjectGrantRequest | undefined;
}

export interface SetTriggerActionsRequest {
  flowType: FlowType;
  triggerType: TriggerType;
  actionIds: string[];
}

function createBaseAddCustomOrgIAMPolicyRequest(): AddCustomOrgIAMPolicyRequest {
  return { orgId: "", userLoginMustBeDomain: false };
}

export const AddCustomOrgIAMPolicyRequest = {
  encode(message: AddCustomOrgIAMPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.userLoginMustBeDomain === true) {
      writer.uint32(16).bool(message.userLoginMustBeDomain);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddCustomOrgIAMPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddCustomOrgIAMPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgId = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.userLoginMustBeDomain = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddCustomOrgIAMPolicyRequest {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      userLoginMustBeDomain: isSet(object.userLoginMustBeDomain) ? Boolean(object.userLoginMustBeDomain) : false,
    };
  },

  toJSON(message: AddCustomOrgIAMPolicyRequest): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.userLoginMustBeDomain !== undefined && (obj.userLoginMustBeDomain = message.userLoginMustBeDomain);
    return obj;
  },

  create(base?: DeepPartial<AddCustomOrgIAMPolicyRequest>): AddCustomOrgIAMPolicyRequest {
    return AddCustomOrgIAMPolicyRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<AddCustomOrgIAMPolicyRequest>): AddCustomOrgIAMPolicyRequest {
    const message = createBaseAddCustomOrgIAMPolicyRequest();
    message.orgId = object.orgId ?? "";
    message.userLoginMustBeDomain = object.userLoginMustBeDomain ?? false;
    return message;
  },
};

function createBaseImportDataOrg(): ImportDataOrg {
  return { orgs: [] };
}

export const ImportDataOrg = {
  encode(message: ImportDataOrg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.orgs) {
      DataOrg.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ImportDataOrg {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseImportDataOrg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgs.push(DataOrg.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ImportDataOrg {
    return { orgs: Array.isArray(object?.orgs) ? object.orgs.map((e: any) => DataOrg.fromJSON(e)) : [] };
  },

  toJSON(message: ImportDataOrg): unknown {
    const obj: any = {};
    if (message.orgs) {
      obj.orgs = message.orgs.map((e) => e ? DataOrg.toJSON(e) : undefined);
    } else {
      obj.orgs = [];
    }
    return obj;
  },

  create(base?: DeepPartial<ImportDataOrg>): ImportDataOrg {
    return ImportDataOrg.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ImportDataOrg>): ImportDataOrg {
    const message = createBaseImportDataOrg();
    message.orgs = object.orgs?.map((e) => DataOrg.fromPartial(e)) || [];
    return message;
  },
};

function createBaseDataOrg(): DataOrg {
  return {
    orgId: "",
    org: undefined,
    iamPolicy: undefined,
    labelPolicy: undefined,
    lockoutPolicy: undefined,
    loginPolicy: undefined,
    passwordComplexityPolicy: undefined,
    privacyPolicy: undefined,
    projects: [],
    projectRoles: [],
    apiApps: [],
    oidcApps: [],
    humanUsers: [],
    machineUsers: [],
    triggerActions: [],
    actions: [],
    projectGrants: [],
    userGrants: [],
    orgMembers: [],
    projectMembers: [],
    projectGrantMembers: [],
    userMetadata: [],
    loginTexts: [],
    initMessages: [],
    passwordResetMessages: [],
    verifyEmailMessages: [],
    verifyPhoneMessages: [],
    domainClaimedMessages: [],
    passwordlessRegistrationMessages: [],
    oidcIdps: [],
    jwtIdps: [],
    secondFactors: [],
    multiFactors: [],
    idps: [],
    userLinks: [],
    domains: [],
    appKeys: [],
    machineKeys: [],
  };
}

export const DataOrg = {
  encode(message: DataOrg, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orgId !== "") {
      writer.uint32(10).string(message.orgId);
    }
    if (message.org !== undefined) {
      AddOrgRequest.encode(message.org, writer.uint32(26).fork()).ldelim();
    }
    if (message.iamPolicy !== undefined) {
      AddCustomOrgIAMPolicyRequest.encode(message.iamPolicy, writer.uint32(34).fork()).ldelim();
    }
    if (message.labelPolicy !== undefined) {
      AddCustomLabelPolicyRequest.encode(message.labelPolicy, writer.uint32(42).fork()).ldelim();
    }
    if (message.lockoutPolicy !== undefined) {
      AddCustomLockoutPolicyRequest.encode(message.lockoutPolicy, writer.uint32(50).fork()).ldelim();
    }
    if (message.loginPolicy !== undefined) {
      AddCustomLoginPolicyRequest.encode(message.loginPolicy, writer.uint32(58).fork()).ldelim();
    }
    if (message.passwordComplexityPolicy !== undefined) {
      AddCustomPasswordComplexityPolicyRequest.encode(message.passwordComplexityPolicy, writer.uint32(66).fork())
        .ldelim();
    }
    if (message.privacyPolicy !== undefined) {
      AddCustomPrivacyPolicyRequest.encode(message.privacyPolicy, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.projects) {
      DataProject.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    for (const v of message.projectRoles) {
      AddProjectRoleRequest.encode(v!, writer.uint32(90).fork()).ldelim();
    }
    for (const v of message.apiApps) {
      DataAPIApplication.encode(v!, writer.uint32(98).fork()).ldelim();
    }
    for (const v of message.oidcApps) {
      DataOIDCApplication.encode(v!, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.humanUsers) {
      DataHumanUser.encode(v!, writer.uint32(114).fork()).ldelim();
    }
    for (const v of message.machineUsers) {
      DataMachineUser.encode(v!, writer.uint32(122).fork()).ldelim();
    }
    for (const v of message.triggerActions) {
      SetTriggerActionsRequest.encode(v!, writer.uint32(130).fork()).ldelim();
    }
    for (const v of message.actions) {
      DataAction.encode(v!, writer.uint32(138).fork()).ldelim();
    }
    for (const v of message.projectGrants) {
      DataProjectGrant.encode(v!, writer.uint32(146).fork()).ldelim();
    }
    for (const v of message.userGrants) {
      AddUserGrantRequest.encode(v!, writer.uint32(154).fork()).ldelim();
    }
    for (const v of message.orgMembers) {
      AddOrgMemberRequest.encode(v!, writer.uint32(162).fork()).ldelim();
    }
    for (const v of message.projectMembers) {
      AddProjectMemberRequest.encode(v!, writer.uint32(170).fork()).ldelim();
    }
    for (const v of message.projectGrantMembers) {
      AddProjectGrantMemberRequest.encode(v!, writer.uint32(178).fork()).ldelim();
    }
    for (const v of message.userMetadata) {
      SetUserMetadataRequest.encode(v!, writer.uint32(186).fork()).ldelim();
    }
    for (const v of message.loginTexts) {
      SetCustomLoginTextsRequest.encode(v!, writer.uint32(194).fork()).ldelim();
    }
    for (const v of message.initMessages) {
      SetCustomInitMessageTextRequest.encode(v!, writer.uint32(202).fork()).ldelim();
    }
    for (const v of message.passwordResetMessages) {
      SetCustomPasswordResetMessageTextRequest.encode(v!, writer.uint32(210).fork()).ldelim();
    }
    for (const v of message.verifyEmailMessages) {
      SetCustomVerifyEmailMessageTextRequest.encode(v!, writer.uint32(218).fork()).ldelim();
    }
    for (const v of message.verifyPhoneMessages) {
      SetCustomVerifyPhoneMessageTextRequest.encode(v!, writer.uint32(226).fork()).ldelim();
    }
    for (const v of message.domainClaimedMessages) {
      SetCustomDomainClaimedMessageTextRequest.encode(v!, writer.uint32(234).fork()).ldelim();
    }
    for (const v of message.passwordlessRegistrationMessages) {
      SetCustomPasswordlessRegistrationMessageTextRequest.encode(v!, writer.uint32(242).fork()).ldelim();
    }
    for (const v of message.oidcIdps) {
      DataOIDCIDP.encode(v!, writer.uint32(250).fork()).ldelim();
    }
    for (const v of message.jwtIdps) {
      DataJWTIDP.encode(v!, writer.uint32(258).fork()).ldelim();
    }
    for (const v of message.secondFactors) {
      AddSecondFactorToLoginPolicyRequest.encode(v!, writer.uint32(266).fork()).ldelim();
    }
    for (const v of message.multiFactors) {
      AddMultiFactorToLoginPolicyRequest.encode(v!, writer.uint32(274).fork()).ldelim();
    }
    for (const v of message.idps) {
      AddIDPToLoginPolicyRequest.encode(v!, writer.uint32(282).fork()).ldelim();
    }
    for (const v of message.userLinks) {
      IDPUserLink.encode(v!, writer.uint32(290).fork()).ldelim();
    }
    for (const v of message.domains) {
      Domain.encode(v!, writer.uint32(298).fork()).ldelim();
    }
    for (const v of message.appKeys) {
      DataAppKey.encode(v!, writer.uint32(306).fork()).ldelim();
    }
    for (const v of message.machineKeys) {
      DataMachineKey.encode(v!, writer.uint32(314).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataOrg {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataOrg();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.orgId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.org = AddOrgRequest.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.iamPolicy = AddCustomOrgIAMPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.labelPolicy = AddCustomLabelPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.lockoutPolicy = AddCustomLockoutPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.loginPolicy = AddCustomLoginPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 8:
          if (tag != 66) {
            break;
          }

          message.passwordComplexityPolicy = AddCustomPasswordComplexityPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.privacyPolicy = AddCustomPrivacyPolicyRequest.decode(reader, reader.uint32());
          continue;
        case 10:
          if (tag != 82) {
            break;
          }

          message.projects.push(DataProject.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag != 90) {
            break;
          }

          message.projectRoles.push(AddProjectRoleRequest.decode(reader, reader.uint32()));
          continue;
        case 12:
          if (tag != 98) {
            break;
          }

          message.apiApps.push(DataAPIApplication.decode(reader, reader.uint32()));
          continue;
        case 13:
          if (tag != 106) {
            break;
          }

          message.oidcApps.push(DataOIDCApplication.decode(reader, reader.uint32()));
          continue;
        case 14:
          if (tag != 114) {
            break;
          }

          message.humanUsers.push(DataHumanUser.decode(reader, reader.uint32()));
          continue;
        case 15:
          if (tag != 122) {
            break;
          }

          message.machineUsers.push(DataMachineUser.decode(reader, reader.uint32()));
          continue;
        case 16:
          if (tag != 130) {
            break;
          }

          message.triggerActions.push(SetTriggerActionsRequest.decode(reader, reader.uint32()));
          continue;
        case 17:
          if (tag != 138) {
            break;
          }

          message.actions.push(DataAction.decode(reader, reader.uint32()));
          continue;
        case 18:
          if (tag != 146) {
            break;
          }

          message.projectGrants.push(DataProjectGrant.decode(reader, reader.uint32()));
          continue;
        case 19:
          if (tag != 154) {
            break;
          }

          message.userGrants.push(AddUserGrantRequest.decode(reader, reader.uint32()));
          continue;
        case 20:
          if (tag != 162) {
            break;
          }

          message.orgMembers.push(AddOrgMemberRequest.decode(reader, reader.uint32()));
          continue;
        case 21:
          if (tag != 170) {
            break;
          }

          message.projectMembers.push(AddProjectMemberRequest.decode(reader, reader.uint32()));
          continue;
        case 22:
          if (tag != 178) {
            break;
          }

          message.projectGrantMembers.push(AddProjectGrantMemberRequest.decode(reader, reader.uint32()));
          continue;
        case 23:
          if (tag != 186) {
            break;
          }

          message.userMetadata.push(SetUserMetadataRequest.decode(reader, reader.uint32()));
          continue;
        case 24:
          if (tag != 194) {
            break;
          }

          message.loginTexts.push(SetCustomLoginTextsRequest.decode(reader, reader.uint32()));
          continue;
        case 25:
          if (tag != 202) {
            break;
          }

          message.initMessages.push(SetCustomInitMessageTextRequest.decode(reader, reader.uint32()));
          continue;
        case 26:
          if (tag != 210) {
            break;
          }

          message.passwordResetMessages.push(SetCustomPasswordResetMessageTextRequest.decode(reader, reader.uint32()));
          continue;
        case 27:
          if (tag != 218) {
            break;
          }

          message.verifyEmailMessages.push(SetCustomVerifyEmailMessageTextRequest.decode(reader, reader.uint32()));
          continue;
        case 28:
          if (tag != 226) {
            break;
          }

          message.verifyPhoneMessages.push(SetCustomVerifyPhoneMessageTextRequest.decode(reader, reader.uint32()));
          continue;
        case 29:
          if (tag != 234) {
            break;
          }

          message.domainClaimedMessages.push(SetCustomDomainClaimedMessageTextRequest.decode(reader, reader.uint32()));
          continue;
        case 30:
          if (tag != 242) {
            break;
          }

          message.passwordlessRegistrationMessages.push(
            SetCustomPasswordlessRegistrationMessageTextRequest.decode(reader, reader.uint32()),
          );
          continue;
        case 31:
          if (tag != 250) {
            break;
          }

          message.oidcIdps.push(DataOIDCIDP.decode(reader, reader.uint32()));
          continue;
        case 32:
          if (tag != 258) {
            break;
          }

          message.jwtIdps.push(DataJWTIDP.decode(reader, reader.uint32()));
          continue;
        case 33:
          if (tag != 266) {
            break;
          }

          message.secondFactors.push(AddSecondFactorToLoginPolicyRequest.decode(reader, reader.uint32()));
          continue;
        case 34:
          if (tag != 274) {
            break;
          }

          message.multiFactors.push(AddMultiFactorToLoginPolicyRequest.decode(reader, reader.uint32()));
          continue;
        case 35:
          if (tag != 282) {
            break;
          }

          message.idps.push(AddIDPToLoginPolicyRequest.decode(reader, reader.uint32()));
          continue;
        case 36:
          if (tag != 290) {
            break;
          }

          message.userLinks.push(IDPUserLink.decode(reader, reader.uint32()));
          continue;
        case 37:
          if (tag != 298) {
            break;
          }

          message.domains.push(Domain.decode(reader, reader.uint32()));
          continue;
        case 38:
          if (tag != 306) {
            break;
          }

          message.appKeys.push(DataAppKey.decode(reader, reader.uint32()));
          continue;
        case 39:
          if (tag != 314) {
            break;
          }

          message.machineKeys.push(DataMachineKey.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataOrg {
    return {
      orgId: isSet(object.orgId) ? String(object.orgId) : "",
      org: isSet(object.org) ? AddOrgRequest.fromJSON(object.org) : undefined,
      iamPolicy: isSet(object.iamPolicy) ? AddCustomOrgIAMPolicyRequest.fromJSON(object.iamPolicy) : undefined,
      labelPolicy: isSet(object.labelPolicy) ? AddCustomLabelPolicyRequest.fromJSON(object.labelPolicy) : undefined,
      lockoutPolicy: isSet(object.lockoutPolicy)
        ? AddCustomLockoutPolicyRequest.fromJSON(object.lockoutPolicy)
        : undefined,
      loginPolicy: isSet(object.loginPolicy) ? AddCustomLoginPolicyRequest.fromJSON(object.loginPolicy) : undefined,
      passwordComplexityPolicy: isSet(object.passwordComplexityPolicy)
        ? AddCustomPasswordComplexityPolicyRequest.fromJSON(object.passwordComplexityPolicy)
        : undefined,
      privacyPolicy: isSet(object.privacyPolicy)
        ? AddCustomPrivacyPolicyRequest.fromJSON(object.privacyPolicy)
        : undefined,
      projects: Array.isArray(object?.projects) ? object.projects.map((e: any) => DataProject.fromJSON(e)) : [],
      projectRoles: Array.isArray(object?.projectRoles)
        ? object.projectRoles.map((e: any) => AddProjectRoleRequest.fromJSON(e))
        : [],
      apiApps: Array.isArray(object?.apiApps) ? object.apiApps.map((e: any) => DataAPIApplication.fromJSON(e)) : [],
      oidcApps: Array.isArray(object?.oidcApps) ? object.oidcApps.map((e: any) => DataOIDCApplication.fromJSON(e)) : [],
      humanUsers: Array.isArray(object?.humanUsers) ? object.humanUsers.map((e: any) => DataHumanUser.fromJSON(e)) : [],
      machineUsers: Array.isArray(object?.machineUsers)
        ? object.machineUsers.map((e: any) => DataMachineUser.fromJSON(e))
        : [],
      triggerActions: Array.isArray(object?.triggerActions)
        ? object.triggerActions.map((e: any) => SetTriggerActionsRequest.fromJSON(e))
        : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => DataAction.fromJSON(e)) : [],
      projectGrants: Array.isArray(object?.projectGrants)
        ? object.projectGrants.map((e: any) => DataProjectGrant.fromJSON(e))
        : [],
      userGrants: Array.isArray(object?.userGrants)
        ? object.userGrants.map((e: any) => AddUserGrantRequest.fromJSON(e))
        : [],
      orgMembers: Array.isArray(object?.orgMembers)
        ? object.orgMembers.map((e: any) => AddOrgMemberRequest.fromJSON(e))
        : [],
      projectMembers: Array.isArray(object?.projectMembers)
        ? object.projectMembers.map((e: any) => AddProjectMemberRequest.fromJSON(e))
        : [],
      projectGrantMembers: Array.isArray(object?.projectGrantMembers)
        ? object.projectGrantMembers.map((e: any) => AddProjectGrantMemberRequest.fromJSON(e))
        : [],
      userMetadata: Array.isArray(object?.userMetadata)
        ? object.userMetadata.map((e: any) => SetUserMetadataRequest.fromJSON(e))
        : [],
      loginTexts: Array.isArray(object?.loginTexts)
        ? object.loginTexts.map((e: any) => SetCustomLoginTextsRequest.fromJSON(e))
        : [],
      initMessages: Array.isArray(object?.initMessages)
        ? object.initMessages.map((e: any) => SetCustomInitMessageTextRequest.fromJSON(e))
        : [],
      passwordResetMessages: Array.isArray(object?.passwordResetMessages)
        ? object.passwordResetMessages.map((e: any) => SetCustomPasswordResetMessageTextRequest.fromJSON(e))
        : [],
      verifyEmailMessages: Array.isArray(object?.verifyEmailMessages)
        ? object.verifyEmailMessages.map((e: any) => SetCustomVerifyEmailMessageTextRequest.fromJSON(e))
        : [],
      verifyPhoneMessages: Array.isArray(object?.verifyPhoneMessages)
        ? object.verifyPhoneMessages.map((e: any) => SetCustomVerifyPhoneMessageTextRequest.fromJSON(e))
        : [],
      domainClaimedMessages: Array.isArray(object?.domainClaimedMessages)
        ? object.domainClaimedMessages.map((e: any) => SetCustomDomainClaimedMessageTextRequest.fromJSON(e))
        : [],
      passwordlessRegistrationMessages: Array.isArray(object?.passwordlessRegistrationMessages)
        ? object.passwordlessRegistrationMessages.map((e: any) =>
          SetCustomPasswordlessRegistrationMessageTextRequest.fromJSON(e)
        )
        : [],
      oidcIdps: Array.isArray(object?.oidcIdps)
        ? object.oidcIdps.map((e: any) => DataOIDCIDP.fromJSON(e))
        : [],
      jwtIdps: Array.isArray(object?.jwtIdps)
        ? object.jwtIdps.map((e: any) => DataJWTIDP.fromJSON(e))
        : [],
      secondFactors: Array.isArray(object?.secondFactors)
        ? object.secondFactors.map((e: any) => AddSecondFactorToLoginPolicyRequest.fromJSON(e))
        : [],
      multiFactors: Array.isArray(object?.multiFactors)
        ? object.multiFactors.map((e: any) => AddMultiFactorToLoginPolicyRequest.fromJSON(e))
        : [],
      idps: Array.isArray(object?.idps) ? object.idps.map((e: any) => AddIDPToLoginPolicyRequest.fromJSON(e)) : [],
      userLinks: Array.isArray(object?.userLinks) ? object.userLinks.map((e: any) => IDPUserLink.fromJSON(e)) : [],
      domains: Array.isArray(object?.domains) ? object.domains.map((e: any) => Domain.fromJSON(e)) : [],
      appKeys: Array.isArray(object?.appKeys) ? object.appKeys.map((e: any) => DataAppKey.fromJSON(e)) : [],
      machineKeys: Array.isArray(object?.machineKeys)
        ? object.machineKeys.map((e: any) => DataMachineKey.fromJSON(e))
        : [],
    };
  },

  toJSON(message: DataOrg): unknown {
    const obj: any = {};
    message.orgId !== undefined && (obj.orgId = message.orgId);
    message.org !== undefined && (obj.org = message.org ? AddOrgRequest.toJSON(message.org) : undefined);
    message.iamPolicy !== undefined &&
      (obj.iamPolicy = message.iamPolicy ? AddCustomOrgIAMPolicyRequest.toJSON(message.iamPolicy) : undefined);
    message.labelPolicy !== undefined &&
      (obj.labelPolicy = message.labelPolicy ? AddCustomLabelPolicyRequest.toJSON(message.labelPolicy) : undefined);
    message.lockoutPolicy !== undefined && (obj.lockoutPolicy = message.lockoutPolicy
      ? AddCustomLockoutPolicyRequest.toJSON(message.lockoutPolicy)
      : undefined);
    message.loginPolicy !== undefined &&
      (obj.loginPolicy = message.loginPolicy ? AddCustomLoginPolicyRequest.toJSON(message.loginPolicy) : undefined);
    message.passwordComplexityPolicy !== undefined && (obj.passwordComplexityPolicy = message.passwordComplexityPolicy
      ? AddCustomPasswordComplexityPolicyRequest.toJSON(message.passwordComplexityPolicy)
      : undefined);
    message.privacyPolicy !== undefined && (obj.privacyPolicy = message.privacyPolicy
      ? AddCustomPrivacyPolicyRequest.toJSON(message.privacyPolicy)
      : undefined);
    if (message.projects) {
      obj.projects = message.projects.map((e) => e ? DataProject.toJSON(e) : undefined);
    } else {
      obj.projects = [];
    }
    if (message.projectRoles) {
      obj.projectRoles = message.projectRoles.map((e) => e ? AddProjectRoleRequest.toJSON(e) : undefined);
    } else {
      obj.projectRoles = [];
    }
    if (message.apiApps) {
      obj.apiApps = message.apiApps.map((e) => e ? DataAPIApplication.toJSON(e) : undefined);
    } else {
      obj.apiApps = [];
    }
    if (message.oidcApps) {
      obj.oidcApps = message.oidcApps.map((e) => e ? DataOIDCApplication.toJSON(e) : undefined);
    } else {
      obj.oidcApps = [];
    }
    if (message.humanUsers) {
      obj.humanUsers = message.humanUsers.map((e) => e ? DataHumanUser.toJSON(e) : undefined);
    } else {
      obj.humanUsers = [];
    }
    if (message.machineUsers) {
      obj.machineUsers = message.machineUsers.map((e) => e ? DataMachineUser.toJSON(e) : undefined);
    } else {
      obj.machineUsers = [];
    }
    if (message.triggerActions) {
      obj.triggerActions = message.triggerActions.map((e) => e ? SetTriggerActionsRequest.toJSON(e) : undefined);
    } else {
      obj.triggerActions = [];
    }
    if (message.actions) {
      obj.actions = message.actions.map((e) => e ? DataAction.toJSON(e) : undefined);
    } else {
      obj.actions = [];
    }
    if (message.projectGrants) {
      obj.projectGrants = message.projectGrants.map((e) => e ? DataProjectGrant.toJSON(e) : undefined);
    } else {
      obj.projectGrants = [];
    }
    if (message.userGrants) {
      obj.userGrants = message.userGrants.map((e) => e ? AddUserGrantRequest.toJSON(e) : undefined);
    } else {
      obj.userGrants = [];
    }
    if (message.orgMembers) {
      obj.orgMembers = message.orgMembers.map((e) => e ? AddOrgMemberRequest.toJSON(e) : undefined);
    } else {
      obj.orgMembers = [];
    }
    if (message.projectMembers) {
      obj.projectMembers = message.projectMembers.map((e) => e ? AddProjectMemberRequest.toJSON(e) : undefined);
    } else {
      obj.projectMembers = [];
    }
    if (message.projectGrantMembers) {
      obj.projectGrantMembers = message.projectGrantMembers.map((e) =>
        e ? AddProjectGrantMemberRequest.toJSON(e) : undefined
      );
    } else {
      obj.projectGrantMembers = [];
    }
    if (message.userMetadata) {
      obj.userMetadata = message.userMetadata.map((e) => e ? SetUserMetadataRequest.toJSON(e) : undefined);
    } else {
      obj.userMetadata = [];
    }
    if (message.loginTexts) {
      obj.loginTexts = message.loginTexts.map((e) => e ? SetCustomLoginTextsRequest.toJSON(e) : undefined);
    } else {
      obj.loginTexts = [];
    }
    if (message.initMessages) {
      obj.initMessages = message.initMessages.map((e) => e ? SetCustomInitMessageTextRequest.toJSON(e) : undefined);
    } else {
      obj.initMessages = [];
    }
    if (message.passwordResetMessages) {
      obj.passwordResetMessages = message.passwordResetMessages.map((e) =>
        e ? SetCustomPasswordResetMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.passwordResetMessages = [];
    }
    if (message.verifyEmailMessages) {
      obj.verifyEmailMessages = message.verifyEmailMessages.map((e) =>
        e ? SetCustomVerifyEmailMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.verifyEmailMessages = [];
    }
    if (message.verifyPhoneMessages) {
      obj.verifyPhoneMessages = message.verifyPhoneMessages.map((e) =>
        e ? SetCustomVerifyPhoneMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.verifyPhoneMessages = [];
    }
    if (message.domainClaimedMessages) {
      obj.domainClaimedMessages = message.domainClaimedMessages.map((e) =>
        e ? SetCustomDomainClaimedMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.domainClaimedMessages = [];
    }
    if (message.passwordlessRegistrationMessages) {
      obj.passwordlessRegistrationMessages = message.passwordlessRegistrationMessages.map((e) =>
        e ? SetCustomPasswordlessRegistrationMessageTextRequest.toJSON(e) : undefined
      );
    } else {
      obj.passwordlessRegistrationMessages = [];
    }
    if (message.oidcIdps) {
      obj.oidcIdps = message.oidcIdps.map((e) => e ? DataOIDCIDP.toJSON(e) : undefined);
    } else {
      obj.oidcIdps = [];
    }
    if (message.jwtIdps) {
      obj.jwtIdps = message.jwtIdps.map((e) => e ? DataJWTIDP.toJSON(e) : undefined);
    } else {
      obj.jwtIdps = [];
    }
    if (message.secondFactors) {
      obj.secondFactors = message.secondFactors.map((e) =>
        e ? AddSecondFactorToLoginPolicyRequest.toJSON(e) : undefined
      );
    } else {
      obj.secondFactors = [];
    }
    if (message.multiFactors) {
      obj.multiFactors = message.multiFactors.map((e) => e ? AddMultiFactorToLoginPolicyRequest.toJSON(e) : undefined);
    } else {
      obj.multiFactors = [];
    }
    if (message.idps) {
      obj.idps = message.idps.map((e) => e ? AddIDPToLoginPolicyRequest.toJSON(e) : undefined);
    } else {
      obj.idps = [];
    }
    if (message.userLinks) {
      obj.userLinks = message.userLinks.map((e) => e ? IDPUserLink.toJSON(e) : undefined);
    } else {
      obj.userLinks = [];
    }
    if (message.domains) {
      obj.domains = message.domains.map((e) => e ? Domain.toJSON(e) : undefined);
    } else {
      obj.domains = [];
    }
    if (message.appKeys) {
      obj.appKeys = message.appKeys.map((e) => e ? DataAppKey.toJSON(e) : undefined);
    } else {
      obj.appKeys = [];
    }
    if (message.machineKeys) {
      obj.machineKeys = message.machineKeys.map((e) => e ? DataMachineKey.toJSON(e) : undefined);
    } else {
      obj.machineKeys = [];
    }
    return obj;
  },

  create(base?: DeepPartial<DataOrg>): DataOrg {
    return DataOrg.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataOrg>): DataOrg {
    const message = createBaseDataOrg();
    message.orgId = object.orgId ?? "";
    message.org = (object.org !== undefined && object.org !== null) ? AddOrgRequest.fromPartial(object.org) : undefined;
    message.iamPolicy = (object.iamPolicy !== undefined && object.iamPolicy !== null)
      ? AddCustomOrgIAMPolicyRequest.fromPartial(object.iamPolicy)
      : undefined;
    message.labelPolicy = (object.labelPolicy !== undefined && object.labelPolicy !== null)
      ? AddCustomLabelPolicyRequest.fromPartial(object.labelPolicy)
      : undefined;
    message.lockoutPolicy = (object.lockoutPolicy !== undefined && object.lockoutPolicy !== null)
      ? AddCustomLockoutPolicyRequest.fromPartial(object.lockoutPolicy)
      : undefined;
    message.loginPolicy = (object.loginPolicy !== undefined && object.loginPolicy !== null)
      ? AddCustomLoginPolicyRequest.fromPartial(object.loginPolicy)
      : undefined;
    message.passwordComplexityPolicy =
      (object.passwordComplexityPolicy !== undefined && object.passwordComplexityPolicy !== null)
        ? AddCustomPasswordComplexityPolicyRequest.fromPartial(object.passwordComplexityPolicy)
        : undefined;
    message.privacyPolicy = (object.privacyPolicy !== undefined && object.privacyPolicy !== null)
      ? AddCustomPrivacyPolicyRequest.fromPartial(object.privacyPolicy)
      : undefined;
    message.projects = object.projects?.map((e) => DataProject.fromPartial(e)) || [];
    message.projectRoles = object.projectRoles?.map((e) => AddProjectRoleRequest.fromPartial(e)) || [];
    message.apiApps = object.apiApps?.map((e) => DataAPIApplication.fromPartial(e)) || [];
    message.oidcApps = object.oidcApps?.map((e) => DataOIDCApplication.fromPartial(e)) || [];
    message.humanUsers = object.humanUsers?.map((e) => DataHumanUser.fromPartial(e)) || [];
    message.machineUsers = object.machineUsers?.map((e) => DataMachineUser.fromPartial(e)) || [];
    message.triggerActions = object.triggerActions?.map((e) => SetTriggerActionsRequest.fromPartial(e)) || [];
    message.actions = object.actions?.map((e) => DataAction.fromPartial(e)) || [];
    message.projectGrants = object.projectGrants?.map((e) => DataProjectGrant.fromPartial(e)) || [];
    message.userGrants = object.userGrants?.map((e) => AddUserGrantRequest.fromPartial(e)) || [];
    message.orgMembers = object.orgMembers?.map((e) => AddOrgMemberRequest.fromPartial(e)) || [];
    message.projectMembers = object.projectMembers?.map((e) => AddProjectMemberRequest.fromPartial(e)) || [];
    message.projectGrantMembers = object.projectGrantMembers?.map((e) => AddProjectGrantMemberRequest.fromPartial(e)) ||
      [];
    message.userMetadata = object.userMetadata?.map((e) => SetUserMetadataRequest.fromPartial(e)) || [];
    message.loginTexts = object.loginTexts?.map((e) => SetCustomLoginTextsRequest.fromPartial(e)) || [];
    message.initMessages = object.initMessages?.map((e) => SetCustomInitMessageTextRequest.fromPartial(e)) || [];
    message.passwordResetMessages =
      object.passwordResetMessages?.map((e) => SetCustomPasswordResetMessageTextRequest.fromPartial(e)) || [];
    message.verifyEmailMessages =
      object.verifyEmailMessages?.map((e) => SetCustomVerifyEmailMessageTextRequest.fromPartial(e)) || [];
    message.verifyPhoneMessages =
      object.verifyPhoneMessages?.map((e) => SetCustomVerifyPhoneMessageTextRequest.fromPartial(e)) || [];
    message.domainClaimedMessages =
      object.domainClaimedMessages?.map((e) => SetCustomDomainClaimedMessageTextRequest.fromPartial(e)) || [];
    message.passwordlessRegistrationMessages =
      object.passwordlessRegistrationMessages?.map((e) =>
        SetCustomPasswordlessRegistrationMessageTextRequest.fromPartial(e)
      ) || [];
    message.oidcIdps = object.oidcIdps?.map((e) => DataOIDCIDP.fromPartial(e)) || [];
    message.jwtIdps = object.jwtIdps?.map((e) => DataJWTIDP.fromPartial(e)) || [];
    message.secondFactors = object.secondFactors?.map((e) => AddSecondFactorToLoginPolicyRequest.fromPartial(e)) || [];
    message.multiFactors = object.multiFactors?.map((e) => AddMultiFactorToLoginPolicyRequest.fromPartial(e)) || [];
    message.idps = object.idps?.map((e) => AddIDPToLoginPolicyRequest.fromPartial(e)) || [];
    message.userLinks = object.userLinks?.map((e) => IDPUserLink.fromPartial(e)) || [];
    message.domains = object.domains?.map((e) => Domain.fromPartial(e)) || [];
    message.appKeys = object.appKeys?.map((e) => DataAppKey.fromPartial(e)) || [];
    message.machineKeys = object.machineKeys?.map((e) => DataMachineKey.fromPartial(e)) || [];
    return message;
  },
};

function createBaseDataOIDCIDP(): DataOIDCIDP {
  return { idpId: "", idp: undefined };
}

export const DataOIDCIDP = {
  encode(message: DataOIDCIDP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.idp !== undefined) {
      AddOrgOIDCIDPRequest.encode(message.idp, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataOIDCIDP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataOIDCIDP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.idp = AddOrgOIDCIDPRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataOIDCIDP {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      idp: isSet(object.idp) ? AddOrgOIDCIDPRequest.fromJSON(object.idp) : undefined,
    };
  },

  toJSON(message: DataOIDCIDP): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.idp !== undefined && (obj.idp = message.idp ? AddOrgOIDCIDPRequest.toJSON(message.idp) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataOIDCIDP>): DataOIDCIDP {
    return DataOIDCIDP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataOIDCIDP>): DataOIDCIDP {
    const message = createBaseDataOIDCIDP();
    message.idpId = object.idpId ?? "";
    message.idp = (object.idp !== undefined && object.idp !== null)
      ? AddOrgOIDCIDPRequest.fromPartial(object.idp)
      : undefined;
    return message;
  },
};

function createBaseDataJWTIDP(): DataJWTIDP {
  return { idpId: "", idp: undefined };
}

export const DataJWTIDP = {
  encode(message: DataJWTIDP, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.idpId !== "") {
      writer.uint32(10).string(message.idpId);
    }
    if (message.idp !== undefined) {
      AddOrgJWTIDPRequest.encode(message.idp, writer.uint32(258).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataJWTIDP {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataJWTIDP();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.idpId = reader.string();
          continue;
        case 32:
          if (tag != 258) {
            break;
          }

          message.idp = AddOrgJWTIDPRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataJWTIDP {
    return {
      idpId: isSet(object.idpId) ? String(object.idpId) : "",
      idp: isSet(object.idp) ? AddOrgJWTIDPRequest.fromJSON(object.idp) : undefined,
    };
  },

  toJSON(message: DataJWTIDP): unknown {
    const obj: any = {};
    message.idpId !== undefined && (obj.idpId = message.idpId);
    message.idp !== undefined && (obj.idp = message.idp ? AddOrgJWTIDPRequest.toJSON(message.idp) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataJWTIDP>): DataJWTIDP {
    return DataJWTIDP.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataJWTIDP>): DataJWTIDP {
    const message = createBaseDataJWTIDP();
    message.idpId = object.idpId ?? "";
    message.idp = (object.idp !== undefined && object.idp !== null)
      ? AddOrgJWTIDPRequest.fromPartial(object.idp)
      : undefined;
    return message;
  },
};

function createBaseExportHumanUser(): ExportHumanUser {
  return {
    userName: "",
    profile: undefined,
    email: undefined,
    phone: undefined,
    password: "",
    hashedPassword: undefined,
    passwordChangeRequired: false,
    requestPasswordlessRegistration: false,
    otpCode: "",
  };
}

export const ExportHumanUser = {
  encode(message: ExportHumanUser, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userName !== "") {
      writer.uint32(10).string(message.userName);
    }
    if (message.profile !== undefined) {
      ExportHumanUser_Profile.encode(message.profile, writer.uint32(18).fork()).ldelim();
    }
    if (message.email !== undefined) {
      ExportHumanUser_Email.encode(message.email, writer.uint32(26).fork()).ldelim();
    }
    if (message.phone !== undefined) {
      ExportHumanUser_Phone.encode(message.phone, writer.uint32(34).fork()).ldelim();
    }
    if (message.password !== "") {
      writer.uint32(42).string(message.password);
    }
    if (message.hashedPassword !== undefined) {
      ExportHumanUser_HashedPassword.encode(message.hashedPassword, writer.uint32(50).fork()).ldelim();
    }
    if (message.passwordChangeRequired === true) {
      writer.uint32(56).bool(message.passwordChangeRequired);
    }
    if (message.requestPasswordlessRegistration === true) {
      writer.uint32(64).bool(message.requestPasswordlessRegistration);
    }
    if (message.otpCode !== "") {
      writer.uint32(74).string(message.otpCode);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportHumanUser {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportHumanUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userName = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.profile = ExportHumanUser_Profile.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.email = ExportHumanUser_Email.decode(reader, reader.uint32());
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.phone = ExportHumanUser_Phone.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.password = reader.string();
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.hashedPassword = ExportHumanUser_HashedPassword.decode(reader, reader.uint32());
          continue;
        case 7:
          if (tag != 56) {
            break;
          }

          message.passwordChangeRequired = reader.bool();
          continue;
        case 8:
          if (tag != 64) {
            break;
          }

          message.requestPasswordlessRegistration = reader.bool();
          continue;
        case 9:
          if (tag != 74) {
            break;
          }

          message.otpCode = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExportHumanUser {
    return {
      userName: isSet(object.userName) ? String(object.userName) : "",
      profile: isSet(object.profile) ? ExportHumanUser_Profile.fromJSON(object.profile) : undefined,
      email: isSet(object.email) ? ExportHumanUser_Email.fromJSON(object.email) : undefined,
      phone: isSet(object.phone) ? ExportHumanUser_Phone.fromJSON(object.phone) : undefined,
      password: isSet(object.password) ? String(object.password) : "",
      hashedPassword: isSet(object.hashedPassword)
        ? ExportHumanUser_HashedPassword.fromJSON(object.hashedPassword)
        : undefined,
      passwordChangeRequired: isSet(object.passwordChangeRequired) ? Boolean(object.passwordChangeRequired) : false,
      requestPasswordlessRegistration: isSet(object.requestPasswordlessRegistration)
        ? Boolean(object.requestPasswordlessRegistration)
        : false,
      otpCode: isSet(object.otpCode) ? String(object.otpCode) : "",
    };
  },

  toJSON(message: ExportHumanUser): unknown {
    const obj: any = {};
    message.userName !== undefined && (obj.userName = message.userName);
    message.profile !== undefined &&
      (obj.profile = message.profile ? ExportHumanUser_Profile.toJSON(message.profile) : undefined);
    message.email !== undefined &&
      (obj.email = message.email ? ExportHumanUser_Email.toJSON(message.email) : undefined);
    message.phone !== undefined &&
      (obj.phone = message.phone ? ExportHumanUser_Phone.toJSON(message.phone) : undefined);
    message.password !== undefined && (obj.password = message.password);
    message.hashedPassword !== undefined && (obj.hashedPassword = message.hashedPassword
      ? ExportHumanUser_HashedPassword.toJSON(message.hashedPassword)
      : undefined);
    message.passwordChangeRequired !== undefined && (obj.passwordChangeRequired = message.passwordChangeRequired);
    message.requestPasswordlessRegistration !== undefined &&
      (obj.requestPasswordlessRegistration = message.requestPasswordlessRegistration);
    message.otpCode !== undefined && (obj.otpCode = message.otpCode);
    return obj;
  },

  create(base?: DeepPartial<ExportHumanUser>): ExportHumanUser {
    return ExportHumanUser.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportHumanUser>): ExportHumanUser {
    const message = createBaseExportHumanUser();
    message.userName = object.userName ?? "";
    message.profile = (object.profile !== undefined && object.profile !== null)
      ? ExportHumanUser_Profile.fromPartial(object.profile)
      : undefined;
    message.email = (object.email !== undefined && object.email !== null)
      ? ExportHumanUser_Email.fromPartial(object.email)
      : undefined;
    message.phone = (object.phone !== undefined && object.phone !== null)
      ? ExportHumanUser_Phone.fromPartial(object.phone)
      : undefined;
    message.password = object.password ?? "";
    message.hashedPassword = (object.hashedPassword !== undefined && object.hashedPassword !== null)
      ? ExportHumanUser_HashedPassword.fromPartial(object.hashedPassword)
      : undefined;
    message.passwordChangeRequired = object.passwordChangeRequired ?? false;
    message.requestPasswordlessRegistration = object.requestPasswordlessRegistration ?? false;
    message.otpCode = object.otpCode ?? "";
    return message;
  },
};

function createBaseExportHumanUser_Profile(): ExportHumanUser_Profile {
  return { firstName: "", lastName: "", nickName: "", displayName: "", preferredLanguage: "", gender: 0 };
}

export const ExportHumanUser_Profile = {
  encode(message: ExportHumanUser_Profile, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.firstName !== "") {
      writer.uint32(10).string(message.firstName);
    }
    if (message.lastName !== "") {
      writer.uint32(18).string(message.lastName);
    }
    if (message.nickName !== "") {
      writer.uint32(26).string(message.nickName);
    }
    if (message.displayName !== "") {
      writer.uint32(34).string(message.displayName);
    }
    if (message.preferredLanguage !== "") {
      writer.uint32(42).string(message.preferredLanguage);
    }
    if (message.gender !== 0) {
      writer.uint32(48).int32(message.gender);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportHumanUser_Profile {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportHumanUser_Profile();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.firstName = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.lastName = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.nickName = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.displayName = reader.string();
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.preferredLanguage = reader.string();
          continue;
        case 6:
          if (tag != 48) {
            break;
          }

          message.gender = reader.int32() as any;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExportHumanUser_Profile {
    return {
      firstName: isSet(object.firstName) ? String(object.firstName) : "",
      lastName: isSet(object.lastName) ? String(object.lastName) : "",
      nickName: isSet(object.nickName) ? String(object.nickName) : "",
      displayName: isSet(object.displayName) ? String(object.displayName) : "",
      preferredLanguage: isSet(object.preferredLanguage) ? String(object.preferredLanguage) : "",
      gender: isSet(object.gender) ? genderFromJSON(object.gender) : 0,
    };
  },

  toJSON(message: ExportHumanUser_Profile): unknown {
    const obj: any = {};
    message.firstName !== undefined && (obj.firstName = message.firstName);
    message.lastName !== undefined && (obj.lastName = message.lastName);
    message.nickName !== undefined && (obj.nickName = message.nickName);
    message.displayName !== undefined && (obj.displayName = message.displayName);
    message.preferredLanguage !== undefined && (obj.preferredLanguage = message.preferredLanguage);
    message.gender !== undefined && (obj.gender = genderToJSON(message.gender));
    return obj;
  },

  create(base?: DeepPartial<ExportHumanUser_Profile>): ExportHumanUser_Profile {
    return ExportHumanUser_Profile.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportHumanUser_Profile>): ExportHumanUser_Profile {
    const message = createBaseExportHumanUser_Profile();
    message.firstName = object.firstName ?? "";
    message.lastName = object.lastName ?? "";
    message.nickName = object.nickName ?? "";
    message.displayName = object.displayName ?? "";
    message.preferredLanguage = object.preferredLanguage ?? "";
    message.gender = object.gender ?? 0;
    return message;
  },
};

function createBaseExportHumanUser_Email(): ExportHumanUser_Email {
  return { email: "", isEmailVerified: false };
}

export const ExportHumanUser_Email = {
  encode(message: ExportHumanUser_Email, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    if (message.isEmailVerified === true) {
      writer.uint32(16).bool(message.isEmailVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportHumanUser_Email {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportHumanUser_Email();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.email = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.isEmailVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExportHumanUser_Email {
    return {
      email: isSet(object.email) ? String(object.email) : "",
      isEmailVerified: isSet(object.isEmailVerified) ? Boolean(object.isEmailVerified) : false,
    };
  },

  toJSON(message: ExportHumanUser_Email): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    message.isEmailVerified !== undefined && (obj.isEmailVerified = message.isEmailVerified);
    return obj;
  },

  create(base?: DeepPartial<ExportHumanUser_Email>): ExportHumanUser_Email {
    return ExportHumanUser_Email.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportHumanUser_Email>): ExportHumanUser_Email {
    const message = createBaseExportHumanUser_Email();
    message.email = object.email ?? "";
    message.isEmailVerified = object.isEmailVerified ?? false;
    return message;
  },
};

function createBaseExportHumanUser_Phone(): ExportHumanUser_Phone {
  return { phone: "", isPhoneVerified: false };
}

export const ExportHumanUser_Phone = {
  encode(message: ExportHumanUser_Phone, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.phone !== "") {
      writer.uint32(10).string(message.phone);
    }
    if (message.isPhoneVerified === true) {
      writer.uint32(16).bool(message.isPhoneVerified);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportHumanUser_Phone {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportHumanUser_Phone();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.phone = reader.string();
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.isPhoneVerified = reader.bool();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExportHumanUser_Phone {
    return {
      phone: isSet(object.phone) ? String(object.phone) : "",
      isPhoneVerified: isSet(object.isPhoneVerified) ? Boolean(object.isPhoneVerified) : false,
    };
  },

  toJSON(message: ExportHumanUser_Phone): unknown {
    const obj: any = {};
    message.phone !== undefined && (obj.phone = message.phone);
    message.isPhoneVerified !== undefined && (obj.isPhoneVerified = message.isPhoneVerified);
    return obj;
  },

  create(base?: DeepPartial<ExportHumanUser_Phone>): ExportHumanUser_Phone {
    return ExportHumanUser_Phone.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportHumanUser_Phone>): ExportHumanUser_Phone {
    const message = createBaseExportHumanUser_Phone();
    message.phone = object.phone ?? "";
    message.isPhoneVerified = object.isPhoneVerified ?? false;
    return message;
  },
};

function createBaseExportHumanUser_HashedPassword(): ExportHumanUser_HashedPassword {
  return { value: "", algorithm: "" };
}

export const ExportHumanUser_HashedPassword = {
  encode(message: ExportHumanUser_HashedPassword, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value !== "") {
      writer.uint32(10).string(message.value);
    }
    if (message.algorithm !== "") {
      writer.uint32(18).string(message.algorithm);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExportHumanUser_HashedPassword {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExportHumanUser_HashedPassword();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.value = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.algorithm = reader.string();
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExportHumanUser_HashedPassword {
    return {
      value: isSet(object.value) ? String(object.value) : "",
      algorithm: isSet(object.algorithm) ? String(object.algorithm) : "",
    };
  },

  toJSON(message: ExportHumanUser_HashedPassword): unknown {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value);
    message.algorithm !== undefined && (obj.algorithm = message.algorithm);
    return obj;
  },

  create(base?: DeepPartial<ExportHumanUser_HashedPassword>): ExportHumanUser_HashedPassword {
    return ExportHumanUser_HashedPassword.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<ExportHumanUser_HashedPassword>): ExportHumanUser_HashedPassword {
    const message = createBaseExportHumanUser_HashedPassword();
    message.value = object.value ?? "";
    message.algorithm = object.algorithm ?? "";
    return message;
  },
};

function createBaseDataAppKey(): DataAppKey {
  return {
    id: "",
    projectId: "",
    appId: "",
    clientId: "",
    type: 0,
    expirationDate: undefined,
    publicKey: Buffer.alloc(0),
  };
}

export const DataAppKey = {
  encode(message: DataAppKey, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.projectId !== "") {
      writer.uint32(18).string(message.projectId);
    }
    if (message.appId !== "") {
      writer.uint32(26).string(message.appId);
    }
    if (message.clientId !== "") {
      writer.uint32(34).string(message.clientId);
    }
    if (message.type !== 0) {
      writer.uint32(40).int32(message.type);
    }
    if (message.expirationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.expirationDate), writer.uint32(50).fork()).ldelim();
    }
    if (message.publicKey.length !== 0) {
      writer.uint32(58).bytes(message.publicKey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataAppKey {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataAppKey();
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

          message.projectId = reader.string();
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.appId = reader.string();
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.clientId = reader.string();
          continue;
        case 5:
          if (tag != 40) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 6:
          if (tag != 50) {
            break;
          }

          message.expirationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag != 58) {
            break;
          }

          message.publicKey = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataAppKey {
    return {
      id: isSet(object.id) ? String(object.id) : "",
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      appId: isSet(object.appId) ? String(object.appId) : "",
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      type: isSet(object.type) ? keyTypeFromJSON(object.type) : 0,
      expirationDate: isSet(object.expirationDate) ? fromJsonTimestamp(object.expirationDate) : undefined,
      publicKey: isSet(object.publicKey) ? Buffer.from(bytesFromBase64(object.publicKey)) : Buffer.alloc(0),
    };
  },

  toJSON(message: DataAppKey): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.appId !== undefined && (obj.appId = message.appId);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.type !== undefined && (obj.type = keyTypeToJSON(message.type));
    message.expirationDate !== undefined && (obj.expirationDate = message.expirationDate.toISOString());
    message.publicKey !== undefined &&
      (obj.publicKey = base64FromBytes(message.publicKey !== undefined ? message.publicKey : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<DataAppKey>): DataAppKey {
    return DataAppKey.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataAppKey>): DataAppKey {
    const message = createBaseDataAppKey();
    message.id = object.id ?? "";
    message.projectId = object.projectId ?? "";
    message.appId = object.appId ?? "";
    message.clientId = object.clientId ?? "";
    message.type = object.type ?? 0;
    message.expirationDate = object.expirationDate ?? undefined;
    message.publicKey = object.publicKey ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseDataMachineKey(): DataMachineKey {
  return { keyId: "", userId: "", type: 0, expirationDate: undefined, publicKey: Buffer.alloc(0) };
}

export const DataMachineKey = {
  encode(message: DataMachineKey, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyId !== "") {
      writer.uint32(10).string(message.keyId);
    }
    if (message.userId !== "") {
      writer.uint32(18).string(message.userId);
    }
    if (message.type !== 0) {
      writer.uint32(24).int32(message.type);
    }
    if (message.expirationDate !== undefined) {
      Timestamp.encode(toTimestamp(message.expirationDate), writer.uint32(34).fork()).ldelim();
    }
    if (message.publicKey.length !== 0) {
      writer.uint32(42).bytes(message.publicKey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataMachineKey {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataMachineKey();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.keyId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 3:
          if (tag != 24) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
        case 4:
          if (tag != 34) {
            break;
          }

          message.expirationDate = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag != 42) {
            break;
          }

          message.publicKey = reader.bytes() as Buffer;
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataMachineKey {
    return {
      keyId: isSet(object.keyId) ? String(object.keyId) : "",
      userId: isSet(object.userId) ? String(object.userId) : "",
      type: isSet(object.type) ? keyTypeFromJSON(object.type) : 0,
      expirationDate: isSet(object.expirationDate) ? fromJsonTimestamp(object.expirationDate) : undefined,
      publicKey: isSet(object.publicKey) ? Buffer.from(bytesFromBase64(object.publicKey)) : Buffer.alloc(0),
    };
  },

  toJSON(message: DataMachineKey): unknown {
    const obj: any = {};
    message.keyId !== undefined && (obj.keyId = message.keyId);
    message.userId !== undefined && (obj.userId = message.userId);
    message.type !== undefined && (obj.type = keyTypeToJSON(message.type));
    message.expirationDate !== undefined && (obj.expirationDate = message.expirationDate.toISOString());
    message.publicKey !== undefined &&
      (obj.publicKey = base64FromBytes(message.publicKey !== undefined ? message.publicKey : Buffer.alloc(0)));
    return obj;
  },

  create(base?: DeepPartial<DataMachineKey>): DataMachineKey {
    return DataMachineKey.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataMachineKey>): DataMachineKey {
    const message = createBaseDataMachineKey();
    message.keyId = object.keyId ?? "";
    message.userId = object.userId ?? "";
    message.type = object.type ?? 0;
    message.expirationDate = object.expirationDate ?? undefined;
    message.publicKey = object.publicKey ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseDataProject(): DataProject {
  return { projectId: "", project: undefined };
}

export const DataProject = {
  encode(message: DataProject, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.projectId !== "") {
      writer.uint32(10).string(message.projectId);
    }
    if (message.project !== undefined) {
      AddProjectRequest.encode(message.project, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataProject {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataProject();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.projectId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.project = AddProjectRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataProject {
    return {
      projectId: isSet(object.projectId) ? String(object.projectId) : "",
      project: isSet(object.project) ? AddProjectRequest.fromJSON(object.project) : undefined,
    };
  },

  toJSON(message: DataProject): unknown {
    const obj: any = {};
    message.projectId !== undefined && (obj.projectId = message.projectId);
    message.project !== undefined &&
      (obj.project = message.project ? AddProjectRequest.toJSON(message.project) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataProject>): DataProject {
    return DataProject.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataProject>): DataProject {
    const message = createBaseDataProject();
    message.projectId = object.projectId ?? "";
    message.project = (object.project !== undefined && object.project !== null)
      ? AddProjectRequest.fromPartial(object.project)
      : undefined;
    return message;
  },
};

function createBaseDataAPIApplication(): DataAPIApplication {
  return { appId: "", app: undefined };
}

export const DataAPIApplication = {
  encode(message: DataAPIApplication, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.appId !== "") {
      writer.uint32(10).string(message.appId);
    }
    if (message.app !== undefined) {
      AddAPIAppRequest.encode(message.app, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataAPIApplication {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataAPIApplication();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.appId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.app = AddAPIAppRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataAPIApplication {
    return {
      appId: isSet(object.appId) ? String(object.appId) : "",
      app: isSet(object.app) ? AddAPIAppRequest.fromJSON(object.app) : undefined,
    };
  },

  toJSON(message: DataAPIApplication): unknown {
    const obj: any = {};
    message.appId !== undefined && (obj.appId = message.appId);
    message.app !== undefined && (obj.app = message.app ? AddAPIAppRequest.toJSON(message.app) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataAPIApplication>): DataAPIApplication {
    return DataAPIApplication.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataAPIApplication>): DataAPIApplication {
    const message = createBaseDataAPIApplication();
    message.appId = object.appId ?? "";
    message.app = (object.app !== undefined && object.app !== null)
      ? AddAPIAppRequest.fromPartial(object.app)
      : undefined;
    return message;
  },
};

function createBaseDataOIDCApplication(): DataOIDCApplication {
  return { appId: "", app: undefined };
}

export const DataOIDCApplication = {
  encode(message: DataOIDCApplication, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.appId !== "") {
      writer.uint32(10).string(message.appId);
    }
    if (message.app !== undefined) {
      AddOIDCAppRequest.encode(message.app, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataOIDCApplication {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataOIDCApplication();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.appId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.app = AddOIDCAppRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataOIDCApplication {
    return {
      appId: isSet(object.appId) ? String(object.appId) : "",
      app: isSet(object.app) ? AddOIDCAppRequest.fromJSON(object.app) : undefined,
    };
  },

  toJSON(message: DataOIDCApplication): unknown {
    const obj: any = {};
    message.appId !== undefined && (obj.appId = message.appId);
    message.app !== undefined && (obj.app = message.app ? AddOIDCAppRequest.toJSON(message.app) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataOIDCApplication>): DataOIDCApplication {
    return DataOIDCApplication.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataOIDCApplication>): DataOIDCApplication {
    const message = createBaseDataOIDCApplication();
    message.appId = object.appId ?? "";
    message.app = (object.app !== undefined && object.app !== null)
      ? AddOIDCAppRequest.fromPartial(object.app)
      : undefined;
    return message;
  },
};

function createBaseDataHumanUser(): DataHumanUser {
  return { userId: "", user: undefined };
}

export const DataHumanUser = {
  encode(message: DataHumanUser, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.user !== undefined) {
      ImportHumanUserRequest.encode(message.user, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataHumanUser {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataHumanUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.user = ImportHumanUserRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataHumanUser {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      user: isSet(object.user) ? ImportHumanUserRequest.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: DataHumanUser): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.user !== undefined && (obj.user = message.user ? ImportHumanUserRequest.toJSON(message.user) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataHumanUser>): DataHumanUser {
    return DataHumanUser.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataHumanUser>): DataHumanUser {
    const message = createBaseDataHumanUser();
    message.userId = object.userId ?? "";
    message.user = (object.user !== undefined && object.user !== null)
      ? ImportHumanUserRequest.fromPartial(object.user)
      : undefined;
    return message;
  },
};

function createBaseDataMachineUser(): DataMachineUser {
  return { userId: "", user: undefined };
}

export const DataMachineUser = {
  encode(message: DataMachineUser, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    if (message.user !== undefined) {
      AddMachineUserRequest.encode(message.user, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataMachineUser {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataMachineUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.userId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.user = AddMachineUserRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataMachineUser {
    return {
      userId: isSet(object.userId) ? String(object.userId) : "",
      user: isSet(object.user) ? AddMachineUserRequest.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: DataMachineUser): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    message.user !== undefined && (obj.user = message.user ? AddMachineUserRequest.toJSON(message.user) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataMachineUser>): DataMachineUser {
    return DataMachineUser.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataMachineUser>): DataMachineUser {
    const message = createBaseDataMachineUser();
    message.userId = object.userId ?? "";
    message.user = (object.user !== undefined && object.user !== null)
      ? AddMachineUserRequest.fromPartial(object.user)
      : undefined;
    return message;
  },
};

function createBaseDataAction(): DataAction {
  return { actionId: "", action: undefined };
}

export const DataAction = {
  encode(message: DataAction, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.actionId !== "") {
      writer.uint32(10).string(message.actionId);
    }
    if (message.action !== undefined) {
      CreateActionRequest.encode(message.action, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataAction {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataAction();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 10) {
            break;
          }

          message.actionId = reader.string();
          continue;
        case 2:
          if (tag != 18) {
            break;
          }

          message.action = CreateActionRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataAction {
    return {
      actionId: isSet(object.actionId) ? String(object.actionId) : "",
      action: isSet(object.action) ? CreateActionRequest.fromJSON(object.action) : undefined,
    };
  },

  toJSON(message: DataAction): unknown {
    const obj: any = {};
    message.actionId !== undefined && (obj.actionId = message.actionId);
    message.action !== undefined &&
      (obj.action = message.action ? CreateActionRequest.toJSON(message.action) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataAction>): DataAction {
    return DataAction.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataAction>): DataAction {
    const message = createBaseDataAction();
    message.actionId = object.actionId ?? "";
    message.action = (object.action !== undefined && object.action !== null)
      ? CreateActionRequest.fromPartial(object.action)
      : undefined;
    return message;
  },
};

function createBaseDataProjectGrant(): DataProjectGrant {
  return { grantId: "", projectGrant: undefined };
}

export const DataProjectGrant = {
  encode(message: DataProjectGrant, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.grantId !== "") {
      writer.uint32(10).string(message.grantId);
    }
    if (message.projectGrant !== undefined) {
      AddProjectGrantRequest.encode(message.projectGrant, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DataProjectGrant {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDataProjectGrant();
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

          message.projectGrant = AddProjectGrantRequest.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DataProjectGrant {
    return {
      grantId: isSet(object.grantId) ? String(object.grantId) : "",
      projectGrant: isSet(object.projectGrant) ? AddProjectGrantRequest.fromJSON(object.projectGrant) : undefined,
    };
  },

  toJSON(message: DataProjectGrant): unknown {
    const obj: any = {};
    message.grantId !== undefined && (obj.grantId = message.grantId);
    message.projectGrant !== undefined &&
      (obj.projectGrant = message.projectGrant ? AddProjectGrantRequest.toJSON(message.projectGrant) : undefined);
    return obj;
  },

  create(base?: DeepPartial<DataProjectGrant>): DataProjectGrant {
    return DataProjectGrant.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<DataProjectGrant>): DataProjectGrant {
    const message = createBaseDataProjectGrant();
    message.grantId = object.grantId ?? "";
    message.projectGrant = (object.projectGrant !== undefined && object.projectGrant !== null)
      ? AddProjectGrantRequest.fromPartial(object.projectGrant)
      : undefined;
    return message;
  },
};

function createBaseSetTriggerActionsRequest(): SetTriggerActionsRequest {
  return { flowType: 0, triggerType: 0, actionIds: [] };
}

export const SetTriggerActionsRequest = {
  encode(message: SetTriggerActionsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.flowType !== 0) {
      writer.uint32(8).int32(message.flowType);
    }
    if (message.triggerType !== 0) {
      writer.uint32(16).int32(message.triggerType);
    }
    for (const v of message.actionIds) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetTriggerActionsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetTriggerActionsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag != 8) {
            break;
          }

          message.flowType = reader.int32() as any;
          continue;
        case 2:
          if (tag != 16) {
            break;
          }

          message.triggerType = reader.int32() as any;
          continue;
        case 3:
          if (tag != 26) {
            break;
          }

          message.actionIds.push(reader.string());
          continue;
      }
      if ((tag & 7) == 4 || tag == 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetTriggerActionsRequest {
    return {
      flowType: isSet(object.flowType) ? flowTypeFromJSON(object.flowType) : 0,
      triggerType: isSet(object.triggerType) ? triggerTypeFromJSON(object.triggerType) : 0,
      actionIds: Array.isArray(object?.actionIds) ? object.actionIds.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: SetTriggerActionsRequest): unknown {
    const obj: any = {};
    message.flowType !== undefined && (obj.flowType = flowTypeToJSON(message.flowType));
    message.triggerType !== undefined && (obj.triggerType = triggerTypeToJSON(message.triggerType));
    if (message.actionIds) {
      obj.actionIds = message.actionIds.map((e) => e);
    } else {
      obj.actionIds = [];
    }
    return obj;
  },

  create(base?: DeepPartial<SetTriggerActionsRequest>): SetTriggerActionsRequest {
    return SetTriggerActionsRequest.fromPartial(base ?? {});
  },

  fromPartial(object: DeepPartial<SetTriggerActionsRequest>): SetTriggerActionsRequest {
    const message = createBaseSetTriggerActionsRequest();
    message.flowType = object.flowType ?? 0;
    message.triggerType = object.triggerType ?? 0;
    message.actionIds = object.actionIds?.map((e) => e) || [];
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
