import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { ProjectMember, ProjectType } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import {
    ProjectGrantMembersCreateDialogComponent,
    ProjectGrantMembersCreateDialogExportType,
} from './project-grant-members-create-dialog/project-grant-members-create-dialog.component';
import { ProjectGrantMembersDataSource } from './project-grant-members-datasource';

@Component({
    selector: 'app-project-grant-members',
    templateUrl: './project-grant-members.component.html',
    styleUrls: ['./project-grant-members.component.scss'],
})
export class ProjectGrantMembersComponent implements AfterViewInit, OnInit {
    @Input() public projectId!: string;
    @Input() public grantId!: string;

    @Input() public type: ProjectType = ProjectType.PROJECTTYPE_GRANTED;

    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectMember.AsObject>;
    public dataSource!: ProjectGrantMembersDataSource;
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    public ProjectType: any = ProjectType;
    public memberRoleOptions: string[] = [];

    constructor(
        private mgmtService: ManagementService,
        private dialog: MatDialog,
        private toast: ToastService,
    ) {
        this.dataSource = new ProjectGrantMembersDataSource(this.mgmtService);
        this.getRoleOptions();
    }

    public ngOnInit(): void {
        this.dataSource.loadGrantMembers(this.projectId, this.grantId, 0, 25);
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadMembersPage()),
            )
            .subscribe();
    }

    public getRoleOptions(): void {
        if (this.type === ProjectType.PROJECTTYPE_GRANTED) {
            this.mgmtService.GetProjectGrantMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                this.toast.showError(error);
            });
        } else if (this.type === ProjectType.PROJECTTYPE_OWNED) {
            this.mgmtService.GetProjectMemberRoles().then(resp => {
                this.memberRoleOptions = resp.toObject().rolesList;
            }).catch(error => {
                this.toast.showError(error);
            });
        }
    }

    private loadMembersPage(): void {
        this.dataSource.loadGrantMembers(
            this.projectId,
            this.grantId,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.mgmtService.RemoveProjectGrantMember(this.projectId, this.grantId, member.userId).then(() => {
                this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERREMOVED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
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
                        this.projectId,
                        this.grantId,
                        userid,
                        dataToAdd.rolesKeyList,
                    );
                })).then(() => {
                    this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERADDED', true);
                }).catch(error => {
                    this.toast.showError(error);
                });
            }
        });
    }

    updateRoles(member: ProjectMember.AsObject, selectionChange: MatSelectChange): void {
        this.mgmtService.ChangeProjectGrantMember(this.projectId, this.grantId, member.userId, selectionChange.value)
            .then((newmember: ProjectMember) => {
                this.toast.showInfo('PROJECT.GRANT.TOAST.PROJECTGRANTMEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    public refreshPage(): void {
        this.selection.clear();
        this.dataSource.loadGrantMembers(this.projectId, this.grantId, this.paginator.pageIndex, this.paginator.pageSize);
    }
}
