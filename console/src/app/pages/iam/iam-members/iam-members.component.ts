import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { IamMember, IamMemberView } from 'src/app/proto/generated/admin_pb';
import { ProjectMember, ProjectType, UserView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ToastService } from 'src/app/services/toast.service';

import { IamMembersDataSource } from './iam-members-datasource';

@Component({
    selector: 'app-iam-members',
    templateUrl: './iam-members.component.html',
    styleUrls: ['./iam-members.component.scss'],
})
export class IamMembersComponent implements AfterViewInit {
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<IamMemberView.AsObject>;
    public dataSource!: IamMembersDataSource;
    public selection: SelectionModel<IamMemberView.AsObject> = new SelectionModel<IamMemberView.AsObject>(true, []);

    public memberRoleOptions: string[] = [];
    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(private adminService: AdminService,
        private dialog: MatDialog,
        private toast: ToastService) {

        this.dataSource = new IamMembersDataSource(this.adminService);
        this.dataSource.loadMembers(0, 25);
        this.getRoleOptions();
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadMembersPage()),
            )
            .subscribe();
    }

    private loadMembersPage(): void {
        this.dataSource.loadMembers(
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public getRoleOptions(): void {
        this.adminService.GetIamMemberRoles().then(resp => {
            this.memberRoleOptions = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(member: IamMemberView.AsObject, selectionChange: MatSelectChange): void {
        this.adminService.ChangeIamMember(member.userId, selectionChange.value)
            .then((newmember: IamMember) => {
                this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }


    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.adminService.RemoveIamMember(member.userId).then(() => {
                this.toast.showInfo('IAM.TOAST.MEMBERREMOVED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
        }));
    }

    public removeMember(member: ProjectMember.AsObject): void {
        this.adminService.RemoveIamMember(member.userId).then(() => {
            this.toast.showInfo('IAM.TOAST.MEMBERREMOVED', true);
        }).catch(error => {
            this.toast.showError(error);
        });
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
                creationType: CreationType.IAM,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

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
        });
    }

    public refreshPage(): void {
        this.selection.clear();
        this.dataSource.loadMembers(this.paginator.pageIndex, this.paginator.pageSize);
    }
}
