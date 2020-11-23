import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatInput } from '@angular/material/input';
import { MatPaginator, PageEvent } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
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
    animations: [
        enterAnimations,
    ],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
    public userGrantSearchKey: UserGrantSearchKey | undefined = undefined;
    public UserGrantSearchKey: any = UserGrantSearchKey;

    public INITIAL_PAGE_SIZE: number = 50;
    @Input() context: UserGrantContext = UserGrantContext.NONE;
    @Input() refreshOnPreviousRoutes: string[] = [];

    public dataSource!: UserGrantsDataSource;
    public selection: SelectionModel<UserGrantView.AsObject> = new SelectionModel<UserGrantView.AsObject>(true, []);
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<UserGrantView.AsObject>;

    @Input() disableWrite: boolean = false;
    @Input() disableDelete: boolean = false;

    @Input() userId: string = '';
    @Input() projectId: string = '';
    @Input() grantId: string = '';
    @ViewChild('input') public filter!: MatInput;

    public grantRoleOptions: string[] = [];
    public projectRoleOptions: ProjectRoleView.AsObject[] = [];
    public routerLink: any = [''];

    public loadedGrantId: string = '';
    public loadedProjectId: string = '';
    public grantToEdit: string = '';

    public UserGrantContext: any = UserGrantContext;

    constructor(
        private userService: ManagementService,
        private mgmtService: ManagementService,
        private toast: ToastService,
    ) { }

    @Input() public displayedColumns: string[] = ['select',
        'user',
        'org',
        'projectId', 'dates', 'roleNamesList'];

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
            case UserGrantContext.NONE:
                this.routerLink = ['/grant-create'];
        }

        this.loadGrantsPage();
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadGrantsPage()),
            )
            .subscribe();
    }

    private loadGrantsPage(filterValue?: string): void {
        let queries: UserGrantSearchQuery[] = [];
        if (this.userGrantSearchKey !== undefined && filterValue) {
            const query = new UserGrantSearchQuery();
            query.setKey(this.userGrantSearchKey);
            query.setMethod(SearchMethod.SEARCHMETHOD_CONTAINS_IGNORE_CASE);
            query.setValue(filterValue);
            queries = [query];
        }

        this.dataSource.loadGrants(
            this.context,
            this.paginator?.pageIndex ?? 0,
            this.paginator?.pageSize ?? this.INITIAL_PAGE_SIZE,
            {
                projectId: this.projectId,
                grantId: this.grantId,
                userId: this.userId,
            },
            queries,
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

    public loadGrantOptions(grant: UserGrantView.AsObject): void {
        this.grantToEdit = grant.id;
        if (grant.grantId && grant.projectId) {
            this.getGrantRoleOptions(grant.grantId, grant.projectId);
        } else if (grant.projectId) {
            this.getProjectRoleOptions(grant.projectId);
        }
    }

    private getGrantRoleOptions(grantId: string, projectId: string): void {
        this.mgmtService.GetGrantedProjectByID(projectId, grantId).then(resp => {
            this.loadedGrantId = grantId;
            this.grantRoleOptions = resp.toObject().roleKeysList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    private getProjectRoleOptions(projectId: string): void {
        this.mgmtService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.loadedProjectId = projectId;
            this.projectRoleOptions = resp.toObject().resultList;
        });
    }

    updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
        this.userService.UpdateUserGrant(grant.id, grant.userId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
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

    public applyFilter(event: Event): void {
        this.selection.clear();
        const filterValue = (event.target as HTMLInputElement).value;

        this.loadGrantsPage(filterValue);
    }

    public setFilter(key: UserGrantSearchKey): void {
        setTimeout(() => {
            if (this.filter) {
                (this.filter as any).nativeElement.focus();
            }
        }, 100);

        if (this.userGrantSearchKey !== key) {
            this.userGrantSearchKey = key;
        } else {
            this.userGrantSearchKey = undefined;
            this.loadGrantsPage();
        }
    }
}
