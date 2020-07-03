import { Injectable } from '@angular/core';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { Metadata } from 'grpc-web';

import { ManagementServicePromiseClient } from '../proto/generated/management_grpc_web_pb';
import {
    Application,
    ApplicationID,
    ApplicationSearchQuery,
    ApplicationSearchRequest,
    ApplicationSearchResponse,
    ApplicationUpdate,
    GrantedProjectSearchRequest,
    OIDCApplicationCreate,
    OIDCConfig,
    OIDCConfigUpdate,
    Project,
    ProjectCreateRequest,
    ProjectGrant,
    ProjectGrantCreate,
    ProjectGrantID,
    ProjectGrantMemberAdd,
    ProjectGrantMemberRemove,
    ProjectGrantMemberRoles,
    ProjectGrantMemberSearchQuery,
    ProjectGrantMemberSearchRequest,
    ProjectGrantSearchRequest,
    ProjectGrantSearchResponse,
    ProjectGrantUpdate,
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
    ProjectRoleSearchQuery,
    ProjectRoleSearchRequest,
    ProjectRoleSearchResponse,
    ProjectSearchQuery,
    ProjectSearchRequest,
    ProjectSearchResponse,
    ProjectUpdateRequest,
    ProjectUserGrantSearchRequest,
    ProjectView,
    UserGrant,
    UserGrantCreate,
    UserGrantSearchQuery,
    UserGrantSearchResponse,
} from '../proto/generated/management_pb';
import { GrpcBackendService } from './grpc-backend.service';
import { GrpcService, RequestFactory, ResponseMapper } from './grpc.service';

@Injectable({
    providedIn: 'root',
})
export class ProjectService {
    constructor(private readonly grpcService: GrpcService, private grpcBackendService: GrpcBackendService) { }

    public async request<TReq, TResp, TMappedResp>(
        requestFn: RequestFactory<ManagementServicePromiseClient, TReq, TResp>,
        request: TReq,
        responseMapper: ResponseMapper<TResp, TMappedResp>,
        metadata?: Metadata,
    ): Promise<TMappedResp> {
        const mappedRequestFn = requestFn(this.grpcService.mgmt).bind(this.grpcService.mgmt);
        const response = await this.grpcBackendService.runRequest(
            mappedRequestFn,
            request,
            metadata,
        );
        return responseMapper(response);
    }

    public async SearchProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectSearchResponse> {
        const req = new ProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchProjects,
            req,
            f => f,
        );
    }

    public async SearchGrantedProjects(
        limit: number, offset: number, queryList?: ProjectSearchQuery[]): Promise<ProjectGrantSearchResponse> {
        const req = new GrantedProjectSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchGrantedProjects,
            req,
            f => f,
        );
    }

    public async GetProjectById(projectId: string): Promise<ProjectView> {
        const req = new ProjectID();
        req.setId(projectId);
        return await this.request(
            c => c.projectByID,
            req,
            f => f,
        );
    }

    public async GetGrantedProjectByID(projectId: string, id: string): Promise<ProjectGrantView> {
        const req = new ProjectGrantID();
        req.setId(id);
        req.setProjectId(projectId);
        return await this.request(
            c => c.getGrantedProjectByID,
            req,
            f => f,
        );
    }

    public async CreateProject(project: ProjectCreateRequest.AsObject): Promise<Project> {
        const req = new ProjectCreateRequest();
        req.setName(project.name);
        return await this.request(
            c => c.createProject,
            req,
            f => f,
        );
    }

    public async UpdateProject(id: string, name: string): Promise<Project> {
        const req = new ProjectUpdateRequest();
        req.setName(name);
        req.setId(id);
        return await this.request(
            c => c.updateProject,
            req,
            f => f,
        );
    }

    public async UpdateProjectGrant(id: string, projectId: string, rolesList: string[]): Promise<ProjectGrant> {
        const req = new ProjectGrantUpdate();
        req.setRoleKeysList(rolesList);
        req.setId(id);
        req.setProjectId(projectId);
        return await this.request(
            c => c.updateProjectGrant,
            req,
            f => f,
        );
    }

    public async DeactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return await this.request(
            c => c.deactivateProject,
            req,
            f => f,
        );
    }

    public async ReactivateProject(projectId: string): Promise<Project> {
        const req = new ProjectID();
        req.setId(projectId);
        return await this.request(
            c => c.reactivateProject,
            req,
            f => f,
        );
    }

    public async SearchProjectGrants(projectId: string, limit: number, offset: number): Promise<ProjectGrantSearchResponse> {
        const req = new ProjectGrantSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        return await this.request(
            c => c.searchProjectGrants,
            req,
            f => f,
        );
    }

    public async GetProjectGrantMemberRoles(): Promise<ProjectGrantMemberRoles> {
        const req = new Empty();
        return await this.request(
            c => c.getProjectGrantMemberRoles,
            req,
            f => f,
        );
    }

    public async AddProjectMember(id: string, userId: string, rolesList: string[]): Promise<Empty> {
        const req = new ProjectMemberAdd();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return await this.request(
            c => c.addProjectMember,
            req,
            f => f,
        );
    }

    public async ChangeProjectMember(id: string, userId: string, rolesList: string[]): Promise<ProjectMember> {
        const req = new ProjectMemberChange();
        req.setId(id);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return await this.request(
            c => c.changeProjectMember,
            req,
            f => f,
        );
    }

    public async AddProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
        rolesList: string[],
    ): Promise<Empty> {
        const req = new ProjectGrantMemberAdd();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        req.setRolesList(rolesList);
        return await this.request(
            c => c.addProjectGrantMember,
            req,
            f => f,
        );
    }

    public async SearchProjectGrantMembers(
        projectId: string,
        grantId: string,
        limit: number,
        offset: number,
        queryList?: ProjectGrantMemberSearchQuery[],
    ): Promise<ProjectMemberSearchResponse> {
        const req = new ProjectGrantMemberSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        return await this.request(
            c => c.searchProjectGrantMembers,
            req,
            f => f,
        );
    }

    public async RemoveProjectGrantMember(
        projectId: string,
        grantId: string,
        userId: string,
    ): Promise<Empty> {
        const req = new ProjectGrantMemberRemove();
        req.setProjectId(projectId);
        req.setGrantId(grantId);
        req.setUserId(userId);
        return await this.request(
            c => c.removeProjectGrantMember,
            req,
            f => f,
        );
    }


    public async CreateProjectGrant(orgId: string, projectId: string, roleKeys: string[]): Promise<ProjectGrant> {
        const req = new ProjectGrantCreate();
        req.setGrantedOrgId(orgId);
        req.setProjectId(projectId);
        req.setRoleKeysList(roleKeys);
        return await this.request(
            c => c.createProjectGrant,
            req,
            f => f,
        );
    }

    public async ReactivateApplication(appId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(appId);
        return await this.request(
            c => c.reactivateApplication,
            req,
            f => f,
        );
    }

    public async DectivateApplication(projectId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setId(projectId);
        return await this.request(
            c => c.deactivateApplication,
            req,
            f => f,
        );
    }

    public async RegenerateOIDCClientSecret(id: string, projectId: string): Promise<any> {
        const req = new ApplicationID();
        req.setId(id);
        req.setProjectId(projectId);
        return await this.request(
            c => c.regenerateOIDCClientSecret,
            req,
            f => f,
        );
    }

    public async SearchProjectRoles(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ProjectRoleSearchQuery[],
    ): Promise<ProjectRoleSearchResponse> {
        const req = new ProjectRoleSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchProjectRoles,
            req,
            f => f,
        );
    }

    public async AddProjectRole(role: ProjectRoleAdd.AsObject): Promise<Empty> {
        const req = new ProjectRoleAdd();
        req.setId(role.id);
        if (role.displayName) {
            req.setDisplayName(role.displayName);
        }
        req.setKey(role.key);
        req.setGroup(role.group);
        return await this.request(
            c => c.addProjectRole,
            req,
            f => f,
        );
    }

    public async BulkAddProjectRole(
        id: string,
        rolesList: ProjectRoleAdd[],
    ): Promise<Empty> {
        const req = new ProjectRoleAddBulk();
        req.setId(id);
        req.setProjectRolesList(rolesList);
        return await this.request(
            c => c.bulkAddProjectRole,
            req,
            f => f,
        );
    }

    public async RemoveProjectRole(projectId: string, key: string): Promise<Empty> {
        const req = new ProjectRoleRemove();
        req.setId(projectId);
        req.setKey(key);
        return await this.request(
            c => c.removeProjectRole,
            req,
            f => f,
        );
    }


    public async ChangeProjectRole(projectId: string, key: string, displayName: string, group: string):
        Promise<ProjectRole> {
        const req = new ProjectRoleChange();
        req.setId(projectId);
        req.setKey(key);
        req.setGroup(group);
        req.setDisplayName(displayName);
        return await this.request(
            c => c.changeProjectRole,
            req,
            f => f,
        );
    }


    public async RemoveProjectMember(id: string, userId: string): Promise<Empty> {
        const req = new ProjectMemberRemove();
        req.setId(id);
        req.setUserId(userId);
        return await this.request(
            c => c.removeProjectMember,
            req,
            f => f,
        );
    }

    public async SearchProjectMembers(projectId: string,
        limit: number, offset: number): Promise<ProjectMemberSearchResponse> {
        const req = new ProjectMemberSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        return await this.request(
            c => c.searchProjectMembers,
            req,
            f => f,
        );
    }

    public async SearchApplications(
        projectId: string,
        limit: number,
        offset: number,
        queryList?: ApplicationSearchQuery[]): Promise<ApplicationSearchResponse> {
        const req = new ApplicationSearchRequest();
        req.setProjectId(projectId);
        req.setLimit(limit);
        req.setOffset(offset);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchApplications,
            req,
            f => f,
        );
    }

    public async GetApplicationById(projectId: string, applicationId: string): Promise<Application> {
        const req = new ApplicationID();
        req.setProjectId(projectId);
        req.setId(applicationId);
        return await this.request(
            c => c.applicationByID,
            req,
            f => f,
        );
    }

    public async GetProjectMemberRoles(): Promise<ProjectMemberRoles> {
        const req = new Empty();
        return await this.request(
            c => c.getProjectMemberRoles,
            req,
            f => f,
        );
    }

    public async ProjectGrantByID(id: string): Promise<ProjectGrant> {
        const req = new ProjectGrantID();
        return await this.request(
            c => c.projectGrantByID,
            req,
            f => f,
        );
    }

    // ********* */

    public async SearchProjectUserGrants(
        projectId: string,
        offset: number,
        limit: number,
        queryList?: UserGrantSearchQuery[],
    ): Promise<UserGrantSearchResponse> {
        const req = new ProjectUserGrantSearchRequest();
        req.setLimit(limit);
        req.setOffset(offset);
        req.setProjectId(projectId);
        if (queryList) {
            req.setQueriesList(queryList);
        }
        return await this.request(
            c => c.searchProjectUserGrants,
            req,
            f => f,
        );
    }

    public async CreateProjectUserGrant(
        projectId: string,
        userId: string,
        roleKeysList: string[],
    ): Promise<UserGrant> {
        const req = new UserGrantCreate();
        req.setProjectId(projectId);
        req.setRoleKeysList(roleKeysList);
        req.setUserId(userId);

        return await this.request(
            c => c.createProjectUserGrant,
            req,
            f => f,
        );
    }

    // ********* */

    public async CreateOIDCApp(app: OIDCApplicationCreate.AsObject): Promise<Application> {
        const req = new OIDCApplicationCreate();
        req.setProjectId(app.projectId);
        req.setName(app.name);
        req.setRedirectUrisList(app.redirectUrisList);
        req.setResponseTypesList(app.responseTypesList);
        req.setGrantTypesList(app.grantTypesList);
        req.setApplicationType(app.applicationType);
        req.setAuthMethodType(app.authMethodType);
        req.setPostLogoutRedirectUrisList(app.postLogoutRedirectUrisList);

        return await this.request(
            c => c.createOIDCApplication,
            req,
            f => f,
        );
    }

    public async UpdateApplication(projectId: string, appId: string, name: string): Promise<Application> {
        const req = new ApplicationUpdate();
        req.setId(appId);
        req.setName(name);
        req.setProjectId(projectId);
        return await this.request(
            c => c.updateApplication,
            req,
            f => f,
        );
    }

    public async UpdateOIDCAppConfig(projectId: string,
        appId: string, oidcConfig: OIDCConfig.AsObject): Promise<OIDCConfig> {
        const req = new OIDCConfigUpdate();
        req.setProjectId(projectId);
        req.setApplicationId(appId);
        req.setRedirectUrisList(oidcConfig.redirectUrisList);
        req.setResponseTypesList(oidcConfig.responseTypesList);
        req.setAuthMethodType(oidcConfig.authMethodType);
        req.setPostLogoutRedirectUrisList(oidcConfig.postLogoutRedirectUrisList);
        req.setGrantTypesList(oidcConfig.grantTypesList);
        req.setApplicationType(oidcConfig.applicationType);
        return await this.request(
            c => c.updateApplicationOIDCConfig,
            req,
            f => f,
        );
    }
}
