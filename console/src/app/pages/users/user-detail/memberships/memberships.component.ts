import { Component, Input, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import {
    ProjectGrantMemberSearchKey,
    ProjectGrantMemberSearchQuery,
    ProjectGrantMemberView,
    ProjectMemberSearchKey,
    ProjectMemberSearchQuery,
    ProjectMemberView,
    User,
} from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

@Component({
    selector: 'app-memberships',
    templateUrl: './memberships.component.html',
    styleUrls: ['./memberships.component.scss'],
})
export class MembershipsComponent implements OnInit {
    usergrants: ProjectGrantMemberView.AsObject[] = [];
    projectmembers: ProjectMemberView.AsObject[] = [];

    @Input() public user: string = '';

    constructor(
        private orgService: OrgService,
        private projectService: ProjectService,
        private mgmtUserService: MgmtUserService,
        private adminService: AdminService,
        private dialog: MatDialog,
        private toast: ToastService,
    ) { }

    ngOnInit(): void {
        // this.loadManager(this.userId);
    }

    public async loadManager(userId: string): Promise<void> {
        console.log('load managers');
        // manager of granted project
        const projectGrantQuery = new ProjectGrantMemberSearchQuery();
        projectGrantQuery.setKey(ProjectGrantMemberSearchKey.PROJECTGRANTMEMBERSEARCHKEY_USER_ID);
        projectGrantQuery.setValue(userId);

        this.usergrants = (await this.mgmtUserService.SearchProjectGrantMembers(100, 0, [projectGrantQuery]))
            .toObject().resultList;
        console.log(this.usergrants);

        // manager of owned project
        const projectMemberQuery = new ProjectMemberSearchQuery();
        projectMemberQuery.setKey(ProjectMemberSearchKey.PROJECTMEMBERSEARCHKEY_USER_ID);
        projectMemberQuery.setValue(userId);

        this.projectmembers = (await this.mgmtUserService.SearchProjectMembers(100, 0, [projectMemberQuery]))
            .toObject().resultList;
        console.log(this.projectmembers);

        // manager of organization
        // const projectMemberQuery = new ProjectMemberSearchQuery();
        // projectMemberQuery.setKey(ProjectMemberSearchKey.PROJECTMEMBERSEARCHKEY_USER_ID);
        // projectMemberQuery.setValue(userId);

        // this.projectmembers = (await this.mgmtUserService.searchor(100, 0, [projectMemberQuery]))
        //     .toObject().resultList;
        // console.log(this.projectmembers);
    }

    public addMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            width: '400px',
            data: {
                user: this.user,
            },
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && resp.creationType !== undefined) {
                console.log(resp);
                switch (resp.creationType) {
                    case CreationType.IAM:
                        this.createIamMember(resp);
                        break;
                    case CreationType.ORG:
                        this.createOrgMember(resp);
                        break;
                    case CreationType.PROJECT_OWNED:
                        this.createOwnedProjectMember(resp);
                        break;
                    case CreationType.PROJECT_GRANTED:
                        this.createGrantedProjectMember(resp);
                        break;
                }
            }
        });
    }

    private createIamMember(response: any): void {
        const users: User.AsObject[] = response.users;
        const roles: string[] = response.roles;

        if (users && users.length && roles && roles.length) {
            Promise.all(users.map(user => {
                return this.adminService.AddIamMember(user.id, roles);
            })).then(() => {
                this.toast.showInfo('IAM.TOAST.MEMBERADDED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    private createOrgMember(response: any): void {
        const users: User.AsObject[] = response.users;
        const roles: string[] = response.roles;

        if (users && users.length && roles && roles.length) {
            Promise.all(users.map(user => {
                return this.orgService.AddMyOrgMember(user.id, roles);
            })).then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    private createGrantedProjectMember(response: any): void {
        const users: User.AsObject[] = response.users;
        const roles: string[] = response.roles;

        if (users && users.length && roles && roles.length) {
            users.forEach(user => {
                return this.projectService.AddProjectGrantMember(
                    response.projectId,
                    response.grantId,
                    user.id,
                    roles,
                ).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            });
        }
    }

    private createOwnedProjectMember(response: any): void {
        const users: User.AsObject[] = response.users;
        const roles: string[] = response.roles;

        if (users && users.length && roles && roles.length) {
            users.forEach(user => {
                return this.projectService.AddProjectMember(response.projectId, user.id, roles)
                    .then(() => {
                        this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
            });
        }
    }
}
