import { SelectionModel } from '@angular/cdk/collections';
import { AfterViewInit, Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTable } from '@angular/material/table';
import { tap } from 'rxjs/operators';
import { ProjectGrant, UserGrant } from 'src/app/proto/generated/management_pb';
import { MgmtUserService } from 'src/app/services/mgmt-user.service';

import { UserGrantsDataSource } from './user-grants-datasource';

@Component({
    selector: 'app-user-grants',
    templateUrl: './user-grants.component.html',
    styleUrls: ['./user-grants.component.scss'],
})
export class UserGrantsComponent implements OnInit, AfterViewInit {
    @Input() userId: string = '';
    public grants: UserGrant.AsObject[] = [];

    public dataSource!: UserGrantsDataSource;
    public selection: SelectionModel<UserGrant.AsObject> = new SelectionModel<UserGrant.AsObject>(true, []);
    @ViewChild(MatPaginator) public paginator!: MatPaginator;
    @ViewChild(MatTable) public table!: MatTable<ProjectGrant.AsObject>;
    constructor(private userService: MgmtUserService) { }

    public displayedColumns: string[] = ['select', 'orgId', 'projectId', 'creationDate', 'changeDate', 'roleNamesList'];

    public ngOnInit(): void {
        this.dataSource = new UserGrantsDataSource(this.userService);
        this.dataSource.loadGrants(this.userId, 0, 25);
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
            this.userId,
            this.paginator.pageIndex,
            this.paginator.pageSize,
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
}
