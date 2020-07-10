import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSelectChange } from '@angular/material/select';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { ProjectGrant, ProjectRoleView, UserGrant } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { UserGrantContext, UserGrantsDataSource } from './user-grants-datasource';

@Component({
    selector: 'app-user-grants',
    templateUrl: './user-grants.component.html',
    styleUrls: ['./user-grants.component.scss'],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
    // @Input() filterValue: string = '';
    // @Input() filter: UserGrantSearchKey = UserGrantSearchKey.USERGRANTSEARCHKEY_USER_ID;
    @Input() context: UserGrantContext = UserGrantContext.USER;
    public grants: UserGrant.AsObject[] = [];

    public dataSource!: UserGrantsDataSource;
    public selection: SelectionModel<UserGrant.AsObject> = new SelectionModel<UserGrant.AsObject>(true, []);
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectGrant.AsObject>;

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
        private userService: MgmtUserService,
        private projectService: ProjectService,
        private toast: ToastService,
    ) { }

    public displayedColumns: string[] = ['select',
        'user',
        'org',
        'projectId', 'creationDate', 'changeDate', 'roleNamesList'];

    public ngOnInit(): void {
        console.log(this.context);
        this.dataSource = new UserGrantsDataSource(this.userService);
        const data = {
            projectId: this.projectId,
            grantId: this.grantId,
            userId: this.userId,
        };

        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                if (this.projectId) {
                    this.getProjectRoleOptions(this.projectId);
                    this.routerLink = ['/grant-create', 'project', this.projectId];
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (data && data.grantId) {
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
        this.dataSource.loadGrants(this.context, 0, 25, data);
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
        console.log(grantId, projectId);

        this.projectService.GetGrantedProjectByID(projectId, grantId).then(resp => {
            this.loadedGrantId = projectId;
            this.grantRoleOptions = resp.toObject().roleKeysList;
        }).catch(error => {
            this.toast.showError(error);
        });
    }

    public getProjectRoleOptions(projectId: string): void {
        console.log(projectId);
        this.projectService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.loadedProjectId = projectId;
            this.projectRoleOptions = resp.toObject().resultList;
        });
    }

    updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                if (grant.id && grant.projectId) {
                    this.userService.UpdateProjectUserGrant(grant.id, grant.projectId, grant.userId, selectionChange.value)
                        .then(() => {
                            this.toast.showInfo('GRANTS.TOAST.UPDATED', true);
                        }).catch(error => {
                            this.toast.showError(error);
                        });
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (this.grantId && this.projectId) {
                    this.userService.updateProjectGrantUserGrant(grant.id,
                        this.grantId, grant.userId, selectionChange.value)
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
}
