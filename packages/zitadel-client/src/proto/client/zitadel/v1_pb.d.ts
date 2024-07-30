import * as jspb from 'google-protobuf'

import * as zitadel_user_pb from '../zitadel/user_pb'; // proto import: "zitadel/user.proto"
import * as zitadel_idp_pb from '../zitadel/idp_pb'; // proto import: "zitadel/idp.proto"
import * as zitadel_org_pb from '../zitadel/org_pb'; // proto import: "zitadel/org.proto"
import * as zitadel_management_pb from '../zitadel/management_pb'; // proto import: "zitadel/management.proto"
import * as zitadel_auth_n_key_pb from '../zitadel/auth_n_key_pb'; // proto import: "zitadel/auth_n_key.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from '../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb'; // proto import: "google/protobuf/timestamp.proto"
import * as validate_validate_pb from '../validate/validate_pb'; // proto import: "validate/validate.proto"


export class AddCustomOrgIAMPolicyRequest extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): AddCustomOrgIAMPolicyRequest;

  getUserLoginMustBeDomain(): boolean;
  setUserLoginMustBeDomain(value: boolean): AddCustomOrgIAMPolicyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddCustomOrgIAMPolicyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddCustomOrgIAMPolicyRequest): AddCustomOrgIAMPolicyRequest.AsObject;
  static serializeBinaryToWriter(message: AddCustomOrgIAMPolicyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddCustomOrgIAMPolicyRequest;
  static deserializeBinaryFromReader(message: AddCustomOrgIAMPolicyRequest, reader: jspb.BinaryReader): AddCustomOrgIAMPolicyRequest;
}

export namespace AddCustomOrgIAMPolicyRequest {
  export type AsObject = {
    orgId: string,
    userLoginMustBeDomain: boolean,
  }
}

export class ImportDataOrg extends jspb.Message {
  getOrgsList(): Array<DataOrg>;
  setOrgsList(value: Array<DataOrg>): ImportDataOrg;
  clearOrgsList(): ImportDataOrg;
  addOrgs(value?: DataOrg, index?: number): DataOrg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ImportDataOrg.AsObject;
  static toObject(includeInstance: boolean, msg: ImportDataOrg): ImportDataOrg.AsObject;
  static serializeBinaryToWriter(message: ImportDataOrg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ImportDataOrg;
  static deserializeBinaryFromReader(message: ImportDataOrg, reader: jspb.BinaryReader): ImportDataOrg;
}

export namespace ImportDataOrg {
  export type AsObject = {
    orgsList: Array<DataOrg.AsObject>,
  }
}

export class DataOrg extends jspb.Message {
  getOrgId(): string;
  setOrgId(value: string): DataOrg;

  getOrg(): zitadel_management_pb.AddOrgRequest | undefined;
  setOrg(value?: zitadel_management_pb.AddOrgRequest): DataOrg;
  hasOrg(): boolean;
  clearOrg(): DataOrg;

  getIamPolicy(): AddCustomOrgIAMPolicyRequest | undefined;
  setIamPolicy(value?: AddCustomOrgIAMPolicyRequest): DataOrg;
  hasIamPolicy(): boolean;
  clearIamPolicy(): DataOrg;

  getLabelPolicy(): zitadel_management_pb.AddCustomLabelPolicyRequest | undefined;
  setLabelPolicy(value?: zitadel_management_pb.AddCustomLabelPolicyRequest): DataOrg;
  hasLabelPolicy(): boolean;
  clearLabelPolicy(): DataOrg;

  getLockoutPolicy(): zitadel_management_pb.AddCustomLockoutPolicyRequest | undefined;
  setLockoutPolicy(value?: zitadel_management_pb.AddCustomLockoutPolicyRequest): DataOrg;
  hasLockoutPolicy(): boolean;
  clearLockoutPolicy(): DataOrg;

  getLoginPolicy(): zitadel_management_pb.AddCustomLoginPolicyRequest | undefined;
  setLoginPolicy(value?: zitadel_management_pb.AddCustomLoginPolicyRequest): DataOrg;
  hasLoginPolicy(): boolean;
  clearLoginPolicy(): DataOrg;

  getPasswordComplexityPolicy(): zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest | undefined;
  setPasswordComplexityPolicy(value?: zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest): DataOrg;
  hasPasswordComplexityPolicy(): boolean;
  clearPasswordComplexityPolicy(): DataOrg;

  getPrivacyPolicy(): zitadel_management_pb.AddCustomPrivacyPolicyRequest | undefined;
  setPrivacyPolicy(value?: zitadel_management_pb.AddCustomPrivacyPolicyRequest): DataOrg;
  hasPrivacyPolicy(): boolean;
  clearPrivacyPolicy(): DataOrg;

  getProjectsList(): Array<DataProject>;
  setProjectsList(value: Array<DataProject>): DataOrg;
  clearProjectsList(): DataOrg;
  addProjects(value?: DataProject, index?: number): DataProject;

  getProjectRolesList(): Array<zitadel_management_pb.AddProjectRoleRequest>;
  setProjectRolesList(value: Array<zitadel_management_pb.AddProjectRoleRequest>): DataOrg;
  clearProjectRolesList(): DataOrg;
  addProjectRoles(value?: zitadel_management_pb.AddProjectRoleRequest, index?: number): zitadel_management_pb.AddProjectRoleRequest;

  getApiAppsList(): Array<DataAPIApplication>;
  setApiAppsList(value: Array<DataAPIApplication>): DataOrg;
  clearApiAppsList(): DataOrg;
  addApiApps(value?: DataAPIApplication, index?: number): DataAPIApplication;

  getOidcAppsList(): Array<DataOIDCApplication>;
  setOidcAppsList(value: Array<DataOIDCApplication>): DataOrg;
  clearOidcAppsList(): DataOrg;
  addOidcApps(value?: DataOIDCApplication, index?: number): DataOIDCApplication;

  getHumanUsersList(): Array<DataHumanUser>;
  setHumanUsersList(value: Array<DataHumanUser>): DataOrg;
  clearHumanUsersList(): DataOrg;
  addHumanUsers(value?: DataHumanUser, index?: number): DataHumanUser;

  getMachineUsersList(): Array<DataMachineUser>;
  setMachineUsersList(value: Array<DataMachineUser>): DataOrg;
  clearMachineUsersList(): DataOrg;
  addMachineUsers(value?: DataMachineUser, index?: number): DataMachineUser;

  getTriggerActionsList(): Array<SetTriggerActionsRequest>;
  setTriggerActionsList(value: Array<SetTriggerActionsRequest>): DataOrg;
  clearTriggerActionsList(): DataOrg;
  addTriggerActions(value?: SetTriggerActionsRequest, index?: number): SetTriggerActionsRequest;

  getActionsList(): Array<DataAction>;
  setActionsList(value: Array<DataAction>): DataOrg;
  clearActionsList(): DataOrg;
  addActions(value?: DataAction, index?: number): DataAction;

  getProjectGrantsList(): Array<DataProjectGrant>;
  setProjectGrantsList(value: Array<DataProjectGrant>): DataOrg;
  clearProjectGrantsList(): DataOrg;
  addProjectGrants(value?: DataProjectGrant, index?: number): DataProjectGrant;

  getUserGrantsList(): Array<zitadel_management_pb.AddUserGrantRequest>;
  setUserGrantsList(value: Array<zitadel_management_pb.AddUserGrantRequest>): DataOrg;
  clearUserGrantsList(): DataOrg;
  addUserGrants(value?: zitadel_management_pb.AddUserGrantRequest, index?: number): zitadel_management_pb.AddUserGrantRequest;

  getOrgMembersList(): Array<zitadel_management_pb.AddOrgMemberRequest>;
  setOrgMembersList(value: Array<zitadel_management_pb.AddOrgMemberRequest>): DataOrg;
  clearOrgMembersList(): DataOrg;
  addOrgMembers(value?: zitadel_management_pb.AddOrgMemberRequest, index?: number): zitadel_management_pb.AddOrgMemberRequest;

  getProjectMembersList(): Array<zitadel_management_pb.AddProjectMemberRequest>;
  setProjectMembersList(value: Array<zitadel_management_pb.AddProjectMemberRequest>): DataOrg;
  clearProjectMembersList(): DataOrg;
  addProjectMembers(value?: zitadel_management_pb.AddProjectMemberRequest, index?: number): zitadel_management_pb.AddProjectMemberRequest;

  getProjectGrantMembersList(): Array<zitadel_management_pb.AddProjectGrantMemberRequest>;
  setProjectGrantMembersList(value: Array<zitadel_management_pb.AddProjectGrantMemberRequest>): DataOrg;
  clearProjectGrantMembersList(): DataOrg;
  addProjectGrantMembers(value?: zitadel_management_pb.AddProjectGrantMemberRequest, index?: number): zitadel_management_pb.AddProjectGrantMemberRequest;

  getUserMetadataList(): Array<zitadel_management_pb.SetUserMetadataRequest>;
  setUserMetadataList(value: Array<zitadel_management_pb.SetUserMetadataRequest>): DataOrg;
  clearUserMetadataList(): DataOrg;
  addUserMetadata(value?: zitadel_management_pb.SetUserMetadataRequest, index?: number): zitadel_management_pb.SetUserMetadataRequest;

  getLoginTextsList(): Array<zitadel_management_pb.SetCustomLoginTextsRequest>;
  setLoginTextsList(value: Array<zitadel_management_pb.SetCustomLoginTextsRequest>): DataOrg;
  clearLoginTextsList(): DataOrg;
  addLoginTexts(value?: zitadel_management_pb.SetCustomLoginTextsRequest, index?: number): zitadel_management_pb.SetCustomLoginTextsRequest;

  getInitMessagesList(): Array<zitadel_management_pb.SetCustomInitMessageTextRequest>;
  setInitMessagesList(value: Array<zitadel_management_pb.SetCustomInitMessageTextRequest>): DataOrg;
  clearInitMessagesList(): DataOrg;
  addInitMessages(value?: zitadel_management_pb.SetCustomInitMessageTextRequest, index?: number): zitadel_management_pb.SetCustomInitMessageTextRequest;

  getPasswordResetMessagesList(): Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest>;
  setPasswordResetMessagesList(value: Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest>): DataOrg;
  clearPasswordResetMessagesList(): DataOrg;
  addPasswordResetMessages(value?: zitadel_management_pb.SetCustomPasswordResetMessageTextRequest, index?: number): zitadel_management_pb.SetCustomPasswordResetMessageTextRequest;

  getVerifyEmailMessagesList(): Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest>;
  setVerifyEmailMessagesList(value: Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest>): DataOrg;
  clearVerifyEmailMessagesList(): DataOrg;
  addVerifyEmailMessages(value?: zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest;

  getVerifyPhoneMessagesList(): Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest>;
  setVerifyPhoneMessagesList(value: Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest>): DataOrg;
  clearVerifyPhoneMessagesList(): DataOrg;
  addVerifyPhoneMessages(value?: zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest, index?: number): zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest;

  getDomainClaimedMessagesList(): Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest>;
  setDomainClaimedMessagesList(value: Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest>): DataOrg;
  clearDomainClaimedMessagesList(): DataOrg;
  addDomainClaimedMessages(value?: zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest, index?: number): zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest;

  getPasswordlessRegistrationMessagesList(): Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest>;
  setPasswordlessRegistrationMessagesList(value: Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest>): DataOrg;
  clearPasswordlessRegistrationMessagesList(): DataOrg;
  addPasswordlessRegistrationMessages(value?: zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest, index?: number): zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest;

  getOidcIdpsList(): Array<DataOIDCIDP>;
  setOidcIdpsList(value: Array<DataOIDCIDP>): DataOrg;
  clearOidcIdpsList(): DataOrg;
  addOidcIdps(value?: DataOIDCIDP, index?: number): DataOIDCIDP;

  getJwtIdpsList(): Array<DataJWTIDP>;
  setJwtIdpsList(value: Array<DataJWTIDP>): DataOrg;
  clearJwtIdpsList(): DataOrg;
  addJwtIdps(value?: DataJWTIDP, index?: number): DataJWTIDP;

  getSecondFactorsList(): Array<zitadel_management_pb.AddSecondFactorToLoginPolicyRequest>;
  setSecondFactorsList(value: Array<zitadel_management_pb.AddSecondFactorToLoginPolicyRequest>): DataOrg;
  clearSecondFactorsList(): DataOrg;
  addSecondFactors(value?: zitadel_management_pb.AddSecondFactorToLoginPolicyRequest, index?: number): zitadel_management_pb.AddSecondFactorToLoginPolicyRequest;

  getMultiFactorsList(): Array<zitadel_management_pb.AddMultiFactorToLoginPolicyRequest>;
  setMultiFactorsList(value: Array<zitadel_management_pb.AddMultiFactorToLoginPolicyRequest>): DataOrg;
  clearMultiFactorsList(): DataOrg;
  addMultiFactors(value?: zitadel_management_pb.AddMultiFactorToLoginPolicyRequest, index?: number): zitadel_management_pb.AddMultiFactorToLoginPolicyRequest;

  getIdpsList(): Array<zitadel_management_pb.AddIDPToLoginPolicyRequest>;
  setIdpsList(value: Array<zitadel_management_pb.AddIDPToLoginPolicyRequest>): DataOrg;
  clearIdpsList(): DataOrg;
  addIdps(value?: zitadel_management_pb.AddIDPToLoginPolicyRequest, index?: number): zitadel_management_pb.AddIDPToLoginPolicyRequest;

  getUserLinksList(): Array<zitadel_idp_pb.IDPUserLink>;
  setUserLinksList(value: Array<zitadel_idp_pb.IDPUserLink>): DataOrg;
  clearUserLinksList(): DataOrg;
  addUserLinks(value?: zitadel_idp_pb.IDPUserLink, index?: number): zitadel_idp_pb.IDPUserLink;

  getDomainsList(): Array<zitadel_org_pb.Domain>;
  setDomainsList(value: Array<zitadel_org_pb.Domain>): DataOrg;
  clearDomainsList(): DataOrg;
  addDomains(value?: zitadel_org_pb.Domain, index?: number): zitadel_org_pb.Domain;

  getAppKeysList(): Array<DataAppKey>;
  setAppKeysList(value: Array<DataAppKey>): DataOrg;
  clearAppKeysList(): DataOrg;
  addAppKeys(value?: DataAppKey, index?: number): DataAppKey;

  getMachineKeysList(): Array<DataMachineKey>;
  setMachineKeysList(value: Array<DataMachineKey>): DataOrg;
  clearMachineKeysList(): DataOrg;
  addMachineKeys(value?: DataMachineKey, index?: number): DataMachineKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataOrg.AsObject;
  static toObject(includeInstance: boolean, msg: DataOrg): DataOrg.AsObject;
  static serializeBinaryToWriter(message: DataOrg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataOrg;
  static deserializeBinaryFromReader(message: DataOrg, reader: jspb.BinaryReader): DataOrg;
}

export namespace DataOrg {
  export type AsObject = {
    orgId: string,
    org?: zitadel_management_pb.AddOrgRequest.AsObject,
    iamPolicy?: AddCustomOrgIAMPolicyRequest.AsObject,
    labelPolicy?: zitadel_management_pb.AddCustomLabelPolicyRequest.AsObject,
    lockoutPolicy?: zitadel_management_pb.AddCustomLockoutPolicyRequest.AsObject,
    loginPolicy?: zitadel_management_pb.AddCustomLoginPolicyRequest.AsObject,
    passwordComplexityPolicy?: zitadel_management_pb.AddCustomPasswordComplexityPolicyRequest.AsObject,
    privacyPolicy?: zitadel_management_pb.AddCustomPrivacyPolicyRequest.AsObject,
    projectsList: Array<DataProject.AsObject>,
    projectRolesList: Array<zitadel_management_pb.AddProjectRoleRequest.AsObject>,
    apiAppsList: Array<DataAPIApplication.AsObject>,
    oidcAppsList: Array<DataOIDCApplication.AsObject>,
    humanUsersList: Array<DataHumanUser.AsObject>,
    machineUsersList: Array<DataMachineUser.AsObject>,
    triggerActionsList: Array<SetTriggerActionsRequest.AsObject>,
    actionsList: Array<DataAction.AsObject>,
    projectGrantsList: Array<DataProjectGrant.AsObject>,
    userGrantsList: Array<zitadel_management_pb.AddUserGrantRequest.AsObject>,
    orgMembersList: Array<zitadel_management_pb.AddOrgMemberRequest.AsObject>,
    projectMembersList: Array<zitadel_management_pb.AddProjectMemberRequest.AsObject>,
    projectGrantMembersList: Array<zitadel_management_pb.AddProjectGrantMemberRequest.AsObject>,
    userMetadataList: Array<zitadel_management_pb.SetUserMetadataRequest.AsObject>,
    loginTextsList: Array<zitadel_management_pb.SetCustomLoginTextsRequest.AsObject>,
    initMessagesList: Array<zitadel_management_pb.SetCustomInitMessageTextRequest.AsObject>,
    passwordResetMessagesList: Array<zitadel_management_pb.SetCustomPasswordResetMessageTextRequest.AsObject>,
    verifyEmailMessagesList: Array<zitadel_management_pb.SetCustomVerifyEmailMessageTextRequest.AsObject>,
    verifyPhoneMessagesList: Array<zitadel_management_pb.SetCustomVerifyPhoneMessageTextRequest.AsObject>,
    domainClaimedMessagesList: Array<zitadel_management_pb.SetCustomDomainClaimedMessageTextRequest.AsObject>,
    passwordlessRegistrationMessagesList: Array<zitadel_management_pb.SetCustomPasswordlessRegistrationMessageTextRequest.AsObject>,
    oidcIdpsList: Array<DataOIDCIDP.AsObject>,
    jwtIdpsList: Array<DataJWTIDP.AsObject>,
    secondFactorsList: Array<zitadel_management_pb.AddSecondFactorToLoginPolicyRequest.AsObject>,
    multiFactorsList: Array<zitadel_management_pb.AddMultiFactorToLoginPolicyRequest.AsObject>,
    idpsList: Array<zitadel_management_pb.AddIDPToLoginPolicyRequest.AsObject>,
    userLinksList: Array<zitadel_idp_pb.IDPUserLink.AsObject>,
    domainsList: Array<zitadel_org_pb.Domain.AsObject>,
    appKeysList: Array<DataAppKey.AsObject>,
    machineKeysList: Array<DataMachineKey.AsObject>,
  }
}

export class DataOIDCIDP extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): DataOIDCIDP;

  getIdp(): zitadel_management_pb.AddOrgOIDCIDPRequest | undefined;
  setIdp(value?: zitadel_management_pb.AddOrgOIDCIDPRequest): DataOIDCIDP;
  hasIdp(): boolean;
  clearIdp(): DataOIDCIDP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataOIDCIDP.AsObject;
  static toObject(includeInstance: boolean, msg: DataOIDCIDP): DataOIDCIDP.AsObject;
  static serializeBinaryToWriter(message: DataOIDCIDP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataOIDCIDP;
  static deserializeBinaryFromReader(message: DataOIDCIDP, reader: jspb.BinaryReader): DataOIDCIDP;
}

export namespace DataOIDCIDP {
  export type AsObject = {
    idpId: string,
    idp?: zitadel_management_pb.AddOrgOIDCIDPRequest.AsObject,
  }
}

export class DataJWTIDP extends jspb.Message {
  getIdpId(): string;
  setIdpId(value: string): DataJWTIDP;

  getIdp(): zitadel_management_pb.AddOrgJWTIDPRequest | undefined;
  setIdp(value?: zitadel_management_pb.AddOrgJWTIDPRequest): DataJWTIDP;
  hasIdp(): boolean;
  clearIdp(): DataJWTIDP;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataJWTIDP.AsObject;
  static toObject(includeInstance: boolean, msg: DataJWTIDP): DataJWTIDP.AsObject;
  static serializeBinaryToWriter(message: DataJWTIDP, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataJWTIDP;
  static deserializeBinaryFromReader(message: DataJWTIDP, reader: jspb.BinaryReader): DataJWTIDP;
}

export namespace DataJWTIDP {
  export type AsObject = {
    idpId: string,
    idp?: zitadel_management_pb.AddOrgJWTIDPRequest.AsObject,
  }
}

export class ExportHumanUser extends jspb.Message {
  getUserName(): string;
  setUserName(value: string): ExportHumanUser;

  getProfile(): ExportHumanUser.Profile | undefined;
  setProfile(value?: ExportHumanUser.Profile): ExportHumanUser;
  hasProfile(): boolean;
  clearProfile(): ExportHumanUser;

  getEmail(): ExportHumanUser.Email | undefined;
  setEmail(value?: ExportHumanUser.Email): ExportHumanUser;
  hasEmail(): boolean;
  clearEmail(): ExportHumanUser;

  getPhone(): ExportHumanUser.Phone | undefined;
  setPhone(value?: ExportHumanUser.Phone): ExportHumanUser;
  hasPhone(): boolean;
  clearPhone(): ExportHumanUser;

  getPassword(): string;
  setPassword(value: string): ExportHumanUser;

  getHashedPassword(): ExportHumanUser.HashedPassword | undefined;
  setHashedPassword(value?: ExportHumanUser.HashedPassword): ExportHumanUser;
  hasHashedPassword(): boolean;
  clearHashedPassword(): ExportHumanUser;

  getPasswordChangeRequired(): boolean;
  setPasswordChangeRequired(value: boolean): ExportHumanUser;

  getRequestPasswordlessRegistration(): boolean;
  setRequestPasswordlessRegistration(value: boolean): ExportHumanUser;

  getOtpCode(): string;
  setOtpCode(value: string): ExportHumanUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ExportHumanUser.AsObject;
  static toObject(includeInstance: boolean, msg: ExportHumanUser): ExportHumanUser.AsObject;
  static serializeBinaryToWriter(message: ExportHumanUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ExportHumanUser;
  static deserializeBinaryFromReader(message: ExportHumanUser, reader: jspb.BinaryReader): ExportHumanUser;
}

export namespace ExportHumanUser {
  export type AsObject = {
    userName: string,
    profile?: ExportHumanUser.Profile.AsObject,
    email?: ExportHumanUser.Email.AsObject,
    phone?: ExportHumanUser.Phone.AsObject,
    password: string,
    hashedPassword?: ExportHumanUser.HashedPassword.AsObject,
    passwordChangeRequired: boolean,
    requestPasswordlessRegistration: boolean,
    otpCode: string,
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

    getAlgorithm(): string;
    setAlgorithm(value: string): HashedPassword;

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
      algorithm: string,
    }
  }

}

export class DataAppKey extends jspb.Message {
  getId(): string;
  setId(value: string): DataAppKey;

  getProjectId(): string;
  setProjectId(value: string): DataAppKey;

  getAppId(): string;
  setAppId(value: string): DataAppKey;

  getClientId(): string;
  setClientId(value: string): DataAppKey;

  getType(): zitadel_auth_n_key_pb.KeyType;
  setType(value: zitadel_auth_n_key_pb.KeyType): DataAppKey;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): DataAppKey;
  hasExpirationDate(): boolean;
  clearExpirationDate(): DataAppKey;

  getPublicKey(): Uint8Array | string;
  getPublicKey_asU8(): Uint8Array;
  getPublicKey_asB64(): string;
  setPublicKey(value: Uint8Array | string): DataAppKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataAppKey.AsObject;
  static toObject(includeInstance: boolean, msg: DataAppKey): DataAppKey.AsObject;
  static serializeBinaryToWriter(message: DataAppKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataAppKey;
  static deserializeBinaryFromReader(message: DataAppKey, reader: jspb.BinaryReader): DataAppKey;
}

export namespace DataAppKey {
  export type AsObject = {
    id: string,
    projectId: string,
    appId: string,
    clientId: string,
    type: zitadel_auth_n_key_pb.KeyType,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    publicKey: Uint8Array | string,
  }
}

export class DataMachineKey extends jspb.Message {
  getKeyId(): string;
  setKeyId(value: string): DataMachineKey;

  getUserId(): string;
  setUserId(value: string): DataMachineKey;

  getType(): zitadel_auth_n_key_pb.KeyType;
  setType(value: zitadel_auth_n_key_pb.KeyType): DataMachineKey;

  getExpirationDate(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpirationDate(value?: google_protobuf_timestamp_pb.Timestamp): DataMachineKey;
  hasExpirationDate(): boolean;
  clearExpirationDate(): DataMachineKey;

  getPublicKey(): Uint8Array | string;
  getPublicKey_asU8(): Uint8Array;
  getPublicKey_asB64(): string;
  setPublicKey(value: Uint8Array | string): DataMachineKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataMachineKey.AsObject;
  static toObject(includeInstance: boolean, msg: DataMachineKey): DataMachineKey.AsObject;
  static serializeBinaryToWriter(message: DataMachineKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataMachineKey;
  static deserializeBinaryFromReader(message: DataMachineKey, reader: jspb.BinaryReader): DataMachineKey;
}

export namespace DataMachineKey {
  export type AsObject = {
    keyId: string,
    userId: string,
    type: zitadel_auth_n_key_pb.KeyType,
    expirationDate?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    publicKey: Uint8Array | string,
  }
}

export class DataProject extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): DataProject;

  getProject(): zitadel_management_pb.AddProjectRequest | undefined;
  setProject(value?: zitadel_management_pb.AddProjectRequest): DataProject;
  hasProject(): boolean;
  clearProject(): DataProject;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataProject.AsObject;
  static toObject(includeInstance: boolean, msg: DataProject): DataProject.AsObject;
  static serializeBinaryToWriter(message: DataProject, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataProject;
  static deserializeBinaryFromReader(message: DataProject, reader: jspb.BinaryReader): DataProject;
}

export namespace DataProject {
  export type AsObject = {
    projectId: string,
    project?: zitadel_management_pb.AddProjectRequest.AsObject,
  }
}

export class DataAPIApplication extends jspb.Message {
  getAppId(): string;
  setAppId(value: string): DataAPIApplication;

  getApp(): zitadel_management_pb.AddAPIAppRequest | undefined;
  setApp(value?: zitadel_management_pb.AddAPIAppRequest): DataAPIApplication;
  hasApp(): boolean;
  clearApp(): DataAPIApplication;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataAPIApplication.AsObject;
  static toObject(includeInstance: boolean, msg: DataAPIApplication): DataAPIApplication.AsObject;
  static serializeBinaryToWriter(message: DataAPIApplication, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataAPIApplication;
  static deserializeBinaryFromReader(message: DataAPIApplication, reader: jspb.BinaryReader): DataAPIApplication;
}

export namespace DataAPIApplication {
  export type AsObject = {
    appId: string,
    app?: zitadel_management_pb.AddAPIAppRequest.AsObject,
  }
}

export class DataOIDCApplication extends jspb.Message {
  getAppId(): string;
  setAppId(value: string): DataOIDCApplication;

  getApp(): zitadel_management_pb.AddOIDCAppRequest | undefined;
  setApp(value?: zitadel_management_pb.AddOIDCAppRequest): DataOIDCApplication;
  hasApp(): boolean;
  clearApp(): DataOIDCApplication;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataOIDCApplication.AsObject;
  static toObject(includeInstance: boolean, msg: DataOIDCApplication): DataOIDCApplication.AsObject;
  static serializeBinaryToWriter(message: DataOIDCApplication, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataOIDCApplication;
  static deserializeBinaryFromReader(message: DataOIDCApplication, reader: jspb.BinaryReader): DataOIDCApplication;
}

export namespace DataOIDCApplication {
  export type AsObject = {
    appId: string,
    app?: zitadel_management_pb.AddOIDCAppRequest.AsObject,
  }
}

export class DataHumanUser extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): DataHumanUser;

  getUser(): zitadel_management_pb.ImportHumanUserRequest | undefined;
  setUser(value?: zitadel_management_pb.ImportHumanUserRequest): DataHumanUser;
  hasUser(): boolean;
  clearUser(): DataHumanUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataHumanUser.AsObject;
  static toObject(includeInstance: boolean, msg: DataHumanUser): DataHumanUser.AsObject;
  static serializeBinaryToWriter(message: DataHumanUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataHumanUser;
  static deserializeBinaryFromReader(message: DataHumanUser, reader: jspb.BinaryReader): DataHumanUser;
}

export namespace DataHumanUser {
  export type AsObject = {
    userId: string,
    user?: zitadel_management_pb.ImportHumanUserRequest.AsObject,
  }
}

export class DataMachineUser extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): DataMachineUser;

  getUser(): zitadel_management_pb.AddMachineUserRequest | undefined;
  setUser(value?: zitadel_management_pb.AddMachineUserRequest): DataMachineUser;
  hasUser(): boolean;
  clearUser(): DataMachineUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataMachineUser.AsObject;
  static toObject(includeInstance: boolean, msg: DataMachineUser): DataMachineUser.AsObject;
  static serializeBinaryToWriter(message: DataMachineUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataMachineUser;
  static deserializeBinaryFromReader(message: DataMachineUser, reader: jspb.BinaryReader): DataMachineUser;
}

export namespace DataMachineUser {
  export type AsObject = {
    userId: string,
    user?: zitadel_management_pb.AddMachineUserRequest.AsObject,
  }
}

export class DataAction extends jspb.Message {
  getActionId(): string;
  setActionId(value: string): DataAction;

  getAction(): zitadel_management_pb.CreateActionRequest | undefined;
  setAction(value?: zitadel_management_pb.CreateActionRequest): DataAction;
  hasAction(): boolean;
  clearAction(): DataAction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataAction.AsObject;
  static toObject(includeInstance: boolean, msg: DataAction): DataAction.AsObject;
  static serializeBinaryToWriter(message: DataAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataAction;
  static deserializeBinaryFromReader(message: DataAction, reader: jspb.BinaryReader): DataAction;
}

export namespace DataAction {
  export type AsObject = {
    actionId: string,
    action?: zitadel_management_pb.CreateActionRequest.AsObject,
  }
}

export class DataProjectGrant extends jspb.Message {
  getGrantId(): string;
  setGrantId(value: string): DataProjectGrant;

  getProjectGrant(): zitadel_management_pb.AddProjectGrantRequest | undefined;
  setProjectGrant(value?: zitadel_management_pb.AddProjectGrantRequest): DataProjectGrant;
  hasProjectGrant(): boolean;
  clearProjectGrant(): DataProjectGrant;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DataProjectGrant.AsObject;
  static toObject(includeInstance: boolean, msg: DataProjectGrant): DataProjectGrant.AsObject;
  static serializeBinaryToWriter(message: DataProjectGrant, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DataProjectGrant;
  static deserializeBinaryFromReader(message: DataProjectGrant, reader: jspb.BinaryReader): DataProjectGrant;
}

export namespace DataProjectGrant {
  export type AsObject = {
    grantId: string,
    projectGrant?: zitadel_management_pb.AddProjectGrantRequest.AsObject,
  }
}

export class SetTriggerActionsRequest extends jspb.Message {
  getFlowType(): FlowType;
  setFlowType(value: FlowType): SetTriggerActionsRequest;

  getTriggerType(): TriggerType;
  setTriggerType(value: TriggerType): SetTriggerActionsRequest;

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
    flowType: FlowType,
    triggerType: TriggerType,
    actionIdsList: Array<string>,
  }
}

export enum FlowType { 
  FLOW_TYPE_UNSPECIFIED = 0,
  FLOW_TYPE_EXTERNAL_AUTHENTICATION = 1,
}
export enum TriggerType { 
  TRIGGER_TYPE_UNSPECIFIED = 0,
  TRIGGER_TYPE_POST_AUTHENTICATION = 1,
  TRIGGER_TYPE_PRE_CREATION = 2,
  TRIGGER_TYPE_POST_CREATION = 3,
}
