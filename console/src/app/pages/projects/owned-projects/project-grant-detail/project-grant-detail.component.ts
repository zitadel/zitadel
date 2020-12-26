import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { ActivatedRoute } from '@angular/router';
import {
    ProjectGrant,
    ProjectGrantMember,
    ProjectGrantMemberView,
    ProjectGrantState,
    ProjectGrantView,
    ProjectRoleView,
    ProjectType,
} from 'src/app/proto/generated/zitadel/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    ProjectGrantMembersCreateDialogComponent,
    ProjectGrantMembersCreateDialogExportType,
} from './project-grant-members-create-dialog/project-grant-members-create-dialog.component';
import { ProjectGrantMembersDataSource } from './project-grant-members-datasource';

@Component({
    selector: 'app-project-grant-detail',
    templateUrl: './project-grant-detail.component.html',
    styleUrls: ['./project-grant-detail.component.scss'],
})
export class ProjectGrantDetailComponent {
    public INITIALPAGESIZE: number = 25;

    public grant!: ProjectGrantView.AsObject;
    public projectid: string = '';
    public grantid: string = '';

    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;

    public isZitadel: boolean = false;
    ProjectGrantState: any = ProjectGrantState;

    public projectRoleOptions: ProjectRoleView.AsObject[] = [];
    public memberRoleOptions: Array<string> = [];

    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<ProjectGrantMemberView.AsObject> = [];
    public dataSource!: ProjectGrantMembersDataSource;
    constructor(
        private mgmtService: ManagementService,
        private route: ActivatedRoute,
        private toast: ToastService,
        private dialog: MatDialog,
    ) {
        this.route.params.subscribe(params => {
            this.projectid = params.projectid;
            this.grantid = params.grantid;

            this.dataSource = new ProjectGrantMembersDataSource(this.mgmtService);
            this.dataSource.loadGrantMembers(params.projectid, params.grantid, 0, this.INITIALPAGESIZE);

            this.getRoleOptions(params.projectid);
            this.getMemberRoleOptions();

            this.changePageFactory = (event?: PageEvent) => {
                return this.dataSource.loadGrantMembers(
                    params.projectid,
                    params.grantid,
                    event?.pageIndex ?? 0,
                    event?.pageSize ?? this.INITIALPAGESIZE,
                );
            };

            this.mgmtService.ProjectGrantByID(this.grantid, this.projectid).then((grant) => {
                this.grant = grant.toObject();
            });
        });
    }

    public changeState(newState: ProjectGrantState): void {
        if (newState === ProjectGrantState.PROJECTGRANTSTATE_ACTIVE) {
            this.mgmtService.ReactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (newState === ProjectGrantState.PROJECTGRANTSTATE_INACTIVE) {
            this.mgmtService.DeactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public getRoleOptions(projectId: string): void {
        this.mgmtService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.projectRoleOptions = resp.toObject().resultList;
        });
    }

    public getMemberRoleOptions(): void {
        this.mgmtService.GetProjectGrantMemberRoles().then(resp => {
            this.memberRoleOptions = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(selectionChange: MatSelectChange): void {
        this.mgmtService.UpdateProjectGrant(this.grant.id, this.grant.projectId, selectionChange.value)
            .then((newgrant: ProjectGrant) => {
                this.toast.showInfo('PROJECT.TOAST.GRANTUPDATED');
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            return this.mgmtService.RemoveProjectGrantMember(this.grant.projectId, this.grant.id, member.userId).then(() => {
                this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERREMOVED', true);
                setTimeout(() => {
                    this.changePage.emit();
                }, 1000);
            }).catch(error => {
                this.toast.showError(error);
            });
        }));
    }

    public async openAddMember(): Promise<any> {
        const keysList = (await this.mgmtService.GetProjectGrantMemberRoles()).toObject();

        const dialogRef = this.dialog.open(ProjectGrantMembersCreateDialogComponent, {
            data: {
                roleKeysList: keysList.rolesList,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe((dataToAdd: ProjectGrantMembersCreateDialogExportType) => {
            if (dataToAdd) {
                Promise.all(dataToAdd.userIds.map((userid: string) => {
                    return this.mgmtService.AddProjectGrantMember(
                        this.grant.projectId,
                        this.grant.id,
                        userid,
                        dataToAdd.rolesKeyList,
                    );
                })).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERADDED', true);
                    setTimeout(() => {
                        this.changePage.emit();
                    }, 3000);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    updateMemberRoles(member: ProjectGrantMember.AsObject, selectionChange: MatSelectChange): void {
        this.mgmtService.ChangeProjectGrantMember(this.grant.projectId, this.grant.id, member.userId, selectionChange.value)
            .then((_: ProjectGrantMember) => {
                this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }
}
