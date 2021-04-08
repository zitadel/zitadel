import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { ActivatedRoute } from '@angular/router';
import { tap } from 'rxjs/operators';
import { CreationType, MemberCreateDialogComponent } from 'src/app/modules/add-member-dialog/member-create-dialog.component';
import { PaginatorComponent } from 'src/app/modules/paginator/paginator.component';
import { Membership, User } from 'src/app/proto/generated/zitadel/user_pb';
import { AdminService } from 'src/app/services/admin.service';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { MembershipDetailDataSource } from './membership-detail-datasource';

@Component({
    selector: 'app-membership-detail',
    templateUrl: './membership-detail.component.html',
    styleUrls: ['./membership-detail.component.scss'],
})
export class MembershipDetailComponent implements AfterViewInit {
    public user!: User.AsObject;

    @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
    @ViewChild(MatTable) public table!: MatTable<Membership.AsObject>;
    public dataSource!: MembershipDetailDataSource;
    public selection: SelectionModel<Membership.AsObject>
        = new SelectionModel<Membership.AsObject>(true, []);

    public memberRoleOptions: string[] = [];

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'memberType', 'displayName', 'creationDate', 'changeDate', 'roles', 'actions'];

    public loading: boolean = false;
    public memberships!: Membership.AsObject[];

    constructor(
        activatedRoute: ActivatedRoute,
        private dialog: MatDialog,
        private toast: ToastService,
        private mgmtService: ManagementService,
        private adminService: AdminService,
    ) {
        activatedRoute.params.subscribe(data => {
            const { id } = data;
            if (id) {
                this.mgmtService.getUserByID(id).then(resp => {
                    if (resp.user) {
                        this.user = resp.user;
                        this.dataSource = new MembershipDetailDataSource(this.mgmtService);
                        this.dataSource.loadMemberships(
                            this.user.id,
                            0,
                            50,
                        );
                    }
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
        this.mgmtService.listUserMemberships(userId, 100, 0, []).then(response => {
            this.memberships = response.resultList;
            this.loading = false;
        });
    }

    public createIamMember(response: any): void {
        const users: User.AsObject[] = response.users;
        const roles: string[] = response.roles;

        if (users && users.length && roles && roles.length) {
            Promise.all(users.map(user => {
                return this.adminService.addIAMMember(user.id, roles);
            })).then(() => {
                this.toast.showInfo('IAM.TOAST.MEMBERADDED', true);
                setTimeout(() => {
                    this.refreshPage();
                }, 1000);
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
                return this.mgmtService.addOrgMember(user.id, roles);
            })).then(() => {
                this.toast.showInfo('ORG.TOAST.MEMBERADDED', true);
                setTimeout(() => {
                    this.refreshPage();
                }, 1000);
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
                return this.mgmtService.addProjectGrantMember(
                    response.projectId,
                    response.grantId,
                    user.id,
                    roles,
                ).then(() => {
                    this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                    setTimeout(() => {
                        this.refreshPage();
                    }, 1000);
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
                return this.mgmtService.addProjectMember(response.projectId, user.id, roles)
                    .then(() => {
                        this.toast.showInfo('PROJECT.TOAST.MEMBERADDED', true);
                        setTimeout(() => {
                            this.refreshPage();
                        }, 1000);
                    }).catch(error => {
                        this.toast.showError(error);
                    });
            });
        }
    }

    public removeMembership(membership: Membership.AsObject): void {
        let prom;

        if (membership.projectId && membership.projectGrantId && membership.userId) {
            prom = this.mgmtService.removeProjectGrantMember(membership.projectId, membership.projectGrantId, membership.userId);
        } else if (membership.projectId && membership.userId) {
            prom = this.mgmtService.removeProjectMember(membership.projectId, membership.userId);
        } else if (membership.orgId && membership.userId) {
            prom = this.mgmtService.removeOrgMember(membership.userId);
        } else if (membership.userId) {
            prom = this.adminService.removeIAMMember(membership.userId);
        }

        if (prom) {
            prom.then(() => {
                this.toast.showInfo('PROJECT.TOAST.MEMBERREMOVED', true);
                this.refreshPage();
            }).catch(error => this.toast.showError(error));
        }
    }

    public refreshPage(): void {
        this.selection.clear();
        this.dataSource.loadMemberships(this.user.id, this.paginator.pageIndex, this.paginator.pageSize);
    }
}
