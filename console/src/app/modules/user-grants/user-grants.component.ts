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

    public roleOptions: ProjectRoleView.AsObject[] = [];
    public routerLink: any = [''];

    public loadedProjectId: string = '';
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
        this.dataSource = new UserGrantsDataSource(this.userService);
        const data = {
            projectId: this.projectId,
            grantId: this.grantId,
            userId: this.userId,
        };
        console.log(this.context);

        switch (this.context) {
            case UserGrantContext.OWNED_PROJECT:
                if (this.projectId) {
                    this.getRoleOptions(this.projectId);
                    this.routerLink = ['/grant-create', 'project', this.projectId];
                }
                break;
            case UserGrantContext.GRANTED_PROJECT:
                if (data && data.grantId) {
                    this.routerLink = ['/grant-create', 'project', this.projectId, 'grant', this.grantId];
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
        console.log(data);
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

    public loadRoleOptions(projectId: string): void {
        if (this.context === UserGrantContext.USER) {
            this.getRoleOptions(projectId);
        }
    }

    public getRoleOptions(projectId: string): void {
        this.projectService.SearchProjectRoles(projectId, 100, 0).then(resp => {
            this.loadedProjectId = projectId;
            this.roleOptions = resp.toObject().resultList;
        });
    }

    updateRoles(grant: UserGrant.AsObject, selectionChange: MatSelectChange): void {
        this.userService.UpdateUserGrant(grant.id, grant.userId, selectionChange.value)
            .then((newmember: UserGrant) => {
                this.toast.showInfo('Grant updated!');
            }).catch(error => {
                this.toast.showError(error.message);
            });
    }

    deleteGrantSelection(): void {
        this.userService.BulkRemoveUserGrant(this.selection.selected.map(grant => grant.id)).then(() => {
            this.toast.showInfo('Grants deleted');
            // this.loadGrantsPage();
            const data = this.dataSource.grantsSubject.getValue();
            console.log(data);
            this.selection.selected.forEach((item) => {
                console.log(item);
                const index = data.findIndex(i => i.id === item.id);
                if (index > -1) {
                    data.splice(index, 1);
                    this.dataSource.grantsSubject.next(data);
                }
            });
            this.selection.clear();
        }).catch(error => {
            this.toast.showError(error.message);
        });
    }
}
