import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

import { PaginatorComponent } from '../paginator/paginator.component';
import { ProjectRoleDetailComponent } from './project-role-detail/project-role-detail.component';
import { ProjectRolesDataSource } from './project-roles-datasource';


@Component({
    selector: 'app-project-roles',
    templateUrl: './project-roles.component.html',
    styleUrls: ['./project-roles.component.scss'],
})
export class ProjectRolesComponent implements AfterViewInit, OnInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @Input() public actionsVisible: boolean = false;
    @ViewChild(PaginatorComponent) public paginator!: PaginatorComponent;
    @ViewChild(MatTable) public table!: MatTable<Role.AsObject>;
    public dataSource!: ProjectRolesDataSource;
    public selection: SelectionModel<Role.AsObject> = new SelectionModel<Role.AsObject>(true, []);
    @Output() public changedSelection: EventEmitter<Array<Role.AsObject>> = new EventEmitter();

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'key', 'displayname', 'group', 'creationDate', 'actions'];

    constructor(private mgmtService: ManagementService, private toast: ToastService, private dialog: MatDialog) {
        this.dataSource = new ProjectRolesDataSource(this.mgmtService);
    }

    public ngOnInit(): void {
        this.dataSource.loadRoles(this.projectId, 0, 25, 'asc');

        this.selection.changed.subscribe(() => {
            this.changedSelection.emit(this.selection.selected);
        });
    }

    public ngAfterViewInit(): void {
        this.paginator.page
            .pipe(
                tap(() => this.loadRolesPage()),
            )
            .subscribe();
    }

    public selectAllOfGroup(group: string): void {
        const groupRoles: Role.AsObject[] = this.dataSource.rolesSubject.getValue()
            .filter(role => role.group === group);
        this.selection.select(...groupRoles);
    }

    private loadRolesPage(): void {
        this.dataSource.loadRoles(
            this.projectId,
            this.paginator.pageIndex,
            this.paginator.pageSize,
        );
    }

    public changePage(): void {
        this.selection.clear();
        this.loadRolesPage();
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.rolesSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.rolesSubject.value.forEach((row: Role.AsObject) => this.selection.select(row));
    }

    public deleteRole(role: Role.AsObject): Promise<any> {
        const index = this.dataSource.rolesSubject.value.findIndex(iter => iter.key === role.key);

        return this.mgmtService.removeProjectRole(this.projectId, role.key).then(() => {
            this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);

            if (index > -1) {
                this.dataSource.rolesSubject.value.splice(index, 1);
                this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
            }
        });
    }

    public removeRole(role: Role.AsObject, index: number): void {
        this.mgmtService
            .removeProjectRole(this.projectId, role.key)
            .then(() => {
                this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);
                this.dataSource.rolesSubject.value.splice(index, 1);
                this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
            })
            .catch(error => {
                this.toast.showError(error);
            });
    }

    public openDetailDialog(role: Role.AsObject): void {
        this.dialog.open(ProjectRoleDetailComponent, {
            data: {
                role,
                projectId: this.projectId,
            },
            width: '400px',
        });
    }

    public refreshPage(): void {
        this.dataSource.loadRoles(this.projectId, this.paginator.pageIndex, this.paginator.pageSize);
    }
}
