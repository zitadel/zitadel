import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { User, UserMembershipSearchResponse, UserMembershipView, UserView } from 'src/app/proto/generated/management_pb';
import { AdminService } from 'src/app/services/admin.service';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { OrgService } from 'src/app/services/org.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { MembershipDetailDataSource } from './membership-detail-datasource';

@Component({
    selector: 'app-membership-detail',
    templateUrl: './membership-detail.component.html',
    styleUrls: ['./membership-detail.component.scss'],
})
export class MembershipDetailComponent implements AfterViewInit {
    public user!: UserView.AsObject;

    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<UserMembershipView.AsObject>;
    public dataSource!: MembershipDetailDataSource;
    public selection: SelectionModel<UserMembershipView.AsObject>
        = new SelectionModel<UserMembershipView.AsObject>(true, []);

    public memberRoleOptions: string[] = [];

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'memberType', 'displayName', 'creationDate', 'changeDate', 'roles'];

    public loading: boolean = false;
    public memberships!: UserMembershipSearchResponse.AsObject;

    constructor(
        private mgmtUserService: MgmtUserService,
        activatedRoute: ActivatedRoute,
        private dialog: MatDialog,
        private toast: ToastService,
        private projectService: ProjectService,
        private orgService: OrgService,
        private adminService: AdminService,
    ) {
        activatedRoute.params.subscribe(data => {
            const { id } = data;
            if (id) {
                this.mgmtUserService.GetUserByID(id).then(user => {
                    this.user = user.toObject();
                    this.dataSource = new MembershipDetailDataSource(this.mgmtUserService);
                    this.dataSource.loadMemberships(
                        this.user.id,
                        0,
                        50,
                    );
                }).catch(err => {
                    console.error(err);
                });
            }
        });
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadMembershipsPage()),
            )
            .subscribe();
    }

    private loadMembershipsPage(): void {
        this.dataSource.loadMemberships(
            this.user.id,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    // public removeSelectedMemberships(): void {
    //     Promise.all(this.selection.selected.map(membership => {
    //         switch (membership.memberType) {
    //             case MemberType.MEMBERTYPE_ORGANISATION:
    //                 return this.orgService.RemoveMyOrgMember(membership.objectId);
    //             case MemberType.MEMBERTYPE_PROJECT:
    //                 return this.projectService.RemoveProjectMember(membership.objectId, this.user.id);
    //             // case MemberType.MEMBERTYPE_PROJECT_GRANT:
    //             //     return this.projectService.RemoveProjectGrantMember(membership.objectId, this.user.id);
    //         }
    //     }));
    // }

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

    public addMember(): void {
        const dialogRef = this.dialog.open(MemberCreateDialogComponent, {
            width: '400px',
            data: {
                user: this.user,
            },
        });

        dialogRef.afterClosed().subscribe(resp => {
            if (resp && resp.creationType !== undefined) {
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

    public async loadManager(userId: string): Promise<void> {
        this.mgmtUserService.SearchUserMemberships(userId, 100, 0, []).then(response => {
            this.memberships = response.toObject();
            this.loading = false;
        });
    }

    public createIamMember(response: any): void {
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

    public refreshPage(): void {
        this.selection.clear();
        this.dataSource.loadMemberships(this.user.id, this.paginator.pageIndex, this.paginator.pageSize);
    }
}
