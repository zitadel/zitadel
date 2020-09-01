import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import {
    ProjectRoleView,
    SearchMethod,
    UserGrant,
    UserGrantSearchKey,
    UserGrantSearchQuery,
    UserGrantView,
} from 'src/app/proto/generated/management_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { UserGrantContext, UserGrantsDataSource } from './user-grants-datasource';

@Component({
    selector: 'app-user-grants',
    templateUrl: './user-grants.component.html',
    styleUrls: ['./user-grants.component.scss'],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
    @Input() context: UserGrantContext = UserGrantContext.USER;
    public grants: UserGrantView.AsObject[] = [];

    public dataSource!: UserGrantsDataSource;
    public selection: SelectionModel<UserGrantView.AsObject> = new SelectionModel<UserGrantView.AsObject>(true, []);
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<UserGrantView.AsObject>;

    @Input() allowCreate: boolean = false;
    @Input() allowDelete: boolean = false;

    @Input() userId: string = '';
    @Input() projectId: string = '';
    @Input() grantId: string = '';

    public grantRoleOptions: string[] = [];
    public projectRoleOptions: ProjectRoleView.AsObject[] = [];
    public routerLink: any = [''];

    public loadedGrantId: string = '';
    public loadedProjectId: string = '';

    public UserGrantContext: any = UserGrantContext;

    constructor(
        private userService: ManagementService,
        private mgmtService: ManagementService,
        private toast: ToastService,
    ) { }

    @Input() public displayedColumns: string[] = ['select',
        'user',
        'org',
        'projectId', 'creationDate', 'changeDate', 'roleNamesList'];

    public ngOnInit(): void {
        this.dataSource = new UserGrantsDataSource(this.userService);

        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                if (this.projectId) {
                    this.getProjectRoleOptions(this.projectId);
                    this.routerLink = ['/grant-create', 'project', this.projectId];
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (this.grantId) {
                    this.routerLink = ['/grant-create', 'project', this.projectId, 'grant', this.grantId];
                    this.getGrantRoleOptions(this.grantId, this.projectId);
                }
                break;
            case UserGrantContext.USER:
                if (this.userId) {
                    this.routerLink = ['/grant-create', 'user', this.userId];
                }
                break;
            default:
                this.routerLink = ['/grant-create'];
        }

        this.dataSource.loadGrants(this.context, 0, 25, {
            projectId: this.projectId,
            grantId: this.grantId,
            userId: this.userId,
        });
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadGrantsPage()),
            )
            .subscribe();
    }

    private loadGrantsPage(): void {
        this.dataSource.loadGrants(
            this.context,
            this.paginator.pageIndex,
            this.paginator.pageSize,
            {
                projectId: this.projectId,
                grantId: this.grantId,
            },
        );
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.grantsSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.grantsSubject.value.forEach(row => this.selection.select(row));
    }

    public getGrantRoleOptions(grantId: string, projectId: string): void {
        this.mgmtService.GetGrantedProjectByID(projectId, grantId).then(resp => {
            this.loadedGrantId = projectId;
            this.grantRoleOptions = resp.toObject().roleKeysList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public getProjectRoleOptions(projectId: string): void {
        this.mgmtService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.loadedProjectId = projectId;
            this.projectRoleOptions = resp.toObject().resultList;
        });
    }

    updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                if (grant.id && grant.projectId) {
                    this.userService.UpdateUserGrant(grant.id, grant.userId, selectionChange.value)
                        .then(() => {
                            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (this.grantId && this.projectId) {
                    const projectQuery: UserGrantSearchQuery = new UserGrantSearchQuery();
                    projectQuery.setKey(UserGrantSearchKey.USERGRANTSEARCHKEY_PROJECT_ID);
                    projectQuery.setMethod(SearchMethod.SEARCHMETHOD_EQUALS);
                    projectQuery.setValue(this.projectId);
                    this.userService.UpdateUserGrant(
                        grant.id, grant.userId, selectionChange.value)
                        .then(() => {
                            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                }
                break;
        }
    }

    deleteGrantSelection(): void {
        this.userService.BulkRemoveUserGrant(this.selection.selected.map(grant => grant.id)).then(() => {
            this.toast.showInfo('GRANTS.TOAST.BULKREMOVED', true);
            const data = this.dataSource.grantsSubject.getValue();
            this.selection.selected.forEach((item) => {
                const index = data.findIndex(i => i.id === item.id);
                if (index > -1) {
                    data.splice(index, 1);
                    this.dataSource.grantsSubject.next(data);
                }
            });
            this.selection.clear();
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public changePage(event?: PageEvent): void {
        this.dataSource.loadGrants(
            this.context,
            event?.pageIndex ?? this.paginator.pageIndex,
            event?.pageSize ?? this.paginator.pageSize,
            {
                projectId: this.projectId,
                grantId: this.grantId,
                userId: this.userId,
            },
        );
    }
}
