import { Component, EventEmitter } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs/operators';
import { Member } from 'src/app/proto/generated/zitadel/member_pb';
import { GrantedProject, Project } from 'src/app/proto/generated/zitadel/project_pb';
import { User } from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { CreationType, MemberCreateDialogComponent } from '../add-member-dialog/member-create-dialog.component';
import { ProjectMembersDataSource, ProjectType } from './project-members-datasource';

@Component({
    selector: 'app-project-members',
    templateUrl: './project-members.component.html',
    styleUrls: ['./project-members.component.scss'],
})
export class ProjectMembersComponent {
    public INITIALPAGESIZE: number = 25;
    public project!: Project.AsObject | GrantedProject.AsObject;
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public grantId: string = '';
    public projectName: string = '';
    public dataSource!: ProjectMembersDataSource;
    public memberRoleOptions: string[] = [];

    public changePageFactory!: Function;
    public changePage: EventEmitter<void> = new EventEmitter();
    public selection: Array<Member.AsObject> = [];

    public ProjectType: any = ProjectType;
    constructor(
        private mgmtService: ManagementService,
        private dialog: MatDialog,
        private toast: ToastService,
        private route: ActivatedRoute) {
        this.route.data.pipe(take(1)).subscribe(data => {
            this.projectType = data.type;

            this.getRoleOptions();

            this.route.params.subscribe(params => {
                this.grantId = params.grantid;
                if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                    this.mgmtService.getProjectByID(params.projectid).then(resp => {
                        if (resp.project) {
                            this.project = resp.project;
                            this.projectName = this.project.name;
                            this.dataSource = new ProjectMembersDataSource(this.mgmtService);
                            this.dataSource.loadMembers(this.project.id, this.projectType, 0, this.INITIALPAGESIZE);

                            this.changePageFactory = (event?: PageEvent) => {
                                return this.dataSource.loadMembers(
                                    (this.project as Project.AsObject).id,
                                    this.projectType,
                                    event?.pageIndex ?? 0,
                                    event?.pageSize ?? this.INITIALPAGESIZE,
                                    this.grantId,
                                );
                            };
                        }
                    });
                } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                    this.mgmtService.getGrantedProjectByID(params.projectid, params.grantid).then(resp => {
                        if (resp.grantedProject) {
                            this.project = resp.grantedProject;
                            this.projectName = this.project.projectName;
                            this.dataSource = new ProjectMembersDataSource(this.mgmtService);
                            this.dataSource.loadMembers(this.project.projectId,
                                this.projectType,
                                0,
                                this.INITIALPAGESIZE,
                                this.grantId,
                            );

                            this.changePageFactory = (event?: PageEvent) => {
                                return this.dataSource.loadMembers(
                                    (this.project as GrantedProject.AsObject).projectId,
                                    this.projectType,
                                    event?.pageIndex ?? 0,
                                    event?.pageSize ?? this.INITIALPAGESIZE,
                                    this.grantId,
                                );
                            };
                        }
                    });
                }
            });
        });
    }

    public getRoleOptions(): void {
        if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
            this.mgmtService.listProjectGrantMemberRoles().then(resp => {
                this.memberRoleOptions = resp.resultList;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
            this.mgmtService.listProjectMemberRoles().then(resp => {
                this.memberRoleOptions = resp.resultList;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.map(member => {
            if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                return this.mgmtService.removeProjectMember((this.project as Project.AsObject).id, member.userId)
                    .then(() => {
                        this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
            } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                return this.mgmtService.removeProjectGrantMember(
                    (this.project as GrantedProject.AsObject).projectId,
                    this.grantId,
                    member.userId,
                ).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        })).then(() => {
            setTimeout(() => {
                this.changePage.emit();
            }, 1000);
        });
    }

    public removeProjectMember(member: Member.AsObject | Member.AsObject): void {
        if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
            this.mgmtService.removeProjectMember((this.project as Project.AsObject).id, member.userId).then(() => {
                setTimeout(() => {
                    this.changePage.emit();
                }, 1000);
                this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
            this.mgmtService.removeProjectGrantMember((this.project as GrantedProject.AsObject).projectId, this.grantId,
                member.userId).then(() => {
                    setTimeout(() => {
                        this.changePage.emit();
                    }, 1000);
                    this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.PROJECT_OWNED,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                            return this.mgmtService.addProjectMember((this.project as Project.AsObject).id, user.id, roles);

                        } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                            return this.mgmtService.addProjectGrantMember(
                                (this.project as GrantedProject.AsObject).projectId,
                                this.grantId,
                                user.id,
                                roles,
                            );
                        }
                    })).then(() => {
                        setTimeout(() => {
                            this.changePage.emit();
                        }, 1000);
                        this.toast.showInfo('PROJECT.TOAST.MEMBERSADDED', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }

    updateRoles(member: Member.AsObject, selectionChange: MatSelectChange): void {
        if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
            this.mgmtService.updateProjectMember((this.project as Project.AsObject).id, member.userId, selectionChange.value)
                .then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERCHANGED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
        } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
            this.mgmtService.updateProjectGrantMember((this.project as GrantedProject.AsObject).projectId,
                this.grantId, member.userId, selectionChange.value)
                .then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERCHANGED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }
}
