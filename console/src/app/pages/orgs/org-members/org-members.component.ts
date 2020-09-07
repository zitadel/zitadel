import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { Org, OrgMemberView, ProjectType, UserView } from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { OrgMembersDataSource } from './org-members-datasource';

@Component({
    selector: 'app-org-members',
    templateUrl: './org-members.component.html',
    styleUrls: ['./org-members.component.scss'],
})
export class OrgMembersComponent implements AfterViewInit {
    public org!: Org.AsObject;
    public projectType: ProjectType = ProjectType.PROJECTTYPE_OWNED;
    public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    public dataSource!: OrgMembersDataSource;
    public selection: SelectionModel<OrgMemberView.AsObject> = new SelectionModel<OrgMemberView.AsObject>(true, []);

    public memberRoleOptions: string[] = [];

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(
        private mgmtService: ManagementService,
        private dialog: MatDialog,
        private toast: ToastService,
    ) {
        this.mgmtService.GetMyOrg().then(org => {
            this.org = org.toObject();
            this.dataSource = new OrgMembersDataSource(this.mgmtService);
            this.dataSource.loadMembers(0, 25);
        });

        this.getRoleOptions();
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadMembersPage()),
            )
            .subscribe();
    }

    public getRoleOptions(): void {
        this.mgmtService.GetOrgMemberRoles().then(resp => {
            this.memberRoleOptions = resp.toObject().rolesList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    updateRoles(member: OrgMemberView.AsObject, selectionChange: MatSelectChange): void {
        this.mgmtService.ChangeMyOrgMember(member.userId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERCHANGED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    private loadMembersPage(): void {
        this.dataSource.loadMembers(
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public removeOrgMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.mgmtService.RemoveMyOrgMember(member.userId).then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERREMOVED', true);
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

    public openAddMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            data: {
                creationType: CreationType.ORG,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: UserView.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.mgmtService.AddMyOrgMember(user.id, roles);
                    })).then(() => {
                        this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
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
