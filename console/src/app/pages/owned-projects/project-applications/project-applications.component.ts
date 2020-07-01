import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTable } from '@angular/material/table';
import { merge } from 'rxjs';
import { tap } from 'rxjs/operators';
import { Application } from 'src/app/proto/generated/management_pb';
import { ProjectService } from 'src/app/services/project.service';
import { ToastService } from 'src/app/services/toast.service';

import { ProjectApplicationsDataSource } from './project-applications-datasource';


@Component({
    selector: 'app-project-applications',
    templateUrl: './project-applications.component.html',
    styleUrls: ['./project-applications.component.scss'],
})
export class ProjectApplicationsComponent implements AfterViewInit, OnInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatSort) public sort!: MatSort;
    @ViewChild(MatTable) public table!: MatTable<Application.AsObject>;
    public dataSource!: ProjectApplicationsDataSource;
    public selection: SelectionModel<Application.AsObject> = new SelectionModel<Application.AsObject>(true, []);

    /** Columns displayed in the table. Columns IDs can be added, removed, or reordered. */
    public displayedColumns: string[] = ['select', 'name'];

    constructor(private projectService: ProjectService, private toast: ToastService) { }

    public ngOnInit(): void {
        this.dataSource = new ProjectApplicationsDataSource(this.projectService);
        this.dataSource.loadApps(this.projectId, 0, 25, 'asc');
    }

    public ngAfterViewInit(): void {
        this.sort.sortChange.subscribe(() => this.paginator.pageIndex = 0);
        merge(this.sort.sortChange, this.paginator.page)
            .pipe(
                tap(() => this.loadRolesPage()),
            )
            .subscribe();

    }

    private loadRolesPage(): void {
        this.dataSource.loadApps(
            this.projectId,
            this.paginator.pageIndex,
            this.paginator.pageSize,
            this.sort.direction,
        );
    }

    public isAllSelected(): boolean {
        const numSelected = this.selection.selected.length;
        const numRows = this.dataSource.appsSubject.value.length;
        return numSelected === numRows;
    }

    public masterToggle(): void {
        this.isAllSelected() ?
            this.selection.clear() :
            this.dataSource.appsSubject.value.forEach((row: Application.AsObject) => this.selection.select(row));
    }
}
