import * as grpcWeb from 'grpc-web';

import * as google_api_annotations_pb from './google/api/annotations_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';
import * as protoc$gen$swagger_options_annotations_pb from './protoc-gen-swagger/options/annotations_pb';
import * as validate_validate_pb from './validate/validate_pb';
import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as authoption_options_pb from './authoption/options_pb';

import {
  AddOrgDomainRequest,
  AddOrgMemberRequest,
  Application,
  ApplicationID,
  ApplicationSearchRequest,
  ApplicationSearchResponse,
  ApplicationUpdate,
  ApplicationView,
  AuthGrantSearchRequest,
  AuthGrantSearchResponse,
  ChangeOrgMemberRequest,
  ChangeRequest,
  Changes,
  ClientSecret,
  CreateUserRequest,
  GrantedProjectSearchRequest,
  Iam,
  MultiFactors,
  OIDCApplicationCreate,
  OIDCConfig,
  OIDCConfigUpdate,
  Org,
  OrgDomain,
  OrgDomainSearchRequest,
  OrgDomainSearchResponse,
  OrgID,
  OrgIamPolicy,
  OrgMember,
  OrgMemberRoles,
  OrgMemberSearchRequest,
  OrgMemberSearchResponse,
  OrgView,
  PasswordAgePolicy,
  PasswordAgePolicyCreate,
  PasswordAgePolicyID,
  PasswordAgePolicyUpdate,
  PasswordComplexityPolicy,
  PasswordComplexityPolicyCreate,
  PasswordComplexityPolicyID,
  PasswordComplexityPolicyUpdate,
  PasswordLockoutPolicy,
  PasswordLockoutPolicyCreate,
  PasswordLockoutPolicyID,
  PasswordLockoutPolicyUpdate,
  PasswordRequest,
  Project,
  ProjectCreateRequest,
  ProjectGrant,
  ProjectGrantCreate,
  ProjectGrantID,
  ProjectGrantMember,
  ProjectGrantMemberAdd,
  ProjectGrantMemberChange,
  ProjectGrantMemberRemove,
  ProjectGrantMemberRoles,
  ProjectGrantMemberSearchRequest,
  ProjectGrantMemberSearchResponse,
  ProjectGrantSearchRequest,
  ProjectGrantSearchResponse,
  ProjectGrantUpdate,
  ProjectGrantUserGrantCreate,
  ProjectGrantUserGrantID,
  ProjectGrantUserGrantSearchRequest,
  ProjectGrantUserGrantUpdate,
  ProjectGrantView,
  ProjectID,
  ProjectMember,
  ProjectMemberAdd,
  ProjectMemberChange,
  ProjectMemberRemove,
  ProjectMemberRoles,
  ProjectMemberSearchRequest,
  ProjectMemberSearchResponse,
  ProjectRole,
  ProjectRoleAdd,
  ProjectRoleAddBulk,
  ProjectRoleChange,
  ProjectRoleRemove,
  ProjectRoleSearchRequest,
  ProjectRoleSearchResponse,
  ProjectSearchRequest,
  ProjectSearchResponse,
  ProjectUpdateRequest,
  ProjectUserGrantID,
  ProjectUserGrantSearchRequest,
  ProjectUserGrantUpdate,
  ProjectView,
  RemoveOrgDomainRequest,
  RemoveOrgMemberRequest,
  SetPasswordNotificationRequest,
  UniqueUserRequest,
  UniqueUserResponse,
  UpdateUserAddressRequest,
  UpdateUserEmailRequest,
  UpdateUserPhoneRequest,
  UpdateUserProfileRequest,
  User,
  UserAddress,
  UserAddressView,
  UserEmail,
  UserEmailID,
  UserEmailView,
  UserGrant,
  UserGrantCreate,
  UserGrantCreateBulk,
  UserGrantID,
  UserGrantRemoveBulk,
  UserGrantSearchRequest,
  UserGrantSearchResponse,
  UserGrantUpdate,
  UserGrantUpdateBulk,
  UserGrantView,
  UserID,
  UserPhone,
  UserPhoneView,
  UserProfile,
  UserProfileView,
  UserSearchRequest,
  UserSearchResponse,
  UserView} from './management_pb';

export class ManagementServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_struct_pb.Struct) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_struct_pb.Struct>;

  getIam(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Iam) => void
  ): grpcWeb.ClientReadableStream<Iam>;

  getUserByID(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserView) => void
  ): grpcWeb.ClientReadableStream<UserView>;

  getUserByEmailGlobal(
    request: UserEmailID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserView) => void
  ): grpcWeb.ClientReadableStream<UserView>;

  searchUsers(
    request: UserSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserSearchResponse) => void
  ): grpcWeb.ClientReadableStream<UserSearchResponse>;

  isUserUnique(
    request: UniqueUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UniqueUserResponse) => void
  ): grpcWeb.ClientReadableStream<UniqueUserResponse>;

  createUser(
    request: CreateUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  deactivateUser(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  reactivateUser(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  lockUser(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  unlockUser(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: User) => void
  ): grpcWeb.ClientReadableStream<User>;

  deleteUser(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  userChanges(
    request: ChangeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Changes) => void
  ): grpcWeb.ClientReadableStream<Changes>;

  applicationChanges(
    request: ChangeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Changes) => void
  ): grpcWeb.ClientReadableStream<Changes>;

  orgChanges(
    request: ChangeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Changes) => void
  ): grpcWeb.ClientReadableStream<Changes>;

  projectChanges(
    request: ChangeRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Changes) => void
  ): grpcWeb.ClientReadableStream<Changes>;

  getUserProfile(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserProfileView) => void
  ): grpcWeb.ClientReadableStream<UserProfileView>;

  updateUserProfile(
    request: UpdateUserProfileRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserProfile) => void
  ): grpcWeb.ClientReadableStream<UserProfile>;

  getUserEmail(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserEmailView) => void
  ): grpcWeb.ClientReadableStream<UserEmailView>;

  changeUserEmail(
    request: UpdateUserEmailRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserEmail) => void
  ): grpcWeb.ClientReadableStream<UserEmail>;

  resendEmailVerificationMail(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getUserPhone(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserPhoneView) => void
  ): grpcWeb.ClientReadableStream<UserPhoneView>;

  changeUserPhone(
    request: UpdateUserPhoneRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserPhone) => void
  ): grpcWeb.ClientReadableStream<UserPhone>;

  resendPhoneVerificationCode(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getUserAddress(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAddressView) => void
  ): grpcWeb.ClientReadableStream<UserAddressView>;

  updateUserAddress(
    request: UpdateUserAddressRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserAddress) => void
  ): grpcWeb.ClientReadableStream<UserAddress>;

  getUserMfas(
    request: UserID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: MultiFactors) => void
  ): grpcWeb.ClientReadableStream<MultiFactors>;

  sendSetPasswordNotification(
    request: SetPasswordNotificationRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  setInitialPassword(
    request: PasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getPasswordComplexityPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordComplexityPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordComplexityPolicy>;

  createPasswordComplexityPolicy(
    request: PasswordComplexityPolicyCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordComplexityPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordComplexityPolicy>;

  updatePasswordComplexityPolicy(
    request: PasswordComplexityPolicyUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordComplexityPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordComplexityPolicy>;

  deletePasswordComplexityPolicy(
    request: PasswordComplexityPolicyID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getPasswordAgePolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordAgePolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordAgePolicy>;

  createPasswordAgePolicy(
    request: PasswordAgePolicyCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordAgePolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordAgePolicy>;

  updatePasswordAgePolicy(
    request: PasswordAgePolicyUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordAgePolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordAgePolicy>;

  deletePasswordAgePolicy(
    request: PasswordAgePolicyID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getPasswordLockoutPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordLockoutPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordLockoutPolicy>;

  createPasswordLockoutPolicy(
    request: PasswordLockoutPolicyCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordLockoutPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordLockoutPolicy>;

  updatePasswordLockoutPolicy(
    request: PasswordLockoutPolicyUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: PasswordLockoutPolicy) => void
  ): grpcWeb.ClientReadableStream<PasswordLockoutPolicy>;

  deletePasswordLockoutPolicy(
    request: PasswordLockoutPolicyID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getOrgByID(
    request: OrgID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgView) => void
  ): grpcWeb.ClientReadableStream<OrgView>;

  getOrgByDomainGlobal(
    request: OrgDomain,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgView) => void
  ): grpcWeb.ClientReadableStream<OrgView>;

  deactivateOrg(
    request: OrgID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Org) => void
  ): grpcWeb.ClientReadableStream<Org>;

  reactivateOrg(
    request: OrgID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Org) => void
  ): grpcWeb.ClientReadableStream<Org>;

  searchMyOrgDomains(
    request: OrgDomainSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgDomainSearchResponse) => void
  ): grpcWeb.ClientReadableStream<OrgDomainSearchResponse>;

  addMyOrgDomain(
    request: AddOrgDomainRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgDomain) => void
  ): grpcWeb.ClientReadableStream<OrgDomain>;

  removeMyOrgDomain(
    request: RemoveOrgDomainRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getMyOrgIamPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgIamPolicy) => void
  ): grpcWeb.ClientReadableStream<OrgIamPolicy>;

  getOrgMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgMemberRoles) => void
  ): grpcWeb.ClientReadableStream<OrgMemberRoles>;

  addMyOrgMember(
    request: AddOrgMemberRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgMember) => void
  ): grpcWeb.ClientReadableStream<OrgMember>;

  changeMyOrgMember(
    request: ChangeOrgMemberRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgMember) => void
  ): grpcWeb.ClientReadableStream<OrgMember>;

  removeMyOrgMember(
    request: RemoveOrgMemberRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchMyOrgMembers(
    request: OrgMemberSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OrgMemberSearchResponse) => void
  ): grpcWeb.ClientReadableStream<OrgMemberSearchResponse>;

  searchProjects(
    request: ProjectSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectSearchResponse>;

  projectByID(
    request: ProjectID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectView) => void
  ): grpcWeb.ClientReadableStream<ProjectView>;

  createProject(
    request: ProjectCreateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Project) => void
  ): grpcWeb.ClientReadableStream<Project>;

  updateProject(
    request: ProjectUpdateRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Project) => void
  ): grpcWeb.ClientReadableStream<Project>;

  deactivateProject(
    request: ProjectID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Project) => void
  ): grpcWeb.ClientReadableStream<Project>;

  reactivateProject(
    request: ProjectID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Project) => void
  ): grpcWeb.ClientReadableStream<Project>;

  searchGrantedProjects(
    request: GrantedProjectSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantSearchResponse>;

  getGrantedProjectByID(
    request: ProjectGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantView) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantView>;

  getProjectMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectMemberRoles) => void
  ): grpcWeb.ClientReadableStream<ProjectMemberRoles>;

  searchProjectMembers(
    request: ProjectMemberSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectMemberSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectMemberSearchResponse>;

  addProjectMember(
    request: ProjectMemberAdd,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectMember) => void
  ): grpcWeb.ClientReadableStream<ProjectMember>;

  changeProjectMember(
    request: ProjectMemberChange,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectMember) => void
  ): grpcWeb.ClientReadableStream<ProjectMember>;

  removeProjectMember(
    request: ProjectMemberRemove,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchProjectRoles(
    request: ProjectRoleSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectRoleSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectRoleSearchResponse>;

  addProjectRole(
    request: ProjectRoleAdd,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectRole) => void
  ): grpcWeb.ClientReadableStream<ProjectRole>;

  bulkAddProjectRole(
    request: ProjectRoleAddBulk,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  changeProjectRole(
    request: ProjectRoleChange,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectRole) => void
  ): grpcWeb.ClientReadableStream<ProjectRole>;

  removeProjectRole(
    request: ProjectRoleRemove,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchApplications(
    request: ApplicationSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ApplicationSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ApplicationSearchResponse>;

  applicationByID(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ApplicationView) => void
  ): grpcWeb.ClientReadableStream<ApplicationView>;

  createOIDCApplication(
    request: OIDCApplicationCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  updateApplication(
    request: ApplicationUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  deactivateApplication(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  reactivateApplication(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: Application) => void
  ): grpcWeb.ClientReadableStream<Application>;

  removeApplication(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  updateApplicationOIDCConfig(
    request: OIDCConfigUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: OIDCConfig) => void
  ): grpcWeb.ClientReadableStream<OIDCConfig>;

  regenerateOIDCClientSecret(
    request: ApplicationID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ClientSecret) => void
  ): grpcWeb.ClientReadableStream<ClientSecret>;

  searchProjectGrants(
    request: ProjectGrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantSearchResponse>;

  projectGrantByID(
    request: ProjectGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantView) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantView>;

  createProjectGrant(
    request: ProjectGrantCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrant) => void
  ): grpcWeb.ClientReadableStream<ProjectGrant>;

  updateProjectGrant(
    request: ProjectGrantUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrant) => void
  ): grpcWeb.ClientReadableStream<ProjectGrant>;

  deactivateProjectGrant(
    request: ProjectGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrant) => void
  ): grpcWeb.ClientReadableStream<ProjectGrant>;

  reactivateProjectGrant(
    request: ProjectGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrant) => void
  ): grpcWeb.ClientReadableStream<ProjectGrant>;

  removeProjectGrant(
    request: ProjectGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  getProjectGrantMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantMemberRoles) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantMemberRoles>;

  searchProjectGrantMembers(
    request: ProjectGrantMemberSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantMemberSearchResponse) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantMemberSearchResponse>;

  addProjectGrantMember(
    request: ProjectGrantMemberAdd,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantMember) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantMember>;

  changeProjectGrantMember(
    request: ProjectGrantMemberChange,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: ProjectGrantMember) => void
  ): grpcWeb.ClientReadableStream<ProjectGrantMember>;

  removeProjectGrantMember(
    request: ProjectGrantMemberRemove,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchUserGrants(
    request: UserGrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<UserGrantSearchResponse>;

  userGrantByID(
    request: UserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantView) => void
  ): grpcWeb.ClientReadableStream<UserGrantView>;

  createUserGrant(
    request: UserGrantCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  updateUserGrant(
    request: UserGrantUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  deactivateUserGrant(
    request: UserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  reactivateUserGrant(
    request: UserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  removeUserGrant(
    request: UserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  bulkCreateUserGrant(
    request: UserGrantCreateBulk,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  bulkUpdateUserGrant(
    request: UserGrantUpdateBulk,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  bulkRemoveUserGrant(
    request: UserGrantRemoveBulk,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  searchProjectUserGrants(
    request: ProjectUserGrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<UserGrantSearchResponse>;

  projectUserGrantByID(
    request: ProjectUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantView) => void
  ): grpcWeb.ClientReadableStream<UserGrantView>;

  createProjectUserGrant(
    request: UserGrantCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  updateProjectUserGrant(
    request: ProjectUserGrantUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  deactivateProjectUserGrant(
    request: ProjectUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  reactivateProjectUserGrant(
    request: ProjectUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  searchProjectGrantUserGrants(
    request: ProjectGrantUserGrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<UserGrantSearchResponse>;

  projectGrantUserGrantByID(
    request: ProjectGrantUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrantView) => void
  ): grpcWeb.ClientReadableStream<UserGrantView>;

  createProjectGrantUserGrant(
    request: ProjectGrantUserGrantCreate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  updateProjectGrantUserGrant(
    request: ProjectGrantUserGrantUpdate,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  deactivateProjectGrantUserGrant(
    request: ProjectGrantUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  reactivateProjectGrantUserGrant(
    request: ProjectGrantUserGrantID,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: UserGrant) => void
  ): grpcWeb.ClientReadableStream<UserGrant>;

  searchAuthGrant(
    request: AuthGrantSearchRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: AuthGrantSearchResponse) => void
  ): grpcWeb.ClientReadableStream<AuthGrantSearchResponse>;

}

export class ManagementServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; });

  healthz(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  ready(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  validate(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_struct_pb.Struct>;

  getIam(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<Iam>;

  getUserByID(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserView>;

  getUserByEmailGlobal(
    request: UserEmailID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserView>;

  searchUsers(
    request: UserSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserSearchResponse>;

  isUserUnique(
    request: UniqueUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UniqueUserResponse>;

  createUser(
    request: CreateUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  deactivateUser(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  reactivateUser(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  lockUser(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  unlockUser(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<User>;

  deleteUser(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  userChanges(
    request: ChangeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Changes>;

  applicationChanges(
    request: ChangeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Changes>;

  orgChanges(
    request: ChangeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Changes>;

  projectChanges(
    request: ChangeRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Changes>;

  getUserProfile(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserProfileView>;

  updateUserProfile(
    request: UpdateUserProfileRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserProfile>;

  getUserEmail(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserEmailView>;

  changeUserEmail(
    request: UpdateUserEmailRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserEmail>;

  resendEmailVerificationMail(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getUserPhone(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserPhoneView>;

  changeUserPhone(
    request: UpdateUserPhoneRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserPhone>;

  resendPhoneVerificationCode(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getUserAddress(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAddressView>;

  updateUserAddress(
    request: UpdateUserAddressRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserAddress>;

  getUserMfas(
    request: UserID,
    metadata?: grpcWeb.Metadata
  ): Promise<MultiFactors>;

  sendSetPasswordNotification(
    request: SetPasswordNotificationRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  setInitialPassword(
    request: PasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getPasswordComplexityPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordComplexityPolicy>;

  createPasswordComplexityPolicy(
    request: PasswordComplexityPolicyCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordComplexityPolicy>;

  updatePasswordComplexityPolicy(
    request: PasswordComplexityPolicyUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordComplexityPolicy>;

  deletePasswordComplexityPolicy(
    request: PasswordComplexityPolicyID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getPasswordAgePolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordAgePolicy>;

  createPasswordAgePolicy(
    request: PasswordAgePolicyCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordAgePolicy>;

  updatePasswordAgePolicy(
    request: PasswordAgePolicyUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordAgePolicy>;

  deletePasswordAgePolicy(
    request: PasswordAgePolicyID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getPasswordLockoutPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordLockoutPolicy>;

  createPasswordLockoutPolicy(
    request: PasswordLockoutPolicyCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordLockoutPolicy>;

  updatePasswordLockoutPolicy(
    request: PasswordLockoutPolicyUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<PasswordLockoutPolicy>;

  deletePasswordLockoutPolicy(
    request: PasswordLockoutPolicyID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getOrgByID(
    request: OrgID,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgView>;

  getOrgByDomainGlobal(
    request: OrgDomain,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgView>;

  deactivateOrg(
    request: OrgID,
    metadata?: grpcWeb.Metadata
  ): Promise<Org>;

  reactivateOrg(
    request: OrgID,
    metadata?: grpcWeb.Metadata
  ): Promise<Org>;

  searchMyOrgDomains(
    request: OrgDomainSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgDomainSearchResponse>;

  addMyOrgDomain(
    request: AddOrgDomainRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgDomain>;

  removeMyOrgDomain(
    request: RemoveOrgDomainRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getMyOrgIamPolicy(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgIamPolicy>;

  getOrgMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgMemberRoles>;

  addMyOrgMember(
    request: AddOrgMemberRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgMember>;

  changeMyOrgMember(
    request: ChangeOrgMemberRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgMember>;

  removeMyOrgMember(
    request: RemoveOrgMemberRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchMyOrgMembers(
    request: OrgMemberSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<OrgMemberSearchResponse>;

  searchProjects(
    request: ProjectSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectSearchResponse>;

  projectByID(
    request: ProjectID,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectView>;

  createProject(
    request: ProjectCreateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Project>;

  updateProject(
    request: ProjectUpdateRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<Project>;

  deactivateProject(
    request: ProjectID,
    metadata?: grpcWeb.Metadata
  ): Promise<Project>;

  reactivateProject(
    request: ProjectID,
    metadata?: grpcWeb.Metadata
  ): Promise<Project>;

  searchGrantedProjects(
    request: GrantedProjectSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantSearchResponse>;

  getGrantedProjectByID(
    request: ProjectGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantView>;

  getProjectMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectMemberRoles>;

  searchProjectMembers(
    request: ProjectMemberSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectMemberSearchResponse>;

  addProjectMember(
    request: ProjectMemberAdd,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectMember>;

  changeProjectMember(
    request: ProjectMemberChange,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectMember>;

  removeProjectMember(
    request: ProjectMemberRemove,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchProjectRoles(
    request: ProjectRoleSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectRoleSearchResponse>;

  addProjectRole(
    request: ProjectRoleAdd,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectRole>;

  bulkAddProjectRole(
    request: ProjectRoleAddBulk,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  changeProjectRole(
    request: ProjectRoleChange,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectRole>;

  removeProjectRole(
    request: ProjectRoleRemove,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchApplications(
    request: ApplicationSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ApplicationSearchResponse>;

  applicationByID(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<ApplicationView>;

  createOIDCApplication(
    request: OIDCApplicationCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  updateApplication(
    request: ApplicationUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  deactivateApplication(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  reactivateApplication(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<Application>;

  removeApplication(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  updateApplicationOIDCConfig(
    request: OIDCConfigUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<OIDCConfig>;

  regenerateOIDCClientSecret(
    request: ApplicationID,
    metadata?: grpcWeb.Metadata
  ): Promise<ClientSecret>;

  searchProjectGrants(
    request: ProjectGrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantSearchResponse>;

  projectGrantByID(
    request: ProjectGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantView>;

  createProjectGrant(
    request: ProjectGrantCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrant>;

  updateProjectGrant(
    request: ProjectGrantUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrant>;

  deactivateProjectGrant(
    request: ProjectGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrant>;

  reactivateProjectGrant(
    request: ProjectGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrant>;

  removeProjectGrant(
    request: ProjectGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  getProjectGrantMemberRoles(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantMemberRoles>;

  searchProjectGrantMembers(
    request: ProjectGrantMemberSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantMemberSearchResponse>;

  addProjectGrantMember(
    request: ProjectGrantMemberAdd,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantMember>;

  changeProjectGrantMember(
    request: ProjectGrantMemberChange,
    metadata?: grpcWeb.Metadata
  ): Promise<ProjectGrantMember>;

  removeProjectGrantMember(
    request: ProjectGrantMemberRemove,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchUserGrants(
    request: UserGrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantSearchResponse>;

  userGrantByID(
    request: UserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantView>;

  createUserGrant(
    request: UserGrantCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  updateUserGrant(
    request: UserGrantUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  deactivateUserGrant(
    request: UserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  reactivateUserGrant(
    request: UserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  removeUserGrant(
    request: UserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  bulkCreateUserGrant(
    request: UserGrantCreateBulk,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  bulkUpdateUserGrant(
    request: UserGrantUpdateBulk,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  bulkRemoveUserGrant(
    request: UserGrantRemoveBulk,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  searchProjectUserGrants(
    request: ProjectUserGrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantSearchResponse>;

  projectUserGrantByID(
    request: ProjectUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantView>;

  createProjectUserGrant(
    request: UserGrantCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  updateProjectUserGrant(
    request: ProjectUserGrantUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  deactivateProjectUserGrant(
    request: ProjectUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  reactivateProjectUserGrant(
    request: ProjectUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  searchProjectGrantUserGrants(
    request: ProjectGrantUserGrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantSearchResponse>;

  projectGrantUserGrantByID(
    request: ProjectGrantUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrantView>;

  createProjectGrantUserGrant(
    request: ProjectGrantUserGrantCreate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  updateProjectGrantUserGrant(
    request: ProjectGrantUserGrantUpdate,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  deactivateProjectGrantUserGrant(
    request: ProjectGrantUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  reactivateProjectGrantUserGrant(
    request: ProjectGrantUserGrantID,
    metadata?: grpcWeb.Metadata
  ): Promise<UserGrant>;

  searchAuthGrant(
    request: AuthGrantSearchRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<AuthGrantSearchResponse>;

}

