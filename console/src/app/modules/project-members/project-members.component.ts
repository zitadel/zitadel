import { SelectionModel } from '@angular/cdk/collections';
import { Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { take } from 'rxjs/operators';
import { ProjectGrantView, ProjectMember, ProjectType, ProjectView, User } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { CreationType, MemberCreateDialogComponent } from '../add-member-dialog/member-create-dialog.component';
import { ProjectMembersDataSource } from './project-members-datasource';

@Component({
    selector: 'app-project-members',
    templateUrl: './project-members.component.html',
    styleUrls: ['./project-members.component.scss'],
})
export class ProjectMembersComponent {
    public project!: ProjectView.AsObject | ProjectGrantView.AsObject;
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;
    public grantId: string = '';
    public projectName: string = '';
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectMember.AsObject>;
    public dataSource!: ProjectMembersDataSource;
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);
    public memberRoleOptions: string[] = [];

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'userId', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(
        private projectService: ProjectService,
        private dialog: MatDialog,
        private toast: ToastService,
        private route: ActivatedRoute) {
        this.route.data.pipe(take(1)).subscribe(data => {
            this.projectType = data.type;

            this.getRoleOptions();

            this.route.params.subscribe(params => {
                this.grantId = params.grantid;
                if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                    this.projectService.GetProjectById(params.projectid).then(project => {
                        this.project = project.toObject();
                        this.projectName = this.project.name;
                        this.dataSource = new ProjectMembersDataSource(this.projectService);
                        this.dataSource.loadMembers(this.project.projectId, this.projectType, 0, 25);
                    });
                } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                    console.log(params.projectid, params.grantid);
                    this.projectService.GetGrantedProjectByID(params.projectid, params.grantid).then(project => {
                        this.project = project.toObject();
                        this.projectName = this.project.projectName;
                        this.dataSource = new ProjectMembersDataSource(this.projectService);
                        this.dataSource.loadMembers(this.project.projectId, this.projectType, 0, 25, this.grantId);
                    });
                }
            });
        });
    }

    public getRoleOptions(): void {
        if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
            this.projectService.GetProjectGrantMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
            this.projectService.GetProjectMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
                return this.projectService.RemoveProjectMember(this.project.projectId, member.userId).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                return this.projectService.RemoveProjectGrantMember(this.project.projectId, this.grantId,
                    member.userId).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
            }
        }));
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.membersSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.membersSubject.value.forEach(row => this.selection.select(row));
    }

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.PROJECT_OWNED,
                projectId: this.project.projectId,
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
                            return this.projectService.AddProjectMember(this.project.projectId, user.id, roles);

                        } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
                            return this.projectService.AddProjectGrantMember(this.project.projectId, this.grantId,
                                user.id, roles);
                        }
                    })).then(() => {
                        this.toast.showInfo('PROJECT.TOAST.MEMBERSADDED', true);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
                }
            }
        });
    }

    updateRoles(member: ProjectMember.AsObject, selectionChange: MatSelectChange): void {
        if (this.projectType === ProjectType.PROJECTTYPE_OWNED) {
            this.projectService.ChangeProjectMember(this.project.projectId, member.userId, selectionChange.value)
                .then((newmember: ProjectMember) => {
                    this.toast.showInfo('Member changed');
                }).catch(error => {
                    this.toast.showError(error);
                });
        } else if (this.projectType === ProjectType.PROJECTTYPE_GRANTED) {
            this.projectService.ChangeProjectGrantMember(this.project.projectId,
                this.grantId, member.userId, selectionChange.value)
                .then((newmember: ProjectMember) => {
                    this.toast.showInfo('Member changed');
                }).catch(error => {
                    this.toast.showError(error);
                });
        }
    }
}
