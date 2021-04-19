import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatInput } from '@angular/material/input';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { enterAnimations } from 'src/app/animations';
import { TextQueryMethod } from 'src/app/proto/generated/zitadel/object_pb';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import {
    UserGrant,
    UserGrantDisplayNameQuery,
    UserGrantOrgNameQuery,
    UserGrantProjectNameQuery,
    UserGrantQuery,
    UserGrantRoleKeyQuery,
} from 'src/app/proto/generated/zitadel/user_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PageEvent, PaginatorComponent } from '../paginator/paginator.component';
import { UserGrantContext, UserGrantsDataSource } from './user-grants-datasource';

export enum UserGrantListSearchKey {
    DISPLAY_NAME,
    ORG_NAME,
    PROJECT_NAME,
    ROLE_KEY,
}
@Component({
    selector: 'app-user-grants',
    templateUrl: './user-grants.component.html',
    styleUrls: ['./user-grants.component.scss'],
    animations: [
        enterAnimations,
    ],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
    public userGrantListSearchKey: UserGrantListSearchKey | undefined = undefined;
    public UserGrantListSearchKey: any = UserGrantListSearchKey;

    public INITIAL_PAGE_SIZE: number = 50;
    @Input() context: UserGrantContext = UserGrantContext.NONE;
    @Input() refreshOnPreviousRoutes: string[] = [];

    public dataSource!: UserGrantsDataSource;
    public selection: SelectionModel<UserGrant.AsObject> = new SelectionModel<UserGrant.AsObject>(true, []);
    @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
    @ViewChild(MatTable) public table!: MatTable<UserGrant.AsObject>;

    @Input() disableWrite: boolean = false;
    @Input() disableDelete: boolean = false;

    @Input() userId: string = '';
    @Input() projectId: string = '';
    @Input() grantId: string = '';
    @ViewChild('input') public filter!: MatInput;

    public grantRoleOptions: string[] = [];
    public projectRoleOptions: Role.AsObject[] = [];
    public routerLink: any = [''];

    public loadedId: string = '';
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
        let queries: UserGrantQuery[] = [];
        if (this.userGrantListSearchKey !== undefined && filterValue) {
            const query = new UserGrantQuery();
            switch (this.userGrantListSearchKey) {
                case UserGrantListSearchKey.DISPLAY_NAME:
                    const ugDnQ = new UserGrantDisplayNameQuery();
                    ugDnQ.setDisplayName(filterValue);
                    ugDnQ.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
                    query.setDisplayNameQuery(ugDnQ);
                    break;
                case UserGrantListSearchKey.ORG_NAME:
                    const ugOnQ = new UserGrantOrgNameQuery();
                    ugOnQ.setOrgName(filterValue);
                    ugOnQ.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
                    query.setOrgNameQuery(ugOnQ);
                    break;
                case UserGrantListSearchKey.PROJECT_NAME:
                    const ugPnQ = new UserGrantProjectNameQuery();
                    ugPnQ.setProjectName(filterValue);
                    ugPnQ.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
                    query.setProjectNameQuery(ugPnQ);
                    break;
                case UserGrantListSearchKey.ROLE_KEY:
                    const ugRkQ = new UserGrantRoleKeyQuery();
                    ugRkQ.setRoleKey(filterValue);
                    ugRkQ.setMethod(TextQueryMethod.TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE);
                    query.setRoleKeyQuery(ugRkQ);
                    break;

            }
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

    public loadGrantOptions(grant: UserGrant.AsObject): void {
        this.grantToEdit = grant.id;
        if (grant.projectGrantId && grant.projectId) {
            this.getGrantRoleOptions(grant.projectGrantId, grant.projectId);
        } else if (grant.projectId) {
            this.getProjectRoleOptions(grant.projectId);
        }
    }

    private getGrantRoleOptions(id: string, projectId: string): void {
        this.mgmtService.getGrantedProjectByID(projectId, id).then(resp => {
            if (resp.grantedProject) {
                this.loadedId = id;
                this.grantRoleOptions = resp.grantedProject?.grantedRoleKeysList;
            }
        }).catch(error => {
            this.grantToEdit = '';
            this.toast.showError(error);
        });
    }

    private getProjectRoleOptions(projectId: string): void {
        this.mgmtService.listProjectRoles(projectId, 100, 0).then(resp => {
            this.loadedProjectId = projectId;
            this.projectRoleOptions = resp.resultList;
        });
    }

    updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
        this.userService.updateUserGrant(grant.id, grant.userId, selectionChange.value)
            .then(() => {
                this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
            }).catch(error => {
                this.toast.showError(error);
            });
    }

    deleteGrantSelection(): void {
        this.userService.bulkRemoveUserGrant(this.selection.selected.map(grant => grant.id)).then(() => {
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

    public setFilter(key: UserGrantListSearchKey): void {
        setTimeout(() => {
            if (this.filter) {
                (this.filter as any).nativeElement.focus();
            }
        }, 100);

        if (this.userGrantListSearchKey !== key) {
            this.userGrantListSearchKey = key;
        } else {
            this.userGrantListSearchKey = undefined;
            this.loadGrantsPage();
        }
    }
}
