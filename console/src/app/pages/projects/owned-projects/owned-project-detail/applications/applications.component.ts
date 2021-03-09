import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTable } from '@angular/material/table';
import { merge, of } from 'rxjs';
import { tap } from 'rxjs/operators';
import { App } from 'src/app/proto/generated/zitadel/app_pb';
import { ManagementService } from 'src/app/services/mgmt.service';

import { ProjectApplicationsDataSource } from './applications-datasource';


@Component({
    selector: 'app-applications',
    templateUrl: './applications.component.html',
    styleUrls: ['./applications.component.scss'],
})
export class ApplicationsComponent implements AfterViewInit, OnInit {
    @Input() public projectId: string = '';
    @Input() public disabled: boolean = false;
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatSort) public sort!: MatSort;
    @ViewChild(MatTable) public table!: MatTable<App.AsObject>;
    public dataSource!: ProjectApplicationsDataSource;
    public selection: SelectionModel<App.AsObject> = new SelectionModel<App.AsObject>(true, []);

    public displayedColumns: string[] = ['select', 'name', 'type'];

    constructor(private mgmtService: ManagementService) { }

    public ngOnInit(): void {
        this.dataSource = new ProjectApplicationsDataSource(this.mgmtService);
        this.dataSource.loadApps(this.projectId, 0, 25);
    }

    public ngAfterViewInit(): void {
        merge(this.sort ? this.sort?.sortChange : of(null), this.paginator.page)
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
            this.dataSource.appsSubject.value.forEach((row: App.AsObject) => this.selection.select(row));
    }

    public refreshPage(): void {
        this.dataSource.loadApps(this.projectId, this.paginator.pageIndex, this.paginator.pageSize);
    }
}
