import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { Role } from 'src/app/proto/generated/zitadel/project_pb';
import { ManagementService } from 'src/app/services/mgmt.service';
import { ToastService } from 'src/app/services/toast.service';

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
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<Role.AsObject>;
    public dataSource!: ProjectRolesDataSource;
    public selection: SelectionModel<Role.AsObject> = new SelectionModel<Role.AsObject>(true, []);
    @Output() public changedSelection: EventEmitter<Array<Role.AsObject>> = new EventEmitter();

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'key', 'displayname', 'group', 'creationDate'];

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

    public deleteSelectedRoles(): Promise<any> {
        const oldState = this.dataSource.rolesSubject.value;
        const indexes = this.selection.selected.map(sel => {
            return oldState.findIndex(iter => iter.key === sel.key);
        });

        return Promise.all(this.selection.selected.map(role => {
            return this.mgmtService.removeProjectRole(this.projectId, role.key);
        })).then(() => {
            this.toast.showInfo('PROJECT.TOAST.ROLEREMOVED', true);
            indexes.forEach(index => {
                if (index > -1) {
                    oldState.splice(index, 1);
                    this.dataSource.rolesSubject.next(this.dataSource.rolesSubject.value);
                }
            });
            this.selection.clear();
        }).catch(error => {
            this.toast.showError(error);
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
