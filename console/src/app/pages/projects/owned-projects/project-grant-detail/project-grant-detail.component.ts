import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { ActivatedRoute } from '@angular/router';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { GrantedProject, ProjectGrantState, Role } from 'src/app/proto/generated/zitadel/project_pb';
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

    public grant!: GrantedProject.AsObject;
    public projectid: string = '';
    public grantid: string = '';

    public disabled: boolean = false;

    public isZitadel: boolean = false;
    ProjectGrantState: any = ProjectGrantState;

    public projectRoleOptions: Role.AsObject[] = [];
    public memberRoleOptions: Array<string> = [];

    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<Member.AsObject> = [];
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

            this.mgmtService.getProjectGrantByID(this.grantid, this.projectid).then((resp) => {
                if (resp.projectGrant) {
                    this.grant = resp.projectGrant;
                }
            });
        });
    }

    public changeState(newState: ProjectGrantState): void {
        if (newState === ProjectGrantState.PROJECT_GRANT_STATE_ACTIVE) {
            this.mgmtService.reactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.REACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (newState === ProjectGrantState.PROJECT_GRANT_STATE_INACTIVE) {
            this.mgmtService.deactivateProjectGrant(this.grantid, this.projectid).then(() => {
                this.toast.showInfo('PROJECT.TOAST.DEACTIVATED', true);
                this.grant.state = newState;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public getRoleOptions(projectId: string): void {
        this.mgmtService.listProjectRoles(projectId, 100, 0).then(resp => {
            this.projectRoleOptions = resp.resultList;
        });
    }

    public getMemberRoleOptions(): void {
        this.mgmtService.listProjectGrantMemberRoles().then(resp => {
            this.memberRoleOptions = resp.resultList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(selectionChange: MatSelectChange): void {
        this.mgmtService.updateProjectGrant(this.grant.grantId, this.grant.projectId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTUPDATED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            return this.mgmtService.removeProjectGrantMember(this.grant.projectId, this.grant.grantId, member.userId)
                .then(() => {
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
        const keysList = (await this.mgmtService.listProjectGrantMemberRoles());

        const dialogRef = this.dialog.open(ProjectGrantMembersCreateDialogComponent, {
            data: {
                roleKeysList: keysList.resultList,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe((dataToAdd: ProjectGrantMembersCreateDialogExportType) => {
            if (dataToAdd) {
                Promise.all(dataToAdd.userIds.map((userid: string) => {
                    return this.mgmtService.addProjectGrantMember(
                        this.grant.projectId,
                        this.grant.grantId,
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

    updateMemberRoles(member: Member.AsObject, selectionChange: MatSelectChange): void {
        this.mgmtService.updateProjectGrantMember(
            this.grant.projectId,
            this.grant.grantId,
            member.userId,
            selectionChange.value,
        ).then(() => {
            this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERCHANGED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
    }
}
