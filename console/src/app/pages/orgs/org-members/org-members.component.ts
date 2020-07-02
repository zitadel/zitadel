import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { Org, ProjectMember, ProjectType, User } from 'src/app/proto/generated/management_pb';
import { OrgService } from 'src/app/services/org.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectMembersDataSource } from './org-members-datasource';

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
    @ViewChild(MatTable) public table!: MatTable<ProjectMember.AsObject>;
    public dataSource!: ProjectMembersDataSource;
    public selection: SelectionModel<ProjectMember.AsObject> = new SelectionModel<ProjectMember.AsObject>(true, []);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'firstname', 'lastname', 'username', 'email', 'roles'];

    constructor(private orgService: OrgService,
        private dialog: MatDialog,
        private toast: ToastService) {
        this.orgService.GetMyOrg().then(org => {
            this.org = org.toObject();
            this.dataSource = new ProjectMembersDataSource(this.orgService);
            this.dataSource.loadMembers(0, 25, 'asc');
        });
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

    public removeProjectMemberSelection(): void {
        Promise.all(this.selection.selected.map(member => {
            return this.orgService.RemoveMyOrgMember(member.userId).then(() => {
                this.toast.showInfo('Removed successfully');
            }).catch(error => {
                this.toast.showError(error.message);
            });
        }));
    }

    public removeMember(member: ProjectMember.AsObject): void {
        this.orgService.RemoveMyOrgMember(member.userId).then(() => {
            this.toast.showInfo('Member removed successfully');
        }).catch(error => {
            this.toast.showError(error.message);
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
                creationType: CreationType.ORG,
            },
            width: '400px',
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp) {
                const users: User.AsObject[] = resp.users;
                const roles: string[] = resp.roles;

                if (users && users.length && roles && roles.length) {
                    Promise.all(users.map(user => {
                        return this.orgService.AddMyOrgMember(user.id, roles);
                    })).then(() => {
                        this.toast.showError('members added');
                    }).catch(error => {
                        this.toast.showError(error.message);
                    });
                }
            }
        });
    }
}
